---
name: address
description: Run the post-open PR loop — wait for review and CI, fix findings, reply, push. Use after you open or update a PR.
---

# Purpose
Run the full PR review-fix cycle in one invocation. This skill owns the **review-and-reply** logic and **delegates** the mechanics to their canonical owners so they never drift. The `commit` skill owns commits, and `/rebase-clean` owns history cleanup and the push.

# Usage
```
/address [PR number | URL | file path]
```
- No argument: detect the current branch's PR via `gh pr view --json number,url`. If the branch has no PR, use review content from the previous conversation.
- PR number or `github.com` URL: fetch that PR's review via both endpoints in Phase 3.
- Otherwise: the path to a file that contains review content.

# Workflow

Run the cycle end to end. Exit early only if CI is already green and no review feedback exists (neither a summary body nor inline comments).

### Phase 1 — Wait for review & CI

Wait for CI (status checks):
```bash
gh pr checks {pr} --watch --interval 30   # blocks until all finish; non-zero if any fail
```

Automated reviews (Copilot, Codex, any other review bot) are not status checks. Wait for them separately. They fire on PR open / ready-for-review and post asynchronously, so wait until the set of bot reviews no longer grows. Do not wait for one specific bot, or you miss a slower reviewer. If none arrives within the timeout, proceed:
```bash
# Review bots configured for the repo, matched case-insensitively against the reviewer login:
#   copilot → copilot-pull-request-reviewer[bot] / Copilot;  codex → chatgpt-codex-connector[bot].
# Extend this pattern as you wire up more reviewers.
BOTS='copilot|codex'
prev=-1
for i in $(seq 1 20); do
  n=$(gh api repos/{owner}/{repo}/pulls/{pr}/reviews --paginate \
        --jq "[.[] | select(.user.login | test(\"$BOTS\";\"i\")) | select(.body|length>0) | .user.login] | unique | length")
  [ "$n" -gt 0 ] && [ "$n" -eq "$prev" ] && break   # ≥1 bot review in, and none new since the last poll
  prev=$n
  sleep 15
done
```

### Phase 2 — Triage CI
If all checks are green, continue. If a check fails, inspect it with `gh pr checks {pr}`, then `gh run view {run_id} --log-failed`. Reproduce and fix the failure locally. Re-run the failed checks locally to verify that they pass. Leave these fixes in the working tree. Commit them with the review fixes in Phase 4.

### Phase 3 — Address & reply

Fetch feedback from **both** endpoints. These queries are login-agnostic, so they already cover every reviewer — Copilot, Codex, and humans alike. Skip items that a prior run already addressed.
```bash
# Summary bodies — empty-body reviews are inline-comment containers, skip them.
gh api repos/{owner}/{repo}/pulls/{pr_number}/reviews --paginate \
  --jq '.[] | select(.body|length>0) | {id, login: .user.login, body}'
# Inline diff comments — where bots put most of their concrete findings.
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments --paginate
```

A bot's summary body is mostly boilerplate — Codex's "here are some suggestions" wrapper, Copilot's overview and per-file table, or a no-op notice (`Copilot was unable to review … quota limit`). Treat those as non-actionable. The real findings are the inline comments. Reply only to items that carry an actual request or question.

For **each** actionable summary body and inline comment:
1. Evaluate — decide whether the feedback is technically correct, follows the project conventions, and is a reasonable tradeoff. Do not blindly accept reviewer feedback — bot or human.
2. Implement only the justified fixes — trivial ones directly, substantial ones as an Implementer brief (`/codex-work`). Verify that each fix introduces no new issue.
3. Reply on GitHub to every comment and summary body. Keep each reply concise and in Japanese. If code snippets or references help, include them:
   - Inline comment → `gh api repos/{owner}/{repo}/pulls/{pr_number}/comments/{comment_id}/replies -f body="<reply>"`. Summary body (no reply thread) → `gh pr comment {pr} --body "<reply>"`.
   - Addressed → state the fix (for example, `対応しました。COMMENT ON COLUMNを追加しています。`). Declined → reasoning the reviewer can accept (for example, `こちらは意図的な設計です。理由：...`).

### Phase 4 — Commit
Follow the **`commit`** skill to commit the CI and review fixes in logical units.

### Phase 5 — Clean up history & push
Run **`/rebase-clean`**: it regroups commits, rebases onto `origin/main`, and pushes with `--force-with-lease`. On this push, an authorized follow-up to an open PR, it runs unattended and self-checks the open PR and the clean worktree. Defer to it: do not re-implement its checks here.

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
