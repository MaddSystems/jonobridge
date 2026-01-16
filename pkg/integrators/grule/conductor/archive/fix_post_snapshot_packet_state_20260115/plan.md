# Implementation Plan: Fix Post-Execution Snapshot Packet State

**Track ID:** `fix_post_snapshot_packet_state_20260115`
**Date:** January 15, 2026
**Priority:** High
**Scope:** `backend/`

## Objective
Fix the issue where post-execution snapshots display the original, unmodified packet state (e.g., `BufferUpdated: false`) instead of the state updated by rule execution.

## Root Cause Analysis
1.  **Override Bug:** `worker.go` and `listener.go` were passing the original `packet` as an override to `ExtractSnapshot`, bypassing the updated `DataContext`. (Fixed in TODO 1.1 - 1.4).
2.  **Grule Wrapper Issue:** The object in `DataContext` is wrapped in `*model.GoValueNode`.
3.  **Marshaling Issue:** `json.Marshal` on the wrapper returns `{}` because the internal fields are unexported or not designed for serialization.

## Strategy
1.  **Unwrap Grule Node:** Use an interface assertion or reflection to get the underlying value from the Grule wrapper.
2.  **Bypass JSON Marshal:** Manually extract fields from the `*IncomingPacket` into a `map[string]interface{}` to ensure reliability and capture the latest state.
3.  **Resolve Dependencies:** To use `*IncomingPacket` inside the `audit` package (avoiding circular imports), we must move the struct definition to a shared package or use reflection-based mapping. (Preferred: Reflection mapping implemented).

## Phase 1: Implementation

### TODO 1.1: Update Worker Capture Logic
**Status:** `[x]` Complete
**File:** `backend/grule/worker.go`
- Pass `nil` to `ExtractSnapshot` for post-captures.

### TODO 1.2: Enhance Debug Logging
**Status:** `[x]` Complete
**File:** `backend/audit/snapshot.go`

### TODO 1.3: Confirm IsPost Persistence
**Status:** `[x]` Complete
**File:** `backend/audit/capture.go` / `backend/audit/db.go`

### TODO 1.4: Fix Listener Snapshot Override
**Status:** `[x]` Complete
**File:** `backend/audit/listener.go`

### TODO 1.5: Deep Debugging
**Status:** `[x]` Complete
**File:** `backend/audit/snapshot.go`
- Added logs for type and raw JSON. Confirmed `*model.GoValueNode` is the wrapper.

### TODO 1.6: Implement Unwrapping and Manual Extraction (New)
**Status:** `[x]` Complete
**File:** `backend/audit/snapshot.go`
- **Solution: Bypass JSON Marshal for Packet Extraction**
    - The code now unwraps the Grule node (via `GetValue()`).
    - It then uses reflection to manually map the fields of `*IncomingPacket` into a `map[string]interface{}`.
    - This bypasses the empty JSON marshaling issue and avoids circular dependencies.
    - Fields extracted: `IMEI`, `Speed`, `BufferUpdated`, `BufferHas10`, etc.
    - Debug logs added to confirm successful unwrapping and extraction.

## Phase 2: Verification

### TODO 2.1: Restart and Test
**Status:** `[x]` Complete
- Restart backend, send packets.

### TODO 2.2: Verify Data
**Status:** `[x]` Complete
- Verify `is_post=1` entries show `BufferUpdated: true`.

### TODO 2.3: Verify Logs
**Status:** `[x]` Complete
- Confirm `[ExtractSnapshot] Unwrapped Grule node` and `[ExtractSnapshot] Manually extracted packet_current`.
- Logs confirm: `[ExtractSnapshot] Manually extracted packet_current: BufferUpdated=true, BufferHas10=true` for post-execution frames.