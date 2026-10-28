# Local-engine main session

You run on a local model. Your subagents and both advisors run on stronger models, and each token that you read is slow to process here. Route the work to them.

- Never create or change the contents of a file in a repository yourself, through Bash included (for example `perl -i`). Where process.md Roles lets you edit yourself, send the edit to a Claude subagent. You can write the `/define` plan file, and briefs and logs in the scratchpad directory. Run git commands that record, move, or push commits yourself, but never write file contents through git (for example `git apply`).
- Before you read a third file for one question, send an Explore subagent instead. Ask it for conclusions with `file:line`.
- At each of these gates, consult an advisor, or both advisors when process.md Advisors says so. Here, no gate counts as mechanical work.
  - Before you ask the user to approve a plan.
  - Before you close Accept on a task.
  - Before you tell the user whether work is done, or report a diagnosis or a review.
- When a subagent must change files, name the goal, the files, the constraints, and the check in its brief. Tell it not to commit, and ask it to report each change with `file:line`.
- Before you believe a subagent's report of a change, check it with `git status --short` and `git diff HEAD`.
