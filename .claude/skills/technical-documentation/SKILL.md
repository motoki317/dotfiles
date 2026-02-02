---
name: Technical Documentation
description: This skill should be used when the user asks to "write documentation", "create README", "API docs", "design document", "specification", "user guide", or needs documentation guidance. Provides comprehensive documentation patterns for developers, teams, and end-users in both English and Japanese.
---

# Technical Documentation

Patterns for README, design documents, API specifications, and user guides.

## Document Types

| Type | Audience | When to Use |
|------|----------|-------------|
| README | Developers, contributors | Main project entry point |
| API Spec | API integrators | REST/GraphQL endpoints |
| Design Doc | Team, reviewers | Major features, architecture |
| User Guide | End users, admins | Help, tutorials, product guides |

## Core Concepts

- **Progressive disclosure**: Quick start -> common cases -> advanced -> edge cases
- **Audience-first**: Match depth to reader's knowledge level
- **Example-driven**: Show working code, not just descriptions

## README Structure

```markdown
# Project Name
[![Badge](https://img.shields.io/badge/...)]

Brief one-line description.

## Features
- Feature 1
- Feature 2

## Quick Start
npm install package-name

## Basic Usage
import { example } from "package-name";
const result = example();

## Documentation
See [full documentation](link).

## Contributing / License
```

## API Documentation Structure

```markdown
# API Reference

## Authentication
Authorization: Bearer YOUR_API_KEY

## Base URL
https://api.example.com/v1

## Endpoints

### GET /users
**Parameters:** limit (int), offset (int)
**Response:** { "users": [...], "total": 100 }
**Errors:** 401 Unauthorized, 429 Rate limited
```

## Design Document Structure

```markdown
# Feature Name Design Document

## Summary
**Problem:** [description]
**Solution:** [approach]
**Scope:** [included/excluded]

## Goals / Non-Goals

## Technical Design
- Architecture
- Data flow
- Components

## Alternatives Considered

## Security Considerations

## Testing Strategy

## Rollout Plan
```

## User Guide Structure

```markdown
# User Guide

## Getting Started
## Core Concepts
## Step-by-Step Tutorials
## Feature Reference
## Troubleshooting / FAQ
## Glossary
```

## Best Practices

**Critical:**
- Write for specific audience's knowledge level
- Start with essentials, reveal complexity gradually
- Verify all code examples compile and run

**High:**
- Make content scannable (headings, bullets, code blocks)
- Include working, copy-pasteable examples
- Use active voice: "Run this command" not "The command can be run"

## Language Guidelines

**English:**
- Active voice, present tense
- Professional but approachable
- Avoid idioms that don't translate

**Japanese:**
- User docs: です・ます調
- Technical specs: である調
- Avoid excessive katakana, ambiguous expressions

## Anti-Patterns

| Avoid | Instead |
|-------|---------|
| Wall of text | Bullets, headings, code blocks |
| Outdated info | Document current state |
| Assuming context | Define terms, link prerequisites |
| Untested examples | Verify code compiles/runs |
| Passive voice | Active voice for clarity |
| Jargon overload | Define terms or use simpler language |

## Constraints

**Must:**
- Verify accuracy against implementation
- Include runnable code examples
- Follow project documentation style

**Avoid:**
- Documenting without reading code
- Adding timestamps to documents
- Duplicating information
