---
name: codex-work
description: Delegate a plan wholesale to OpenAI Codex for one autonomous implementation run. Use at the Implement step of the Orchestrator loop.
argument-hint: "[task-description]"
---

# Codex work

Dispatch the Implementer: one autonomous Codex run through the whole plan. Codex commits as it goes. Flags: `codex-run --help`.

```bash
codex-run work -C <repo> --log /tmp/codex-work-<slug>.jsonl < brief.md  # one run per plan; unique log path
codex-run work --resume <id> -C <repo> < followup.md                    # send fixes, finish, or steer (id in the banner)
pkill -INT -f 'codex exec .*<repo>'                                     # interrupt a live run before resuming it
```

- **Brief** = the `/define` plan path plus only what Codex cannot learn from the repo: in-session decisions, branch, scope bounds. Never restate `~/.codex/AGENTS.md` or the wrapper's preamble.
- Before — branch off the default branch. A repo rooted at `$HOME` needs a `git worktree`. The wrapper refuses a writable run there.
- During — run it in the background. Tail the `--log` file. Do not edit the checkout. Give parallel runs separate worktrees.
- As soon as a run goes wrong, **steer** it: interrupt, then `--resume` with the correction. The session keeps its history, so write "drop X, do Y instead", not a re-brief. You lose only the unfinished turn.
- After — stdout is the STATUS report that you judge at Accept. Exit 3 = finished but not COMPLETE → `--resume`. On any other non-zero exit, the turn failed: check `git status` for partial edits. Never resume a session while another run still writes to it.
