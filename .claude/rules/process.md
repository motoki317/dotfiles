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
   - Every finding ends fixed or skipped on measured evidence. If a finding depends on an unmeasured fact, measure that fact yourself before you decide.
   - Report each skip to the user and in the PR body, never as a scope exclusion in a later brief.
   - Findings go back to the Implementer, not into your own editor.
   - After the last task, if the change is high-risk or its behavior lacks test coverage, run `/quality-review` yourself.
5. **Ship** — `/pr`: it tidies, opens the PR, and runs `/address`. A push to a branch with an open PR is a follow-up, not a new Ship.

Auto-advance every step except two stops: plan approval and Ship. Ask for plan approval again when a requirement or the scope changes.

## Advisors
Two advisors sit outside the loop: Codex (`/codex-advise`) and Claude. Consult one before you commit to a hard-to-reverse or uncertain approach, assumption, or "done", and at any design choice the code cannot settle. Consult both when the question spans build and design, or before you skip a finding one of them raised. Skip them for mechanical work. Report findings to the user at their severity, with your counter-evidence.

To reach the Claude advisor, run `env | grep '^ANTHROPIC_BASE_URL='`:

- Set to a local server: every subagent here is that engine, whatever its model label. Call `claude-upstream` once per gate; never fan out.
- Unset: use a subagent with `model: "opus"`, not `claude-upstream`.

## Implementer loop
1. **Implement**.
2. **Verify** — the repo's checks, then `/quality-review`.
3. Settle the findings — fix, or skip with evidence in the report.
4. **Commit** green units (`/commit`). Commit a fix to an earlier commit as `git commit --fixup=<sha>`.

At the end of each dispatched run, **Tidy** your task's commits (`/rebase-clean`) and report with evidence.

Ask only when one of these blocks you:
- a value or scope boundary that you cannot infer
- missing credentials
- a destructive or irreversible action

## Memory
Auto-memory is off. Keep durable knowledge in a git-tracked file of its repo: agent guidance in `CLAUDE.md`/`AGENTS.md`, the rest in the repo's docs.
