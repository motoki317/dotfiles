---
argument-hint: [file-path]
description: Markdown text update command
---

# Purpose
Output results from other commands (/define, /ask, /investigate, etc.) as markdown files.

# Rules
- Retrieve previous command execution results
- Determine output filename based on context (use specified path if provided)
- Never include revision history or discussion process
- Never add timestamps to documents
- Verify code examples are correct before including

# File Mapping
| Source Command | Default Output |
|----------------|----------------|
| /define | EXECUTION.md |
| /ask | RESEARCH.md |
| /investigate | RESEARCH.md |
| other | MEMO.md |

User-specified file path takes precedence.

# Workflow
1. Identify the previous command and its output
2. Determine appropriate output file location
3. Write cleaned, formatted content (no revision history)

# Constraints
- Use context-appropriate filename
- Avoid including consideration process or discussion history
