---
argument-hint: [message]
description: Requirements definition command
---

# Purpose
Define requirements before implementation: clarify constraints, design policy, and specifications. Read-only; never implements.

# Rules
- Never modify, create, or delete files, and never implement — requirements only.
- Flag technically impossible requests; prioritize technical validity over preference.
- Ask without limit until requirements are clear; present questions before assuming.
- Drive every choice through AskUserQuestion, always offering a (Recommended) option.

# Workflow
1. **Analyze** — parse the request, identify constraints, draft candidate questions.
2. **Investigate** in parallel — spawn the read-only agents that fit:
   - `explore` — files, patterns
   - `design` — architecture, dependencies
   - `database` — schema, optimization
   - `general-purpose` — requirements analysis, estimation
   - `fact-check` — external sources
   - `validator` — cross-validation

   Also view the request through the 3–5 relevant decision lenses in `$HOME/.claude/skills/decision-analysis/references/lens-catalog.md` (feasibility, risk, impact, alternatives, security, …) to surface requirement questions and NFR coverage — fold these into the investigation, don't spawn a separate lens fan-out.
3. **Clarify** — score each candidate question by design-branching impact, irreversibility, investigation-impossibility, and effort (1–5 each); ask the highest-scoring first, and don't proceed without clear answers.
4. **Verify** — validate the user's decisions against technical evidence.
5. **Document** — produce the requirements spec and task breakdown for /execute handoff (format below).

# Output Format
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
