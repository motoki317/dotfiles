---
name: address
description: Run a PR's full review-fix cycle — wait for review and CI, fix failures, evaluate and reply to comments, then commit and push. Use whenever a PR was just opened or updated — the post-Ship phase of the Orchestrator loop.
---

# Purpose
Run the full PR review-fix cycle in one invocation. This skill owns the **review-and-reply** logic and **delegates** the mechanics to their canonical owners so they never drift: committing → the `commit` skill; history cleanup + push → `/rebase-clean`.

# Usage
```
/address [PR number | URL | file path]
```
- No argument: detect the current branch's PR via `gh pr view --json number,url`; if none, use review content from the previous conversation.
- PR number or `github.com` URL: fetch that PR's review via both endpoints in Phase 3.
- Otherwise: a file path containing review content.

# Workflow

Run the cycle end to end. Exit early only if there is nothing to do — no review feedback (neither a summary body nor inline comments) and CI already green.

### Phase 1 — Wait for review & CI

Wait for CI (status checks):
```bash
gh pr checks {pr} --watch --interval 30   # blocks until all finish; non-zero if any fail
```

Copilot's review is not a status check — wait for it separately, until a `copilot` review with a non-empty body exists; if none arrives within the timeout, proceed:
```bash
for i in $(seq 1 20); do
  n=$(gh api repos/{owner}/{repo}/pulls/{pr}/reviews \
        --jq '[.[] | select(.user.login | test("copilot";"i")) | select(.body|length>0)] | length')
  [ "$n" -gt 0 ] && break
  sleep 15
done
```

### Phase 2 — Triage CI
All green → continue. Failing → inspect (`gh pr checks {pr}`, then `gh run view {run_id} --log-failed`), reproduce and fix locally, and re-run the failed checks locally to confirm they pass. These fixes stay in the working tree and are committed with the review fixes in Phase 4.

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
1. **Evaluate** — technically correct? aligned with project conventions? a reasonable tradeoff? Don't blindly accept reviewer feedback.
2. **Implement** only the justified fixes — trivial ones directly, substantial ones as an Implementer brief (`/execute`); verify each introduces no new issue.
3. **Reply** on GitHub to every comment and summary body — concise, in Japanese, with code snippets or references when they help:
   - Inline comment → `gh api repos/{owner}/{repo}/pulls/{pr_number}/comments/{comment_id}/replies -f body="<reply>"`; summary body (no reply thread) → `gh pr comment {pr} --body "<reply>"`.
   - Addressed → state the fix (e.g. `対応しました。COMMENT ON COLUMN を追加しています。`); declined → reasoning the reviewer can accept (e.g. `こちらは意図的な設計です。理由: ...`).

### Phase 4 — Commit
Commit the CI and review fixes by following the **`commit`** skill, grouped into logical units.

### Phase 5 — Clean up history & push
Run **`/rebase-clean`**: it regroups commits, rebases onto the latest `main`, and pushes with `--force-with-lease`. On this authorized open-PR follow-up push it runs unattended and self-checks the open PR and clean worktree — defer to it; don't re-implement its checks here.

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

PR: <url>
```
