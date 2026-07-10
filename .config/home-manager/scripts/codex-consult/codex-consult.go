// codex-consult — ask OpenAI Codex one question; only its final answer goes to stdout (so
// it pipes), while everything Codex emits streams verbatim to a log (temp file, or -l) you
// can tail. A non-zero exit with no `turn.completed` in the log means the run failed. The
// sandbox defaults to read-only; raise it to run tests (workspace-write) or edit
// (danger-full-access).
//
// Built by home-manager (buildGoModule) from ~/.config/home-manager/scripts and installed on
// PATH as `codex-consult`; edits take effect on the next `home-manager switch`. Stdlib only.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// preamble frames every consult so callers pass only their task. Codex starts cold
// with no conversation history, so the reviewer role and the output discipline are
// fixed here once instead of being retyped into every prompt. It is deliberately
// generic across review and design questions — anything review-specific (e.g. "cite
// the line") would mis-shape a design consult. The cold variant is the default; the
// --context variant swaps only the one clause that would otherwise be false once
// --context has prepended the session transcript (see session-transcript).
const preambleRole = `You are an independent, cross-model reviewer and advisor, from a different model family than the agent consulting you (Claude). Your value is the outside view: catching blind spots a Claude-only analysis would share.`

const preambleCold = ` You receive no conversation history — reason only from the task below and any repository you can read.`

const preambleContext = ` The task below includes a redacted transcript of the consulting agent's current session, provided as context; reason from it, from the task, and from any repository you can read. It reflects one agent's framing — hold your independent outside view rather than simply ratifying it.`

const preambleDiscipline = `

Be direct and decisive. Separate real defects from speculative risks, prefer concrete and minimal recommendations over sweeping rewrites, and if you find nothing material, say so plainly instead of inventing nits. Close with a clear verdict or recommendation.`

// defaultContextPrompt is the review request used when --context is given with no brief: the
// low-friction "advise me" call that mirrors a zero-arg advisor, asking Codex to review the
// session trajectory rather than a specific question. Any piped brief replaces it.
const defaultContextPrompt = `Review the session transcript above — the task, the approach taken, and the work done so far. Surface blind spots, unstated assumptions, risks, and anything wrong, missing, or worth reconsidering before I proceed. Be direct and specific, and end with a clear verdict on whether the current direction is sound.`

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
	withContext := false

	// Input is flags plus a stdin prompt — there is no positional prompt, so a brief's
	// shell metacharacters can never be interpreted. -C/-s/-m/-l take a value.
	argv := os.Args[1:]
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
		case a == "-x" || a == "--context":
			withContext = true
		case a == "-h" || a == "--help":
			usage(os.Stderr)
			return 0
		case strings.HasPrefix(a, "-"):
			fmt.Fprintf(os.Stderr, "codex-consult: unknown option: %s\n", a)
			usage(os.Stderr)
			return 2
		default:
			fmt.Fprintf(os.Stderr, "codex-consult: unexpected argument %q — the prompt is read from stdin\n", a)
			usage(os.Stderr)
			return 2
		}
	}

	// The prompt is read only from stdin, so a brief's backticks, `$(…)`, `$VAR`, and `!`
	// are never shell-interpreted. --context prepends the session transcript to it below.
	prompt, ok := readStdin()
	if !ok {
		return 1
	}
	prompt = strings.TrimSpace(prompt)

	// A cold review needs a brief — with no session and no question, Codex has nothing to
	// review. Under --context an empty brief is instead the low-friction "advise me" call:
	// fall back to a fixed session-trajectory prompt (assigned below).
	if prompt == "" && !withContext {
		fmt.Fprintln(os.Stderr, "codex-consult: empty prompt — pipe a brief via stdin (e.g. `codex-consult … < brief.md`), or pass --context for a briefless session review")
		return 2
	}

	// --context prepends the current session's transcript as context. We run the
	// session-transcript helper rather than parsing the JSONL here, so this shim stays
	// unaware of the transcript format; a failed extraction aborts the consult, so Codex is
	// never told a transcript is present when it is not.
	if withContext {
		transcript, err := sessionTranscript()
		if err != nil {
			fmt.Fprintf(os.Stderr, "codex-consult: --context: %v\n", err)
			return 1
		}
		if prompt == "" {
			prompt = defaultContextPrompt
		}
		prompt = transcript + "\n\n===== END SESSION TRANSCRIPT · REVIEW REQUEST BELOW =====\n\n" + prompt
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
	// Codex writes its final message here via --output-last-message; we read it back for
	// stdout. Use a unique temp file, never <log>.answer: runs that share a --log path — reused
	// sequentially or run concurrently — would otherwise share one answer file and print each
	// other's (or a previous run's) verdict. Removed on exit; the verdict itself goes to stdout.
	answerFile, err := os.CreateTemp("", "codex-consult-answer-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "codex-consult: cannot create answer file: %v\n", err)
		return 1
	}
	answerPath := answerFile.Name()
	answerFile.Close()
	defer os.Remove(answerPath)

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

	preamble := preambleRole + preambleCold + preambleDiscipline
	if withContext {
		preamble = preambleRole + preambleContext + preambleDiscipline
	}
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
		if rc == 0 {
			// Codex exited 0 but wrote no verdict — never report success with no answer.
			rc = 1
		}
	}
	return rc
}

// readStdin returns the whole piped prompt, or "" when stdin is a terminal (an interactive
// run with nothing piped). The prompt comes only from stdin — never a positional argument —
// so backticks, `$(…)`, `$VAR`, and `!` in a code-heavy brief cannot be shell-interpreted. A
// slow or large pipe (e.g. `cat brief.md <(git diff)`) is read in full, never truncated; an
// empty pipe yields "" — a cold review then errors, while --context treats it as the
// briefless "advise me" call.
func readStdin() (string, bool) {
	if isTerminal(os.Stdin) {
		return "", true
	}
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "codex-consult: error reading stdin: %v\n", err)
		return "", false
	}
	return string(b), true
}

// sessionTranscript runs the session-transcript helper and returns its stdout (the redacted
// current-session transcript). Its own stderr (progress or error) passes straight through. A
// non-zero exit or empty output is an error, so --context never proceeds with no context.
func sessionTranscript() (string, error) {
	cmd := exec.Command("session-transcript")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("session-transcript failed (%v) — is it on PATH and are you in a Claude Code session?", err)
	}
	if strings.TrimSpace(out.String()) == "" {
		return "", fmt.Errorf("session-transcript produced no transcript")
	}
	return out.String(), nil
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

  codex-consult --context                                       # briefless: review the current session ("advise me")
  codex-consult --context [options] < brief.md                  # session + a specific question
  codex-consult [options] < brief.md                            # cold review; brief required, read from stdin
  cat brief.md <(git diff main) | codex-consult [options]       # question + diff

  -C, --repo <dir>      working root for Codex (default: current directory)
  -s, --sandbox <mode>  read-only (default) | workspace-write | danger-full-access
  -m, --model <model>   model override (default: your ~/.codex config)
  -l, --log <path>      log file for Codex's raw JSONL events (default: a temp file).
                        Pick a path to tail progress while Codex runs.
  -v, --verbose         also tee Codex's output to stderr
  -x, --context         prepend the current session's transcript as context (runs
                        session-transcript; aborts if extraction fails). Makes the
                        brief optional — with none, Codex reviews the session.
  -h, --help            this help

stdout = final answer only.   stderr = log path (at start) + token usage (at end).
The log is Codex's raw --json event stream, written live — tail it to watch progress.
`)
}
