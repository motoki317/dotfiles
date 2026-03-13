---
argument-hint: "[PR number | URL]"
description: End-to-end PR review cycle — address comments, commit fixes, rebase-clean, and force-push
allowed-tools: [Bash, Read, Edit, Write, Grep, Glob, Agent, AskUserQuestion]
---

# Purpose
Complete the full PR review-fix cycle in one invocation: fetch review comments, evaluate and implement fixes, reply to reviewers, clean up commit history, and force-push.

# Usage
```
/review-fix [PR number | URL]
```

**Argument interpretation**:
- No argument: Detect PR for current branch via `gh pr view --json number,url`
- Number: PR number
- URL: Extract PR number from URL

# Workflow

### Phase 1: Address Review Comments

Fetch and process review comments:

```bash
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments
```

For each comment:
1. **Evaluate** — Is it technically correct? Does it align with project conventions?
2. **Implement** justified fixes
3. **Reply** to each comment on GitHub:
   ```bash
   gh api repos/{owner}/{repo}/pulls/{pr_number}/comments/{comment_id}/replies -f body="<reply>"
   ```
   - Addressed: Briefly explain fix (in Japanese)
   - Not addressed: Explain reasoning (in Japanese)

### Phase 2: Commit Fixes

Stage and commit changes following Conventional Commits:
- Group related fixes into logical commits
- Message format: `fix(review): <description>`
- Append `Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>`

### Phase 3: Rebase Clean

Present commit reorganization plan to user for approval, then:

```bash
MERGE_BASE=$(git merge-base main HEAD)
git reset --soft $MERGE_BASE
git reset HEAD
```

Re-commit in logical units (layer order: proto → domain → infra → UI):
1. One commit per logical change unit
2. Fix commits absorbed into parent commits
3. Generated code included with its source commit

### Phase 4: Push

```bash
git push --force-with-lease
```

### Phase 5: Summary

Output:
```markdown
## Review Fix Summary

### Comments Addressed
| Comment | Fix | Reply |
|---------|-----|-------|

### Comments Declined
| Comment | Reason | Reply |
|---------|--------|-------|

### Final Commits
- commit1: description
- commit2: description

PR: <url>
```

# Rules
- Always ask user approval before rebase-clean (Phase 3)
- Use `--force-with-lease` never `--force`
- If no review comments exist, report and exit
- If rebase conflicts occur, stop and ask user
