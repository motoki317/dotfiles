# Compatibility — ISO/IEC 25010:2023

> Degree to which a product, system or component can exchange information with other products, systems or components, and/or perform its required functions while sharing the same common environment and resources.

## Official sub-characteristic definitions (ISO/IEC 25010:2023)
- **Co-existence** — Degree to which a product can perform its required functions efficiently while sharing a common environment and resources with other products, without detrimental impact on any other product.
- **Interoperability** — Degree to which a system, product or component can exchange information with other products and mutually use the information that has been exchanged.

## What to look for (review guidance — not part of ISO/IEC 25010)
- **Co-existence** — global/shared mutable state, port or file contention, singletons or caches that affect co-running processes.
- **Interoperability** — schema/format mismatches, protocol or API version assumptions, encoding/charset issues, contract drift with upstream producers or downstream consumers.
