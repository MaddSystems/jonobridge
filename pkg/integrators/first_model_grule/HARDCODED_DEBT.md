# Hardcoded Rule Debt Analysis

**Date:** December 17, 2025
**Status:** Critical Technical Debt Identified

## Overview
The Grule Engine codebase contains explicit, hardcoded dependencies on specific business rules (specifically "Jammer" detection). This contradicts the purpose of a generic rule engine, where logic should be defined in GRL files, not in the Go source code.

**Critical Finding:** The hardcoded rule name `"Jammer Real - Detección Avanzada con Buffer Circular"` **does not exist** in any of the rule templates provided in `external-web/rules_templates`. This means the engine is likely trying to audit or track a rule that is either obsolete or manually inserted into the database, leading to potential "phantom" logic or failures in the audit system.

## 1. Hardcoded Locations

### A. `engine/grule_worker.go`
This file contains the most significant debt. It manually attempts to construct audit trails ("movie frames") for a specific rule.

*   **Line 587:** Explicit rule name hardcoding.
    ```go
    ruleName = "Jammer Real - Detección Avanzada con Buffer Circular"
    ```
*   **Lines 421, 444:** Checks for a specific alert ID string.
    ```go
    state.IsAlertSent("jammer_real_mercury_2025")
    ```
*   **Lines 491-494:** Manually constructs a `jammer_metrics` map for the audit log, pulling specific keys from the state.
    ```go
    "jammer_metrics": map[string]interface{}{
        "avg_speed_90min": state.GetCounter("jammer_avg_speed_90min"),
        ...
    }
    ```

### B. `engine/persistent_state.go`
The state engine, which should be a generic key-value store, has specific struct fields and logic for Jammer detection.

*   **Struct Definition (Lines 29-31):**
    ```go
    JammerPositions     int64
    JammerAvgSpeed90min int64
    JammerAvgGsm5       int64
    ```
*   **Logic Interception (Lines 163, 179, 194):** Generic methods like `IsAlertSent` have `if` statements to hijack behavior for `"jammer_real_mercury_2025"`.
*   **Dedicated Functions:**
    *   `UpdateJammerHistoryExact` (Line 415) - Completely custom logic mirroring a specific GPSGate script.
    *   `CalculateJammerMetricsIfReady` (Line 565).

### C. `engine/audit/types.go` & `engine/audit/db.go`
*   Contains specific JSON field mappings and comments referencing "Jammer Alert".

## 2. Safety Analysis: Can this be removed?

### Risk: High (Functionality Breakage)
Removing this code **will break the current Jammer Detection logic** if the GRL rules rely on the Go code to perform the heavy lifting (calculating averages, managing history buffers).

*   **Why it's unsafe to simply delete:** The GRL rules likely call methods like `state.UpdateJammerHistoryExact(...)` or access properties `state.JammerAvgSpeed90min`. If you remove the Go backing code, the rules will fail to compile or execute.
*   **Why the current state is bad:** The engine is not generic. Adding a new complex rule (e.g., "Fuel Theft") would require modifying the Go source code to add `FuelTheftPositions` struct fields, which is non-scalable.

### Discrepancy Verification
You noted that `external-web/rules_templates` does not contain the rule name `"Jammer Real - Detección Avanzada con Buffer Circular"`.
*   **Verified:** The templates contain names like `Jammer_Initialization_Component1`, `DEFCON6_FireAlert`, etc.
*   **Implication:** The code at `grule_worker.go:587` is likely executing logic for a rule that doesn't match your actual loaded rules. The `getRuleIDByName` call likely returns 0 or an error, and the "Progress Audit" (Movie Mode) might be generating invalid data or not working at all for your actual rules.

## 3. Recommendations

1.  **Immediate Action (Cleanup):**
    *   The code in `grule_worker.go` that generates "movie frames" for `"Jammer Real - Detección Avanzada..."` is likely dead code or buggy for your current rule set. **It is safe to remove the logic that filters specifically for this string**, but you must ensure the generic audit mechanism can still capture the state.

2.  **Refactoring Strategy (Make it Generic):**
    *   **State:** Replace `JammerPositions int64` with a `map[string]int64` or `map[string]interface{}` in `PersistentState`.
    *   **Logic:** Move the logic from `UpdateJammerHistoryExact` into the GRL rule itself, or (if too complex) create a *generic* `HistoryBuffer` helper in Go that any rule can instantiate by name (e.g., `buffer.Add("jammer_history", value)`).
    *   **Audit:** The audit system should dump the entire `state.Variables` map instead of hardcoding `jammer_metrics`.

3.  **Conclusion:**
    *   Do **not** simply delete the `PersistentState` fields yet, as your current active rules (in the DB) likely depend on them.
    *   **Do** delete/refactor the `grule_worker.go` hardcoded rule name check, as it refers to a non-existent rule name.
