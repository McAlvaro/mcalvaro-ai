---
name: spec-unificada
description: "Trigger: Al diseñar una feature o fix complejo. Especificación técnica unificada para Gentle AI Lite con arquitectura, contrato y checklist TDD."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "3.0"
  delegate_only: true
---

# Spec Unificada (Gentle AI Lite)

Skill de diseño y arquitectura para el rol **Arquitecto**.
Filosofía central: **CONCEPTS > CODE**. Razonamiento profundo, contratos estrictos y cero burocracia documental.

## Objetivo
Explorar el repositorio y producir una especificación unificada concisa (máximo 2 páginas) que sirva de contrato inequívoco para el rol Obrero (Builder) sin ensuciar el repositorio con archivos intermedios.

## Herramientas de Apoyo
- **`codegraph`**: Navegación estructural, análisis de dependencias e impacto de símbolos en el codebase.
- **`context7`**: Consulta de documentación oficial y mejores prácticas de frameworks/librerías de terceros.
- **`engram`**: Consulta selectiva al inicio (`mem_search`) de decisiones o convenciones arquitectónicas previas relevantes.

## Estructura de la Spec Unificada

Toda Spec Unificada debe contener exactamente estas tres secciones:

### 1. Alcance (Scope)
- **Meta principal:** Definición concisa del cambio o funcionalidad.
- **En alcance (IN Scope):** Lista explícita de lo que se implementará.
- **Fuera de alcance (OUT of Scope):** Límites claros para prevenir scope creep.

### 2. Arquitectura y Contratos
- **Modelos y tipos de datos:** Estructuras, campos e invariantes.
- **Firmas e interfaces:** Funciones, métodos, entradas y salidas.
- **Casos de borde y errores:** Comportamiento determinista ante fallos.

### 3. Checklist Atómico de TDD (2 a 5 Unidades)
Para cada unidad de trabajo:
1. **RED:** Test unitario o de integración que falla primero.
2. **GREEN:** Implementación mínima necesaria para pasar el test.
3. **REFACTOR:** Limpieza de código y cumplimiento de convenciones.

## Entrega al Orquestador
Entrega la Spec Unificada directamente en tu respuesta al Orquestador para la puerta de aprobación del usuario. No crees archivos de tareas temporales en disco.
