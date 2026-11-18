# Runbook

- Write each step as one checkbox line that names what the operator acts on (the workflow template, the job) and where they act: the UI path they use, or the command. Link each PR, page, and dashboard that a step names.
- If a step changes a file in a repository, link the PR that makes the change. Do not describe the change in prose.
- Before you write a read-only command down, run it. Check a command that changes state with its dry-run flag or on the local stack, never on the real target.
- Keep a sentence only if the operator needs it to do the step or to tell that the step worked.
- Name each section by when the operator reads it ("Before the release", "If the migration fails").
- In chat, give the operator only the next step.
