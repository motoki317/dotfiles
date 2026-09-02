---
name: commit
description: Stage and commit in logical units with WHY-focused messages. Use when committing changes.
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

**Types**: `feat:` (MINOR) | `fix:` (PATCH) | `refactor:` | `perf:` | `test:` | `docs:` | `style:` | `build:` | `ci:` | `chore:`
**Breaking Changes**: Append `!` or add `BREAKING CHANGE:` footer

## Workflow

1. **Analyze** — Review diffs; ensure tests pass and warnings are resolved
2. **Group** — One commit per logical unit; separate structural (refactor) from behavioral (feat/fix) changes; use `git add -p` for partial staging
3. **Write** — Imperative description + WHY in body (code=HOW, tests=WHAT, commit=WHY, comments=WHY NOT)
4. **Commit**

## Example

```
feat(auth): add OAuth2 support for GitHub login

Users requested GitHub authentication to avoid creating another
account. OAuth2 chosen over OAuth1 for simpler flow and better
security with short-lived tokens.

Closes #142
```
