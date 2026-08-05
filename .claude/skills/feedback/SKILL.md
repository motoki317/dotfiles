---
name: feedback
description: Multi-viewpoint quality review (ISO/IEC 25010) of a change or plan, run by parallel agents. Use after implementing a change, to pressure-test a plan, or when the user asks for feedback, a quality check, or a critique of recent work.
---

# Purpose
Review recent work across the ISO/IEC 25010:2023 quality model: one agent per characteristic materially at stake, in parallel, then synthesize. Unlike `/code-review` (correctness bugs on the raw diff), this adapts its viewpoints to the change and covers the whole quality model.

# Step 1 — Target and intent
**Target**: explicit argument (PR number/URL, path, or `plan`) → else the branch's change since it diverged (`git diff $(git merge-base HEAD <base>)..HEAD` plus uncommitted edits) → else files edited this session → else ask.

**Intent** (what the change was meant to do): from the conversation, PR/issue, or ask. Pass it to every agent — they start with no conversation history, and Functional Suitability and Interaction Capability are judged against intent, not the code alone.

# Step 2 — Select viewpoints
Run the three core viewpoints for almost any change; add each conditional one when the change matches its "include when". Record skips in the Coverage line so gaps are visible. When unsure, include it — a clean review is cheap, a missed dimension is not.

**Core (always):**
- **Functional Suitability** — graded against the intent.
- **Maintainability** — the cost of the next change.
- **Reliability** — every code path has failure modes.

**Conditional — include when the change involves…**
- **Security** — untrusted or external input; authentication/authorization; secrets, tokens, or crypto; (de)serialization; file, SQL, shell, or network access; permission changes; PII or other sensitive data.
- **Performance Efficiency** — loops or recursion over large/unbounded data; DB queries or N+1; hot, batch, or real-time paths; heavy allocation or copying; caching; large payloads; anything latency- or throughput-sensitive.
- **Interaction Capability** — any user- or operator-facing surface: GUI, CLI flags/help, a public or library API signature, error messages, or log output meant to be read.
- **Compatibility** — external systems or APIs; shared data formats or schemas; protocol/version interop or backward compatibility; resources shared with co-running processes.
- **Flexibility** — configuration or feature flags; environment/platform portability; scaling; swappable dependencies or plugin points; install, upgrade, or migration paths.
- **Safety** — destructive or irreversible operations (deletes, overwrites, bulk mutations, schema migrations, money movement); data-loss potential; safety-critical domains.

When tests are central to the change, also spawn a second Maintainability agent scoped to the Testability sub-characteristic.

# Step 3 — Run in parallel
Launch one `reviewer` agent per selected viewpoint, all in one message (sequential spawning risks timeout). The rubric, not the agent type, carries the expertise; `reviewer` cannot edit and keeps Bash inspection-only, so review leaves the worktree untouched. Give each agent only its own slice:
1. the change (diff or file list),
2. the intent,
3. its rubric — the full path `$HOME/.claude/skills/feedback/references/<viewpoint>.md` (kebab-case, e.g. `security.md`); a freshly spawned agent can't resolve a path relative to this skill,
4. the ask: report findings as `{severity, file:line, problem, concrete fix}` plus any commendable aspect, judging only the change, not pre-existing unrelated issues.

# Step 4 — Synthesize
Normalize to `{severity, file:line, problem, fix, characteristic}`. One root issue can surface under several lenses (an unhandled error path is Reliability + Security + Functional Suitability) — report it once, under the aptest characteristic, cross-noting the others. Group by characteristic, rank by severity. A finding conditional on an unmeasured fact ("slow if the table is large") is unfinished — resolve the condition with a count, a timing, or an `EXPLAIN`, or rank the worst case. The report is the deliverable — fixing belongs to the Implementer loop (`~/.claude/rules/process.md`).

# Planning mode (target = `plan`)
Reviewing a plan, not code: apply the same viewpoints as forward-looking questions ("does the plan account for failure modes? security surface? scaling?"). Keep it to a few agents.

# Output
```markdown
## Quality Review — ISO/IEC 25010 (target: <what>, intent: <one line>)

Coverage: ran <viewpoints>; skipped <viewpoint — why>.

### Critical (fix before merge)
- [<Characteristic>/<sub>] <problem> — `file:line`
  Fix: <concrete proposal>

### Warning (fix recommended)
- [<Characteristic>/<sub>] <problem> — `file:line`
  Fix: <concrete proposal>

### Good
- [<Characteristic>] <commendable aspect>
```

---
Source — sub-characteristic definitions quoted in `references/` are from ISO/IEC 25010:2023 (SQuaRE — Product quality model): https://www.iso.org/standard/78176.html
