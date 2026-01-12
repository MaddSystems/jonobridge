# Plan: Cleanup Unused Cond5 Flags

## Objective
Remove unused `Cond5` boolean flags (`Cond5Processed`, `Cond5Passed`, `Cond5Failed`) and their associated channels from the Go backend after verifying they are not used in active rule files.

## Phases

### [x] Phase 1: Verification
**Goal:** Confirm `Cond5` flags are not used in any GRL files.

- [x] **Task:** Search for `Cond5` in `frontend/rules_templates/` and `frontend/static/`.

### [x] Phase 2: Go Code Cleanup
**Goal:** Remove the unused fields and logic from the Go backend.

- [x] **Task:** Refactor `engine/property.go`
    - Remove `Cond5Processed`, `Cond5Failed`, `Cond5Passed` fields.
    - Remove initialization of `updatedChannels` for these keys.
    - Remove `Get`/`Set` methods for these fields.
- [x] **Task:** Refactor `engine/grule_worker.go`
    - Remove `Cond5Processed`, `Cond5Passed`, `Cond5Failed` from `PacketWrapper` struct.
    - Remove these fields from `PacketWrapper` initialization in `workerRoutine`.
    - Remove these fields from sync logic in `executeRulesForPacket`.
    - Remove any logging references to these fields.

### [x] Phase 3: Verification
**Goal:** Ensure the system builds and is clean.

- [x] **Task:** Run `go build ./...` to ensure no compilation errors.
- [x] **Task:** Verify no remaining references via `grep`.
