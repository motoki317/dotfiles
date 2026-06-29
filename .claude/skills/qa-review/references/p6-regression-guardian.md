# P6 — 回帰デグレ番人 (Regression guardian)

> Did something that used to work break? Peripheral features, post-reload behaviour.

Looks outward from the change to what it might have broken — the hardest to hold in one's head. Scope from the blast radius: shared modules, callers, common components, shared data.

## What to probe
- **Peripheral features** — screens/commands/endpoints sharing code or data with the change but not its target; still behave unchanged?
- **Reload / navigation** — F5, back/forward, deep-link in, re-open after the action; state restored, not lost or doubled?
- **Pre-existing happy paths** — re-run core flows that worked before; no silent breakage?
- **Shared resources** — global config, shared cache, common validation the change touches; still correct for other consumers?
- **Cross-feature** — does it change a default, schema, or shared format another feature depends on? (cross-note `/feedback` Compatibility)

## Expected result
Everything that worked before still works, and state survives reload/navigation. Name which existing flows were checked, so coverage is auditable.
