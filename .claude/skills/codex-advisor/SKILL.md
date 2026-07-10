---
name: codex-advisor
description: Cross-model second opinion from OpenAI Codex: `codex-consult --context` reviews the current session, or pipe a brief. Triggers on "ask codex", "codex review", "second opinion", or the "Cross-model review" rules.
---

# Codex advisor / reviewer

Cross-model outside view from OpenAI Codex. `codex-consult` reads the prompt from stdin; `--help` lists flags.

```bash
codex-consult --context                                    # review the current session ("advise me", no brief)
codex-consult --context < brief.md                         # session + a specific question
codex-consult -C <repo> --log /tmp/codex.jsonl < brief.md  # cold review of a repo (no session)
```

- `--context` sends the current session to OpenAI. Codex sees your actions, not your thinking or injected context — put load-bearing reasoning in the brief.
- Read-only by default; add `-s workspace-write` to let Codex edit (only when asked). stdout is the verdict.
- Implemented in Go under `~/.config/home-manager/scripts/{codex-consult,session-transcript}`; home-manager builds them and puts them on PATH.
