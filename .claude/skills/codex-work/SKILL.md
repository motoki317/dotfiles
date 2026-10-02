---
name: codex-work
description: Delegate a plan wholesale to OpenAI Codex for one autonomous implementation run. Use at the Implement step of the Orchestrator loop.
argument-hint: "[task-description]"
---

# Codex work

Dispatch the Implementer: one autonomous Codex run through the whole plan. Codex commits as it goes. Flags: `codex-run --help`.

```bash
codex-run work -C <repo> < brief.md                # one run per plan
codex-run work --resume <id> < followup.md         # send fixes, finish, steer, or retry (id in the banner)
pkill -INT -f 'codex exec .*-C <repo> --sandbox'   # interrupt a fresh run
pkill -INT -f 'codex exec resume <id>'             # interrupt a resumed run
```

- **Brief** = the `/define` plan path plus only what Codex cannot learn from the repo: in-session decisions, branch, scope bounds. Never restate `~/.codex/AGENTS.md` or the wrapper's preamble.
- Before — branch off the default branch. A repo rooted at `$HOME` needs a `git worktree`. The wrapper refuses a writable run there.
- During — run it with `run_in_background: true` and `timeout: 7200000`. Then wait for its completion notification. Do not poll the run's log, sleep on it, or arm a Monitor on it. Read the log only to answer a progress question from the user, or to diagnose a failed or killed run. The stderr banner gives its path. Do not edit the checkout. Give parallel runs separate worktrees.
- If you learn during a run that its brief is wrong, **steer** it. Interrupt it, wait for its completion notification, then `--resume` with the correction. The session keeps its history, so write "drop X, do Y instead", not a re-brief.
- After — stdout is the STATUS report that you judge at Accept. Exit 3 = finished but not COMPLETE: read the report, settle each blocker, then `--resume` with the resolution. If the run failed or the background time limit killed it, check `git status` for partial edits, then `--resume` to retry. Never resume a session while another run still writes to it.
