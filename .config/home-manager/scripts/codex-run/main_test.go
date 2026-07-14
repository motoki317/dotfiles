package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunUsesStatelessContextLimitOverrides(t *testing.T) {
	tmp := t.TempDir()
	argsPath := filepath.Join(tmp, "args")
	logPath := filepath.Join(tmp, "events.jsonl")
	fakeCodex := filepath.Join(tmp, "codex")
	script := `#!/bin/sh
previous=
for argument do
	printf '%s\n' "$argument" >> "$CODEX_ARGS_FILE"
	if [ "$previous" = "--output-last-message" ]; then
		printf 'review complete\n' > "$argument"
	fi
	previous=$argument
done
printf '%s\n' '{"type":"turn.completed","usage":{"input_tokens":1,"output_tokens":1}}'
`
	if err := os.WriteFile(fakeCodex, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CODEX_ARGS_FILE", argsPath)
	t.Setenv("PATH", tmp+string(os.PathListSeparator)+os.Getenv("PATH"))

	prompt, err := os.CreateTemp(tmp, "prompt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := prompt.WriteString("review this"); err != nil {
		t.Fatal(err)
	}
	if _, err := prompt.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	previousArgs, previousStdin := os.Args, os.Stdin
	os.Args, os.Stdin = []string{"codex-run", "advise", "-l", logPath}, prompt
	t.Cleanup(func() {
		os.Args, os.Stdin = previousArgs, previousStdin
		prompt.Close()
	})

	if code := run(); code != 0 {
		t.Fatalf("run() = %d, want 0", code)
	}
	got, err := os.ReadFile(argsPath)
	if err != nil {
		t.Fatal(err)
	}
	arguments := "\n" + string(got)
	want := "\n--config\nmodel_context_window=272000\n--config\nmodel_auto_compact_token_limit=240000\n"
	if !strings.Contains(arguments, want) || strings.Contains(arguments, "\n--profile\n") {
		t.Fatalf("codex arguments:\n%s\nwant context-limit overrides without --profile", got)
	}
}

// TestSessionAnnouncerPrintsEarlyAndForwardsVerbatim pins the crash-robustness fix: the watcher
// must surface the session id the moment thread.started flows past — even when that event is split
// across two writes — while forwarding every byte to the log unchanged. A corrupted log or a
// missed id would defeat the resume feature.
func TestSessionAnnouncerPrintsEarlyAndForwardsVerbatim(t *testing.T) {
	var log, banner strings.Builder
	printed := false
	ann := &sessionAnnouncer{sink: &log, out: &banner, mode: "advise", printed: &printed}

	stream := `{"type":"thread.started","thread_id":"sess-xyz"}` + "\n" +
		`{"type":"turn.completed","usage":{"input_tokens":1,"output_tokens":1}}` + "\n"
	// Split mid thread.started line: the first chunk carries the type marker but not yet the id,
	// so it must NOT print — the id only becomes parseable once the second chunk completes the line.
	cut := 30
	if _, err := ann.Write([]byte(stream[:cut])); err != nil {
		t.Fatal(err)
	}
	if printed {
		t.Fatal("printed on a partial thread.started line (no id yet); want no print until the line completes")
	}
	if _, err := ann.Write([]byte(stream[cut:])); err != nil {
		t.Fatal(err)
	}

	if log.String() != stream {
		t.Fatalf("log = %q, want the stream forwarded verbatim %q", log.String(), stream)
	}
	if !printed {
		t.Fatal("session id never printed after thread.started completed")
	}
	b := banner.String()
	if !strings.Contains(b, "session: sess-xyz") {
		t.Fatalf("banner = %q, want it to name the session id", b)
	}
	if !strings.Contains(b, "codex-run advise --resume sess-xyz") {
		t.Fatalf("banner = %q, want the ready-to-run headless resume command", b)
	}
	// A later chunk (more of the same stream) must not reprint.
	before := banner.Len()
	if _, err := ann.Write([]byte(`{"type":"item.completed"}` + "\n")); err != nil {
		t.Fatal(err)
	}
	if banner.Len() != before {
		t.Fatalf("banner grew on a later write (%q) — the session id must print exactly once", banner.String()[before:])
	}
}

// TestRunResumeBuildsResumeCommand pins the resume path: it must invoke `codex exec resume
// <id>`, drop the flags resume rejects (-C/--sandbox/--color), and send only the follow-up
// prompt — never the preamble, which the resumed session already holds in its history.
func TestRunResumeBuildsResumeCommand(t *testing.T) {
	tmp := t.TempDir()
	argsPath := filepath.Join(tmp, "args")
	stdinPath := filepath.Join(tmp, "stdin")
	logPath := filepath.Join(tmp, "events.jsonl")
	fakeCodex := filepath.Join(tmp, "codex")
	script := `#!/bin/sh
previous=
for argument do
	printf '%s\n' "$argument" >> "$CODEX_ARGS_FILE"
	if [ "$previous" = "--output-last-message" ]; then
		printf 'resumed answer\n' > "$argument"
	fi
	previous=$argument
done
cat > "$CODEX_STDIN_FILE"
printf '%s\n' '{"type":"thread.started","thread_id":"abc-123"}'
printf '%s\n' '{"type":"turn.completed","usage":{"input_tokens":1,"output_tokens":1}}'
`
	if err := os.WriteFile(fakeCodex, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CODEX_ARGS_FILE", argsPath)
	t.Setenv("CODEX_STDIN_FILE", stdinPath)
	t.Setenv("PATH", tmp+string(os.PathListSeparator)+os.Getenv("PATH"))

	prompt, err := os.CreateTemp(tmp, "prompt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := prompt.WriteString("please continue"); err != nil {
		t.Fatal(err)
	}
	if _, err := prompt.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	previousArgs, previousStdin := os.Args, os.Stdin
	os.Args, os.Stdin = []string{"codex-run", "advise", "--resume", "abc-123", "-l", logPath}, prompt
	t.Cleanup(func() {
		os.Args, os.Stdin = previousArgs, previousStdin
		prompt.Close()
	})

	if code := run(); code != 0 {
		t.Fatalf("run() = %d, want 0", code)
	}
	got, err := os.ReadFile(argsPath)
	if err != nil {
		t.Fatal(err)
	}
	arguments := "\n" + string(got)
	if !strings.Contains(arguments, "\nexec\nresume\nabc-123\n") {
		t.Fatalf("codex arguments:\n%s\nwant `exec resume abc-123`", got)
	}
	for _, forbidden := range []string{"\n-C\n", "\n--sandbox\n", "\n--color\n"} {
		if strings.Contains(arguments, forbidden) {
			t.Fatalf("codex arguments:\n%s\nresume must not pass %q (rejected by `exec resume`)", got, strings.TrimSpace(forbidden))
		}
	}
	stdin, err := os.ReadFile(stdinPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(stdin)) != "please continue" {
		t.Fatalf("codex stdin = %q, want only the follow-up prompt (no preamble re-injected)", stdin)
	}
}
