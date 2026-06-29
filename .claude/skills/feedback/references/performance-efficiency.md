# Performance Efficiency — ISO/IEC 25010:2023

> Degree to which a product performs its functions within specified time and throughput parameters and is efficient in the use of resources (such as CPU, memory, storage, network devices, energy, materials, etc.) under specified conditions.

## Official sub-characteristic definitions (ISO/IEC 25010:2023)
- **Time behaviour** — Degree to which the response time and throughput rates of a product or system, when performing its functions, meet requirements.
- **Resource utilization** — Degree to which the amounts and types of resources used by a product or system, when performing its functions, meet requirements.
- **Capacity** — Degree to which the maximum limits of a product or system parameter meet requirements.

## What to look for (review guidance — not part of ISO/IEC 25010)
- **Time behaviour** — blocking calls on hot paths, sequential awaits that could be concurrent, N+1 queries, missing pagination/streaming, work repeated per-item that could be hoisted.
- **Resource utilization** — unbounded memory growth, leaks, redundant allocations/copies, connections or file handles held open, missing reuse/pooling.
- **Capacity** — assumptions about list size, payload size, concurrency, or connection counts that break under realistic load.
