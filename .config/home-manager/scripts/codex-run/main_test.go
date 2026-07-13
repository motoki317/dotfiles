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
