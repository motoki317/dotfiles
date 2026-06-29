# Functional Suitability — ISO/IEC 25010:2023

> Degree to which a product or system provides functions that meet stated and implied needs when used under specified conditions.

## Official sub-characteristic definitions (ISO/IEC 25010:2023)
- **Functional completeness** — Degree to which the set of functions covers all the specified tasks and intended users' objectives.
- **Functional correctness** — Degree to which a product or system provides accurate results when used by intended users.
- **Functional appropriateness** — Degree to which the functions facilitate the accomplishment of specified tasks and objectives.

## What to look for (review guidance — not part of ISO/IEC 25010)
Judge against the supplied **intent**, not the code's own apparent purpose.
- **Completeness** — cases in the intent the change omits; inputs, modes, or states the requirement implies but the code never handles.
- **Correctness** — wrong formulas, off-by-one, boundary/empty/overflow cases, incorrect return or status values, rounding and unit errors.
- **Appropriateness** — logic that technically works but adds awkward steps, solves the wrong shape of the problem, or doesn't fit how the task is actually performed.
