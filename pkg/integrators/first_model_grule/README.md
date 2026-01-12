# Grule Engine Universal - Sistema de Auditoría

Sistema de auditoría universal para Grule Rule Engine que captura ejecuciones en tiempo real sin código hardcoded.

## 🎯 Características Principales

- **✅ Auditoría Universal**: `actions.RecordExecution()` funciona con cualquier regla GRL
- **✅ Sin Duplicados**: 1 fila por IMEI en tabla `alert_summary`
- **✅ Captura en Tiempo Real**: Datos capturados DURANTE la ejecución (no post-mortem)
- **✅ Thread-Safe**: `sync.RWMutex` para concurrencia segura
- **✅ API REST + Frontend**: Backend Go + Frontend Flask con Bootstrap 5

## 📁 Estructura del Proyecto

```
newgrule/
├── main.go                    # API REST (endpoints /api/audit/*)
├── go.mod                     # Dependencias Go
├── Dockerfile                 # Multi-stage build
├── k8s-deployment.yaml        # Deployment + Service + Ingress
├── deploy.sh                  # Script de deployment automatizado
├── engine/
│   ├── grule_worker.go        # Motor de ejecución de reglas
│   ├── rule_loader.go         # Carga de reglas desde MySQL
│   ├── property.go            # Sistema de properties
│   ├── memory_buffer.go       # Buffer en memoria
│   ├── persistent_state.go    # Estado persistente
│   └── audit/
│       ├── types.go           # RuleExecution, IMEISummary, AlertDetail
│       ├── capture.go         # ExecutionCapture con StartCapture/FinishCapture
│       └── db.go              # SaveExecutions, GetIMEISummaries, GetAlertDetails
├── actions/
│   ├── actions.go             # ActionsHelper con RecordExecution()
│   ├── alerts.go              # SendTelegram, SendEmail, Log
│   └── commands.go            # CutEngine, RestoreEngine, SendRawHex
├── external-web/
│   ├── main.py                # Flask app con proxy API
│   ├── requirements.txt       # Flask, requests
│   └── templates/
│       ├── base.html          # Layout base con Bootstrap 5
│       ├── index.html         # Dashboard principal
│       ├── audit_summary.html # Resumen de IMEIs (sin duplicados)
│       └── audit_details.html # Detalles de alertas por IMEI
└── rules_templates/
    ├── speed_alert.grl        # Regla de velocidad con RecordExecution()
    └── jammer_detection.grl   # Reglas de jamming con RecordExecution()
```

## 🗄️ Base de Datos

### Tabla: `alert_summary`
```sql
CREATE TABLE IF NOT EXISTS alert_summary (
    imei VARCHAR(20) PRIMARY KEY,
    last_alert_date DATETIME,
    total_alerts_24h INT DEFAULT 0,
    alert_types JSON,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Tabla: `alert_details`
```sql
CREATE TABLE IF NOT EXISTS alert_details (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    imei VARCHAR(20),
    rule_name VARCHAR(100),
    description TEXT,
    salience INT,
    status VARCHAR(20),
    timestamp DATETIME,
    duration_ms INT,
    conditions JSON,
    actions JSON,
    alert_fired BOOLEAN,
    INDEX idx_imei_timestamp (imei, timestamp),
    INDEX idx_alert_fired (alert_fired)
);
```

**Optimización clave**: Solo se insertan filas en `alert_details` cuando `alertFired = true`.

## 🚀 API REST Endpoints

### `GET /api/audit/summary?limit=100`
Retorna lista de IMEIs con alertas (sin duplicados)

**Respuesta:**
```json
{
  "success": true,
  "count": 5,
  "data": [
    {
      "imei": "123456789012345",
      "last_alert_date": "2025-12-10T17:30:00Z",
      "total_alerts_24h": 12,
      "alert_types": ["SpeedAlert", "EnviarAlertaTelegram"]
    }
  ]
}
```

### `GET /api/audit/details?imei=123456789012345&limit=100`
Retorna historial de alertas para un IMEI específico

**Respuesta:**
```json
{
  "success": true,
  "imei": "123456789012345",
  "count": 12,
  "data": [
    {
      "id": 1001,
      "rule_name": "SpeedAlert",
      "description": "Vehículo en movimiento detectado",
      "salience": 100,
      "timestamp": "2025-12-10T17:30:00Z",
      "conditions": {"Speed": 85, "Latitude": 19.432, "Longitude": -99.133},
      "actions": ["SendTelegram", "Log"],
      "alert_fired": true,
      "duration_ms": 15
    }
  ]
}
```

### `GET /api/health`
Health check del servicio

## 📝 Ejemplo de Regla GRL

```groovy
rule SpeedAlert "Alerta de Velocidad" salience 100 {
    when
        Jono.Speed > 0 &&
        Jono.Latitude != 0.0
    then
        // Sistema UNIVERSAL de captura
        actions.RecordExecution(
            "SpeedAlert",                                    // ruleName
            "Vehículo en movimiento detectado",            // description
            100,                                            // salience
            map[string]interface{}{                        // conditions
                "Speed": Jono.Speed,
                "Latitude": Jono.Latitude,
                "IMEI": Jono.IMEI,
            },
            []string{"SendTelegram", "Log"},              // actions
            true                                           // alertFired
        );
        
        actions.SendTelegram("🚗 Alerta: " + Jono.IMEI);
        Retract("SpeedAlert");
}
```

## 🐳 Docker Build

```bash
# Build local
docker build -t grule-universal:latest .

# Build y push a registry
docker build -t your-registry/grule-universal:latest .
docker push your-registry/grule-universal:latest
```

## ☸️ Kubernetes Deployment

### 1. Crear secrets
```bash
kubectl create secret generic grule-secrets -n testgrule \
  --from-literal=mysql-host=mysql.default.svc.cluster.local:3306 \
  --from-literal=mysql-user=grule_user \
  --from-literal=mysql-password=your_password \
  --from-literal=telegram-bot-token=your_bot_token \
  --from-literal=mqtt-broker-host=tcp://mqtt:1883
```

### 2. Deployment automático
```bash
# Editar deploy.sh con tu registry
REGISTRY="your-registry"

# Ejecutar deployment
./deploy.sh
```

### 3. Acceso local
```bash
# Port forwarding
kubectl port-forward -n testgrule deployment/grule-universal 8080:8080 5000:5000

# Backend API: http://localhost:8080
# Frontend Web: http://localhost:5000
```

## 🧪 Testing Local

### Backend Go
```bash
# Compilar
go build -o grule-engine

# Configurar variables
export GRULE_AUDIT_ENABLED=Y
export GRULE_AUDIT_LEVEL=ALL
export MYSQL_HOST=localhost:3306
export MYSQL_USER=root
export MYSQL_PASSWORD=password

# Ejecutar
./grule-engine
```

### Frontend Flask
```bash
cd external-web

# Instalar dependencias
pip install -r requirements.txt

# Configurar
export GRULE_API_URL=http://localhost:8080
export FLASK_PORT=5000

# Ejecutar
python main.py
```

## 📊 Diferencias con Sistema Anterior

| Aspecto | Sistema Anterior | Sistema Universal |
|---------|------------------|-------------------|
| Captura | Post-mortem (buildExecutionSteps) | Durante ejecución (RecordExecution) |
| Código | Hardcoded para cada regla | Universal para todas las reglas |
| Base de datos | Miles de duplicados IMEI | 1 fila por IMEI (summary) |
| Agregar reglas | Modificar código Go | Solo crear archivo .grl |
| Precisión | Muestra reglas incorrectas | Muestra ejecución real |
| Tablas | rule_executions, execution_steps, execution_context | alert_summary, alert_details |

## 🔧 Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `GRULE_AUDIT_ENABLED` | Activar auditoría (Y/N) | `Y` |
| `GRULE_AUDIT_LEVEL` | Nivel de auditoría (ALL/ALERT_ONLY) | `ALL` |
| `API_PORT` | Puerto del backend Go | `8080` |
| `FLASK_PORT` | Puerto del frontend Flask | `5000` |
| `GRULE_API_URL` | URL del backend para Flask | `http://localhost:8080` |
| `MYSQL_HOST` | Host de MySQL | `localhost:3306` |
| `MYSQL_USER` | Usuario de MySQL | `root` |
| `MYSQL_PASSWORD` | Password de MySQL | - |
| `MYSQL_DATABASE` | Base de datos | `grule` |
| `TELEGRAM_BOT_TOKEN` | Token de Telegram | - |
| `MQTT_BROKER_HOST` | Host del broker MQTT | `tcp://localhost:1883` |

## 📈 Métricas

- **Binary size**: 14 MB (optimizado con `-ldflags="-w -s"`)
- **Memory**: 256 Mi (request), 512 Mi (limit)
- **CPU**: 250m (request), 500m (limit)
- **Startup time**: ~10 segundos

## 🎨 Frontend Features

- ✅ Bootstrap 5 con diseño moderno
- ✅ Actualización automática cada 15-30 segundos
- ✅ Modal para ver condiciones JSON completas
- ✅ Tabla responsive con búsqueda y filtros
- ✅ Health check visual del backend
- ✅ Sin duplicados de IMEIs en vista summary

## 📞 Soporte

Para agregar nuevas reglas:
1. Crear archivo `.grl` en `rules_templates/`
2. Incluir llamada a `actions.RecordExecution()` en la sección `then`
3. Subir regla a MySQL tabla `fleet_rules`
4. No se requiere modificar código Go

---

**Autor**: Sistema desarrollado con arquitectura universal para eliminar hardcoding  
**Versión**: 2.0 (Universal Audit System)  
**Fecha**: Diciembre 2025


Tablas de auditoria:


curl -X 'GET' \
  'https://jonobridge.madd.com.mx/grule/api/audit/progress/timeline?imei=864352046167177&rule_name=Jammer%20Real%20-%20Detecci%C3%B3n%20Avanzada%20con%20Buffer%20Circular&limit=500' \
  -H 'accept: application/json'


SELECT id, rule_id, rule_name, components_executed, component_details, step_number, stage_reached, stop_reason, buffer_size, metrics_ready, geofence_eval, context_snapshot, execution_time 
FROM rule_execution_state 
WHERE imei = '864352046167177' 
  AND rule_name = 'Jammer Real - Detección Avanzada con Buffer Circular' 
ORDER BY execution_time ASC 
LIMIT 500;

1. Endpoint para el Grid Principal (Resumen)
El primer jqGrid (#imeis-grid) usa el endpoint /summary.

Endpoint: /summary

Propósito: Proporcionar una vista agregada y resumida del progreso de las reglas por cada IMEI. Es una lista de "IMEIs y su último estado".

Consulta: Se activa al presionar "Cargar IMEIs".

Ejemplo: .../progress/summary?rule_name=MiRegla

2. Endpoint para el Modal (Detalle/Timeline)
El segundo jqGrid (#timeline-grid) dentro del modal usa el endpoint /timeline.

Endpoint: /timeline

Propósito: Proporcionar la secuencia detallada y cronológica de todos los frames (toda la historia) para un IMEI específico.

Consulta: Se activa al hacer clic en una fila del primer grid.

Ejemplo: .../progress/timeline?imei=860...&limit=1000&rule_name=MiRegla

SELECT 
    imei,
    rule_name,
    MAX(step_number) AS max_step,
    COUNT(*) AS total_frames,
    MAX(execution_time) AS last_frame_time
FROM rule_execution_state
WHERE rule_name = 'Jammer Real - Detección Avanzada con Buffer Circular'
GROUP BY imei, rule_name
ORDER BY last_frame_time DESC
LIMIT 100;

memory_debug

----
Sistema de Detección de Jammer "DEFCON"

  Este documento describe el funcionamiento del sistema de detección de Jammer basado en reglas, que sigue un flujo secuencial de "DEFCONs"
  (Niveles de Condición de Defensa) para determinar si se debe disparar una alerta. Cada paso en la secuencia es una regla individual que se
  audita, permitiendo una trazabilidad completa del proceso de decisión.

  Flujo de Detección (DEFCON 1-6)

  El sistema evalúa cada trama de datos de un vehículo a través de una serie de reglas, donde cada una representa un "DEFCON". Una trama debe pasar
  exitosamente cada DEFCON para avanzar al siguiente. Si una condición no se cumple, la secuencia se detiene y se audita la razón exacta, evitando
  falsos positivos y proporcionando una visibilidad clara del proceso.

  La secuencia lógica es la siguiente:

   1. DEFCON 0 (Preparación):
       * Componente: UpdateCircularBuffer
       * Acción: Se ejecuta incondicionalmente para cada trama. Limpia los flags de la ejecución anterior y actualiza el buffer de memoria con las
         últimas 10 posiciones del vehículo.
       * Auditoría: Component1_UpdateCircularBuffer.

   2. DEFCON 1 (Puerta de Entrada):
       * Componente: EvaluateOnlyOnInvalid
       * Condición: Verifica si la trama actual tiene una posición GPS inválida Y si el buffer de memoria ya está lleno (10 posiciones).
       * Resultado:
           * PASA: Si ambas condiciones son verdaderas, avanza a DEFCON 2.
           * FALLA: Si la trama es válida o el buffer está incompleto, la secuencia se detiene.
       * Auditoría: Component2_EvaluateOnlyOnInvalid (PASA), Jammer_Stop_Valid (FALLA), o Jammer_Stop_Buffer (FALLA).

   3. DEFCON 2 (Chequeo Offline):
       * Componente: CheckOfflineStatus
       * Condición: Verifica si el vehículo ha estado sin reportar una posición válida durante al menos 5 minutos.
       * Resultado:
           * PASA: Si está offline, avanza a DEFCON 3.
           * FALLA: Si no cumple el tiempo offline, la secuencia se detiene.
       * Auditoría: Component3_CheckOfflineStatus (PASA) o Jammer_Stop_Offline (FALLA).

   4. DEFCON 3 (Cálculo de Métricas):
       * Componente: CalculateBufferMetrics
       * Acción: Calcula la velocidad promedio de los últimos 90 minutos y el nivel de señal GSM promedio de las últimas 5 posiciones almacenadas
         en el buffer.
       * Resultado: Siempre avanza a DEFCON 4 para que las métricas puedan ser evaluadas.
       * Auditoría: Component4_CalculateBufferMetrics.

   5. DEFCON 4 (Verificación de Umbrales):
       * Componente: CheckMetricThresholds
       * Condición: Compara las métricas calculadas con los umbrales predefinidos (Velocidad >= 10 km/h y Señal GSM >= 9).
       * Resultado:
           * PASA: Si ambas métricas cumplen, avanza a DEFCON 5.
           * FALLA: Si alguna de las métricas no cumple, la secuencia se detiene.
       * Auditoría: Component5_CheckMetricThresholds (PASA) o Jammer_Stop_Thresholds (FALLA).

   6. DEFCON 5 (Verificación de Geocercas):
       * Componente: CheckGeofenceExclusion
       * Condición: Verifica que el vehículo no se encuentre dentro de ninguna zona segura predefinida (Taller, CLIENTES, Resguardo).
       * Resultado:
           * PASA: Si está fuera de todas las zonas seguras, avanza a DEFCON 6.
           * FALLA: Si está dentro de alguna zona segura, la secuencia se detiene.
       * Auditoría: Component6_CheckGeofenceExclusion (PASA) o Jammer_Stop_Geofence (FALLA).

   7. DEFCON 6 (Alerta Final):
       * Componente: Disparo de Alerta.
       * Condición: Se activa si todos los DEFCONs anteriores (1 al 5) han pasado.
       * Acción: Envía una notificación de alerta detallada a través de Telegram y marca la alerta como "enviada" para evitar duplicados.
       * Auditoría: Component7_FireAlert.

  Cómo Interpretar la Auditoría ("Movie Frames")

  En la interfaz de "Progress Audit", cada trama que inicia una evaluación de Jammer generará una serie de entradas. Para entender el proceso,
  debes observar la columna rule_name:

   * Verás una secuencia de entradas como Component1_..., Component2_..., etc., para cada DEFCON que se haya cumplido.
   * Si la secuencia llega hasta Component7_FireAlert, significa que se detectó un Jammer y se envió la alerta.
   * Si la secuencia se detiene antes, la última entrada será Jammer_Stop_.... El nombre de esta regla te indicará exactamente por qué el sistema
     decidió que no era un evento de Jammer (ej. Jammer_Stop_Valid, Jammer_Stop_Thresholds, etc.).

  Esta trazabilidad permite un análisis preciso del comportamiento del sistema y facilita la depuración y el ajuste de los umbrales.


 El problema: La última línea es Changed("Jono");. Esto es una orden directa al motor de reglas que dice: "¡Atención! Los datos de Jono
            * El motor obedece y re-evalúa todo.
             ¿Por qué falla y por qué es un bucle?
      La razón exacta, como te comenté en el turno anterior, está en la regla DEFCON0_Check_Valid que modificamos. Vamos a seguir el flujo exacto que
  causa el bucle:

       * Ciclo 1:
           * `when`: La condición Jono.BufferUpdated && !Jono.CurrentlyInvalid se cumple (el buffer se actualizó y la trama es válida).
           * `then`: Se ejecuta el bloque de la regla.
           * El problema: La última línea es Changed("Jono");. Esto es una orden directa al motor de reglas que dice: "¡Atención! Los datos de Jono
             han cambiado. Vuelve a revisar todas las reglas desde el principio para ver si ahora se cumple alguna nueva".

       * Ciclo 2:
           * El motor obedece y re-evalúa todo.
           * Vuelve a mirar la regla DEFCON0_Check_Valid.
           * `when`: La condición Jono.BufferUpdated && !Jono.CurrentlyInvalid sigue siendo cierta. No hemos hecho nada en el bloque then para cambiar
             estos valores.
           * `then`: Como la condición se cumple de nuevo, la regla se ejecuta otra vez.
           * El problema: Se vuelve a ejecutar Changed("Jono");, ordenando otra re-evaluación.

       * Ciclos 3 al 100:
           * El proceso se repite idénticamente. La regla se llama a sí misma indirectamente una y otra vez.

       * Ciclo 101:
           * El motor de reglas tiene un fusible de seguridad (MaxCycle = 100). Al ver que ha dado 100 vueltas sin que nada cambie, detiene la
             ejecución, genera el mensaje de error que ves, y la función eng.Execute devuelve ese err.
           * El código en grule_worker.go:415 captura ese err y lo imprime en tu log.

 /model gemini/gemini-3.0-pro     