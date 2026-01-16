# Implementation Plan: Explain and Enhance Buffer Snapshot Visibility

**Track ID:** `explain_buffer_snapshot_20260116`
**Date:** January 16, 2026
**Priority:** Medium
**Status:** Proposed
**Scope:** `backend/capabilities/buffer/`

## Goal
Explain how the buffer content is currently captured in audit snapshots and implement a modification to ensure **full visibility** of the buffer state, regardless of business logic filters (like the 90-minute window).

## Current Mechanism
1.  **Interface:** `BufferCapability` implements `SnapshotProvider`.
2.  **Method:** `GetSnapshotData(imei)` is called by the Audit system.
3.  **Logic:** Currently calls `GetEntriesInTimeWindow90Min(imei)`.
    *   **Limitation:** This filters entries based on `time.Now()`. If test data is older than 90 minutes, the snapshot returns an empty list `[]`, hiding the actual buffer state.

## Strategy
To show the *true* contents of the buffer in the snapshot (useful for debugging and audit):
1.  **Add Method:** Implement `GetAllEntries()` in `FixedCircularBuffer` to return the raw 10 entries without time filtering.
2.  **Update Capability:** Modify `GetSnapshotData` in `BufferCapability` to use `GetAllEntries()` instead of the filtered getter.

## Implementation Steps

### 1. Add `GetAllEntries` to Circular Buffer
**File:** `backend/capabilities/buffer/circular.go`

```go
// GetAllEntries returns all entries currently in the buffer without filtering
func (b *FixedCircularBuffer) GetAllEntries() []BufferEntry {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    
    // Return copy of slice up to current size
    result := make([]BufferEntry, b.size)
    copy(result, b.entries[:b.size])
    return result
}
```

### 2. Update `GetSnapshotData`
**File:** `backend/capabilities/buffer/capability.go`

Change the call from `GetEntriesInTimeWindow90Min` to `GetAllEntries` (or exposed via Manager).

*Note: You might need to add `GetAllEntries(imei)` to `manager.go` as well to bridge the call.*

**File:** `backend/capabilities/buffer/manager.go` (Likely needs update)
```go
func (m *BufferManager) GetAllEntries(imei string) []BufferEntry {
    b := m.GetOrCreateBuffer(imei)
    return b.GetAllEntries()
}
```

**File:** `backend/capabilities/buffer/capability.go`
```go
// GetSnapshotData implements SnapshotProvider
func (c *BufferCapability) GetSnapshotData(imei string) map[string]interface{} {
    if c == nil { return nil }

    // CHANGED: Use GetAllEntries to see raw state
    entries := c.manager.GetAllEntries(imei) 
    
    // ... mapping code remains the same ...
}
```

## Verification
*   **Restart Backend.**
*   **Send Test Data:** Use data with old timestamps (>90 mins ago).
*   **Check Audit:** Verify that `buffer_circular` in the snapshot now shows the entries instead of an empty list.
