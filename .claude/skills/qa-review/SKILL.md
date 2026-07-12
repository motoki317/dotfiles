---
name: qa-review
description: Use after implementing or changing a feature to actively QA-test it by exercising the running software. Trigger when the user asks for QA, to test or "try to break" a feature, hunt for defects, or verify behaviour against spec.
---

# Purpose
Actively QA-test a just-implemented feature. Seven skeptical personas design adversarial scenarios in parallel (Phase A), then the orchestrator executes them against the live target and reports defects with evidence (Phase B). Shared stance: *don't trust "it should work" — be maliciously skeptical.* Two phases because parallel mutation corrupts shared state, so design fans out but mutating execution is serialized. (`/feedback` reasons statically about the diff; this exercises the running software.)

# Step 1 — Scope
- **Target** — what to test: explicit arg (PR/path/feature) → else `git diff $(git merge-base HEAD <base>)..HEAD` + uncommitted → else ask.
- **Intent / acceptance criteria** — the basis for "correct," from the conversation, PBI/issue, or spec; ask if absent. Every pass/fail is judged against this; record its source as each finding's **Test Basis**.
- **How to exercise it** — URL + test env, CLI, API + auth, and/or DB access. Drive a web/Electron UI with `agent-browser` (load `agent-browser skills get dogfood` first). If nothing runs, say so: findings become `by-inspection` / `could-not-verify`, never `executed`.

# Step 2 — Phase A: design scenarios (parallel)
Spawn all seven as read-only `reviewer` agents in one message — a fixed checklist, each leaving ≥1 concrete check, tactics adapted to the surface; the rubric, not the agent type, carries each persona's expertise. Give each only: the change + how to exercise it, the intent (its Test Basis), and its rubric at the full path `$HOME/.claude/skills/qa-review/references/<slug>.md`. Each returns `{persona, title, priority, basis, preconditions, test_data, steps, expected, evidence_to_collect}` — design only, no fixes. With a UI, explore it read-only via `agent-browser` to ground the design (read-only is parallel-safe; mutation is not).

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
- Operate the real surface: drive a UI with `agent-browser`, else CLI / API / queries. P4 reads the authoritative store, not the screen.
- Serial across scenarios, but one scenario may fire concurrent sessions on purpose — that race / double-submit (P2/P3/P4) is the test.
- Tag each finding `executed` (reproduced, with evidence) / `by-inspection` (read, not run) / `could-not-verify`. Write unknown as unknown.
- Report defects here; don't fix them mid-run (a fix would invalidate the execution). Remediation, if any, is Step 5 and direct-invocation only. Never declare "passed" or release-approved — a human gates release.

# Step 4 — Report
Normalize to `{severity, persona, area, expected vs actual, repro, evidence, status, basis}`; one defect → report once under its aptest persona, cross-noting. Group by severity, then persona.

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

# Step 5 — Close the loop
<!-- Keep this phase in sync with feedback's Step 5: same fix → re-verify → commit → push-if-open-PR shape. qa-review adds the executed-only gate below; feedback omits it (static findings are always re-checkable). -->
Close the loop by default once the report is done. Stop at the report instead — and let the caller remediate — only when the caller says "analysis only" (e.g. `/execute`'s Verify owns its own fix loop and says so). The signal is the caller's word, not a guess about who invoked you.

Be more cautious than `/feedback`: a dynamic defect is fixed by changing behaviour, so close the loop **only for findings you `executed` and can re-execute**. If nothing runs — findings are `by-inspection` / `could-not-verify` — report and stop; never auto-fix and push a change you cannot re-exercise to confirm. Even a clean re-run never means "release-approved" — a human still gates that.

1. **Fix the justified Blockers** (and any Major worth fixing) — re-judge each for validity first, then fix with scoped sub-agents. Never route the fix through `/execute`: its Verify re-invokes this skill, which would loop. Leave declined items in the report with a reason.
2. **Re-exercise, bounded.** Re-run the affected scenarios against the live target — re-execution, not re-inspection — to confirm the fix resolved the defect without regression. If the fix introduces a new Blocker, do at most one more fix → re-exercise pass, then stop and report what remains.
3. **Commit** the fixes by following the `commit` skill. Never on `main` or a detached HEAD; if the worktree holds unrelated edits, stage only the fix files.
4. **Push — only to an open PR, from a clean worktree.** If `gh pr view --json number,state` shows an open PR and no unrelated uncommitted edits remain, clean up history and push via `/rebase-clean` (it rebuilds commits from the whole worktree after a soft reset, so stray edits would be swept in; force-pushed with `--force-with-lease`), then point the agent at `/address` for the review and CI the push re-triggers. If unrelated edits remain, stop after the commit and report that the PR push needs a clean or stashed worktree. With no open PR, stop after the commit — first publishing a branch and opening a PR stay user-triggered (`~/.claude/rules/process.md`). Reporting a defect fixed is not declaring the feature released.

# Migration mode
Migrating or backfilling existing data, not a greenfield feature: the basis is the existing canonical spec (not new requirements), and "correct" = behaves as specified. P5 and P7 lead; still run the rest. Add a requirement-coverage view (testable / testable-pending / not-testable) and a config + design-pattern coverage list, disclosing which were exercised.

---
Source — the seven persona rubrics in `references/` adapt Nexta's "7 QA personas" method: https://zenn.dev/nexta_/articles/be13a2395a5d2a
