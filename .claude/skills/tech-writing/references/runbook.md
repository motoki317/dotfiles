# Runbook

- Write each step as one checkbox line that names what the operator acts on (the workflow template, the job) and where they act: the UI path they use, or the command. Link each PR, page, and dashboard that a step names.
- If a step changes a file in a repository, link the PR that makes the change. Do not describe the change in prose.
- Before you write a read-only command down, run it. Check a command that changes state with its dry-run flag or on the local stack, never on the real target.
- Keep a step only if this operator must do it. Leave out a step that another person or a script can do, and hand it to that person or script. Keep a sentence only if the operator needs it to do the step or to tell that the step worked.
- If the steps are split across posts or comments that the operator can open in any order, make each one complete. Repeat a shared step instead of a pointer to another post.
- Name each section by when the operator reads it ("Before the release", "If the migration fails").
- In chat, give the operator only the next step.
