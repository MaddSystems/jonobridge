# Implementation Plan: Universal Alert Deduplication

**Track ID:** `universal_alert_deduplication_20260109`
**Date:** January 9, 2026
**Priority:** High
**Scope:** `backend/` and `frontend/`

## Objective
Implement a universal, rule-name-based deduplication mechanism to prevent duplicate alerts caused by the Grule engine's multi-cycle evaluation (intra-packet duplicates).

## Strategy
Use a "Local-First, Global-Sync" pattern. A temporary map in the `StateWrapper` tracks alert states for the duration of a single packet's processing.

1.  **Local State:** `StateWrapper` holds `AlertStates map[string]bool`.
2.  **Lazy Loading:** `IsAlertSentForRule` checks local map -> global guard -> caches result locally.
3.  **Write-Through:** `MarkAlertSentForRule` updates global guard -> updates local map.

## Phase 1: Backend Implementation

### TODO 1.1: Update StateWrapper
**Status:** `[x]` Complete
**File:** `backend/grule/context_builder.go`
Add the local state map to the `StateWrapper` struct and initialize it.

```go
type StateWrapper struct {
    // ... existing fields
    AlertStates map[string]bool // Universal Alert State (keyed by rule name)
}

// In Build():
stateWrapper.AlertStates = make(map[string]bool)
```

### TODO 1.2: Implement Universal Methods
**Status:** `[x]` Complete
**File:** `backend/grule/context_builder.go`
Implement the deduplication logic with lazy loading (avoiding the need for manifest loading).

```go
func (s *StateWrapper) IsAlertSentForRule(ruleName string) bool {
    // 1. Check local state (fast, intra-packet)
    if sent, exists := s.AlertStates[ruleName]; exists && sent {
        log.Printf("🛡️ [StateWrapper] BLOCKED (local): IMEI=%s, Rule=%s", s.imei, ruleName)
        return true
    }

    // 2. Check global guard (inter-packet)
    if s.Alrt != nil && s.Alrt.IsAlertSent(s.imei, ruleName) {
        // Cache true result locally for subsequent cycles
        s.AlertStates[ruleName] = true
        log.Printf("🛡️ [StateWrapper] BLOCKED (global): IMEI=%s, Rule=%s", s.imei, ruleName)
        return true
    }

    return false
}

func (s *StateWrapper) MarkAlertSentForRule(ruleName string) bool {
    if s.Alrt == nil { return false }

    // 1. Mark global
    success := s.Alrt.MarkAlertSent(s.imei, ruleName)

    // 2. Mark local (stop future cycles)
    if success {
        s.AlertStates[ruleName] = true
    }
    return success
}
```

## Phase 2: Rules Update

### TODO 2.1: Update GRL Templates
**Status:** `[x]` Complete
**File:** `frontend/rules_templates/jammer_wargames.grl`
Refactor alert rules to use the new methods and Rule Name as the ID.

**Before:**
```grl
!state.IsAlertSent("jammer_real_mercury_2025")
state.MarkAlertSent("jammer_real_mercury_2025")
```

**After:**
```grl
!state.IsAlertSentForRule("DEFCON4_JammerAlert_Fire")
state.MarkAlertSentForRule("DEFCON4_JammerAlert_Fire")
```

## Phase 3: Verification

### TODO 3.1: Test for Duplicates
**Status:** `[ ]` Pending User Verification
- Run `tests/send_multiple.py`
- Verify logs show "BLOCKED (local)" for re-evaluations.
- Verify exactly 1 alert per IMEI/Rule in `alert_details` table.
