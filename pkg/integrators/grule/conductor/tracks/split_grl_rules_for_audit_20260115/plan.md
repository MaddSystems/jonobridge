# Implementation Plan: Split GRL Rules for Granular Audit

**Track ID:** `split_grl_rules_for_audit_20260115`
**Date:** January 15, 2026
**Priority:** High
**Scope:** `backend/`

## Objective
Enable granular post-execution auditing for *every* rule defined in a single GRL resource file. Currently, the backend loads an entire GRL file (containing multiple rules like DEFCON0, DEFCON1, etc.) into a single KnowledgeBase (KB). The worker loop executes this KB once, meaning the "post-execution" snapshot only captures the state after *all* rules in that file have executed (or after the engine cycle finishes).

To capture the state changes *between* rules (e.g., after DEFCON0 but before DEFCON1), we must execute them sequentially in separate cycles.

## Strategy
1.  **Parse GRL Content:** Instead of feeding the entire GRL string to the `RuleBuilder`, we will parse the GRL content to extract individual `rule Name { ... }` blocks.
2.  **Individual KnowledgeBases:** For each extracted rule, we will create a separate `RuleKB` instance containing only that single rule.
3.  **Order by Manifest:** We will use the existing `AuditManifest` to assign the correct execution order to each individual rule based on its name (e.g., `DEFCON0_Surveillance` -> Order 1).
4.  **Worker Execution Loop:** The existing `worker.go` loop already iterates through `ruleKBs` and captures a snapshot after each one. By splitting the rules, this loop will naturally capture a snapshot after *each rule* completes.

## Risk Assessment
-   **Breaking Change:** This modifies how rules are loaded. If the regex parsing fails or is too fragile, rules might not load.
-   **Dependencies:** If rules depend on each other within the same cycle (e.g., one rule sets a variable that another rule reads *in the same Rete network*), splitting them into separate KBs/Cycles *might* affect behavior if `Retract` or complex Rete state was relied upon. However, our rules primarily communicate via the `IncomingPacket` and `StateWrapper` in the `DataContext`, which persists across these execution cycles in `worker.go`. This "Chain of Responsibility" pattern is actually safer and clearer.
-   **Performance:** Slightly higher overhead due to creating multiple KBs and starting the engine multiple times. Given the low number of rules (5-10), this is negligible.

## Phase 1: Implementation

### TODO 1.1: Implement Rule Splitting in `loadRulesFromSlice`
**Status:** `[ ]` Pending
**File:** `backend/main.go`
-   Import `regexp`.
-   Define a regex to capture individual rules: `(?s)rule\s+([a-zA-Z0-9_]+)\s+.*?\s*\{.*?\}[\s\n]*` (Needs careful testing to handle comments/strings braces).
-   Iterate through matches.
-   Build a `RuleKB` for each match.
-   Lookup `Order` from manifest using the extracted rule name.
-   Append to the `kbs` slice.
-   Sort `kbs` slice by `Order`.

### TODO 1.2: Validate Regex Robustness
**Status:** `[ ]` Pending
-   Ensure the regex handles the standard GRL format used in `jammer_wargames.grl`.
-   Verify it captures the full body including nested braces if any (though standard GRL rules usually don't have nested braces logic that trips simple regex, we should be careful). *Refinement: A simple regex might be brittle with nested braces. A bracket-counting parser is safer, but for now, we assume standard indentation/formatting or simple structure.*

## Phase 2: Verification

### TODO 2.1: Restart and Test
**Status:** `[ ]` Pending
-   Restart backend.
-   Check logs to see "Loaded X individual rules".
-   Send test packets.

### TODO 2.2: Verify Snapshots
**Status:** `[ ]` Pending
-   Query DB for `is_post=1`.
-   We should now see multiple post-execution snapshots for the *same packet execution flow*, one for each rule (DEFCON0, DEFCON1, etc.), showing the progressive state updates.

## Rollback Plan
-   Revert changes to `backend/main.go` to restore the "bulk load" behavior.
