---
name: address
description: Run the post-open PR loop — wait for review and CI, fix findings, reply, push. Use after you open or update a PR.
---

Target: the PR in the argument, else `gh pr view --json number`. Repeat steps 1–3 until a round pushes nothing.

## 1. Wait

Read the head SHA, the draft state, and the base: `gh pr view <pr> --json headRefOid,isDraft,baseRefName`. If the head was pushed after the PR opened, comment `@codex review`.

Run in the background:

- CI: `gh pr checks <pr> --watch --fail-fast --interval 30`
- Reviews:
  ```bash
  sig() { for ep in issues/<pr>/comments pulls/<pr>/comments pulls/<pr>/reviews; do gh api "repos/{owner}/{repo}/$ep" --paginate --jq '.[] | "\(.id) \(.updated_at // .submitted_at)"'; done | sort; }
  start=$(sig); for i in $(seq 1 30); do sleep 30; [ "$(sig)" != "$start" ] && break; done
  ```

Each time the review loop ends, run the step 2 commands to check each bot. Rerun the loop until each bot is done:

- Codex: its "Codex Review Summary" comment shows `Completed` for the head SHA, or Codex posts a usage-limit notice. On `Failed`, comment `@codex review` once per head.
- Copilot: its review has the head SHA as `commit_id`. Wait for it only if `gh api repos/{owner}/{repo}/rules/branches/<base> --jq '.[] | select(.type == "copilot_code_review") | .parameters'` prints a rule. The rule skips a draft unless its `review_draft_pull_requests` is true, and a head pushed after the PR opened unless its `review_on_push` is true.

If the loop ends without a change, stop waiting for the bots that are not done. Go to step 2 when CI ends and no bot is left to wait for.

## 2. Read

```bash
gh api repos/{owner}/{repo}/pulls/<pr>/reviews --paginate --jq '.[] | select(.body | length > 0) | {id, login: .user.login, commit_id, body}'
gh api repos/{owner}/{repo}/pulls/<pr>/comments --paginate
gh api repos/{owner}/{repo}/issues/<pr>/comments --paginate
```
Skip bot summaries and the comments that you already replied to. For a failed check, read `gh run view <id> --log-failed`.

## 3. Fix and reply

Fix each failed check, and each comment that is correct, fits the project conventions, and is worth its cost. Reproduce a failed check locally before and after the fix. If you fixed anything, commit it (`/commit`) and run `/rebase-clean`, which force-pushes. If `/rebase-clean` stops, go to step 4 without replying.

Reply in Japanese to each comment with the fix (`対応しました。COMMENT ON COLUMNを追加しています。`) or a reason the reviewer can accept. Reply to an inline comment with `gh api repos/{owner}/{repo}/pulls/<pr>/comments/<id>/replies -f body="<reply>"`, and to any other comment with `gh pr comment <pr> --body "<reply>"`. Resolve the thread of each fixed inline comment. Find the thread by the `databaseId` of its first comment:
```bash
gh api graphql -f query='{repository(owner:"<owner>",name:"<repo>"){pullRequest(number:<pr>){reviewThreads(first:100){nodes{id isResolved comments(first:1){nodes{databaseId}}}}}}}'
gh api graphql -f query='mutation{resolveReviewThread(input:{threadId:"<thread id>"}){thread{isResolved}}}'
```

## 4. Report

Report the PR URL, the CI state at the start and at the end, and the bots that you stopped waiting for. Add one row per comment: the fix or the reason, and the reply link.
