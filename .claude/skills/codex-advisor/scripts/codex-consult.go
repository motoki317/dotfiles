//usr/bin/env sh -c 'd="${XDG_CACHE_HOME:-$HOME/.cache}/codex-consult"; b="$d/bin-$(uname -m)"; [ -x "$b" ] && [ "$b" -nt "$0" ] || ( mkdir -p "$d" && cd "$(dirname "$0")" && go build -o "$b" "$(basename "$0")" ) || exit 1; exec "$b" "$@"' "$0" "$@"; exit "$?"

// codex-consult — ask OpenAI Codex one question and print only its final answer.
//
// stdout carries ONLY the final answer, so it stays pipeable. Everything Codex emits is
// written verbatim to one log (a temp file, or -l <path>), whose path prints on stderr
// at startup. The wrapper passes `--json`, so Codex flushes one event per step instead
// of buffering the whole run; because the log is its raw stdout, those events land as
// they happen and a caller can tail the log to watch progress. The log is the raw record
// — no reformatting here — so it can't drift from what Codex actually emits. A run that
// exits non-zero with no `turn.completed` event failed (the reason is a `turn.failed`
// event in the log); a long silent gap with the process still alive may be stuck, though
// a pure-reasoning phase legitimately emits nothing between turn.started and the answer.
//
// The sandbox defaults to read-only because a consult should observe, not modify; raise
// it only to let Codex run tests (workspace-write) or apply changes (danger-full-access).
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
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// preamble frames every consult so callers pass only their task. Codex starts cold
// with no conversation history, so the reviewer role and the output discipline are
// fixed here once instead of being retyped into every prompt. It is deliberately
// generic across review and design questions — anything review-specific (e.g. "cite
// the line") would mis-shape a design consult.
const preamble = `You are an independent, cross-model reviewer and advisor, from a different model family than the agent consulting you (Claude). Your value is the outside view: catching blind spots a Claude-only analysis would share. You receive no conversation history — reason only from the task below and any repository you can read.

Be direct and decisive. Separate real defects from speculative risks, prefer concrete and minimal recommendations over sweeping rewrites, and if you find nothing material, say so plainly instead of inventing nits. Close with a clear verdict or recommendation.`

// stdinWait bounds how long we wait for the FIRST byte of piped input when the prompt
// was already supplied as an argument. It exists so a backgrounded run, which inherits
// an open pipe that never delivers data or EOF, cannot stall before Codex even starts.
// Only the first byte is time-bounded; once data is present the full read has no
// deadline, so a large or slow pipe is read in full rather than truncated.
const stdinWait = 2 * time.Second

func main() { os.Exit(run()) }

func run() int {
	sandbox := "read-only"
	workdir := "."
	if wd, err := os.Getwd(); err == nil {
		workdir = wd
	}
	model := ""
	logPath := ""
	verbose := false

	// Flags precede the prompt. -C/-s/-m/-l take a value; everything after the first
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
		case a == "-l" || a == "--log":
			v, ok := needVal(i, "-l/--log")
			if !ok {
				return 2
			}
			logPath, i = v, i+1
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
	stdinPrompt, ok := readStdin(argPrompt != "")
	if !ok {
		return 1
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

	// Open the log. Default to a temp file; honor -l so the caller can pick a path it
	// knows in advance and tail while Codex runs.
	var logf *os.File
	var err error
	if logPath == "" {
		logf, err = os.CreateTemp("", "codex-consult-*")
		if err == nil {
			logPath = logf.Name()
		}
	} else {
		logf, err = os.Create(logPath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "codex-consult: cannot create log file: %v\n", err)
		return 1
	}
	answerPath := logPath + ".answer"

	// Print the log path up front, not just at exit, so a caller that backgrounds this
	// run (and reads stderr incrementally) learns where to watch immediately.
	fmt.Fprintf(os.Stderr, "codex-consult: streaming events → %s\n", logPath)

	// Announce token usage and the log path on the way out — normal or error — so the
	// caller can always find the log. A signal terminates the process before this runs,
	// but the log is already on disk, so only the banner is lost.
	defer announce(logPath, sandbox)

	// --json makes Codex emit progress events as they happen; --output-last-message
	// still writes the final answer to its own file, so stdout stays clean.
	cmdArgs := []string{
		"exec",
		"--json",
		"-C", workdir,
		"--sandbox", sandbox,
		"--color", "never",
		"--skip-git-repo-check",
		"--output-last-message", answerPath,
	}
	if model != "" {
		cmdArgs = append(cmdArgs, "--model", model)
	}
	cmdArgs = append(cmdArgs, "-") // read the prompt from stdin

	cmd := exec.Command("codex", cmdArgs...)
	cmd.Stdin = strings.NewReader(preamble + "\n\n---\n\n" + prompt)
	// Codex's raw output is the log. -v also tees it to stderr for a live foreground view.
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

	if answer, err := os.ReadFile(answerPath); err == nil && len(bytes.TrimSpace(answer)) > 0 {
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

// readStdin returns the piped prompt, or "" when stdin is a terminal/character device
// (an interactive run or `</dev/null`). A goroutine blocks on the first byte so a
// backgrounded run with an idle inherited pipe (no data, no EOF) cannot stall here; only
// that first-byte wait is time-bounded when the prompt is already in argv (haveArg). Once
// a byte is seen the rest is read without a deadline, so a slow or large pipe such as
// `git diff | …` is read in full, never truncated. With no argument, stdin is the only
// prompt source, so it waits for the first byte however long that takes.
func readStdin(haveArg bool) (string, bool) {
	if isTerminal(os.Stdin) {
		return "", true
	}
	br := bufio.NewReader(os.Stdin)
	type result struct {
		s   string
		err error
	}
	hasData := make(chan bool, 1)
	full := make(chan result, 1)
	go func() {
		if _, err := br.Peek(1); err != nil {
			hasData <- false // empty pipe (EOF) or read error: nothing piped
			return
		}
		hasData <- true
		b, err := io.ReadAll(br) // includes the peeked byte; bufio retains it
		full <- result{string(b), err}
	}()
	var timeout <-chan time.Time
	if haveArg {
		timeout = time.After(stdinWait)
	}
	select {
	case ok := <-hasData:
		if !ok {
			return "", true
		}
		r := <-full
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "codex-consult: error reading stdin: %v\n", r.err)
			return "", false
		}
		return r.s, true
	case <-timeout:
		return "", true // argument prompt present; don't stall on an idle pipe
	}
}

// announce prints the log path and Codex's token usage to stderr. Tokens come from the
// last turn.completed event in the log; "n/a" when the run produced none (e.g. failed).
func announce(logPath, sandbox string) {
	fmt.Fprintf(os.Stderr, "── codex (%s) · tokens: %s · log: %s ──\n", sandbox, lastUsage(logPath), logPath)
}

// lastUsage scans the log backward for the most recent turn.completed event and reports
// its token counts. Reading one known field of one event type is the wrapper's only look
// inside Codex's JSON; the log itself stays a verbatim copy.
func lastUsage(logPath string) string {
	data, err := os.ReadFile(logPath)
	if err != nil {
		return "n/a"
	}
	lines := strings.Split(string(data), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if !strings.Contains(lines[i], `"turn.completed"`) {
			continue
		}
		var ev struct {
			Usage *struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
			} `json:"usage"`
		}
		if json.Unmarshal([]byte(lines[i]), &ev) == nil && ev.Usage != nil {
			return fmt.Sprintf("%d in / %d out", ev.Usage.InputTokens, ev.Usage.OutputTokens)
		}
	}
	return "n/a"
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
  -l, --log <path>      log file for Codex's raw JSONL events (default: a temp file).
                        Pick a path to tail progress while Codex runs.
  -v, --verbose         also tee Codex's output to stderr
  -h, --help            this help

stdout = final answer only.   stderr = log path (at start) + token usage (at end).
The log is Codex's raw --json event stream, written live — tail it to watch progress.
`)
}
