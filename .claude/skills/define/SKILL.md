---
name: define
description: Clarify constraints, then write a requirements spec and task breakdown to a gitignored plan file. Use at the Plan step.
argument-hint: "[message]"
---

# Rules
- Never implement, and never touch project files — the plan file (step 5) is the only file you write.
- Ask the user only about requirement and value decisions that investigation cannot settle. Pose each question through AskUserQuestion, and always offer a (Recommended) option.
- Some problems only the user can fix: missing access, or a broken system outside the repo. If you find one, report it and end your turn. Do not work around it.
- Settle every design decision that more than one task depends on, and every new interface, schema, module boundary, or dependency. Record each one in Technical Specifications. Leave the details inside one task, such as its local names, to the Implementer.
- Cut the work into tasks for `/codex-work`. Each task must make one change: a behavior, a refactor, or a migration. It must end with a check that proves the change and leave the repo's checks passing.

# Workflow
1. Analyze — identify constraints and draft candidate questions.
2. Investigate what you have not read yet. If the scope is small, investigate inline. For broad or independent surfaces, spawn read-only agents in parallel:
   - `Explore` for files and patterns
   - `reviewer` for a lensed read (architecture, schema, risk)
   - `general-purpose` for external sources and estimation

   Also view the request through the 3–5 relevant decision lenses in `$HOME/.claude/skills/decision-analysis/references/lens-catalog.md`. Use them to surface requirement questions and NFR coverage. Apply them yourself, or name them in the prompts of the agents above. Do not spawn other agents for the lenses.
3. Clarify — first ask the questions whose answers most change the design, then those whose answers are hardest to reverse. Do not proceed without clear answers.
4. Verify the request and the user's decisions against technical evidence. If one conflicts with the evidence, show the conflict and ask again.
5. Document — write the requirements spec and task breakdown (format below) to `<git root>/docs/plans/<YYYY-MM-DD>-<slug>.md`. The file is what post-compaction agents, the Implementer, and the user read. If the directory does not exist, create it: Git ignores it globally through `~/.config/git/ignore`. Reply with the file path and a short summary, not the full document. Then end your turn for plan approval.

# Output Format
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
