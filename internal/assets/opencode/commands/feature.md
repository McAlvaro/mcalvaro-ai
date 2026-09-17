---
description: Diseña e implementa una nueva feature con la Tríada Limpia (Arquitecto + Obrero)
agent: orchestrator
---

The user wants to develop a new feature: "$ARGUMENTS".

Follow the Clean Triad Workflow:
1. **Explore & Architect:** Invoke the `architect` subagent to explore the codebase and produce a concise Unified Spec (Scope, Interfaces/Contracts, and Atomic TDD Checklist).
2. **Review Gate:** Present the Unified Spec clearly to the user and ask for approval. STOP and wait for the user's confirmation.
3. **Build & Test:** Once approved by the user, invoke the `builder` subagent to implement the spec using strict TDD (RED -> GREEN -> REFACTOR).
4. **Deterministic Verify:** Confirm that the local test harness (compiler, linter, test runner via `gentle-ai test`) passed with exit code 0. Deliver the final result.
