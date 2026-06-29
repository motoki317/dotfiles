# Maintainability — ISO/IEC 25010:2023

> Degree of effectiveness and efficiency with which a product or system can be modified to improve it, correct it or adapt it to changes in environment, and in requirements.

## Official sub-characteristic definitions (ISO/IEC 25010:2023)
- **Modularity** — Degree to which a system or computer program is composed of discrete components such that a change to one component has minimal impact on other components.
- **Reusability** — Degree to which a product can be used as an asset in more than one system, or in building other assets.
- **Analysability** — Degree of effectiveness and efficiency with which it is possible to assess the impact on a product or system of an intended change to one or more of its parts, to diagnose a product for deficiencies or causes of failures, or to identify parts to be modified.
- **Modifiability** — Degree to which a product or system can be effectively and efficiently modified without introducing defects or degrading existing product quality.
- **Testability** — Degree of effectiveness and efficiency with which test criteria can be established for a system, product or component and tests can be performed to determine whether those criteria have been met.

## What to look for (review guidance — not part of ISO/IEC 25010)
- **Modularity** — tight coupling, leaky abstractions, edits that ripple across unrelated modules.
- **Reusability** — duplicated logic, context-locked helpers that could be generalized.
- **Analysability** — tangled control flow, unclear naming, missing logs/observability, hidden state.
- **Modifiability** — fragile code, hidden side effects, missing tests around the changed behaviour.
- **Testability** — hard-wired dependencies with no seams, non-deterministic behaviour, new behaviour without tests.
