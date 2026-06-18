//usr/bin/env sh -c 'd="${XDG_CACHE_HOME:-$HOME/.cache}/codex-consult"; b="$d/bin-$(uname -m)"; [ -x "$b" ] && [ "$b" -nt "$0" ] || ( mkdir -p "$d" && cd "$(dirname "$0")" && go build -o "$b" "$(basename "$0")" ) || exit 1; exec "$b" "$@"' "$0" "$@"; exit "$?"

// codex-consult — ask OpenAI Codex one question and print only its final answer.
//
// stdout carries ONLY the final answer, so it stays pipeable. The full transcript
// is written to a temp log whose path (and Codex's token usage) is reported on
// stderr on normal and error exits; an interrupt skips that banner, but the log is
// already on disk. The sandbox defaults to read-only because a consult should
// observe, not modify; raise it only to let Codex run tests (workspace-write) or
// apply changes (danger-full-access).
//
// The leading `//usr/bin/env sh -c ...` line makes this file self-executing: chmod
// +x and run it directly. It is a // comment to Go, but a shell command when execve
// finds no #! header — it compiles this source to a cached per-arch binary (only
// when missing or stale) and execs it. That preserves exact exit codes (unlike
// `go run`, which collapses every failure to 1) and is instant after the first
// build. Stdlib only, so the build needs no module or dependencies; building from
// the source's own directory keeps the caller's working module out of the picture.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// preamble frames every consult so callers pass only their task. Codex starts cold
// with no conversation history, so the reviewer role and the output discipline are
// fixed here once instead of being retyped into every prompt. It is deliberately
// generic across review and design questions — anything review-specific (e.g. "cite
// the line") would mis-shape a design consult.
const preamble = `You are an independent, cross-model reviewer and advisor, from a different model family than the agent consulting you (Claude). Your value is the outside view: catching blind spots a Claude-only analysis would share. You receive no conversation history — reason only from the task below and any repository you can read.

Be direct and decisive. Separate real defects from speculative risks, prefer concrete and minimal recommendations over sweeping rewrites, and if you find nothing material, say so plainly instead of inventing nits. Close with a clear verdict or recommendation.`

func main() { os.Exit(run()) }

func run() int {
	sandbox := "read-only"
	workdir := "."
	if wd, err := os.Getwd(); err == nil {
		workdir = wd
	}
	model := ""
	verbose := false

	// Flags precede the prompt. -C/-s/-m take a value; everything after the first
	// bare word (or `--`) is the prompt, so a question may contain spaces.
	argv := os.Args[1:]
	var promptArgs []string
	needVal := func(i int, flag string) (string, bool) {
		if i+1 >= len(argv) {
			fmt.Fprintf(os.Stderr, "codex-consult: %s needs a value\n", flag)
			return "", false
		}
		return argv[i+1], true
	}
	for i := 0; i < len(argv); i++ {
		a := argv[i]
		switch {
		case a == "-C" || a == "--repo":
			v, ok := needVal(i, "-C/--repo")
			if !ok {
				return 2
			}
			workdir, i = v, i+1
		case a == "-s" || a == "--sandbox":
			v, ok := needVal(i, "-s/--sandbox")
			if !ok {
				return 2
			}
			sandbox, i = v, i+1
		case a == "-m" || a == "--model":
			v, ok := needVal(i, "-m/--model")
			if !ok {
				return 2
			}
			model, i = v, i+1
		case a == "-v" || a == "--verbose":
			verbose = true
		case a == "-h" || a == "--help":
			usage(os.Stderr)
			return 0
		case a == "--":
			promptArgs = append(promptArgs, argv[i+1:]...)
			i = len(argv)
		case strings.HasPrefix(a, "-"):
			fmt.Fprintf(os.Stderr, "codex-consult: unknown option: %s\n", a)
			fmt.Fprintln(os.Stderr, "(if this is your question and it starts with '-', put it after '--')")
			usage(os.Stderr)
			return 2
		default:
			promptArgs = append(promptArgs, argv[i:]...)
			i = len(argv)
		}
	}

	// Combine the argument and piped stdin (argument leads), so `... "question"`,
	// `echo brief | ...`, and `git diff | ... "review"` all work. Read stdin only
	// when it is not a terminal, so an interactive run never blocks on the keyboard.
	argPrompt := strings.Join(promptArgs, " ")
	stdinPrompt := ""
	if !isTerminal(os.Stdin) {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "codex-consult: error reading stdin: %v\n", err)
			return 1
		}
		stdinPrompt = string(b)
	}
	var prompt string
	switch {
	case argPrompt != "" && strings.TrimSpace(stdinPrompt) != "":
		prompt = argPrompt + "\n\n" + stdinPrompt
	case argPrompt != "":
		prompt = argPrompt
	default:
		prompt = stdinPrompt
	}
	if strings.TrimSpace(prompt) == "" {
		fmt.Fprintln(os.Stderr, "codex-consult: empty prompt — pass a question as an argument or via stdin")
		return 2
	}

	// read-only stops writes, not reads: a $HOME working root lets Codex scan the
	// whole home tree. Warn rather than block, but make the exposure visible.
	if home, err := os.UserHomeDir(); err == nil && samePath(workdir, home) {
		fmt.Fprintln(os.Stderr, "codex-consult: warning — working root is $HOME; Codex can read your whole home tree. Pass -C <repo> to scope it.")
	}

	logf, err := os.CreateTemp("", "codex-consult-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "codex-consult: cannot create log file: %v\n", err)
		return 1
	}
	logPath := logf.Name()
	lastPath := logPath + ".answer"

	// Announce the log path on the way out — normal or error — so the caller can
	// always read Codex's full transcript. A signal terminates the process before
	// this defer runs, but the log file is already on disk, so only the banner is
	// lost; forwarding signals to the child to print it is not worth the complexity.
	defer announce(logPath, sandbox)

	cmdArgs := []string{
		"exec",
		"-C", workdir,
		"--sandbox", sandbox,
		"--color", "never",
		"--skip-git-repo-check",
		"--output-last-message", lastPath,
	}
	if model != "" {
		cmdArgs = append(cmdArgs, "--model", model)
	}
	cmdArgs = append(cmdArgs, "-") // read the prompt from stdin

	cmd := exec.Command("codex", cmdArgs...)
	cmd.Stdin = strings.NewReader(preamble + "\n\n---\n\n" + prompt)
	var sink io.Writer = logf
	if verbose {
		sink = io.MultiWriter(logf, os.Stderr)
	}
	cmd.Stdout, cmd.Stderr = sink, sink

	rc := 0
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			rc = ee.ExitCode()
		} else {
			fmt.Fprintf(os.Stderr, "codex-consult: failed to run codex: %v\n", err)
			rc = 127
		}
	}
	logf.Close()

	if answer, err := os.ReadFile(lastPath); err == nil && len(bytes.TrimSpace(answer)) > 0 {
		os.Stdout.Write(answer)
		if !bytes.HasSuffix(answer, []byte("\n")) {
			fmt.Println()
		}
	} else {
		fmt.Fprintf(os.Stderr, "codex-consult: no final message captured (exit %d) — showing log tail\n", rc)
		printTail(logPath, 30)
	}
	return rc
}

// announce prints the log path and Codex's token usage to stderr.
func announce(logPath, sandbox string) {
	tokens := "n/a"
	if data, err := os.ReadFile(logPath); err == nil {
		re := regexp.MustCompile(`(?i)tokens used\s*([0-9][0-9,]*)`)
		if ms := re.FindAllStringSubmatch(string(data), -1); len(ms) > 0 {
			tokens = ms[len(ms)-1][1]
		}
	}
	fmt.Fprintf(os.Stderr, "── codex (%s) · tokens: %s · full log: %s ──\n", sandbox, tokens, logPath)
}

func printTail(path string, n int) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	fmt.Fprintln(os.Stderr, strings.Join(lines, "\n"))
}

func isTerminal(f *os.File) bool {
	st, err := f.Stat()
	return err == nil && st.Mode()&os.ModeCharDevice != 0
}

func samePath(a, b string) bool {
	ra, ea := filepath.EvalSymlinks(a)
	rb, eb := filepath.EvalSymlinks(b)
	if ea != nil {
		ra, _ = filepath.Abs(a)
	}
	if eb != nil {
		rb, _ = filepath.Abs(b)
	}
	return ra == rb
}

func usage(w io.Writer) {
	fmt.Fprint(w, `codex-consult — ask OpenAI Codex one question; print only its final answer.

  codex-consult [options] "your question"
  echo "your brief" | codex-consult [options]
  git diff | codex-consult [options] "review this for races"   # argument + stdin combined

  -C, --repo <dir>      working root for Codex (default: current directory)
  -s, --sandbox <mode>  read-only (default) | workspace-write | danger-full-access
  -m, --model <model>   model override (default: your ~/.codex config)
  -v, --verbose         also stream Codex's full output to stderr
  -h, --help            this help

stdout = final answer only.   stderr = token usage + path to the full log.
`)
}
