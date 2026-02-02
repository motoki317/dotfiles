---
argument-hint: [upstream-url]
description: Upstream PR preparation and review command
---

# Purpose
Review and prepare changes before submitting PRs to upstream OSS repositories. Auto-fetch contribution guidelines, analyze changes, evaluate tests, and generate compliant PR metadata.

# Rules
- **Read-only operation**: Analyze and report only, no file modifications
- Auto-fetch CONTRIBUTING.md (check root, .github/, docs/)
- Verify gh CLI authentication before PR history operations
- Never create PR automatically; only provide review and suggestions

# Workflow
1. **Preflight**:
   - Check Serena memories for contribution patterns
   - Verify gh CLI authentication: `gh auth status`
   - Detect upstream from git remotes (prefer `upstream`, fallback to `origin`)
   - Get current branch and pending changes
   - Compare with upstream default branch

2. **Gather** (parallel):
   - Fetch CONTRIBUTING.md from upstream
   - Analyze code changes (quality-assurance agent)
   - Evaluate test coverage (tests agent)
   - Fetch author's past PRs: `gh pr list --author @me --repo upstream --state all --limit 20`

3. **Synthesize**:
   - Generate PR title and description following contribution guide
   - Determine local verification commands
   - Detect change types and generate manual QA checklist

4. **Validate**: Cross-validate guideline compliance with code review findings

# Agents (all read-only)
- **guidelines**: Parse CONTRIBUTING.md and extract requirements
- **changes**: Review code changes for quality
- **tests**: Evaluate test coverage
- **history**: Analyze author's past PR feedback
- **metadata**: Generate compliant PR title/description
- **verify**: Determine verification commands, detect change types
- **validator**: Cross-validate findings

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
