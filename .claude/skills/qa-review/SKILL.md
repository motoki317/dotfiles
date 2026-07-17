---
name: qa-review
description: QA-test a feature by exercising the running software with adversarial scenarios. Use after implementing or changing a feature, or when the user asks for QA, to "try to break" something, hunt for defects, or verify behaviour against spec.
---

# Purpose
Skeptical personas design adversarial scenarios in parallel (Phase A); the orchestrator executes them against the live target and reports defects with evidence (Phase B). Design fans out, but execution is serial — parallel mutation corrupts shared state. Stance: don't trust "it should work". (`/feedback` reasons statically about the diff; this exercises the running software.)

# Step 1 — Scope
- **Target** — what to test: explicit arg (PR/path/feature) → else `git diff $(git merge-base HEAD <base>)..HEAD` + uncommitted → else ask.
- **Intent / acceptance criteria** — the basis for "correct", from the conversation, PBI/issue, or spec; ask if absent. Every pass/fail is judged against it; record its source as each finding's **Test Basis**.
- **How to exercise it** — URL + test env, CLI, API + auth, and/or DB access. Drive a web/Electron UI with `agent-browser` (load `agent-browser skills get dogfood` first). If nothing runs, say so: findings become `by-inspection` / `could-not-verify`, never `executed`.

# Step 2 — Phase A: design scenarios (parallel)
Select the personas whose angle applies to the change — a narrow fix warrants a few, a new user-facing surface most or all seven; record skips in Coverage. Spawn the selected personas as `reviewer` agents in one message — tactics adapted to the surface, each leaving ≥1 concrete check; the rubric, not the agent type, carries each persona's expertise. Give each only: the change + how to exercise it, the intent (its Test Basis), and its rubric at the full path `$HOME/.claude/skills/qa-review/references/<slug>.md`. Each returns `{persona, title, priority, basis, preconditions, test_data, steps, expected, evidence_to_collect}` — design only, no fixes. With a UI, ground the design by exploring it read-only via `agent-browser` (read-only is parallel-safe; mutation is not).

| # | Persona — angle | slug |
|---|-----------------|------|
| P1 | New user — careless: misclicks, empty submits, impatient retries | `p1-new-user` |
| P2 | Veteran operator — fast/bulk: keyboard, shortcuts, batch/concurrent | `p2-veteran-operator` |
| P3 | Malicious operator — boundary/invalid/out-of-permission, double-submit | `p3-malicious-operator` |
| P4 | Data-integrity auditor — verify the authoritative state, not the surface | `p4-data-integrity-auditor` |
| P5 | Migration specialist — legacy data: missing fields, format/encoding, counts | `p5-migration-specialist` |
| P6 | Regression guardian — did existing/peripheral behaviour break? | `p6-regression-guardian` |
| P7 | Spec skeptic — reconcile primary spec against actual behaviour | `p7-spec-skeptic` |

New feature → P3, P4 lead. Migration → P5, P7 lead (see Migration mode).

# Step 3 — Phase B: execute (serial)
Orchestrator owns one session and runs the consolidated suite:
- Dedup across personas; highest `priority` first; destructive ones (P3, P5) last or on disposable data. If capped, disclose it in Coverage — silent truncation reads as full coverage.
- Operate the real surface: UI via `agent-browser`, else CLI / API / queries. P4 reads the authoritative store, not the screen.
- Serial across scenarios, but one scenario may fire concurrent sessions on purpose — that race / double-submit (P2/P3/P4) is the test.
- Tag each finding `executed` (reproduced, with evidence) / `by-inspection` (read, not run) / `could-not-verify`. Write unknown as unknown.
- Report defects; don't fix them mid-run (a fix would invalidate the execution). Never declare "passed" or release-approved — a human gates release.

# Step 4 — Report
Normalize to `{severity, persona, area, expected vs actual, repro, evidence, status, basis}`; one defect → report once under its aptest persona, cross-noting. Group by severity, then persona. The report is the deliverable — fixing belongs to the Implementer loop, and a dynamic defect counts as fixed only when its scenario re-executes clean.

```markdown
## QA Review — <feature> (intent: <one line>; exercised via: <UI/CLI/static>)
Coverage: personas <…>; execution <…>; skipped <… — why>.

### Blocker (must fix before release)
- [P<n>] <expected vs actual> — repro: <steps> — `<evidence>` — status: executed
  Basis: <criterion / spec heading>

### Major / Minor
- [P<n>] <expected vs actual> — repro: <steps> — status: <executed|by-inspection>

### Verified good
- [P<n>] <what was exercised and held up>

### Could not verify
- [P<n>] <what, and why — no test env, missing spec>
```
# Migration mode
Migrating or backfilling existing data, not a greenfield feature: the basis is the existing canonical spec, and "correct" = behaves as specified. P5 and P7 lead; still run the rest. Add a requirement-coverage view (testable / testable-pending / not-testable) and a config + design-pattern coverage list, disclosing which were exercised.

---
Source — the seven persona rubrics in `references/` adapt Nexta's "7 QA personas" method: https://zenn.dev/nexta_/articles/be13a2395a5d2a
