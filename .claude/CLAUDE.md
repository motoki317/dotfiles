# General Instructions

## Looking up documentations

- Use `context7` MCP to lookup libraries' official documentations, whenever necessary.

## Choosing solutions

- Prefer "simple" over "easy".
- Prefer "systematic" problem solving over rabbit hole of "workarounds".
- Respect "Occam's Razor" and do NOT make a large amount of rules / configurations.

## Using Subagents (Task tool)

- Do NOT use subagents to continue deep research into the main task.
Instead, continue in the main context so users can see the main research progress.
- Use subagents ONLY when executing self-contained research tasks.
- When using subagents, prepare the instruction steps clearly in the main context, so they do not get lost. NEVER spawn subagents without clear step-by-step instructions on the task.
- Use subagents in parlalel, if some simple tasks can be parallelized easily.
When doing so, prepare instructions in the main context BEFORE spawning subagents, so the subagents' behaviors are well-predictable.
