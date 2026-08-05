---
name: address
description: Run the post-open PR loop — wait for review and CI, fix findings, reply, push. Use after opening or updating a PR.
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

Automated reviews (Copilot, Codex, any other review bot) are not status checks — wait for them separately. They fire on PR open / ready-for-review and post asynchronously, so wait until the set of bot reviews stops growing rather than for one specific bot; a slower reviewer is then never missed. If none arrives within the timeout, proceed:
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
All green → continue. Failing → inspect (`gh pr checks {pr}`, then `gh run view {run_id} --log-failed`), reproduce and fix locally, and re-run the failed checks locally to confirm they pass. These fixes stay in the working tree and are committed with the review fixes in Phase 4.

### Phase 3 — Address & reply

Fetch feedback from **both** endpoints. These queries are login-agnostic, so they already cover every reviewer — Copilot, Codex, and humans alike; skip items already addressed on a prior run:
```bash
# Summary bodies — empty-body reviews are inline-comment containers, skip them.
gh api repos/{owner}/{repo}/pulls/{pr_number}/reviews --paginate \
  --jq '.[] | select(.body|length>0) | {id, login: .user.login, body}'
# Inline diff comments — where bots put most of their concrete findings.
gh api repos/{owner}/{repo}/pulls/{pr_number}/comments --paginate
```

A bot's summary body is mostly boilerplate — Codex's "here are some suggestions" wrapper, Copilot's overview and per-file table, or a no-op notice (`Copilot was unable to review … quota limit`). Treat those as non-actionable; the real findings are the inline comments. Reply only to items that carry an actual request or question.

For **each** actionable summary body and inline comment:
1. **Evaluate** — technically correct? aligned with project conventions? a reasonable tradeoff? Don't blindly accept reviewer feedback — bot or human.
2. **Implement** only the justified fixes — trivial ones directly, substantial ones as an Implementer brief (`/codex-work`); verify each introduces no new issue.
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
