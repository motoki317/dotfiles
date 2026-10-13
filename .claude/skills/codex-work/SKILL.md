---
name: codex-work
description: Delegate one task from a plan to OpenAI Codex, which implements, verifies, and commits it. Use at the Implement step of the Orchestrator loop.
argument-hint: "[plan file] [task ID]"
---

# Codex work

Flags: `codex-run --help`.

```bash
codex-run work -C <repo> < brief.md                # first task
codex-run work --resume <id> < next.md             # next task, fixes, a steer, or a retry (id in the banner)
pkill -INT -f 'codex exec .*-C <repo> --sandbox'   # interrupt a fresh run
pkill -INT -f 'codex exec resume <id>'             # interrupt a resumed run
```

- **Brief** = the absolute path of the `/define` plan file and the task ID. Add only what Codex cannot learn from the repo: in-session decisions and scope bounds. For a brief without a plan file, such as an `/address` review fix, name the Tidy base. For a fix to an accepted commit, name its SHA as the `--fixup` target.
- Before a plan's first task, branch off the default branch. If two plans run at the same time, give each its own worktree.
- Dispatch the bare command with `run_in_background: true` and no `timeout`: no output redirect, pipe, or trailing `echo`. The completion notice then carries the exit code, and its output file holds the STATUS report and the log path.
- Wait for that notice: do not poll the log, sleep on it, or arm a Monitor. Read the log only to answer a progress question from the user, or to diagnose a failed or killed run. A goal check-in is not a progress question. During a run, do not edit the checkout.
- If you learn during a run that its brief is wrong, **steer** it. Interrupt it, wait for its completion notice, then `--resume` with the correction. The session keeps its history, so write "drop X, do Y instead", not a re-brief.
- Exit 3 = finished but not COMPLETE: read the report, settle each blocker, then `--resume` with the resolution. If the run failed or was killed, check `git status` for partial edits, then `--resume` to retry.
- After you accept a task, update the plan's remaining tasks with what it revealed. Then send the next task with `--resume`.
