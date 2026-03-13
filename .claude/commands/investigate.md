---
argument-hint: <error-or-symptom>
description: Multi-hypothesis root cause investigation — enumerate hypotheses, verify systematically, avoid premature fixes
allowed-tools: [Bash, Read, Grep, Glob, Agent, AskUserQuestion]
---

# Purpose
Investigate bugs and production issues by enumerating multiple hypotheses, ranking by likelihood, and systematically verifying each before implementing any fix. Avoids the trap of fixing the most obvious symptom.

# Usage
```
/investigate <description of error, symptom, or issue>
```

# Rules
- **NEVER implement a fix until root cause is confirmed with evidence**
- Enumerate at least 3 hypotheses before investigating any
- Investigate most likely hypothesis first, but don't tunnel-vision
- If first hypothesis is disproven, re-rank remaining hypotheses with new information
- Ask user for scope boundaries (which repos/services are available) if unclear

# Workflow

### Phase 1: Gather Context

Collect available evidence in parallel:
- Error messages, stack traces, logs
- Recent changes: `git log --oneline -20`, `git diff HEAD~5..HEAD --stat`
- Related code paths (grep for error strings, function names)
- Configuration that might be relevant

Ask user:
- When did this start? (deploy, config change, traffic spike?)
- Is it reproducible? (always, intermittent, specific conditions?)
- What is the investigation boundary? (which repos/services/infra available?)

### Phase 2: Hypothesize

List **at least 3** hypotheses ranked by likelihood:

```markdown
## Hypotheses
1. **[Most likely]** Description — because [evidence]
2. **[Possible]** Description — because [evidence]
3. **[Less likely]** Description — because [evidence]
```

Present to user before proceeding.

### Phase 3: Investigate Systematically

For each hypothesis (starting with most likely):

1. **Define** what evidence would confirm or refute it
2. **Gather** that specific evidence (read code, run commands, check logs)
3. **Verdict**: Confirmed / Refuted / Inconclusive
4. If refuted, update remaining hypothesis rankings with new information

```markdown
### Hypothesis 1: [description]
**Testing**: [what we're looking for]
**Evidence**: [what we found]
**Verdict**: Confirmed ✅ / Refuted ❌ / Inconclusive ⚠️
```

### Phase 4: Root Cause Report

Once confirmed:

```markdown
## Root Cause
**Cause**: Clear description of the root cause
**Evidence**: What confirmed it
**Impact**: What is affected

## Recommended Fix
- Description of the fix
- Files to change
- Risk assessment

## Rejected Hypotheses
| Hypothesis | Why Rejected |
|------------|--------------|
```

### Phase 5: Fix (with permission)

Ask user: "Shall I implement the fix?"

If yes:
1. Implement the minimal fix
2. If tests exist, run them to verify
3. Do NOT commit — leave for user to review

# Anti-patterns to Avoid
- Fixing the first symptom without confirming root cause
- Investigating only one hypothesis
- Assuming the most recent change is the cause without evidence
- Iterating on the same failed approach instead of re-examining assumptions
