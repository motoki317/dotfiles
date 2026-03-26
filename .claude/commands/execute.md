---
argument-hint: [task-description]
description: Task execution command
---

# Purpose
Execute tasks by delegating detailed work to sub-agents while focusing on policy decisions and orchestration.

# Rules
- Delegate detailed work to specialized sub-agents
- Execute independent tasks in parallel
- Verify sub-agent outputs before integration
- Prefer basic tools (Read/Edit/Write) over Codex MCP when sufficient

# Workflow
1. **Analyze**: Identify tasks, select agents, find parallel opportunities
2. **Decompose**: Break into manageable units with clear boundaries
3. **Structure**: Define parallel vs sequential tasks and dependencies
4. **Assign**: Delegate with detailed instructions and context
5. **Consolidate**: Verify and combine results

# Agents
| Agent | Purpose |
|-------|---------|
| quality | Syntax, type, format verification |
| security | Vulnerability detection |
| test | Test creation, coverage |
| docs | Documentation updates |
| review | Post-implementation review |
| performance | Performance optimization |
| database | Database design and optimization |
| infrastructure | Infrastructure design |
| git | Git workflow design |
| validator | Cross-validation |

# Codex Usage
**Allowed**: Code generation (new files/functions), code modification
**Prohibited**: Research/analysis, quality/security verification, tests, docs, reviews

Rules:
- Prefer basic tools when sufficient
- One clear, small task per call
- No multi-file edits in single call

# Delegation Requirements
- Specific scope and expected deliverables
- Target file paths
- Reference implementations (specific paths)
- Memory check: `list_memories` for patterns

# Constraints
- Delegate detailed work to sub-agents
- Execute independent tasks in parallel
- Verify outputs before integration
- Avoid implementing detailed logic directly
