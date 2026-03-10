---
aliases: [address, address-review, ar]
description: Review PR comments, evaluate validity, and implement justified fixes
source: https://github.com/kawarimidoll/dotfiles/commit/070c374426ff44d1167d9a53c40a9161013d6cd8
---

# Purpose
Review PR comments, critically evaluate each suggestion, and implement only technically valid fixes.

# Usage
```
/address [PR number | URL | file path]
```

**Argument interpretation**:
- No argument: Use review content from previous conversation
- Number only: PR number → fetch via `gh api /repos/{owner}/{repo}/pulls/{number}/comments`
- URL with github.com: Extract PR number and fetch comments
- Otherwise: File path containing review content

# Workflow
1. Parse argument and fetch/load review comments
2. Evaluate each comment's validity:
   - Is it technically correct?
   - Does it align with project conventions?
   - Is it a reasonable implementation tradeoff?
3. Implement only justified fixes
4. Reply to each review comment on GitHub (see Reply section)
5. Output summary

# Reply to Review Comments
After evaluating and implementing fixes, reply to **each** review comment on GitHub via:
```
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments/{comment_id}/replies -f body="<reply>"
```

Reply guidelines:
- **Addressed**: Briefly explain the fix applied (e.g., "対応しました。`COMMENT ON COLUMN` を追加しています。")
- **Not addressed**: Explain why with clear technical reasoning (e.g., "こちらは意図的な設計です。理由: ...")
- Keep replies concise and in Japanese
- Include relevant code snippets or references when helpful

# Output Format
```markdown
## Review Response Summary

### All Comments
- [ ] Comment 1 summary
- [ ] Comment 2 summary

### Addressed
| Comment | Fix Applied | Reply |
|---------|-------------|-------|
| Comment 1 | Description of change | Reply sent |

### Not Addressed
| Comment | Reason | Reply |
|---------|--------|-------|
| Comment X | Why it was declined | Reply sent |
```

# Rules
- Do not blindly accept all reviewer feedback
- Independently evaluate technical validity
- When declining, provide clear reasoning the reviewer can accept
- Verify fixes don't introduce new issues
- Always reply to every review comment on GitHub after evaluation
