# Implementation Plan: Fix Post-Execution Snapshot Packet State

**Track ID:** `fix_post_snapshot_packet_state_20260115`
**Date:** January 15, 2026
**Priority:** High
**Scope:** `backend/`

## Objective
Fix the issue where post-execution snapshots display the original, unmodified packet state (e.g., `BufferUpdated: false`) instead of the state updated by rule execution.

## Root Cause Analysis
**What Should Happen:**
In `worker.go`, the rule execution updates the `IncomingPacket` within the `DataContext` (e.g., setting `BufferUpdated = true`). The post-execution snapshot should reflect these changes.

**What Is Happening:**
Currently, `worker.go` passes the original `packet` object as an override to `audit.ExtractSnapshot(..., packet)`.
Inside `ExtractSnapshot`, this override takes precedence:
```go
if packetOverride != nil {
    packetObj = packetOverride // Uses original, unmodified packet
} else {
    packetObj = dc.Get("IncomingPacket") // Updates ignored
}
```
This causes the snapshot to capture the pre-execution state even after the rule has executed.
The same issue applies to `listener.go` where `ExtractSnapshot` is called with `l.packet`, preventing the capture of any updates made to the DataContext during the flow.

## Strategy
**Corrective Action:** Stop passing the original packet as an override during post-execution capture (in `worker.go`) AND in the listener (in `listener.go`). Pass `nil` instead, forcing `ExtractSnapshot` to retrieve the modified `IncomingPacket` from the `DataContext`.

## Phase 1: Implementation

### TODO 1.1: Update Worker Capture Logic
**Status:** `[x]` Complete
**File:** `backend/grule/worker.go`
- Change `audit.ExtractSnapshot(dataContext, packet.IMEI, packet)` to `audit.ExtractSnapshot(dataContext, packet.IMEI, nil)`.
- Add a debug log: `log.Printf("[Worker] Capturing post-snapshot for '%s' (using dc state, no override)", rkb.RuleName)`.

### TODO 1.2: Enhance Debug Logging
**Status:** `[x]` Complete
**File:** `backend/audit/snapshot.go`
- Add log to confirm which packet source is being used (`packetOverride` vs `DataContext`).
- Add specific logging for `BufferUpdated` and `BufferHas10` flags immediately after extraction, including the IMEI for tracing.
  ```go
  if current, ok := snapshot["packet_current"].(map[string]interface{}); ok {
      log.Printf("[Snapshot] Extracted packet_current for IMEI %s: BufferUpdated=%v, BufferHas10=%v", imei, current["BufferUpdated"], current["BufferHas10"])
  }
  ```

### TODO 1.3: Confirm IsPost Persistence
**Status:** `[x]` Complete
**File:** `backend/audit/capture.go` / `backend/audit/db.go`
- Verify that `SaveProgressAudit` includes the `IsPost` field in the `INSERT` query.
- No DB migration is needed, just ensure the table structure supports it.

### TODO 1.4: Fix Listener Snapshot Override
**Status:** `[ ]` Pending
**File:** `backend/audit/listener.go`
- Change `ExtractSnapshot(dc, imei, l.packet)` to `ExtractSnapshot(dc, imei, nil)`.
- Add debug log: `log.Printf("[Listener] Preparing snapshot for '%s' (using DataContext state, no override)", rule.RuleName)`.

## Phase 2: Verification

### TODO 2.1: Restart and Test
**Status:** `[ ]` Pending
- Restart the backend to apply changes.
- Send a sequence of 35 test packets to simulate a full flow and fill the buffer.

### TODO 2.2: Verify Data
**Status:** `[ ]` Pending
- **Database Check:** Query `rule_execution_state` to verify post-execution snapshots.
    ```sql
    SELECT step_number, stage_reached, is_post, context_snapshot 
    FROM rule_execution_state 
    WHERE rule_name = 'DEFCON0_Surveillance' 
    ORDER BY execution_time DESC LIMIT 5;
    ```
- **Expectation:** Post-execution entries (`is_post=1`) must show `BufferUpdated: true` and `BufferHas10: true` (once buffer is full).

### TODO 2.3: Verify Logs
**Status:** `[ ]` Pending
- Check logs for the message: `[Snapshot] Using IncomingPacket from DataContext`.
- Verify listener logs show: `[Listener] Preparing snapshot... (using DataContext state)`.
- Verify the enhanced logs show `BufferUpdated=true` for the correct IMEI.

## Notes
- **Non-Firing Rules:** The current implementation captures a snapshot even if the rule logic didn't "fire" (execute consequences), as long as the rule was processed. This is acceptable as it preserves the state. If strict "fired" only capturing is needed later, we can compare pre/post snapshots, but for now, we capture always to show state evolution.