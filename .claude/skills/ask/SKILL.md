---
name: ask
description: Answer questions about the project with file:line evidence from code and docs — read-only, never modifies files. Use when the user asks how or why something works, or weighs a "should we" question that needs an investigated answer rather than a change.
argument-hint: "[question]"
---

# Rules
- Never modify, create, or delete files.
- Answer only from factual investigation; mark inferences as inferences. Never confirm the user's assumptions without verification — technical accuracy over agreement.
- Cite `file:line` for every finding.

# Workflow
1. **Investigate** — read the relevant code and docs; inline reading suffices for a scoped question. Fan out 1–3 read-only agents (`Explore` for structure sweeps, `reviewer` for a lensed read, `general-purpose` for external sources) only when the surfaces are broad or independent.
2. **Answer** — the direct conclusion first, then the `file:line` evidence behind it, then what remains unknown or would sharpen the answer.

For a deliberative question ("should we…", a trade-off rather than a lookup), also weigh the options through the 3–5 relevant lenses in `$HOME/.claude/skills/decision-analysis/references/lens-catalog.md` and present the trade-off each option accepts.
