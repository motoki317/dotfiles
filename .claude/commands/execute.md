---
argument-hint: [task-description]
description: Task execution command
---

# Purpose
Execute tasks by decomposing them and delegating detailed work to specialized sub-agents, while keeping policy decisions and orchestration in the main loop.

# Rules
- Delegate detailed work to sub-agents; don't implement detailed logic directly.
- Run independent tasks in parallel.
- Verify sub-agent outputs before integrating them.
- Prefer basic tools (Read/Edit/Write) over the Codex MCP when they suffice.
- Don't declare the task done until the Verify step passes; show evidence, not assertions.

# Workflow
1. **Analyze** — identify tasks, select agents, find parallel opportunities.
2. **Decompose** — break work into units with clear boundaries.
3. **Structure** — define parallel vs sequential tasks and their dependencies.
4. **Assign** — delegate with: scope and expected deliverables, target file paths, reference implementations (specific paths), and a memory check (`list_memories`) for prior patterns.
5. **Consolidate** — integrate sub-agent outputs into a coherent whole.
6. **Verify** — before declaring done: run the available tests/build/lint, invoke the `feedback` skill for static review, and the `qa-review` skill when there is runnable behavior — for analysis only; tell each not to run its close-the-loop phase, since this step owns the fix loop. Fix critical and major findings (re-delegate as needed) and re-check. Report the evidence — commands run and results — and don't finish with unresolved critical or major findings.

# Agents
Spawn `general-purpose` for work units and `reviewer` (read-only) for checks and cross-validation. The step-4 assignment carries the specialization, not the agent type.

# Codex Usage
Use the Codex MCP only for code generation (new files/functions) and code modification — never for research, verification, tests, docs, or reviews. Keep each call to one small task; no multi-file edits in a single call.
