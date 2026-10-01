---
name: qa-review
description: Exercise the running software with adversarial scenarios. Use after you change a feature, or on any request to QA, break, or verify behavior against spec.
---

# Purpose
Skeptical personas design adversarial scenarios in parallel (Phase A). The orchestrator executes them against the live target and reports defects with evidence (Phase B). Design fans out, but execution is serial — parallel mutation corrupts shared state. Stance: do not trust "it should work". (`/feedback` reasons statically about the diff, but this skill exercises the running software.)

# Step 1 — Scope
- **Target** — what to test: explicit arg (PR/path/feature) → else `git diff $(git merge-base HEAD <base>)..HEAD` + uncommitted → else ask.
- **Intent / acceptance criteria** — the basis for "correct", from the conversation, PBI/issue, or spec. If it is absent, ask. Judge every pass/fail against it. Record its source as each finding's **Test Basis**.
- **How to exercise it** — URL + test env, CLI, API + auth, DB access, or a combination of these. Drive a web/Electron UI with `agent-browser` (load `agent-browser skills get dogfood` first). If nothing runs, say so: findings become `by-inspection` / `could-not-verify`, never `executed`.

# Step 2 — Phase A: design scenarios (parallel)
Select the personas whose angle applies to the change. For a narrow fix, select a few. For a new user-facing surface, select most or all seven. Record skips in Coverage.

Spawn the selected personas as `reviewer` agents in one message. Each persona adapts its tactics to the surface and designs ≥1 concrete check. Each persona's expertise is in its rubric, not in the agent type. Give each only:
- the change + how to exercise it
- the intent (its Test Basis)
- its rubric, at the full path `$HOME/.claude/skills/qa-review/references/<slug>.md`

Each returns `{persona, title, priority, basis, preconditions, test_data, steps, expected, evidence_to_collect}` — design only, no fixes.

If the target has a UI, each persona explores it read-only via `agent-browser` to ground its design. Read-only access is parallel-safe, but mutation is not.

| # | Persona — angle | slug |
|---|-----------------|------|
| P1 | New user — careless: misclicks, empty submits, impatient retries | `p1-new-user` |
| P2 | Veteran operator — fast/bulk: keyboard, shortcuts, batch/concurrent | `p2-veteran-operator` |
| P3 | Malicious operator — boundary/invalid/out-of-permission, double-submit | `p3-malicious-operator` |
| P4 | Data-integrity auditor — verify the authoritative state, not the surface | `p4-data-integrity-auditor` |
| P5 | Migration specialist — legacy data: missing fields, format/encoding, counts | `p5-migration-specialist` |
| P6 | Regression guardian — did existing/peripheral behavior break? | `p6-regression-guardian` |
| P7 | Spec skeptic — reconcile primary spec against actual behavior | `p7-spec-skeptic` |

New feature → P3, P4 lead. Migration → P5, P7 lead (see Migration mode).

# Step 3 — Phase B: execute (serial)
Orchestrator owns one session and runs the consolidated suite:
- Dedup the scenarios across personas. Run the highest `priority` first. Run destructive scenarios (P3, P5) last, or on disposable data. If you cap the suite, disclose the cap in Coverage — silent truncation reads as full coverage.
- Operate the real surface. For a UI, use `agent-browser`. Otherwise, use the CLI, the API, or queries. P4 reads the authoritative state, not the screen.
- Run scenarios serially, but one scenario can fire concurrent sessions on purpose. That race / double-submit (P2/P3/P4) is the test.
- Tag each finding `executed` (reproduced, with evidence) / `by-inspection` (read, not run) / `could-not-verify`. Write unknown as unknown.
- Report defects. Do not fix them mid-run (a fix invalidates the execution). Never declare "passed" or release-approved — a human gates release.

# Step 4 — Report
Normalize the findings to `{severity, persona, area, expected vs actual, repro, evidence, status, basis}`. Report each defect once, under its aptest persona, with cross-notes. Group by severity, then persona. The report is the deliverable — fixes belong to the Implementer loop. A dynamic defect counts as fixed only when its scenario re-executes clean.

```markdown
## QA Review — <feature>
Intent: <one line>. Exercised via: <UI/CLI/static>.

Coverage:
- Personas: <…>
- Execution: <…>
- Skipped: <… — why>

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
Migration mode covers a change that migrates or backfills existing data, not a greenfield feature. The basis is the existing canonical spec, and "correct" means that the target behaves as the spec states. P5 and P7 lead. Still run the rest. Add a requirement-coverage view (testable / testable-pending / not-testable) and a config + design-pattern coverage list. Disclose which items you exercised.

---
Source — the seven persona rubrics in `references/` adapt Nexta's "7 QA personas" method: https://zenn.dev/nexta_/articles/be13a2395a5d2a
