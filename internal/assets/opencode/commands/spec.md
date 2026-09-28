---
description: Diseña la Spec Unificada para una idea o requerimiento sin implementar código todavía
agent: orchestrator
---

The user wants to explore and architect a specification: "$ARGUMENTS".

Follow the Clean Spec Workflow:
1. **Explore & Architect:** Invoke the `architect` subagent to inspect the codebase, verify dependencies, and generate a concise Unified Spec (Scope, Interfaces/Contracts, and Atomic TDD Checklist).
2. **Review:** Present the resulting Unified Spec clearly to the user. Do NOT proceed to implementation until the user explicitly requests it (e.g. via `/feature`).
