---
name: define
description: Define requirements before implementation — clarify constraints and design policy, then write a requirements spec and task breakdown for /execute to a gitignored plan file. Use at the Plan step of the operating loop, or whenever a request's scope needs pinning down before building. Never implements; the plan file is its only write.
argument-hint: "[message]"
---

# Rules
- Never implement, and never touch project files — the requirements document (step 5) is the only file you write.
- Flag technically impossible requests; prioritize technical validity over preference.
- Ask about requirement and value decisions you can't infer; settle derivable details (naming, structure, style) yourself and record them. Present open questions before assuming.
- Pose the questions you do ask through AskUserQuestion, always offering a (Recommended) option.

# Workflow
1. **Analyze** — parse the request, identify constraints, draft candidate questions.
2. **Investigate** — inline when the scope is small; for broad or independent surfaces spawn read-only agents in parallel: `Explore` for files and patterns, `reviewer` for a lensed read (architecture, schema, risk), `general-purpose` for external sources and estimation.

   Also view the request through the 3–5 relevant decision lenses in `$HOME/.claude/skills/decision-analysis/references/lens-catalog.md` (feasibility, risk, impact, alternatives, security, …) to surface requirement questions and NFR coverage — fold these into the investigation, don't spawn a separate lens fan-out.
3. **Clarify** — score each candidate question by design-branching impact, irreversibility, investigation-impossibility, and effort (1–5 each); ask the highest-scoring first, and don't proceed without clear answers.
4. **Verify** — validate the user's decisions against technical evidence.
5. **Document** — write the requirements spec and task breakdown (format below) to `<git root>/docs/plans/<YYYY-MM-DD>-<slug>.md` in the project. Chat output dies with the context window; the file is what post-compaction agents, the /execute delegate (Codex via `codex-run work`), and human reviewers read. Create the directory if needed, as it is ignored globally in `~/.config/git/ignore`. Reply with the file path and a short summary, not the full document.

# Output Format
The structure of the plan file (the chat reply is just its path and a summary):
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
- Execute handoff: decisions, references, constraints
```
