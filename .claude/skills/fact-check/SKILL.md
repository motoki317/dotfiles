---
name: Fact Check
description: This skill should be used when the user asks to "verify claims", "fact check", "validate documentation", "check sources", or needs verification of external source references. Provides patterns for systematic fact verification using Context7 and WebSearch.
---

# Fact Verification Patterns

## Tools

### Context7 (Primary for libraries)
- `resolve-library-id` - Get library ID first
- `query-docs` - Verify API behavior, best practices, deprecations

### Web Tools (Fallback)
- `WebSearch` - Standards, specifications, general technical facts
- `WebFetch` - Verify specific documentation URLs

## Verification Process

### 1. Extract Claims
Identify verifiable assertions:
- Library API claims: "useState returns a tuple"
- Documentation references: "according to the React docs"
- Standard compliance: "follows WCAG 2.1 AA"
- Version-specific behavior: "in React 18, Suspense..."

### 2. Select Source
| Claim Type | Source |
|------------|--------|
| Library/framework API | Context7 (trust score 7+) |
| Web standard/spec | WebSearch (w3.org, MDN) |
| Specific URL cited | WebFetch |
| General technical fact | WebSearch |

### 3. Assess Confidence
- **90-100**: Exact match with authoritative source
- **80-89**: Strong match, minor wording differences
- **70-79**: Partial match, some details unverified
- **60-69**: Weak match, significant uncertainty
- **<60**: No match or contradictory evidence

**Threshold**: Flag claims with confidence below 80

## Discrepancy Report Format
```
Claim: [original assertion]
Source: [where claim was made]
Verification: [Context7/WebSearch result]
Evidence: [actual information from source]
Confidence: [0-100]
Recommendation: [correction or note]
```

## Critical Rules
- Always verify against authoritative sources before flagging
- Context7 is primary for library/framework claims
- Document evidence source for every verification
- Note version context for version-specific claims
- Cross-reference disputed claims with multiple sources

## Anti-Patterns
- Marking verified without actual source check
- Single source reliance for disputed claims
- Ignoring version context differences
- Over-verifying obvious facts
