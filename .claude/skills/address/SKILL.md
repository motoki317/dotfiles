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
- Number only: PR number → fetch its review via both endpoints in Phase 3.
- `github.com` URL: extract the PR number and fetch its review via both endpoints in Phase 3.
- Otherwise: treat as a file path containing review content.

# Workflow

The cycle runs end to end. Exit early only if there is nothing to do — no review feedback to address (neither a summary body nor inline comments) and CI already green.

### Phase 1 — Wait for review & CI

Wait for CI (status checks):
```bash
gh pr checks {pr} --watch --interval 30   # blocks until all finish; non-zero if any fail
```

Copilot's review is not a status check — wait for it separately, until a `copilot` review with a non-empty body exists:
```bash
for i in $(seq 1 20); do
  n=$(gh api repos/{owner}/{repo}/pulls/{pr}/reviews \
        --jq '[.[] | select(.user.login | test("copilot";"i")) | select(.body|length>0)] | length')
  [ "$n" -gt 0 ] && break
  sleep 15
done
```
If none arrives within the timeout, proceed.

### Phase 2 — Triage CI

Check the CI result from Phase 1.
- **All green** → continue.
- **Failing** → inspect the failure (`gh pr checks {pr}`, then `gh run view {run_id} --log-failed`), reproduce and fix locally, and re-run the failed checks locally to confirm they pass. These fixes stay in the working tree and are committed together with the review fixes in Phase 4.

### Phase 3 — Address & reply

Fetch feedback from **both** endpoints; skip items already addressed on a prior run:
```bash
# Summary bodies — empty-body reviews are inline-comment containers, skip them.
gh api repos/{owner}/{repo}/pulls/{pr_number}/reviews --paginate \
  --jq '.[] | select(.body|length>0) | {id, login: .user.login, body}'
# Inline diff comments.
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments --paginate
```

For **each** summary body and inline comment:
1. **Evaluate** — technically correct? aligned with project conventions? a reasonable tradeoff?
2. **Implement** only the justified fixes; verify each introduces no new issue.
3. **Reply**:
   - Inline comment → `gh api repos/{owner}/{repo}/pulls/{pr_number}/comments/{comment_id}/replies -f body="<reply>"`
   - Summary body (no reply thread) → `gh pr comment {pr} --body "<reply>"`
   - **Addressed**: state the fix (e.g. `対応しました。COMMENT ON COLUMN を追加しています。`).
   - **Not addressed**: give reasoning the reviewer can accept (e.g. `こちらは意図的な設計です。理由: ...`).
   - Concise, in Japanese; add code snippets or references when they help.

### Phase 4 — Commit the fixes

Commit the CI and review fixes by following the **`commit`** skill. Group related changes into logical units; let the `commit` skill own the message format and the WHY-focused body.

### Phase 5 — Clean up history & push

Run **`/rebase-clean`**. It regroups commits into logical units, rebases onto the latest `main`, and pushes with `--force-with-lease`. On this authorized close-the-loop push it runs unattended (no plan prompt, resolves conflicts itself) and self-checks an open PR and clean worktree before pushing — defer to it for all of that.

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
- Always reply to every review comment and to the summary body on GitHub after evaluation.
- Confirm CI failures are fixed locally before committing.
- If there is nothing to do — no summary body, no review comments, and CI already green — report and exit.
- History cleanup, conflict resolution, and the pre-push safety checks belong to `/rebase-clean` — do not re-implement them here.
