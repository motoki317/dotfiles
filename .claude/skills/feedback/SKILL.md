---
name: feedback
description: Use after implementing a change, or to pressure-test a plan, for a structured multi-viewpoint quality review (ISO/IEC 25010) run by parallel specialist agents. Trigger when the user asks for feedback, a quality check, or a critique of recent work.
---

# Purpose
Review recent work — by default the change just produced this session — for quality across the ISO/IEC 25010:2023 model. Spawn one agent per characteristic that is materially at stake, in parallel, then synthesize their findings. Unlike `/code-review` (correctness bugs on the raw diff), this adapts its viewpoints to the change and covers the whole quality model.

# Step 1 — Target and intent
**Target** (what to review): explicit argument (PR number/URL, path, or `plan`) → else the branch's change since it diverged: `git diff $(git merge-base HEAD <base>)..HEAD` plus any uncommitted edits → else files edited this session → else ask.

**Intent** (what the change was meant to do): take it from the conversation, PR/issue, or ask. Pass it to every agent — they start with no conversation history, and Functional Suitability and Interaction Capability are judged against intent, not the code alone.

# Step 2 — Select viewpoints and agents
Run the three **core** viewpoints for almost any change. Add each **conditional** viewpoint when the change matches its "include when"; the broader the change, the more apply. Skip the rest and record skips in the Coverage line, so gaps are visible rather than silent. Err toward including a viewpoint when unsure — a clean review is cheap, a missed dimension is not.

**Core (always):**
- **Functional Suitability** → `general-purpose` — graded against the intent.
- **Maintainability** → `code-quality`.
- **Reliability** → `quality-assurance` — every code path has failure modes.

**Conditional — include when the change involves…**
- **Security** → `security` — untrusted or external input; authentication/authorization; secrets, tokens, or crypto; (de)serialization; file, SQL, shell, or network access; access-control or permission changes; PII or other sensitive data.
- **Performance Efficiency** → `performance` — loops or recursion over large/unbounded data; DB queries or N+1; hot, batch, or real-time paths; heavy allocation or copying; caching; large payloads; anything latency- or throughput-sensitive.
- **Interaction Capability** → `general-purpose` — any user- or operator-facing surface: GUI, CLI flags/help, a public or library API signature, error messages, or log output meant to be read.
- **Compatibility** → `design` — integration with external systems or APIs; shared data formats or schemas; protocol/version interop or backward compatibility; resources shared with co-running processes.
- **Flexibility** → `design` — configuration or feature flags; environment/platform portability; scaling behaviour; swappable dependencies or plugin points; install, upgrade, or migration paths.
- **Safety** → `general-purpose` (or `security`) — destructive or irreversible operations (deletes, overwrites, bulk mutations, schema migrations, money movement); data-loss potential; irreversible external side effects; safety-critical domains.

When tests are central to the change, also spawn `test` for Maintainability's Testability sub-characteristic.

# Step 3 — Run in parallel
Launch all selected agents in one message (sequential spawning risks timeout). Give each agent only its own slice:
1. the change (diff or file list),
2. the intent,
3. its viewpoint rubric: `$HOME/.claude/skills/feedback/references/<viewpoint>.md` (kebab-case, e.g. `security.md`) — lists the sub-characteristics to review against and the official source. Pass this full path; a freshly spawned agent can't resolve one relative to the skill,
4. the ask: report findings as `{severity, file:line, problem, concrete fix}`; review only the change, not pre-existing unrelated issues.

# Step 4 — Synthesize
Normalize the agents' differing output shapes into `{severity, file:line, problem, fix, characteristic}`. One root issue can surface under several lenses (an unhandled error path is Reliability + Security + Functional Correctness) — report it once, under the most apt characteristic, cross-noting the others. Group by characteristic, rank by severity.

# Planning mode (target = `plan`)
Reviewing a plan, not code: apply the same viewpoints as forward-looking questions ("does the plan ensure completeness? account for failure modes? security surface? scaling?"). Keep it to a few agents — a secondary, lighter path.

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
