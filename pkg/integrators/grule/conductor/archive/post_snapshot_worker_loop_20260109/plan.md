# Implementation Plan: Post-Execution Snapshot Capture via Worker Loop

**Track ID:** `post_snapshot_worker_loop_20260109`
**Date:** January 09, 2026
**Priority:** Medium-High
**Scope:** `backend/`

## Objective
Capture audit snapshots after each rule's `then` block executes (post-state), without modifying any `.grl` files.

## Strategy
1. Load rules as ordered `RuleKB` slice (name + KB).
2. Execute rules sequentially in worker loop.
3. Capture snapshot immediately after each successful execution.
4. Use shared `DataContext` to get post-state naturally.
5. Optional: reduce/disable pre-captures in listener.

## Phase 1: Backend Implementation

### TODO 1.1: Add RuleKB Struct & Update Loading
**Status:** `[x]` Complete
**File:** `backend/main.go`
- Define `RuleKB` struct
- Modify `loadRulesFromSlice` to accept `manifest` and return `[]RuleKB` with sorting by manifest `Order`
- Update caller in `main()` to pass `ruleKBs` to worker

### TODO 1.2: Implement Ordered Execution Loop with Post-Capture
**Status:** `[x]` Complete
**File:** `backend/grule/worker.go`
- Change current rule execution to explicit loop over `w.ruleKBs`
- After `ExecuteWithContext`, call `audit.ExtractSnapshot` + `audit.Capture`
- Populate entry with manifest metadata
- **Important:** Ensure `audit.FinishCapture(imei)` is called at the end of the packet processing loop (if applicable/needed) to flush any pending buffers.

### TODO 1.3: Add IsPost Flag & Handle in Audit
**Status:** `[x]` Complete
**Files:** `backend/audit/types.go`, `backend/audit/capture.go`, `backend/audit/db.go`
- Add `IsPost bool` to `AuditEntry`
- Update `Capture` logic to save to DB (add `is_post` BOOLEAN column to `rule_execution_state`)
- Run migration or backfill existing rows (default to `false`)

### TODO 1.4: Optional Listener Cleanup
**Status:** `[ ]` Pending  
**File:** `backend/audit/listener.go`  
- Comment out or conditionalize pre-snapshot capture  
- Or keep for pre-metadata (tag as `IsPost: false`)

## Phase 2: Verification

### TODO 2.1: Functional Testing
**Status:** `[ ]` Pending  
- Send test packets (e.g. 35-frame sequence)  
- Verify DEFCON0 snapshots show `BufferUpdated: true`, `BufferHas10: true` (if applicable)  
- Check logs: updates visible in same-rule snapshot  
- Confirm ordering matches manifest (DEFCON0 → DEFCON1 etc.)

### TODO 2.2: DB & UI Validation
**Status:** `[ ]` Pending  
- Check `rule_execution_state` table: post entries present, `IsPost` flag correct  
- UI movie view: state evolution looks correct per rule

### TODO 2.3: Performance & Duplicates
**Status:** `[ ]` Pending  
- Confirm compatibility with Universal Alert Deduplication (no extra duplicate alerts generated)  
- Measure overhead (should be negligible)

## Phase 3: Deployment

### TODO 3.1: Final Integration
**Status:** `[ ]` Pending  
- Update `backend/main.go` to use the new `loadRulesFromSlice` signature.
- Deploy changes to `worker.go`.
- Test with real MQTT packets to ensure end-to-end stability.

## Dependencies
- Must be after deduplication implementation (fewer cycles = cleaner audits)
- Assumes current worker uses separate KBs per rule