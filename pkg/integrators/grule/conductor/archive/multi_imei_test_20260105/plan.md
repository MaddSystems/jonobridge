# Plan - Multi-IMEI Defcon Testing

Create a new test script `tests/send_multiple.py` to simulate 4 IMEIs with different Defcon targets to verify concurrent per-IMEI state isolation.

## Phase 1: Setup and Scaffolding
- [x] Task: Create `tests/send_multiple.py` and copy utility functions (`crc`, `charcounter`, `identifier`, `payload`) from `send_periodic4defcon4.py`.
- [x] Task: Define the IMEI configuration data structure containing the 4 test IMEIs and their specific signal parameters (GSM strength, etc.).
- [x] Task: Conductor - User Manual Verification 'Phase 1: Setup and Scaffolding' (Protocol in workflow.md)

## Phase 2: Implementation of Sequential Multi-Sending
- [x] Task: Implement the Warmup Loop (Phase 1 of the simulation) to send 11 valid packets for all 4 IMEIs.
- [x] Task: Implement the Simulation Loop (Phase 2 of the simulation) to send 24 invalid packets with varying signal parameters per IMEI.
- [x] Task: Add console logging to track the progress of each IMEI independently during execution.
- [x] Task: Conductor - User Manual Verification 'Phase 2: Implementation of Sequential Multi-Sending' (Protocol in workflow.md)

## Phase 3: Verification (Manual)
- [x] Task: Run `tests/send_multiple.py` and verify console output shows interleaved packets for all 4 IMEIs.
- [x] Task: Manually verify in the Audit UI that IMEI_1/2 reach Defcon 4, while IMEI_3 and IMEI_4 hit their respective targets.
- [x] Task: Conductor - User Manual Verification 'Phase 3: Verification (Manual)' (Protocol in workflow.md)
