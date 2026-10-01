---
name: codex-advise
description: Cross-model second opinion from OpenAI Codex. Use at Orchestrator consult gates or when asked for a codex or second opinion.
---

# Codex advisor

Get a cross-model outside view from OpenAI Codex. `codex-run` reads the prompt from stdin. Flags: `codex-run --help`.

```bash
codex-run advise --context                                    # review this session, no brief
codex-run advise --context < brief.md                         # session + a specific question
codex-run advise -C <repo> --log /tmp/codex.jsonl < brief.md  # cold review of a repo, no session
codex-run advise --resume <id> < followup.md                  # continue a review (id in the banner; only after it exits)
```

- Gate occasions (Orchestrator loop):
  - when you lock in an approach or interpretation
  - when you build on a load-bearing assumption
  - when you declare non-trivial work done
  - when you are stuck or change approach

  Hard-to-reverse or still uncertain decisions make the consult mandatory. Each call costs minutes, tokens, and egress to OpenAI. Skip the consult for mechanical or low-stakes work.
- `--context` sends the session transcript to OpenAI. Codex sees your actions, not your thinking. Put load-bearing reasoning and the file names in the brief.
- Run a gate call (`~/.claude/rules/process.md`) in the foreground. Run only a non-gating opinion in the background.
- Codex runs read-only unless asked. `-s workspace-write` lets Codex run tests. stdout is the verdict.
