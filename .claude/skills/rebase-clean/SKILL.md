---
name: rebase-clean
description: Fold and regroup commits into reviewer-readable logical units, rebase them onto origin/main, and force-push to an open PR. The Implementer regroups only its current task's commits. Use to tidy trial-and-error history — the Tidy step.
allowed-tools: [Bash, Read]
---

# Rebase Clean

Where a step writes `<base>` or `<orig>`, substitute its literal SHA. When a step says to stop and report, the tidy failed: report the failure and the step where it happened.

## How to group commits
The unit is one **logical functional change**.

- Fold each fix-up or trial-and-error commit into the commit that it corrects.
- Put a feature and its tests in one commit. Put a bugfix and its regression test in one commit.
- Put a refactor in its own commit.
- Put generated or derived code in the commit of its source.
- Order the commits so that each dependency lands before the code that uses it.
- If one functional commit is too large to review, split it by technical layer (proto, domain, infra, UI).

## Steps

### 1. Choose the base
Use the first case that applies:
1. If you hold only the Implementer seat, `<base>` is your current task's Tidy base. Reread the task's entry in the plan file for it. If there is no plan file, take it from the brief. Skip steps 7 and 8. If neither names a Tidy base, do not tidy: stop and report.
2. Otherwise (the Orchestrator, or a session that holds both seats), run `git fetch origin main`. If origin has no main branch, stop and report. `<base>` is the output of `git merge-base origin/main HEAD`.

### 2. Check the preconditions
- If `git log <base>..HEAD` lists no commits, end the skill.
- If `git status --porcelain` prints anything, stop and report.

### 3. Survey
Record `<orig>`, the output of `git rev-parse HEAD`, and put it in your report.
```bash
git log --reverse --oneline <base>..HEAD
```
Read every commit with `git show <sha>`. Plan the new commits by the rules in "How to group commits".

### 4. Reset to the base
```bash
git reset -N <base>
```

### 5. Recommit
Make the planned commits in order:
- Write each message in the format of the `commit` skill. Do not run the tests before each commit: the worktree holds the final tree, not the tree of that commit.
- To split a file between commits, use the piped `git add -p` from the `commit` skill.
- `git add -p` cannot split a file that the commits added: its whole content is one hunk. To commit a file as it was at an original commit, run `git restore --source=<sha> <path>`, then stage it and commit. Afterward, run `git restore --source=<orig> <path>`. It writes the final version back, and deletes the file if `<orig>` has none.
- Carry over the `Co-Authored-By` trailers of the commits that each new commit replaces.
- If a planned commit contains exactly the changes of one original commit, reuse its message and author: `git commit -C <sha>`.

### 6. Verify
```bash
git status --porcelain
git diff <orig>
```
If either command prints anything, run `git reset --hard <orig>`, then stop and report. Otherwise, put the output of `git log --oneline <base>..HEAD` in your report.

### 7. Rebase onto origin/main
```bash
git rebase origin/main
```
If the rebase stops on a conflict, resolve it so that the result keeps the intent of both sides. If no resolution keeps both, run `git rebase --abort`, then stop and report. After the rebase completes, run the repo's checks. If they fail, do not push: leave the rebased branch in place, then stop and report.

### 8. Push
If `gh pr view --json state --jq .state` prints `OPEN`, run `git push --force-with-lease`. Otherwise, end the skill without a push.
