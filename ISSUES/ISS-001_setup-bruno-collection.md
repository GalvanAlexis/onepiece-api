---
estado: En Progreso
prioridad: Alta
creado: 2026-06-03
---

# Setup de Colección Bruno y CLI Integración

## 📝 Descripción
El proyecto necesita una herramienta moderna y offline-first para probar y documentar la API. Se utilizará **Bruno** como cliente principal. 

Se debe inicializar la colección de Bruno en el repositorio y configurar los comandos CLI (`bru`) necesarios para poder automatizar o ejecutar pruebas contra la API desde la terminal, previendo integraciones futuras de CI/CD.

## ✅ Criterios de Aceptación
- [ ] Existe un directorio para la colección de Bruno en el root del proyecto.
- [ ] Se añade un comando en el `Makefile` para correr tests usando Bruno CLI en modo seguro (`--sandbox=developer` si fuera estrictamente necesario, aunque por defecto usar Safe Mode v3.0.0).
- [ ] La colección está excluida o adaptada correctamente en `.gitignore` (para no pushear variables de entorno sensibles).

## 🏗️ Implementación Técnica
- Ejecutar inicialización de la colección de Bruno.
- Actualizar `Makefile`.
- Añadir README o actualización en documentación sobre cómo ejecutar los tests con `bru cli`.

## 📌 Notas Adicionales
- Documentación base provista: https://docs.usebruno.com/llms.txt
- Recordatorio de Bruno CLI v3.0.0: Safe Mode por defecto.
