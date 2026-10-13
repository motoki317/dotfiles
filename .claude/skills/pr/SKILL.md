---
name: pr
description: Push the current branch, open a pull request, and run /address. Use when the user asks for a PR — the Ship step, always user-triggered.
argument-hint: "[base-branch]"
allowed-tools: [Bash, Read, Grep, Agent]
---

Base: the argument, else `main`.

1. Run `/rebase-clean <base>`. It does not push a branch that has no open PR.
2. Write the title in Conventional Commits, under 70 characters. For the body, follow the repo's PR template or project PR skill. If the repo has neither, use `## Motivation`, `## Changes`, and `## Test plan`. After the sections, add the line `@codex review`.
3. Write the body to a file. Never push to `main`. Push and open the PR:
   ```bash
   git push --force-with-lease -u origin HEAD
   gh pr create --base <base> --title "<title>" --body-file <file>
   ```
4. Unless the user said otherwise, run `/address` on the new PR.

To change the body of an open PR, edit the live body. Fetch it with `gh pr view <pr> --json body -q .body`, edit it, and apply it with `gh pr edit <pr> --body-file <file>`. Then confirm that `grep -c user-attachments` still prints the count from before.
