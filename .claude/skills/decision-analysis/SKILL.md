---
name: decision-analysis
description: Analyze an open decision through the lenses that best expose its blind spots (feasibility, risk, impact, alternatives, cost, …), run in parallel, then synthesize a recommendation with explicit trade-offs. Use *before* deciding — when the user weighs options, compares approaches, assesses feasibility/risk/impact, or thinks a problem through. Not for reviewing an existing change (feedback) or exercising running software (qa-review).
---

# Step 1 — Frame the decision
State the decision in one line, and what a good answer delivers: a go/no-go, a ranked set of options, or a recommended approach. Capture the hard constraints (deadline, stack, budget) up front — they bound every lens.

# Step 2 — Select lenses
Pick the 3–5 lenses from `$HOME/.claude/skills/decision-analysis/references/lens-catalog.md` that would most expose blind spots in *this* problem — the broader or more irreversible the decision, the more apply, but cap at five or the synthesis blurs into a survey. Record the skipped ones so coverage is visible.

# Step 3 — Investigate in parallel
Spawn one agent per lens, in a single message (sequential spawning risks timeout). Give each only the Step 1 framing, its lens, and the catalog path above (a freshly spawned agent can't resolve a path relative to this skill). Read-only analysis, not implementation: `Explore` for codebase-grounded lenses, `reviewer` otherwise.

# Step 4 — Synthesize
Reconcile: where lenses agree, where they conflict, which unknowns block a confident call. One factor can surface under several lenses (a risky migration is Risk + Feasibility + Cost) — name it once, under the lens that owns it. Close with a recommendation, the trade-offs it accepts, and concrete next steps; surface conflicts and critical unknowns rather than papering over them — a confident wrong call is the failure mode this guards against.

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
