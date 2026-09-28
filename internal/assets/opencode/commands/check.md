---
description: Ejecuta el Test Harness determinista local (compilador, linter, tests)
agent: orchestrator
---

The user wants to run the project's test harness: "$ARGUMENTS".

1. Execute the native test harness using the Bash tool: `mcalvaro-ai test` (or `gentle-ai test`).
2. If tests pass (exit code 0), report success concisely with the step count and timing.
3. If tests fail (exit code non-zero), report the failing step and stack trace for immediate diagnosis.
