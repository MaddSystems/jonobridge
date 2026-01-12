# Plan: Refactor 'Jono' to 'IncomingPacket'

## Objective
Refactor the term "Jono" to "IncomingPacket" across the `pkg/integrators/grule` service. This includes renaming the object key in the Grule DataContext, updating all Rule (.grl) files, and cleaning up Go variable names and log messages.
**Note:** `models.JonoModel` cannot be renamed as it resides in the shared `common` library which is outside the current workspace scope. Refactoring will focus on local usage and the Rule Engine interface.

## Phases

### [x] Phase 1: Rule Engine Interface Refactor
**Goal:** Change the key used to inject the packet data into the Grule engine from "Jono" to "IncomingPacket".

- [x] **Task:** Update `engine/grule_worker.go`
    - Change `dataCtx.Add("Jono", &wrapper)` to `dataCtx.Add("IncomingPacket", &wrapper)`.
    - Update logging to reflect this change.

### [x] Phase 2: Rule Files (.grl) Update
**Goal:** Update all logic files to use the new `IncomingPacket` key.

- [x] **Task:** Update `frontend/rules_templates/jammer_wargames.grl`
    - Replace all occurrences of `Jono.` with `IncomingPacket.`.
- [x] **Task:** Update `frontend/rules_templates/speed_alert_2.grl`
    - Replace all occurrences of `Jono.` with `IncomingPacket.`.
- [x] **Task:** Update `frontend/static/jammer_real.grl`
    - Replace all occurrences of `Jono.` with `IncomingPacket.`.

### [x] Phase 3: Go Code Variable & Function Renaming
**Goal:** Align internal Go code terminology with the new concept.

- [x] **Task:** Refactor `engine/grule_worker.go`
    - Rename variables `jono` (type `*models.JonoModel`) to `incomingPacket` or `packet`.
    - Rename `ProcessJonoMessage` to `ProcessPacketMessage` (or similar).
- [x] **Task:** Refactor `engine/rule_loader.go`
    - Update `SaveProcessedFactWithContext` parameters and logs.
- [x] **Task:** Refactor `main.go`
    - Update calls to `ProcessJonoMessage`.
- [x] **Task:** Refactor `actions` package if necessary (check logs/strings).

### [x] Phase 4: Verification (Manual)
**Goal:** Ensure the system builds and runs correctly. **To be performed by the human.**

- [x] **Task:** Run `go build ./...` to ensure no compilation errors.
- [x] **Task:** Run existing tests (e.g., `tests/send_periodic4defcon4.py` or unit tests if available) to verify rule execution.
