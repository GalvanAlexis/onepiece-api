# 📋 ISSUES — Tracker Local

Este directorio reemplaza herramientas externas (Notion, Jira, etc.) como tracker de tareas del proyecto `onepiece-api`.

---

## ¿Cómo crear un issue?

1. Copiar el template: `ISS-000_template.md`
2. Renombrar con el número correlativo: `ISS-001_nombre-descriptivo.md`
3. Completar todos los campos del template
4. Cambiar el estado a `Backlog`

---

## Formato de nombre de archivo

```
ISS-XXX_descripcion-breve-en-kebab-case.md
```

**Ejemplos:**
```
ISS-001_endpoint-get-characters.md
ISS-002_middleware-rate-limiting.md
ISS-003_seed-data-frutas-del-diablo.md
```

---

## Estados del ciclo de vida

| Estado | Descripción | Rama Git asociada |
|---|---|---|
| `Backlog` | Definido, no iniciado | — |
| `En Progreso` | Trabajo activo | `feature/ISS-XXX` |
| `En Review` | PR abierto, esperando merge | `feature/ISS-XXX` |
| `Hecho` | PR mergeado a `main` | — |
| `Cancelado` | Descartado (con justificación) | — |

---

## Índice de Issues

| ID | Título | Estado | Prioridad |
|---|---|---|---|
| ISS-000 | Template (no usar) | — | — |

> Actualizar esta tabla cada vez que se crea o cierra un issue.

---

## Reglas para el Agente IA

- Todo cambio de código debe tener un issue asociado.
- Si no existe el issue, crearlo **antes** de abrir la rama.
- Al abrir el PR, el issue debe pasar a `En Review`.
- Al confirmar el merge, actualizar a `Hecho`.
- Features complejas requieren un `.spec.md` dentro del issue (ver `Documentacion/README.md`).
