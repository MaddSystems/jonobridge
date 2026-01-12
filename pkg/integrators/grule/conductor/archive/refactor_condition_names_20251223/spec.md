# Specification: Refactor Condition Flags to Descriptive Names

## Context
The current implementation uses abstract, sequential naming for condition flags (`Cond1Passed`, `Cond2Passed`, `Cond3Passed`, `Cond4Passed`) in the `PacketWrapper` and `Property` structs. These names are opaque and require constant cross-referencing with rule files to understand their meaning (e.g., that "Cond1" means "GPS Status V"). This reduces code readability and maintainability.

## Goals
1.  **Improve Readability:** Replace abstract names with self-documenting names that describe the business logic.
2.  **Enhance Maintainability:** Make the Go code and GRL rules easier to understand without needing deep context.
3.  **Consistency:** Ensure the new names are used consistently across the Go backend (structs, logs) and Rule files (.grl).

## Naming Mapping

| Current Name | New Name | Description |
| :--- | :--- | :--- |
| `Cond1Passed` | `PositionInvalidDetected` | Indicates GPS status is "V" (Void/Invalid). |
| `Cond2Passed` | `MovingWithWeakSignal` | Indicates the pattern of moving vehicle + low signal strength. |
| `Cond3Passed` | `OutsideAllSafeZones` | Indicates vehicle is outside known safe zones (workshop, clients, etc.). |
| `Cond4Passed` | `JammerPatternFullyConfirmed` | High-level confirmation of the Jammer detection chain. |

*Note: Also applies to corresponding `...Failed` and `...Processed` flags if they exist and are part of the same logical grouping.*

## Scope

### In Scope
-   **Go Code:**
    -   `engine/property.go`: Update `Property` struct fields and getter/setter methods.
    -   `engine/grule_worker.go`: Update `PacketWrapper` struct fields, initialization, and sync logic.
-   **Rule Files (.grl):**
    -   `frontend/rules_templates/*.grl`
    -   `frontend/static/*.grl`
-   **Logging:** Update log messages to reflect new names.

### Out of Scope
-   External API contracts (unless these internal flags are exposed directly, which they shouldn't be for external consumption).

## Implementation Details

### 1. Property Struct Update (`engine/property.go`)
Rename fields and accessors.
- `GetCond1Passed()` -> `GetPositionInvalidDetected()`
- `SetCond1Passed()` -> `SetPositionInvalidDetected()`
- ... and so on for all 4 conditions.

### 2. PacketWrapper Struct Update (`engine/grule_worker.go`)
Rename fields.
- `Cond1Passed` -> `PositionInvalidDetected`
- Sync logic in `executeRulesForPacket` must map Property getters to Wrapper fields correctly.

### 3. Rule Updates
Search and replace in GRL files:
- `IncomingPacket.Cond1Passed` -> `IncomingPacket.PositionInvalidDetected`
- `property.SetCond1Passed(...)` -> `property.SetPositionInvalidDetected(...)`

## Verification Plan
1.  **Compilation:** `go build ./...` must pass.
2.  **Logic Check:** Review changed GRL files to ensure the logic flow remains identical, just with different names.