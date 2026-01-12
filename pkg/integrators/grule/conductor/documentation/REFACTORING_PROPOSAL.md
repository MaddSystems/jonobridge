# Refactoring Proposal: Modular "Lego Blocks" Architecture for GRULE Rule Engine

**Date:** January 6, 2026  
**Author:** Architecture Team  
**Version:** 1.0  

---

## Executive Summary

This document proposes a comprehensive refactoring of the GRULE rule engine codebase to achieve:

1. **Modular "Lego Block" Architecture** - Isolated, testable functional units
2. **LLM-Friendly JSON Schema** - Machine-readable capability manifest for AI-powered rule generation
3. **Domain Agnostic Design** - Extensibility beyond GPS trackers to any IoT device
4. **Pattern Compatibility** - GRULE-compatible design patterns for maintainability

---

## Table of Contents

1. [Current State Analysis](#1-current-state-analysis)
2. [Four Architecture Scenarios](#2-four-architecture-scenarios)
3. [Recommended Solution](#3-recommended-solution)
4. [JSON Schema for LLM Integration](#4-json-schema-for-llm-integration)
5. [Implementation Roadmap](#5-implementation-roadmap)
6. [Pattern Analysis: Factory vs Rule Engine](#6-pattern-analysis-factory-vs-rule-engine)

---

## 1. Current State Analysis

### 1.1 Current File Structure
```
engine/
├── grule_worker.go       (~700 lines) - Mixed responsibilities
├── persistent_state.go   (~650 lines) - State + Geofences + Counters + Jammer
├── memory_buffer.go      (~280 lines) - Buffer logic (relatively clean)
├── property.go           (~470 lines) - Execution flags (well structured)
├── rule_loader.go        (~680 lines) - DB + API + Rules
└── audit/
    ├── capture.go
    ├── db.go
    └── types.go
```

### 1.2 Identified Problems

| Problem | Impact | Files Affected |
|---------|--------|----------------|
| **Mixed Responsibilities** | Hard to test/debug | `persistent_state.go`, `grule_worker.go` |
| **Large File Sizes** | LLM context overflow | All files >400 lines |
| **Coupled Domain Logic** | Not reusable for IoT | Geofences mixed with GPS state |
| **No Capability Discovery** | LLM can't know what's available | No metadata/schema |
| **Hardcoded Jammer Logic** | Not composable | `persistent_state.go` |

### 1.3 Current GRL Rule Analysis (jammer_wargames.grl)

```grl
// Current rule uses these capabilities:
IncomingPacket.BufferUpdated          // → MemoryBuffer brick
state.UpdateMemoryBuffer(...)          // → MemoryBuffer brick
state.IsOfflineFor(5)                  // → Timing brick
state.IsInsideGroup("Taller", ...)    // → Geofence brick
state.JammerAvgSpeed90min             // → Metrics brick
state.MarkAlertSent(...)              // → AlertState brick
actions.SendTelegram(...)             // → Actions brick
```

**Observation:** The rule implicitly uses 6 different "bricks" but they're scattered across files.

---

## 2. Four Architecture Scenarios

### Scenario A: Domain-Centric Modules (Vertical Slicing)

```
bricks/
├── geofence/
│   ├── geofence.go           # Core logic
│   ├── geofence_grule.go     # GRULE adapter
│   ├── geofence_test.go
│   └── schema.json           # LLM capability manifest
├── buffer/
│   ├── circular_buffer.go
│   ├── buffer_grule.go
│   └── schema.json
├── timing/
│   ├── offline_detector.go
│   ├── timing_grule.go
│   └── schema.json
├── metrics/
│   ├── calculator.go
│   ├── metrics_grule.go
│   └── schema.json
└── alerts/
    ├── spam_guard.go
    ├── alerts_grule.go
    └── schema.json
```

**Pros:**
- ✅ Each "brick" is self-contained
- ✅ Easy to test in isolation
- ✅ LLM can load only relevant schemas
- ✅ Clear ownership and responsibility

**Cons:**
- ❌ Many small packages (import complexity)
- ❌ Shared state management becomes complex
- ❌ GRULE DataContext injection needs coordination

---

### Scenario B: Layer-Based Architecture (Horizontal Slicing)

```
engine/
├── core/
│   ├── packet.go            # Input data structures
│   ├── context.go           # Execution context
│   └── registry.go          # Brick registry
├── bricks/
│   ├── interface.go         # Brick interface definition
│   ├── geofence.go
│   ├── buffer.go
│   ├── timing.go
│   ├── metrics.go
│   └── alerts.go
├── adapters/
│   ├── grule_adapter.go     # Single GRULE integration point
│   └── rest_adapter.go      # API for rule management
├── persistence/
│   ├── mysql.go
│   └── memory.go
└── schema/
    └── capabilities.json     # Single unified schema
```

**Pros:**
- ✅ Clean separation of concerns
- ✅ Single point of GRULE integration
- ✅ Easy to swap persistence layer
- ✅ Unified schema for LLM

**Cons:**
- ❌ Bricks share same package (potential coupling)
- ❌ Less intuitive for adding new domains
- ❌ Schema becomes very large

---

### Scenario C: Plugin Architecture with Factory Pattern

```
engine/
├── core/
│   ├── brick_interface.go
│   ├── brick_factory.go     # Factory creates brick instances
│   ├── brick_registry.go    # Registry of available bricks
│   └── context.go
├── plugins/
│   ├── gps/
│   │   ├── geofence_brick.go
│   │   ├── speed_brick.go
│   │   └── manifest.json
│   ├── iot/
│   │   ├── temperature_brick.go
│   │   ├── humidity_brick.go
│   │   └── manifest.json
│   └── common/
│       ├── buffer_brick.go
│       ├── timing_brick.go
│       └── manifest.json
├── grule/
│   ├── data_context_builder.go
│   ├── rule_executor.go
│   └── schema_generator.go   # Auto-generates LLM schema
└── persistence/
    └── state_manager.go
```

**Pros:**
- ✅ Extensible for any domain (GPS, IoT, etc.)
- ✅ Factory pattern creates configured instances
- ✅ Auto-discovery of capabilities
- ✅ Perfect for LLM: load domain-specific manifest

**Cons:**
- ❌ More complex initial setup
- ❌ Factory boilerplate overhead
- ❌ Plugin discovery at runtime

---

### Scenario D: Capability-Oriented Architecture (Recommended)

```
engine/
├── capabilities/                    # "LEGO BRICKS"
│   ├── interface.go                 # Unified Capability interface
│   ├── registry.go                  # Auto-registration
│   │
│   ├── geofence/                    # BRICK: Geofence
│   │   ├── capability.go            # Implements Capability interface
│   │   ├── functions.go             # Pure logic (testable)
│   │   ├── state.go                 # State management
│   │   └── manifest.yaml            # LLM schema (human-readable)
│   │
│   ├── buffer/                      # BRICK: Circular Buffer
│   │   ├── capability.go
│   │   ├── circular.go
│   │   └── manifest.yaml
│   │
│   ├── timing/                      # BRICK: Time-based Detection
│   │   ├── capability.go
│   │   ├── offline.go
│   │   ├── duration.go
│   │   └── manifest.yaml
│   │
│   ├── metrics/                     # BRICK: Metric Calculations
│   │   ├── capability.go
│   │   ├── averages.go
│   │   └── manifest.yaml
│   │
│   ├── alerts/                      # BRICK: Alert Management
│   │   ├── capability.go
│   │   ├── spam_guard.go
│   │   ├── channels.go              # Telegram, Email, Webhook
│   │   └── manifest.yaml
│   │
│   └── state/                       # BRICK: Persistent State
│       ├── capability.go
│       ├── counters.go
│       ├── flags.go
│       └── manifest.yaml
│
├── grule/                           # GRULE Integration Layer
│   ├── context_builder.go           # Builds DataContext from capabilities
│   ├── executor.go                  # Rule execution
│   ├── worker.go                    # Packet processing (simplified)
│   └── loader.go                    # Rule loading
│
├── adapters/                        # External Adapters
│   ├── gps_tracker.go               # GPS tracker input adapter
│   ├── iot_sensor.go                # IoT sensor input adapter (future)
│   └── mqtt.go                      # MQTT adapter (future)
│
├── persistence/                     # Persistence Layer
│   ├── mysql.go
│   ├── redis.go                     # Future: Redis for hot state
│   └── interface.go
│
├── audit/                           # Audit (unchanged)
│   ├── capture.go
│   ├── db.go
│   └── types.go
│
└── schema/                          # LLM Integration
    ├── generator.go                 # Generates unified schema from manifests
    ├── capabilities.json            # Auto-generated unified schema
    └── examples/                    # Example rules for LLM training
        ├── jammer_detection.grl
        ├── speed_alert.grl
        └── geofence_timing.grl
```

**Pros:**
- ✅ **Each capability is a complete "Lego brick"**
- ✅ **YAML manifests are LLM-readable and human-editable**
- ✅ **Auto-generates unified JSON schema for LLM context**
- ✅ **Clear separation: capabilities vs GRULE integration vs persistence**
- ✅ **Extensible: add new domains by adding adapters**
- ✅ **Registry pattern for discovery**
- ✅ **Testable: each capability tests independently**

**Cons:**
- ❌ Requires upfront investment in capability interface design
- ❌ Migration requires careful planning

---

## 3. Recommended Solution

### **Scenario D: Capability-Oriented Architecture**

This is the recommended approach because:

1. **LLM Optimization**: YAML manifests per capability → easy to load partial context
2. **Composability**: True "Lego blocks" that compose without coupling
3. **Testability**: Each capability is independently testable
4. **Extensibility**: Adding IoT sensors only requires new adapters + capabilities
5. **GRULE Compatibility**: Single integration layer manages DataContext

### 3.1 Capability Interface Design

```go
// engine/capabilities/interface.go

package capabilities

import (
    "github.com/hyperjumptech/grule-rule-engine/ast"
)

// Capability represents a functional "Lego brick" that can be used in rules
type Capability interface {
    // Identity
    Name() string                      // e.g., "geofence", "buffer"
    Version() string                   // Semantic versioning
    Description() string               // Human-readable description
    
    // GRULE Integration
    GetDataContextName() string        // e.g., "geo", "buffer"
    GetGRULEFunctions() []GRULEFunction // Functions exposed to rules
    
    // State Management
    Initialize(imei string) error
    Reset() error
    GetSnapshot() map[string]interface{}
    
    // Schema
    GetManifest() Manifest
}

// GRULEFunction describes a function available in GRL rules
type GRULEFunction struct {
    Name        string           // e.g., "IsInsideGroup"
    Description string
    Parameters  []Parameter
    ReturnType  string           // "bool", "int64", "float64", "string"
    Example     string           // GRL usage example
}

// Parameter describes a function parameter
type Parameter struct {
    Name        string
    Type        string
    Description string
    Required    bool
}

// Manifest is the LLM-readable capability description
type Manifest struct {
    Name        string          `yaml:"name"`
    Version     string          `yaml:"version"`
    Description string          `yaml:"description"`
    Category    string          `yaml:"category"` // "geofence", "timing", "metrics"
    Functions   []GRULEFunction `yaml:"functions"`
    StateVars   []StateVariable `yaml:"state_variables"`
    Examples    []RuleExample   `yaml:"examples"`
}

// StateVariable describes readable/writable state
type StateVariable struct {
    Name        string `yaml:"name"`
    Type        string `yaml:"type"`
    Description string `yaml:"description"`
    Readable    bool   `yaml:"readable"`
    Writable    bool   `yaml:"writable"`
}

// RuleExample shows how to use this capability in GRL
type RuleExample struct {
    Name        string `yaml:"name"`
    Description string `yaml:"description"`
    GRL         string `yaml:"grl"`
}
```

### 3.2 Example: Geofence Capability

```go
// engine/capabilities/geofence/capability.go

package geofence

import (
    "github.com/jonobridge/grule-engine/engine/capabilities"
)

type GeofenceCapability struct {
    imei       string
    db         *sql.DB
    prevInside map[string]bool
    mutex      sync.RWMutex
}

func New(db *sql.DB) *GeofenceCapability {
    return &GeofenceCapability{
        db:         db,
        prevInside: make(map[string]bool),
    }
}

func (g *GeofenceCapability) Name() string { return "geofence" }
func (g *GeofenceCapability) Version() string { return "1.0.0" }
func (g *GeofenceCapability) Description() string {
    return "Geofence detection and timing capabilities"
}

func (g *GeofenceCapability) GetDataContextName() string {
    return "geo"
}

func (g *GeofenceCapability) GetGRULEFunctions() []capabilities.GRULEFunction {
    return []capabilities.GRULEFunction{
        {
            Name:        "IsInsideGroup",
            Description: "Check if vehicle is inside any geofence of a group",
            Parameters: []capabilities.Parameter{
                {Name: "groupName", Type: "string", Description: "Name of geofence group", Required: true},
                {Name: "lat", Type: "float64", Description: "Current latitude", Required: true},
                {Name: "lon", Type: "float64", Description: "Current longitude", Required: true},
            },
            ReturnType: "bool",
            Example:    `geo.IsInsideGroup("Taller", packet.Latitude, packet.Longitude)`,
        },
        {
            Name:        "MinutesInsideGeofence",
            Description: "Time spent inside a specific geofence in minutes",
            Parameters: []capabilities.Parameter{
                {Name: "geofenceID", Type: "string", Description: "Geofence identifier", Required: true},
                {Name: "currentlyInside", Type: "bool", Description: "Is currently inside", Required: true},
            },
            ReturnType: "float64",
            Example:    `geo.MinutesInsideGeofence("loading_zone", true) > 30`,
        },
    }
}

// Actual implementation methods
func (g *GeofenceCapability) IsInsideGroup(groupName string, lat, lon float64) bool {
    // ... implementation from persistent_state.go
}

func (g *GeofenceCapability) MinutesInsideGeofence(id string, currentlyInside bool) float64 {
    // ... implementation from persistent_state.go
}
```

### 3.3 Geofence Manifest (YAML)

```yaml
# engine/capabilities/geofence/manifest.yaml

name: geofence
version: "1.0.0"
description: |
  Geofence detection and timing capabilities.
  Allows rules to check if a device is inside specific areas
  and track how long it has been inside.

category: location

grule_context_name: geo

functions:
  - name: IsInsideGroup
    description: Check if the device is inside ANY geofence belonging to a named group
    parameters:
      - name: groupName
        type: string
        description: The name of the geofence group (e.g., "Taller", "CLIENTES")
        required: true
      - name: lat
        type: float64
        description: Current latitude of the device
        required: true
      - name: lon
        type: float64
        description: Current longitude of the device
        required: true
    return_type: bool
    example: |
      geo.IsInsideGroup("Taller", packet.Latitude, packet.Longitude)

  - name: IsInsideCircle
    description: Check if device is inside a circular geofence
    parameters:
      - name: lat
        type: float64
        description: Device latitude
        required: true
      - name: lon
        type: float64
        description: Device longitude
        required: true
      - name: centerLat
        type: float64
        description: Circle center latitude
        required: true
      - name: centerLon
        type: float64
        description: Circle center longitude
        required: true
      - name: radiusMeters
        type: float64
        description: Circle radius in meters
        required: true
    return_type: bool
    example: |
      geo.IsInsideCircle(packet.Latitude, packet.Longitude, -34.603722, -58.381592, 500)

  - name: MinutesInsideGeofence
    description: Get how many minutes the device has been inside a geofence
    parameters:
      - name: geofenceID
        type: string
        description: Unique identifier of the geofence
        required: true
      - name: currentlyInside
        type: bool
        description: Whether the device is currently inside (triggers timer start/stop)
        required: true
    return_type: float64
    example: |
      geo.MinutesInsideGeofence("loading_zone_A", true) > 30

  - name: EnteredGeofence
    description: Returns true only on the first packet after entering a geofence
    parameters:
      - name: geofenceID
        type: string
        description: Geofence identifier
        required: true
      - name: inside
        type: bool
        description: Current inside status
        required: true
    return_type: bool
    example: |
      geo.EnteredGeofence("warehouse", geo.IsInsideGroup("Warehouse", packet.Latitude, packet.Longitude))

state_variables: []  # This capability doesn't expose direct state, only functions

examples:
  - name: Alert on long parking
    description: Send alert if vehicle is parked in loading zone for more than 30 minutes
    grl: |
      rule LongParking "Alert on excessive parking time" salience 100 {
          when
              geo.MinutesInsideGeofence("loading_zone", 
                  geo.IsInsideGroup("LoadingZones", packet.Latitude, packet.Longitude)) > 30 &&
              !alerts.IsAlertSent("long_parking")
          then
              alerts.SendTelegram("Vehicle " + packet.IMEI + " parked >30min in loading zone");
              alerts.MarkAlertSent("long_parking");
      }

  - name: Geofence entry notification
    description: Notify when vehicle enters a specific area
    grl: |
      rule EnterWarehouse "Notify warehouse entry" salience 100 {
          when
              geo.EnteredGeofence("warehouse", 
                  geo.IsInsideGroup("Warehouse", packet.Latitude, packet.Longitude))
          then
              alerts.Log("Vehicle entered warehouse");
      }
```

---

## 4. JSON Schema for LLM Integration

### 4.1 Unified Capabilities Schema (Auto-generated)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "GRULE Capabilities Schema",
  "description": "Available capabilities for rule generation in the GRULE engine",
  "version": "1.0.0",
  "generated_at": "2026-01-06T12:00:00Z",
  
  "context_objects": {
    "packet": {
      "description": "Incoming data packet from device",
      "type": "object",
      "properties": {
        "IMEI": { "type": "string", "description": "Device identifier" },
        "Speed": { "type": "int64", "description": "Speed in km/h" },
        "Latitude": { "type": "float64", "description": "GPS latitude" },
        "Longitude": { "type": "float64", "description": "GPS longitude" },
        "Altitude": { "type": "int64", "description": "Altitude in meters" },
        "GSMSignalStrength": { "type": "int64", "description": "GSM signal 0-31" },
        "PositioningStatus": { "type": "string", "description": "A=valid, V=invalid" },
        "Datetime": { "type": "time.Time", "description": "Packet timestamp" },
        "Direction": { "type": "int64", "description": "Heading in degrees" },
        "NumberOfSatellites": { "type": "int64", "description": "Visible satellites" },
        "BufferUpdated": { "type": "bool", "description": "Flag: buffer was updated", "writable": true },
        "BufferHas10": { "type": "bool", "description": "Flag: buffer has 10 entries", "writable": true },
        "MetricsReady": { "type": "bool", "description": "Flag: metrics calculated", "writable": true },
        "AlertFired": { "type": "bool", "description": "Flag: alert was fired", "writable": true }
      }
    }
  },
  
  "capabilities": {
    "geo": {
      "name": "geofence",
      "version": "1.0.0",
      "description": "Geofence detection and timing",
      "functions": [
        {
          "name": "IsInsideGroup",
          "signature": "geo.IsInsideGroup(groupName string, lat float64, lon float64) bool",
          "description": "Check if device is inside any geofence of a group",
          "example": "geo.IsInsideGroup(\"Taller\", packet.Latitude, packet.Longitude)"
        },
        {
          "name": "MinutesInsideGeofence",
          "signature": "geo.MinutesInsideGeofence(id string, inside bool) float64",
          "description": "Minutes spent inside a geofence",
          "example": "geo.MinutesInsideGeofence(\"zone1\", true) > 30"
        }
      ]
    },
    
    "buffer": {
      "name": "buffer",
      "version": "1.0.0",
      "description": "Circular buffer for position history",
      "functions": [
        {
          "name": "UpdateMemoryBuffer",
          "signature": "buffer.Update(speed int64, gsm int64, datetime time.Time, status string, lat float64, lon float64) bool",
          "description": "Add entry to buffer, returns true when buffer has 10 entries",
          "example": "packet.BufferHas10 = buffer.Update(packet.Speed, packet.GSMSignalStrength, packet.Datetime, packet.PositioningStatus, packet.Latitude, packet.Longitude)"
        },
        {
          "name": "HasExactly10",
          "signature": "buffer.HasExactly10() bool",
          "description": "Check if buffer has exactly 10 positions"
        },
        {
          "name": "GetSize",
          "signature": "buffer.GetSize() int",
          "description": "Get current buffer size"
        }
      ]
    },
    
    "metrics": {
      "name": "metrics",
      "version": "1.0.0",
      "description": "Metric calculations from buffer data",
      "functions": [
        {
          "name": "GetAverageSpeed90Min",
          "signature": "metrics.GetAverageSpeed90Min() int64",
          "description": "Average speed from buffer entries within 90 minutes"
        },
        {
          "name": "GetAverageGSMLast5",
          "signature": "metrics.GetAverageGSMLast5() int64",
          "description": "Average GSM signal from last 5 buffer entries"
        }
      ],
      "state_variables": [
        {
          "name": "AvgSpeed90min",
          "type": "int64",
          "description": "Cached average speed (90 min window)",
          "readable": true,
          "writable": true
        },
        {
          "name": "AvgGsm5",
          "type": "int64",
          "description": "Cached average GSM (last 5)",
          "readable": true,
          "writable": true
        }
      ]
    },
    
    "timing": {
      "name": "timing",
      "version": "1.0.0",
      "description": "Time-based detection functions",
      "functions": [
        {
          "name": "IsOfflineFor",
          "signature": "timing.IsOfflineFor(minutes int) bool",
          "description": "Check if device has been offline for X minutes",
          "example": "timing.IsOfflineFor(5)"
        }
      ]
    },
    
    "alerts": {
      "name": "alerts",
      "version": "1.0.0",
      "description": "Alert management and spam prevention",
      "functions": [
        {
          "name": "IsAlertSent",
          "signature": "alerts.IsAlertSent(alertID string) bool",
          "description": "Check if a specific alert was already sent"
        },
        {
          "name": "MarkAlertSent",
          "signature": "alerts.MarkAlertSent(alertID string) bool",
          "description": "Mark an alert as sent to prevent spam"
        },
        {
          "name": "ResetAlert",
          "signature": "alerts.ResetAlert(alertID string) bool",
          "description": "Reset alert flag to allow re-sending"
        }
      ]
    },
    
    "state": {
      "name": "state",
      "version": "1.0.0",
      "description": "Persistent state management (counters, flags)",
      "functions": [
        {
          "name": "GetCounter",
          "signature": "state.GetCounter(name string) int64",
          "description": "Get current value of a named counter"
        },
        {
          "name": "IncCounter",
          "signature": "state.IncCounter(name string) int64",
          "description": "Increment counter and return new value"
        },
        {
          "name": "ResetCounter",
          "signature": "state.ResetCounter(name string) bool",
          "description": "Reset counter to zero"
        }
      ]
    },
    
    "actions": {
      "name": "actions",
      "version": "1.0.0",
      "description": "Output actions (alerts, commands)",
      "functions": [
        {
          "name": "SendTelegram",
          "signature": "actions.SendTelegram(message string)",
          "description": "Send message via Telegram",
          "example": "actions.SendTelegram(\"Alert: \" + packet.IMEI)"
        },
        {
          "name": "SendEmail",
          "signature": "actions.SendEmail(subject string, body string)",
          "description": "Send email notification"
        },
        {
          "name": "Log",
          "signature": "actions.Log(message string)",
          "description": "Log a message for debugging"
        },
        {
          "name": "CutEngine",
          "signature": "actions.CutEngine(imei string)",
          "description": "Send engine cut command to device"
        },
        {
          "name": "Audit",
          "signature": "actions.Audit(ruleName string, description string, salience int64, alertFired bool)",
          "description": "Record rule execution for auditing"
        },
        {
          "name": "CastString",
          "signature": "actions.CastString(value interface{}) string",
          "description": "Convert any value to string for concatenation"
        }
      ]
    }
  },
  
  "rule_template": {
    "description": "Template for creating a new rule",
    "structure": {
      "format": "rule RuleName \"Description\" salience N { when CONDITION then ACTIONS }",
      "salience_guide": {
        "1000": "Highest priority - initialization rules",
        "900": "High priority - validation rules",
        "800": "Medium-high - calculation rules",
        "700": "Medium - evaluation rules",
        "600": "Low - alert/action rules"
      }
    },
    "example": "rule SpeedAlert \"Alert when speed exceeds limit\" salience 700 {\n    when\n        packet.Speed > 120 && !alerts.IsAlertSent(\"speed_alert\")\n    then\n        actions.SendTelegram(\"Speed alert: \" + actions.CastString(packet.Speed) + \" km/h\");\n        alerts.MarkAlertSent(\"speed_alert\");\n}"
  },
  
  "composition_patterns": {
    "sequential_defcon": {
      "description": "Sequential validation pattern (like jammer detection)",
      "pattern": [
        "DEFCON0 (salience 1000): Initialize/update buffer",
        "DEFCON1 (salience 900): First condition check",
        "DEFCON2 (salience 800): Calculate metrics (if DEFCON1 passed)",
        "DEFCON3 (salience 700): Additional validations",
        "DEFCON4 (salience 600): Final alert (if all conditions passed)"
      ]
    },
    "guard_pattern": {
      "description": "Check alert not already sent before sending",
      "example": "!alerts.IsAlertSent(\"alert_id\") && CONDITION"
    }
  }
}
```

### 4.2 LLM Prompt Template

```markdown
# Rule Generation Context

You are generating a GRL rule for the GRULE engine. Below is the capabilities schema.

## Available Objects

### packet (IncomingPacket)
- `packet.IMEI` (string) - Device identifier
- `packet.Speed` (int64) - Speed in km/h
- `packet.Latitude` (float64) - GPS latitude
- `packet.Longitude` (float64) - GPS longitude
- `packet.GSMSignalStrength` (int64) - Signal strength 0-31
- `packet.PositioningStatus` (string) - "A" (valid) or "V" (invalid)
- `packet.Datetime` (time.Time) - Packet timestamp

### Flags (read/write on packet)
- `packet.BufferUpdated` (bool)
- `packet.BufferHas10` (bool)
- `packet.MetricsReady` (bool)
- `packet.AlertFired` (bool)

## Available Capabilities

### geo (Geofence)
- `geo.IsInsideGroup(groupName, lat, lon)` → bool
- `geo.MinutesInsideGeofence(id, inside)` → float64

### buffer
- `buffer.Update(speed, gsm, datetime, status, lat, lon)` → bool
- `buffer.HasExactly10()` → bool

### metrics
- `metrics.GetAverageSpeed90Min()` → int64
- `metrics.GetAverageGSMLast5()` → int64

### timing
- `timing.IsOfflineFor(minutes)` → bool

### alerts
- `alerts.IsAlertSent(alertID)` → bool
- `alerts.MarkAlertSent(alertID)` → bool

### state
- `state.GetCounter(name)` → int64
- `state.IncCounter(name)` → int64

### actions
- `actions.SendTelegram(message)`
- `actions.Log(message)`
- `actions.CastString(value)` → string
- `actions.Audit(ruleName, description, salience, alertFired)`

## User Request
{USER_REQUEST}

## Generate a GRL rule following this format:
```grl
rule RuleName "Description" salience N {
    when
        CONDITION
    then
        ACTIONS;
}
```
```

---

## 5. Implementation Roadmap

### Phase 1: Foundation (Week 1-2)

1. **Create capability interface** (`engine/capabilities/interface.go`)
2. **Create base registry** (`engine/capabilities/registry.go`)
3. **Extract geofence capability** (from `persistent_state.go`)
4. **Extract buffer capability** (from `memory_buffer.go`)

### Phase 2: Core Capabilities (Week 3-4)

5. **Extract timing capability** (offline detection)
6. **Extract metrics capability** (averages calculation)
7. **Extract alerts capability** (from `alerts.go` + spam guard from `persistent_state.go`)
8. **Extract state capability** (counters, flags)

### Phase 3: GRULE Integration (Week 5)

9. **Refactor context builder** (auto-inject all capabilities)
10. **Simplify grule_worker.go** (remove scattered logic)
11. **Update existing rules** to new naming convention

### Phase 4: LLM Integration (Week 6)

12. **Create schema generator** (YAML → JSON)
13. **Create example rules library**
14. **Document prompt templates**

### Phase 5: IoT Extensibility (Future)

15. **Create IoT adapter interface**
16. **Implement temperature sensor capability**
17. **Implement humidity sensor capability**

---

## 6. Pattern Analysis: Factory vs Rule Engine

### 6.1 Factory Pattern

```go
// CapabilityFactory creates configured capability instances
type CapabilityFactory struct {
    db       *sql.DB
    config   *Config
}

func (f *CapabilityFactory) Create(capType string, imei string) Capability {
    switch capType {
    case "geofence":
        return geofence.New(f.db)
    case "buffer":
        return buffer.New(imei)
    case "metrics":
        return metrics.New()
    default:
        return nil
    }
}
```

**Verdict:** ✅ **Compatible with GRULE** - Use for creating capability instances per request.

### 6.2 Registry Pattern (Recommended)

```go
// CapabilityRegistry manages all available capabilities
type CapabilityRegistry struct {
    capabilities map[string]Capability
    mu           sync.RWMutex
}

func (r *CapabilityRegistry) Register(cap Capability) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.capabilities[cap.Name()] = cap
}

func (r *CapabilityRegistry) BuildDataContext() *ast.DataContext {
    dc := ast.NewDataContext()
    for name, cap := range r.capabilities {
        dc.Add(cap.GetDataContextName(), cap)
    }
    return dc
}
```

**Verdict:** ✅ **Best for GRULE** - Auto-registers capabilities into DataContext.

### 6.3 Rule Engine Pattern (GRULE Native)

GRULE already IS a rule engine. Our architecture should complement it, not compete.

**Our approach:**
- **Capabilities** = Building blocks (functions/state)
- **GRULE** = Rule execution engine
- **Registry** = Connects capabilities to GRULE DataContext

### 6.4 Pattern Recommendation

| Pattern | Use Case | GRULE Compatibility |
|---------|----------|---------------------|
| **Registry** | Capability discovery & DataContext building | ✅ Excellent |
| **Factory** | Creating capability instances with config | ✅ Good |
| **Strategy** | Swappable algorithms (e.g., distance calculation) | ✅ Good |
| **Observer** | Not needed - GRULE handles event flow | ⚠️ Redundant |

**Final Recommendation:** Use **Registry + Factory** pattern combination.

---

## 7. Conclusion

The **Capability-Oriented Architecture (Scenario D)** provides the optimal balance of:

1. **Modularity** - True "Lego blocks" that are independently testable
2. **LLM Friendliness** - YAML manifests + auto-generated JSON schema
3. **GRULE Compatibility** - Registry pattern for seamless DataContext integration
4. **Extensibility** - Easy to add new domains (IoT) without modifying core
5. **Maintainability** - Small, focused files (~100-200 lines each)

### Key Takeaways

- **Split by functionality, not by technical layer**
- **Each capability owns its state, functions, and manifest**
- **Registry auto-builds GRULE DataContext**
- **JSON schema enables LLM rule generation**
- **YAML manifests are human-editable and version-controllable**

---

## Appendix A: Migration Checklist

| Current File | Target Capability | Priority |
|--------------|-------------------|----------|
| `persistent_state.go` (IsInsideGroup, etc.) | `geofence/` | High |
| `persistent_state.go` (SecondsInsideGeofence) | `geofence/` | High |
| `memory_buffer.go` | `buffer/` | High |
| `persistent_state.go` (GetAverageSpeed90Min) | `metrics/` | Medium |
| `persistent_state.go` (IsOfflineFor) | `timing/` | Medium |
| `persistent_state.go` (IsAlertSent, MarkAlertSent) | `alerts/` | High |
| `persistent_state.go` (Counters) | `state/` | Medium |
| `alerts.go` (SendTelegram, etc.) | `alerts/` or `actions/` | High |
| `property.go` | Keep as-is (execution flags) | Low |

---

## Appendix B: File Size Guidelines

| File Type | Max Lines | Rationale |
|-----------|-----------|-----------|
| Capability main | 150 | Single responsibility |
| Functions file | 200 | Grouped related functions |
| Manifest YAML | 100 | Human readable |
| Test file | 300 | Comprehensive coverage |

---

*End of Proposal*
