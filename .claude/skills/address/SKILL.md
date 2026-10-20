---
name: address
description: Run the post-open PR loop — wait for review and CI, fix findings, reply, push. Use after you open or update a PR.
---

Target: the PR in the argument, else `gh pr view --json number`.

Run rounds. When a round with a push ends, start the next round on the new head. A round with nothing to push is the last one.

## 1. Wait

Read the head SHA, the draft state, and the base: `gh pr view <pr> --json headRefOid,isDraft,baseRefName`. Codex reviews only when the PR opens or leaves draft, or on an `@codex review` comment. If you pushed to the PR after it opened, comment `@codex review`.

Run two waits in the background:

- CI: `gh pr checks <pr> --watch --fail-fast --interval 30`.
- Reviews: this loop stops when a comment or review on the PR is added or edited, or after 15 minutes.
  ```bash
  sig() { for ep in issues/<pr>/comments pulls/<pr>/comments pulls/<pr>/reviews; do gh api "repos/{owner}/{repo}/$ep" --paginate --jq '.[] | "\(.id) \(.updated_at // .submitted_at)"'; done | sort; }
  start=$(sig); for i in $(seq 1 30); do sleep 30; [ "$(sig)" != "$start" ] && break; done
  ```

Each time the review loop stops, read the comments with the commands in step 2. Then decide which bots are done:

- Codex (`chatgpt-codex-connector[bot]`) edits one "Codex Review Summary" comment in place. Codex is done when that comment shows `Completed` for the head SHA, or when Codex posts a usage-limit notice. On the first `Failed` for a head, comment `@codex review` again. On a second `Failed`, Codex is done.
- Copilot (inline comments as `Copilot`) is done when `gh api repos/{owner}/{repo}/pulls/<pr>/reviews --paginate --jq '.[] | select(.user.login == "copilot-pull-request-reviewer[bot]") | .commit_id'` prints the head SHA. It reviews only if `gh api repos/{owner}/{repo}/rules/branches/<base> --jq '.[] | select(.type == "copilot_code_review") | .parameters'` prints a rule. With a rule, it reviews when the PR opens or leaves draft. It reviews a draft only with `review_draft_pull_requests`, and a new push only with `review_on_push`. If it skips a draft, tell the user that marking the PR ready starts Copilot and needs another `/address`.

If a bot that reviews the head is not done, run the review loop again. If the loop ran the full 15 minutes, stop waiting for that bot and name it in the step 4 report. Go to step 2 after CI ends and each bot is done or reported.

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
