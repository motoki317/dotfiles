---
name: address
description: Run the post-open PR loop — wait for review and CI, fix findings, reply, push. Use after you open or update a PR.
---

Target: the PR in the argument, else `gh pr view --json number`.

Run rounds. When a round with a push ends, start the next round on the new head. A round with nothing to push is the last one.

## 1. Wait

Read the head SHA, the draft state, and the base: `gh pr view <pr> --json headRefOid,isDraft,baseRefName`. Run each wait in the background:

- CI: `gh pr checks <pr> --watch --fail-fast --interval 30`.
- Codex: poll its latest status line.
  ```bash
  gh api repos/{owner}/{repo}/issues/<pr>/comments --paginate --jq '.[] | select(.user.login == "chatgpt-codex-connector[bot]") | .updated_at + " " + (.body | match(".*(Code Review\\*\\*|usage limits).*").string)' | sort | tail -1
  ```
  Codex is done when the line shows `Completed` or `Failed` with the first 7 characters of the head SHA. Codex is also done on a usage-limit line whose leading `updated_at` is later than the push of the head. If the line shows `Failed`, or does not show the head's SHA within 5 minutes, comment `@codex review` once per head. Then wait for that run.
- Copilot (inline comments as `Copilot`) is done when `gh api repos/{owner}/{repo}/pulls/<pr>/reviews --paginate --jq '.[] | select(.user.login == "copilot-pull-request-reviewer[bot]") | .commit_id'` prints the head SHA. Read its review rule with `gh api repos/{owner}/{repo}/rules/branches/<base> --jq '.[] | select(.type == "copilot_code_review") | .parameters'`. No output means no automatic review. Otherwise Copilot reviews when the PR opens or leaves draft. It reviews a draft only with `review_draft_pull_requests`, and a new push only with `review_on_push`. If Copilot will not review the head, do not wait for it. If it skips a draft, tell the user that marking the PR ready starts Copilot and needs another `/address`.

Go to step 2 only after all waits end.

## 2. Read

Skip the items that an earlier round addressed.
```bash
gh api repos/{owner}/{repo}/pulls/<pr>/reviews --paginate --jq '.[] | select(.body | length > 0) | {id, login: .user.login, body}'
gh api repos/{owner}/{repo}/pulls/<pr>/comments --paginate
gh api repos/{owner}/{repo}/issues/<pr>/comments --paginate
```
A bot's summary body or status comment is boilerplate. The real findings are the inline comments. For a failed check, read `gh run view <id> --log-failed`.

## 3. Fix and reply

For each actionable item:
1. Evaluate whether it is correct, fits the project conventions, and is a reasonable tradeoff.
2. Fix only the justified items. Fix config-only or single-file items yourself, and send other code through `/codex-work`. Reproduce a failed check locally, and run it again after the fix.

Commit your own fixes (`commit` skill), then run `/rebase-clean`. It rebases and force-pushes. If it stops, report and do not reply.

Otherwise, reply to each item concisely in Japanese, for example `対応しました。COMMENT ON COLUMNを追加しています。` For a declined item, give a reason the reviewer can accept. Reply to an inline comment with `gh api repos/{owner}/{repo}/pulls/<pr>/comments/<id>/replies -f body="<reply>"`, and to any other item with `gh pr comment <pr> --body "<reply>"`. Then resolve the thread of each fixed inline comment. Find the thread by the `databaseId` of its first comment:
```bash
gh api graphql -f query='{repository(owner:"<owner>",name:"<repo>"){pullRequest(number:<pr>){reviewThreads(first:100){nodes{id isResolved comments(first:1){nodes{databaseId}}}}}}}'
gh api graphql -f query='mutation{resolveReviewThread(input:{threadId:"<id>"}){thread{isResolved}}}'
```

## 4. Report

Report the CI state at the start and at the end, and the PR URL. Add one row per item: the fix or the reason, and the reply link.
