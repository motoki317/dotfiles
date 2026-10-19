# Single source of truth

- Flags or fields that must agree, such as `is_running` and `is_failed`: use one enum or status field. Derive any existing flags from it, or delete them.
- Branches that differ only in their data: put the data in one table keyed by the case. Read the table in one code path.
- Switches over the same cases in several files: use one handler map or one method per variant.
- Checks that callers repeat before each transition, such as `if started and not closed`: put the checks and the transitions in one state machine.
- A list that another list must match: keep one list, and derive the others from it. If you cannot derive them, add a test that fails when they differ.
- A stored value that other values determine: compute it from those values instead of storing it. If a query needs it stored, make the database derive it, for example with a generated column.

Refactor only the duplicates that your change touches. Leave the others as they are, and mention them in your reply.
