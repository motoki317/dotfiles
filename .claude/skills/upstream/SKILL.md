---
name: upstream
description: Review a contribution against the upstream repo contribution guide and draft compliant PR metadata. Use before you submit a PR to an external repository. This skill is read-only.
argument-hint: "[upstream-url]"
---

# Upstream contribution review

Read-only: analyze and report. Never create the PR or modify files.

1. Orient — detect the upstream repo (the argument, else the `upstream` remote, else `origin`). Compare the current branch against the upstream repo's default branch.
2. Gather these inputs, in parallel where they are independent:
   - the contribution guide (in the root, `.github/`, or `docs/`) and any PR template. If the repo has no guide, fall back to generic OSS practice.
   - the diff and its test coverage.
   - your past PR feedback in that repo (`gh pr list --author @me --repo <upstream> --state all`) — reviewers repeat themselves, so recurring feedback weighs heaviest.
3. Report:
   - findings against the guide, code quality, and test coverage — each pass/fail/warn with location.
   - draft PR title and description in the repo's required format.
   - local verification commands (lint, test, build) plus a manual QA checklist where the change warrants one (UI, API, cross-component).
   - overall status: ready / needs work / blocked.
