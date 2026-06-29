# Safety — ISO/IEC 25010:2023

> Capability of a product under defined conditions to avoid a state in which human life, health, property, or the environment is endangered.

(New characteristic in the 2023 edition. Most relevant when the change performs destructive, irreversible, or high-stakes operations.)

## Official sub-characteristic definitions (ISO/IEC 25010:2023)
- **Operational constraint** — Degree to which a product or system constrains its operation to within safe parameters or states when encountering operational hazard.
- **Risk identification** — Degree to which a product can identify a course of events or operations that can expose life, property or environment to unacceptable risk.
- **Fail safe** — Degree to which a product can automatically place itself in a safe operating mode, or to revert to a safe condition in the event of a failure.
- **Hazard warning** — Degree to which a product or system provides warnings of unacceptable risks to operations or internal controls so that they can react in sufficient time to sustain safe operations.
- **Safe integration** — Degree to which a product can maintain safety during and after integration with one or more components.

## What to look for (review guidance — not part of ISO/IEC 25010)
- **Operational constraint** — destructive or bulk operations with no bounds, guards, or safe-mode under risky conditions.
- **Risk identification** — acting without first detecting the dangerous condition (e.g. deleting without checking what is matched).
- **Fail safe** — a failure that leaves the system in a destructive or partially-applied state; no safe default on error.
- **Hazard warning** — destructive actions with no confirmation, dry-run, or warning before the point of no return.
- **Safe integration** — a new component that can trigger destructive paths in others; migrations without a tested rollback.
