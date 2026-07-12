---
name: execute
description: Execute a task by decomposing it and delegating detailed work to sub-agents, keeping policy decisions and orchestration in the main loop, then verifying with tests and the feedback/qa-review skills. Use at the Implement step of the operating loop, for any multi-part build.
argument-hint: "[task-description]"
---

# Rules
- Delegate detailed work to sub-agents; don't implement detailed logic directly.
- Run independent tasks in parallel.
- Verify sub-agent outputs before integrating them.
- Prefer basic tools (Read/Edit/Write) over `codex-run work` when they suffice.
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
`codex-run work -C <repo> < task.md` delegates one implementation unit to OpenAI Codex (cross-model; workspace-write sandbox — see `codex-run --help`). Use it only for code generation (new files/functions) and code modification — never for standalone research, verification, docs, or review tasks. A work unit still includes its own tests and runs the smallest relevant check.
- One small task per call, scoped to a few files at most. Write the task brief like a step-4 assignment: scope, target paths, references, and the `/define` plan-file path.
- Codex edits the working tree directly: don't edit the same checkout while a call runs; parallel calls need disjoint files or separate worktrees.
- Non-zero exit = failed turn, and partial edits may remain — check `git status`/`git diff` before integrating or retrying.
