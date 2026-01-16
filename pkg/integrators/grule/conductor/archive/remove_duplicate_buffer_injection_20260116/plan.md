# Implementation Plan: Remove Duplicate Buffer Data from Snapshot

**Track ID:** `remove_duplicate_buffer_injection_20260116`
**Date:** January 16, 2026
**Priority:** Low
**Status:** Proposed
**Scope:** `backend/audit/snapshot.go`

## Goal
Eliminate data duplication in the audit snapshot by removing the manual injection of buffer data into `packet_current`. The buffer data is already correctly exposed at the top level of the snapshot as `buffer_circular` via the `SnapshotProvider` interface.

## Current State
Snapshots contain buffer data twice:
1.  **`buffer_circular`** (Root Level): Correct, automatic from `SnapshotProvider`.
2.  **`buffer`** (Nested in `packet_current`): Redundant, from manual extraction code.

## Strategy
Remove the manual extraction logic in `ExtractSnapshot` that calls `getBufferData` and assigns it to `extracted["buffer"]`.

## Implementation Steps

### 1. Remove Manual Injection Code
**File:** `backend/audit/snapshot.go`

Locate and remove this specific block within `ExtractSnapshot`:

```go
			// ADD BUFFER CONTENTS TO PACKET_CURRENT
			if bufferData := getBufferData(dc, imei); bufferData != nil {
				extracted["buffer"] = bufferData
			}
```

### 2. Remove Helper Function (Cleanup)
**File:** `backend/audit/snapshot.go`

If `getBufferData` is no longer used after step 1, remove the `getBufferData` function entirely to keep the code clean.

## Verification
*   **Restart Backend:** Ensure changes take effect.
*   **Send Test Packet:** Trigger a rule execution.
*   **Check Audit Log/DB:** Verify the snapshot JSON.
    *   ✅ `buffer_circular` should exist at the root.
    *   ❌ `packet_current.buffer` should NOT exist.
