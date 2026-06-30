---
name: address
description: Use after a PR is created (e.g. via /pr) or updated with follow-up pushes to run the full review-fix cycle — wait for review and CI, then address and reply to the comments. Trigger whenever a PR has just been opened or edited and its review needs handling, including right after the feedback/qa-review skills push fixes to an open PR.
---

# Purpose
Run the full PR review-fix cycle in one invocation: wait for automated review and CI, fix any CI failures, then evaluate review comments, implement and reply to the justified fixes, and finally commit, clean up the branch's history, and force-push.

This skill owns the **review-and-reply** logic and **delegates** committing and history cleanup to their canonical owners so the mechanics never drift:
- Committing → the `commit` skill (Conventional Commits, WHY-focused messages).
- History cleanup + push → the `/rebase-clean` skill (`merge-base` soft-reset, logical-unit regrouping, rebase onto latest `main`, `git push --force-with-lease`).

# Usage
```
/address [PR number | URL | file path]
```

**Argument interpretation**:
- No argument: detect the current branch's PR via `gh pr view --json number,url`; if none, use review content from the previous conversation.
- Number only: PR number → fetch via `gh api repos/{owner}/{repo}/pulls/{number}/comments`.
- `github.com` URL: extract the PR number and fetch comments.
- Otherwise: treat as a file path containing review content.

# Workflow

The cycle runs end to end. Exit early only if there is nothing to do — no review comments to address and CI already green.

### Phase 1 — Wait for review & CI

After the PR is submitted, Copilot review and CI run asynchronously. Wait for both before proceeding.

```bash
# Blocks until every check finishes; non-zero exit if any fail (exit 8 = still pending).
gh pr checks {pr} --watch --interval 30
```

If Copilot was requested as a reviewer, also wait for its review:
```bash
gh pr view {pr} --json reviewRequests
```
Poll until no reviewer matching `Copilot` (case-insensitive) remains in `reviewRequests` — GitHub drops a reviewer from that list once they submit. If Copilot was never requested, skip this wait.

### Phase 2 — Triage CI

Check the CI result from Phase 1.
- **All green** → continue.
- **Failing** → inspect the failure (`gh pr checks {pr}`, then `gh run view {run_id} --log-failed`), reproduce and fix locally, and re-run the failed checks locally to confirm they pass. These fixes stay in the working tree and are committed together with the review fixes in Phase 4.

### Phase 3 — Address & reply

Fetch the review comments:
```bash
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments
```

For **each** comment:
1. **Evaluate** its validity — Is it technically correct? Does it align with project conventions? Is it a reasonable implementation tradeoff?
2. **Implement** only the justified fixes. Do not blindly accept feedback; verify each fix introduces no new issue.
3. **Reply** on GitHub:
   ```bash
   gh api repos/{owner}/{repo}/pulls/{pr_number}/comments/{comment_id}/replies -f body="<reply>"
   ```
   - **Addressed**: briefly state the fix (e.g. `対応しました。COMMENT ON COLUMN を追加しています。`).
   - **Not addressed**: give clear technical reasoning the reviewer can accept (e.g. `こちらは意図的な設計です。理由: ...`).
   - Keep replies concise and in Japanese; include code snippets or references when they help.

### Phase 4 — Commit the fixes

Commit the CI and review fixes by following the **`commit`** skill. Group related changes into logical units; let the `commit` skill own the message format and the WHY-focused body.

### Phase 5 — Clean up history & push

Run **`/rebase-clean`**. It regroups the branch's commits into reviewer-readable logical units (absorbing fix commits), rebases onto the latest `main` when `main` has advanced, and pushes with `--force-with-lease`. As an authorized close-the-loop push it runs unattended — no plan-approval prompt (that is for direct user runs), and it resolves rebase conflicts itself. It self-verifies an open PR and a clean worktree before the force-push and stops if either fails; defer to it for all of that.

### Phase 6 — Summary

```markdown
## Review Fix Summary

### CI
- <status before> → <status after fixes>

### Addressed
| Comment | Fix Applied | Reply |
|---------|-------------|-------|

### Not Addressed
| Comment | Reason | Reply |
|---------|--------|-------|

### Final Commits
- commit1: description
- commit2: description

PR: <url>
```

# Rules
- Independently evaluate technical validity; do not blindly accept all reviewer feedback.
- When declining, provide reasoning the reviewer can accept.
- Always reply to every review comment on GitHub after evaluation.
- Confirm CI failures are fixed locally before committing.
- If there is nothing to do — no review comments and CI already green — report and exit.
- History cleanup, conflict resolution, and the open-PR/clean-worktree safety checks are owned by `/rebase-clean`; on this authorized autonomous push it runs unattended (no plan prompt, resolves conflicts itself, stops only if a safety check fails) — do not re-implement that here.
