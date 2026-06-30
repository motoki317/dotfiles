---
name: rebase-clean
description: Regroup a branch's messy work-in-progress commits into clean, reviewer-readable logical units, rebase onto the latest main, and force-push. Use when a branch's history shows trial-and-error (fix-ups, reverts, WIP commits) that should be tidied before or after opening a PR. The feedback/qa-review/address skills invoke it to publish close-the-loop fixes.
allowed-tools: [Bash, Read, AskUserQuestion]
---

# Rebase Clean

Rewrite a branch's history so a reviewer reads a clean narrative, not your trial-and-error: regroup the commits into logical units, rebase onto the latest `main`, and force-push with `--force-with-lease`.

## Preconditions
- On a branch other than `main`, with commits ahead of `main`.
- **Clean worktree**, apart from changes you mean to fold into the commits: this skill rebuilds every commit from the whole worktree after the soft reset, so stray edits get swept in. If unrelated edits remain, stop and report.
- **An open PR for this branch** — required only on an autonomous (non-direct) run, before the force-push: `gh pr view --json number,state`. `~/.claude/rules/loop.md` authorizes autonomous force-pushes only as follow-ups to an already-open PR; with none, the first push stays user-triggered — commit and stop.

## How to group commits
The unit is one **logical functional change**, named for what it does for a reader — not for which layer it touches.

- A feature ships **with its tests** in the same commit; a bugfix with its regression test. The reader sees the behaviour and its proof together.
- Absorb fix-up and trial-and-error commits into the parent they correct, so the history hides the stumbles.
- Put generated or derived code in the same commit as the source it comes from.
- Order the commits so a dependency lands before what builds on it, when that aids reading.

Split by technical layer (proto / domain / infra / UI) only as a fallback — when one functional commit would be too large to review, or genuinely mixes unrelated concerns. Layering is a readability escape hatch, not the default.

## Approval and autonomy
Default to running unattended: regroup, rebase, resolve conflicts, and force-push without prompting. Ask only when a user invoked this skill **directly** and the choice genuinely needs them — a conflict where either resolution discards real work, or a destructive step whose intent you cannot infer. The close-the-loop callers (`/feedback`, `/qa-review`, `/address`) never prompt.

Autonomy is fenced by the Preconditions self-checks, not by prompts. The open-PR and clean-worktree checks are what justify skipping approval on a non-direct run; if either fails, stop and report rather than push.

## Steps

### 1. Survey — reset target is `merge-base`, not `main`
`main` may have advanced; resetting to it would swallow main's newer commits. Reset to the merge-base instead.
```bash
MERGE_BASE=$(git merge-base main HEAD)
git log --oneline $MERGE_BASE..HEAD   # this branch's own commits
git log --oneline $MERGE_BASE..main   # how far main has advanced (may be empty)
```
Read every commit; identify the fix-ups and the logical functional units to regroup into.

### 2. Plan
Group per the principle above. Present the plan for approval only when the Approval rule calls for it; otherwise proceed.

### 3. Soft-reset to the merge-base
```bash
git reset --soft $MERGE_BASE
git reset HEAD
```
Reset to `$MERGE_BASE`, never `main`, so only this branch's own changes are restaged.

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
On conflict, resolve it yourself from both sides' intent and verify the result (build/test if available) before continuing. Stop only when the correct resolution is genuinely undeterminable and proceeding would lose real work — then report it, consulting the user if one invoked you directly.

### 7. Push
```bash
git log --oneline main..HEAD
git push --force-with-lease
```
