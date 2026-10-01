---
name: define
description: Clarify constraints, then write a requirements spec and task breakdown to a gitignored plan file. Use at the Plan step. This skill never implements.
argument-hint: "[message]"
---

# Rules
- Never implement, and never touch project files — the requirements document (step 5) is the only file you write.
- Flag technically impossible requests. Prioritize technical validity over preference.
- Ask about requirement and value decisions that you cannot infer. Settle derivable details (naming, structure, style) yourself. Record them. Present open questions before you make assumptions.
- Pose the questions that you do ask through AskUserQuestion. Always offer a (Recommended) option.

# Workflow
1. Analyze — parse the request. Identify constraints. Draft candidate questions.
2. Investigate — if the scope is small, investigate inline. For broad or independent surfaces, spawn read-only agents in parallel:
   - `Explore` for files and patterns
   - `reviewer` for a lensed read (architecture, schema, risk)
   - `general-purpose` for external sources and estimation

   Also view the request through the 3–5 relevant decision lenses in `$HOME/.claude/skills/decision-analysis/references/lens-catalog.md` (for example, feasibility, risk, impact, alternatives, security). Use them to surface requirement questions and NFR coverage. Apply the lenses within the investigation: do not spawn a separate lens fan-out.
3. Clarify — ask these questions first:
   - questions whose answers most change the design
   - questions whose answers are hardest to reverse
   - questions whose answers investigation cannot settle

   Do not proceed without clear answers.
4. Verify the user's decisions against technical evidence.
5. Document — write the requirements spec and task breakdown (format below) to `<git root>/docs/plans/<YYYY-MM-DD>-<slug>.md` in the project. Chat output disappears with the context window. The file is what post-compaction agents, the /codex-work delegate (Codex via `codex-run work`), and human reviewers read. If the directory does not exist, create it: Git ignores it globally through `~/.config/git/ignore`. Reply with the file path and a short summary, not the full document.

# Output Format
The structure of the plan file:
```
## Requirements Document
- Summary: One-sentence request, background, expected outcomes
- Current State: Existing system, tech stack
- Functional Requirements: FR-001 format (mandatory/optional)
- Non-Functional Requirements: Performance, security, maintainability
- Technical Specifications: Design policies, impact scope
- Constraints: Technical, operational
- Test Requirements: Unit, integration, acceptance criteria
- Outstanding Issues: Unresolved questions

## Task Breakdown
- Dependency graph
- Phased tasks with files and overview
- Implementer handoff: decisions, references, constraints
```
