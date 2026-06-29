# Flexibility — ISO/IEC 25010:2023

> Degree to which a product can be adapted to changes in its requirements, contexts of use or system environment.

(Renamed from "Portability" in the 2011 edition, with Scalability added as a sub-characteristic.)

## Official sub-characteristic definitions (ISO/IEC 25010:2023)
- **Adaptability** — Degree to which a product or system can effectively and efficiently be adapted for or transferred to different hardware, software or other operational or usage environments.
- **Scalability** — Degree to which a product can handle growing or shrinking workloads or to adapt its capacity to handle variability.
- **Installability** — Degree of effectiveness and efficiency with which a product or system can be successfully installed and/or uninstalled in a specified environment.
- **Replaceability** — Degree to which a product can replace another specified software product for the same purpose in the same environment.

## What to look for (review guidance — not part of ISO/IEC 25010)
- **Adaptability** — hardcoded paths, OS/arch assumptions, environment-specific behaviour baked into logic.
- **Scalability** — super-linear cost on growth, stateful bottlenecks, no horizontal-scale path.
- **Installability** — undocumented dependencies, broken setup/teardown, no clean uninstall or down-migration.
- **Replaceability** — lock-in from leaked vendor specifics, no abstraction over a swappable dependency.
