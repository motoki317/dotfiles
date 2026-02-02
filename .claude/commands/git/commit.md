---
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git add:*), Bash(git commit:*), Bash(git restore:*), Bash(git show:*)
description: Stage meaningful diffs and create Conventional Commits with WHY-focused messages
model: sonnet
---

## Current Git Context

### Working Directory Status
!`git status --short`

### Unstaged Changes
!`git diff --stat`

### Staged Changes
!`git diff --cached --stat`

### Recent Commits
!`git log --oneline -10`

### Current Branch
!`git branch --show-current`

## Conventional Commits Format

```
<type>[scope]: <description>

[body: explain WHY]

[footer]
```

**Types**:
- `fix:` Bug fix (PATCH)
- `feat:` New feature (MINOR)
- `refactor:` Code restructuring, no behavior change
- `perf:` Performance improvement
- `test:` Adding/correcting tests
- `docs:` Documentation only
- `style:` Formatting, whitespace
- `build:` Build system, dependencies
- `ci:` CI configuration
- `chore:` Maintenance

**Breaking Changes**: Append `!` or add `BREAKING CHANGE:` footer

## Commit Message Philosophy

- **Code** describes HOW
- **Tests** describe WHAT
- **Commit log** describes **WHY**
- **Comments** describe WHY NOT

## Workflow

1. **Analyze**: Review changes, show detailed diffs if needed
2. **Plan**: Group related changes, separate unrelated ones
3. **Stage**: Use `git add -p` for partial staging when needed
4. **Type**: Determine commit type from staged changes
5. **Message**: Write imperative description (< 50 chars) + WHY in body
6. **Commit**: Execute

## Quality Checklist

- [ ] Staged changes are logically cohesive (single purpose)
- [ ] Commit type accurately reflects change nature
- [ ] Description is imperative, under 50 characters
- [ ] Body explains WHY, not just WHAT
- [ ] No unrelated changes included
- [ ] Breaking changes properly indicated

## Examples

```
feat(auth): add OAuth2 support for GitHub login

Users requested GitHub authentication to avoid creating another
account. OAuth2 chosen over OAuth1 for simpler flow and better
security with short-lived tokens.

Closes #142
```

```
fix(parser): handle empty input without panic

Parser assumed non-empty input, causing crashes in automated
pipelines where empty files occasionally appear. Added defensive
handling following the robustness principle.

Fixes #87
```
