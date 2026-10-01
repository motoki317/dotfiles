---
name: commit
description: Stage and commit in logical units with WHY-focused messages. Use when you commit changes.
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git add:*), Bash(git commit:*), Bash(git restore:*), Bash(git show:*)
---

# Git Commit

## Context

!`git status --short`
!`git diff --stat`
!`git diff --cached --stat`
!`git log --oneline -10`
!`git branch --show-current`

## Conventional Commits

```
<type>[scope]: <imperative description, < 50 chars>

[body: explain WHY, not WHAT]

[footer]
```

Types: `feat:` (MINOR) | `fix:` (PATCH) | `refactor:` | `perf:` | `test:` | `docs:` | `style:` | `build:` | `ci:` | `chore:`
Breaking Changes: Append `!`, or add a `BREAKING CHANGE:` footer.

## Workflow

1. Analyze — Review the diffs. Make sure that the tests pass and that no warnings remain.
2. Group — Make one commit per logical unit. Separate structural (refactor) changes from behavioral (feat/fix) changes. To split one file across units, first read `git diff <file>`. Then answer `git add -p` from a pipe: `printf 'y\nn\n' | git add -p <file>`. A bare `git add -p` has no TTY here and stages nothing.
3. Write — Make the description imperative. Put the WHY and WHY-NOT in the body (code=HOW, tests=WHAT).
4. Commit

## Example

```
feat(auth): add OAuth2 support for GitHub login

Users requested GitHub authentication so that they do not need
another account. OAuth2 gives a simpler flow than OAuth1, and
its short-lived tokens give better security.

Closes #142
```
