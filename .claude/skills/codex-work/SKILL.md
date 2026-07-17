---
name: codex-work
description: Delegate a task's implementation wholesale to OpenAI Codex (`codex-run work`) — the whole requirements document in one long autonomous run — keeping policy, orchestration, and acceptance in the main loop. Use at the Implement step of the Orchestrator loop, for any multi-part build.
argument-hint: "[task-description]"
---

# Codex work

Dispatch the Implementer: `codex-run work` runs OpenAI Codex through the whole plan in one autonomous run, committing as it goes (`danger-full-access` by default). `codex-run --help` lists the flags.

```bash
codex-run work -C <repo> --log /tmp/codex-work-<slug>.jsonl < brief.md   # one run per plan; unique log path
codex-run work --resume <id> -C <repo> < findings.md                     # send fixes back, or finish an interrupted turn
```

- Brief = the `/define` plan-file path plus only what Codex can't infer from the repo: in-session decisions, the branch, scope boundaries. Don't restate `~/.codex/AGENTS.md` or the wrapper's preamble (Implementer seat, safety rails, report format).
- Branch off the default branch first. A repo rooted at `$HOME` (this dotfiles repo) must be delegated as a `git worktree` — the wrapper refuses a writable run there.
- Run it in the background (a run can outlast foreground timeouts) and tail `--log`; never edit the checkout while a run is active, and give parallel runs separate worktrees.
- Exit `3` = turn finished but STATUS is not COMPLETE → `--resume` to finish. Other non-zero = turn failed → read stderr and the log tail; partial edits may remain (`git status`). `--resume` (id in the stderr banner) keeps the session's context — prefer it to a cold re-delegate, but never resume a session another run is still writing.
- stdout is the final report: STATUS line plus the requirement-by-requirement audit that Accept judges.
