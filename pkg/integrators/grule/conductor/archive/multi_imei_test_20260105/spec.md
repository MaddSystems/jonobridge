# Specification - Multi-IMEI Defcon Testing (`send_multiple.py`)

## Overview
Create a new test script `tests/send_multiple.py` that simulates 4 distinct IMEIs reporting telemetry simultaneously. The goal is to verify that the Grule Engine correctly maintains per-IMEI state and handles concurrent message streams without cross-contamination.

## Functional Requirements
- **Multi-Device Simulation:** Simulate 4 unique IMEIs in a single script using a sequential loop with short sleeps to maintain approximate concurrency.
- **Specific Outcomes:**
    - **IMEI_1 & IMEI_2:** Reach **Defcon 4** (Invalid packets for > 5 minutes + Low GSM).
    - **IMEI_3:** Reach **Defcon 2** (Invalid packets, but manipulated signals to stay at Defcon 2).
    - **IMEI_4:** Reach **Defcon 3** (Invalid packets, but manipulated signals to stay at Defcon 3).
- **Timing:** 
    - Use a fixed 15-second interval for all devices.
    - Offset the start times of each IMEI slightly to avoid a single burst of 4 packets exactly at the same millisecond.
- **Phases:**
    - **Warmup:** All 4 IMEIs send 11 valid packets to fill their buffers and activate `BufferHas10`.
    - **Simulation:** All 4 IMEIs send invalid packets (`Status V`), with varying signals to hit target Defcon levels.

## Non-Functional Requirements
- **Code Reuse:** Reuse the `payload`, `crc`, and `identifier` logic from `tests/send_periodic4defcon4.py`.
- **Maintainability:** Use a clean data structure (e.g., a list of dictionaries) to define the specific behavior and state for each IMEI.

## Acceptance Criteria
- [ ] `tests/send_multiple.py` is created.
- [ ] The script executes without errors, sending packets for 4 different IMEIs.
- [ ] Logs show packets being sent for all 4 IMEIs in an interleaved fashion.
- [ ] (Manual) Verification via Audit UI shows the correct Defcon progression for each independent IMEI.

## Out of Scope
- True multi-threading/multiprocessing (Sequential loop is sufficient for 4 devices).
- Automated verification of the engine state (Manual verification via UI/DB is expected).
