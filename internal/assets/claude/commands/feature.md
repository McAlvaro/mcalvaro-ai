---
description: Diseña e implementa una nueva feature con la Tríada Limpia (Arquitecto + Obrero)
---

The user wants to develop a new feature: "$ARGUMENTS".

Follow the Clean Triad Workflow:
1. **Explore & Architect:** Delegate to the `architect` subagent (or adopt the Architect role) to explore the codebase and produce a concise Unified Spec (Scope, Interfaces/Contracts, and Atomic TDD Checklist).
2. **Review Gate:** Present the Unified Spec clearly to the user and ask for approval. STOP and wait for confirmation.
3. **Build & Test:** Once approved by the user, delegate to the `builder` subagent (or adopt the Builder role) to implement the spec using strict TDD (RED -> GREEN -> REFACTOR).
4. **Deterministic Verify:** Run the local test harness (`gentle-ai test`) until exit code 0. Deliver the final result.
