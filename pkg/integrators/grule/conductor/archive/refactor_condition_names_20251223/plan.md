# Plan: Refactor Condition Flags to Descriptive Names

## Objective
Refactor abstract condition flags (`Cond1Passed`, etc.) to descriptive names (`PositionInvalidDetected`, etc.) across the Go backend and GRL rule files to improve code readability and maintainability.

## Phases

### [x] Phase 1: Go Struct & Method Refactor
**Goal:** Update the underlying Go structs and methods to use the new descriptive names.

- [x] **Task:** Refactor `engine/property.go`
    - Rename `Cond1Passed` -> `PositionInvalidDetected` (field & methods)
    - Rename `Cond2Passed` -> `MovingWithWeakSignal` (field & methods)
    - Rename `Cond3Passed` -> `OutsideAllSafeZones` (field & methods)
    - Rename `Cond4Passed` -> `JammerPatternFullyConfirmed` (field & methods)
    - *Note: Check for `...Failed` and `...Processed` variants and rename similarly if applicable.*
- [x] **Task:** Refactor `engine/grule_worker.go`
    - Update `PacketWrapper` struct fields.
    - Update initialization logic in `workerRoutine`.
    - Update sync logic in `executeRulesForPacket` (mapping Property -> Wrapper).
    - Update logging statements.

### [x] Phase 2: Rule File Updates
**Goal:** Update all GRL files to use the new `IncomingPacket` property names.

- [x] **Task:** Update `frontend/rules_templates/jammer_wargames.grl`
    - Replace `IncomingPacket.Cond1Passed` with `IncomingPacket.PositionInvalidDetected`.
    - Replace `IncomingPacket.Cond2Passed` with `IncomingPacket.MovingWithWeakSignal`.
    - Replace `IncomingPacket.Cond3Passed` with `IncomingPacket.OutsideAllSafeZones`.
    - Replace `IncomingPacket.Cond4Passed` with `IncomingPacket.JammerPatternFullyConfirmed`.
- [x] **Task:** Update `frontend/static/jammer_real.grl`
    - Apply same replacements if applicable (check if this rule file uses these specific flags).
- [x] **Task:** Update `frontend/rules_templates/speed_alert_2.grl`
    - Check if this file uses any of these flags and update if so.

### [x] Phase 3: Verification
**Goal:** Ensure the system builds and runs correctly.

- [x] **Task:** Run `go build ./...` to ensure no compilation errors.
- [x] **Task:** Manual verification of rule logic flow by reading the updated GRL files.
