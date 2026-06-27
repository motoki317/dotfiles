---
name: codex-advisor
description: Independent cross-model second opinion from OpenAI Codex — an outside view that catches blind spots a Claude-only review (including the built-in advisor) shares. Invoke on your own judgment when a decision is genuinely uncertain or hard to reverse, when you diverge from the built-in advisor, or before finishing non-trivial work; also on explicit request ("ask codex", "codex review", "get a second opinion"). Read-only by default; full-access only when asked to apply changes.
---

# Codex advisor / reviewer

`codex-consult` (on `PATH`) consults OpenAI Codex as an **independent, cross-model** reviewer. Codex is a different model family from Claude, so it catches blind spots a Claude-only review — including the built-in `advisor` tool — shares. The two are complementary: `advisor` is transcript-aware; Codex starts cold and reads the repo.

When to invoke is in the user's CLAUDE.md ("Cross-Model Review"): pair it with the built-in `advisor` at substantive checkpoints. This file covers how to run it.

## Running it

Write the task to a file and pipe it in, anchored to the repo, with a known log path:

```bash
codex-consult -C <repo> --log /tmp/codex.jsonl < /tmp/brief.md
```

The notes below adjust this one form; you rarely need anything else.

- **Prompt = task only.** The preamble already gives Codex its reviewer role and tells it to close with a verdict, so the brief carries only context and the question. Pipe it via stdin (a file written with the Write tool): stdin is never shell-interpreted, so backticks, `$(…)`, `$VAR`, and `!` in a code-heavy prompt cannot corrupt it. Reserve an inline `"question"` argument for a short one-liner free of those characters. To review a diff, append it: `cat /tmp/brief.md <(git diff main) | codex-consult -C <repo> --log /tmp/codex.jsonl`.
- **Scope the root.** Keep `-C` on the project, not a broad parent like `$HOME`: it bounds what Codex reads, and in a write sandbox what it may change. The tool warns when the root is `$HOME`.
- **Watch a long run.** A consult can take minutes. Run the command in the background (Bash `run_in_background`) and poll `--log`, which is Codex's raw `--json` event stream written live (one JSON object per line: a `command_execution` event per tool Codex runs, `agent_message` for its narration, a final `turn.completed`). New lines mean progress; parse them or just watch the file grow. A silent gap is only a concern when it is long *and* the process is still alive, because a pure-reasoning phase emits nothing between `turn.started` and the answer.
- **Apply changes (only when asked).** A consult is read-only by design. Raise the sandbox only when the user asks Codex to *apply* edits: `-s workspace-write` keeps changes inside the repo; the `cdy` alias (`codex --dangerously-bypass-approvals-and-sandbox`) grants full access for an interactive session.
- **Other flags** (`--help`): `-m <model>` overrides the model, `-v` also streams the transcript to stderr.

## Handling the result

stdout is the verdict — relay it. Weigh it as you weigh the built-in advisor's: evaluate, don't apply blindly, and when the two disagree surface the conflict to the user rather than silently picking a side.

A successful run emits a `turn.completed` event in the log and exits 0. A non-zero exit with no `turn.completed` means Codex failed — the reason is a `turn.failed` (or `error`) event in the log. stderr reports the log path (at start) and token usage (at end).

The wrapper is a thin shim around the self-building Go source at `~/.claude/skills/codex-advisor/scripts/codex-consult.go`; edit that to change the helper itself.
