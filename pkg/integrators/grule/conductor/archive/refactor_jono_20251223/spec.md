# Specification: Refactor 'Jono' to 'IncomingPacket'

## Context
The term "Jono" is currently used as a key in the Grule DataContext and within rule files (.grl) to represent the packet data being processed. This terminology is obscure and does not clearly convey the purpose of the object. The decision has been made to refactor this term to "IncomingPacket" to better reflect that it represents the data packet that has just arrived for processing.

## Goals
1.  **Improve Clarity:** Replace the ambiguous term "Jono" with the descriptive "IncomingPacket".
2.  **Standardize Terminology:** Ensure consistent usage of "IncomingPacket" across the Rule Engine interface and rule definitions.
3.  **Maintain Functionality:** Ensure that the system continues to process packets correctly after the rename.

## Scope
The refactoring affects the `pkg/integrators/grule` service.

### In Scope
-   **Rule Engine Interface:** Updating `engine/grule_worker.go` to inject the packet wrapper into the `DataContext` using the key "IncomingPacket".
-   **Rule Definitions (.grl):** Updating all `.grl` files in `frontend/rules_templates/` and `frontend/static/` to reference `IncomingPacket` instead of `Jono`.
-   **Internal Go Code:** Renaming local variables and functions in `engine/grule_worker.go`, `engine/rule_loader.go`, and `main.go` to align with the "IncomingPacket" terminology.
-   **Logging:** Updating log messages to use "IncomingPacket" or "Packet" instead of "Jono".

### Out of Scope
-   **Shared Library:** The `models.JonoModel` struct definition resides in `github.com/MaddSystems/jonobridge/common`, which is outside the current workspace. It will **not** be renamed.
-   **External APIs:** Any external API contracts that are not part of the internal rule processing flow.

## Technical Details

### 1. DataContext Injection
**Current:**
```go
dataCtx.Add("Jono", &wrapper)
```
**New:**
```go
dataCtx.Add("IncomingPacket", &wrapper)
```

### 2. Rule Files (.grl)
All rules must be updated to use the new root object name.
**Example Change:**
*Before:* `when Jono.Speed > 0`
*After:* `when IncomingPacket.Speed > 0`

### 3. Variable Naming
-   `jono *models.JonoModel` -> `incomingPacket *models.JonoModel` (or just `packet`)
-   `ProcessJonoMessage` -> `ProcessPacketMessage`

## Verification Plan
1.  **Compilation:** The Go code must compile without errors (`go build ./...`).
2.  **Rule Validation:** The existing tests (e.g., `tests/send_periodic4defcon4.py`) should be run manually (Phase 4) to confirm that rules are still firing correctly with the new terminology.
