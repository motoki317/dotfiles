---
name: codex-work
description: Delegate one task from a plan to OpenAI Codex, which implements, verifies, and commits it. Use at the Implement step of the Orchestrator loop.
argument-hint: "[task-description]"
---

# Codex work

Flags: `codex-run --help`.

```bash
codex-run work -C <repo> < brief.md                # first task
codex-run work --resume <id> < next.md             # next task, fixes, a steer, or a retry (id in the banner)
pkill -INT -f 'codex exec .*-C <repo> --sandbox'   # interrupt a fresh run
pkill -INT -f 'codex exec resume <id>'             # interrupt a resumed run
```

- **Brief** = the absolute path of the `/define` plan file and the task ID. Add only what Codex cannot learn from the repo: in-session decisions and scope bounds. Name the Tidy base in every brief and follow-up: the commit that HEAD pointed to when you dispatched the current task. Never restate `~/.codex/AGENTS.md` or the wrapper's preamble.
- Before the first task, branch off the default branch. A repo rooted at `$HOME` needs a `git worktree`.
- Dispatch each run with `run_in_background: true` and `timeout: 7200000`. Then wait for its completion notification. Do not poll the run's log, sleep on it, or arm a Monitor on it. Read the log only to answer a progress question from the user, or to diagnose a failed or killed run. The stderr banner gives its path.
- During a run, do not edit the checkout. If two plans run at the same time, give each its own worktree.
- If you learn during a run that its brief is wrong, **steer** it. Interrupt it, wait for its completion notification, then `--resume` with the correction. The session keeps its history, so write "drop X, do Y instead", not a re-brief.
- After a run — stdout is the STATUS report that you judge at Accept. Exit 3 = finished but not COMPLETE: read the report, settle each blocker, then `--resume` with the resolution. If the run failed or the background time limit killed it, check `git status` for partial edits, then `--resume` to retry. Never resume a session while another run still writes to it.
- After you accept a task, update the remaining tasks in the plan file with what this task revealed. If an update changes a requirement or the scope, the plan needs approval again: stop and ask the user. Then send the next task to the same session with `--resume`.
