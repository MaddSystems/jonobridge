# Cómo Programar: Extendiendo el Backend de GRULE

Esta guía explica cómo extender la funcionalidad del backend. La arquitectura está diseñada como un conjunto de "bloques de Lego" modulares (Capabilities) que implementan el **Patrón Strategy**, permitiendo lógica intercambiable y testeable.

## 🏗️ Resumen de Arquitectura

Para agregar una nueva funcionalidad, típicamente sigues este flujo de datos:

1.  **Adapters (`adapters/`)**: Por donde entran los datos (ej. paquete GPS).
2.  **Logic Flags (`grule/packet.go`)**: Donde se guardan estados intermedios para el motor de reglas.
3.  **Capabilities (`capabilities/`)**: Donde reside la lógica compleja (la "Estrategia").
4.  **Persistence (`persistence/`)**: Cómo se guardan los datos.
5.  **Audit (`audit/`)**: Cómo se registra la ejecución.

---

## � Extendiendo el Adaptador GPS Tracker

El `GPSTrackerAdapter` en `adapters/gps_tracker.go` analiza las cargas útiles JSON entrantes del Protocolo Jono. Para agregar soporte para nuevos campos del protocolo, sigue estos pasos:

### Paso 1: Actualizar la Estructura IncomingPacket

Agrega nuevos campos a la estructura `IncomingPacket` en `grule/packet.go` para almacenar los datos adicionales:

**Archivo:** `backend/grule/packet.go`

```go
type IncomingPacket struct {
    // ... campos existentes ...
    
    // Nuevos campos del Protocolo Jono
    Altitude           int64    // De packet.Altitude
    Direction          int64    // De packet.Direction  
    HDOP               int64    // De packet.HDOP
    NumberOfSatellites int64    // De packet.NumberOfSatellites
    Mileage            int64    // De packet.Mileage
    RunTime            int64    // De packet.RunTime
    
    // Entradas analógicas (AD1-AD10)
    AnalogInputs       map[string]string // De packet.AnalogInputs
    
    // Estados de puertos
    InputPortStatus    map[string]interface{} // De packet.InputPortStatus
    OutputPortStatus   map[string]interface{} // De packet.OutputPortStatus
    IoPortStatus       map[string]int         // De packet.IoPortStatus
    
    // Información de estación base
    CellID             string   // De packet.BaseStationInfo.CellID
    LAC                string   // De packet.BaseStationInfo.LAC
    MCC                string   // De packet.BaseStationInfo.MCC
    MNC                string   // De packet.BaseStationInfo.MNC
    
    // Banderas del sistema
    SystemFlag         map[string]interface{} // De packet.SystemFlag
    
    // Información de evento
    EventCode          map[string]interface{} // De packet.EventCode
    
    // Datos adicionales de sensores
    Temperature        float64 // De packet.TemperatureSensor.Value (si está disponible)
    Humidity           float64 // De packet.TemperatureAndHumiditySensor.Humidity (si está disponible)
}
```

### Paso 2: Actualizar el Adaptador GPS Tracker

Modifica `adapters/gps_tracker.go` para analizar y poblar los nuevos campos:

**Archivo:** `backend/adapters/gps_tracker.go`

```go
packet := &grule.IncomingPacket{
    IMEI:              jono.IMEI,
    Speed:             speedKmH,
    GSMSignalStrength: gsm,
    Datetime:          p.Datetime,
    PositioningStatus: p.PositioningStatus,
    Latitude:          p.Latitude,
    Longitude:         p.Longitude,
    
    // Analizar nuevos campos
    Altitude:           int64(p.Altitude),
    Direction:          int64(p.Direction),
    HDOP:               int64(p.HDOP),
    NumberOfSatellites: int64(p.NumberOfSatellites),
    Mileage:            p.Mileage,
    RunTime:            p.RunTime,
    
    // Analizar estructuras complejas
    AnalogInputs:       p.AnalogInputs,
    InputPortStatus:    p.InputPortStatus,
    OutputPortStatus:   p.OutputPortStatus,
    IoPortStatus:       p.IoPortStatus,
    SystemFlag:         p.SystemFlag,
    EventCode:          p.EventCode,
    
    // Analizar información de estación base
    CellID:             p.BaseStationInfo.CellID,
    LAC:                p.BaseStationInfo.LAC,
    MCC:                p.BaseStationInfo.MCC,
    MNC:                p.BaseStationInfo.MNC,
    
    // Analizar datos de sensores (con verificaciones seguras)
    Temperature:        parseTemperature(p),
    Humidity:           parseHumidity(p),
    
    // Inicializar banderas en false
    BufferUpdated:           false,
    BufferHas10:             false,
    IsOfflineFor5Min:        false,
    PositionInvalidDetected: false,
    MetricsReady:            false,
    MovingWithWeakSignal:    false,
    OutsideAllSafeZones:     false,
}
```

### Paso 3: Manejar Campos Opcionales

Algunos campos pueden ser nulos u opcionales. Agrega análisis seguro:

```go
// Análisis seguro para campos opcionales
altitude := int64(0)
if p.Altitude != nil {
    altitude = int64(*p.Altitude)
}

direction := int64(0)
if p.Direction != nil {
    direction = int64(*p.Direction)
}

// Para campos de mapa, verifica si existen
analogInputs := make(map[string]string)
if p.AnalogInputs != nil {
    analogInputs = p.AnalogInputs
}

// Funciones auxiliares para datos de sensores
func parseTemperature(p *models.JonoPacket) float64 {
    if p.TemperatureSensor != nil && p.TemperatureSensor.Value != nil {
        return *p.TemperatureSensor.Value
    }
    return 0.0
}

func parseHumidity(p *models.JonoPacket) float64 {
    if p.TemperatureAndHumiditySensor != nil && p.TemperatureAndHumiditySensor.Humidity != nil {
        return *p.TemperatureAndHumiditySensor.Humidity
    }
    return 0.0
}
```

### Paso 4: Actualizar Snapshots de Auditoría

Para incluir los nuevos campos en los snapshots de auditoría, actualiza la función `ExtractSnapshot` en `audit/snapshot.go`. Puede que necesites agregar una función auxiliar para campos de mapa:

**Agrega esta función auxiliar a `backend/audit/snapshot.go`:**

```go
func getFieldMap(v reflect.Value, name string) map[string]interface{} {
	f := v.FieldByName(name)
	if f.IsValid() && f.Kind() == reflect.Map {
		if f.Type().Key().Kind() == reflect.String {
			result := make(map[string]interface{})
			for _, key := range f.MapKeys() {
				if value := f.MapIndex(key); value.IsValid() && value.CanInterface() {
					result[key.String()] = value.Interface()
				}
			}
			return result
		}
	}
	return make(map[string]interface{})
}
```

**Luego actualiza el mapa extraído:**

```go
extracted := map[string]interface{}{
    // ... campos existentes ...
    
    // Nuevos campos para auditoría
    "Altitude":           getFieldInt(v, "Altitude"),
    "Direction":          getFieldInt(v, "Direction"),
    "HDOP":               getFieldInt(v, "HDOP"),
    "NumberOfSatellites": getFieldInt(v, "NumberOfSatellites"),
    "Mileage":            getFieldInt(v, "Mileage"),
    "RunTime":            getFieldInt(v, "RunTime"),
    "AnalogInputs":       getFieldMap(v, "AnalogInputs"),
    "CellID":             getFieldString(v, "CellID"),
    "LAC":                getFieldString(v, "LAC"),
    "MCC":                getFieldString(v, "MCC"),
    "MNC":                getFieldString(v, "MNC"),
    "Temperature":        getFieldFloat(v, "Temperature"),
    "Humidity":           getFieldFloat(v, "Humidity"),
}
```

---

## �🚀 Paso a Paso: Agregando una Nueva Funcionalidad

Supongamos que quieres agregar una funcionalidad para rastrear el **Nivel de Combustible**.

### Paso 1: Actualizar Logic Flags (`grule/packet.go`)

La estructura `IncomingPacket` mantiene el estado usado por el Motor de Reglas. Agrega nuevas banderas (flags) aquí para guardar datos o resultados de decisiones.

**Archivo:** `backend/grule/packet.go`

```go
type IncomingPacket struct {
    // ... campos existentes ...
    
    // Nueva Logic Flag
    FuelLevelCritical bool // true si combustible < 10%
}
```

### Paso 2: Crear una Nueva Capability (`capabilities/`)

Las Capabilities encapsulan lógica. Cada capability es una "Estrategia" que implementa una interfaz común.

1.  **Crear Carpeta:** `backend/capabilities/fuel/`
2.  **Definir Estrategia:** Implementar la lógica.

**Archivo:** `backend/capabilities/fuel/capability.go`

```go
package fuel

// FuelCapability implementa la lógica para monitoreo de combustible
type FuelCapability struct {
    threshold float64
}

func NewFuelCapability() *FuelCapability {
    return &FuelCapability{threshold: 10.0}
}

// CheckLevel es el método de estrategia llamado por las reglas
func (f *FuelCapability) CheckLevel(currentLevel float64) bool {
    return currentLevel < f.threshold
}

// GetSnapshotData implementa SnapshotProvider para Auditoría (Ver Paso 4)
func (f *FuelCapability) GetSnapshotData(imei string) map[string]interface{} {
    return map[string]interface{}{
        "fuel_threshold": f.threshold,
    }
}
```

3.  **Registrar Capability:** Agregarla al `StateWrapper` para que las reglas puedan verla.

**Archivo:** `backend/grule/context_builder.go`

```go
type StateWrapper struct {
    // ...
    Fuel *fuel.FuelCapability // Registrar nueva capability
}

// En la función Build():
func (cb *ContextBuilder) Build(...) {
    // ...
    state := &StateWrapper{
        Fuel: fuel.NewFuelCapability(),
    }
}
```

### Paso 3: Persistencia (`persistence/`) (Si es necesario)

Si tu capability necesita guardar estado en una DB:

1.  Define la interfaz en `backend/persistence/interface.go`.
2.  Impleméntala en `backend/persistence/mysql.go`.

### Paso 4: Actualizar Sistema de Auditoría (`audit/`)

Para asegurar que tus nuevos datos aparezcan en los snapshots de auditoría (la UI de "Película"), tu capability debe implementar la interfaz `SnapshotProvider`.

**Interfaz:** `backend/capabilities/interface.go`

```go
type SnapshotProvider interface {
    GetSnapshotData(imei string) map[string]interface{}
}
```

**Implementación:**
¡Ya hiciste esto en el **Paso 2**! El método `GetSnapshotData` retorna un mapa de datos. El sistema de Auditoría descubre automáticamente todas las capabilities en `StateWrapper` que implementan esta interfaz y fusiona sus datos en el snapshot.

**Concepto Clave:** El sistema de Auditoría es **Declarativo**. No cambias el código de auditoría; solo expones datos desde tu capability.

### Paso 5: Crear/Actualizar Reglas (`.grl`)

Ahora puedes usar tu nueva bandera y capability en las reglas.

**Archivo:** `frontend/rules_templates/my_rules.grl`

```grl
rule CheckFuelLevel "Verificar si el combustible es crítico" salience 100 {
    when
        !IncomingPacket.FuelLevelCritical &&
        state.Fuel.CheckLevel(IncomingPacket.FuelLevel) // Llamar a tu capability
    then
        IncomingPacket.FuelLevelCritical = true;
        actions.Log("¡Combustible crítico!");
        actions.CaptureSnapshot("CheckFuelLevel"); // Capturar snapshot explícitamente
}
```

---

## 🧠 Conceptos Centrales

### Patrón Strategy
Usamos el Patrón Strategy para hacer que las capabilities sean intercambiables. Por ejemplo, `capabilities/buffer/` podría tener una estrategia `CircularBuffer` o una estrategia `RedisBuffer`. Las reglas solo llaman a `state.Buffer.Add()`, sin importarles la implementación.

### Captura de Auditoría Explícita
**No** usamos captura automática en segundo plano. Debes llamar explícitamente a `actions.CaptureSnapshot("NombreRegla")` en tu regla GRL cuando ocurran cambios de estado significativos. Esto previene registros duplicados y te da control total sobre *cuándo* vale la pena guardar un snapshot.
