# REFACTOR Phase — Improve the Code

Goal: improve the structure of the code while every test stays green. The green bar is your
safety net: it lets you change shape freely because any regression shows up instantly.

## Rules

- Only refactor from a green bar. If tests aren't passing, you're still in GREEN.
- Run the whole suite after every change — refactoring is only safe while the net holds.
- Repeatable: do as many small refactors as the code warrants (RED and GREEN happen once each;
  this phase loops).
- Commit after each successful refactor, separately from the GREEN commit (see Tidy First below).

## Primary goal: remove duplication

Duplication is the clearest signal that structure is missing. Look for:

- Logic duplicated between test and production code
- Magic numbers or strings repeated across the code
- Similar patterns that want a shared abstraction
- Names that hide intent rather than reveal it

## Structural vs. behavioral — Tidy First

A refactor changes *structure*, never *behavior*. Keeping the two in separate commits is the
heart of Kent Beck's "Tidy First": a `refactor:` commit that a reader can trust touches no
behavior is far easier to review, and far safer to revert or `git bisect` through, than one
that smuggles a behavior change inside a rename.

**Allowed (structural only):**
- Renaming variables / methods
- Extracting functions / methods
- Moving code to a better location
- Improving formatting
- Simplifying expressions

**Not allowed here — these are the next RED:**
- Adding new functionality
- Fixing a bug (write a failing test for it first)
- Any change to observable behavior

## Workflow

```
1. Identify ONE small improvement
2. Make the change (< 20 lines)
3. Run ALL tests
4. Pass?  → commit with "refactor:"
5. Fail?  → revert immediately (git checkout -- .), take a smaller step
6. More improvements? → repeat
7. Satisfied? → start the next cycle at RED
```

## Step size limits

- Lines changed: < 20 per commit
- Files touched: 1 per commit, ideally
- Scope: one rename, one extraction, or one move — not several at once

Small steps mean a failed test points at a single change, so reverting costs you seconds, not a
session.

## Common patterns

| Pattern | When to use |
|---------|-------------|
| Extract Method | A block does one identifiable thing |
| Rename | A name doesn't reveal intent |
| Inline | An abstraction adds no value |
| Extract Variable | An expression is complex or repeated |

## Safety

Because GREEN was committed, the last good state is always one command away. If a refactor
breaks tests and the fix isn't obvious, don't debug the broken refactor — revert and take a
smaller step:

```
git checkout -- .
```

## Exit checklist

- [ ] Every change was purely structural
- [ ] Tests still green
- [ ] Each commit's diff is < 20 lines
- [ ] Committed (separate `refactor:` commits)

→ Next behavior starts a new cycle at RED (`$HOME/.claude/skills/tdd/references/red.md`).
