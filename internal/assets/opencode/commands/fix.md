---
description: Diagnostica y corrige un bug con la Tríada Limpia (Arquitecto + Obrero)
agent: orchestrator
---

The user wants to fix a bug: "$ARGUMENTS".

Follow the Clean Triad Workflow:
1. **Root Cause Analysis:** Invoke the `architect` subagent to locate the root cause in the codebase, identify failure conditions, and propose a minimal fix with a regression test.
2. **Review Gate:** Present the diagnosis and the proposed fix briefly to the user. STOP and wait for confirmation.
3. **Fix & Test:** Once approved, invoke the `builder` subagent to write the failing test first, apply the minimal fix, and run the local test harness until exit code 0.
4. **Deterministic Verify:** Confirm that the local test harness (via `gentle-ai test`) passed with exit code 0. Deliver the result.
