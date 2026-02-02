---
name: tdd
description: Guide Test-Driven Development using Kent Beck's Red-Green-Refactor cycle. Use when writing tests, implementing features via TDD, or following plan.md test instructions.
---

# TDD (Test-Driven Development)

Kent Beck's Red-Green-Refactor cycle with Tidy First principles.

## Workflow

```
/tdd:red -> write failing test -> /tdd:green -> pass test -> /git:commit
                                                                  |
      <- next feature <- /tdd:red <- satisfied? <- /tdd:refactor (repeat)
```

## Phases

| Phase | Command | Action |
|-------|---------|--------|
| RED | `/tdd:red` | Write ONE small failing test |
| GREEN | `/tdd:green` | Make it pass with minimal code, commit |
| REFACTOR | `/tdd:refactor` | Improve structure (no behavior change), commit each step |

## Core Principles

- **One test at a time**: Each RED adds exactly ONE failing test
- **Minimal code**: GREEN writes just enough to pass
- **Never skip REFACTOR**: Complete all three phases
- **Tidy First**: Separate structural changes (refactor) from behavioral changes (feat/fix)
- **Small commits**: Commit after GREEN, commit after EACH refactor

## GREEN Phase Strategies

| Confidence | Strategy | Use When |
|------------|----------|----------|
| Low | **Fake It** | Return constant, generalize later |
| High | **Obvious Implementation** | Solution is clear |
| Generalizing | **Triangulation** | Add test to break a fake |

## Quality Standards

- Eliminate duplication between test and production code
- Express intent through clear naming
- Keep methods small and focused
- Run ALL tests after EVERY change
