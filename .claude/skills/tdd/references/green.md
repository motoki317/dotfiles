# GREEN Phase — Make the Test Pass

Goal: make the one failing test from RED pass with the **minimal** code that does so.

Quality is *deliberately* deferred to REFACTOR. Separating "make it work" from "make it
good" is what keeps each step small and reversible: GREEN proves the behavior exists, REFACTOR
improves its shape without re-asking whether it works. Mixing the two is how a green bar turns
into an unreviewable diff.

## Rules

- Occurs exactly once per cycle (RED → **GREEN** → REFACTOR).
- Add nothing the failing test does not demand. Unbuilt functionality has no test, so it is
  unproven — and it is the next RED, not this GREEN.
- Tolerate ugly code for now: duplication, poor names, hard-coded values are all fine here.
  REFACTOR exists precisely so you don't have to get the design right under a red bar.

## Beck's Three Strategies

The strategy is chosen *here*, by how confident you are in the implementation — not in RED.

### 1. Fake It (Till You Make It)

Return a constant that satisfies the test, then generalize on later cycles.

```python
# Test: assert add(1, 1) == 2
def add(a, b):
    return 2  # the literal the test expects
```

Use when the path to a real implementation is unclear. A passing test — even faked — turns an
open design question into a concrete starting point you can triangulate away from.

### 2. Obvious Implementation

Type the real solution directly.

```python
def add(a, b):
    return a + b
```

Use when the solution is clear and small enough to get right in one go. **If it fails, fall
back to Fake It** — an unexpected red bar means your "obvious" was a guess, so shrink the step.

### 3. Triangulation

When you faked it and need to force generality, add a second test the fake cannot satisfy.

```python
# add() currently fakes `return 2`.
def test_add_two_and_three():
    assert add(2, 3) == 5  # the constant 2 can no longer pass — forces real arithmetic
```

Use when the right abstraction isn't yet visible. Two concrete examples pin down a
generalization that one example left ambiguous.

## Choosing Among Them

| Reach for… | When |
|------------|------|
| **Fake It** | Implementation needs more than one new function, the algorithm is unclear after reading the test, several conditionals are involved, or Obvious already failed |
| **Obvious Implementation** | A single expression suffices, the pattern matches existing code, or it's a direct translation of the test |
| **Triangulation** | You faked it and now need to generalize, or the abstraction is still unclear |

## Commit

A passing test is a **safe state** — the one point in the cycle you can always return to.
Commit now (`feat:` / `fix:`), before touching structure in REFACTOR. See the commit cadence
and Tidy First rationale in `$HOME/.claude/skills/tdd/SKILL.md`.

## Exit checklist

- [ ] All tests pass (run the whole suite, not just the new one)
- [ ] No code beyond what the test required
- [ ] Strategy was a deliberate choice, not a reflex
- [ ] Committed — the safe checkpoint exists

→ Proceed to REFACTOR (`$HOME/.claude/skills/tdd/references/refactor.md`).
