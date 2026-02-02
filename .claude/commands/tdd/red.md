---
description: Start the RED phase - Write a small, focused failing test
---

# TDD RED Phase: Write a Failing Test

Write ONE small, focused failing test.

## Rules
- **No commit** during RED phase (not in safe state)
- This phase occurs exactly once per cycle
- Write ONE failing test or add ONE failing assertion

## Writing Your Test

**Test naming**: Describe behavior, not implementation
- Good: `shouldAuthenticateValidUser`, `test_empty_input_returns_error`
- Bad: `testAuth`, `test1`

**Keep it small**: Test ONE specific behavior. If thinking "and also...", that's a second test.

**Must fail for the right reason**: Missing implementation (expected), NOT syntax errors or typos.

## Strategy Preview (for GREEN phase)

| Confidence | Strategy | Description |
|------------|----------|-------------|
| Uncertain | Fake It | Return a constant to pass |
| Confident | Obvious Implementation | Type the real solution |
| Generalizing | Triangulation | Add test to break a fake |

## When to Split Tests
- Testing different behavior (not just another case)
- New assertion requires different setup
- Unclear which assertion would fail first

## Checklist
- [ ] Descriptive name explaining expected behavior
- [ ] Small and focused on ONE thing
- [ ] Test FAILS (verify by running)
- [ ] Fails for the RIGHT reason

## Next
Once test fails correctly, proceed to `/tdd:green`.
