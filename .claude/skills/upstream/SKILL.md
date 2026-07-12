---
name: upstream
description: Prepare a contribution to an upstream OSS repo — fetch the contribution guide, review the diff and tests against it, and draft compliant PR metadata. Use before submitting a PR to an external repository. Read-only; never opens the PR.
argument-hint: "[upstream-url]"
---

# Rules
- Read-only: analyze and report only, no file modifications.
- Auto-fetch CONTRIBUTING.md (check root, `.github/`, `docs/`).
- Verify `gh` CLI auth before any PR-history operation.
- Never create the PR automatically; only review and suggest.

# Workflow
1. **Preflight** — verify `gh auth status`; detect upstream from remotes (prefer `upstream`, else `origin`); get the current branch and pending changes; compare against the upstream default branch.
2. **Gather** in parallel — spawn the read-only agents below: `guidelines`, `changes`, `tests`, `history`.
3. **Synthesize** — `metadata` drafts the PR title/description per the guide; `verify` determines local verification commands, detects change types (table below), and generates the matching manual QA checklist.
4. **Validate** — `validator` cross-checks guideline compliance against the review findings.

# Agents (read-only)
`Agent` is the role label used in the workflow; `Spawn as` is the real agent type to pass to the Agent tool.

| Agent | Spawn as | Role |
|-------|----------|------|
| guidelines | general-purpose | Parse CONTRIBUTING.md, extract requirements |
| changes | reviewer | Review the diff for quality |
| tests | reviewer | Evaluate test coverage |
| history | general-purpose | Analyze the author's past PR feedback (`gh pr list --author @me --repo upstream --state all --limit 20`) |
| metadata | general-purpose | Generate compliant PR title/description |
| verify | general-purpose | Determine verification commands, detect change types |
| validator | reviewer | Cross-validate findings |

# Change Type Detection
| Type | File Patterns |
|------|---------------|
| UI | *.css, *.scss, *.jsx, *.tsx, *.vue, *.svelte |
| API | **/api/**, **/routes/**, **/handlers/**, *.openapi.* |
| Integration | Changes spanning 3+ modules or external services |

# Output Format
```
## Upstream Review Summary
- Repo: owner/repo
- Branch: feature-branch
- Overall Score: XX/100
- Status: ready | needs_work | blocked

## Checklist
### Contribution Guidelines
- [pass|fail|warn] Guideline item

### Code Quality
- [pass|fail|warn] Issue with location

### Test Coverage
- [pass|fail|warn] Test evaluation

### Past Review Patterns
- [info] Recurring feedback to address

## PR Metadata
- Title: Suggested title
- Description: Suggested description

## Local Verification
- lint: npm run lint
- test: npm test
- build: npm run build

## Manual Verification (if applicable)
### UI Changes
- Verify visual layout, responsive behavior, dark mode, accessibility

### API Changes
- Verify endpoints, error responses, auth, payload structure

### Integration Changes
- Test end-to-end flows, cross-component data, state persistence

## Recommended Actions
- [high] Critical action
- [medium] Recommended improvement
- [low] Optional enhancement
```

# Failure Handling
- CONTRIBUTING.md not found → use generic OSS best practices
- gh CLI fails → skip PR history, note in report
- Upstream detection ambiguous → ask user to specify
