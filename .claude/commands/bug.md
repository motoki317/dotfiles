---
argument-hint: [error-message]
description: Root cause investigation command
---

# Purpose
Identify the root cause of an error or anomaly from evidence. Read-only; suggests fixes but never applies them.

# Rules
- Never modify, create, or delete files, and never implement fixes — suggestions only.
- Judge from facts (logs first), never from user speculation; verify before accepting a claim.
- Build the evidence chain from symptom to cause before concluding; never force a cause the evidence doesn't support.
- Report honestly when the cause cannot be identified.

# Workflow
1. **Analyze** — classify the error (syntax, runtime, logic), locate where it occurs, gather logs.
2. **Investigate** in parallel — spawn the read-only agents that fit:
   - `quality-assurance` — error tracking, stack-trace analysis
   - `general-purpose` — log analysis, dependency errors
   - `explore` — error locations, related code paths
   - `fact-check` — external source verification
3. **Gather** — collect runtime info and check resources (see checklist).
4. **Report** — identify the root cause with confidence metrics.

# Investigation Checklist
- Error location details
- Dependencies and imports
- Config files and recent changes
- Runtime environment (OS, versions, env vars)
- Resource state (disk, memory, network)

# Output Format
```
## Overview
Summary of error and investigation

## Log Analysis
Critical log information, error context

## Code Analysis
Relevant code, identified issues

## Root Cause
- Direct cause
- Underlying cause
- Conditions

## Metrics
- Confidence: 0-100
- Log Utilization: 0-100
- Objectivity: 0-100

## Impact
Scope, similar errors

## Recommendations
Fix suggestions (no implementation), prevention

## Further Investigation
Unclear points, next steps
```
