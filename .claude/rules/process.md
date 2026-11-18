# Process

## Roles
Two seats run every change request. For questions, reviews, and diagnosis, inspect and report — do not implement unless asked.
- **Orchestrator** — Claude Code, always: owns understanding, the design, and the outcome. It writes prose (docs, skills, prompts, PR text) and config-only or single-file edits itself or with Claude subagents, never Codex. Other code goes to the Implementer.
- **Implementer** — usually Codex (`codex-run work`, whose dispatch preamble names the seat): owns the build loop of each task and never pushes, ships, or touches external service state. An interactive Codex session holds both seats.

## Orchestrator loop
1. **Explore**.
2. **Plan** — `/define`.
3. **Implement** — send the Implementer one plan task at a time (`/codex-work`). Before a task's first run, write the Implementer's HEAD SHA into its plan entry as its **Tidy base**. Its fix, steer, and retry runs keep that base.
4. **Accept** — after each task, reconcile its report and its commits since the Tidy base against the plan. Rerun the repo's checks.
   - While you reconcile, never trust an exit code or a fluent summary.
   - Every finding ends fixed or skipped on measured evidence. Report each skip to the user and in the PR body, never as a scope exclusion in a later brief.
   - Findings go back to the Implementer, not into your own editor.
   - After the last task, if the change is high-risk or its behavior lacks test coverage, run `/quality-review` yourself.
5. **Ship** — `/pr`: it tidies, opens the PR, and runs `/address`. A push to a branch with an open PR is a follow-up, not a new Ship.

Auto-advance every step except two stops: plan approval and Ship. If the user authorized Ship at plan approval, do not stop at Ship. Ask for plan approval again when a requirement or the scope changes.

## Advisors
Two advisors sit outside the loop: Codex (`/codex-advise`) and Claude. Consult one before you commit to a hard-to-reverse or uncertain approach, assumption, or "done", and at any design choice the code cannot settle. Consult both when the question spans build and design, or before you skip a finding one of them raised. Report findings to the user at their severity, with your counter-evidence.

To reach the Claude advisor, use a `reviewer` subagent with `model: "opus"`.

## Implementer loop
1. **Implement**.
2. **Verify** — the repo's checks, then `/quality-review`.
3. Settle the findings — fix, or skip with evidence in the report.
4. **Commit** green units (`/commit`). Commit a fix to an earlier commit as `git commit --fixup=<sha>`.

At the end of each dispatched run, **Tidy** your task's commits (`/rebase-clean`) and report with evidence.

## Memory
Auto-memory is off. When the user gives an instruction that outlives the current step, even as an aside, write it down in the same turn: in the repo's `CLAUDE.md`/`AGENTS.md` if it holds for the repo, else in the work's plan file. If the work has no plan file, create one under `docs/plans/`. After a compaction, reread the plan file before you act. Keep other durable knowledge in the repo's docs.
