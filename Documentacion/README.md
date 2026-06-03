# 🤖 Guía del Agente IA — `onepiece-api`

> Este documento es la fuente de verdad del agente **Antigravity** para trabajar en este repositorio.
> Leerlo completo al inicio de cada sesión que involucre cambios de código.

---

## 1. Identidad del Repositorio

| Campo | Valor |
|---|---|
| **Proyecto** | One Piece API |
| **Stack** | Go 1.23 · Gin · PostgreSQL 16 · Redis 7 · Docker |
| **Arquitectura** | Clean Architecture + DDD |
| **Rama protegida** | `main` — **NUNCA hacer push directo** |
| **Tracker de issues** | `ISSUES/` (local, dentro del repo) |
| **Documentación** | `Documentacion/` (local, dentro del repo) |
| **Notion** | ❌ No se usa en este proyecto |

---

## 2. Reglas Críticas (No Negociables)

### 🚫 Git
- **NUNCA** ejecutar `git push origin main`.
- **NUNCA** hacer commits directamente sobre `main`.
- Todo trabajo va en una rama `feature/ISS-XXX`.
- **SIEMPRE** abrir un PR con `gh pr create` antes de mergear.

### 📋 Issues
- Todo trabajo debe estar vinculado a un issue en `ISSUES/`.
- Si no existe el issue para la tarea solicitada, crearlo antes de comenzar.
- Actualizar el estado del issue (`En Progreso` → `En Review` → `Hecho`) a medida que avanza.

### 📝 Documentación
- No usar Notion. Toda documentación vive en `Documentacion/`.
- Las decisiones de arquitectura o patrones nuevos se documentan en `Documentacion/`.

---

## 3. Arranque de Sesión (Checklist)

Al iniciar una sesión de trabajo con código:

```
1. mem_context          → recuperar contexto de sesiones previas
2. Leer ISSUES/         → identificar el issue activo
3. git status           → verificar estado del repositorio
4. Verificar rama       → debe ser feature/ISS-XXX, nunca main
```

Si el agente se encuentra en `main` y hay trabajo pendiente:

```bash
git checkout -b feature/ISS-XXX
```

---

## 4. Ciclo de Trabajo (SDD Adaptado)

### Para bugs simples (≤ 1 archivo, claro y acotado)
```
sequential-thinking → fix → make test → commit → PR
```

### Para features o bugs complejos
```
Explorer    → Leer el código relevante, entender el flujo
Proposer    → Definir el enfoque de solución
Spec Writer → Crear ISSUES/ISS-XXX_descripcion.md con el .spec
Designer    → Definir interfaces/tipos antes de implementar
Implementer → Escribir código ESTRICTAMENTE según el .spec.md
Verifier    → make test · make lint → confirmar que todo pasa
```

> El `.spec.md` dentro de `ISSUES/` es la **única fuente de verdad** para el Implementer.
> Si hay contradicción entre el spec y lo que se pide verbalmente, consultar al usuario.

---

## 5. Estructura de Capas (No Romper)

```
Handler → UseCase → Repository → Database
```

| Capa | Paquete | Regla |
|---|---|---|
| Entidades | `internal/domain/` | Sin dependencias externas |
| Lógica | `internal/usecase/` | Solo depende de interfaces del dominio |
| HTTP | `internal/handler/` | Solo parsea y delega. Sin lógica de negocio |
| Datos | `internal/repository/` | Única capa que conoce pgx/sqlc |
| Infraestructura | `pkg/` | Database connection, logger |

---

## 6. Comandos Esenciales

```bash
# Levantar infra local
make docker-up

# Correr la API
make run

# Migraciones
make migrate-up
make migrate-down
make migrate-create name=<nombre>

# Generar código sqlc (tras modificar queries SQL)
make sqlc-gen

# Quality Gates — deben pasar SIEMPRE antes de cerrar un issue
make test
make lint
```

---

## 7. Flujo Git con `gh` CLI

```bash
# 1. Actualizar main
git checkout main && git pull origin main

# 2. Crear rama del issue
git checkout -b feature/ISS-XXX

# 3. Hacer cambios, luego commitear (infinitivo, descriptivo)
git add .
git commit -m "feat: descripción del cambio ISS-XXX"

# 4. Subir rama
git push origin feature/ISS-XXX

# 5. Abrir PR
gh pr create \
  --title "feat: descripción [ISS-XXX]" \
  --body "Closes ISS-XXX. Descripción de los cambios realizados." \
  --base main \
  --head feature/ISS-XXX
```

---

## 8. Cierre de Sesión (Checklist)

Antes de terminar una sesión con cambios de código:

```
[ ] make test     → todos los tests pasan
[ ] make lint     → sin errores de lint
[ ] git commit    → mensaje descriptivo en infinitivo + referencia ISS
[ ] gh pr create  → PR abierto apuntando a main
[ ] ISSUES/       → estado del issue actualizado
[ ] mem_session_summary → resumen guardado en memoria del agente
```

---

## 9. Gestión de Issues

Los issues viven en `ISSUES/` con el formato:

```
ISSUES/ISS-XXX_descripcion-breve.md
```

Copiar `ISSUES/ISS-000_template.md` para crear uno nuevo.

### Estados válidos

| Estado | Significado |
|---|---|
| `Backlog` | Definido pero no iniciado |
| `En Progreso` | Rama creada, trabajo activo |
| `En Review` | PR abierto, esperando revisión |
| `Hecho` | PR mergeado a main |
| `Cancelado` | Descartado con justificación |

---

## 10. Memoria Persistente

Guardar con `mem_save` (formato What/Why/Where/Learned) tras:
- Bugs resueltos que involucren patrones no obvios
- Decisiones de arquitectura tomadas con el usuario
- Gotchas descubiertos en el stack (pgx, sqlc, Gin, Redis)

Proyecto a usar en `mem_save`: `GalvanAlexis/onepiece-api`
