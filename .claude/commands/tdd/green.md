---
description: Complete the GREEN phase - Make the test pass with minimal code
---

# TDD GREEN Phase: Make the Test Pass

Make the failing test pass with **minimal code**.

## Rules
- This phase occurs exactly once per cycle
- Do NOT add unnecessary functionality
- Ignore code quality temporarily (that's for REFACTOR)

## Beck's Three Strategies

### 1. Fake It (Till You Make It)
**Use when**: Unsure how to implement

```python
# Test: assert add(1, 1) == 2
def add(a, b):
    return 2  # Just return the expected constant
```

### 2. Obvious Implementation
**Use when**: Solution is clear and you're confident

```python
def add(a, b):
    return a + b  # Obviously correct
```

**Warning**: If it fails, fall back to Fake It.

### 3. Triangulation
**Use when**: You faked it and need to generalize

```python
# First test passed with fake: return 2
# Add second test to force generalization:
def test_add_two_and_three():
    assert add(2, 3) == 5  # This breaks the fake!
```

## Strategy Selection

**Fake It when**:
- Implementation requires > 1 new function
- Algorithm unclear after reading test
- Multiple conditionals needed
- Obvious Implementation already failed

**Obvious Implementation when**:
- Single expression suffices
- Pattern matches existing codebase
- Direct translation of test

**Triangulation when**:
- Faked it and need to generalize
- Abstraction not yet clear

## Checklist
- [ ] Test passes (run all tests)
- [ ] Implementation is minimal
- [ ] Used appropriate strategy

## Commit
GREEN = SAFE. Run `/git:commit` to save this checkpoint.

## Next
After committing, proceed to `/tdd:refactor`.
