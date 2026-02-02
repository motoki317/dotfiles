---
argument-hint: [message]
description: Requirements definition command
---

# Purpose
Conduct detailed requirements definition before implementation, clarifying technical constraints, design policies, and specifications.

# Rules
- **Never modify, create, or delete files** (read-only)
- Never implement code; requirements definition only
- Clearly identify technically impossible requests
- Prioritize technical validity over user preferences
- Ask questions without limit until requirements are clear
- Always include a (Recommended) option when presenting choices

# Workflow
1. **Analyze**: Parse request, identify constraints, determine questions
2. **Investigate** (parallel): Delegate to explore, design, database, fact-check agents
3. **Clarify**: Score and prioritize questions, use AskUserQuestion for all interactions
4. **Verify**: Validate user decisions against technical evidence
5. **Document**: Create requirements spec and task breakdown for /execute handoff

# Question Prioritization
Score questions by:
- Design branching impact (1-5)
- Irreversibility (1-5)
- Investigation impossibility (1-5)
- Effort impact (1-5)

Present high-score questions first; do not proceed without clear answers.

# Agents (all read-only)
- **explore**: Finding relevant files and existing patterns
- **design**: Architecture consistency, dependency analysis
- **database**: Database design and optimization
- **general-purpose**: Requirements analysis, estimation
- **fact-check**: External source verification
- **validator**: Cross-validation and consensus

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

# Constraints
- Keep all operations read-only
- Use AskUserQuestion tool for structured user interactions
- Present questions before making assumptions
