# Code — elaborations

Detail behind `~/.claude/rules/code.md`; load when reviewing code or when an always-loaded rule needs its edge cases.

## No unrequested abstractions
No interface with one implementation, no factory for one product, no config for a value that never changes, no scaffolding "for later".

## Diff shape
Deletion over addition. Fewest files, shortest clear diff — once you understand the problem. Follow the repository's conventions unless they conflict with the always-loaded rules: match its public idioms, not its incidental verbosity.

## Comments
WHY and WHY-NOT only — a constraint, a tradeoff, a cut corner's ceiling and upgrade path (`// global lock; per-account locks if throughput matters`). Write code commentless first. A comment (unlike a commit message) reads from the current tree alone, where "now", "no longer", "previously" are the diff talking.

## Output
Every paragraph defending a simplification is complexity smuggled back in as prose. Exempt: verification evidence, and explanation the user explicitly asked for — give those in full. Complex request → ship the lazy version and question the rest in the same response.
