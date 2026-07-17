---
name: codex-work
description: Delegate a task's implementation wholesale to OpenAI Codex (`codex-run work`) — the whole requirements document in one long autonomous run — keeping policy, orchestration, and acceptance in the main loop. Use at the Implement step of the Orchestrator loop, for any multi-part build.
argument-hint: "[task-description]"
---

# Rules
- Delegate the whole plan in one run — don't subdivide into per-file tasks, don't write the logic yourself. Split only when the plan exceeds one run's context, at plan-defined, independently verifiable milestones — never finer.
- Never edit the checkout while a run is active; parallel runs need separate worktrees.

# Workflow
1. **Prepare** — read the `/define` plan file end to end; settle open policy decisions yourself. Branch off the default branch (Codex commits as it goes). A repo rooted at `$HOME` (this dotfiles repo) needs a `git worktree` checkout — the wrapper refuses a writable run at `$HOME`. Brief = the plan-file path plus only what Codex can't infer from the repo: in-session decisions, the branch, scope boundaries. Don't restate `~/.codex/AGENTS.md` (auto-loaded) or the wrapper's preamble (Implementer seat, safety rails, report format).
2. **Delegate** — one long background call (it can outlast foreground timeouts), unique log path per run.
   ```bash
   codex-run work -C <repo> --log /tmp/codex-work-<slug>.jsonl < brief.md
   ```
   `danger-full-access` by default: Codex commits locally, fetches deps, runs tests. Let it run to completion.
3. **Monitor** — tail `--log`; don't interrupt a healthy run. The stderr banner prints `session: <id>` and a ready-to-run `--resume` line; resuming keeps that session's context, so prefer it to a cold re-delegate — but never resume a session another run is still writing. Exit `3` = turn finished, report STATUS not COMPLETE → resume to finish, or re-delegate citing its checkpoint commits. Other non-zero = turn failed or never launched → read stderr and the log tail; partial edits may remain (`git status`/`git log`). `--resume` also recovers a transient failure (network drop, timeout, kill); cold re-delegate only when the approach itself was wrong. Recurring failures → change approach; don't retry harder.
4. **Accept** — never infer success from an exit code or a fluent summary: on top of the loop's reconcile-and-rerun, audit the requirement-by-requirement report against the plan file. Escalate to `feedback` / `qa-review` when the change is high-risk or externally visible behavior lacks acceptance coverage. Fix trivial findings directly; send substantial ones back via `--resume` with a brief citing them. Report the commands and results; don't finish with critical or major findings open.

# Agents
Spawn `reviewer` (read-only) for verification lenses, `general-purpose` for glue around the delegation (environment setup, data gathering). Implementation belongs to Codex.
