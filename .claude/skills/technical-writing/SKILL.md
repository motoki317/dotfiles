---
name: Technical Writing
description: This skill should be used when the user asks to "write blog post", "technical article", "tutorial", "explain concept", or needs guidance on technical writing for external audiences. Provides patterns for technical blogs and articles in both English and Japanese.
---

# Technical Writing

Patterns for technical blogs and articles for external audiences.

## Article Types

| Type | Audience | Length | Structure |
|------|----------|--------|-----------|
| Tutorial | Learners | 1500-3000 | Problem -> Steps -> Example -> Troubleshooting |
| Concept | Understanding seekers | 1000-2500 | Hook -> Core -> Examples -> Misconceptions |
| Comparison | Decision makers | 1500-2500 | Criteria -> Options -> Features -> Recommendation |
| Case Study | Leaders | 1500-3000 | Challenge -> Solution -> Results -> Lessons |
| Opinion | Experienced devs | 800-1500 | Thesis -> Arguments -> Counterarguments -> Action |

## Core Principles

- **Hook early**: Problem statement, surprising fact, or relatable story in first paragraph
- **Inverted pyramid**: Lead with conclusion, most important info first
- **Show don't tell**: Working code, diagrams, before/after comparisons
- **One idea per section**: Split if covering multiple ideas
- **Respect reader time**: Be concise, provide TL;DR for long articles

## Title Patterns

- **How-to**: "How to Implement Rate Limiting with Redis"
- **Number list**: "5 Things Every Developer Should Know About TypeScript Generics"
- **Comparison**: "REST vs GraphQL: Which Should You Choose in 2025?"
- **Problem-solution**: "Solving N+1 Queries with DataLoader"
- **Deep dive**: "Understanding React Reconciliation: A Deep Dive"
- **Lessons**: "What I Learned Building a Real-Time Collaboration System"

## Language Guidelines

**English:**
- Conversational but professional
- First person ("I found...") or second person ("You can...")
- Vary sentence length

**Japanese:**
- Technical articles: です・ます調
- Personal blogs: flexible
- English technical terms as katakana where appropriate

## Writing Workflow

1. **Ideation**: Core message (one sentence), audience, article type
2. **Outline**: Section headings, key points, examples needed
3. **Draft**: Body first (not intro), include tested code, transitions
4. **Edit**: Read aloud, cut 20%, verify claims, check code
5. **Polish**: Proofread, optimize title/headings, add meta description

## Best Practices

**Critical:**
- Start with compelling hook, never "In this article..."
- Test all code examples before publishing
- Verify all technical claims

**High:**
- Explain code context before/after snippets
- Support claims with evidence (benchmarks, examples)
- Each section has one clear purpose

## Anti-Patterns

| Avoid | Instead |
|-------|---------|
| Burying the lede | Most important point first |
| Assuming motivation | Explain why reader should care |
| Code without context | Explain purpose before/after |
| Unsubstantiated claims | Benchmarks or reasoned arguments |
| Clickbait title | Title matches content scope |
| Wall of code | Break up with explanations |
| Generic intro | Hook that captures attention |

## Constraints

**Must:**
- Verify all technical claims
- Test all code examples
- Write for target audience level

**Avoid:**
- Jargon without explanation
- Untested code examples
- Overly complex explanations
