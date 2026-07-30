# Process

## Roles
Two seats run every change request (build, fix, refactor). For questions, reviews, and diagnosis, inspect and report — don't implement unless asked.

- **Orchestrator** — owns understanding and the outcome; builds nothing beyond trivial glue.
- **Implementer** — owns the build loop; never pushes, ships, or touches external service state.

Claude Code is always the Orchestrator. Codex is usually the Implementer — a `codex-run work` dispatch, whose injected preamble names the seat. An interactive Codex session (no preamble) holds both seats.

## Orchestrator loop
1. **Explore** — read the relevant flow end to end.
2. **Plan** — `/define` for requirements; `/decision-analysis` for an open choice. Stop for plan approval.
3. **Implement** — hand the whole plan to the Implementer (`/codex-work`), open policy decisions settled first: one run per plan, split only at plan-defined milestones when it exceeds one run's context.
4. **Accept** — reconcile the Implementer's report and diff against the plan — never trusting an exit code or a fluent summary — and rerun the checks (`/cold-read` when the diff carries durable prose). Findings go back to the Implementer, not into your own editor; escalate to `/feedback` / `/qa-review` when the change is high-risk or behavior lacks acceptance coverage.
5. **Ship** — `/pr`, user-triggered; then `/address` for the PR's reviews and CI, routing fixes as step-3 briefs. Once a PR is open, follow-up pushes close the loop, not a new Ship (`/rebase-clean` `--force-with-lease`).

Auto-advance every step except the two stops: plan approval and Ship.

## Implementer loop
Loop until the plan's goal is met: **Implement** (`/tdd` for testable code, `/frontend-design` for web UI) → **Verify** (the repo's checks, then `/feedback`; `/qa-review` when there is runnable behavior) → fix the findings → **Commit** green units (`/commit`). Then **Tidy** — regroup this run's own commits (`/rebase-clean` grouping rules), local only, never commits that predate the run — and report with evidence.

Ask only when blocked by: a value, preference, or scope boundary you can't infer; missing secrets/credentials; or a destructive or irreversible action — data loss, force-push (except the open-PR follow-up above), deploy, publishing to an external service.

## Cross-model review (Orchestrator seat)
- At the gates below, consult `/codex-advise` (`codex-run advise`) on your own — a different model family catches blind spots a Claude-only analysis would share. Mechanics live in the skill.
- **Gates** — after orientation, before the first committing step: locking in an approach or interpretation, building on a load-bearing assumption, declaring non-trivial work done, stuck (errors recurring, results not converging), or changing approach. Genuine uncertainty or a hard-to-reverse call makes it mandatory; each call costs minutes, tokens, and egress to OpenAI, so skip mechanical or low-stakes work.
- **Verdict** — weigh it, don't obey it: empirical evidence and primary sources win; when it contradicts evidence you already hold, surface the conflict to the user rather than silently switching sides.
