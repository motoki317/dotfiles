---
name: decision-analysis
description: Use *before* deciding, choosing, or committing to an approach — analyze a problem or decision through the few lenses that best expose its blind spots (feasibility, risk, impact, alternatives, cost, …), run in parallel by sub-agents, then synthesize a recommendation with explicit trade-offs. Trigger when the user weighs options, asks whether to do something, compares approaches, assesses feasibility/risk/impact, or asks to think a problem through. Not for reviewing an existing change (use feedback) or exercising running software (use qa-review).
---

# Purpose
Deliberate on an open decision before it is made. Decompose the problem into the 3–5 lenses that best expose its blind spots, investigate each in parallel, reconcile the findings, and recommend with the trade-offs made explicit. This runs *before* a choice is committed — `/feedback` reviews a change that already exists, and `/qa-review` exercises software that already runs.

# Step 1 — Frame the decision
State the decision in one line, and what a good answer delivers: a go/no-go, a ranked set of options, or a recommended approach. Capture the hard constraints (deadline, stack, budget) up front — they bound every lens.

# Step 2 — Select lenses
Pick the 3–5 lenses from `$HOME/.claude/skills/decision-analysis/references/lens-catalog.md` that would most reveal blind spots in *this* problem. The broader or more irreversible the decision, the more lenses apply. Record the ones you skip, so coverage is visible rather than silently narrow. Cap at five — past that the synthesis blurs into a survey.

# Step 3 — Investigate in parallel
Spawn one agent per lens, in a single message (sequential spawning risks timeout). Give each only: the framing from Step 1, its assigned lens, and the catalog path above so it can read that lens's rubric — a freshly spawned agent cannot resolve a path relative to this skill. Agents are read-only; this is analysis, not implementation. Use `Explore` for codebase-grounded lenses, `reviewer` otherwise.

# Step 4 — Synthesize
Reconcile across lenses: where they agree, where they conflict, and which unknowns block a confident call. One factor can surface under several lenses (a risky migration is Risk + Feasibility + Cost) — name it once, under the lens that owns it. Close with a recommendation, the trade-offs it accepts, and concrete next steps. Surface conflicts and critical unknowns instead of papering over them; a confident wrong call is the failure mode this guards against.

# Output
```markdown
## Decision Analysis — <decision>
Lenses: ran <…>; skipped <… — why>.

### Findings by lens
- **<Lens>**: <what it surfaced> — <evidence / `file:line`>

### Synthesis
Consensus: <…>. Conflicts: <…>. Critical unknowns: <…>.

### Recommendation
<the call>, accepting <trade-off>. Next: <concrete steps>.
```
