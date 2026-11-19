---
name: pr
description: Push the current branch, open a pull request, and run /address. Use at the Ship step, when the user asks for a PR or authorized Ship at plan approval.
argument-hint: "[base-branch]"
allowed-tools: [Bash, Read, Grep, Agent]
---

Base: the argument, else `main`.

1. Run `/rebase-clean <base>`.
2. Write the title in Conventional Commits, under 70 characters. For the body, follow the repo's PR template or project PR skill. If the repo has neither, use `## Motivation`, `## Changes`, and `## Test plan`. The body tells the reviewer why the change exists, what changed, how you verified it, and which findings you skipped. Keep any other fact out of the body. If a later reader of the code needs that fact and the change does not hold it yet, add it to the code, a comment, or the docs, and run step 1 again. Load `/tech-writing` before you write the body, and run `/cold-read` on the draft.
3. If the change alters a UI, show each changed screen or state before and after the change, in a section of the body. Take the before screens from a build of `origin/<base>`, for example in a `git worktree`. Push the PNGs to a secret gist in one Bash call, and embed their raw URLs:
   ```bash
   export GH_TOKEN=$(gh auth token --user <account that opens the PR>)
   url=$(gh gist create -d "<repo> PR screenshots" <a one-line README.md>)
   git -c credential.helper='!gh auth git-credential' clone -q "https://gist.github.com/${url##*/}.git" <dir>
   cp <PNGs> <dir> && git -C <dir> add . && git -C <dir> commit -qm screenshots
   git -C <dir> -c credential.helper='!gh auth git-credential' push -q
   # In the body: ![after](https://gist.githubusercontent.com/<account that opens the PR>/<gist id>/raw/<file>.png)
   ```
4. Write the body to a file, and end it with the line `@codex review`. Never push to `main`. Push and open the PR:
   ```bash
   git push --force-with-lease -u origin HEAD
   gh pr create --base <base> --title "<title>" --body-file <file>
   ```
5. Unless the user said otherwise, run `/address` on the new PR.

To change the body of an open PR, fetch it with `gh pr view <pr> --json body -q .body > <file>`, and count its images with `grep -oE 'user-attachments|gist\.githubusercontent' <file> | wc -l`. Edit the file, apply it with `gh pr edit <pr> --body-file <file>`, fetch the body again, and confirm that the count did not change.
