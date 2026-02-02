---
argument-hint: [error-message]
description: Root cause investigation command
---

# Purpose
Identify root causes from error messages and anomalous behavior, providing fact-based analysis without performing fixes.

# Rules
- **Never modify, create, or delete files**
- Never implement fixes; provide suggestions only
- Prioritize log analysis as primary information source
- Judge from facts, not user speculation
- Report honestly if cause cannot be identified

# Workflow
1. **Analyze**: Classify error type (syntax, runtime, logic), locate occurrence, gather logs
2. **Investigate** (parallel): Delegate to quality-assurance, explore, general-purpose, fact-check agents
3. **Gather**: Collect runtime info, check resources
4. **Report**: Compile findings with confidence metrics, identify root cause

# Agents (all read-only)
- **quality-assurance**: Error tracking, stack trace analysis
- **general-purpose**: Log analysis, dependency errors
- **explore**: Finding error locations, related code paths
- **fact-check**: External source verification

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

# Constraints
- Keep all operations read-only
- Build evidence chain from symptom to cause before concluding
- Never accept user speculation without verification
- Never force contrived causes when evidence is insufficient
