# Specification: Cleanup Unused Cond5 Flags

## Context
The `Property` struct and `PacketWrapper` in the Go backend contain boolean flags and channels for `Cond5` (`Cond5Processed`, `Cond5Passed`, `Cond5Failed`). These appear to be vestiges of an older rule set or reserved fields that are not currently used by the active rule files (specifically `jammer_wargames.grl`). To improve code cleanliness and reduce memory footprint, these unused fields should be removed.

## Goals
1.  **Code Cleanup:** Remove unused `Cond5` boolean fields from `Property` and `PacketWrapper` structs.
2.  **Resource Optimization:** Remove initialization of associated channels in `Property`.
3.  **Verification:** Ensure no active rules rely on these fields before removal.

## Scope

### In Scope
-   **Go Code:**
    -   `engine/property.go`: Remove `Cond5Processed`, `Cond5Passed`, `Cond5Failed` fields, methods, and channel initializations.
    -   `engine/grule_worker.go`: Remove `Cond5Processed`, `Cond5Passed`, `Cond5Failed` from `PacketWrapper` struct and sync logic.
-   **Verification:**
    -   Check `frontend/rules_templates/jammer_wargames.grl` for any references to `Cond5`.
    -   Check other GRL files just in case.

### Out of Scope
-   Any other conditions (`PositionInvalidDetected`, `MovingWithWeakSignal`, etc.) - these must remain touched.

## Implementation Details

### 1. Verify Usage
Search for `Cond5` in all `.grl` files. If found, stop and report.

### 2. Property Struct Update (`engine/property.go`)
- Remove fields: `Cond5Processed`, `Cond5Failed`, `Cond5Passed`.
- Remove channel keys: `"Cond5Processed"`, `"Cond5Failed"`, `"Cond5Passed"`.
- Remove Get/Set methods for these fields.

### 3. PacketWrapper Struct Update (`engine/grule_worker.go`)
- Remove fields: `Cond5Processed`, `Cond5Passed`, `Cond5Failed`.
- Remove sync logic in `executeRulesForPacket`.
- Remove initialization in `workerRoutine`.

## Verification Plan
1.  **Compilation:** `go build ./...` must pass.
2.  **Grep Check:** Ensure `Cond5` strings are gone from the Go codebase (except maybe comments if any).
