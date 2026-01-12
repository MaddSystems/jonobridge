# Visor de Auditoría con jqGrid

Sistema de visualización de auditoría usando jqGrid para mostrar alertas de reglas Grule.

## 📊 Arquitectura

```
┌─────────────────────────────────────────────────────────────────┐
│ Frontend (Flask + jqGrid)                                       │
│ templates/audit_summary.html                                    │
└─────────────────────────────────────────────────────────────────┘
                              ↓ HTTP GET
┌─────────────────────────────────────────────────────────────────┐
│ Proxy (Flask main.py)                                           │
│ /api/audit/grid → API_BASE_URL/api/audit/grid                  │
│ /api/audit/details → API_BASE_URL/api/audit/details            │
└─────────────────────────────────────────────────────────────────┘
                              ↓ HTTP GET
┌─────────────────────────────────────────────────────────────────┐
│ Backend (Go main.go)                                            │
│ auditGridHandler() → audit.GetIMEISummariesPaginated()         │
│ auditDetailsHandler() → audit.GetAlertDetails()                │
└─────────────────────────────────────────────────────────────────┘
                              ↓ SQL Query
┌─────────────────────────────────────────────────────────────────┐
│ MySQL Database                                                  │
│ - alert_summary (1 fila por IMEI)                              │
│ - alert_details (historial completo)                           │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🗂️ Tablas Utilizadas

### 1. `alert_summary` - Grid Principal
**Vista consolidada sin duplicados (1 fila por IMEI)**

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `imei` | VARCHAR(20) | IMEI del dispositivo (PK) |
| `last_alert_date` | DATETIME(6) | Última alerta disparada |
| `total_alerts_24h` | INT | Total alertas últimas 24h |
| `alert_types` | JSON | `{"SpeedAlert": 5, "Jammer": 2}` |
| `last_rule_executed` | VARCHAR(100) | Última regla ejecutada |
| `last_alert_location` | VARCHAR(100) | "lat,lon" |

**Endpoint:** `GET /api/audit/grid`

**Parámetros jqGrid:**
- `page`: Número de página (1, 2, 3...)
- `rows`: Registros por página (25, 50, 100)
- `sidx`: Columna para ordenar (`imei`, `last_alert_date`, `total_alerts_24h`)
- `sord`: Orden (`ASC` o `DESC`)
- `searchText`: Texto de búsqueda en IMEI o regla

**Respuesta:**
```json
{
  "page": 1,
  "total": 10,
  "records": 250,
  "rows": [
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

---

### 2. `alert_details` - Modal de Detalles
**Historial completo de alertas por IMEI**

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | BIGINT | ID único |
| `imei` | VARCHAR(20) | IMEI del dispositivo |
| `alert_date` | DATETIME(6) | Fecha/hora de la alerta |
| `rule_name` | VARCHAR(100) | Nombre de la regla |
| `rule_description` | VARCHAR(255) | Descripción |
| `salience` | INT | Prioridad |
| `conditions_snapshot` | JSON | Condiciones evaluadas |
| `actions_executed` | JSON | `["SendTelegram", "Log"]` |
| `telegram_sent` | BOOLEAN | Si se envió Telegram |
| `latitude` | DECIMAL(10,6) | Latitud GPS |
| `longitude` | DECIMAL(10,6) | Longitud GPS |
| `speed` | INT | Velocidad en km/h |

**Endpoint:** `GET /api/audit/details?imei=123456789012345`

**Respuesta:**
```json
{
  "success": true,
  "imei": "123456789012345",
  "count": 15,
  "data": [
    {
      "id": 12345,
      "imei": "123456789012345",
      "alert_date": "2025-12-10T15:30:00.123456Z",
      "rule_name": "SpeedAlert",
      "rule_description": "Velocidad > 120 km/h",
      "salience": 100,
      "conditions": {"Speed": 135, "Latitude": 40.416775},
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

## 🎨 Características del Visor

### Grid Principal (alert_summary)
- ✅ **Paginación**: 25/50/100 registros por página
- ✅ **Ordenamiento**: Por cualquier columna (IMEI, fecha, alertas)
- ✅ **Búsqueda**: Filtra por IMEI o nombre de regla
- ✅ **Sin Duplicados**: 1 fila por IMEI
- ✅ **Click en IMEI**: Abre modal con historial completo

### Modal de Detalles (alert_details)
- ✅ **Historial Completo**: Todas las alertas del IMEI
- ✅ **Timeline**: Ordenado por fecha descendente
- ✅ **Mapa GPS**: Link directo a Google Maps
- ✅ **Indicadores**: Badge de Telegram enviado/no enviado
- ✅ **Scroll Vertical**: Para muchas alertas

### Formatters Personalizados
```javascript
// IMEI con badge clickeable
function imeiFormatter(cellvalue, options, rowObject) {
    return '<span class="badge-imei">' + cellvalue + '</span>';
}

// Contador de alertas con color
function alertCountFormatter(cellvalue, options, rowObject) {
    if (cellvalue > 10) {
        return '<span class="badge-alerts">⚠️ ' + cellvalue + '</span>';
    }
    return '<span class="badge bg-warning">' + cellvalue + '</span>';
}

// Tipos de alertas como lista
function alertTypesFormatter(cellvalue, options, rowObject) {
    let html = '';
    for (const [rule, count] of Object.entries(cellvalue)) {
        html += '<span>' + rule + ': ' + count + '</span> ';
    }
    return html;
}

// Ubicación con link a Google Maps
function locationFormatter(cellvalue, options, rowObject) {
    const coords = cellvalue.split(',');
    const mapsUrl = 'https://www.google.com/maps?q=' + coords[0] + ',' + coords[1];
    return '<a href="' + mapsUrl + '" target="_blank">📍 Ver mapa</a>';
}
```

---

## 🚀 Uso

### 1. Acceder al Dashboard
```
http://localhost:5000/audit
```

### 2. Buscar por IMEI
```
[IMEI: 123456789]  [🔍 Buscar]  [Limpiar]  [📥 Exportar CSV]
```

### 3. Ver Detalles
**Click en cualquier IMEI** → Se abre modal con historial completo

### 4. Exportar Datos
**Click en "Exportar CSV"** → Descarga CSV con datos actuales

---

## 🔧 Configuración

### Variables de Entorno (Flask)
```bash
export API_BASE_URL="https://jonobridge.madd.com.mx/grule"
export FLASK_PORT=5000
```

### Variables de Entorno (Go)
```bash
export GRULE_AUDIT_ENABLED=Y
export GRULE_AUDIT_LEVEL=ALL
export PORTAL_ENDPOINT=/grule
```

---

## 📝 Endpoints Disponibles

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/audit` | GET | Página principal del visor |
| `/api/audit/grid` | GET | Grid principal (jqGrid) |
| `/api/audit/details` | GET | Detalles por IMEI |
| `/api/audit/summary` | GET | Resumen completo (no paginado) |

---

## 🧪 Testing

### 1. Verificar que hay datos
```sql
SELECT COUNT(*) FROM alert_summary;
SELECT COUNT(*) FROM alert_details;
```

### 2. Ver último IMEI con alertas
```sql
SELECT * FROM alert_summary ORDER BY last_alert_date DESC LIMIT 1;
```

### 3. Test endpoint grid
```bash
curl "http://localhost:8080/grule/api/audit/grid?page=1&rows=25&sidx=last_alert_date&sord=DESC"
```

### 4. Test endpoint details
```bash
curl "http://localhost:8080/grule/api/audit/details?imei=123456789012345&limit=50"
```

---

## 🎯 Ventajas vs Sistema Anterior

| Aspecto | Visor Anterior | Visor jqGrid |
|---------|----------------|--------------|
| **Duplicados** | Miles de filas por IMEI | 1 fila por IMEI |
| **Performance** | Carga todo en memoria | Paginación server-side |
| **Búsqueda** | Cliente (lento) | Servidor (MySQL índices) |
| **Ordenamiento** | Cliente | Servidor |
| **Exportar** | No disponible | CSV completo |
| **Responsivo** | Limitado | Bootstrap 5 + jqGrid |

---

## 📚 Dependencias

### Frontend
- **jQuery**: 3.6.0
- **jqGrid**: 4.15.4 (free-jqgrid)
- **Bootstrap**: 5.3.0
- **Font Awesome**: 6.0.0

### Backend
- **Go**: 1.21+
- **Flask**: 2.0+
- **MySQL**: 8.0+

---

## 🔒 Seguridad

### CORS Configurado
```go
"Access-Control-Allow-Origin": "*"
"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS"
"Access-Control-Allow-Headers": "Content-Type"
```

### SQL Injection Protegido
```go
// Whitelist de columnas permitidas
allowedSortColumns := map[string]bool{
    "imei": true,
    "last_alert_date": true,
    "total_alerts_24h": true,
}
```

### Prepared Statements
```go
query := "SELECT * FROM alert_summary WHERE imei LIKE ? LIMIT ? OFFSET ?"
db.Query(query, "%"+searchText+"%", limit, offset)
```

---

## 📞 Soporte

Para reportar bugs o sugerencias:
- Backend Go: `main.go` y `engine/audit/db.go`
- Frontend Flask: `main.py` y `templates/audit_summary.html`
- Documentación: Este README
