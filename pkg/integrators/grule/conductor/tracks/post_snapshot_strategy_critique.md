# Post-Snapshot Strategy Critique & Validation Track

## 1. Context & Objective
The goal is to capture a **Post-Execution Snapshot** of the rule engine state (specifically `IncomingPacket.BufferUpdated`) *immediately after* a specific rule (`DEFCON0_Surveillance`) executes its actions. This is to verify that the rule correctly updated the state (setting `BufferUpdated = true`).

## 2. Current Architecture & Failure Analysis

### 2.1. Listener (`ExecuteRuleEntry`)
- **Mechanism:** The `AuditListener` implements `ExecuteRuleEntry`.
- **Timing:** This hook fires **BEFORE** the rule's `then` block is executed.
- **Result:** The snapshot captures the state *before* modification. Thus, `BufferUpdated` is always `false` (because the rule only fires if it's false).
- **Verdict:** This hook is suitable for "Pre-Execution" audit but incapable of capturing the result of the rule's logic.

### 2.2. Worker Loop Post-Capture
- **Mechanism:** The `Worker.Process` loop iterates through `RuleKB` entries and calls `audit.Capture` after `eng.ExecuteWithContext` returns.
- **Timing:** This fires after the **entire KnowledgeBase** (containing multiple rules) has finished execution.
- **Configuration Issue:** The `RuleKB` structure maps to a *file/group* name (e.g., "Jammer-Wargames-DEFCON"), not individual rule names.
- **Result:**
    1. It is too coarse-grained (post-KB, not post-Rule).
    2. It fails to find metadata because it looks up the Group Name in the Manifest (which keys by Rule Name), resulting in "No manifest meta".
- **Verdict:** This approach is fundamentally flawed for rule-level granularity and is currently misconfigured.

### 2.3. Recent Attempt (Packet Resolution)
- **Change:** Added explicit packet resolution (`l.packet`) to `AuditListener` to fix data visibility issues (pointer vs value).
- **Outcome:** This fixed the ability to *see* the packet data (avoiding `GoValueNode` logs), but did not change the *timing* of the capture. It confirmed that we are capturing the correct object, but at the wrong time.

## 3. Strategy Critique

The current strategy of relying on `ExecuteRuleEntry` for validation is impossible because of the temporal order of events. The strategy of using the Worker loop for "Post" capture is structurally misaligned with the requirement of per-rule validation.

**Missing Component:** A true "Post-Rule-Execution" hook.

## 4. Proposed Plan (The Track)

### 4.1. Validate `GruleEngineListener` Capabilities
We need to determine if `github.com/hyperjumptech/grule-rule-engine` (v1.20.4) supports a post-execution listener method.
- **Action:** Attempt to implement `DidExecuteRule` (or similarly named methods like `ExecuteRuleExit`) in `AuditListener`.
- **Validation:** Compile the code. If it compiles, the interface supports it.

### 4.2. Implementation Path A: Interface Exists
If `DidExecuteRule` exists:
1.  Implement `DidExecuteRule` in `backend/audit/listener.go`.
2.  Move the "Post-Execution" capture logic from `Worker.go` to this new method.
3.  Ensure `IsPost: true` is set.
4.  Remove the dead code in `Worker.go`.

### 4.3. Implementation Path B: Interface Missing (Contingency)
If the interface does not support a post-hook:
- **Option 1 (Intrusive):** Inject a `CapturePost()` call into the `then` block of the GRL rules. This violates the "declarative/invisible" audit goal but is guaranteed to work.
- **Option 2 (Engine Fork/Mod):** Too expensive for this task.
- **Option 3 (Single-Rule Execution):** Refactor `Worker` to load and execute rules individually. This would break the Rete network's ability to chain inferences (e.g., DEFCON0 -> DEFCON1 in one cycle) and is likely non-viable for the product's logic.

**Recommendation:** Proceed with testing **Path A**. If it fails, fallback to **Option 1 (Intrusive)** as a temporary measure or re-evaluate the requirement for "Post" snapshots (perhaps "Pre" of the *next* rule is sufficient?).

## 5. Next Steps
1.  Modify `backend/audit/listener.go` to add `DidExecuteRule` (signature guess: `DidExecuteRule(ctx context.Context, cycle uint64, entry *ast.RuleEntry, err error)`).
2.  Attempt compilation.
3.  If successful, refine implementation and deploy.
4.  If compilation fails, report and switch to "Intrusive GRL" strategy.
