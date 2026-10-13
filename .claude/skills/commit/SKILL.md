---
name: commit
description: Stage and commit in logical units with WHY-focused messages. Use when you commit changes.
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git add:*), Bash(git commit:*), Bash(git restore:*), Bash(git show:*)
---

## Context

!`git status --short`
!`git diff --stat`
!`git diff --cached --stat`
!`git log --oneline -10`
!`git branch --show-current`

## Rules

- Make one commit per logical unit. Keep a refactor apart from a behavior change.
- To split one file across units, first read `git diff <file>`. Then answer `git add -p` from a pipe: `printf 'y\nn\n' | git add -p <file>`. A bare `git add -p` has no TTY here and stages nothing.
- Write the message in Conventional Commits: `type(scope): imperative summary`, under 50 characters. The body says WHY and WHY-NOT. Code says HOW, and tests say WHAT.
