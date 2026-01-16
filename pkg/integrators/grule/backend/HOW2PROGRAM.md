# How to Program: Extending the GRULE Backend

This guide explains how to extend the backend's functionality. The architecture is designed as a set of modular "Lego bricks" (Capabilities) that implement the **Strategy Pattern**, allowing for interchangeable and testable logic.

## 🏗️ Architecture Overview

To add a new feature, you typically follow this data flow:

1.  **Adapters (`adapters/`)**: Where data enters (e.g., GPS packet).
2.  **Logic Flags (`grule/packet.go`)**: Where intermediate states are stored for the rule engine.
3.  **Capabilities (`capabilities/`)**: Where the complex logic resides (the "Strategy").
4.  **Persistence (`persistence/`)**: How data is saved.
5.  **Audit (`audit/`)**: How execution is recorded.

---

## � Extending the GPS Tracker Adapter

The `GPSTrackerAdapter` in `adapters/gps_tracker.go` parses incoming Jono Protocol JSON payloads. To add support for new fields from the protocol, follow these steps:

### Step 1: Update the IncomingPacket Struct

Add new fields to the `IncomingPacket` struct in `grule/packet.go` to store the additional data:

**File:** `backend/grule/packet.go`

```go
type IncomingPacket struct {
    // ... existing fields ...
    
    // New Jono Protocol fields
    Altitude           int64    // From packet.Altitude
    Direction          int64    // From packet.Direction  
    HDOP               int64    // From packet.HDOP
    NumberOfSatellites int64    // From packet.NumberOfSatellites
    Mileage            int64    // From packet.Mileage
    RunTime            int64    // From packet.RunTime
    
    // Analog inputs (AD1-AD10)
    AnalogInputs       map[string]string // From packet.AnalogInputs
    
    // Port statuses
    InputPortStatus    map[string]interface{} // From packet.InputPortStatus
    OutputPortStatus   map[string]interface{} // From packet.OutputPortStatus
    IoPortStatus       map[string]int         // From packet.IoPortStatus
    
    // Base station info
    CellID             string   // From packet.BaseStationInfo.CellID
    LAC                string   // From packet.BaseStationInfo.LAC
    MCC                string   // From packet.BaseStationInfo.MCC
    MNC                string   // From packet.BaseStationInfo.MNC
    
    // System flags
    SystemFlag         map[string]interface{} // From packet.SystemFlag
    
    // Event information
    EventCode          map[string]interface{} // From packet.EventCode
    
    // Additional sensor data
    Temperature        float64 // From packet.TemperatureSensor.Value (if available)
    Humidity           float64 // From packet.TemperatureAndHumiditySensor.Humidity (if available)
}
```

### Step 2: Update the GPS Tracker Adapter

Modify `adapters/gps_tracker.go` to parse and populate the new fields:

**File:** `backend/adapters/gps_tracker.go`

```go
packet := &grule.IncomingPacket{
    IMEI:              jono.IMEI,
    Speed:             speedKmH,
    GSMSignalStrength: gsm,
    Datetime:          p.Datetime,
    PositioningStatus: p.PositioningStatus,
    Latitude:          p.Latitude,
    Longitude:         p.Longitude,
    
    // Parse new fields
    Altitude:           int64(p.Altitude),
    Direction:          int64(p.Direction),
    HDOP:               int64(p.HDOP),
    NumberOfSatellites: int64(p.NumberOfSatellites),
    Mileage:            p.Mileage,
    RunTime:            p.RunTime,
    
    // Parse complex structures
    AnalogInputs:       p.AnalogInputs,
    InputPortStatus:    p.InputPortStatus,
    OutputPortStatus:   p.OutputPortStatus,
    IoPortStatus:       p.IoPortStatus,
    SystemFlag:         p.SystemFlag,
    EventCode:          p.EventCode,
    
    // Parse base station info
    CellID:             p.BaseStationInfo.CellID,
    LAC:                p.BaseStationInfo.LAC,
    MCC:                p.BaseStationInfo.MCC,
    MNC:                p.BaseStationInfo.MNC,
    
    // Parse sensor data (with safe checks)
    Temperature:        parseTemperature(p),
    Humidity:           parseHumidity(p),
    
    // Initialize flags to false
    BufferUpdated:           false,
    BufferHas10:             false,
    IsOfflineFor5Min:        false,
    PositionInvalidDetected: false,
    MetricsReady:            false,
    MovingWithWeakSignal:    false,
    OutsideAllSafeZones:     false,
}
```

### Step 3: Handle Optional Fields

Some fields may be null or optional. Add safe parsing:

```go
// Safe parsing for optional fields
altitude := int64(0)
if p.Altitude != nil {
    altitude = int64(*p.Altitude)
}

direction := int64(0)
if p.Direction != nil {
    direction = int64(*p.Direction)
}

// For map fields, check if they exist
analogInputs := make(map[string]string)
if p.AnalogInputs != nil {
    analogInputs = p.AnalogInputs
}

// Helper functions for sensor data
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

### Step 4: Update Audit Snapshots

To include the new fields in audit snapshots, update the `ExtractSnapshot` function in `audit/snapshot.go`. You may need to add a helper function for map fields:

**Add this helper function to `backend/audit/snapshot.go`:**

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

**Then update the extracted map:**

```go
extracted := map[string]interface{}{
    // ... existing fields ...
    
    // New fields for audit
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

## �🚀 Step-by-Step: Adding a New Feature

Let's assume you want to add a feature to track **Fuel Level**.

### Step 1: Update Logic Flags (`grule/packet.go`)

The `IncomingPacket` struct holds the state used by the Rule Engine. Add new flags here to store data or decision results.

**File:** `backend/grule/packet.go`

```go
type IncomingPacket struct {
    // ... existing fields ...
    
    // New Logic Flag
    FuelLevelCritical bool // true if fuel < 10%
}
```

### Step 2: Create a New Capability (`capabilities/`)

Capabilities encapsulate logic. Each capability is a "Strategy" that implements a common interface.

1.  **Create Folder:** `backend/capabilities/fuel/`
2.  **Define Strategy:** Implement the logic.

**File:** `backend/capabilities/fuel/capability.go`

```go
package fuel

// FuelCapability implements the logic for fuel monitoring
type FuelCapability struct {
    threshold float64
}

func NewFuelCapability() *FuelCapability {
    return &FuelCapability{threshold: 10.0}
}

// CheckLevel is the strategy method called by rules
func (f *FuelCapability) CheckLevel(currentLevel float64) bool {
    return currentLevel < f.threshold
}

// GetSnapshotData implements SnapshotProvider for Audit (See Step 4)
func (f *FuelCapability) GetSnapshotData(imei string) map[string]interface{} {
    return map[string]interface{}{
        "fuel_threshold": f.threshold,
    }
}
```

3.  **Register Capability:** Add it to the `StateWrapper` so rules can see it.

**File:** `backend/grule/context_builder.go`

```go
type StateWrapper struct {
    // ...
    Fuel *fuel.FuelCapability // Register new capability
}

// In Build() function:
func (cb *ContextBuilder) Build(...) {
    // ...
    state := &StateWrapper{
        Fuel: fuel.NewFuelCapability(),
    }
}
```

### Step 3: Persistence (`persistence/`) (If needed)

If your capability needs to save state to a DB:

1.  Define the interface in `backend/persistence/interface.go`.
2.  Implement it in `backend/persistence/mysql.go`.

### Step 4: Update Audit System (`audit/`)

To ensure your new data appears in the audit snapshots (the "Movie" UI), your capability must implement the `SnapshotProvider` interface.

**Interface:** `backend/capabilities/interface.go`

```go
type SnapshotProvider interface {
    GetSnapshotData(imei string) map[string]interface{}
}
```

**Implementation:**
You already did this in **Step 2**! The `GetSnapshotData` method returns a map of data. The Audit system automatically discovers all capabilities in `StateWrapper` that implement this interface and merges their data into the snapshot.

**Key Concept:** The Audit system is **Declarative**. You don't change the audit code; you just expose data from your capability.

### Step 5: Create/Update Rules (`.grl`)

Now you can use your new flag and capability in the rules.

**File:** `frontend/rules_templates/my_rules.grl`

```grl
rule CheckFuelLevel "Check if fuel is critical" salience 100 {
    when
        !IncomingPacket.FuelLevelCritical &&
        state.Fuel.CheckLevel(IncomingPacket.FuelLevel) // Call your capability
    then
        IncomingPacket.FuelLevelCritical = true;
        actions.Log("Fuel is critical!");
        actions.CaptureSnapshot("CheckFuelLevel"); // Explicitly capture snapshot
}
```

---

## 🧠 Core Concepts

### Strategy Pattern
We use the Strategy Pattern to make capabilities interchangeable. For example, `capabilities/buffer/` could have a `CircularBuffer` strategy or a `RedisBuffer` strategy. The rules just call `state.Buffer.Add()`, not caring about the implementation.

### Explicit Audit Capture
We do **not** use automatic background capturing. You must explicitly call `actions.CaptureSnapshot("RuleName")` in your GRL rule when significant state changes occur. This prevents duplicate logs and gives you full control over *when* a snapshot is worth saving.
