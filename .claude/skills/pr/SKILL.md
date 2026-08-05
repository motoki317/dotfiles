---
name: pr
description: Push the current branch and open a pull request derived from the commit history. Use when the user asks for a PR — the Ship step, always user-triggered.
argument-hint: "[base-branch]"
allowed-tools: [Bash, Read, Grep, Agent]
---

# Usage
```
/pr [base-branch]
```

**Argument interpretation**:
- No argument: Use `main` as base branch
- Branch name: Use specified branch as base

# Workflow

### Step 1: Preflight Check

```bash
CURRENT=$(git branch --show-current)
BASE=${1:-main}
```

- Verify current branch is not `main`
- Verify there are commits ahead of base: `git log --oneline $BASE..HEAD`
- Check for uncommitted changes — if present, ask user whether to commit first

### Step 2: Analyze Changes

Gather context in parallel:
- `git log --oneline $BASE..HEAD` — all commits in this branch
- `git diff --stat $BASE..HEAD` — files changed summary
- `git diff $BASE..HEAD` — full diff for understanding the changes

Identify:
- The primary change type (feat, fix, refactor, etc.)
- Scope of changes (which modules/areas affected)
- The "why" behind the changes (from commit messages and code)

### Step 3: Generate PR Metadata

**Title**: Under 70 chars, Conventional Commits style (`feat(scope): description`)

**Body**: Use this template:
```markdown
## Motivation
Why this PR is needed — the problem, user pain, or business context driving the change.

## Summary
- Bullet points explaining WHAT changed

## Changes
- Key changes grouped by logical unit

## Test plan
- [ ] How to verify the changes

@codex review
```

Then end the body with your standard attribution footer. Always include the `@codex review` line which triggers Codex's automated review

### Step 4: Push and Create PR

```bash
git push -u origin $CURRENT
gh pr create --base $BASE --title "..." --body "$(cat <<'EOF'
...
EOF
)"
```

### Step 5: Report

Output the PR URL.

# Rules
- Never push to `main` directly
- Always use `--base` flag to specify target branch
- If `gh pr create` fails due to existing PR, show the existing PR URL instead
- Keep title concise — details go in the body
- Derive content from commits, not speculation
