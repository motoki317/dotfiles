---
name: tdd
description: Guide Test-Driven Development using Kent Beck's Red-Green-Refactor cycle. Use when writing tests, implementing features via TDD, or following plan.md test instructions.
---

# TDD (Test-Driven Development)

Kent Beck's Red-Green-Refactor cycle with Tidy First principles. This file owns the rules that span the whole cycle; **read each phase's detail file as you enter that phase**, so the per-phase guidance is in front of you while you act on it.

| Phase | Runs | Does | Detail |
|-------|------|------|--------|
| RED | once per cycle | Write ONE small failing test | `$HOME/.claude/skills/tdd/references/red.md` |
| GREEN | once per cycle | Make it pass with minimal code, then commit | `$HOME/.claude/skills/tdd/references/green.md` |
| REFACTOR | loops | Improve structure with no behavior change, commit each step | `$HOME/.claude/skills/tdd/references/refactor.md` |

REFACTOR loops until the code is clean; then the next behavior starts a new cycle back at RED.

## Cross-phase invariants

These hold across every phase — the phase files don't restate them:

- **One test at a time.** Each RED adds exactly one failing test; one behavior per test keeps every failure a single, unambiguous signal.
- **Run ALL tests after every change.** The suite, not just the test at hand, is what tells you a change was safe.
- **Tidy First — separate structural from behavioral changes.** Behavior is `feat:`/`fix:` (committed at GREEN); structure is `refactor:` (committed during REFACTOR). Never mix them in one commit: a reviewer must be able to trust a `refactor:` touched no behavior.
- **Commit cadence.** Never during RED (a red bar isn't safe); at GREEN — the checkpoint you can always return to; after *each* refactor step.
- **Never skip REFACTOR.** A cycle isn't done at the green bar; it's done after the cleanup — skipping it is how a passing test today becomes unmaintainable code tomorrow.
