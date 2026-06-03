# 🎓 Introducción a las APIs Backend (Guía para Devs)

¡Bienvenido al mundo del Backend! Construir una API (Application Programming Interface) es muy distinto a hacer una página web clásica. Aquí no trabajamos con píxeles, colores ni animaciones, sino con **datos, seguridad y rendimiento**.

A continuación, respondemos las preguntas clave para entender cómo funciona este ecosistema.

---

## 1. ¿Tiene UI? ¿Se puede ver en algún lado?

**No, una API no tiene una Interfaz de Usuario (UI) gráfica tradicional.** 
Una API es un puente de comunicación entre dos computadoras. En lugar de devolver una página HTML bonita con CSS, devuelve **datos crudos**, generalmente en formato **JSON** (JavaScript Object Notation).

### Ejemplo de lo que "se ve":
Si pides un personaje a la API de One Piece, verás algo así:
```json
{
  "id": 1,
  "name": "Monkey D. Luffy",
  "crew": "Straw Hat Pirates",
  "bounty": 3000000000
}
```

### ¿Dónde y cómo se puede ver?
Aunque no tenga botones, los desarrolladores "ven" y prueban la API usando herramientas especiales:
1. **Clientes HTTP**: Programas como [Postman](https://www.postman.com/), [Insomnia](https://insomnia.rest/) o [Bruno](https://www.usebruno.com/). Vos ponés la URL (`http://localhost:8080/characters`) y la herramienta te muestra el JSON que responde.
2. **Swagger / OpenAPI**: Es lo más parecido a una UI. Muchas APIs autogeneran una página web de documentación interactiva (Swagger) donde podés ver todos los endpoints (URLs disponibles) y tocando un botón `Try it out`, ejecutar la petición desde el navegador.
3. **En Producción**: El usuario final jamás ve la API. El Frontend (una app web en React o una app móvil en Flutter) consume este JSON y se encarga de dibujar la tarjeta del personaje con su foto y botones.

---

## 2. El Ecosistema de un Desarrollador Backend

El día a día de un backend dev no ocurre en el navegador, ocurre en la consola y el editor. Usamos este ecosistema:

* **Lenguaje y Framework**: En este proyecto usamos **Go** (Golang) y su framework web **Gin**. Go es famoso por ser ridículamente rápido y fácil de compilar.
* **Base de Datos**: Todo debe guardarse en algún lado. Usamos bases de datos relacionales como **PostgreSQL** para la data permanente (personajes, arcos) y en memoria como **Redis** para cosas rápidas (caché, contadores).
* **Docker y Docker Compose**: Para no tener que instalar PostgreSQL y Redis manualmente en tu PC y ensuciar tu sistema, usamos Docker. Levanta todo en "contenedores" aislados con un solo comando (`make docker-up`).
* **Migraciones**: Scripts SQL versionados. En lugar de crear tablas a mano en la base de datos, escribimos scripts (`migrations/`) que el sistema ejecuta, asegurando que la base de datos sea igual en la PC de todos los devs.
* **Testing**: Un backend sin tests es una bomba de tiempo. Se escriben "Unit Tests" (pruebas unitarias) que simulan peticiones para garantizar que el código hace lo que debe.

---

## 3. ¿Cuál es el trabajo de la API? ¿Qué debo verificar?

Imaginá que un restaurante es el software completo. 
- El cliente sentado en la mesa es el **Frontend** (o app móvil).
- La cocina llena de ingredientes y heladeras es la **Base de Datos**.
- **La API es el Mozo**. 

El trabajo de la API es recibir la orden (Request), ir a la cocina a buscarla de forma segura, y devolverle el plato terminado (Response) al cliente.

### ¿Qué tenés que verificar que funcione?
Cuando programes en este repo, debés asegurar 4 cosas:

1. **El Contrato (Formato)**: Si la documentación dice que devuelves un campo `bounty` como número entero, no podés devolverlo como un texto (`"3000"`). Romperías la app del Frontend.
2. **La Lógica de Negocio**: Que las reglas se cumplan. (Ej: Un personaje muerto no puede actualizar su recompensa actual).
3. **Seguridad e Identidad**: ¿El usuario que pide borrar un personaje envió un Token JWT válido? ¿Es Administrador? Si no, rechazar la petición.
4. **Los Códigos de Estado HTTP**: La API debe comunicarse universalmente usando números.
   * `200 OK`: Salió todo bien.
   * `201 Created`: Salió todo bien y el registro se creó (ej: alta de personaje).
   * `400 Bad Request`: El cliente mandó datos inválidos (ej: edad = -5).
   * `401 / 403`: Problemas de autenticación o permisos.
   * `404 Not Found`: Pidió el ID 9999 y no existe.
   * **`500 Internal Server Error`**: Tu código crasheó. **Esto es lo único que nunca debe pasar**.

---

## 4. El Mercado en Buenos Aires (2026)

### ¿Para qué se usan las APIs?
Hoy **todo es una API**. El modelo de aplicaciones "monolíticas" (donde el backend escupía HTML directamente) está obsoleto. 
Hoy las empresas crean una única API central, y esa misma API alimenta a:
- La página web (React/Next.js)
- La app de Android (Kotlin/Jetpack)
- La app de iOS (Swift)
- Integraciones B2B (otros sistemas corporativos que consumen sus datos).

### ¿Qué tecnologías pide el mercado porteño/argentino hoy?
En el ecosistema local (Mercado Libre, Ualá, Globant, Fintechs, Startups), el stack se divide así:

1. **Go (Golang)**: Está en auge masivo. Mercado Libre lo usa como estándar. Se elige por su alta concurrencia (soporta miles de usuarios a la vez) y eficiencia de costos en servidores de AWS.
2. **Node.js (NestJS / TypeScript)**: Sigue reinando en agencias y startups semilla porque permite que el equipo de frontend también meta mano en el backend fácilmente.
3. **Python (FastAPI)**: Cuando la empresa tiene áreas de Datos, IA o Machine Learning, Python es obligatorio.
4. **Java (Spring Boot) / C# (.NET)**: El pan de cada día en bancos tradicionales y corporaciones masivas. Estables pero pesados.

### ¿Cómo se mantienen y evolucionan?
- **Cloud Computing**: Ya nadie tiene servidores físicos. Todo vive en la nube (AWS, Google Cloud).
- **CI/CD (Integración y Despliegue Continuo)**: Cuando un dev hace `merge` a `main`, una computadora en la nube compila el código, corre los tests, y si todo pasa, actualiza la API en vivo sin que los usuarios noten un corte.
- **Observabilidad (Datadog / Grafana)**: Herramientas gigantescas que monitorean la API 24/7. Si el endpoint de login tarda más de 200 milisegundos, o si salta un error 500, suena una alarma en el celular del desarrollador de guardia.
