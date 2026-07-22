---
name: codex-work
description: Delegate a task's implementation wholesale to OpenAI Codex (`codex-run work`) — the whole requirements document in one long autonomous run — keeping policy, orchestration, and acceptance in the main loop. Use at the Implement step of the Orchestrator loop, for any multi-part build.
argument-hint: "[task-description]"
---

# Codex work

Dispatch the Implementer: one autonomous Codex run through the whole plan, committing as it goes. Flags: `codex-run --help`.

```bash
codex-run work -C <repo> --log /tmp/codex-work-<slug>.jsonl < brief.md  # one run per plan; unique log path
codex-run work --resume <id> -C <repo> < followup.md                    # send fixes, finish, or steer (id in the banner)
pkill -INT -f 'codex exec .*<repo>'                                     # interrupt a live run before resuming it
```

- **Brief** = the `/define` plan path plus only what the repo can't tell Codex: in-session decisions, branch, scope bounds. Never restate `~/.codex/AGENTS.md` or the wrapper's preamble.
- **Before** — branch off the default branch. A repo rooted at `$HOME` needs a `git worktree`; the wrapper refuses a writable run there.
- **During** — run it in the background and tail `--log`. Don't edit the checkout; give parallel runs separate worktrees.
- **Steer** the moment a run goes wrong — interrupt, then `--resume` with the correction. History survives, so write "drop X, do Y instead", not a re-brief; you lose only the unfinished turn.
- **After** — stdout is the STATUS report Accept judges. Exit 3 = finished but not COMPLETE → `--resume`. Other non-zero = turn failed; check `git status` for partial edits. Never resume a session another run is still writing.
