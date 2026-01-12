# Backend Refactor: Jammer Wargames

**Track ID:** `backend_jammer_refactor_20260106`  
**Date:** January 6, 2026  
**Status:** Planning  
**Estimated Duration:** 7 weeks  
**Priority:** High

---

## 🎯 Objective

Create a **NEW complete standalone solution** in `backend/` folder implementing only what `jammer_wargames.grl` needs. No garbage code.

**CRITICAL: Original code stays 100% UNTOUCHED**
- DO NOT modify `main.go`
- DO NOT modify `Dockerfile`
- DO NOT modify `engine/`
- DO NOT modify `actions/`
- DO NOT modify `build.sh` or `deploy.sh`

Goals:
- LLM-powered rule generation via JSON schema
- Independent testing of each capability
- Easy extension for IoT devices beyond GPS trackers
- Smaller, focused files (~100-200 lines each)
- Complete standalone deployment

---

## System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         FRONTEND (Flask)                            │
│  frontend/main.py                                                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────┐ │
│  │ Rule Templates  │  │ Rule CRUD UI    │  │ Audit Dashboards    │ │
│  │ rules_templates/│  │ form.html       │  │ audit*.html         │ │
│  │ *.grl files     │  │ index.html      │  │ progress_audit.html │ │
│  └────────┬────────┘  └────────┬────────┘  └─────────┬───────────┘ │
│           │                    │                      │             │
│           └────────────────────┼──────────────────────┘             │
│                                │ HTTP API Calls                     │
└────────────────────────────────┼────────────────────────────────────┘
                                 │
         ┌───────────────────────┴───────────────────────┐
         │                                               │
         ▼                                               ▼
┌─────────────────────────────────┐    ┌─────────────────────────────────┐
│   ORIGINAL (100% UNTOUCHED)     │    │   NEW BACKEND (STANDALONE)      │
│   grule/                        │    │   grule/backend/                │
│   ├── main.go                   │    │   ├── main.go        ← NEW      │
│   ├── Dockerfile                │    │   ├── Dockerfile     ← NEW      │
│   ├── engine/                   │    │   ├── build.sh       ← NEW      │
│   ├── actions/                  │    │   ├── deploy.sh      ← NEW      │
│   └── ... (unchanged)           │    │   ├── go.mod         ← NEW      │
│                                 │    │   ├── capabilities/             │
│   Port: 8081                    │    │   ├── grule/                    │
│   Port: 8081                    │    │   Port: 8081 (same)             │
│                                 │    │                                 │
│                                 │    │                                 │
└─────────────────────────────────┘    └─────────────────────────────────┘
```

### Replacement Strategy
1. **Original** (`grule/`) - Stays 100% as-is, serves as reference during development
2. **New Backend** (`grule/backend/`) - Complete standalone, same port 8081

### Safe Implementation
- Build new implementation in `backend/`
- If new works → replace old implementation
- If new fails → keep old implementation
- Old code available for reference during development

---

## 📁 Target Directory Structure

```
grule/
├── Dockerfile           # ❌ DO NOT TOUCH
├── build.sh             # ❌ DO NOT TOUCH
├── deploy.sh            # ❌ DO NOT TOUCH
├── go.mod               # ❌ DO NOT TOUCH
├── go.sum               # ❌ DO NOT TOUCH
├── main.go              # ❌ DO NOT TOUCH
├── engine/              # ❌ DO NOT TOUCH
├── actions/             # ❌ DO NOT TOUCH
│
├── backend/             # ✅ NEW COMPLETE STANDALONE SOLUTION
│   ├── main.go              # NEW - Entry point (port 8081)
│   ├── Dockerfile           # NEW - Standalone container build
│   ├── build.sh             # NEW - Build script
│   ├── go.mod               # NEW - Independent Go module
│   ├── go.sum               # NEW - Dependencies
│   │
│   ├── capabilities/
│   │   ├── interface.go     # Capability interface definition
│   │   ├── registry.go      # Capability registry
│   │   ├── geofence/
│   │   │   ├── capability.go
│   │   │   ├── functions.go   # IsInsideGroup, IsInsideCircle
│   │   │   └── manifest.yaml
│   │   ├── buffer/
│   │   │   ├── capability.go
│   │   │   ├── circular.go    # FixedCircularBuffer
│   │   │   ├── manager.go     # BufferManager
│   │   │   └── manifest.yaml
│   │   ├── metrics/
│   │   │   ├── capability.go
│   │   │   ├── averages.go    # GetAverageSpeed90Min, GetAverageGSMLast5
│   │   │   └── manifest.yaml
│   │   ├── timing/
│   │   │   ├── capability.go
│   │   │   ├── offline.go     # IsOfflineFor
│   │   │   └── manifest.yaml
│   │   └── alerts/
│   │       ├── capability.go
│   │       ├── channels.go    # SendTelegram, Log
│   │       ├── spam_guard.go  # MarkAlertSent
│   │       └── manifest.yaml
│   │
│   ├── grule/
│   │   ├── context_builder.go
│   │   ├── executor.go
│   │   ├── worker.go
│   │   ├── loader.go
│   │   └── packet.go          # IncomingPacket struct
│   │
│   ├── adapters/
│   │   ├── interface.go
│   │   └── gps_tracker.go
│   │
│   ├── persistence/
│   │   ├── interface.go
│   │   └── mysql.go
│   │
│   ├── audit/
│   │   ├── capture.go
│   │   ├── db.go
│   │   └── types.go
│   │
│   ├── schema/
│   │   ├── generator.go
│   │   └── capabilities.json   # Auto-generated for LLM
│
└── frontend/            # SEPARATE APP (not containerized)
    ├── main.py
    └── ...
```

---

## 🐳 Containerization (Docker)

### ORIGINAL STAYS UNTOUCHED
```
grule/
├── Dockerfile      # ❌ DO NOT MODIFY - keeps building original
├── build.sh        # ❌ DO NOT MODIFY
└── deploy.sh       # ❌ DO NOT MODIFY
```

### NEW STANDALONE DOCKER IN backend/

**File:** `backend/Dockerfile`
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o grule-backend

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/grule-backend .
EXPOSE 8081
CMD ["./grule-backend"]
```

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

### Replacement Containers
| Container | Image | Port | Source | Status |
|-----------|-------|------|--------|--------|
| Original | `maddsystems/grule:x.x.x` | 8081 | `grule/` | Reference |
| New | `maddsystems/grule-backend:x.x.x` | 8081 | `grule/backend/` | Replacement |

**Note:** Kubernetes Deployment and Service are managed by jonobridge, not created here.

---

## ✅ IMPLEMENT: Functions Used by `jammer_wargames.grl`

### From `state` (PersistentState)

| Function | Signature | Rule |
|----------|-----------|------|
| `UpdateMemoryBuffer` | `(speed, gsm int64, datetime time.Time, posStatus string, lat, lon float64) bool` | DEFCON0 |
| `IsOfflineFor` | `(minutes int64) bool` | DEFCON0 |
| `GetAverageSpeed90Min` | `(imei string) int64` | DEFCON2 |
| `GetAverageGSMLast5` | `(imei string) int64` | DEFCON2 |
| `IsInsideGroup` | `(groupName string, lat, lon float64) bool` | DEFCON3 |
| `MarkAlertSent` | `(alertID string) bool` | DEFCON4 |

| Public Field | Type | Rule |
|--------------|------|------|
| `JammerAvgSpeed90min` | `int64` | DEFCON2 |
| `JammerAvgGsm5` | `int64` | DEFCON2 |
| `JammerAlertSent` | `bool` | DEFCON4 |

### From `actions` (ActionsHelper)

| Function | Signature | Rule |
|----------|-----------|------|
| `Log` | `(message string)` | ALL |
| `Audit` | `(ruleName, desc string, salience int64, alertFired bool)` | DEFCON1-4 |
| `CastString` | `(v interface{}) string` | DEFCON2, DEFCON4 |
| `SendTelegram` | `(message string)` | DEFCON4 |

### From `IncomingPacket` (PacketWrapper)

| Field | Type | Usage |
|-------|------|-------|
| `IMEI` | `string` | then |
| `Speed` | `int64` | then |
| `GSMSignalStrength` | `int64` | then |
| `Datetime` | `time.Time` | then |
| `PositioningStatus` | `string` | when |
| `Latitude` | `float64` | when, then |
| `Longitude` | `float64` | when, then |
| `BufferUpdated` | `bool` | when, then |
| `BufferHas10` | `bool` | when, then |
| `IsOfflineFor5Min` | `bool` | when, then |
| `PositionInvalidDetected` | `bool` | when, then |
| `MetricsReady` | `bool` | when, then |
| `MovingWithWeakSignal` | `bool` | when, then |
| `OutsideAllSafeZones` | `bool` | when, then |

---

## ❌ DO NOT IMPLEMENT: Unused Code

### Actions (not used in jammer_wargames.grl)
- `SendEmail()`
- `SendWebhook()`
- `CutEngine()`
- `RestoreEngine()`
- `SendRawHex()`
- `Concat()`

### State Functions (not used)
- `EnteredGeofence()`
- `SecondsInsideGeofence()`
- `MinutesInsideGeofence()`
- `HoursInsideGeofence()`
- `IsAlertSent()` - rule reads `JammerAlertSent` directly
- `ResetAlert()`
- `IncCounter()`
- `GetCounter()`
- `ResetCounter()`
- `SetCounter()`
- `UpdateJammerHistoryExact()`
- `CalculateJammerMetricsIfReady()`
- `GetBotonPanicoExecuted()`
- `SetBotonPanicoExecuted()`
- `GetCond2Checked()` / `SetCond2Checked()`
- `GetCond3Checked()` / `SetCond3Checked()`
- `GetCond4Checked()` / `SetCond4Checked()`

### State Fields (not used)
- `BotonPanicoExecuted`
- `ExcesoVelocidadExecuted`
- `JammerPositions`
- `JammerProcessed`
- `Cond2Checked`, `Cond3Checked`, `Cond4Checked`
- `prevInside` map

### Packet Fields (not used)
- `DebugProcessed`
- `ResetProcessed`
- `PositionInvalidDetectedProcessed`
- `PositionInvalidDetectedFailed`
- `MovingWithWeakSignalProcessed`
- `MovingWithWeakSignalFailed`
- `OutsideAllSafeZonesProcessed`
- `OutsideAllSafeZonesFailed`
- `JammerPatternFullyConfirmed`
- `JammerPatternFullyConfirmedProcessed`
- `JammerPatternFullyConfirmedFailed`
- `CurrentlyInvalid`
- `EvaluationSkipped`
- `AlertFired`
- `HistoryUpdated`

### Entire Files (not needed)
- `actions/commands.go` - MQTT commands not used

---

## 📊 Code Reduction

| Component | `engine/` lines | `backend/` lines | Reduction |
|-----------|-----------------|------------------|-----------|
| state.go | ~657 | ~150 | 77% |
| buffer.go | ~280 | ~200 | 29% |
| actions.go | ~100 | ~40 | 60% |
| alerts.go | ~130 | ~50 | 62% |
| commands.go | ~110 | 0 | 100% |
| packet.go | ~250 | ~50 | 80% |
| worker.go | ~693 | ~200 | 71% |
| **Total** | **~2,220** | **~690** | **69%** |

---

## 🔗 Internal Dependencies

```
buffer.go
    └── Used by state.go (UpdateMemoryBuffer, GetAverage*)

state.go
    ├── Uses buffer.go
    └── Uses MySQL (geofence queries, alert flags)

actions.go
    ├── Uses audit/ (Audit function)
    └── Uses Telegram API (SendTelegram)

packet.go
    └── Standalone struct

worker.go
    ├── Uses state.go
    ├── Uses actions.go
    ├── Uses packet.go
    └── Uses grule-rule-engine
```

---

## 🖥️ Frontend Integration

**Note:** Frontend is a **separate Flask app** (not containerized). Backend exposes Swagger API endpoints.

### API Endpoints (Same API, Same Port)

Both backends expose the **same API endpoints** on port 8081:

| Endpoint | Original | New Backend |
|----------|----------|-------------|
| `GET /api/rules` | ✅ | ✅ |
| `POST /api/rules` | ✅ | ✅ |
| `POST /api/validate` | ✅ | ✅ |
| `GET /api/audit/*` | ✅ | ✅ |
| `GET /api/schema/capabilities` | ❌ | ✅ NEW |

### Frontend Configuration

**File:** `frontend/main.py`

```python
# Same URL for both backends (same port 8081)
BACKEND_URL = "http://localhost:8081"
```

### No Frontend Code Changes Required
- API is identical
- Same port (8081)
- Just deploy new container to replace old one

### Template Compatibility

The `jammer_wargames.grl` template works **as-is** with new backend:
```grl
// These context names are preserved:
state.UpdateMemoryBuffer(...)
state.IsInsideGroup("Taller", ...)
state.JammerAvgSpeed90min
actions.SendTelegram(...)
```

---

## ✅ Success Criteria

- [ ] `backend/` folder is a complete standalone solution
- [ ] Original code 100% unchanged (no modifications to `main.go`, `Dockerfile`, `engine/`, etc.)
- [ ] `backend/go.mod` is independent Go module
- [ ] `cd backend && go build` compiles without errors
- [ ] `backend/Dockerfile` builds standalone container
- [ ] `jammer_wargames.grl` executes successfully with new backend
- [ ] DEFCON 0→4 progression works correctly
- [ ] Telegram alert fires on DEFCON 4
- [ ] Audit trail recorded
- [ ] API endpoints identical to original (frontend compatible)
- [ ] JSON schema auto-generated from YAML manifests
- [ ] New backend can replace old backend (same port 8081)

---

## Key Constraints

1. **Original code 100% UNTOUCHED** - Do NOT modify anything outside `backend/`
2. **Complete standalone solution** - Own `main.go`, `Dockerfile`, `go.mod`
3. **Replacement deployment** - New container replaces old when ready
4. **API compatible** - Same endpoints as original for frontend compatibility
5. **Jammer-only scope** - Only implement functions used by `jammer_wargames.grl`
6. **Same port (8081)** - Drop-in replacement for original backend
7. **Safe rollback** - Keep old code as reference, can revert if needed
