---
name: rebase-clean
description: Fold and regroup commits into reviewer-readable logical units, rebase them onto the target branch, and force-push to an open PR. Use to tidy trial-and-error history.
allowed-tools: [Bash, Read]
---

Where a step writes `<base>` or `<orig>`, substitute its literal SHA. `<target>` is the branch in the argument, else the open PR's `baseRefName`, else `main`. When a step says to stop and report, the tidy failed: report the failure and the step where it happened.

## Grouping

Make the fewest commits that a first-time reviewer can follow, each one logical change, dependencies first.
- Commit a feature with its tests, a bugfix with its regression test, and generated code with its source. A refactor stands alone.
- Fold every fix into the commit that it corrects. A fix whose target lies below `<base>` stays a `fixup!` commit.
- One change repeated across services or packages is one commit.

## Steps

### 1. Base
If you hold only the Implementer seat, `<base>` is your current task's Tidy base. Reread it in the task's plan entry. Without a plan file, take it from the brief. If neither names a Tidy base, stop and report. Skip steps 6 and 7.

Otherwise, run `git fetch origin <target>`. If origin has no `<target>` branch, stop and report. `<base>` is the output of `git merge-base origin/<target> HEAD`.

### 2. Preconditions
If `git log <base>..HEAD` lists no commits, end the skill. If `git status --porcelain` prints anything, stop and report.

### 3. Survey
Record `<orig>`, the output of `git rev-parse HEAD`, in your report. Read each commit of `git log --reverse --oneline <base>..HEAD` with `git show <sha>`. Plan the new commits by the rules in "Grouping".

### 4. Regroup
If the plan only folds `fixup!` commits, run `git rebase --autosquash <base>`. On a conflict, resolve it toward the tree of `<orig>` and run `git rebase --continue`. If you cannot, run `git rebase --abort`, then stop and report.

To merge, split, or reorder commits, run `git reset -N <base>`. Then make the planned commits in order:
- Write each message in the format of the `commit` skill. Carry over the `Co-Authored-By` trailers of the commits that it replaces.
- For a commit with exactly the changes of one original commit, run `git commit -C <sha>`.
- Do not run tests per commit: the worktree holds the final tree.
- To split a file, use the piped `git add -p` from the `commit` skill. A file that the commits added is one hunk. To commit it as it was at `<sha>`, run `git restore --source=<sha> <path>`, stage it, and commit. Then run `git restore --source=<orig> <path>`.

### 5. Verify
If `git status --porcelain` or `git diff <orig>` prints anything, run `git reset --hard <orig>`, then stop and report. Otherwise, put `git log --oneline <base>..HEAD` in your report.

### 6. Rebase onto origin/<target>
Run `git rebase origin/<target>`. On a conflict, resolve it so that the result keeps the intent of both sides. If no resolution keeps both, run `git rebase --abort`, then stop and report. Then run the repo's checks. If they fail, leave the rebased branch in place, then stop and report.

### 7. Push
If `gh pr view --json state --jq .state` prints `OPEN`, run `git push --force-with-lease`.
