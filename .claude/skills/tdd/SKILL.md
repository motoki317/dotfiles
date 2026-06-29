---
name: tdd
description: Guide Test-Driven Development using Kent Beck's Red-Green-Refactor cycle. Use when writing tests, implementing features via TDD, or following plan.md test instructions.
---

# TDD (Test-Driven Development)

Kent Beck's Red-Green-Refactor cycle with Tidy First principles. This file is the orchestrator:
it owns the rules that span the whole cycle. Each phase has a detail file under
`$HOME/.claude/skills/tdd/references/` (full paths in the table below) — **read the phase file
when you enter that phase**, so the per-phase guidance is in front of you while you act on it.

## The cycle

```
        ┌──────────────────────────────────────────────────────┐
        ↓                                                        │
   RED ──→ GREEN ──→ commit ──→ REFACTOR ──→ commit each step ───┘
 write a   make it    (safe     improve         (loops until
 failing   pass,      state)    structure,       satisfied)
  test     minimally            no behavior
                                  change
```

## Phases

| Phase | Runs | Does | Detail |
|-------|------|------|--------|
| RED | once per cycle | Write ONE small failing test | `$HOME/.claude/skills/tdd/references/red.md` |
| GREEN | once per cycle | Make it pass with minimal code, then commit | `$HOME/.claude/skills/tdd/references/green.md` |
| REFACTOR | loops | Improve structure with no behavior change, commit each step | `$HOME/.claude/skills/tdd/references/refactor.md` |

Work the phases in order, reading each phase's detail file as you enter it. A full cycle is
RED → GREEN → REFACTOR; REFACTOR loops until the code is clean, then the next behavior starts a
new cycle back at RED.

## Cross-phase invariants

These hold across every phase — the phase files don't restate them:

- **One test at a time.** Each RED adds exactly one failing test. One behavior per test keeps
  every failure a single, unambiguous signal.
- **Run ALL tests after every change.** Not just the test you're working on — the suite is what
  tells you a change was safe.
- **Tidy First — separate structural from behavioral changes.** Behavior changes are `feat:` /
  `fix:` (committed at GREEN); structure changes are `refactor:` (committed during REFACTOR).
  Never mix them in one commit: a reviewer must be able to trust a `refactor:` touched no behavior.
- **Commit cadence.** No commit during RED (a red bar isn't safe). Commit at GREEN — that's the
  checkpoint you can always return to. Commit after *each* refactor.
- **Never skip REFACTOR.** Skipping it is how a passing test today becomes unmaintainable code
  tomorrow. A cycle isn't done at the green bar; it's done after the cleanup.

## GREEN strategies (summary)

Pick by how confident you are in the implementation. Full detail, code examples, and selection
criteria are in `$HOME/.claude/skills/tdd/references/green.md`.

| Confidence | Strategy | In one line |
|------------|----------|-------------|
| Low | **Fake It** | Return a constant, generalize on later cycles |
| High | **Obvious Implementation** | Type the real solution; fall back to Fake It if it fails |
| Generalizing | **Triangulation** | Add a test the fake can't satisfy, forcing real logic |
