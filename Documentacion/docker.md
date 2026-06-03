# 🐳 Docker y Alternativas para Bajos Recursos

Este proyecto utiliza **Docker** a través de `docker-compose` para simplificar la configuración del entorno de desarrollo. En lugar de instalar bases de datos en tu computadora y lidiar con diferentes versiones, Docker crea "contenedores" aislados.

## 📦 ¿Para qué usamos Docker en este repo?

Nuestro archivo `docker-compose.yml` levanta los siguientes servicios:
1. **PostgreSQL 16**: Nuestra base de datos relacional principal.
2. **Redis 7**: Base de datos en memoria utilizada para Caché y Rate Limiting.

Al ejecutar `make docker-up`, descargamos y encendemos estas dos bases de datos listas para que la API se conecte a ellas.

---

## 💻 El Problema en Windows (Bajos Recursos)

**Docker Desktop en Windows** utiliza virtualización (WSL2 o Hyper-V) para correr un motor Linux de fondo. Esto consume una cantidad significativa de memoria RAM (a menudo entre 2GB y 4GB solo estando inactivo) y puede volver lentas a las PCs con 8GB de RAM o procesadores antiguos.

Dado que **Go (Golang)** es increíblemente liviano, la API en sí no consume casi nada. El problema es exclusivamente Docker.

---

## 🚀 Cómo prescindir de Docker (Modo Nativo)

Si tu computadora sufre al correr Docker, puedes optar por **no utilizarlo en absoluto**. Para ello, debes instalar las dependencias directamente en tu Windows. A esto se le llama trabajar de forma **Nativa**.

### 1. Instalar PostgreSQL Nativamente
1. Descarga el instalador de [PostgreSQL para Windows](https://www.postgresql.org/download/windows/).
2. Sigue los pasos del instalador. Anota la contraseña que le asignas al superusuario `postgres`.
3. Abre **pgAdmin** (se instala junto a Postgres) o la terminal `psql`.
4. Crea una base de datos llamada `onepiece_db`.
5. Crea un rol/usuario llamado `onepiece` con la contraseña `onepiece_secret`.
   - *Alternativa:* Simplemente cambia tu archivo `.env` en la raíz del proyecto para que la URL de conexión coincida con el usuario y contraseña que configuraste.

### 2. Instalar Redis Nativamente
Redis no soporta Windows oficialmente desde hace tiempo. Sin embargo, existe una alternativa gratuita al 100% compatible diseñada para Windows: **Memurai**.

1. Descarga [Memurai Developer](https://www.memurai.com/).
2. Instálalo. Automáticamente se configurará como un servicio de Windows corriendo en el puerto `6379` (exactamente igual que Redis).
3. ¡Listo! Tu archivo `.env` funcionará sin cambios.

### 3. Modificar el Flujo de Trabajo
Una vez que tengas Postgres y Memurai corriendo como servicios nativos de Windows, ya no necesitas usar los comandos de Docker del `Makefile`.

**Flujo Anterior (con Docker):**
```bash
make docker-up
make migrate-up
make run
```

**Nuevo Flujo (Nativo):**
```bash
# Ya no ejecutas docker-up porque Postgres y Redis arrancan solos al prender tu PC.
make migrate-up
make run
```

### Notas adicionales
- Los tests de regresión (`tests/`) utilizan la librería `testcontainers-go`. Si no tienes Docker, estos tests **fallarán**. Para desarrollo local, simplemente escribe tu código y corre la API. La verificación profunda de los testcontainers la hará el servidor de Integración Continua (CI/CD) como GitHub Actions antes de aceptar tus Pull Requests.
