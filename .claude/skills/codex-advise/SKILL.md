---
name: codex-advise
description: Cross-model second opinion from OpenAI Codex. Use at Orchestrator consult gates, on an open design choice, or when asked for a codex or second opinion.
---

# Codex advisor

`codex-run` reads the brief from stdin, and stdout is the verdict. Flags: `codex-run --help`.

```bash
codex-run advise -C <repo> --log <scratchpad>/advise-<slug>.jsonl < brief.md
codex-run advise --resume <id> < followup.md   # only after the run exits (its banner shows <id>)
```

Codex sees only the brief and the repo. In the brief, name the files to read first and the decision, and state your load-bearing reasoning. Ask numbered questions, and request findings by severity, each with file:line evidence.

Act on the gated decision only after you read the verdict.
