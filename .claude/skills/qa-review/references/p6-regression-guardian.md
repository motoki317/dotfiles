# P6: Regression guardian (回帰デグレ番人)

> Did something that used to work break? Peripheral features, post-reload behavior.

P6 looks outward from the change to what it might break, which is the hardest thing to keep in mind. Scope from the blast radius: shared modules, callers, common components, shared data.

## What to probe
- Peripheral features — screens/commands/endpoints that share code or data with the change but are not its target. Do they still behave unchanged?
- Reload / navigation — F5, back/forward, deep-link in, re-open after the action. Is the state restored, not lost or doubled?
- Pre-existing happy paths — re-run core flows that worked before. Does anything break silently?
- Shared resources — global config, shared cache, common validation the change touches. Are they still correct for other consumers?
- Cross-feature — does it change a default, schema, or shared format another feature depends on? (cross-note `/feedback` Compatibility)

## Expected result
Everything that worked before still works, and state survives reload/navigation. Name the existing flows that the scenarios check, so that coverage is auditable.
