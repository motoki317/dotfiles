---
argument-hint: [question]
description: Question and inquiry command
---

# Purpose
Answer project questions with evidence from the code and docs. Read-only; never modifies files.

# Rules
- Never modify, create, or delete files.
- Answer only from factual investigation; mark inferences as inferences, not facts.
- Cite `file:line` for every finding.
- Report confidence and remaining unknowns honestly; never confirm or justify the user's assumptions without verification — prioritize technical accuracy over agreement.

# Workflow
1. **Analyze** — classify the question (architecture, implementation, debugging, design).
2. **Investigate** in parallel — spawn the read-only agents that fit: `Explore` for files and structure, `reviewer` for a lensed read (architecture, performance, quality), `general-purpose` for external sources.
3. **Synthesize** — compile findings with `file:line` evidence and confidence metrics. For a deliberative question ("should we…", "which option…", a trade-off rather than a lookup), also weigh the options through the 3–5 relevant lenses in `$HOME/.claude/skills/decision-analysis/references/lens-catalog.md` and present the trade-offs; a factual lookup skips this.
4. **Self-evaluate** — flag gaps and unclear points.

# Output Format
```
## Question
Restate the user's question

## Investigation
Evidence-based findings with file:line references
- `path/to/file.ts:42` — finding

## Conclusion
Direct answer based on evidence

## Trade-offs (deliberative questions only)
Options weighed through the relevant lenses, with the trade-off each accepts

## Metrics
- Confidence: 0-100
- Evidence Coverage: 0-100

## Recommendations
Suggested actions (no implementation)

## Unclear Points
Information gaps that would improve the answer
```
