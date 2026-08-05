---
name: rebase-clean
description: Regroup unshipped commits into reviewer-readable logical units and push where authorized. Use to tidy trial-and-error history — the Tidy step.
allowed-tools: [Bash, Read]
---

# Rebase Clean

Rewrite the unshipped history so a reviewer reads a clean narrative, not your trial-and-error.

## Target — the unshipped work
`git fetch origin main` first, then pick the base by where you stand:
- **Feature branch**: `BASE=$(git merge-base main HEAD)` — never `main` itself, which may have advanced and would swallow its newer commits.
- **Default branch** (`main`): `BASE=origin/main` — the unpushed commits are the unit. If `origin/main` has commits you lack, `git rebase origin/main` first (conflicts per step 5); rebuilding from the worktree would otherwise revert them.

## Preconditions
- Commits exist in `$BASE..HEAD`.
- **Clean worktree**, apart from changes you mean to fold in: every commit is rebuilt from the whole worktree after the soft reset, so stray edits get swept in. If unrelated edits remain, stop and report.

## How to group commits
The unit is one **logical functional change**, named for what it does for a reader — not for which layer it touches.

- A feature ships with its tests in the same commit; a bugfix with its regression test — the reader sees the behaviour and its proof together.
- Absorb fix-up and trial-and-error commits into the parent they correct.
- Generated or derived code goes with the source it comes from.
- Order commits so a dependency lands before what builds on it, when that aids reading.

Split by technical layer (proto / domain / infra / UI) only when one functional commit would be too large to review or genuinely mixes unrelated concerns — a readability escape hatch, not the default.

## Steps

### 1. Survey
```bash
ORIG=$(git rev-parse HEAD)      # for the identity check in step 4
git log --oneline $BASE..HEAD   # the unshipped commits
```
Read every commit; identify the fix-ups and the logical units to regroup into.

### 2. Soft-reset to the base
```bash
git reset --soft $BASE
git reset HEAD
```

### 3. Recommit
`git add` each logical unit and commit in Conventional Commits form (follow the `commit` skill for message style; end new messages with your standard `Co-Authored-By` footer). A commit that survives unchanged keeps its message and authorship via `git commit -C <sha>`.

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
Resolve conflicts yourself from both sides' intent and verify the result (build/test if available); stop and report only when the correct resolution is genuinely undeterminable and either choice discards real work.

### 6. Push — only what's already published
- Feature branch with an open PR (`gh pr view --json number,state`): `git push --force-with-lease` — `~/.claude/rules/process.md` authorizes autonomous force-pushes only as follow-ups to an already-open PR.
- Otherwise — no PR, or on the default branch — stop after the rebuild: the first publish is Ship, and Ship is user-triggered.
