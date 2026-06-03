---
estado: Hecho
prioridad: Alta
creado: 2026-06-03
---

# Tests E2E de Regresión en la API

## 📝 Descripción
Es necesario implementar una suite de tests automatizados en Go dentro de la carpeta `tests/` para asegurar que futuros arreglos o nuevas features no introduzcan regresiones en los endpoints existentes. Estos tests atacarán directamente a la capa HTTP (endpoints).

## ✅ Criterios de Aceptación
- [ ] Refactor del router en `main.go` para que sea testeable de forma aislada.
- [ ] Implementación de `tests/health_test.go` para probar el endpoint `/health`.
- [ ] Implementación de `tests/character_test.go` para probar `/api/v1/characters` con paginación y por ID.
- [ ] Los tests corren exitosamente con `make test`.

## 🏗️ Implementación Técnica
Ver especificación en la propuesta del plan de implementación.
