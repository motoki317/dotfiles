---
name: execute
description: Execute a task by delegating the implementation wholesale to OpenAI Codex (`codex-run work`) — the entire requirements document in one long autonomous run — keeping policy decisions, orchestration, and monitoring in the main loop, then verifying the result with tests and the feedback/qa-review skills. Use at the Implement step of the operating loop, for any multi-part build.
argument-hint: "[task-description]"
---

# Rules
- Codex is the primary implementer — assume it is smarter and more thorough than you at coding and at testing its own work. Delegate the whole requirements document in one long run; don't subdivide it into per-file tasks and don't implement detailed logic yourself. Only when the plan clearly exceeds one run's context, split at plan-defined, independently verifiable milestones — never finer.
- Only trivial glue (a config line, a rename) may use Read/Edit/Write directly; anything larger goes to Codex.
- Never edit the checkout while a Codex run is active; parallel runs need separate worktrees.
- Don't declare the task done until the Verify step passes; show evidence, not assertions.

# Workflow
1. **Prepare** — read the `/define` plan file end to end and settle any open policy decisions yourself. Choose or create the working branch (branch off the default branch first — Codex commits as it goes). For a repo rooted at `$HOME` (this dotfiles repo), delegate from a `git worktree` checkout — the wrapper refuses a writable run rooted at `$HOME`. Write the brief: the plan-file path plus only what Codex can't infer from the repo — in-session decisions, the branch, scope boundaries. Codex auto-loads the coding principles, conventions, and skill playbooks from the git-tracked `~/.codex/AGENTS.md`, and the wrapper's preamble fixes the commit/never-push policy and report format; don't restate either.
2. **Delegate** — one long-running call, in the background (it can exceed foreground command timeouts), with a unique log path per run:
   ```bash
   codex-run work -C <repo> --log /tmp/codex-work-<slug>.jsonl < brief.md
   ```
   Full permission by default (`danger-full-access`): Codex commits locally, fetches dependencies, and runs tests. Let it run to completion.
3. **Monitor** — tail the `--log` file to follow progress; don't interrupt a healthy run. Every run's stderr banner prints its Codex `session: <id>` and a ready-to-run `codex-run work --resume <id>`; resuming continues that same session with its context intact, so prefer it over a cold re-delegate whenever the session is still worth continuing (never resume one while another run is still writing it). Exit codes: `3` = the turn finished but the report's STATUS is not COMPLETE — resume the session to finish the remainder, or re-delegate citing its checkpoint commits. Other non-zero = the turn failed or never launched: read stderr and the log tail (partial edits may remain — `git status`/`git log`). For a transient failure (network drop, timeout, killed process), `--resume <id>` recovers the run without losing context; reserve a cold re-delegate or a fix-forward brief for when the approach itself was wrong. If failures recur, switch approach per the error-handling rules.
4. **Verify** — Claude's review of Codex's work, after the run completes: compare Codex's requirement-by-requirement audit against the plan file — never infer success from the exit code or a fluent summary — and run the available tests/build/lint yourself. Then invoke the `feedback` skill for static review, and the `qa-review` skill when there is runnable behavior — for analysis only; tell each not to run its close-the-loop phase, since this step owns the fix loop. Fix critical and major findings: trivial ones directly, substantial ones re-delegated as a new `codex-run work` brief citing the findings. Re-check and report the evidence — commands run and results — and don't finish with unresolved critical or major findings.

# Agents
Spawn `reviewer` (read-only) for verification lenses and cross-checks, and `general-purpose` for glue around the delegation (environment setup, data gathering). Implementation belongs to Codex.
