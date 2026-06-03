# 🏴‍☠️ One Piece API

API REST del universo **One Piece**. Consulta personajes, frutas del diablo, tripulaciones, arcos, episodios y más.

> Construida en **Go 1.23** · **Gin** · **PostgreSQL 16** · **Redis 7** · **Docker**

---

## 📋 Requisitos previos

Antes de comenzar, asegurate de tener instalado:

| Herramienta | Versión mínima | Instalación |
|---|---|---|
| [Go](https://go.dev/dl/) | 1.23+ | `winget install GoLang.Go` |
| [Docker Desktop](https://www.docker.com/products/docker-desktop/) | — | [docker.com](https://www.docker.com) |
| [GitHub CLI (`gh`)](https://cli.github.com/) | 2.x+ | `winget install GitHub.cli` |
| [Make](https://www.gnu.org/software/make/) | — | `winget install GnuWin32.Make` |
| [golang-migrate](https://github.com/golang-migrate/migrate) | — | Ver docs oficiales |
| [sqlc](https://sqlc.dev/) | — | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` |
| [golangci-lint](https://golangci-lint.run/) | — | Ver docs oficiales |

---

## 🚀 Setup inicial

### 1. Clonar el repositorio

```bash
git clone https://github.com/GalvanAlexis/onepiece-api.git
cd onepiece-api
```

### 2. Configurar variables de entorno

```bash
cp .env.example .env
# Editá .env si necesitás cambiar puertos o credenciales
```

### 3. Levantar la infraestructura (PostgreSQL + Redis)

```bash
make docker-up
```

### 4. Ejecutar migraciones

```bash
make migrate-up
```

### 5. Correr la API

```bash
make run
```

La API estará disponible en `http://localhost:8080`.

---

## 🛠️ Comandos disponibles (`make`)

```bash
make build          # Compila el binario en bin/
make run            # Corre la API en modo desarrollo
make clean          # Limpia el directorio bin/

make docker-up      # Levanta PostgreSQL y Redis con Docker
make docker-down    # Detiene y elimina los contenedores
make docker-all     # Levanta todos los servicios (incluye API en Docker)

make migrate-up     # Aplica las migraciones pendientes
make migrate-down   # Revierte la última migración
make migrate-create name=<nombre>  # Crea una nueva migración

make sqlc-gen       # Genera el código Go a partir de las queries SQL

make test           # Corre la suite de tests unitarios de Go con -race
make test-api       # Corre los tests automatizados de la API usando Bruno CLI
make lint           # Corre golangci-lint
make tidy           # Ejecuta go mod tidy
```

---

## 📂 Estructura del proyecto

```
onepiece-api/
├── cmd/
│   └── api/            # Punto de entrada (main.go)
├── internal/
│   ├── config/         # Carga de configuración (viper)
│   ├── domain/         # Entidades y contratos (interfaces de repositorio)
│   ├── handler/        # Controladores HTTP (Gin)
│   ├── middleware/      # Middlewares (auth, rate-limit, logging)
│   ├── repository/     # Implementaciones de acceso a datos (pgx + sqlc)
│   └── usecase/        # Lógica de negocio
├── migrations/         # Archivos SQL de migraciones (golang-migrate)
├── pkg/
│   ├── database/       # Conexión a PostgreSQL y Redis
│   └── logger/         # Logger estructurado (zap)
├── seeds/              # Datos iniciales de One Piece
├── Documentacion/      # Documentación interna del proyecto
├── ISSUES/             # Tracker local de issues y tareas
├── bruno/              # Colección de endpoints y tests de la API (Bruno)
├── docker-compose.yml
├── Makefile
└── sqlc.yaml
```

---

## 🌿 Flujo de trabajo con Git y GitHub CLI

> ⚠️ **`main` es una rama protegida. Jamás hagas `git push origin main` directamente.**

Todo el trabajo se realiza a través de Pull Requests desde ramas de feature.

### Flujo completo paso a paso

#### 1. Asegurarte de estar en `main` actualizado

```bash
git checkout main
git pull origin main
```

#### 2. Crear tu rama de trabajo

El nombre de la rama debe referenciar el número de issue:

```bash
git checkout -b feature/ISS-001
```

#### 3. Hacer tus cambios y commitear

Los mensajes de commit deben estar en **infinitivo** y ser descriptivos:

```bash
git add .
git commit -m "feat: agregar endpoint GET /characters ISS-001"
```

#### 4. Subir tu rama al remoto

```bash
git push origin feature/ISS-001
```

#### 5. Crear el Pull Request con `gh` CLI

```bash
gh pr create \
  --title "feat: agregar endpoint GET /characters [ISS-001]" \
  --body "Closes ISS-001. Implementa el listado paginado de personajes con filtros." \
  --base main \
  --head feature/ISS-001
```

#### 6. Después del merge, limpiar la rama local

```bash
git checkout main
git pull origin main
git branch -d feature/ISS-001
```

### Comandos `gh` útiles

```bash
gh pr list                  # Ver PRs abiertos
gh pr view <número>         # Ver detalle de un PR
gh pr checks <número>       # Ver el estado de los checks de CI
gh issue list               # Listar issues (si se usan en GitHub)
gh repo view --web          # Abrir el repo en el navegador
```

---

## 🏗️ Convenciones de código

Este proyecto sigue **Clean Architecture** en capas:

```
Handler → UseCase → Repository → Database
```

- **`domain/`**: Solo entidades e interfaces. **Sin dependencias externas.**
- **`usecase/`**: Lógica de negocio pura. Depende de interfaces, no de implementaciones.
- **`handler/`**: Solo parsea requests y delega al usecase. Sin lógica de negocio.
- **`repository/`**: Única capa que conoce la base de datos (pgx/sqlc).

### Reglas generales

- Un paquete = una responsabilidad.
- No importar capas internas desde capas externas (respetar el sentido del flujo).
- Todas las queries SQL nuevas deben agregarse en `sqlc/query/` y regenerar con `make sqlc-gen`.
- Migraciones: **siempre** agregar el archivo `.down.sql` correspondiente.

---

## 📖 Documentación adicional

- [`Documentacion/`](./Documentacion/) — Documentación técnica del proyecto
- [`ISSUES/`](./ISSUES/) — Tracker local de tareas e issues

### 🧪 Pruebas de API con Bruno
Usamos [Bruno](https://www.usebruno.com/) como cliente y runner de tests offline-first. La colección completa se encuentra en la carpeta `bruno/`.
- Para utilizar la colección gráficamente, descarga la App de Bruno y abre el directorio `bruno/`.
- Para correr los tests desde la terminal, requiere Node.js: `npx @usebruno/cli run bruno --env local` o usar el atajo `make test-api`.
- Por defecto Bruno v3.0.0 corre en Safe Mode.

---

## 🤝 Contribución

1. Revisá los issues abiertos en [`ISSUES/`](./ISSUES/).
2. Asignate un issue (cambiar estado a `En Progreso`).
3. Seguí el flujo Git descrito arriba.
4. Asegurate de que `make test` y `make lint` pasen **antes** de abrir el PR.
