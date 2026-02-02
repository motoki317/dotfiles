---
description: Execute the REFACTOR phase - Improve code quality while keeping tests green
---

# TDD REFACTOR Phase: Improve the Code

Improve code quality while keeping all tests passing.

## Rules
- Only refactor when all tests pass (you should be GREEN)
- Run tests after EVERY small change
- This phase is repeatable - do multiple refactors as needed
- Commit after EACH successful refactor

## Primary Goal: Remove Duplication

Look for:
- Duplicated logic between test and production code
- Magic numbers/strings in multiple places
- Similar code patterns to extract
- Unclear naming that hides intent

## Allowed Changes (Structural Only)
- Renaming variables/methods
- Extracting functions/methods
- Moving code to better locations
- Improving formatting
- Simplifying expressions

## NOT Allowed (Save for Next RED)
- Adding new functionality
- Fixing bugs (write a failing test first)
- Changing behavior

## Workflow
```
1. Identify ONE small improvement
2. Make the change (< 20 lines)
3. Run ALL tests
4. Pass? → /git:commit with "refactor:"
5. Fail? → Revert immediately, take smaller step
6. More improvements? → Repeat
7. Satisfied? → /tdd:red for next cycle
```

## Step Size Limits
- Lines changed: < 20 per commit
- Files touched: 1 per commit (ideally)
- Scope: One rename, one extraction, or one move

## Common Patterns
| Pattern | When to Use |
|---------|-------------|
| Extract Method | Block does one identifiable thing |
| Rename | Name doesn't reveal intent |
| Inline | Abstraction adds no value |
| Extract Variable | Expression is complex or repeated |

## Checklist
- [ ] Change is purely structural
- [ ] Tests still pass
- [ ] Diff is < 20 lines
- [ ] Ready to commit

## Safety
- GREEN = Safe, can always revert to last commit
- If refactor breaks tests: `git checkout -- .`
