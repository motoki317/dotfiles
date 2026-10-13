# Correctness

Judge against the intent, not the code's apparent purpose. Look for:
- a case the intent states that the change omits, and a requirement that no test checks
- wrong results at boundaries: empty, max, off-by-one, units
- failure paths: error, timeout, partial write, retry, crash mid-operation, a swallowed error
- untrusted input that reaches a sink (SQL, shell, path, HTML), missing authz, secrets in logs

# Simplicity

Judge each addition by the ladder in `$HOME/.claude/rules/code.md`. Look for:
- an abstraction with one implementation or one caller
- configurable code used one way, and a parameter with the same value at every call site
- logic that an existing helper covers (name the helper), and a forwarding wrapper
- dead code, and production code that only tests use
- a check that repeats one the caller already made
- a test that asserts a contract another test at the same layer asserts
- verbosity copied from the repo's incidental style, not its public idioms

Do not propose to cut:
- behavior that a requirement needs, or the only test of a requirement
- a tenant/authz/PII check, or error handling that prevents data loss

# Risk

Cover each bullet whose trigger (before the colon) the change matches.
- Security — auth, tokens, crypto, deserialization, permissions, PII: another tenant's record by ID, a client-only check, an unverified token or signature, unbounded input (DoS, ReDoS).
- Performance — loops or queries over unbounded data, hot paths: N+1 queries, missing pagination, unbounded memory. Judge SQL by its execution plan, not its text. A correlated subquery in an UPDATE/DELETE can turn one pass into per-row rescans. PostgreSQL materializes a CTE referenced more than once into an unindexable buffer. `EXPLAIN` shows these plan shapes even on an empty database.
- Compatibility — a schema, data format, or API that other code or a deployed version reads: old code on new data, new code on old data.
- Safety — deletes, overwrites, bulk mutations, migrations: does it check what it matches, is a second run safe, and is there a tested way back?
- Interaction — UI, CLI, API signatures, error messages: keyboard and screen-reader reach, and a remedy in each error message.
