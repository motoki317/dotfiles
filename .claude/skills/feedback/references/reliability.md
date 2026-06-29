# Reliability — ISO/IEC 25010:2023

> Degree to which a system, product or component performs specified functions under specified conditions for a specified period of time.

## Official sub-characteristic definitions (ISO/IEC 25010:2023)
- **Faultlessness** — Degree to which a system, product or component performs specified functions without fault under normal operation.
- **Availability** — Degree to which a system, product or component is operational and accessible when required for use.
- **Fault tolerance** — Degree to which a system, product or component operates as intended despite the presence of hardware or software faults.
- **Recoverability** — Degree to which, in the event of an interruption or a failure, a product or system can recover the data directly affected and re-establish the desired state of the system.

## What to look for (review guidance — not part of ISO/IEC 25010)
- **Faultlessness** — latent bugs and unchecked assumptions on the happy path, ignored error returns, swallowed exceptions.
- **Availability** — single points of failure, blocking initialization, missing health/readiness signals, deadlock/livelock risk.
- **Fault tolerance** — no handling of downstream failures, missing timeouts/retries/backoff/circuit-breakers, no partial-failure handling.
- **Recoverability** — non-idempotent retries, no transaction/rollback, in-flight state lost on crash, no resume.
