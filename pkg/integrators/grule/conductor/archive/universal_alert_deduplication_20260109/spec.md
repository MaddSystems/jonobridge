# Specification: Universal Alert Deduplication

**Track ID:** `universal_alert_deduplication_20260109`
**Date:** January 9, 2026
**Status:** Approved
**Related Files:** `backend/grule/context_builder.go`, `frontend/rules_templates/jammer_wargames.grl`

## 1. Problem Definition

### The Issue
The Grule rule engine evaluates rules in multiple cycles per packet until no more rules fire. This "Rete algorithm" behavior causes **duplicate alerts** for the same IMEI within the same packet processing session.

**Symptoms:**
- A single incoming packet triggers the same rule multiple times (once per cycle).
- The Global Guard (`AlertsCapability`) prevents *inter-packet* duplicates (e.g., Packet A vs Packet B) but fails to prevent *intra-packet* duplicates because it doesn't block re-evaluation within the same transaction/context efficiently or logically for this specific engine behavior.

### Root Cause
1.  **Per-Packet State:** Each packet creates a new `StateWrapper`.
2.  **No Local Cache:** The `StateWrapper` lacks a memory of what happened in previous cycles of *this specific packet execution*.
3.  **Global Latency/Logic:** Relying solely on the global guard for intra-packet cycles is inefficient and conceptually wrong; the engine needs to know "I just fired this rule 1 millisecond ago in Cycle 1, don't fire it in Cycle 2".

## 2. Proposed Solution

Implement a **"Local-First, Global-Sync"** deduplication pattern.

### Architecture

1.  **Local State (The Cache):**
    The `StateWrapper` will hold a `map[string]bool` called `AlertStates`. This exists *only* for the lifetime of the packet processing.

2.  **Lazy Loading (The Read Path):**
    When a rule checks `IsAlertSentForRule("RuleName")`:
    -   **Step 1 (Local):** Check `AlertStates`. If `true`, return `true` (Blocked).
    -   **Step 2 (Global):** If not in local map, check the Global Guard (`AlertsCapability`).
    -   **Step 3 (Cache):** If Global Guard says `true`, update `AlertStates` so next time it's a fast local check.

3.  **Write-Through (The Write Path):**
    When a rule calls `MarkAlertSentForRule("RuleName")`:
    -   **Step 1 (Global):** Call Global Guard `MarkAlertSent`.
    -   **Step 2 (Local):** If successful, set `AlertStates["RuleName"] = true`. This ensures that if the engine runs another cycle for this packet, the Local Check (Step 1 of Read Path) will block it.

## 3. Technical Implementation

### Backend (`backend/grule/context_builder.go`)

**Struct Update:**
```go
type StateWrapper struct {
    // ... existing fields
    AlertStates map[string]bool // Ephemeral map for intra-packet deduplication
}
```

**Initialization:**
```go
// In ContextBuilder.Build():
stateWrapper.AlertStates = make(map[string]bool)
```

**Logic Methods:**
```go
func (s *StateWrapper) IsAlertSentForRule(ruleName string) bool {
    // 1. Local Check
    if sent, exists := s.AlertStates[ruleName]; exists && sent {
        log.Printf("🛡️ [StateWrapper] BLOCKED (local): %s", ruleName)
        return true
    }
    // 2. Global Check & Cache
    if s.Alrt != nil && s.Alrt.IsAlertSent(s.imei, ruleName) {
        s.AlertStates[ruleName] = true
        return true
    }
    return false
}

func (s *StateWrapper) MarkAlertSentForRule(ruleName string) bool {
    // 1. Global Mark
    if s.Alrt == nil || !s.Alrt.MarkAlertSent(s.imei, ruleName) {
        return false
    }
    // 2. Local Mark
    s.AlertStates[ruleName] = true
    return true
}
```

### Rules (`.grl`)

**Standardization:**
Rules will no longer use arbitrary ID strings (like `"jammer_real_mercury_2025"`). They will use their own **Rule Name** as the unique identifier.

```grl
rule DEFCON4_JammerAlert_Fire "DEFCON 4: JAMMER ALERT" {
    when
        !state.IsAlertSentForRule("DEFCON4_JammerAlert_Fire")
    then
        state.MarkAlertSentForRule("DEFCON4_JammerAlert_Fire");
        // ... actions
}
```

## 4. Benefits
1.  **Eliminates Duplicates:** Solves the root cause of multi-cycle re-evaluation.
2.  **Performance:** Reduces calls to the global guard (and potentially DB) by serving repeated checks from the map.
3.  **Clean Code:** Decouples specific rule logic from the context builder (no need to pre-load specific rules).
4.  **Extensible:** Any new rule just needs to use the standard method; no backend changes required.
