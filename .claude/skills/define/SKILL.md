---
name: define
description: Clarify constraints, then write a requirements spec and task breakdown to a gitignored plan file. Use at the Plan step.
argument-hint: "[message]"
---

# Rules
- Never implement, and never touch project files. The plan file is the only file you write.
- Ask the user only about requirement and value decisions that investigation cannot settle. First ask the questions whose answers most change the design, then those whose answers are hardest to reverse. Pose each question through AskUserQuestion, and always offer a (Recommended) option.
- If the request or a user decision conflicts with the technical evidence, show the conflict and ask again.
- Some problems only the user can fix: missing access, or a broken system outside the repo. If you find one, report it and end your turn. Do not work around it.
- Settle every design decision that more than one task depends on, and every new interface, schema, module boundary, or dependency. Record each one in Technical Specifications, with the standard approach or upstream prior art that you found, or "none found". Leave the details inside one task, such as its local names, to the Implementer.
- Before you settle a new interface, schema, or module boundary, read `$HOME/.claude/refs/single-source-of-truth.md`.
- Cut the work into tasks for the Implementer. Each task must make one change: a behavior, a refactor, or a migration. It must end with a check that proves the change and leave the repo's checks passing.

# Output
Write the plan in this format to `<git root>/docs/plans/<YYYY-MM-DD>-<slug>.md`. Reply with the path and the decisions that the user must approve, not the document. Then end your turn for plan approval.

```
## Requirements Document
- Summary: One-sentence request, background, expected outcomes
- Current State: Existing system, tech stack
- Functional Requirements: FR-001 format, each marked mandatory or optional
- Non-Functional Requirements: Performance, security, maintainability
- Technical Specifications: Design policies, impact scope
- Constraints: Technical, operational
- Test Requirements: Unit, integration, acceptance criteria
- Outstanding Issues: open questions that block no task

## Task Breakdown
- Dependency graph between the tasks
- Per task: an ID such as T1, the FR IDs it covers, the change it makes, files, the check that proves it, and a `Tidy base:` line that the Orchestrator fills at dispatch
```
