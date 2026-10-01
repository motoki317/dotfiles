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

Argument interpretation:
- No argument: Use `main` as the base branch.
- Branch name: Use the specified branch as the base.

# Workflow

### Step 1: Preflight Check

```bash
CURRENT=$(git branch --show-current)
BASE=${1:-main}
```

- Verify that the current branch is not `main`.
- Verify that there are commits ahead of the base: `git log --oneline $BASE..HEAD`
- If uncommitted changes exist, ask the user whether to commit them first.

### Step 2: Analyze Changes

Gather context in parallel:
- `git log --oneline $BASE..HEAD` — all commits in this branch
- `git diff --stat $BASE..HEAD` — files changed summary
- `git diff $BASE..HEAD` — full diff, to understand the changes

### Step 3: Generate PR Metadata

Title: write it in Conventional Commits style (`feat(scope): description`), in fewer than 70 characters.

Body: Use this template:
```markdown
## Motivation
Why this PR is needed — the problem, user pain, or business context driving the change.

## Changes
- Key changes grouped by logical unit

## Test plan
- [ ] How to verify the changes

@codex review
```

Then end the body with your standard attribution footer. Always include the `@codex review` line, which triggers Codex's automated review.

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
- When you update an existing PR body, never overwrite it from a local scratch file. The live body can contain screenshots or notes that the user added directly on GitHub. Fetch the live body first (`gh pr view <n> --json body -q .body`). Edit on top of it. Apply the result with `--body-file`. Then verify that the attachments are still in the body (`grep -c user-attachments`).
- Never push to `main` directly
- If `gh pr create` fails because a PR already exists, show the existing PR URL instead
