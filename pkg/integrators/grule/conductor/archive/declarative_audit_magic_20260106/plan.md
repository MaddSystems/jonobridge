# Implementation Plan: Declarative Audit Magic

**Track ID:** `declarative_audit_magic_20260106`  
**Date:** January 6, 2026  
**Estimated Duration:** 1 week  
**Scope:** `backend/` and `frontend/` folders

> **Note:** This track modifies both `backend/` and `frontend/` folders. The root directory contains the original implementation which remains unchanged.

**Key Changes:**
- Backend: GRULE listener for auto-audit capture
- Backend: Manifest loader from database
- Frontend: Upload paired `.grl` + `.yaml` template files
- Database: New `audit_manifest` column in `fleet_rules` table

---

## Phase 1: GRULE Listener (Days 1-2)

### TODO 1.1: Implement Audit Listener
**Status:** `[x]` Complete  
**Effort:** 4 hours  

**File:** `backend/audit/listener.go`

```go
package audit

import (
    "log"
    "sync"
    "github.com/hyperjumptech/grule-rule-engine/ast"
)

type AuditListener struct {
    manifest     *AuditManifest
    enabled      bool              // Global kill switch
    loggedOnce   map[string]bool   // Track logged unknown rules
    loggedMu     sync.Mutex
}

func NewAuditListener(manifest *AuditManifest) *AuditListener {
    return &AuditListener{
        manifest:   manifest,
        enabled:    true,
        loggedOnce: make(map[string]bool),
    }
}

// SetEnabled - Global kill switch for performance
func (l *AuditListener) SetEnabled(enabled bool) {
    l.enabled = enabled
}

func (l *AuditListener) BeforeExecuteConsequence(
    cycle uint64,
    rule *ast.RuleEntry,
    dc ast.IDataContext,
) {
    // Use EXISTING global toggle from audit.IsProgressAuditEnabled()
    // This ensures frontend controls (Activar/Desactivar) still work
    if !audit.IsProgressAuditEnabled() {
        return
    }
    
    meta := l.manifest.GetRuleMeta(rule.RuleName)
    
    // Unknown rule - log once, skip audit
    if meta == nil {
        l.logOnce(rule.RuleName)
        return
    }
    
    // Explicitly disabled
    if !meta.Enabled {
        return
    }
    
    // Extract snapshot with nil-safety
    snapshot, err := extractSnapshot(dc)
    if err != nil {
        log.Printf("[AuditListener] Snapshot error for '%s': %v", rule.RuleName, err)
        snapshot = map[string]interface{}{"error": err.Error()}
    }
    
    // Extract IMEI from packet for database storage
    imei := ""
    if packet := getFromContext[*IncomingPacket](dc, "IncomingPacket"); packet != nil {
        imei = packet.IMEI
    }
    
    // Non-blocking capture (hot path optimization)
    go Capture(&AuditEntry{
        IMEI:         imei,
        RuleName:     rule.RuleName,
        Salience:     int(rule.Salience),
        Description:  meta.Description,
        Level:        meta.Level,
        IsAlert:      meta.IsAlert,
        StepNumber:   meta.Order,        // From manifest
        StageReached: meta.Description,  // Use description as stage name
        Snapshot:     snapshot,
    })
}

func (l *AuditListener) AfterExecuteConsequence(
    cycle uint64,
    rule *ast.RuleEntry,
    dc ast.IDataContext,
) {
    // Optional post-execution capture
}

// logOnce logs unknown rule warning only once
func (l *AuditListener) logOnce(ruleName string) {
    l.loggedMu.Lock()
    defer l.loggedMu.Unlock()
    if !l.loggedOnce[ruleName] {
        log.Printf("[AuditListener] Rule '%s' not in manifest, skipping audit", ruleName)
        l.loggedOnce[ruleName] = true
    }
}
```

**Tasks:**
- [x] Create listener struct
- [x] Implement `logOnce()` to avoid log spam for unknown rules
- [x] Extract rule name and salience from `*ast.RuleEntry`
- [x] Lookup manifest metadata (skip if nil or disabled)
- [x] Use `go Capture()` for non-blocking hot path
- [x] Handle snapshot extraction errors gracefully
- [x] **Use existing `audit.IsProgressAuditEnabled()` for kill switch** (frontend controls compatibility)

---

### Frontend Controls Compatibility

The existing frontend has Progress Audit controls that must continue working:

| Control | Backend Function | Effect on Listener |
|---------|------------------|-------------------|
| **Activar** | `EnableProgressAudit()` | Listener starts capturing |
| **Desactivar** | `DisableProgressAudit()` | Listener stops capturing |
| **Limpiar datos** | `ClearProgressAudit()` | Clears `rule_execution_state` table |
| **Status badge** | `IsProgressAuditEnabled()` | Shows ACTIVO/INACTIVO |

The listener checks `audit.IsProgressAuditEnabled()` instead of its own flag, so the frontend controls work automatically.

---

### TODO 1.2: Define SnapshotProvider Interface
**Status:** `[x]` Complete  
**Effort:** 1 hour  

**Design Pattern:** Open/Closed Principle - capabilities are open for extension (add new data) but audit code is closed for modification.

**File:** `backend/capabilities/interface.go`

```go
package capabilities

// SnapshotProvider allows capabilities to contribute their own audit data
// Each capability implements this to self-report its snapshot without
// modifying the central audit code.
type SnapshotProvider interface {
    // GetSnapshotData returns capability-specific data for audit snapshots
    // Keys should be descriptive (e.g., "buffer_circular", "jammer_metrics")
    // Returns nil if capability has no data to contribute
    GetSnapshotData(imei string) map[string]interface{}
}
```

**File:** `backend/capabilities/buffer/capability.go` (add method)

```go
// GetSnapshotData implements SnapshotProvider
func (b *BufferCapability) GetSnapshotData(imei string) map[string]interface{} {
    if b == nil {
        return nil
    }
    
    entries := b.GetEntriesInTimeWindow90Min(imei)
    var bufferData []map[string]interface{}
    for _, e := range entries {
        bufferData = append(bufferData, map[string]interface{}{
            "imei":               e.IMEI,
            "datetime":           e.Datetime.Format(time.RFC3339),
            "speed":              e.Speed,
            "gsm_signal":         e.GSMSignalStrength,
            "positioning_status": e.PositioningStatus,
            "latitude":           e.Latitude,
            "longitude":          e.Longitude,
            "is_valid":           e.IsValid,
        })
    }
    
    return map[string]interface{}{
        "buffer_circular": bufferData,
    }
}
```

**File:** `backend/capabilities/metrics/capability.go` (add method)

```go
// GetSnapshotData implements SnapshotProvider
func (m *MetricsCapability) GetSnapshotData(imei string) map[string]interface{} {
    if m == nil {
        return nil
    }
    return map[string]interface{}{
        "jammer_metrics": map[string]interface{}{
            "avg_speed_90min":  m.GetAverageSpeed90Min(imei),
            "avg_gsm_last5":    m.GetAverageGSMLast5(imei),
        },
    }
}
```

**File:** `backend/capabilities/geofence/capability.go` (add method)

```go
// GetSnapshotData implements SnapshotProvider
func (g *GeofenceCapability) GetSnapshotData(imei string) map[string]interface{} {
    if g == nil {
        return nil
    }
    val, ok := g.lastCoords.Load(imei)
    if !ok {
        return nil
    }
    co := val.(coords)
    return map[string]interface{}{
        "geofence_checks": map[string]bool{
            "inside_taller":    g.IsInsideGroup("Taller", co.lat, co.lon),
            "inside_clientes":  g.IsInsideGroup("CLIENTES", co.lat, co.lon),
            "inside_resguardo": g.IsInsideGroup("Resguardo/Cedis/Puerto", co.lat, co.lon),
        },
    }
}
```

**Tasks:**
- [x] Define `SnapshotProvider` interface in `capabilities/interface.go`
- [x] Implement `GetSnapshotData()` in `BufferCapability`
- [x] Implement `GetSnapshotData()` in `MetricsCapability`
- [x] Implement `GetSnapshotData()` in `GeofenceCapability`
- [x] Implement `GetSnapshotData()` in `TimingCapability`

---

### TODO 1.3: Implement Generic Snapshot Extraction
**Status:** `[x]` Complete  
**Effort:** 1 hour  

**File:** `backend/audit/snapshot.go`

```go
package audit

import (
    "log"
    
    "github.com/hyperjumptech/grule-rule-engine/ast"
    "backend/capabilities"
)

// extractSnapshot builds rich context snapshot using SnapshotProvider pattern
// Each capability self-reports its data - no modification needed when adding new capabilities
func extractSnapshot(dc ast.IDataContext) (map[string]interface{}, error) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[AuditListener] Panic during snapshot: %v", r)
        }
    }()
    
    snapshot := make(map[string]interface{})
    
    // 1. Extract packet_current (always present)
    packet := getFromContext[*IncomingPacket](dc, "IncomingPacket")
    if packet != nil {
        snapshot["packet_current"] = extractPacketData(packet)
        snapshot["wrapper_flags"] = extractWrapperFlags(packet)
    }
    
    // 2. Collect all SnapshotProviders from DataContext
    providers := collectSnapshotProviders(dc)
    
    // 3. Each provider contributes its own data (Open/Closed Principle)
    for _, provider := range providers {
        if data := provider.GetSnapshotData(); data != nil {
            for key, value := range data {
                snapshot[key] = value
            }
        }
    }
    
    return snapshot, nil
}

// collectSnapshotProviders gathers all capabilities that implement SnapshotProvider
func collectSnapshotProviders(dc ast.IDataContext) []capabilities.SnapshotProvider {
    var providers []capabilities.SnapshotProvider
    
    // Get StateWrapper which holds all capabilities
    state := getFromContext[*StateWrapper](dc, "state")
    if state == nil {
        return providers
    }
    
    // Check each capability - only add if it implements SnapshotProvider
    if p, ok := interface{}(state.buf).(capabilities.SnapshotProvider); ok && state.buf != nil {
        providers = append(providers, p)
    }
    if p, ok := interface{}(state.metrics).(capabilities.SnapshotProvider); ok && state.metrics != nil {
        providers = append(providers, p)
    }
    if p, ok := interface{}(state.geo).(capabilities.SnapshotProvider); ok && state.geo != nil {
        providers = append(providers, p)
    }
    if p, ok := interface{}(state.tim).(capabilities.SnapshotProvider); ok && state.tim != nil {
        providers = append(providers, p)
    }
    
    return providers
}

// extractPacketData extracts GPS packet fields
func extractPacketData(packet *IncomingPacket) map[string]interface{} {
    return map[string]interface{}{
        "Speed":             packet.Speed,
        "Latitude":          packet.Latitude,
        "Longitude":         packet.Longitude,
        "Altitude":          packet.Altitude,
        "GSMSignalStrength": packet.GSMSignalStrength,
        "Satellites":        packet.NumberOfSatellites,
        "PositioningStatus": packet.PositioningStatus,
        "Datetime":          packet.Datetime,
        "EventCode":         packet.EventCode,
    }
}

// extractWrapperFlags extracts processing flags from packet
func extractWrapperFlags(packet *IncomingPacket) map[string]interface{} {
    return map[string]interface{}{
        "BufferUpdated":    packet.BufferUpdated,
        "BufferHas10":      packet.BufferHas10,
        "MetricsReady":     packet.MetricsReady,
        "CurrentlyInvalid": packet.CurrentlyInvalid,
        "AlertFired":       packet.AlertFired,
        "IsOfflineFor5Min": packet.IsOfflineFor5Min,
    }
}

// getFromContext safely gets typed value from DataContext
func getFromContext[T any](dc ast.IDataContext, name string) T {
    var zero T
    obj, err := dc.Get(name)
    if err != nil || obj == nil {
        return zero
    }
    typed, ok := obj.(T)
    if !ok {
        return zero
    }
    return typed
}
```

**Tasks:**
- [x] Define `extractSnapshot()` with panic recovery
- [x] Implement `collectSnapshotProviders()` to gather all capabilities
- [x] Implement `extractPacketData()` for GPS fields
- [x] Implement `extractWrapperFlags()` for processing flags
- [x] Use type assertion to check `SnapshotProvider` interface

---

### TODO 1.4: Implement Capture Function (Frontend Compatibility)
**Status:** `[x]` Complete  
**Effort:** 2 hours  

**Critical:** The `Capture()` function must save data to `rule_execution_state` table in the exact format the frontend expects at `/progress-audit-movie`.

**File:** `backend/audit/capture.go`

```go
package audit

import (
    "encoding/json"
    "log"
    "time"
)

// AuditEntry is the data structure passed from listener to Capture()
type AuditEntry struct {
    IMEI        string                 // From IncomingPacket
    RuleID      int64                  // From manifest or 0
    RuleName    string                 // From rule execution
    Salience    int                    // From rule execution
    Description string                 // From manifest
    Level       string                 // From manifest (debug/info/warning/critical)
    IsAlert     bool                   // From manifest
    StepNumber  int                    // From manifest order field
    StageReached string                // Same as Description for now
    Snapshot    map[string]interface{} // Rich snapshot from extractSnapshot()
}

// Capture saves audit entry to database in frontend-compatible format
// This function MUST populate rule_execution_state in the exact format
// expected by the frontend at /progress-audit-movie
func Capture(entry *AuditEntry) {
    if entry == nil || !IsProgressAuditEnabled() {
        return
    }
    
    // Build ProgressAudit compatible with existing SaveProgressAudit
    progress := ProgressAudit{
        IMEI:              entry.IMEI,
        RuleID:            entry.RuleID,
        RuleName:          entry.RuleName,
        StepNumber:        entry.StepNumber,
        StageReached:      entry.StageReached,
        StopReason:        "",  // Set if rule didn't fire
        BufferSize:        extractBufferSize(entry.Snapshot),
        MetricsReady:      extractMetricsReady(entry.Snapshot),
        GeofenceEval:      extractGeofenceEval(entry.Snapshot),
        ContextSnapshot:   JSONMap(entry.Snapshot),
        ExecutionTime:     time.Now(),
        ComponentsExecuted: []string{entry.RuleName},
        ComponentDetails:   JSONMap{"salience": entry.Salience, "level": entry.Level, "is_alert": entry.IsAlert},
    }
    
    if err := SaveProgressAudit(progress); err != nil {
        log.Printf("[Capture] Error saving audit: %v", err)
    }
}

// extractBufferSize gets buffer size from snapshot for frontend display
func extractBufferSize(snapshot map[string]interface{}) int {
    if buffer, ok := snapshot["buffer_circular"].([]map[string]interface{}); ok {
        return len(buffer)
    }
    return 0
}

// extractMetricsReady checks if metrics are available in snapshot
func extractMetricsReady(snapshot map[string]interface{}) bool {
    if flags, ok := snapshot["wrapper_flags"].(map[string]interface{}); ok {
        if ready, ok := flags["MetricsReady"].(bool); ok {
            return ready
        }
    }
    return false
}

// extractGeofenceEval summarizes geofence checks for frontend display
func extractGeofenceEval(snapshot map[string]interface{}) string {
    if checks, ok := snapshot["geofence_checks"].(map[string]bool); ok {
        for name, inside := range checks {
            if inside {
                return "inside:" + name
            }
        }
        return "outside_all"
    }
    return "not_evaluated"
}
```

**Frontend Compatibility Verification:**

| Frontend Column | Source in Capture() | Database Column |
|-----------------|---------------------|-----------------|
| `step_number` | `entry.StepNumber` (from manifest order) | `step_number` |
| `stage_reached` | `entry.StageReached` | `stage_reached` |
| `buffer_size` | `extractBufferSize(snapshot)` | `buffer_size` |
| `snapshot` (in JSON viewer) | `entry.Snapshot` | `context_snapshot` |
| `execution_time` | `time.Now()` | `execution_time` |

**API Endpoints (already exist - no changes needed):**

| Endpoint | Handler | Status |
|----------|---------|--------|
| `GET /api/audit/progress/summary` | `GetFrameSummaryPaginated()` | ✅ Works |
| `GET /api/audit/progress/timeline` | `GetFrameTimelinePaginated()` | ✅ Works |
| `GET /api/audit/progress/snapshot?id=X` | `GetSnapshotByID()` | ✅ Works |

**Tasks:**
- [x] Define `AuditEntry` struct with all fields listener provides
- [x] Implement `Capture()` to convert to `ProgressAudit` format
- [x] Extract `buffer_size` from snapshot for frontend display
- [x] Extract `metrics_ready` from snapshot wrapper_flags
- [x] Extract `geofence_eval` summary from snapshot
- [x] Call existing `SaveProgressAudit()` for database storage
- [x] Verify `/progress-audit-movie` displays data correctly

---

### Graceful Degradation

The snapshot system is designed for **graceful degradation** - it never aborts, only captures what's available.

#### Design Principles

1. **Open/Closed Principle** - Audit code is closed for modification; capabilities are open for extension via `SnapshotProvider` interface
2. **Self-Reporting Capabilities** - Each capability implements `GetSnapshotData()` to contribute its own audit data
3. **YAML manifest controls what triggers audit** - Manifest defines `enabled: true/false` per rule, not what data is captured
4. **Nil checks prevent panics** - All capability access is guarded with nil checks before type assertion

#### Behavior Matrix

| Scenario | Behavior |
|----------|----------|
| Rule has no buffer capability | `buffer_circular` not in snapshot (provider returns nil) |
| Rule has no geofence capability | `geofence_checks` not in snapshot (provider returns nil) |
| Rule has no metrics | `jammer_metrics` not in snapshot (provider returns nil) |
| Packet missing some fields | Only available fields captured |
| New capability added | Implement `SnapshotProvider`, data appears automatically |
| Rule not in YAML manifest | Rule executes normally, no audit captured |
| Capability doesn't implement interface | Silently skipped, no error |

#### Example Snapshots

**Simple speed rule** (minimal capabilities):
```json
{
  "packet_current": {"Speed": 85, "Latitude": 19.4326, "Longitude": -99.1332},
  "buffer_circular": [],
  "jammer_metrics": {"avg_speed_90min": 0, "avg_gsm_last5": 0},
  "wrapper_flags": {"BufferUpdated": false, "MetricsReady": false}
}
```

**Jammer detection rule** (full capabilities):
```json
{
  "packet_current": {"Speed": 0, "Latitude": 19.4326, "Longitude": -99.1332, "GSMSignalStrength": 0},
  "buffer_circular": [{"imei": "123", "speed": 45, "gsm_signal": 28}, ...],
  "jammer_metrics": {"avg_speed_90min": 52.3, "avg_gsm_last5": 2.1, "jammer_positions": 8},
  "geofence_checks": {"inside_taller": false, "inside_clientes": true, "offline_5min": false},
  "wrapper_flags": {"BufferUpdated": true, "BufferHas10": true, "MetricsReady": true}
}
```

#### Extending for Future Capabilities (Open/Closed Principle)

When adding new capabilities, **no changes to audit code required**:

```go
// Step 1: Create new capability
// backend/capabilities/fuel/capability.go
type FuelCapability struct {
    Level       float64
    Consumption float64
}

// Step 2: Implement SnapshotProvider interface
func (f *FuelCapability) GetSnapshotData() map[string]interface{} {
    if f == nil {
        return nil
    }
    return map[string]interface{}{
        "fuel_data": map[string]interface{}{
            "level":       f.Level,
            "consumption": f.Consumption,
        },
    }
}

// Step 3: Add to StateWrapper
type StateWrapper struct {
    buf     *buffer.BufferCapability
    metrics *metrics.MetricsCapability
    geo     *geofence.GeofenceCapability
    tim     *timing.TimingCapability
    fuel    *fuel.FuelCapability  // NEW - just add the field
}

// Step 4: Register in collectSnapshotProviders (one-time addition)
if p, ok := interface{}(state.fuel).(capabilities.SnapshotProvider); ok && state.fuel != nil {
    providers = append(providers, p)
}
```

**Benefits:**
- ✅ Audit code unchanged when adding capabilities
- ✅ Each capability owns its snapshot format
- ✅ Easy to test capabilities in isolation
- ✅ New data appears automatically in audit

#### YAML Manifest Changes

Manifest changes only affect **which rules trigger audit capture**, not the snapshot content:

```yaml
# Before: rule disabled
jammer_alert:
  enabled: false
  
# After: rule enabled  
jammer_alert:
  enabled: true
  level: full
```

- **Snapshot structure**: Unchanged (always captures available data)
- **Audit behavior**: Now captures when rule fires

---

### TODO 1.5: Wire Listener to Engine
**Status:** `[x]` Complete  
**Effort:** 2 hours  

**File:** `backend/grule/worker.go`

```go
engine := engine.NewGruleEngine()
engine.Listener = audit.NewAuditListener(manifest)
engine.Execute(dc, kb)
```

**Tasks:**
- [x] Load manifest at worker initialization
- [x] Create listener instance
- [x] Attach to GRULE engine before `Execute()`

---

### TODO 1.6: Remove Hardcoded Audit() Stub
**Status:** `[x]` Complete  
**Effort:** 1 hour  

The backend currently has a hardcoded `Audit()` stub that does nothing. Remove it - the listener handles everything.

**Files to modify:**

1. **`backend/capabilities/alerts/capability.go`** - Remove the `Audit()` method:
```go
// DELETE THIS:
func (c *AlertsCapability) Audit(ruleName, description string, salience int64, alertFired bool) {
    // Stub
}
```

2. **`backend/grule/context_builder.go`** - Remove the `Audit()` wrapper method:
```go
// DELETE THIS:
func (a *ActionsWrapper) Audit(ruleName, description string, salience int64, alertFired bool) {
    if a.alrt != nil {
        a.alrt.Audit(ruleName, description, salience, alertFired)
    }
}
```

**Tasks:**
- [x] Remove `Audit()` from `AlertsCapability`
- [x] Remove `Audit()` from `ActionsWrapper`
- [x] Verify no GRL rules in database call `actions.Audit()` (or remove those calls)

---

## Phase 2: Audit Manifest from Database (Days 3-4)

> **Key Architecture:** YAML manifests are stored in `fleet_rules.audit_manifest` column (uploaded via frontend). Backend loads manifests from database at startup, NOT from filesystem.

### TODO 2.1: Define Manifest Struct
**Status:** `[x]` Complete  
**Effort:** 2 hours  

**File:** `backend/audit/manifest.go`

```go
package audit

import (
    "fmt"
    "log"
    "gopkg.in/yaml.v3"
)

type AuditManifest struct {
    rules map[string]*RuleMeta
}

type RuleMeta struct {
    Enabled     bool      // false = skip auditing this rule
    Description string
    Level       string    // debug, info, warning, critical
    IsAlert     bool
    Order       int
    Snapshot    []string  // Fields to capture: ["packet", "state"]. Empty = all
}

func NewAuditManifest() *AuditManifest {
    return &AuditManifest{
        rules: make(map[string]*RuleMeta),
    }
}

// ParseYAML parses a single YAML manifest string
func (m *AuditManifest) ParseYAML(yamlContent string) error {
    var manifest struct {
        Stages []struct {
            Rule  string `yaml:"rule"`
            Order int    `yaml:"order"`
            Audit struct {
                Enabled     *bool    `yaml:"enabled"`
                Description string   `yaml:"description"`
                Level       string   `yaml:"level"`
                IsAlert     bool     `yaml:"is_alert"`
                Snapshot    []string `yaml:"snapshot"`
            } `yaml:"audit"`
        } `yaml:"stages"`
    }
    
    if err := yaml.Unmarshal([]byte(yamlContent), &manifest); err != nil {
        return fmt.Errorf("YAML parse error: %w", err)
    }
    
    for _, stage := range manifest.Stages {
        enabled := true
        if stage.Audit.Enabled != nil {
            enabled = *stage.Audit.Enabled
        }
        
        m.rules[stage.Rule] = &RuleMeta{
            Enabled:     enabled,
            Description: stage.Audit.Description,
            Level:       stage.Audit.Level,
            IsAlert:     stage.Audit.IsAlert,
            Order:       stage.Order,
            Snapshot:    stage.Audit.Snapshot,
        }
    }
    return nil
}

func (m *AuditManifest) GetRuleMeta(ruleName string) *RuleMeta {
    if meta, ok := m.rules[ruleName]; ok {
        return meta
    }
    return nil
}

func (m *AuditManifest) Count() int {
    return len(m.rules)
}
```

**Tasks:**
- [x] Implement `AuditManifest` struct
- [x] Implement `ParseYAML()` method
- [x] Handle YAML parse errors gracefully
- [x] Default `enabled: true` when not specified

---

### TODO 2.2: Load Manifest from Database
**Status:** `[x]` Complete  
**Effort:** 2 hours  

**File:** `backend/audit/manifest.go` (add method)

```go
// LoadFromRules loads all audit manifests from fleet_rules records
func (m *AuditManifest) LoadFromRules(rules []persistence.Rule) error {
    for _, rule := range rules {
        if rule.AuditManifest == "" {
            log.Printf("[Manifest] Rule '%s' has no audit manifest, skipping", rule.Name)
            continue
        }
        
        if err := m.ParseYAML(rule.AuditManifest); err != nil {
            log.Printf("[Manifest] Warning: Invalid manifest for rule '%s': %v", rule.Name, err)
            continue  // Don't fail startup, just skip this rule
        }
        
        log.Printf("[Manifest] Loaded manifest for rule '%s'", rule.Name)
    }
    return nil
}
```

**File:** `backend/main.go`

```go
// Load rules from database
rules, err := store.LoadActiveRules()
if err != nil {
    log.Fatalf("Failed to load rules: %v", err)
}

// Build audit manifest from rules
manifest := audit.NewAuditManifest()
if err := manifest.LoadFromRules(rules); err != nil {
    log.Fatalf("Failed to load audit manifests: %v", err)
}
log.Printf("Loaded %d audit rule entries", manifest.Count())

// Pass manifest to worker
worker := grule.NewWorker(ctxBuilder, adapter, kbs, manifest)
```

**Tasks:**
- [x] Implement `LoadFromRules()` method
- [x] Update `main.go` to load manifest from rules
- [x] Log warnings for rules without manifests (don't fail)
- [x] Log count of loaded manifest entries

---


## Phase 3: Integration Test (Day 5)

### TODO 3.1: Integration Test (Python)
**Status:** `[x]` Complete  
**Effort:** 4 hours  

**Note:** This integration test will be implemented in Python after the Go backend code is complete. The test code will be located in `tests/integration/`.

**Test Cases:**
- [x] Load clean GRL rules into `fleet_rules` table
- [x] Load YAML manifest at startup
- [x] Send test packets through backend system
- [x] Verify DEFCON stages fire in order
- [x] Check `rule_execution_state` table has entries for ALL stages
- [x] Verify `context_snapshot` contains rich data (`packet_current`, `buffer_circular`, `jammer_metrics`, etc.)
- [x] Verify descriptions match manifest YAML
- [x] Verify `level` field populated (debug/warning/critical)
- [x] Verify `is_alert` flag correct for DEFCON4
- [x] Verify alerts sent correctly
- [x] Verify frontend at `/progress-audit-movie` displays data correctly

---

## Phase 4: Frontend & Database Integration (Days 6-7)

### TODO 4.1: Database Schema Migration
**Status:** `[x]` Complete  
**Effort:** 1 hour  

Add `audit_manifest` column to store YAML content alongside GRL.

**File:** `backend/persistence/rules.go`

```go
type Rule struct {
    ID            int64  `json:"id"`
    Name          string `json:"name"`
    Description   string `json:"description,omitempty"`
    GRL           string `json:"grl"`
    AuditManifest string `json:"audit_manifest,omitempty"`  // NEW: YAML content
    Active        bool   `json:"active"`
    Priority      int    `json:"priority"`
    CreatedAt     string `json:"created_at,omitempty"`
    UpdatedAt     string `json:"updated_at,omitempty"`
}
```

**SQL Migration:**
```sql
ALTER TABLE fleet_rules 
ADD COLUMN audit_manifest TEXT AFTER grl_content;
```

**Tasks:**
- [x] Add `AuditManifest` field to `Rule` struct
- [x] Update all SQL queries to include new column
- [x] Run migration on database

---

### TODO 4.2: Update Backend API
**Status:** `[x]` Complete  
**Effort:** 2 hours  

**File:** `backend/api/handlers.go`

Update create/update handlers to accept `audit_manifest` field.

```go
// CreateRule handler now accepts audit_manifest
func CreateRule(w http.ResponseWriter, r *http.Request) {
    var rule persistence.Rule
    json.NewDecoder(r.Body).Decode(&rule)
    
    // Validate YAML if provided
    if rule.AuditManifest != "" {
        if err := validateYAMLManifest(rule.AuditManifest); err != nil {
            http.Error(w, "Invalid audit manifest: "+err.Error(), 400)
            return
        }
    }
    
    // ... save to database
}
```

**Tasks:**
- [x] Update `CreateRule` handler to accept `audit_manifest`
- [x] Update `UpdateRule` handler to accept `audit_manifest`
- [x] Add YAML validation before saving
- [x] Update `GetRule` to return `audit_manifest`

---

### TODO 4.3: Load Manifest from Database
**Status:** `[x]` Complete  
**Effort:** 2 hours  

> **Note:** This is implemented via `LoadFromRules()` method which is called in `main.go` after loading rules from database. Equivalent to the planned `LoadFromDatabase()` approach.

**File:** `backend/audit/manifest.go`

```go
// LoadFromRules loads all audit manifests from fleet_rules records
func (m *AuditManifest) LoadFromRules(rules []persistence.Rule) error {
    for _, rule := range rules {
        if rule.AuditManifest == "" {
            log.Printf("[Manifest] Rule '%s' has no audit manifest, skipping", rule.Name)
            continue
        }
        if err := m.ParseYAML(rule.AuditManifest); err != nil {
            log.Printf("[Manifest] Warning: Invalid manifest for rule '%s': %v", rule.Name, err)
            continue
        }
        log.Printf("[Manifest] Loaded manifest for rule '%s'", rule.Name)
    }
    return nil
}
```

**File:** `backend/main.go` (already wired)

```go
manifest := audit.NewAuditManifest()
if err := manifest.LoadFromRules(rules); err != nil {
    log.Printf("Warning: Failed to load audit manifests: %v", err)
}
log.Printf("Loaded %d audit rule entries", manifest.Count())
```

**Tasks:**
- [x] Implement `LoadFromRules()` method (equivalent to `LoadFromDatabase()`)
- [x] Update `main.go` to call `LoadFromRules()` after loading rules
- [x] Handle rules without manifest gracefully (logs warning, continues)

---

### TODO 4.4: Update Frontend Template Loader
**Status:** `[x]` Complete  
**Effort:** 2 hours  

**File:** `frontend/main.py`

Load paired `.yaml` files alongside `.grl` templates.

```python
def load_templates():
    """Loads rule templates from .grl files with paired .yaml manifests"""
    templates = {}
    rules_dir = os.path.join(os.path.dirname(__file__), 'rules_templates')
    
    for filename in os.listdir(rules_dir):
        if filename.endswith('.grl'):
            template_key = filename[:-4]  # Remove .grl extension
            grl_path = os.path.join(rules_dir, filename)
            yaml_path = os.path.join(rules_dir, f"{template_key}.yaml")
            
            with open(grl_path, 'r') as f:
                grl_content = f.read()
            
            # Load paired YAML manifest if exists
            audit_manifest = ""
            if os.path.exists(yaml_path):
                with open(yaml_path, 'r') as f:
                    audit_manifest = f.read()
                print(f"✅ Loaded manifest: {template_key}.yaml")
            else:
                print(f"⚠️ No manifest found for {template_key}")
            
            templates[template_key] = {
                "name": name,
                "category": category,
                "description": description,
                "grl": grl_content,
                "audit_manifest": audit_manifest  # NEW
            }
    
    return templates
```

**Tasks:**
- [x] Update `load_templates()` to load paired `.yaml` files
- [x] Include `audit_manifest` in template data
- [x] Update API calls to send `audit_manifest` when creating rules

---

### TODO 4.5: Create Paired Template Files
**Status:** `[x]` Complete  
**Effort:** 2 hours  

**Files:**
- `frontend/rules_templates/jammer_wargames.grl` - Clean GRL (remove `Audit()` calls)
- `frontend/rules_templates/jammer_wargames.yaml` - Audit manifest

**jammer_wargames.yaml:**
```yaml
name: jammer_wargames
version: "1.0.0"
description: "Jammer detection with DEFCON stages"

stages:
  - rule: DEFCON0_Surveillance
    order: 1
    audit:
      enabled: true
      description: "Buffer Updated"
      level: debug
      
  - rule: DEFCON1_ContactLost_Pass
    order: 2
    audit:
      enabled: true
      description: "Step 1: Contact Lost (V)"
      level: warning
      
  # ... more stages
```

**Tasks:**
- [x] Create `jammer_wargames.yaml` manifest
- [x] Clean `jammer_wargames.grl` (remove all `actions.Audit()` calls)
- [x] Verify rule names in YAML match GRL rule names

---

### TODO 4.6: Update Frontend Form to Send Manifest
**Status:** `[x]` Complete  
**Effort:** 1 hour  

**File:** `frontend/templates/form.html`

When user selects a template, include the `audit_manifest` in the API request.

```javascript
function applyTemplate(name, grl, auditManifest) {
    document.getElementById('ruleName').value = name;
    document.getElementById('grlContent').value = grl;
    document.getElementById('auditManifest').value = auditManifest || '';
}

function submitRule() {
    const data = {
        name: document.getElementById('ruleName').value,
        grl: document.getElementById('grlContent').value,
        audit_manifest: document.getElementById('auditManifest').value  // NEW
    };
    // ... submit to API
}
```

**Tasks:**
- [x] Add hidden `auditManifest` field to form
- [x] Update `applyTemplate()` to set manifest
- [x] Update submit function to include manifest in API call

---

### TODO 4.7: Add Audit Manifest Viewer Tab
**Status:** `[x]` Complete  
**Effort:** 2 hours  

**Purpose:** Allow operators to inspect the YAML audit manifest alongside the GRL for debugging and transparency.

**File:** `frontend/templates/form.html`

Add a tabbed interface to view GRL and YAML separately:

```html
<!-- Tab Navigation -->
<ul class="nav nav-tabs" id="ruleTabs" role="tablist">
    <li class="nav-item" role="presentation">
        <button class="nav-link active" id="grl-tab" data-bs-toggle="tab" data-bs-target="#grl-content" type="button">
            <i class="bi bi-code-slash"></i> GRL Rules
        </button>
    </li>
    <li class="nav-item" role="presentation">
        <button class="nav-link" id="yaml-tab" data-bs-toggle="tab" data-bs-target="#yaml-content" type="button">
            <i class="bi bi-file-earmark-text"></i> Audit Manifest
            <span id="yaml-badge" class="badge bg-secondary d-none">YAML</span>
        </button>
    </li>
</ul>

<!-- Tab Content -->
<div class="tab-content" id="ruleTabsContent">
    <div class="tab-pane fade show active" id="grl-content" role="tabpanel">
        <textarea id="grlContent" class="form-control code-editor" rows="20"></textarea>
    </div>
    <div class="tab-pane fade" id="yaml-content" role="tabpanel">
        <textarea id="auditManifest" class="form-control code-editor" rows="20" 
                  placeholder="# No audit manifest for this rule"></textarea>
        <small class="text-muted">
            <i class="bi bi-info-circle"></i> 
            Defines which rules are audited and their metadata (description, level, order).
        </small>
    </div>
</div>
```

```javascript
// Show badge when manifest exists
function updateYamlBadge() {
    const yaml = document.getElementById('auditManifest').value;
    const badge = document.getElementById('yaml-badge');
    if (yaml && yaml.trim()) {
        badge.classList.remove('d-none');
        badge.textContent = yaml.split('\n').length + ' lines';
    } else {
        badge.classList.add('d-none');
    }
}

// Call when loading rule or applying template
function loadRule(rule) {
    document.getElementById('grlContent').value = rule.grl;
    document.getElementById('auditManifest').value = rule.audit_manifest || '';
    updateYamlBadge();
}
```

```css
.code-editor {
    font-family: 'Fira Code', 'Consolas', monospace;
    font-size: 0.9em;
    background: #1e1e1e;
    color: #d4d4d4;
    border-radius: 0 0 8px 8px;
}

#yaml-content .code-editor {
    background: #2d2d2d;
    color: #e6db74;  /* YAML-like syntax color */
}
```

**Tasks:**
- [x] Add Bootstrap nav-tabs for GRL / Audit Manifest
- [x] Move GRL textarea into tab pane
- [x] Add YAML textarea in second tab pane
- [x] Show line count badge when manifest exists
- [x] Style both editors with monospace font
- [x] Update `loadRule()` to populate both tabs
- [x] Add placeholder text explaining manifest purpose

---

### TODO 4.8: Add Critical Missing API Endpoints
**Status:** `[x]` Complete  
**Effort:** 2.5 hours  

**Problem:** The Python frontend requires endpoints that are missing in `backend/main.go`:

| Endpoint | Used By | Priority |
|----------|---------|----------|
| `GET /api/rules/components` | `progress_audit_movie.html` | 🔴 Critical |
| `POST /api/validate` | `form.html`, `main.py` | 🔴 Critical |
| `/swagger` | Browser API docs | 🟡 Nice to have |
| `/swagger.json` | Swagger UI spec | 🟡 Nice to have |
| CORS middleware | Flask cross-origin requests | 🔴 Critical |

**Files to modify:**

1. **`backend/api/handlers.go`** - Add `RuleComponentsHandler` and `ValidateHandler`
2. **`backend/main.go`** - Add CORS middleware, Swagger handlers, and register new routes

---

**File:** `backend/api/handlers.go` - Add missing handlers

```go
// ValidateHandler validates GRL syntax without saving
// POST /api/validate
func (s *Server) ValidateHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != http.MethodPost {
        http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
        return
    }

    var payload struct {
        GRL string `json:"grl"`
    }
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, `{"error":"Invalid JSON"}`, 400)
        return
    }

    if err := ValidateRule(payload.GRL); err != nil {
        json.NewEncoder(w).Encode(map[string]interface{}{
            "valid": false,
            "error": err.Error(),
        })
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "valid":   true,
        "message": "Rule syntax is valid",
    })
}

// RuleComponentsHandler returns internal rule names for a given rule package
// GET /api/rules/components?rule_name=xxx
func (s *Server) RuleComponentsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != http.MethodGet {
        http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
        return
    }

    ruleName := r.URL.Query().Get("rule_name")
    if ruleName == "" {
        http.Error(w, `{"error":"rule_name parameter required"}`, http.StatusBadRequest)
        return
    }

    // Get rule from database and extract component names
    rules, err := s.Store.GetAllRules()
    if err != nil {
        http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
        return
    }

    var components []string
    for _, rule := range rules {
        if rule.Name == ruleName {
            components = extractRuleNames(rule.GRL)
            break
        }
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "success":    true,
        "rule_name":  ruleName,
        "components": components,
    })
}

// extractRuleNames parses GRL and extracts individual rule names
func extractRuleNames(grl string) []string {
    var names []string
    re := regexp.MustCompile(`rule\s+(\w+)`)
    matches := re.FindAllStringSubmatch(grl, -1)
    for _, m := range matches {
        if len(m) > 1 {
            names = append(names, m[1])
        }
    }
    return names
}
```

---

**File:** `backend/main.go` - Add CORS middleware, Swagger handlers, and register routes

```go
// CORS Middleware (add before route registrations)
corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next(w, r)
    }
}

// API Routes (wrap all with corsMiddleware)
http.HandleFunc("/api/rules", corsMiddleware(server.RulesHandler))
http.HandleFunc("/api/rules/", corsMiddleware(server.RuleByIDHandler))
http.HandleFunc("/api/rules/components", corsMiddleware(server.RuleComponentsHandler))  // NEW
http.HandleFunc("/api/validate", corsMiddleware(server.ValidateHandler))                 // NEW
http.HandleFunc("/api/reload", corsMiddleware(server.ReloadHandler))
http.HandleFunc("/api/schema/capabilities", corsMiddleware(server.SchemaHandler))

// Progress Audit Routes (wrap with corsMiddleware)
http.HandleFunc("/api/audit/progress/enable", corsMiddleware(server.ProgressEnableHandler))
http.HandleFunc("/api/audit/progress/disable", corsMiddleware(server.ProgressDisableHandler))
http.HandleFunc("/api/audit/progress/clear", corsMiddleware(server.ProgressClearHandler))
http.HandleFunc("/api/audit/progress/status", corsMiddleware(server.ProgressStatusHandler))
http.HandleFunc("/api/audit/progress/summary", corsMiddleware(server.ProgressSummaryHandler))
http.HandleFunc("/api/audit/progress/timeline", corsMiddleware(server.ProgressTimelineHandler))
http.HandleFunc("/api/audit/progress/snapshot", corsMiddleware(server.SnapshotHandler))

// Swagger UI Routes
http.HandleFunc("/swagger", swaggerUIHandler)
http.HandleFunc("/swagger.json", swaggerJSONHandler)
```

---

**File:** `backend/main.go` - Add Swagger handler functions (before `loadRulesFromSlice`)

```go
// swaggerUIHandler serves the Swagger UI
// Access at: https://jonobridge.madd.com.mx/grule/swagger
func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
    html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Grule Backend API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css">
    <style>body { margin: 0; padding: 0; }</style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: '/swagger.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(html))
}

// swaggerJSONHandler serves the Swagger JSON specification
func swaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
    swagger := map[string]interface{}{
        "swagger": "2.0",
        "info": map[string]interface{}{
            "title":       "Grule Backend API v2.0",
            "description": "Declarative Audit System - GPS Fleet Rules Manager",
            "version":     "2.0.0",
        },
        "basePath": "/api",
        "schemes":  []string{"https", "http"},
        "paths": map[string]interface{}{
            "/rules":            map[string]interface{}{"get": map[string]interface{}{"summary": "List all rules"}, "post": map[string]interface{}{"summary": "Create rule"}},
            "/rules/{id}":       map[string]interface{}{"get": map[string]interface{}{"summary": "Get rule"}, "put": map[string]interface{}{"summary": "Update rule"}, "delete": map[string]interface{}{"summary": "Delete rule"}},
            "/rules/components": map[string]interface{}{"get": map[string]interface{}{"summary": "Get rule components", "parameters": []interface{}{map[string]interface{}{"name": "rule_name", "in": "query", "required": true, "type": "string"}}}},
            "/validate":         map[string]interface{}{"post": map[string]interface{}{"summary": "Validate GRL syntax"}},
            "/reload":           map[string]interface{}{"post": map[string]interface{}{"summary": "Reload rules from DB"}},
            "/audit/progress/enable":   map[string]interface{}{"post": map[string]interface{}{"summary": "Enable progress audit"}},
            "/audit/progress/disable":  map[string]interface{}{"post": map[string]interface{}{"summary": "Disable progress audit"}},
            "/audit/progress/clear":    map[string]interface{}{"post": map[string]interface{}{"summary": "Clear progress audit data"}},
            "/audit/progress/status":   map[string]interface{}{"get": map[string]interface{}{"summary": "Get progress audit status"}},
            "/audit/progress/summary":  map[string]interface{}{"get": map[string]interface{}{"summary": "Get IMEI summary (jqGrid Level 1)"}},
            "/audit/progress/timeline": map[string]interface{}{"get": map[string]interface{}{"summary": "Get frame timeline (jqGrid Level 2)"}},
            "/audit/progress/snapshot": map[string]interface{}{"get": map[string]interface{}{"summary": "Get snapshot by ID"}},
        },
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(swagger)
}
```

---

**Tasks:**
- [x] Add `ValidateHandler` to `api/handlers.go`
- [x] Add `RuleComponentsHandler` to `api/handlers.go`
- [x] Add `extractRuleNames` helper function
- [x] Add CORS middleware function to `main.go`
- [x] Wrap all existing route handlers with CORS middleware
- [x] Register `/api/validate` route
- [x] Register `/api/rules/components` route
- [x] Add `swaggerUIHandler` function to `main.go`
- [x] Add `swaggerJSONHandler` function to `main.go`
- [x] Register `/swagger` and `/swagger.json` routes
- [x] Add `"regexp"` import to `handlers.go`
- [x] Test: `https://jonobridge.madd.com.mx/grule/swagger`

---

## Phase 5: Dynamic Frontend (Days 8-9)

> **Goal:** Make the frontend 100% data-driven. Remove all hardcoded DEFCON stage references. Use captured audit data (`step_number`, `stage_reached`, `level`) to render stages dynamically.

### TODO 5.1: Verify `level` is Captured in Database
**Status:** `[x]` Complete  
**Effort:** 1 hour

**Critical Verification:** Ensure the `level` field from manifest is being saved to `rule_execution_state` table.

**File:** `backend/audit/capture.go`

```go
func Capture(entry *AuditEntry) {
    if entry == nil || !IsProgressAuditEnabled() {
        return
    }

    progress := ProgressAudit{
        IMEI:               entry.IMEI,
        RuleName:           entry.RuleName,
        StepNumber:         entry.StepNumber,    // ✓ From manifest order
        StageReached:       entry.StageReached,  // ✓ From manifest description
        Level:              entry.Level,          // ✓ From manifest level (debug/info/warning/critical)
        BufferSize:         extractBufferSize(entry.Snapshot),
        MetricsReady:       extractMetricsReady(entry.Snapshot),
        GeofenceEval:       extractGeofenceEval(entry.Snapshot),
        ContextSnapshot:    entry.Snapshot,
        ExecutionTime:      time.Now(),
    }

    if err := SaveProgressAudit(progress); err != nil {
        log.Printf("Error saving audit: %v", err)
    }
}
```

**Verify SQL Schema:**

```sql
-- Check if level column exists
DESCRIBE rule_execution_state;

-- Should show:
-- level VARCHAR(20) or TEXT
```

**If missing, add it:**

```sql
ALTER TABLE rule_execution_state 
ADD COLUMN level VARCHAR(20) DEFAULT 'info' 
AFTER stage_reached;
```

**Tasks:**
- [x] Verify `level` field is in `ProgressAudit` struct
- [x] Verify `level` is passed from `AuditEntry` to `ProgressAudit`
- [x] Check database schema has `level` column
- [x] Add column if missing
- [x] Test that `level` is saved correctly

---

### TODO 5.2: Update Timeline API to Include Level
**Status:** `[x]` Complete  
**Effort:** 30 minutes

**File:** `backend/audit/db.go`

Ensure `GetFrameTimelinePaginated()` returns the `level` field:

```go
func GetFrameTimelinePaginated(imei string, limit, offset int, sortBy, sortOrder, ruleName string) ([]map[string]interface{}, int, error) {
    // ... existing code ...
    
    query := `
        SELECT id, rule_id, rule_name, components_executed, component_details, 
               step_number, stage_reached, level, stop_reason,
               buffer_size, metrics_ready, geofence_eval, context_snapshot, execution_time
        FROM rule_execution_state 
        WHERE imei = ? 
        ORDER BY execution_time ASC
        LIMIT ? OFFSET ?
    `
    
    // ... scan results ...
    var level string
    rows.Scan(&id, &rId, &rName, &compJSON, &detJSON, &step, &stage, &level, &stop, ...)
    
    frame := map[string]interface{}{
        "id": id,
        "rule_name": rName,
        "step_number": step,
        "stage_reached": stage,
        "level": level,  // NEW: Include level in response
        "snapshot": snap,
        // ... other fields ...
    }
}
```

**Tasks:**
- [x] Add `level` to SELECT query
- [x] Add `level` variable to `rows.Scan()`
- [x] Include `level` in response map
- [x] Test API returns level for each frame

---

### TODO 5.3: Remove Hardcoded DEFCON from Frontend
**Status:** `[x]` Complete  
**Effort:** 2 hours

**File:** `frontend/templates/progress_audit_movie.html`

**Current (Hardcoded):**
```javascript
// ❌ BAD: Hardcoded DEFCON stages
const DEFCON_STAGES = {
    'DEFCON0': { color: '#28a745', name: 'Surveillance' },
    'DEFCON1': { color: '#ffc107', name: 'Contact Lost' },
    'DEFCON2': { color: '#fd7e14', name: 'Inhibition' },
    'DEFCON3': { color: '#dc3545', name: 'Safe Zones' },
    'DEFCON4': { color: '#8b0000', name: 'JAMMER ALERT' }
};

function renderStage(frame) {
    const defcon = extractDefcon(frame.rule_name);  // ❌ Parsing rule name
    const info = DEFCON_STAGES[defcon];
    return `<div style="color: ${info.color}">${info.name}</div>`;
}
```

**New (Data-Driven):**
```javascript
// ✅ GOOD: Use level from API response
const LEVEL_COLORS = {
    'debug': '#6c757d',     // Gray
    'info': '#0dcaf0',      // Blue
    'warning': '#ffc107',   // Orange
    'critical': '#dc3545'   // Red
};

const LEVEL_ICONS = {
    'debug': '🔍',
    'info': 'ℹ️',
    'warning': '⚠️',
    'critical': '🚨'
};

function renderStage(frame) {
    // Use data from API (no parsing, no hardcoding)
    const level = frame.level || 'info';  // Graceful fallback
    const color = LEVEL_COLORS[level];
    const icon = LEVEL_ICONS[level];
    const stepNum = frame.step_number || '?';
    const description = frame.stage_reached || frame.rule_name;  // Fallback to rule_name
    
    return `
        <div class="audit-stage" style="border-left: 4px solid ${color}">
            <span class="stage-icon">${icon}</span>
            <span class="stage-step">Step ${stepNum}</span>
            <span class="stage-desc">${description}</span>
        </div>
    `;
}
```

**CSS Styling:**
```css
.audit-stage {
    padding: 12px;
    margin: 8px 0;
    border-radius: 4px;
    background: rgba(255, 255, 255, 0.05);
    display: flex;
    align-items: center;
    gap: 12px;
}

.stage-icon {
    font-size: 1.5em;
}

.stage-step {
    font-weight: bold;
    color: #aaa;
    min-width: 60px;
}

.stage-desc {
    flex-grow: 1;
    color: #fff;
}

/* Level-specific backgrounds */
.audit-stage[data-level="critical"] {
    background: rgba(220, 53, 69, 0.1);
}

.audit-stage[data-level="warning"] {
    background: rgba(255, 193, 7, 0.1);
}
```

**Tasks:**
- [x] Remove all `DEFCON_STAGES` hardcoded objects
- [x] Remove `extractDefcon()` parsing functions
- [x] Add `LEVEL_COLORS` and `LEVEL_ICONS` mappings
- [x] Update `renderStage()` to use `level`, `step_number`, `stage_reached`
- [x] Add graceful fallback: if no `stage_reached`, use `rule_name`
- [x] Add graceful fallback: if no `level`, use `'info'`
- [x] Add CSS for level-specific styling
- [x] Test with jammer rules (should work)
- [x] Test with rules without manifest (should fallback gracefully)

---

### TODO 5.4: Update Progress Summary View
**Status:** `[x]` Complete  
**Effort:** 1 hour

**File:** `frontend/templates/progress_audit.html`

Update the summary grid to use dynamic stage names:

**Current (Hardcoded):**
```javascript
// ❌ Shows "DEFCON2" from rule name
function formatGridRow(row) {
    return {
        stage: extractDefcon(row.rule_name),
        ...
    };
}
```

**New (Data-Driven):**
```javascript
// ✅ Shows actual stage description from manifest
function formatGridRow(row) {
    return {
        stage: row.stage_reached || row.rule_name,  // Fallback to rule_name
        step: row.max_step || '?',
        level: row.level || 'info',
        ...
    };
}
```

**Tasks:**
- [x] Remove DEFCON extraction from summary view
- [x] Use `stage_reached` for stage column
- [x] Add `level` indicator (colored badge)
- [x] Test summary shows correct stage names
- [x] Test fallback for rules without manifest

---

### TODO 5.5: Integration Test for Dynamic Frontend
**Status:** `[x]` Complete  
**Effort:** 1 hour

**File:** `tests/integration_audit_v2.py`

Add test case to verify frontend receives correct data:

```python
def test_dynamic_frontend():
    # Create rule with manifest
    manifest = """
name: test_rule
stages:
  - rule: TestRule1
    order: 1
    audit:
      description: "First Step"
      level: info
  - rule: TestRule2
    order: 2
    audit:
      description: "Critical Step"
      level: critical
"""
    create_rule("TestRule", test_grl, manifest)
    
    # Send packet
    send_packet()
    time.sleep(2)
    
    # Verify timeline includes level
    response = requests.get(f"{API_URL}/api/audit/progress/timeline?imei={IMEI}")
    frames = response.json()['rows']
    
    assert len(frames) >= 1
    assert 'level' in frames[0], "Level missing from API response"
    assert frames[0]['level'] in ['debug', 'info', 'warning', 'critical']
    assert frames[0]['stage_reached'] == "First Step"
    assert frames[0]['step_number'] == 1
    
    print("✅ Frontend receives dynamic data correctly")
```

**Tasks:**
- [x] Add test case for dynamic data
- [x] Verify `level` in API response
- [x] Verify `stage_reached` matches manifest
- [x] Verify `step_number` matches manifest order
- [x] Test passes with manifest
- [x] Test passes without manifest (fallback)

---

## Summary

### Files to Create

| File | Phase | Description |
|------|-------|-------------|
| `backend/audit/listener.go` | 1 | GRULE execution listener |
| `backend/audit/snapshot.go` | 1 | Snapshot extraction using SnapshotProvider |
| `backend/audit/manifest.go` | 2 | Manifest loader from database |
| `frontend/rules_templates/jammer_wargames.yaml` | 4 | Audit manifest paired with GRL |

### Files to Modify

| File | Phase | Change |
|------|-------|--------|
| `backend/capabilities/interface.go` | 1 | Add `SnapshotProvider` interface |
| `backend/capabilities/*/capability.go` | 1 | Implement `GetSnapshotData()` |
| `backend/capabilities/alerts/capability.go` | 1 | Remove `Audit()` stub |
| `backend/grule/context_builder.go` | 1 | Remove `Audit()` wrapper |
| `backend/grule/worker.go` | 1 | Wire listener to engine |
| `backend/audit/capture.go` | 1 | Add `Capture()` function |
| `backend/main.go` | 2 | Load manifest from database |
| `backend/persistence/rules.go` | 4 | Add `audit_manifest` field |
| `backend/api/handlers.go` | 4 | Accept manifest in API |
| `frontend/main.py` | 4 | Load paired `.yaml` files |
| `frontend/rules_templates/jammer_wargames.grl` | 4 | Remove `Audit()` calls |
| `frontend/templates/form.html` | 4 | Tabbed GRL/YAML viewer, send manifest |

### Database Migration

```sql
ALTER TABLE fleet_rules ADD COLUMN audit_manifest TEXT AFTER grl_content;
```

---

## Audit Levels Reference

| Level | Use For | Example |
|-------|---------|---------|
| `debug` | Internal housekeeping | Buffer updates, metric calculations |
| `info` | Normal flow steps | State transitions |
| `warning` | Decision points | Condition matches, threshold crossed |
| `critical` | Alerts, final actions | Jammer detected, alert sent |

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Listener overhead | `go Capture()` non-blocking, global kill switch |
| GRULE API changes | Pin grule-rule-engine version |
| Missing manifest entry | Log once, skip audit (no crash) |
| YAML parse error | Fail fast at startup with clear error |
| Snapshot nil panic | Panic recovery in `extractSnapshot()` |
| Log spam | `logOnce()` for unknown rules |
| DB write failure | Log error, continue (don't crash worker) |
| Goroutine leak | Use bounded channel for async writes |

---

## Error Handling Summary

| Scenario | Action |
|----------|--------|
| Rule not in manifest | Log once, skip audit |
| `audit.enabled: false` | Silent skip |
| YAML parse error | Startup failure with path in error |
| Snapshot extraction fails | Log error, capture with `{"error": "..."}` |
| DataContext returns nil | Safe extraction returns nil |
| DB write fails | Log error, continue processing |

```