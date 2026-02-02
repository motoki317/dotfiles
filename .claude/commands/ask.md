---
argument-hint: [question]
description: Question and inquiry command
---

# Purpose
Provide accurate, evidence-based answers to project questions through fact-based investigation. Read-only mode; never modifies files.

# Rules
- **Never modify, create, or delete files**
- Base answers on factual investigation from code and documentation
- Report confidence levels and unclear points honestly
- Never justify user assumptions; prioritize technical accuracy
- Provide file:line references for all findings

# Workflow
1. **Analyze**: Understand the question, classify type (architecture, implementation, debugging, design)
2. **Investigate** (parallel): Delegate to explore, design, performance, fact-check agents
3. **Synthesize**: Compile findings with confidence metrics
4. **Self-evaluate**: Identify gaps, append feedback section

# Agents (all read-only)
- **explore**: Finding files, exploring codebase structure
- **design**: System design, architecture, API structure
- **performance**: Performance bottlenecks, optimization
- **quality-assurance**: Code quality evaluation
- **code-quality**: Code complexity analysis
- **fact-check**: External source verification

# Output Format
```
## Question
Restate the user's question

## Investigation
Evidence-based findings with file:line references
- Source 1: `path/to/file.ts:42` - finding
- Source 2: `path/to/other.ts:15` - finding

## Conclusion
Direct answer based on evidence

## Metrics
- Confidence: 0-100
- Evidence Coverage: 0-100

## Recommendations
Suggested actions (no implementation)

## Unclear Points
Information gaps that would improve the answer
```

# Constraints
- Keep all operations read-only
- Cite specific file:line references
- Distinguish between facts and inferences
- Never confirm user assumptions without verification
