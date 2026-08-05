# Process

## Roles
Two seats run every change request; for questions, reviews, and diagnosis, inspect and report — don't implement unless asked.
- **Orchestrator** — Claude Code, always: owns understanding and the outcome; builds nothing beyond trivial glue.
- **Implementer** — usually Codex (`codex-run work`; the dispatch preamble names the seat): owns the build loop; never pushes, ships, or touches external service state. An interactive Codex session holds both seats.

## Orchestrator loop
1. **Explore** — read the relevant flow end to end.
2. **Plan** — `/define` for requirements; `/decision-analysis` for an open choice.
3. **Implement** — settle open policy decisions, then hand the whole plan — the brief — to the Implementer (`/codex-work`).
4. **Accept** — reconcile the report and diff against the plan — never trusting an exit code or a fluent summary — and rerun the repo's checks. Every finding ends fixed or skipped on measured evidence: when it hinges on an unmeasured fact, measure it yourself before deciding; report each skip to the user and in the PR body, never as a scope exclusion in a later brief. Findings go back to the Implementer, not into your own editor; run `/feedback` / `/qa-review` yourself when the change is high-risk or its behavior lacks test coverage.
5. **Ship** — `/pr`, then `/address`; follow-up pushes close the loop, not a new Ship.

Auto-advance every step except the two stops: plan approval (end of Plan) and Ship. Consult `/codex-advise` before committing to an approach, an assumption, or a "done" that is hard to reverse or still uncertain — report its findings to the user at the advisor's severity, your counter-evidence attached; skip it for mechanical work.

## Implementer loop
**Implement** → **Verify** (the repo's checks, then `/feedback`; `/qa-review` for runnable behavior) → settle the findings (fix, or skip with evidence in the report) → **Commit** green units (`/commit`); at the end, **Tidy** this run's own commits (`/rebase-clean`) and report with evidence. Ask only when blocked on a value or scope boundary you can't infer, missing credentials, or a destructive or irreversible action.
