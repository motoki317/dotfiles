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
4. Output summary

# Output Format
```markdown
## Review Response Summary

### All Comments
- [ ] Comment 1 summary
- [ ] Comment 2 summary

### Addressed
| Comment | Fix Applied |
|---------|-------------|
| Comment 1 | Description of change |

### Not Addressed
| Comment | Reason |
|---------|--------|
| Comment X | Why it was declined |
```

# Rules
- Do not blindly accept all reviewer feedback
- Independently evaluate technical validity
- When declining, provide clear reasoning the reviewer can accept
- Verify fixes don't introduce new issues
