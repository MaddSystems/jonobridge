# Explicación de la Modularidad del Backend: Arquitectura como Lego Bricks

## Introducción

El backend del proyecto Grule ha sido diseñado con un enfoque **modular como "lego bricks"**, donde cada carpeta representa un **bloque funcional independiente**. Esta arquitectura permite:

- **Intercambiabilidad**: Reemplazar o actualizar un módulo sin afectar otros.
- **Extensibilidad**: Agregar nuevas funcionalidades como nuevos bricks.
- **Mantenibilidad**: Cada módulo tiene responsabilidades claras y acopladas de manera laxa.
- **Testabilidad**: Probar módulos individualmente.
- **Reutilización**: Usar bricks en diferentes contextos o proyectos.

El ensamblaje ocurre en `main.go`, donde se registran y conectan los bricks a través de interfaces. A continuación, se explica cada carpeta y su rol en esta arquitectura.

## Carpetas Raíz y Punto de Entrada

### `main.go`
**Propósito**: Punto de entrada principal que ensambla todos los bricks.
- Inicializa el registro de capacidades.
- Configura MQTT, HTTP server y worker GRULE.
- Registra cada capacidad en el registry para que estén disponibles en el contexto de reglas.
- **Modularidad**: Actúa como el "pegamento" que conecta bricks sin lógica propia.

### `go.mod`, `Dockerfile`, `build.sh`
**Propósito**: Configuración del proyecto y despliegue.
- `go.mod`: Dependencias independientes del backend.
- `Dockerfile`: Contenedor standalone.
- `build.sh`: Script de construcción.
- **Modularidad**: Permiten desplegar el backend como un brick completo e independiente.

## `adapters/`
**Propósito**: Adaptadores para diferentes tipos de dispositivos trackers (GPS, IoT, etc.).
- `gps_tracker.go`: Parsea payloads JSON de trackers GPS, convirtiéndolos en `IncomingPacket` para GRULE.
- **Razón de ser**: Abstrae el parsing de datos entrantes, permitiendo soporte para nuevos tipos de trackers sin cambiar el core.
- **Modularidad**: Implementa `TrackerAdapter` interface. Se puede reemplazar por adaptadores Bluetooth, LoRa, etc., como bricks intercambiables.

## `api/`
**Propósito**: Handlers HTTP para la API REST.
- `handlers.go`: Endpoints para recargar reglas, subir archivos, consultar estado.
- **Razón de ser**: Proporciona interfaz externa para gestión de reglas y auditoría.
- **Modularidad**: Los handlers llaman a funciones inyectadas (como `ReloadFunc`), permitiendo composición flexible. Es un brick que se conecta al worker y persistencia.

## `audit/`
**Propósito**: Sistema completo de auditoría y logging de ejecución de reglas.
- `listener.go`: Listener GRULE que captura eventos de reglas.
- `capture.go`: Lógica de captura de snapshots y alertas.
- `db.go`: Persistencia de auditoría en MySQL.
- `types.go`: Estructuras de datos para auditoría.
- `manifest.go`: Manejo de manifests declarativos.
- `snapshot.go`: Extracción de snapshots.
- **Razón de ser**: Proporciona trazabilidad completa de decisiones de reglas sin modificar las reglas mismas.
- **Modularidad**: Implementa `SnapshotProvider` para contribuir datos a auditoría. Se puede deshabilitar o reemplazar por bricks alternativos (e.g., logging a Kafka).

### Refactorización del Sistema de Auditoría (Enero 2026)
Se ha implementado una estrategia de **"Captura Explícita Pura"** para eliminar duplicados y ruido:
1.  **Sin Captura Automática**: Se eliminaron los listeners automáticos y las capturas post-ejecución en el worker.
2.  **Control Manual**: Las reglas GRL deciden cuándo capturar usando `actions.CaptureSnapshot("RuleName")`.
3.  **SnapshotProvider**: Las capabilities exponen sus datos implementando esta interfaz, y el sistema de auditoría los agrega automáticamente al snapshot.

## `capabilities/`
**Propósito**: Núcleo modular del sistema. Cada subcarpeta es un brick funcional.
- `interface.go`: Define `Capability` interface que todos los bricks deben implementar.
- `registry.go`: Registry que registra y construye el DataContext GRULE con todos los bricks.
- **Razón de ser**: Permite que reglas accedan a funcionalidades específicas (alertas, buffers, etc.) de manera desacoplada.
- **Modularidad**: Cada capacidad es un brick independiente con su propio `manifest.yaml` para auto-descripción. Se registran en runtime, permitiendo configuración dinámica.

### Estrategia de Diseño para Capabilities
Las capabilities implementan el **patrón Strategy** combinado con un **Registry** para lograr modularidad extrema:

- **Patrón Strategy**: Cada capability implementa la interface `Capability`, definiendo métodos como `Name()`, `Initialize()`, `GetSnapshot()`. Esto permite que diferentes estrategias (e.g., buffer circular vs. buffer lineal) sean intercambiables sin cambiar el código cliente (reglas GRULE).
- **Registry Pattern**: El `Registry` actúa como un contenedor que registra capabilities en runtime. En `main.go`, se registran bricks específicos (alerts, buffer, etc.), y el registry construye el `DataContext` GRULE inyectando cada capability bajo su nombre de contexto.
- **Auto-descripción**: Cada capability incluye un `manifest.yaml` que describe sus funciones, parámetros y ejemplos, permitiendo generación automática de esquemas JSON para herramientas externas (como LLMs para generación de reglas).
- **Beneficios**: Permite composición dinámica; por ejemplo, un sistema de detección de jammer puede usar solo `buffer` y `alerts`, mientras que uno de logística añade `geofence` y `timing`. Nuevas capabilities se agregan creando una subcarpeta y registrándolas, sin modificar código existente.

### `capabilities/alerts/`
**Brick: Gestión de Alertas**
- `capability.go`: Implementa envío de alertas con deduplicación.
- `spam_guard.go`: Previene alertas duplicadas.
- `channels.go`: Manejo de canales de notificación.
- **Razón de ser**: Maneja notificaciones de eventos críticos (jammers, etc.).
- **Modularidad**: Brick que se registra en registry y aporta métodos como `MarkAlertSentForRule()` a reglas.

### `capabilities/buffer/`
**Brick: Buffer Circular de Paquetes**
- `capability.go`: Gestión del buffer de posiciones GPS.
- `circular.go`: Implementación del buffer circular.
- `manager.go`: Lógica de actualización y consulta.
- **Razón de ser**: Mantiene historial de posiciones para detectar patrones (e.g., movimiento).
- **Modularidad**: Brick que aporta estado persistente por IMEI, intercambiable con otros tipos de buffers.

### `capabilities/geofence/`
**Brick: Geofencing**
- `capability.go`: Verificación de posiciones dentro/fuera de zonas.
- `functions.go`: Funciones auxiliares para cálculos geoespaciales.
- **Razón de ser**: Detecta si un vehículo está en zonas seguras o peligrosas.
- **Modularidad**: Brick que consulta geofences desde persistencia, extensible a nuevos tipos de zonas.

### `capabilities/metrics/`
**Brick: Métricas y Estadísticas**
- `capability.go`: Cálculo de promedios y métricas.
- `averages.go`: Funciones para rolling averages.
- **Razón de ser**: Proporciona datos estadísticos para reglas (e.g., velocidad promedio).
- **Modularidad**: Brick que acumula métricas en memoria, reemplazable por bricks con persistencia externa.

### `capabilities/timing/`
**Brick: Temporización y Offline Detection**
- `capability.go`: Seguimiento de timestamps y detección de offline.
- `offline.go`: Lógica para detectar dispositivos desconectados.
- **Razón de ser**: Maneja aspectos temporales como timeouts y estados offline.
- **Modularidad**: Brick que usa timers, extensible a bricks con calendarios o zonas horarias.

## `grule/`
**Propósito**: Integración con el motor de reglas GRULE.
- `worker.go`: Procesa payloads, ejecuta reglas en loop ordenado.
- `context_builder.go`: Construye el DataContext con capacidades registradas.
- `packet.go`: Definición de `IncomingPacket` y estructuras.
- **Razón de ser**: Puente entre datos entrantes y ejecución de reglas.
- **Modularidad**: El worker inyecta el registry de capacidades en el contexto, permitiendo que reglas usen bricks dinámicamente.

## `persistence/`
**Propósito**: Capa de abstracción para almacenamiento de estado.
- `interface.go`: Define `StateStore` interface.
- `mysql.go`: Implementación MySQL.
- `memory.go`: Implementación en memoria para testing.
- `rules.go`: Gestión de reglas en DB.
- **Razón de ser**: Abstrae el almacenamiento, permitiendo cambiar de MySQL a Redis sin afectar otros bricks.
- **Modularidad**: Implementa interfaces, por lo que bricks como geofence pueden consultar datos sin conocer el backend de storage.

## `schema/`
**Propósito**: Generador de esquemas JSON para capacidades.
- `generator.go`: Lee `manifest.yaml` de cada capacidad y genera esquema JSON.
- **Razón de ser**: Proporciona documentación automática de APIs de capacidades para herramientas externas (e.g., LLM generation).
- **Modularidad**: Escanea capacidades registradas, actuando como un brick meta que describe otros bricks.

## Conclusión

Esta arquitectura permite construir sistemas complejos ensamblando bricks simples. Por ejemplo, para agregar una nueva capacidad (como "weather"), solo se crea una nueva subcarpeta en `capabilities/` con su `manifest.yaml`, se implementa la interface, y se registra en `main.go`. El resto del sistema permanece intacto, demostrando la verdadera modularidad como lego bricks.