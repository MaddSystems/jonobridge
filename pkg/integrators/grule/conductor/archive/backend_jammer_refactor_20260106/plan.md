# Implementation Plan: Backend Jammer Refactor

**Track ID:** `backend_jammer_refactor_20260106`  
**Date:** January 6, 2026  
**Estimated Duration:** 7 weeks

---

## Phase 0: Foundation Setup (Week 1 - Days 1-4)

### TODO 0.1: Create Backend Directory Structure
**Status:** `[x]` Complete
**Effort:** 2 hours  

```bash
mkdir -p backend/capabilities/{geofence,buffer,timing,metrics,alerts}
mkdir -p backend/grule
mkdir -p backend/adapters
mkdir -p backend/persistence
mkdir -p backend/audit
mkdir -p backend/schema
```

**CRITICAL:** The `backend/` folder is a **COMPLETE STANDALONE SOLUTION** with its own `go.mod`, `main.go`, `Dockerfile`, etc. It does NOT share anything with the original code.

---

### TODO 0.2: Create Backend go.mod
**Status:** `[x]` Complete
**Effort:** 30 minutes  

**File:** `backend/go.mod`

```go
module github.com/jonobridge/grule-backend

go 1.24

require (
    github.com/hyperjumptech/grule-rule-engine v1.15.0
    github.com/go-sql-driver/mysql v1.7.1
    // other dependencies
)
```

---

### TODO 0.3: Create Backend main.go
**Status:** `[x]` Complete
**Effort:** 3 hours  

**File:** `backend/main.go`

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    log.Println("Starting GRULE Backend (Jammer) on port 8081...")
    
    // Initialize capabilities
    // Initialize MySQL connection
    // Initialize GRULE engine
    // Start HTTP server
    
    http.ListenAndServe(":8081", nil)
}
```

**Port:** 8081 (same as original - replacement strategy)

---

### TODO 0.4: Create Backend Dockerfile
**Status:** `[x]` Complete
**Effort:** 1 hour  

**File:** `backend/Dockerfile`

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o grule-backend

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/grule-backend .
EXPOSE 8081
CMD ["./grule-backend"]
```

---

### TODO 0.5: Create Backend build.sh
**Status:** `[x]` Complete
**Effort:** 30 minutes  

**File:** `backend/build.sh`

```bash
#!/bin/bash
set -e

echo "Building backend (jammer-only)..."
CGO_ENABLED=0 go build -o grule-backend

echo "Building Docker image..."
docker build -t maddsystems/grule-backend:1.0.0 .

echo "Done!"
```

---

### TODO 0.6: Define Capability Interface
**Status:** `[x]` Complete
**Effort:** 4 hours  

**File:** `backend/capabilities/interface.go`

```go
package capabilities

type Capability interface {
    Name() string
    Version() string
    GetDataContextName() string
    Initialize(imei string) error
    GetSnapshot() map[string]interface{}
}
```

---

### TODO 0.7: Create Capability Registry
**Status:** `[x]` Complete
**Effort:** 3 hours  

**File:** `backend/capabilities/registry.go`

**Functions:**
- [ ] `NewRegistry() *Registry`
- [ ] `Register(cap Capability) error`
- [ ] `Get(name string) Capability`
- [ ] `BuildDataContext(imei string) *ast.DataContext`

---

### TODO 0.8: Create Persistence Interface
**Status:** `[~]` In Progress
**Effort:** 2 hours  

**File:** `backend/persistence/mysql.go`
- [ ] Copy connection logic from `engine/rule_loader.go:initMySQL()`
- [ ] Implement `MySQLStateStore` using existing `vehicle_rule_state` table

---

## Phase 1: Geofence Capability (Week 1 - Days 5-7)

### TODO 1.1: Extract Geofence Capability (JAMMER ONLY)
**Status:** `[ ]` Not Started  
**Effort:** 4 hours  
**Source:** `engine/persistent_state.go`

**IMPLEMENT (used by jammer_wargames.grl):**
| Function | Source Lines | Target File |
|----------|--------------|-------------|
| `IsInsideCircle()` | ~89-100 | `backend/capabilities/geofence/functions.go` |
| `IsInsideGroup()` | ~361-427 | `backend/capabilities/geofence/functions.go` |

**DO NOT IMPLEMENT (not used):**
- `EnteredGeofence()` - not in jammer rule
- `SecondsInsideGeofence()` - not in jammer rule
- `MinutesInsideGeofence()` - not in jammer rule
- `HoursInsideGeofence()` - not in jammer rule

**Files to Create:**
- [ ] `backend/capabilities/geofence/capability.go`
- [ ] `backend/capabilities/geofence/functions.go`
- [ ] `backend/capabilities/geofence/manifest.yaml`

---

### TODO 1.2: Create Geofence Manifest
**Status:** `[ ]` Not Started  
**Effort:** 1 hour  

**File:** `backend/capabilities/geofence/manifest.yaml`

```yaml
name: geofence
version: "1.0.0"
description: Geofence detection (jammer-only subset)
grule_context_name: geo

functions:
  - name: IsInsideGroup
    description: Check if device is inside any geofence of a named group
    parameters:
      - name: groupName
        type: string
      - name: lat
        type: float64
      - name: lon
        type: float64
    return_type: bool
    example: 'geo.IsInsideGroup("Taller", packet.Latitude, packet.Longitude)'
```

---

## Phase 2: Buffer Capability (Week 2 - Days 1-3)

### TODO 2.1: Extract Buffer Capability
**Status:** `[x]` Complete
**Effort:** 4 hours  
**Source:** `engine/memory_buffer.go`

**IMPLEMENT (all needed for jammer):**
| Struct/Function | Source Lines | Target File |
|-----------------|--------------|-------------|
| `BufferEntry` | ~9-18 | `backend/capabilities/buffer/types.go` |
| `FixedCircularBuffer` | ~20-25 | `backend/capabilities/buffer/circular.go` |
| `BufferManager` | ~27-31 | `backend/capabilities/buffer/manager.go` |
| `AddEntry()` | ~35-78 | `backend/capabilities/buffer/circular.go` |
| `GetEntriesInTimeWindow90Min()` | ~81-99 | `backend/capabilities/buffer/circular.go` |
| `InitializeFixedBuffers()` | ~134-146 | `backend/capabilities/buffer/manager.go` |
| `GetOrCreateBuffer()` | ~149-165 | `backend/capabilities/buffer/manager.go` |
| `cleanupInactiveBuffers()` | ~183-203 | `backend/capabilities/buffer/manager.go` |
| Public functions | ~207-280 | `backend/capabilities/buffer/manager.go` |

**Files to Create:**
- [ ] `backend/capabilities/buffer/capability.go`
- [ ] `backend/capabilities/buffer/types.go`
- [ ] `backend/capabilities/buffer/circular.go`
- [ ] `backend/capabilities/buffer/manager.go`
- [ ] `backend/capabilities/buffer/manifest.yaml`

---

## Phase 3: Metrics Capability (Week 2 - Days 4-5)

### TODO 3.1: Extract Metrics Capability
**Status:** `[x]` Complete
**Effort:** 3 hours  
**Source:** `engine/memory_buffer.go`, `engine/persistent_state.go`

**IMPLEMENT (used by jammer_wargames.grl):**
| Function | Source File | Lines |
|----------|-------------|-------|
| `CalculateAverageSpeed90Min()` | memory_buffer.go | ~103-115 |
| `CalculateAverageGSMLast5()` | memory_buffer.go | ~118-132 |
| `GetAverageSpeed90Min()` | persistent_state.go | ~651-653 |
| `GetAverageGSMLast5()` | persistent_state.go | ~655-657 |

**State Variables:**
- `AvgSpeed90min` (int64)
- `AvgGsm5` (int64)

**Files to Create:**
- [ ] `backend/capabilities/metrics/capability.go`
- [ ] `backend/capabilities/metrics/averages.go`
- [ ] `backend/capabilities/metrics/manifest.yaml`

---

## Phase 4: Timing Capability (Week 3 - Days 1-2)

### TODO 4.1: Extract Timing Capability
**Status:** `[x]` Complete
**Effort:** 2 hours  
**Source:** `engine/persistent_state.go`

**IMPLEMENT (used by jammer_wargames.grl):**
| Function | Lines |
|----------|-------|
| `IsOfflineFor()` | ~621-642 |
| `lastValidPosition` field | ~44 |
| `currentPacketTime` field | ~45 |

**Files to Create:**
- [ ] `backend/capabilities/timing/capability.go`
- [ ] `backend/capabilities/timing/offline.go`
- [ ] `backend/capabilities/timing/manifest.yaml`

---

## Phase 5: Alerts Capability (Week 3 - Days 3-5)

### TODO 5.1: Extract Alerts Capability (JAMMER ONLY)
**Status:** `[ ]` Not Started  
**Effort:** 4 hours  
**Source:** `actions/actions.go`, `actions/alerts.go`, `engine/persistent_state.go`

**IMPLEMENT (used by jammer_wargames.grl):**
| Function | Source File | Lines |
|----------|-------------|-------|
| `MarkAlertSent()` | persistent_state.go | ~193-205 |
| `SendTelegram()` | alerts.go | ~38-60 |
| `Log()` | alerts.go | ~91-115 |
| `Audit()` | actions.go | ~28-63 |
| `CastString()` | actions.go | ~93-95 |

**State Variable:**
- `JammerAlertSent` (bool)

**DO NOT IMPLEMENT (not used):**
- `SendEmail()` - not in jammer rule
- `SendWebhook()` - not in jammer rule
- `CutEngine()` - not in jammer rule
- `RestoreEngine()` - not in jammer rule
- `SendRawHex()` - not in jammer rule
- `IsAlertSent()` - rule reads field directly
- `ResetAlert()` - not in jammer rule
- `Concat()` - rule uses + operator

**Files to Create:**
- [ ] `backend/capabilities/alerts/capability.go`
- [ ] `backend/capabilities/alerts/channels.go`
- [ ] `backend/capabilities/alerts/spam_guard.go`
- [ ] `backend/capabilities/alerts/manifest.yaml`

---

## Phase 6: GRULE Integration Layer (Week 4)

### TODO 6.1: Create Context Builder
**Status:** `[ ]` Not Started  
**Effort:** 4 hours  
**Source:** `engine/grule_worker.go:executeRulesForPacket()`

**File:** `backend/grule/context_builder.go`

```go
// Build DataContext with all registered capabilities
// IMPORTANT: Use same context names for template compatibility
dc := ast.NewDataContext()
dc.Add("IncomingPacket", wrapper)
dc.Add("state", state)        // Same name as original
dc.Add("actions", actionsHelper)  // Same name as original
```

---

### TODO 6.2: Create Simplified Worker
**Status:** `[ ]` Not Started  
**Effort:** 5 hours  
**Source:** `engine/grule_worker.go`

**File:** `backend/grule/worker.go`

**Goals:**
- Simplify from ~700 lines to ~200 lines
- Keep WorkerPool pattern
- Remove all unused flag handling

---

### TODO 6.3: Create IncomingPacket Struct (MINIMAL)
**Status:** `[ ]` Not Started  
**Effort:** 1 hour  
**Source:** `engine/grule_worker.go` (PacketWrapper)

**File:** `backend/grule/packet.go`

**IMPLEMENT (used by jammer_wargames.grl):**
```go
type IncomingPacket struct {
    IMEI              string
    Speed             int64
    GSMSignalStrength int64
    Datetime          time.Time
    PositioningStatus string
    Latitude          float64
    Longitude         float64
    
    BufferUpdated           bool
    BufferHas10             bool
    IsOfflineFor5Min        bool
    PositionInvalidDetected bool
    MetricsReady            bool
    MovingWithWeakSignal    bool
    OutsideAllSafeZones     bool
}
```

**DO NOT IMPLEMENT (not used):**
- All `*Processed` flags
- All `*Failed` flags
- `JammerPatternFullyConfirmed`
- `CurrentlyInvalid`
- `EvaluationSkipped`
- `AlertFired`
- `HistoryUpdated`

---

## Phase 7: Adapters (Week 5 - Days 1-3)

### TODO 7.1: Create GPS Tracker Adapter
**Status:** `[ ]` Not Started  
**Effort:** 3 hours  
**Source:** `engine/grule_worker.go:ProcessPacketMessage()`

**File:** `backend/adapters/gps_tracker.go`

**Logic to Extract:**
- JSON parsing from `ProcessPacketMessage()` (~315-360)
- Speed conversion (m/s to km/h)
- Positioning status validation

---

## Phase 8: Schema Generation (Week 5 - Days 4-5)

### TODO 8.1: Create Schema Generator
**Status:** `[ ]` Not Started  
**Effort:** 4 hours  

**File:** `backend/schema/generator.go`

**Functions:**
- [ ] `GenerateFromManifests(dir string) ([]byte, error)`
- [ ] Parse all `manifest.yaml` files
- [ ] Generate unified `capabilities.json`

---

## Phase 9: Audit Module (Week 6 - Days 1-2)

### TODO 9.1: Copy Audit Module
**Status:** `[ ]` Not Started  
**Effort:** 2 hours  
**Source:** `engine/audit/`

**Files to Create:**
- [ ] `backend/audit/capture.go` - Copy from `engine/audit/capture.go`
- [ ] `backend/audit/db.go` - Copy from `engine/audit/db.go`
- [ ] `backend/audit/types.go` - Copy from `engine/audit/types.go`

---

## Phase 10: Integration Testing (Week 6 - Days 3-7)

### TODO 10.1: Integration Test
**Status:** `[ ]` Not Started  
**Effort:** 4 hours  

**Test Cases:**
- [ ] `TestFullJammerDetectionFlow` - End-to-end with jammer_wargames.grl
- [ ] `TestDEFCON0to4Progression` - All rules fire in order
- [ ] `TestTelegramAlertOnDEFCON4` - Alert sends correctly
- [ ] `TestAPICompatibility` - Same endpoints as original

---

### TODO 10.2: Replacement Validation Test
**Status:** `[ ]` Not Started  
**Effort:** 2 hours  

**Verify:**
- [ ] New backend runs on port 8081
- [ ] API endpoints match original exactly
- [ ] Same packets processed correctly
- [ ] Can swap container image without frontend changes
- [ ] Rollback to old image works if needed

---

## Phase 11: API & Documentation (Week 7 - Days 5-7)

**Note:** Frontend is a separate Flask app (not containerized). Backend exposes Swagger API.

### TODO 11.1: Add Schema Endpoint
**Status:** `[ ]` Not Started  
**Effort:** 2 hours  

**New Endpoint:** `GET /api/schema/capabilities`

```json
{
  "version": "1.0.0",
  "capabilities": {
    "geo": { "functions": [...] },
    "buffer": { "functions": [...] }
  }
}
```

---

## Quick Reference: Source → Target Mapping (JAMMER ONLY)

| Original File | Lines | Target Capability | Target File |
|---------------|-------|-------------------|-------------|
| `persistent_state.go` | 89-100 | geofence | `geofence/functions.go` |
| `persistent_state.go` | 361-427 | geofence | `geofence/functions.go` |
| `persistent_state.go` | 193-205 | alerts | `alerts/spam_guard.go` |
| `persistent_state.go` | 621-642 | timing | `timing/offline.go` |
| `persistent_state.go` | 583-603 | buffer | (wrapper) |
| `persistent_state.go` | 651-657 | metrics | `metrics/averages.go` |
| `memory_buffer.go` | 9-203 | buffer | `buffer/*.go` |
| `memory_buffer.go` | 103-132 | metrics | `metrics/averages.go` |
| `alerts.go` | 38-60 | alerts | `alerts/channels.go` |
| `alerts.go` | 91-115 | alerts | `alerts/channels.go` |
| `actions.go` | 28-63 | alerts | `alerts/capability.go` |
| `actions.go` | 93-95 | alerts | `alerts/capability.go` |
| `grule_worker.go` | 73-207 | grule | `grule/worker.go` |
| `grule_worker.go` | 378-515 | grule | `grule/executor.go` |

---

## Files NOT to Create (not used by jammer_wargames.grl)

| Capability | Functions | Reason |
|------------|-----------|--------|
| geofence | `EnteredGeofence`, `*InsideGeofence` timing | Not in rule |
| state | `IncCounter`, `GetCounter`, `ResetCounter` | Not in rule |
| alerts | `SendEmail`, `SendWebhook`, `CutEngine`, `RestoreEngine` | Not in rule |
| commands | All MQTT commands | Not in rule |

---

## Code Reduction Summary

| Component | `engine/` lines | `backend/` lines | Reduction |
|-----------|-----------------|------------------|-----------|
| geofence | ~150 | ~50 | 67% |
| buffer | ~280 | ~200 | 29% |
| metrics | ~50 | ~40 | 20% |
| timing | ~30 | ~25 | 17% |
| alerts | ~250 | ~80 | 68% |
| grule/worker | ~693 | ~200 | 71% |
| **Total** | **~1,453** | **~595** | **59%** |

---

## Checklist

### Phase 0: Foundation (STANDALONE INFRASTRUCTURE)
- [ ] Create directory structure
- [ ] **Create `backend/go.mod`** (independent module)
- [ ] **Create `backend/main.go`** (port 8081)
- [ ] **Create `backend/Dockerfile`** (standalone container)
- [ ] **Create `backend/build.sh`**
- [ ] Define Capability interface
- [ ] Create Registry
- [ ] Create persistence layer

### Phase 1-5: Capabilities
- [ ] Geofence capability (IsInsideGroup only)
- [ ] Buffer capability (full)
- [ ] Metrics capability (full)
- [ ] Timing capability (IsOfflineFor only)
- [ ] Alerts capability (4 functions only)

### Phase 6-8: Integration
- [ ] Context builder
- [ ] Simplified worker
- [ ] GPS tracker adapter
- [ ] Schema generator

### Phase 9: Audit
- [ ] Copy audit module

### Phase 10-11: Testing & Documentation
- [ ] Integration tests
- [ ] Replacement validation test (same port 8081)
- [ ] Schema endpoint
- [ ] Documentation

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Breaking existing rules | Original code 100% UNTOUCHED - available as reference |
| Frontend incompatibility | API endpoints remain identical, same port |
| Performance regression | Benchmark before deploying to production |
| New implementation fails | Keep old code, rollback to old container image |
| Database conflicts | Both use same tables, only one runs at a time |

---

## Verification Commands

```bash
# Build new backend
cd backend
go build -o grule-backend

# Run new backend locally (same port as original)
./grule-backend  # Listens on :8081

# Test API
curl http://localhost:8081/api/rules

# Docker build
cd backend
docker build -t maddsystems/grule-backend:1.0.0 .

# Deployment is managed by jonobridge
```
