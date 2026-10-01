---
name: decision-analysis
description: Analyze an open decision through parallel lenses and synthesize a recommendation. Use before you decide, to weigh options, feasibility, risk, or impact. Do not use it to review an existing change.
---

# Step 1 — Frame the decision
State the decision in one line, and what a good answer delivers: a go/no-go, a ranked set of options, or a recommended approach. Capture the hard constraints (deadline, stack, budget) up front — they bound every lens.

# Step 2 — Select lenses
Pick the 3–5 lenses from `$HOME/.claude/skills/decision-analysis/references/lens-catalog.md` that expose the most blind spots in *this* problem. The broader or more irreversible the decision, the more lenses apply. Cap the count at five, because more lenses turn the synthesis into a survey. Record the skipped ones so coverage is visible.

# Step 3 — Investigate in parallel
Spawn one agent per lens, in a single message so they run concurrently. Give each only the Step 1 framing, its lens, and the catalog path above. A freshly spawned agent cannot resolve a path relative to this skill. The agents do read-only analysis, not implementation. Use `Explore` for a codebase-grounded lens, and `reviewer` for any other lens.

# Step 4 — Synthesize
Reconcile: where lenses agree, where they conflict, which unknowns block a confident call. One factor can surface under several lenses (a risky migration is Risk + Feasibility + Cost). Name it once, under the lens that owns it. Close with a recommendation, its accepted trade-offs, and concrete next steps. Surface conflicts and critical unknowns. Do not hide them: a confident wrong call is the failure mode that this rule guards against.

# Output
```markdown
## Decision Analysis — <decision>
Lenses run: <…>. Lenses skipped: <… — why>.

### Findings by lens
- **<Lens>**: <what it surfaced> — <evidence / `file:line`>

### Synthesis
Consensus: <…>. Conflicts: <…>. Critical unknowns: <…>.

### Recommendation
<the call>. Accepted trade-offs: <trade-offs>. Next: <concrete steps>.
```
