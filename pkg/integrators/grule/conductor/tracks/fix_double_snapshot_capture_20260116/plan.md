# Implementation Plan: Fix Double Snapshot Capture & Alert Spam

**Track ID:** `fix_double_snapshot_capture_20260116`
**Date:** January 16, 2026
**Priority:** High
**Status:** Proposed
**Scope:** `backend/`

## Goal
Eliminate duplicate snapshots and alert spam by enforcing "Pure Explicit Capture". The system currently captures up to 3 snapshots per rule execution (Listener, Worker, Manual), causing noise and database bloat. We will switch to **exclusively** using the manual `actions.CaptureSnapshot()` calls inside GRL files.

## Problem Diagnosis

The current implementation violates the "explicit control" principle, resulting in 4 distinct bugs:

### BUG #1: DUPLICATE SNAPSHOTS
*   **Cause:** The Listener captures a pre-snapshot, the Worker captures a post-snapshot, and the GRL rule calls `actions.CaptureSnapshot()` manually.
*   **Result:** Triple redundancy for every rule execution.

### BUG #2: LISTENER CAPTURES EVERYTHING
*   **Cause:** The Listener captures on *every* rule evaluation attempt, even if the rule condition is false and the rule doesn't fire.
*   **Result:** Massive log noise with empty/irrelevant snapshots.

### BUG #3: REDUNDANT GRL CALLS
*   **Cause:** The manual `actions.CaptureSnapshot()` calls in the GRL files are executing in addition to the automatic backend captures.
*   **Result:** Intended as the *only* capture method, but currently acting as a third duplicate.

### BUG #4: ALERT SPAM (No Guard Check)
*   **Cause:** The Worker's post-capture logic (lines 95-115 in `worker.go`) blindly captures after every execution without checking if the rule actually fired or if an alert was already sent.
*   **Result:** Continuous alert snapshots for DEFCON4 rules on every packet.

## Strategy: Pure Explicit Capture (Option A)

We will remove all automatic backend capturing and rely solely on the explicit calls within the GRL rules. This gives full control to the rule author.

### Phase 1: Implementation

#### TODO 1.1: Disable Listener in Worker
**Status:** `[ ]` Pending
**File:** `backend/grule/worker.go`
*   Comment out or remove the Listener initialization and attachment (approx lines 70-75).
*   Ensure `eng.Listeners` is not set or is empty.

#### TODO 1.2: Disable Worker Post-Capture
**Status:** `[ ]` Pending
**File:** `backend/grule/worker.go`
*   Comment out or remove the automatic post-execution capture block inside the worker loop (approx lines 95-115).
*   The worker should strictly execute the rule and move on.

#### TODO 1.3: Verify GRL Explicit Calls
**Status:** `[ ]` Pending
**File:** `frontend/rules_templates/jammer_wargames.grl` (and others)
*   Verify that `actions.CaptureSnapshot("RuleName")` exists at the end of the `then` block for all relevant rules.
*   (This is likely already done, but verification is required).

## Phase 2: Verification

### TODO 2.1: Restart and Test
**Status:** `[ ]` Pending
*   Restart the backend.
*   Send a test packet sequence.

### TODO 2.2: Verify Single Capture
**Status:** `[ ]` Pending
*   Check `rule_execution_state` table.
*   Confirm exactly **ONE** entry per fired rule per packet.
*   Confirm NO entries for rules that evaluated to false.
