---
name: rebase-clean
description: Regroup a branch's messy commits into clean, reviewer-readable logical units, rebase onto the latest main, and force-push. Use when a branch's history shows trial-and-error (fix-ups, reverts, WIP commits) worth tidying before or after opening a PR; the feedback/qa-review/address skills invoke it to publish close-the-loop fixes.
allowed-tools: [Bash, Read, AskUserQuestion]
---

# Rebase Clean

Rewrite a branch's history so a reviewer reads a clean narrative, not your trial-and-error.

## Preconditions
- On a branch other than `main`, with commits ahead of `main`.
- **Clean worktree**, apart from changes you mean to fold in: every commit is rebuilt from the whole worktree after the soft reset, so stray edits get swept in. If unrelated edits remain, stop and report.
- **An open PR for this branch** — required only on an autonomous (non-direct) run, checked before the force-push: `gh pr view --json number,state`. `~/.claude/rules/process.md` authorizes autonomous force-pushes only as follow-ups to an already-open PR; with none, the first push stays user-triggered — commit and stop.

## How to group commits
The unit is one **logical functional change**, named for what it does for a reader — not for which layer it touches.

- A feature ships with its tests in the same commit; a bugfix with its regression test — the reader sees the behaviour and its proof together.
- Absorb fix-up and trial-and-error commits into the parent they correct.
- Generated or derived code goes with the source it comes from.
- Order commits so a dependency lands before what builds on it, when that aids reading.

Split by technical layer (proto / domain / infra / UI) only when one functional commit would be too large to review or genuinely mixes unrelated concerns — a readability escape hatch, not the default.

## Approval and autonomy
Run unattended by default: regroup, rebase, resolve conflicts, force-push. The Preconditions self-checks (open PR, clean worktree), not prompts, fence this autonomy — if either fails, stop and report rather than push. Ask only when a user invoked this skill **directly** and the choice genuinely needs them (a conflict where either resolution discards real work, a destructive step whose intent you cannot infer); the close-the-loop callers (`/feedback`, `/qa-review`, `/address`) never prompt.

## Steps

### 1. Survey — reset target is `merge-base`, not `main`
`main` may have advanced; resetting to it would swallow main's newer commits.
```bash
MERGE_BASE=$(git merge-base main HEAD)
git log --oneline $MERGE_BASE..HEAD   # this branch's own commits
git log --oneline $MERGE_BASE..main   # how far main has advanced (may be empty)
```
Read every commit; identify the fix-ups and the logical units to regroup into.

### 2. Plan
Group per the principle above; present the plan for approval only when the Approval rule calls for it.

### 3. Soft-reset to the merge-base
```bash
git reset --soft $MERGE_BASE
git reset HEAD
```
`$MERGE_BASE`, never `main`, so only this branch's own changes are restaged.

### 4. Recommit
`git add` each logical unit and commit in Conventional Commits form (follow the `commit` skill for message style). End each message with:
`Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>`

### 5. Verify
```bash
git log --oneline $MERGE_BASE..HEAD
git status   # working tree clean
```
Confirm every change landed in a commit.

### 6. Rebase onto latest main (only if main advanced)
```bash
git fetch origin main
git rebase main
```
Resolve conflicts yourself from both sides' intent and verify the result (build/test if available) before continuing. Stop only when the correct resolution is genuinely undeterminable and proceeding would lose real work — then report, consulting the user if one invoked you directly.

### 7. Push
```bash
git log --oneline main..HEAD
git push --force-with-lease
```
