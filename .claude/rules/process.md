# Process

## Roles
Two seats run every change request. For questions, reviews, and diagnosis, inspect and report — do not implement unless asked.
- **Orchestrator** — Claude Code, always: owns understanding, the design, and the outcome, and builds nothing beyond trivial glue.
- **Implementer** — usually Codex (`codex-run work`, whose dispatch preamble names the seat): owns the build loop of each task and never pushes, ships, or touches external service state. An interactive Codex session holds both seats.

## Orchestrator loop
1. **Explore** — read the relevant flow end to end.
2. **Plan** — `/define` for requirements, `/decision-analysis` for an open choice.
3. **Implement** — settle open policy decisions. Then send the Implementer one task from the plan at a time (`/codex-work`).
4. **Accept** — after each task, reconcile the report and the diff since you dispatched the task against the plan. Rerun the repo's checks.
   - While you reconcile, never trust an exit code or a fluent summary.
   - Every finding ends fixed or skipped on measured evidence. If a finding depends on an unmeasured fact, measure that fact yourself before you decide.
   - Report each skip to the user and in the PR body, never as a scope exclusion in a later brief.
   - Findings go back to the Implementer, not into your own editor.
   - After the last task, if the change is high-risk or its behavior lacks test coverage, run `/feedback`, `/qa-review`, or both yourself.
5. **Ship** — `/pr`, then `/address`. Follow-up pushes close the loop, not a new Ship.

Auto-advance every step except the two stops: plan approval (end of Plan) and Ship. Before you commit to a hard-to-reverse or still-uncertain approach, assumption, or "done", consult `/codex-advise`. Report its findings to the user at the advisor's severity, with your counter-evidence attached. Skip `/codex-advise` for mechanical work.

## Implementer loop
1. **Implement**.
2. **Verify** — the repo's checks, then `/feedback`, and `/qa-review` for runnable behavior.
3. Settle the findings — fix, or skip with evidence in the report.
4. **Commit** green units (`/commit`).

At the end of each run, **Tidy** the commits that your current task added (`/rebase-clean`) and report with evidence.

Ask only when one of these blocks you:
- a value or scope boundary that you cannot infer
- missing credentials
- a destructive or irreversible action

## Memory
Auto-memory is off: keep nothing durable in harness-local state. Put knowledge to keep in a git-tracked file of the repo that it belongs to. Agent guidance goes in `CLAUDE.md`/`AGENTS.md`, and everything else goes in the repo's docs.
