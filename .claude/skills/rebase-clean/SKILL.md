---
name: rebase-clean
description: Regroup unshipped commits into reviewer-readable logical units and push where authorized. Use to tidy trial-and-error history — the Tidy step.
allowed-tools: [Bash, Read]
---

# Rebase Clean

Rewrite the unshipped history so a reviewer reads a clean narrative, not your trial-and-error.

## Target — the unshipped work
Run `git fetch origin main` first. Then pick the base for your current branch:
- Feature branch: `BASE=$(git merge-base main HEAD)`. Never use `main` itself: if `main` advanced, the rebuild reverts its newer commits.
- Default branch (`main`): `BASE=origin/main` — the unpushed commits are the unit. If `origin/main` has commits that you lack, run `git rebase origin/main` first (resolve conflicts as in step 5). If you skip this rebase, the rebuild from the worktree reverts those commits.

## Preconditions
- Commits exist in `$BASE..HEAD`.
- The worktree is clean, apart from changes that you mean to fold in. After the soft reset, you rebuild every commit from the whole worktree, so the new commits include any stray edits. If unrelated edits remain, stop and report.

## How to group commits
The unit is one **logical functional change**, named for what it does for a reader — not for which layer it touches.

- A feature ships with its tests in the same commit, and a bugfix with its regression test. The reader then sees the behavior and its proof together.
- Absorb fix-up and trial-and-error commits into the parent they correct.
- Generated or derived code goes with the source it comes from.
- If it helps the reader, order commits so that a dependency lands before what builds on it.

A split by technical layer (proto / domain / infra / UI) is a readability escape hatch, not the default. Use it only if one functional commit is too large to review or mixes unrelated concerns.

## Steps

### 1. Survey
```bash
ORIG=$(git rev-parse HEAD)      # for the identity check in step 4
git log --oneline $BASE..HEAD   # the unshipped commits
```
Read every commit. Identify the fix-ups and the logical units to regroup into.

### 2. Soft-reset to the base
```bash
git reset --soft $BASE
git reset HEAD
```

### 3. Recommit
`git add` each logical unit and commit it in Conventional Commits form. Follow the `commit` skill for the message style. End new messages with your standard `Co-Authored-By` footer. For a commit that stays unchanged, keep its message and authorship with `git commit -C <sha>`.

### 4. Verify
```bash
git log --oneline $BASE..HEAD
git status       # working tree clean
git diff $ORIG   # empty — same tree, only history changed
```

### 5. Rebase onto latest main (feature branch only, if main advanced)
```bash
git rebase main
```
Resolve conflicts yourself from both sides' intent. If a build or test command is available, run it to verify the result. Stop and report only if the correct resolution is undeterminable and either choice discards real work.

### 6. Push — only what is already published
- Feature branch with an open PR (`gh pr view --json number,state`): `git push --force-with-lease` — `~/.claude/rules/process.md` authorizes autonomous force-pushes only as follow-ups to an already-open PR.
- Otherwise — no PR, or on the default branch — stop after the rebuild: the first publish is Ship, and Ship is user-triggered.
