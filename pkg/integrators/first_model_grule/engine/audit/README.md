# Sistema de Auditoría Universal - Grule Engine

Sistema de auditoría thread-safe que captura ejecuciones de reglas en tiempo real sin código hardcoded.

## 🎯 Características

- **✅ Captura Universal**: Funciona con cualquier regla GRL mediante `actions.RecordExecution()`
- **✅ Sin Duplicados**: 1 fila por IMEI en tabla `alert_summary`
- **✅ Tiempo Real**: Captura datos DURANTE la ejecución (no post-mortem)
- **✅ Thread-Safe**: `sync.RWMutex` para concurrencia segura
- **✅ Filtrado Inteligente**: Solo guarda alertas (`AlertFired = true`)

---

## 🗄️ Tablas de Base de Datos

### 1. **`alert_summary`** - Dashboard Ejecutivo (sin duplicados)

**Propósito**: Vista consolidada de alertas por IMEI (1 fila por dispositivo)

```sql
CREATE TABLE IF NOT EXISTS alert_summary (
    imei VARCHAR(20) PRIMARY KEY,                  -- SIN DUPLICADOS
    last_alert_date DATETIME(6) NOT NULL,          -- Última alerta
    total_alerts_24h INT DEFAULT 0,                -- Conteo últimas 24h
    alert_types JSON,                              -- {"SpeedAlert": 5, "Jammer": 2}
    last_rule_executed VARCHAR(100),               -- Última regla ejecutada
    last_alert_location VARCHAR(100),              -- "lat,lon"
    updated_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_last_alert (last_alert_date),
    INDEX idx_total_alerts (total_alerts_24h)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**Operaciones:**
- `INSERT ... ON DUPLICATE KEY UPDATE` → Actualiza en lugar de duplicar
- Se actualiza automáticamente con cada nueva alerta
- Conteo de alertas recalculado desde últimas 24 horas

**Funciones que escriben:**
- `SaveExecutions()` → `updateSummary()` (db.go:139)

---

### 2. **`alert_details`** - Historial Completo de Alertas

**Propósito**: Registro detallado de cada alerta disparada (auditoría completa)

```sql
CREATE TABLE IF NOT EXISTS alert_details (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    imei VARCHAR(20) NOT NULL,
    alert_date DATETIME(6) NOT NULL,
    rule_name VARCHAR(100) NOT NULL,
    rule_description VARCHAR(255),
    salience INT,
    conditions_snapshot JSON,                      -- Valores evaluados (Speed, GPS, etc.)
    actions_executed JSON,                         -- ["SendTelegram", "CutEngine"]
    telegram_sent BOOLEAN DEFAULT false,
    latitude DECIMAL(10, 6),
    longitude DECIMAL(10, 6),
    speed INT,
    created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6),
    INDEX idx_imei_date (imei, alert_date),
    INDEX idx_rule_name (rule_name),
    INDEX idx_alert_date (alert_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**Operaciones:**
- Solo inserta cuando `AlertFired = true`
- Una fila por cada alerta disparada
- Soporta múltiples alertas por IMEI (historial completo)

**Funciones que escriben:**
- `SaveExecutions()` (db.go:93)

---

## 📊 Flujo de Captura

```
┌─────────────────────────────────────────────────────────────┐
│ 1. INICIO EJECUCIÓN                                         │
│    audit.StartCapture(imei) → ExecutionCapture en memoria   │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. EJECUCIÓN DE REGLAS (grule_worker.go)                   │
│    Cada regla GRL ejecuta:                                  │
│    actions.RecordExecution(ruleName, desc, salience,        │
│                            conditions, actions, alertFired)  │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. CAPTURA EN MEMORIA (capture.go)                         │
│    ExecutionCapture.RecordExecution(exec)                   │
│    → Guarda en []RuleExecution (thread-safe)               │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. FINALIZACIÓN (audit.FinishCapture)                      │
│    → Verificar GRULE_AUDIT_ENABLED=Y                        │
│    → Aplicar filtro GRULE_AUDIT_LEVEL                       │
│    → Llamar SaveExecutions(imei, executions)                │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. PERSISTENCIA EN MySQL (db.go)                           │
│    a) Filtrar solo alertFired=true                          │
│    b) INSERT INTO alert_details (cada alerta)              │
│    c) INSERT ... ON DUPLICATE KEY UPDATE alert_summary     │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔧 Configuración

### Variables de Entorno

```bash
# Habilitar/Deshabilitar auditoría
GRULE_AUDIT_ENABLED=Y              # Y/N

# Nivel de captura
GRULE_AUDIT_LEVEL=ALL              # ALL, ERROR, NONE
# - ALL: Captura todas las ejecuciones
# - ERROR: Solo captura cuando hay alerta (alertFired=true)
# - NONE: No captura nada
```

### Ejemplo de Inicialización (main.go)

```go
import (
    "github.com/jonobridge/grule-integrator/engine"
    "github.com/jonobridge/grule-integrator/engine/audit"
)

func main() {
    // Inicializar engine
    engine.Initialize()
    
    // Inicializar auditoría si está habilitada
    if os.Getenv("GRULE_AUDIT_ENABLED") == "Y" {
        db := engine.GetDB()
        audit.InitDB(db)
        log.Println("✅ Sistema de auditoría inicializado")
    }
}
```

---

## 📝 Uso en Reglas GRL

### Sintaxis Completa

```groovy
rule SpeedAlert "Alerta de Velocidad" salience 100 {
    when
        pw.Speed > 120
    then
        // 1. Ejecutar acciones reales
        actions.SendTelegram("Velocidad excedida: " + pw.Speed + " km/h");
        actions.Log("ALERT: Speed=%d, IMEI=%s", pw.Speed, pw.IMEI);
        
        // 2. Registrar auditoría (UNIVERSAL)
        actions.RecordExecution(
            "SpeedAlert",                               // ruleName
            "Alerta de velocidad > 120 km/h",          // description
            100,                                        // salience
            {
                "Speed": pw.Speed,                      // Valores REALES evaluados
                "Latitude": pw.Latitude,
                "Longitude": pw.Longitude,
                "IMEI": pw.IMEI,
                "EventCode": pw.EventCode
            },
            ["SendTelegram", "Log"],                    // Acciones ejecutadas
            true                                        // alertFired (true=alerta)
        );
        
        Retract("SpeedAlert");
}
```

### Parametros de `RecordExecution()`

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `ruleName` | `string` | Nombre único de la regla |
| `description` | `string` | Descripción legible |
| `salience` | `int` | Prioridad de ejecución |
| `conditions` | `map[string]interface{}` | Valores evaluados (Speed, GPS, etc.) |
| `actions` | `[]string` | Lista de acciones ejecutadas |
| `alertFired` | `bool` | `true` = alerta crítica, `false` = solo log |

---

## 🔍 Consultas SQL Útiles

### Ver IMEIs con más alertas (últimas 24h)
```sql
SELECT imei, total_alerts_24h, last_rule_executed, last_alert_date
FROM alert_summary
ORDER BY total_alerts_24h DESC
LIMIT 20;
```

### Historial de alertas por IMEI
```sql
SELECT alert_date, rule_name, rule_description, 
       conditions_snapshot, actions_executed
FROM alert_details
WHERE imei = '123456789012345'
ORDER BY alert_date DESC
LIMIT 50;
```

### Alertas por tipo de regla
```sql
SELECT rule_name, COUNT(*) as total
FROM alert_details
WHERE alert_date >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
GROUP BY rule_name
ORDER BY total DESC;
```

### Limpiar alertas antiguas (>30 días)
```sql
DELETE FROM alert_details
WHERE alert_date < DATE_SUB(NOW(), INTERVAL 30 DAY);
```

---

## 📡 API REST Endpoints

### `GET /api/audit/summary?limit=100`
Retorna lista de IMEIs con alertas (sin duplicados)

**Respuesta:**
```json
{
  "success": true,
  "data": [
    {
      "imei": "123456789012345",
      "last_alert_date": "2025-12-10T15:30:00.123456Z",
      "total_alerts_24h": 15,
      "alert_types": {"SpeedAlert": 10, "JammerDetected": 5},
      "last_rule_executed": "SpeedAlert",
      "last_alert_location": "40.416775,-3.703790"
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
  "data": [
    {
      "id": 12345,
      "imei": "123456789012345",
      "alert_date": "2025-12-10T15:30:00.123456Z",
      "rule_name": "SpeedAlert",
      "rule_description": "Velocidad > 120 km/h",
      "salience": 100,
      "conditions": {
        "Speed": 135,
        "Latitude": 40.416775,
        "Longitude": -3.703790
      },
      "actions": ["SendTelegram", "Log"],
      "telegram_sent": true,
      "latitude": 40.416775,
      "longitude": -3.703790,
      "speed": 135
    }
  ]
}
```

---

## 🧪 Testing

### 1. Verificar tablas creadas
```bash
mysql -u root -p grule -e "SHOW TABLES LIKE 'alert_%';"
```

### 2. Simular alerta
```go
// En test_rule.grl
rule TestAlert "Test de Auditoría" salience 999 {
    when
        pw.Speed > 0
    then
        actions.RecordExecution(
            "TestAlert",
            "Test de captura universal",
            999,
            {"Speed": pw.Speed, "IMEI": pw.IMEI},
            ["Log"],
            true  // Disparar alerta
        );
        Retract("TestAlert");
}
```

### 3. Verificar inserción
```sql
SELECT * FROM alert_summary WHERE imei = '123456789012345';
SELECT * FROM alert_details WHERE imei = '123456789012345' ORDER BY alert_date DESC LIMIT 1;
```

---

## 🚀 Optimizaciones

### Índices Creados

**`alert_summary`:**
- `PRIMARY KEY (imei)` → Búsquedas rápidas por IMEI
- `INDEX idx_last_alert (last_alert_date)` → Ordenar por fecha
- `INDEX idx_total_alerts (total_alerts_24h)` → Top alertas

**`alert_details`:**
- `INDEX idx_imei_date (imei, alert_date)` → Historial por IMEI
- `INDEX idx_rule_name (rule_name)` → Estadísticas por regla
- `INDEX idx_alert_date (alert_date)` → Limpieza de antiguos

### Mantenimiento Automático

```sql
-- Crear evento para limpiar alertas antiguas (ejecutar 1 vez)
CREATE EVENT IF NOT EXISTS cleanup_old_alerts
ON SCHEDULE EVERY 1 DAY
DO
DELETE FROM alert_details
WHERE alert_date < DATE_SUB(NOW(), INTERVAL 30 DAY)
LIMIT 10000;
```

---

## 📚 Archivos del Sistema

| Archivo | Responsabilidad |
|---------|----------------|
| `types.go` | Estructuras de datos (`RuleExecution`, `IMEISummary`, `AlertDetail`) |
| `capture.go` | Captura en memoria thread-safe (`StartCapture`, `FinishCapture`) |
| `db.go` | Persistencia MySQL (`SaveExecutions`, `GetIMEISummaries`) |
| `../actions/actions.go` | Helper para reglas GRL (`ActionsHelper.RecordExecution`) |

---

## 🔒 Thread Safety

El sistema usa `sync.RWMutex` en múltiples niveles:

1. **`ExecutionCapture.mu`** (capture.go): Protege `[]RuleExecution` durante append concurrente
2. **`globalCaptures` + `capturesMutex`** (capture.go): Protege map de capturas activas
3. **Worker Pool** (grule_worker.go): 1 goroutine por IMEI, procesamiento secuencial

**Garantía**: Múltiples workers pueden capturar auditoría simultáneamente sin race conditions.

---

## 📊 Diferencias con Sistema Anterior

| Aspecto | Sistema Anterior | Sistema Universal |
|---------|------------------|-------------------|
| **Captura** | Post-mortem hardcoded | Tiempo real universal |
| **Código** | 1 función por regla | 1 función para todas |
| **Duplicados** | Miles de filas IMEI | 1 fila en summary |
| **Agregar reglas** | Modificar código Go | Solo crear .grl |
| **Condiciones** | Valores ficticios | Valores reales evaluados |
| **Performance** | N consultas UPDATE | 1 UPSERT por IMEI |

---

## 📞 Soporte

Para reportar bugs o sugerencias, revisar:
- `engine/audit/capture.go` - Lógica de captura
- `engine/audit/db.go` - Persistencia MySQL
- `actions/actions.go` - Interface para reglas GRL
