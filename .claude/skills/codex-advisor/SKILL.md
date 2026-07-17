---
name: codex-advisor
description: Cross-model second opinion from OpenAI Codex — `codex-run advise --context` reviews the current session, or pipe a brief. Triggers on "ask codex", "codex review", "second opinion", or the "Cross-model review" rules.
---

# Codex advisor

Cross-model outside view from OpenAI Codex. `codex-run advise` reads the prompt from stdin; `codex-run --help` lists the flags.

```bash
codex-run advise --context                                    # review the current session ("advise me", no brief)
codex-run advise --context < brief.md                         # session + a specific question
codex-run advise -C <repo> --log /tmp/codex.jsonl < brief.md  # cold review of a repo (no session)
```

- `--context` sends the session to OpenAI. Codex sees your actions, not your thinking or injected context — put load-bearing reasoning in the brief and point Codex at the files.
- A gate call (`~/.claude/rules/process.md`) runs foreground — wait for the verdict; background only a non-gating opinion.
- Read-only by default; `-s workspace-write` lets Codex run tests — only when asked. stdout is the verdict.
- A dropped call prints `session: <id>`; `codex-run advise --resume <id> < followup.md` continues that review — a network drop, or a follow-up on the same context — instead of a cold call. Resume only after the call exits.
- Source: `~/.config/home-manager/scripts/{codex-run,session-transcript}` (Go, built onto PATH by home-manager).
