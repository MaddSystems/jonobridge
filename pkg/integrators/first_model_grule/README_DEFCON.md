# Jammer Detection System: Project "Wargames" (DEFCON 0-4)

This document outlines the multi-stage Jammer Detection logic, strictly adhering to the following escalation strategy:

## The Logic Map (DEFCON Levels)

The system operates on a 5-level escalation scale (0 to 4).

### **DEFCON 0: Surveillance (Buffering & Normal State)**
*   **Status:** Normal Monitoring.
*   **Logic:** Infinite loop. Every single packet (valid or invalid) enters here to update the memory buffer (last 10 positions).
*   **Condition:** `rule_execution_state->step_number = 0`.
*   **Exit to 1:** If `CurrentPacket == Invalid` **AND** `TimeSinceLastValid > 5 minutes`.
*   **Rule:** `DEFCON0_Surveillance`

### **DEFCON 1: Tracker Contact is Lost**
*   **Status:** Anomaly Detected.
*   **User Requirement:** "Activated when we receive V (invalid) for 5 minutes".
*   **Logic:** The device is confirmed "Offline" or sending invalid data for a sustained period (offline timer > 5 min).
*   **GPSgate Map:** Corresponds to expression *"5 min han pasado desde la ultima Posición valida"*.
*   **Exit to 2:** If `Buffer.Length == 10` (Full data set available for heuristic analysis).
*   **Rule:** `DEFCON1_ContactLost_Pass`

### **DEFCON 2: Inhibition Detected (High Speed + Signal -> No GPS)**
*   **Status:** High Probability of Jammer.
*   **Logic:** Signature Analysis. Analyze the buffer to determine *why* it went offline.
*   **GPSgate Map:** Corresponds to the Script logic.
    *   `AverageSpeed (90m window) >= 10 km/h` In short: It means "Average speed of the last 10 positions, provided they occurred within the last 90 minutes."
    *   `AverageSignal (last 5) >= 9`
*   **Meaning:** "The vehicle was moving fast and had good signal, then suddenly died." This is the signature of a Jammer.
*   **Exit to 3:** If thresholds are met.
*   **Rule:** `DEFCON2_Inhibition`

### **DEFCON 3: Outside Safe Zones**
*   **Status:** Strategic Clearance.
*   **Logic:** Ensure the vehicle is not in a known "safe zone" where signal loss is expected.
*   **GPSgate Map:** Corresponds to Expressions *"Taller: Fuera", "Clientes: Fuera", "Resguardo: Fuera"*.
*   **Exit to 4:** If **NOT** inside any of these zones.
*   **Rule:** `DEFCON3_SafeZones`

### **DEFCON 4: Jammer Detected**
*   **Status:** Action Stage (Alert Fired).
*   **Logic:** Fire the alert.
*   **GPSgate Map:** Corresponds to *"Notifications: Telegram"*.
*   **Action:** Send Telegram message and log the event.
*   **Rule:** `DEFCON4_JammerDetected`

---

## Technical Implementation

### Rule Engine (GRL)
The logic is implemented in `rules_templates/jammer_wargames.grl`.
*   **Salience 1000:** DEFCON 0 (Buffer Update)
*   **Salience 900:** DEFCON 1 (Contact Lost)
*   **Salience 800:** DEFCON 2 (Inhibition Analysis)
*   **Salience 700:** DEFCON 3 (Safe Zone Check)
*   **Salience 600:** DEFCON 4 (Alert)

### Auditing (Progress Audit)
The "Movie Frames" audit viewer (`progress_audit_movie.html`) and the Go worker (`grule_worker.go`) have been updated to visualize these specific stages using the exact terminology above.

### Testing
A test script `send_periodic4defcon4.py` is provided to simulate the full escalation chain:
1.  **Phase 1:** Sends 10 valid packets with High Speed (100 km/h) and Strong Signal to fill the buffer (DEFCON 0).
2.  **Phase 2:** Sends 25 invalid packets over ~6 minutes to trigger the "Offline > 5 min" condition (DEFCON 1) and escalate through DEFCON 2, 3, and 4.

**Note:** The test script sends `speed=100` even during the invalid phase to ensure the buffer's average speed remains high enough to pass the DEFCON 2 check.
