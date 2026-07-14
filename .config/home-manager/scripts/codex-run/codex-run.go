// codex-run — run OpenAI Codex for one non-interactive turn; only its final answer goes to
// stdout (so it pipes), while everything Codex emits streams verbatim to a log (temp file,
// or -l) you can tail. A non-zero exit with no `turn.completed` in the log means the run
// failed. Each run prints its Codex session id: `--resume <id>`/`--last` continues that session
// (retry a turn that died on a network error), and `codex resume <id>` opens it in the Codex TUI
// — but only after the run exits, since two writers on one session interleave its history. The
// mode fixes the role and default sandbox per call site: `advise` reviews
// without writing (read-only), `work` implements a delegated task — up to a whole
// requirements document — with full access (danger-full-access), committing locally as it
// goes. The house rules and skills are not injected here: codex auto-loads them from
// ~/.codex/AGENTS.md (git-tracked), which every run — including interactive ones — sees.
//
// Built by home-manager (buildGoModule) from ~/.config/home-manager/scripts and installed on
// PATH as `codex-run`; edits take effect on the next `home-manager switch`. Stdlib only.
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

// Preambles frame every run so callers pass only their task. Codex starts cold with no
// conversation history, so the role and the output discipline are fixed here once per mode
// instead of being retyped into every prompt. The advise preamble is deliberately generic
// across review and design questions — anything review-specific (e.g. "cite the line")
// would mis-shape a design consult. Each mode's cold variant is the default; the --context
// variant swaps only the one clause that would otherwise be false once --context has
// prepended the session transcript (see session-transcript).
const adviseRole = `You are an independent, cross-model reviewer and advisor, from a different model family than the agent consulting you (Claude). Your value is the outside view: catching blind spots a Claude-only analysis would share.`

const adviseCold = ` You receive no conversation history — reason only from the task below and any repository you can read.`

const adviseContext = ` The task below includes a redacted transcript of the consulting agent's current session, provided as context; reason from it, from the task, and from any repository you can read. It reflects one agent's framing — hold your independent outside view rather than simply ratifying it.`

const adviseDiscipline = `

Be direct and decisive. Separate real defects from speculative risks, prefer concrete and minimal recommendations over sweeping rewrites, and if you find nothing material, say so plainly instead of inventing nits. Close with a clear verdict or recommendation.`

const workRole = `You are an autonomous implementation agent from a different model family than the orchestrating agent (Claude), entrusted with a delegated implementation — typically an entire requirements document — inside a larger operating loop. Claude reviews the result after you finish; deliver work that survives that review.`

const workCold = ` You receive no conversation history — work from the task below, the repository you are in, and the house assets named in your AGENTS.md instructions.`

const workContext = ` The task below includes a redacted transcript of the orchestrating agent's current session, provided as background only — the delegated task itself is authoritative.`

const workDiscipline = `

Deliver the complete implementation: work non-interactively through every requirement in the task, tests included — prefer finishing over stopping to ask. Resolve routine implementation details from the requirements, the repository, and the house assets; make only narrow, reversible assumptions and never invent product policy. Stop and report instead when missing information would change externally visible behavior, when requirements contradict each other, when credentials or external services are missing, or when an irreversible or destructive action would be needed — after finishing the unblocked requirements first.

Commit locally as you go: after each coherent, independently reviewable unit, run the relevant checks and commit that green unit before starting the next — do not defer all commits to the end, and do not commit known failures. Stage only what you changed; never reset, checkout, stash, or discard work that is not yours. Never run git push or change remote configuration, and never create or update pull requests, releases, deployments, packages, or any other external service state, even if repository instructions ask for it; read-only network access and dependency downloads are fine. Treat the selected repository as the writable project scope and do not modify user files outside it (temporary files and dependency caches excepted).

Before declaring completion, re-read the authoritative requirements and audit every requirement against the implementation, tests, and commits. Begin the final response with exactly STATUS: COMPLETE, STATUS: PARTIAL, or STATUS: BLOCKED — COMPLETE only when every requirement is implemented and verification passed. Then report: requirement-by-requirement status, what was implemented, commit hashes and subjects, exact verification commands and their results, assumptions, deviations, blockers, and the house assets you read.`

// defaultContextPrompt is the review request used when `advise --context` is given with no
// brief: the low-friction "advise me" call that mirrors a zero-arg advisor, asking Codex to
// review the session trajectory rather than a specific question. Any piped brief replaces
// it. work has no equivalent — an implementer with no task has nothing to do.
const defaultContextPrompt = `Review the session transcript above — the task, the approach taken, and the work done so far. Surface blind spots, unstated assumptions, risks, and anything wrong, missing, or worth reconsidering before I proceed. Be direct and specific, and end with a clear verdict on whether the current direction is sound.`

// defaultResumePrompt is sent when --resume/--last is given with no piped follow-up: the
// network-failure retry, where the goal is simply to finish the interrupted turn. The role and
// discipline are already in the resumed session's history, so only this nudge is needed. Any
// piped stdin replaces it (an advise follow-up question, or a corrective instruction to work).
const defaultResumePrompt = `The previous turn was interrupted before it finished. Continue from where you left off and complete it, then give your final response.`

func main() { os.Exit(run()) }

func run() int {
	// The mode comes first so each call site's intent — and with it the default sandbox —
	// is explicit; an omitted mode can never silently select write access.
	argv := os.Args[1:]
	if len(argv) == 0 {
		usage(os.Stderr)
		return 2
	}
	if argv[0] == "-h" || argv[0] == "--help" {
		usage(os.Stderr)
		return 0
	}
	mode := argv[0]
	if mode != "advise" && mode != "work" {
		fmt.Fprintf(os.Stderr, "codex-run: unknown mode %q — expected advise or work\n", mode)
		usage(os.Stderr)
		return 2
	}
	argv = argv[1:]

	// work gets full access by default: workspace-write leaves .git read-only (commits fail
	// on index.lock) and blocks the network (dependency fetches fail), so a long autonomous
	// run needs the unsandboxed mode. -C still tells Codex its workspace; -s overrides.
	sandbox := "read-only"
	if mode == "work" {
		sandbox = "danger-full-access"
	}
	workdir := "."
	workdirSet := false
	if wd, err := os.Getwd(); err == nil {
		workdir = wd
	}
	model := ""
	logPath := ""
	verbose := false
	withContext := false
	resumeID := ""
	resumeLast := false

	// After the mode: flags plus a stdin prompt — there is no positional prompt, so a
	// brief's shell metacharacters can never be interpreted. -C/-s/-m/-l take a value.
	needVal := func(i int, flag string) (string, bool) {
		if i+1 >= len(argv) {
			fmt.Fprintf(os.Stderr, "codex-run: %s needs a value\n", flag)
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
			workdir, workdirSet, i = v, true, i+1
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
		case a == "--resume":
			v, ok := needVal(i, "--resume")
			if !ok {
				return 2
			}
			resumeID, i = v, i+1
		case a == "--last":
			resumeLast = true
		case a == "-h" || a == "--help":
			usage(os.Stderr)
			return 0
		case strings.HasPrefix(a, "-"):
			fmt.Fprintf(os.Stderr, "codex-run: unknown option: %s\n", a)
			usage(os.Stderr)
			return 2
		default:
			fmt.Fprintf(os.Stderr, "codex-run: unexpected argument %q — the prompt is read from stdin\n", a)
			usage(os.Stderr)
			return 2
		}
	}

	// Normalize the resume selectors. --resume takes a session id; --last picks the most
	// recent. Fold `--resume --last` (the value lands as the string "--last") into --last so
	// either spelling works, then reject the ambiguous or malformed combinations.
	if resumeID == "--last" {
		resumeID, resumeLast = "", true
	}
	resuming := resumeID != "" || resumeLast
	if resumeID != "" && resumeLast {
		fmt.Fprintln(os.Stderr, "codex-run: pass either --resume <id> or --last, not both")
		return 2
	}
	if strings.HasPrefix(resumeID, "-") {
		fmt.Fprintf(os.Stderr, "codex-run: --resume needs a session id, got %q\n", resumeID)
		return 2
	}
	if resuming && withContext {
		fmt.Fprintln(os.Stderr, "codex-run: --context cannot be combined with --resume/--last — the resumed session already holds its history")
		return 2
	}

	// With full access, an implicit working root is how accidents happen: work must name
	// its repository rather than inherit whatever directory the caller happened to be in.
	// A resume inherits the session's recorded cwd, so it needs no -C (workdir only scopes
	// which session --last selects and gives the run its workspace, both via cmd.Dir below).
	if mode == "work" && !workdirSet && !resuming {
		fmt.Fprintln(os.Stderr, "codex-run: work requires an explicit -C <repo>")
		return 2
	}

	// The prompt is read only from stdin, so a brief's backticks, `$(…)`, `$VAR`, and `!`
	// are never shell-interpreted. --context prepends the session transcript to it below.
	prompt, ok := readStdin()
	if !ok {
		return 1
	}
	prompt = strings.TrimSpace(prompt)

	// work always needs a brief — --context supplies background, not a task. A cold advise
	// needs one too: with no session and no question, Codex has nothing to review. Under
	// `advise --context` an empty brief is instead the low-friction "advise me" call: fall
	// back to a fixed session-trajectory prompt (assigned below).
	if !resuming && prompt == "" && (mode == "work" || !withContext) {
		if mode == "work" {
			fmt.Fprintln(os.Stderr, "codex-run: empty prompt — pipe the delegated task via stdin (e.g. `codex-run work … < task.md`)")
		} else {
			fmt.Fprintln(os.Stderr, "codex-run: empty prompt — pipe a brief via stdin (e.g. `codex-run advise … < brief.md`), or pass --context for a briefless session review")
		}
		return 2
	}

	// --context prepends the current session's transcript as context. We run the
	// session-transcript helper rather than parsing the JSONL here, so this shim stays
	// unaware of the transcript format; a failed extraction aborts the run, so Codex is
	// never told a transcript is present when it is not.
	if withContext {
		transcript, err := sessionTranscript()
		if err != nil {
			fmt.Fprintf(os.Stderr, "codex-run: --context: %v\n", err)
			return 1
		}
		if prompt == "" {
			prompt = defaultContextPrompt
		}
		divider := "\n\n===== END SESSION TRANSCRIPT · REVIEW REQUEST BELOW =====\n\n"
		if mode == "work" {
			divider = "\n\n===== END SESSION TRANSCRIPT · DELEGATED TASK BELOW =====\n\n"
		}
		prompt = transcript + divider + prompt
	}

	// A resume with no piped follow-up is the network-failure retry: nudge the session to
	// finish where it stopped. (--context is incompatible with resume, so this cannot collide
	// with the transcript prompt above.)
	if resuming && prompt == "" {
		prompt = defaultResumePrompt
	}

	// read-only stops writes, not reads: a working root at (or above) $HOME lets Codex scan
	// the whole home tree — make the exposure visible, but proceed. A writable sandbox there
	// would make everything you own project scope, so that is an error, not a warning — even
	// for a repo rooted at $HOME, which must be delegated via a git worktree instead. A resume
	// inherits the session's own recorded sandbox and cwd — workdir only scopes --last here and
	// cannot broaden that scope — so the guard does not apply to it.
	if !resuming {
		if home, err := os.UserHomeDir(); err == nil && coversHome(workdir, home) {
			if sandbox != "read-only" {
				fmt.Fprintln(os.Stderr, "codex-run: working root covers $HOME with a writable sandbox — pass -C <repo>; for a repo rooted at $HOME, pass a git worktree of it")
				return 2
			}
			fmt.Fprintln(os.Stderr, "codex-run: warning — working root covers $HOME; Codex can read your whole home tree. Pass -C <repo> to scope it.")
		}
	}

	// Open the log. Default to a temp file; honor -l so the caller can pick a path it
	// knows in advance and tail while Codex runs.
	var logf *os.File
	var err error
	if logPath == "" {
		logf, err = os.CreateTemp("", "codex-run-*")
		if err == nil {
			logPath = logf.Name()
		}
	} else {
		logf, err = os.Create(logPath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "codex-run: cannot create log file: %v\n", err)
		return 1
	}
	// Codex writes its final message here via --output-last-message; we read it back for
	// stdout. Use a unique temp file, never <log>.answer: runs that share a --log path — reused
	// sequentially or run concurrently — would otherwise share one answer file and print each
	// other's (or a previous run's) verdict. Removed on exit; the verdict itself goes to stdout.
	answerFile, err := os.CreateTemp("", "codex-run-answer-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "codex-run: cannot create answer file: %v\n", err)
		return 1
	}
	answerPath := answerFile.Name()
	answerFile.Close()
	defer os.Remove(answerPath)

	// Print the log path up front, not just at exit, so a caller that backgrounds this
	// run (and reads stderr incrementally) learns where to watch immediately.
	fmt.Fprintf(os.Stderr, "codex-run: streaming events → %s\n", logPath)

	// Close with the token-usage line on the way out — normal or error. The session id and
	// resume hints are printed earlier instead (sessionAnnouncer, below), the instant Codex
	// emits them: a signal (Ctrl-C, hangup) kills the process before this deferred banner runs,
	// so an id shown only here would be lost exactly when a resume is most needed. `printed`
	// tells announce the id is already out; it reprints only as a backstop. A resumed run's real
	// sandbox is the session's own, not this mode's default, so label it "resumed".
	sandboxLabel := sandbox
	if resuming {
		sandboxLabel = "resumed"
	}
	printed := false
	defer announce(logPath, mode, sandboxLabel, &printed)

	// --json makes Codex emit progress events as they happen; --output-last-message
	// still writes the final answer to its own file, so stdout stays clean. A resume reuses
	// this whole pipeline (log, answer capture, token banner, STATUS gate) — it only swaps in
	// the `exec resume` subcommand and drops -C/--sandbox, which resume rejects because the
	// session already recorded both; cwd comes from cmd.Dir below instead.
	var cmdArgs []string
	if resuming {
		cmdArgs = []string{"exec", "resume"}
		if resumeLast {
			cmdArgs = append(cmdArgs, "--last")
		} else {
			cmdArgs = append(cmdArgs, resumeID)
		}
		cmdArgs = append(cmdArgs,
			"--config", "model_context_window=272000",
			"--config", "model_auto_compact_token_limit=240000",
			"--json",
			"--skip-git-repo-check",
			"--output-last-message", answerPath,
		)
	} else {
		cmdArgs = []string{
			"exec",
			"--config", "model_context_window=272000",
			"--config", "model_auto_compact_token_limit=240000",
			"--json",
			"-C", workdir,
			"--sandbox", sandbox,
			"--color", "never",
			"--skip-git-repo-check",
			"--output-last-message", answerPath,
		}
	}
	if model != "" {
		cmdArgs = append(cmdArgs, "--model", model)
	}
	cmdArgs = append(cmdArgs, "-") // read the prompt from stdin

	cmd := exec.Command("codex", cmdArgs...)
	if resuming {
		// The session already holds the preamble in its history; re-injecting the role would
		// duplicate it. Send only the follow-up prompt. cmd.Dir gives resume its workspace and
		// scopes which recorded session --last selects.
		cmd.Stdin = strings.NewReader(prompt)
		cmd.Dir = workdir
	} else {
		role, cold, context, discipline := adviseRole, adviseCold, adviseContext, adviseDiscipline
		if mode == "work" {
			role, cold, context, discipline = workRole, workCold, workContext, workDiscipline
		}
		preamble := role + cold + discipline
		if withContext {
			preamble = role + context + discipline
		}
		cmd.Stdin = strings.NewReader(preamble + "\n\n---\n\n" + prompt)
	}
	// Codex's raw output is the log. -v also tees it to stderr for a live foreground view.
	var sink io.Writer = logf
	if verbose {
		sink = io.MultiWriter(logf, os.Stderr)
	}
	// Surface the session id the moment Codex's first thread.started event flows past — early,
	// so a run later killed by a signal has already shown how to resume it — while forwarding
	// every byte to the log unchanged.
	ann := &sessionAnnouncer{sink: sink, out: os.Stderr, mode: mode, printed: &printed}
	cmd.Stdout, cmd.Stderr = ann, ann

	rc := 0
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			rc = ee.ExitCode()
		} else {
			fmt.Fprintf(os.Stderr, "codex-run: failed to run codex: %v\n", err)
			rc = 127
		}
	}
	logf.Close()

	// Codex's exit code is trusted but verified: a run whose log holds no turn.completed
	// failed even if the process exited 0 (e.g. a stream dropped mid-turn).
	if rc == 0 && !turnCompleted(logPath) {
		fmt.Fprintln(os.Stderr, "codex-run: codex exited 0 but the log has no turn.completed — treating the run as failed")
		rc = 1
	}

	answer, aerr := os.ReadFile(answerPath)
	if aerr == nil && len(bytes.TrimSpace(answer)) > 0 {
		os.Stdout.Write(answer)
		if !bytes.HasSuffix(answer, []byte("\n")) {
			fmt.Println()
		}
	} else {
		fmt.Fprintf(os.Stderr, "codex-run: no final message captured (exit %d) — showing log tail\n", rc)
		printTail(logPath, 30)
		if rc == 0 {
			// Codex exited 0 but wrote no verdict — never report success with no answer.
			rc = 1
		}
	}

	// A writable sandbox is not transactional: a turn that dies mid-task leaves whatever
	// it had already edited. The caller must see that before trusting the tree or retrying.
	// Checked after every rc adjustment above, so no failed writable run skips the warning.
	if rc != 0 && sandbox != "read-only" {
		fmt.Fprintf(os.Stderr, "codex-run: turn failed (exit %d) — partial edits may remain; inspect git status/diff\n", rc)
	}

	// Transport success is not semantic completion: the work preamble demands a STATUS
	// marker, and a clean turn that reports PARTIAL or BLOCKED — or omits the marker —
	// must not exit 0, or the orchestrator would take an unfinished implementation as
	// done. Exit 3 keeps it distinct from a failed turn.
	if mode == "work" && rc == 0 {
		if st := answerStatus(answer); !strings.HasPrefix(st, "COMPLETE") {
			fmt.Fprintf(os.Stderr, "codex-run: work run's STATUS marker is %q, not COMPLETE — implementation incomplete; read the report and the commits made so far (exit 3)\n", st)
			rc = 3
		}
	}
	return rc
}

// answerStatus extracts the STATUS marker the work preamble requires as the first line of
// the final response, or "" when absent. Markdown emphasis is stripped rather than trimmed:
// `**STATUS**: COMPLETE` puts the emphasis mid-line, where a Trim can't reach it. The caller
// matches COMPLETE as a prefix, so a decorated-but-complete marker ("COMPLETE — all done")
// still counts; PARTIAL, BLOCKED, and a missing marker stay fail-safe.
func answerStatus(answer []byte) string {
	first, _, _ := strings.Cut(strings.TrimSpace(string(answer)), "\n")
	first = strings.Map(func(r rune) rune {
		if r == '*' || r == '_' || r == '`' || r == '#' {
			return -1
		}
		return r
	}, first)
	if s, ok := strings.CutPrefix(strings.TrimSpace(first), "STATUS:"); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

// turnCompleted reports whether the log holds a turn.completed event — the cross-check that
// keeps an exit-0 process whose turn never finished from being reported as success.
func turnCompleted(logPath string) bool {
	data, err := os.ReadFile(logPath)
	return err == nil && bytes.Contains(data, []byte(`"turn.completed"`))
}

// readStdin returns the whole piped prompt, or "" when stdin is a terminal (an interactive
// run with nothing piped). The prompt comes only from stdin — never a positional argument —
// so backticks, `$(…)`, `$VAR`, and `!` in a code-heavy brief cannot be shell-interpreted. A
// slow or large pipe (e.g. `cat brief.md <(git diff)`) is read in full, never truncated; an
// empty pipe yields "" — work and a cold advise then error, while `advise --context` treats
// it as the briefless "advise me" call.
func readStdin() (string, bool) {
	if isTerminal(os.Stdin) {
		return "", true
	}
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "codex-run: error reading stdin: %v\n", err)
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

// announce prints the closing token-usage line — the one fact known only at completion; tokens
// come from the last turn.completed event, "n/a" when the run produced none (e.g. failed). The
// session id and resume hints are surfaced earlier, the instant Codex emits them (see
// sessionAnnouncer), so a signal-killed run still shows how to resume. announce reprints them
// here only as a backstop: the rare run whose live stream-watch missed the id but whose completed
// log still holds it. The log path is not repeated — it was printed at start.
func announce(logPath, mode, sandbox string, printed *bool) {
	fmt.Fprintf(os.Stderr, "── codex %s (%s) · tokens: %s ──\n", mode, sandbox, lastUsage(logPath))
	if !*printed {
		if tid := threadID(logPath); tid != "" {
			printSession(os.Stderr, mode, tid)
		}
	}
}

// printSession writes the session id and how to resume it. The headless retry (`codex-run
// --resume`) leads, since recovering a dropped run is the common case; `codex resume <id>` opens
// the same session in the Codex TUI — both only after this run exits, never against a session
// another run is still writing. Takes a writer so the stream watcher (stderr) and tests share it.
func printSession(w io.Writer, mode, tid string) {
	fmt.Fprintf(w, "   session: %s\n", tid)
	fmt.Fprintf(w, "   resume:  codex-run %s --resume %s  (headless retry)  ·  codex resume %s  (TUI, after this run exits)\n", mode, tid, tid)
}

// sessionAnnouncer wraps the log sink and prints the session id and resume hints (to out) the
// moment Codex's first thread.started event flows past — early, so a run later killed by a signal
// (Ctrl-C, hangup), which never reaches the deferred banner, has still shown how to resume it.
// thread.started is Codex's opening event, so it lands in the first bytes; once seen — or 64 KiB
// pass without it, meaning the run created no session — scanning stops. Every byte is forwarded to
// the sink unchanged, so the log stays a verbatim copy. *printed lets the deferred announce know
// the id is already out. Its Write runs in os/exec's copier goroutine; cmd.Run joins that goroutine
// before returning, so announce's later read of *printed sees this write (no race).
type sessionAnnouncer struct {
	sink    io.Writer // the log (and stderr under -v); receives every byte verbatim
	out     io.Writer // where the session banner goes (os.Stderr in prod)
	mode    string
	printed *bool
	buf     []byte // accumulates the opening bytes until thread.started is found
	giveUp  bool   // set once found or once 64 KiB passed without it
}

func (s *sessionAnnouncer) Write(p []byte) (int, error) {
	n, err := s.sink.Write(p)
	if !s.giveUp && !*s.printed {
		s.buf = append(s.buf, p[:n]...)
		if tid := scanThreadID(string(s.buf)); tid != "" {
			printSession(s.out, s.mode, tid)
			*s.printed, s.giveUp, s.buf = true, true, nil
		} else if len(s.buf) > 64<<10 {
			s.giveUp, s.buf = true, nil
		}
	}
	return n, err
}

// threadID returns the session/thread UUID from the log's first thread.started event — the id
// codex resume (TUI) and codex exec resume (headless) both accept, and the name of the durable
// rollout under ~/.codex/sessions. "" when the run died before the session was created.
func threadID(logPath string) string {
	data, err := os.ReadFile(logPath)
	if err != nil {
		return ""
	}
	return scanThreadID(string(data))
}

// scanThreadID pulls the thread_id from the first complete thread.started line in data — the
// shared parser for both the finished log (threadID) and the live stream (sessionAnnouncer).
// A line still mid-write fails to unmarshal and is skipped, so a straddled event is picked up
// on the next chunk. Mirrors lastUsage: one known field of one event type, no format ownership.
func scanThreadID(data string) string {
	for _, line := range strings.Split(data, "\n") {
		if !strings.Contains(line, `"thread.started"`) {
			continue
		}
		var ev struct {
			ThreadID string `json:"thread_id"`
		}
		if json.Unmarshal([]byte(line), &ev) == nil && ev.ThreadID != "" {
			return ev.ThreadID
		}
	}
	return ""
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

// coversHome reports whether root is $HOME or an ancestor of it (/, /home, …) — either way
// a sandbox rooted there reaches everything the user owns.
func coversHome(root, home string) bool {
	r, h := resolvePath(root), resolvePath(home)
	return r == h || strings.HasPrefix(h, strings.TrimSuffix(r, "/")+"/")
}

// resolvePath makes a path canonical for comparison. Abs before EvalSymlinks: a relative
// path (e.g. `-C .`) resolves to a relative result, which would never match the absolute
// home path — silently skipping the $HOME guard.
func resolvePath(p string) string {
	if a, err := filepath.Abs(p); err == nil {
		p = a
	}
	if r, err := filepath.EvalSymlinks(p); err == nil {
		return r
	}
	return p
}

func usage(w io.Writer) {
	fmt.Fprint(w, `codex-run — run OpenAI Codex for one non-interactive turn; print only its final answer.

  codex-run advise --context                                    # briefless: review the current session ("advise me")
  codex-run advise --context [options] < brief.md               # session + a specific question
  codex-run advise [options] < brief.md                         # cold review; brief required, read from stdin
  cat brief.md <(git diff main) | codex-run advise [options]    # question + diff
  codex-run work -C <repo> [options] < task.md                  # implement a delegated task (-C and task required)
  codex-run work --resume <id> -C <repo>                        # continue a session that died mid-turn (retry)
  codex-run <mode> --last [options] [< followup.md]             # continue the most recent session for this repo

The mode fixes Codex's role and default sandbox per call:
  advise  independent reviewer/advisor; read-only sandbox
  work    autonomous implementer for anything up to a whole requirements document;
          danger-full-access sandbox so it can commit locally (never push) and fetch
          deps — pass -s workspace-write for containment (blocks commits and network).
          An auto-appended preamble points it at ~/.claude rules and skills and demands
          a STATUS marker: exit 0 = STATUS: COMPLETE, exit 3 = the turn finished but
          the report's STATUS is not COMPLETE (PARTIAL, BLOCKED, missing, malformed).
          A repo rooted at $HOME must be delegated via a git worktree of it.

  -C, --repo <dir>      working root for Codex (default: current directory)
  -s, --sandbox <mode>  read-only | workspace-write | danger-full-access (default: per mode)
  -m, --model <model>   model override (default: your ~/.codex config)
  -l, --log <path>      log file for Codex's raw JSONL events (default: a temp file).
                        Pick a path to tail progress while Codex runs.
  -v, --verbose         also tee Codex's output to stderr
  -x, --context         prepend the current session's transcript as context (runs
                        session-transcript; aborts if extraction fails). Under advise,
                        makes the brief optional — with none, Codex reviews the session.
      --resume <id>     continue an existing session (its session id is printed in an earlier
                        run's banner) instead of starting fresh. Sandbox and cwd inherit from that
                        session; a piped stdin prompt is sent as the follow-up, else the
                        session is nudged to finish where it stopped. Skips the preamble
                        (already in history) and -C/-s (rejected by resume). work keeps its
                        STATUS gate; -C still scopes the run's workspace.
      --last            like --resume, but for the most recent session in the -C repo (no id).
  -h, --help            this help

stdout = final answer only.   stderr = log path + session id + resume hints (at start); token usage (at end).
The log is Codex's raw --json event stream, written live — tail it to watch progress.
Each run prints its session id: retry a run that died with --resume/--last, or open it in the
Codex TUI with "codex resume <id>" — but only after it exits (two writers on one session
interleave its history). A failed work turn can leave partial edits — inspect git status first.
`)
}
