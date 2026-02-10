---
name: commit
description: Stage meaningful diffs and create commits with WHY-focused messages. Use when agent needs to commit code changes.
---

# Git Commit

Use `/commit` to stage meaningful diffs and create commits with WHY-focused messages.

## Commit Discipline

Only commit when:
1. ALL tests are passing
2. ALL compiler/linter warnings resolved
3. Change represents a single logical unit
4. Message states whether structural or behavioral change

## Best Practices

- Small, frequent commits over large, infrequent ones
- WHY-focused messages explaining the reason for change
- Separate structural changes (refactor) from behavioral changes (feat/fix)
