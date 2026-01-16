# Implementation Plan: Explicit Per-Rule Post-Execution Snapshot Capture

**Track ID:** `explicit_per_rule_snapshot_20260115`
**Date:** January 15, 2026
**Priority:** High
**Status:** Proposed
**Scope:** `backend/` + `.grl` files (manual changes)

## Goal
Enable explicit, rule-controlled post-execution snapshots by adding `actions.CaptureSnapshot();` at the end of each `then` block in the GRL files.
This captures the state after the rule's actions (e.g., `BufferUpdated = true`, `BufferHas10 = true`), providing accurate per-rule "after" snapshots for the audit movie UI.

## Benefits
*   **True post-state per rule:** Shows rule-specific changes.
*   **Full control:** Rule author decides exactly when to capture.
*   **No backend injection or parsing:** Zero risk of syntax errors.
*   **No fake entries:** All snapshots from real rules.
*   **Easy to remove later:** Delete lines → zero overhead.
*   **Compatible with LLM-generated rules:** Instruct LLM to include it.

## Strategy
1.  Add `CaptureSnapshot()` method to `ActionsWrapper` in `context_builder.go`.
2.  Manually add `actions.CaptureSnapshot();` as the last line of each `then` block in relevant `.grl` files.
3.  Keep listener for pre-snapshots (optional: disable if confusing).
4.  Disable `POST_PACKET` final snapshot to avoid duplication.

## Implementation Steps

### 1. Add CaptureSnapshot to ActionsWrapper
*   **File:** `backend/grule/context_builder.go`

```go
func (a *ActionsWrapper) CaptureSnapshot() {
    dc := ast.NewDataContext() // ← Replace with actual dc access (add field to wrapper if needed)
    imei := a.imei
    ruleName := "unknown"      // Optional: improve later with context

    snapshot, err := audit.ExtractSnapshot(dc, imei, nil)
    if err != nil {
        log.Printf("[CaptureSnapshot] Error: %v", err)
        return
    }

    meta := manifest.GetRuleMeta(ruleName)
    entry := &audit.AuditEntry{
        IMEI:         imei,
        RuleName:     ruleName,
        Description:  meta.Description,
        Level:        meta.Level,
        IsAlert:      meta.IsAlert,
        StepNumber:   meta.Order,
        StageReached: meta.Description,
        Snapshot:     snapshot,
        IsPost:       true,
    }
    audit.Capture(entry)
    log.Printf("[CaptureSnapshot] Captured for %s", ruleName)
}
```

### 2. Modify GRL Files – Add the Call Manually
*   **File:** `frontend/rules_templates/jammer_wargames.grl` (and others)
*   **Action:** Add `actions.CaptureSnapshot();` as the last line of each `then` block.

**Example:**
```grl
rule DEFCON0_Surveillance "DEFCON 0: Surveillance" salience 1000 {
    when
        !IncomingPacket.BufferUpdated
    then
        // 1. Actualizar Buffer y Offline status
        // Nota: El reset de flags se recomienda hacer en Go antes de Execute()
        IncomingPacket.BufferHas10 = state.UpdateMemoryBuffer(IncomingPacket.Speed, IncomingPacket.GSMSignalStrength, IncomingPacket.Datetime, IncomingPacket.PositioningStatus, IncomingPacket.Latitude, IncomingPacket.Longitude);
        IncomingPacket.IsOfflineFor5Min = state.IsOfflineFor(5);
        
        // 2. Marcar como actualizado para activar las siguientes reglas
        IncomingPacket.BufferUpdated = true;
        actions.Log("[DEFCON 0] Surveillance Active. Buffer Updated.");
        actions.CaptureSnapshot("DEFCON0_Surveillance");
}
```

### 3. Disable POST_PACKET
*   **File:** `backend/grule/worker.go`
*   **Action:** Comment out or remove the final snapshot capture.

```go
// // Disabled to avoid duplication with per-rule snapshots
// finalSnapshot, err := audit.ExtractSnapshot(...)
```

### 4. Optional: Disable Listener Pre-Snapshots
*   **File:** `backend/audit/listener.go`
*   **Action:** If pre-snapshots (false values) are confusing, comment out.

```go
// Capture(&AuditEntry{ ... }) // ← Comment for pre-only
```

## Verification
*   User will verify the changes.
