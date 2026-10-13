---
name: quality-review
description: Review a change or plan against its intent with parallel reviewers, then exercise its runnable behavior. Use after a change, to pressure-test a plan, or on any request for a review, critique, QA, or quality check.
---

# Review

Target: the argument, else `git diff $(git merge-base HEAD <base>)` plus each untracked file. `<base>` is the task's Tidy base in the Implementer seat, else `origin/main`. Intent: the plan task, the acceptance criteria, or the user's request.

Spawn at most three reviewers in one message. In Claude, use the `reviewer` agent type. In Codex, use the default agent type: the message's read-only line constrains it. The lenses are the sections of `$HOME/.claude/skills/quality-review/rubric.md`. Use Correctness and Simplicity always. If the change matches a trigger in the Risk section, add Risk. For a diff over about 1,000 lines, split it into at most three code areas instead. Give each reviewer one area with the chosen lenses. Send each reviewer this message:

    Read-only: do not edit files or change state. Read only your lenses in the rubric.
    Target: {diff command, or paths}
    Intent: {intent}
    Lens: {lens names} in $HOME/.claude/skills/quality-review/rubric.md
    Return one line per finding — `severity · file:line · failure scenario · fix` — or `No findings`. Severity is `high`, `medium`, or `low`. Judge only the change.

Merge findings that share a root cause. A finding conditional on an unmeasured fact ("slow if the table is large") is unfinished. Resolve the condition with a count, a timing, or an `EXPLAIN`, or report it at its worst-case severity.

To re-check, send each fix diff to the agents that raised its finding (Claude `SendMessage`, Codex `followup_task` or v1 `send_input`). Do not spawn a new set, and do not re-check a lens that had no findings.

For a plan target, ask each lens what the plan will get wrong, with one or two reviewers.

# Exercise

Settle the review's findings first. Then, if the change has runnable behavior and you can run its target (UI, CLI, API, DB), run adversarial scenarios. Cover each angle that applies:
- empty, out-of-order, and repeated input
- boundary values, a denied permission, a double submit
- the authoritative state: query the DB, file, or queue, not the screen
- neighbors of the change, after a reload or restart
- behavior against the spec, not the code

Drive a web UI with `agent-browser`. Run scenarios one at a time on shared state, or in parallel on disjoint surfaces or data. Give each parallel surface its own agent (not `reviewer`) and its own `agent-browser --session`. Fix nothing mid-run. Report each scenario as `executed` with evidence, or `not run` with the reason. A unit test that covers the scenario counts as `not run`.
