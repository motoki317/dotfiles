---
name: investigate
description: Root-cause investigation — enumerate competing hypotheses, verify each with evidence, suggest a fix but never apply it. Use when an error, failing test, or anomaly needs diagnosing before any fix is written; implementation hands off to /execute.
argument-hint: "<error-or-symptom>"
allowed-tools: [Bash, Read, Grep, Glob, Agent, AskUserQuestion]
---

# Rules
- Read-only: never modify files or apply a fix; suggest the change and hand it to /execute.
- Judge from evidence (logs and code first), never from speculation; build the chain from symptom to cause, and report honestly when it can't be confirmed.
- Enumerate at least 3 hypotheses before investigating any; test the most likely first, and re-rank the rest whenever evidence refutes one.
- Ask the user only for what you can't infer — scope boundaries (which repos/services are in play), reproduction conditions, or when the issue began.

# Workflow
1. **Gather** in parallel — collect error messages, stack traces, and logs; recent changes (`git log --oneline -20`, `git diff HEAD~5..HEAD --stat`); related code paths; relevant config. Spawn the read-only agents that fit:
   - `Explore` — error locations, related code paths
   - `reviewer` — stack-trace and error analysis
   - `general-purpose` — log, dependency, and external-source checks
2. **Hypothesize** — list at least 3 hypotheses ranked by likelihood, each with the evidence that suggests it.
3. **Verify systematically** — for each hypothesis (most likely first): define what evidence would confirm or refute it, gather that evidence, and record a verdict (Confirmed / Refuted / Inconclusive). Re-rank the remaining hypotheses when new information lands.
4. **Report** — once a cause is confirmed, give the root cause, its evidence chain, impact, a suggested fix for /execute, and the rejected hypotheses.

# Investigation Checklist
- Error location, stack trace, and the failing code path
- Dependencies, imports, and recent changes (commits, config, deploys)
- Runtime environment (OS, versions, env vars)
- Resource state (disk, memory, network)

# Output Format
```
## Overview
Error summary and how it was investigated

## Hypotheses
1. [Most likely] Description — because [evidence]
2. [Possible]    Description — because [evidence]
3. [Less likely] Description — because [evidence]

## Verification
### Hypothesis N: [description]
- Testing: what would confirm or refute it
- Evidence: what was found
- Verdict: Confirmed / Refuted / Inconclusive

## Root Cause
- Direct cause, underlying cause, and the conditions that trigger it
- Evidence that confirmed it
- Confidence: 0–100

## Impact
Scope and similar at-risk code

## Recommended Fix
Suggested change and files to touch (implement via /execute); risk assessment

## Rejected Hypotheses
| Hypothesis | Why rejected |
|------------|--------------|

## Further Investigation
Unresolved points and next steps (when the root cause isn't confirmed)
```

# Anti-patterns
- Fixing the first symptom without confirming the root cause
- Investigating only one hypothesis, or tunnel-visioning the most recent change
- Iterating on a failed approach instead of re-examining assumptions
