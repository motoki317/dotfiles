---
name: codex-advise
description: Cross-model second opinion from OpenAI Codex. Use at Orchestrator consult gates or when asked for a codex or second opinion.
---

# Codex advisor

Cross-model outside view from OpenAI Codex. The prompt is read from stdin. Flags: `codex-run --help`.

```bash
codex-run advise --context                                    # review this session, no brief
codex-run advise --context < brief.md                         # session + a specific question
codex-run advise -C <repo> --log /tmp/codex.jsonl < brief.md  # cold review of a repo, no session
codex-run advise --resume <id> < followup.md                  # continue a review (id in the banner; only after it exits)
```

- Gate occasions (Orchestrator loop): locking in an approach or interpretation, building on a load-bearing assumption, declaring non-trivial work done, stuck or changing approach. Hard-to-reverse or genuinely uncertain calls make the consult mandatory; each call costs minutes, tokens, and egress to OpenAI — skip mechanical or low-stakes work.
- `--context` sends the session transcript to OpenAI. Codex sees your actions, not your thinking — put load-bearing reasoning in the brief and name the files.
- A gate call (`~/.claude/rules/process.md`) runs foreground; background only a non-gating opinion.
- Read-only unless asked; `-s workspace-write` lets Codex run tests. stdout is the verdict.
