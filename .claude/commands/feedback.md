---
argument-hint: [previous-command]
description: Review command for Claude Code's recent work
---

# Purpose
Multi-faceted review of Claude Code's work within the same session. Auto-selects review mode based on previous command.

# Rules
- Launch all agents simultaneously (timeout avoidance)
- Auto-select mode based on previous command
- Review only changed code in execute mode, not existing issues
- Provide concrete fix proposals, not abstract theories
- Include specific file:line references in all feedback

# Modes

## After /define
**Target**: Execution plan from conversation
**Review**: Step granularity, dependencies, risk identification, completeness, feasibility
**Agents**: plan, estimation, fact-check

## After /execute
**Target**: Files modified via Edit/Write tools
**Agents**: quality, security, design, docs, performance, test, fact-check
**Focus**: Naming, DRY, readability, OWASP Top 10, architecture patterns, test coverage

## After /bug
**Target**: Investigation results from conversation
**Review**: Evidence collection, hypothesis validity, root cause accuracy, log utilization
**Metrics**: Confidence, Log Utilization, Objectivity
**Agents**: quality-assurance, general-purpose, explore, fact-check

## After /ask
**Target**: Answer and evidence from conversation
**Review**: Evidence citation quality, conclusion validity, reference accuracy
**Metrics**: Confidence, Evidence Coverage
**Agents**: explore, quality-assurance, code-quality, fact-check

## General (other)
**Target**: Recent Claude Code work
**Agents**: review, complexity, memory, fact-check

# Output Format
```
## Feedback Results (Mode: {mode})

### Evaluation Scores
- {Metric1}: XX/100
- {Metric2}: XX/100
- Overall: XX/100

### Critical (Immediate Fix Required)
- [Category] Issue: Location
- Problem: Description
- Fix: Proposal

### Warning (Fix Recommended)
- [Category] Issue: Location
- Recommendation: Proposal

### Good Practice
[Category] Commendable aspects

### Fact Check Results
- Verified claims: ...
- Flagged claims (confidence < 80): ...
- Unverifiable claims: ...

### Recommended Actions
- [High] Action
- [Medium] Action
- [Low] Action
```

# Constraints
- All agents must run in parallel (sequential execution causes timeout)
- Provide suggested improvements for each issue identified
- Never provide feedback without code analysis
