---
description: Diagnostica y corrige un bug con la Tríada Limpia (Arquitecto + Obrero)
---

# /fix

1. **Arquitecto:** Análisis de causa raíz y propuesta de corrección mínima con test de regresión.
2. **Revisión Humana:** Presenta el diagnóstico al usuario y espera aprobación.
3. **Obrero (Builder):** Escribe el test primero (RED), aplica la corrección (GREEN) y refactoriza.
4. **Verificación Determinista:** Ejecuta `gentle-ai test` en la terminal local hasta exit code 0.
