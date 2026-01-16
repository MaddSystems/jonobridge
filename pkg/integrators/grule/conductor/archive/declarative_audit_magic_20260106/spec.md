# Declarative Audit Magic

**Track ID:** `declarative_audit_magic_20260106`  
**Date:** January 6, 2026  
**Status:** Planning  
**Estimated Duration:** 1 week  
**Priority:** High  
**Scope:** `backend/` and `frontend/` folders

> **Note:** This track modifies both `backend/` and `frontend/` folders. The root directory contains the original implementation which remains unchanged.

**Key Changes:**
- Backend: GRULE listener for auto-audit capture
- Backend: Manifest loader from database  
- Frontend: Upload paired `.grl` + `.yaml` template files
- Database: New `audit_manifest` column in `fleet_rules` table

---

## 🎯 Objective

Eliminate hardcoded `actions.Audit()` calls from GRL rules by implementing a **declarative audit system** that automatically captures rule execution with rich metadata.

### Goals
- Rules contain **pure business logic** only
- Audit metadata defined separately in YAML (paired with GRL)
- GRULE listener auto-captures all rule executions
- Remove hardcoded `Audit()` stub from backend code
- Frontend uploads paired `.grl` + `.yaml` files
- Manifests stored in database alongside rules

---

## The Problem

### Current State (Hardcoded)
```grl
rule DEFCON1_ContactLost_Pass salience 900 {
    when
        IncomingPacket.PositioningStatus == "V"
    then
        IncomingPacket.PositionInvalidDetected = true;
        actions.Log("[DEFCON 1] Positive: Status V detected.");
        actions.Audit("DEFCON1_ContactLost_Pass", "Step 1: Contact Lost", 900, false);  // ❌ Hardcoded
}
```

### Problems
1. Rule writer must remember to add `Audit()` calls
2. Audit strings duplicated (rule name appears twice)
3. Business logic mixed with observability
4. LLM must generate audit calls correctly
5. Changing audit format requires editing all rules

---

## The Solution

### Target State (Declarative)

**Rule YAML Definition:**
```yaml
name: jammer_wargames
version: "1.0.0"
description: "Jammer detection with DEFCON stages"

stages:
  - rule: DEFCON0_Surveillance
    order: 1
    audit:
      description: "Surveillance Active"
      level: info
      
  - rule: DEFCON1_ContactLost_Pass
    order: 2
    audit:
      description: "Step 1: Contact Lost (V)"
      level: warning
      
  - rule: DEFCON4_JammerAlert_Fire
    order: 5
    audit:
      description: "CRITICAL: JAMMER DETECTED"
      level: critical
      is_alert: true
```

**Clean GRL (No Audit Calls):**
```grl
rule DEFCON1_ContactLost_Pass salience 900 {
    when
        IncomingPacket.PositioningStatus == "V"
    then
        IncomingPacket.PositionInvalidDetected = true;
        actions.Log("[DEFCON 1] Positive: Status V detected.");
        // ✅ No Audit() call - handled by listener
}
```

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         RULE DEFINITION                             │
│                                                                     │
│  ┌─────────────────────┐         ┌─────────────────────────────┐   │
│  │  Rule YAML          │         │  Audit Manifest             │   │
│  │  (stages, requires) │────────▶│  (descriptions, levels)     │   │
│  └──────────┬──────────┘         └──────────────┬──────────────┘   │
│             │                                    │                  │
│             │ Compiler                           │ Loaded at        │
│             ▼                                    │ startup          │
│  ┌─────────────────────┐                         │                  │
│  │  GRL Rules          │                         │                  │
│  │  (pure logic)       │                         │                  │
│  └──────────┬──────────┘                         │                  │
└─────────────┼────────────────────────────────────┼──────────────────┘
              │                                    │
              ▼                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         GRULE ENGINE                                │
│                                                                     │
│  ┌─────────────────────┐    ┌───────────────────────────────────┐  │
│  │  Knowledge Base     │───▶│  Audit Listener                   │  │
│  │  (compiled rules)   │    │  BeforeExecuteConsequence()       │  │
│  └─────────────────────┘    │  → Lookup manifest by rule name   │  │
│                             │  → Capture with rich metadata     │  │
│                             │  → Write to audit table           │  │
│                             └───────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         AUDIT OUTPUT                                │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  progress_audit table                                        │   │
│  │  + rule_name (from listener)                                 │   │
│  │  + salience (from listener)                                  │   │
│  │  + description (from manifest)                               │   │
│  │  + level (from manifest)                                     │   │
│  │  + is_alert (from manifest)                                  │   │
│  │  + context snapshot (from DataContext)                       │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 📁 New Files

```
backend/
├── audit/
│   ├── listener.go          # NEW - GRULE execution listener
│   ├── snapshot.go          # NEW - Snapshot extraction
│   └── manifest.go          # NEW - Audit manifest loader (from DB)

frontend/
└── rules_templates/
    └── jammer_wargames.yaml  # NEW - Paired audit manifest
```

## 🗑️ Files to Modify

| File | Change |
|------|--------|
| `backend/capabilities/alerts/capability.go` | Remove `Audit()` stub method |
| `backend/grule/context_builder.go` | Remove `Audit()` wrapper method |
| `backend/grule/worker.go` | Wire listener to engine |
| `backend/main.go` | Load manifest from database at startup |
| `backend/persistence/rules.go` | Add `AuditManifest` field to Rule struct |
| `backend/api/handlers.go` | Accept YAML manifest in create/update API |
| `frontend/main.py` | Load paired `.yaml` files with templates |
| `frontend/rules_templates/jammer_wargames.grl` | Remove `actions.Audit()` calls |
| `frontend/templates/form.html` | Add hidden manifest field, update JS |

## 🗄️ Database Migration

```sql
ALTER TABLE fleet_rules 
ADD COLUMN audit_manifest TEXT AFTER grl_content;
```

---

## YAML Schema Definition

### Rule Definition Schema

```yaml
# JSON Schema for rule YAML validation
$schema: "http://json-schema.org/draft-07/schema#"
type: object
required: [name, version, stages]

properties:
  name:
    type: string
    description: "Unique rule name"
    
  version:
    type: string
    pattern: "^\\d+\\.\\d+\\.\\d+$"
    
  description:
    type: string
    
  capabilities:
    type: array
    items:
      type: string
    description: "Required capabilities (geofence, buffer, etc.)"
    
  stages:
    type: array
    items:
      type: object
      required: [rule, order]
      properties:
        rule:
          type: string
          description: "GRL rule name"
        order:
          type: integer
          minimum: 1
        requires:
          type: array
          items:
            type: string
          description: "Previous stages that must pass"
        condition:
          type: string
          description: "High-level condition (for LLM generation)"
        audit:
          type: object
          properties:
            enabled:
              type: boolean
              default: true
              description: "Set to false to skip auditing this rule"
            description:
              type: string
            level:
              type: string
              enum: [debug, info, warning, critical]
            is_alert:
              type: boolean
              default: false
            snapshot:
              type: array
              items:
                type: string
              description: "Fields to capture: [packet, state, buffer]. Default: all"
              default: ["packet", "state"]
```

---

## Selective Auditing

**Default: ALL stages have audit enabled.** This provides full visibility for debugging complex rules. Use `enabled: false` only when you explicitly want to hide a stage.

### Audit Levels
| Level | Use For | Example |
|-------|---------|--------|
| `debug` | Internal housekeeping, calculations | Buffer updates, metric calculations |
| `info` | Normal flow steps | State transitions |
| `warning` | Important decision points | Condition matches, threshold crossed |
| `critical` | Alerts, final actions | Jammer detected, alert sent |

### Example: Disable Specific Stage
```yaml
stages:
  - rule: InternalCalculation
    order: 3
    audit:
      enabled: false  # Explicitly hide this stage from audit
```

---

## GRULE Listener Implementation

### Listener Interface

```go
// audit/listener.go
package audit

import (
    "log"
    "github.com/hyperjumptech/grule-rule-engine/ast"
)

type AuditListener struct {
    manifest   *AuditManifest
    loggedOnce map[string]bool
    loggedMu   sync.Mutex
}

func NewAuditListener(manifest *AuditManifest) *AuditListener {
    return &AuditListener{
        manifest:   manifest,
        loggedOnce: make(map[string]bool),
    }
}

// Called BEFORE rule consequence executes
func (l *AuditListener) BeforeExecuteConsequence(
    cycle uint64,
    rule *ast.RuleEntry,
    dc ast.IDataContext,
) {
    // Use EXISTING global toggle - ensures frontend controls work
    // (Activar/Desactivar buttons in Progress Audit dashboard)
    if !IsProgressAuditEnabled() {
        return
    }
    
    ruleName := rule.RuleName
    salience := rule.Salience
    
    // Lookup enrichment from manifest
    meta := l.manifest.GetRuleMeta(ruleName)
    
    // Skip if rule not in manifest or explicitly disabled
    if meta == nil {
        // Unknown rule - log once and skip (no audit)
        log.Printf("[AuditListener] Rule '%s' not in manifest, skipping audit", ruleName)
        return
    }
    
    if !meta.Enabled {
        return  // Explicitly disabled
    }
    
    // Extract snapshot with nil-safety
    snapshot, err := extractSnapshot(dc)
    if err != nil {
        log.Printf("[AuditListener] Error extracting snapshot for '%s': %v", ruleName, err)
        snapshot = map[string]interface{}{"error": err.Error()}
    }
    
    // Capture audit entry (async to avoid blocking hot path)
    entry := &AuditEntry{
        RuleName:    ruleName,
        Salience:    salience,
        Description: meta.Description,
        Level:       meta.Level,
        IsAlert:     meta.IsAlert,
        Snapshot:    snapshot,
    }
    
    // Non-blocking capture
    go Capture(entry)
}

// Called AFTER rule consequence executes
func (l *AuditListener) AfterExecuteConsequence(
    cycle uint64,
    rule *ast.RuleEntry,
    dc ast.IDataContext,
) {
    // Optional: capture execution result/errors
}
```

---

## Snapshot Extraction

```go
// extractSnapshot safely extracts data from DataContext
func extractSnapshot(dc ast.IDataContext) (map[string]interface{}, error) {
    snapshot := make(map[string]interface{})
    
    // Get IncomingPacket
    if packet, err := dc.Get("IncomingPacket"); err == nil && packet != nil {
        snapshot["packet"] = safeExtract(packet)
    }
    
    // Get state wrapper
    if state, err := dc.Get("state"); err == nil && state != nil {
        snapshot["state"] = safeExtract(state)
    }
    
    return snapshot, nil
}

// safeExtract handles nil values and panics
func safeExtract(v interface{}) interface{} {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[AuditListener] Panic extracting snapshot: %v", r)
        }
    }()
    
    if v == nil {
        return nil
    }
    
    // Type-specific extraction
    switch t := v.(type) {
    case *IncomingPacket:
        return map[string]interface{}{
            "IMEI": t.IMEI,
            "Speed": t.Speed,
            // ... other fields
        }
    default:
        return fmt.Sprintf("%+v", v)
    }
}
```

---

## Example: Clean GRL

```grl
rule DEFCON1_ContactLost_Pass salience 900 {
    when
        IncomingPacket.PositioningStatus == "V"
    then
        IncomingPacket.PositionInvalidDetected = true;
        actions.Log("[DEFCON 1] Positive: Status V detected.");
        // ✅ No Audit() call - handled by listener
}
```

## Example: Manifest File

**File:** `frontend/rules_templates/jammer_wargames.yaml`

```yaml
name: jammer_wargames
version: "1.0.0"

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
    
  - rule: DEFCON2_CalculateMetrics
    order: 3
    audit:
      enabled: true
      description: "Step 2: Metrics Calculated"
      level: debug
    
  - rule: DEFCON2_Inhibition_Pass
    order: 4
    audit:
      enabled: true
      description: "Step 2: Inhibition Confirmed"
      level: warning
    
  - rule: DEFCON3_SafeZones_Pass
    order: 5
    audit:
      enabled: true
      description: "Step 3: Outside Safe Zones"
      level: warning
    
  - rule: DEFCON4_JammerAlert_Fire
    order: 6
    audit:
      enabled: true
      description: "CRITICAL: JAMMER DETECTED"
      level: critical
      is_alert: true
```

---

## ✅ Success Criteria

- [ ] GRULE listener captures all rule executions automatically
- [ ] Audit manifest loads from database at startup
- [ ] Rule names matched to manifest for rich metadata
- [ ] GRL rules work WITHOUT `actions.Audit()` calls
- [ ] Hardcoded `Audit()` stub removed from backend code
- [ ] Database has `audit_manifest` column in `fleet_rules` table
- [ ] Frontend loads paired `.grl` + `.yaml` template files
- [ ] Frontend sends manifest to API when creating rules
- [ ] Integration test passes with auto-audit

---

## Key Constraints

1. **Clean Implementation** - No hardcoded `Audit()` calls anywhere in backend
2. **Frontend Controls Compatible** - Uses existing `IsProgressAuditEnabled()` toggle
3. **Manifest Required** - Rules without manifest entry logged once (no audit for unknown rules)
4. **No GRL Syntax Changes** - Standard GRULE syntax preserved
5. **Non-blocking** - Audit capture runs async (`go Capture()`) to avoid hot-path delays
6. **Nil-safe** - Snapshot extraction handles nil values and panics gracefully

---

## Frontend Controls Compatibility

The existing Progress Audit dashboard controls continue to work:

| Control | Action | Effect on Listener |
|---------|--------|-------------------|
| **Activar** (Enable) | Sets `progressAuditEnabled = true` | Listener starts capturing |
| **Desactivar** (Disable) | Sets `progressAuditEnabled = false` | Listener stops capturing |
| **Limpiar datos** (Clear) | Truncates `rule_execution_state` | Clears audit history |
| **Status badge** | Reads `IsProgressAuditEnabled()` | Shows ACTIVO/INACTIVO |

The listener uses `audit.IsProgressAuditEnabled()` instead of its own flag, so no frontend changes are required.

---

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Rule not in manifest | Log warning once, skip audit |
| `audit.enabled: false` | Silent skip |
| Snapshot extraction fails | Log error, capture with `{"error": "..."}` |
| YAML parse error | Startup failure with clear error message |
| DB write fails | Log error, continue (don't crash worker) |
| Nil in DataContext | Safe extraction returns nil, no panic |

---

## Performance Considerations

### Hot Path Optimization

1. **Async DB writes** - `go Capture()` so listener doesn't block rule execution
2. **O(1) manifest lookup** - `map[string]*RuleMeta` for instant access
3. **Global kill switch** - `listener.SetEnabled(false)` disables all auditing
4. **Level filtering** - Option to skip `debug` level in production
