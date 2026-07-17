---
name: upstream
description: Prepare a contribution to an upstream OSS repo — review the diff and tests against its contribution guide and draft compliant PR metadata. Use before submitting a PR to an external repository. Read-only; never opens the PR.
argument-hint: "[upstream-url]"
---

# Upstream contribution review

Read-only: analyze and report; never create the PR or modify files.

1. **Orient** — detect the upstream repo (the argument, else the `upstream` remote, else `origin`); compare the current branch against its default branch.
2. **Gather**, in parallel where independent:
   - the contribution guide (root, `.github/`, `docs/`; fall back to generic OSS practice if absent) and any PR template;
   - the diff and its test coverage;
   - your past PR feedback in that repo (`gh pr list --author @me --repo <upstream> --state all`) — reviewers repeat themselves, so recurring feedback weighs heaviest.
3. **Report**:
   - findings against the guide, code quality, and test coverage — each pass/fail/warn with location;
   - draft PR title and description in the repo's required format;
   - local verification commands (lint, test, build) plus a manual QA checklist where the change warrants one (UI, API, cross-component);
   - overall status: ready / needs work / blocked.
