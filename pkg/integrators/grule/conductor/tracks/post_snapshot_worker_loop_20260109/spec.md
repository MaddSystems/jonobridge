# Specification: Post-Execution Snapshot Capture via Worker Loop

**Track ID:** `post_snapshot_worker_loop_20260109`
**Date:** January 09, 2026
**Status:** Proposed
**Related Files:** `backend/main.go`, `backend/grule/worker.go`, `backend/audit/snapshot.go`, `backend/audit/capture.go`, `backend/audit/listener.go`, `backend/audit/types.go`

## 1. Problem Definition

### The Issue
Current audit snapshots (captured in `AuditListener.ExecuteRuleEntry`) reflect the **pre-execution state** of each rule. This means flags and state modified inside the rule's `then` block (e.g. `BufferUpdated = true`, buffer entry added, `BufferHas10 = true`) are **not visible** in that rule's snapshot — even though logs clearly show the updates occurred.

**Symptoms:**
- **Example:** Snapshot shows `BufferUpdated: false` while the Log shows `"Buffer Updated"`.
- DEFCON0 snapshots always show `BufferUpdated: false`, `BufferHas10: false` despite the rule updating them.
- Human reviewers (and boss) see misleading/incomplete audit "movie" — state appears unchanged within a rule.
- Post-state only appears in the next rule's pre-snapshot (if any), breaking intuitive timeline.

### Root Cause
- Grule listener hook `ExecuteRuleEntry` fires **before** the consequence (`then` block) executes.
- No native post-execution hook exists in Grule.
- Snapshots are tied to the listener instead of post-rule execution point.

## 2. Proposed Solution

Capture snapshots **immediately after** each individual rule execution in the worker's processing loop. Since:
- Rules are loaded as separate `KnowledgeBase` instances
- They share the same `DataContext` (`dc`)
- Execution is sequential

… the `dc` reflects the **post-state** of each rule right after its `ExecuteWithContext` call.

### Key Principles
- **No GRL modification** — zero changes to rules
- **Universal** — applies automatically to all rules
- **Ordered execution** — respect manifest `Order` for correct DEFCON sequencing
- **Hybrid mode** — optional retention of listener for pre-metadata
- **Minimal changes** — leverage existing per-rule-KB design

## 3. Technical Implementation

### Backend (`backend/main.go`)

**New Struct:**
```go
type RuleKB struct {
    Name string
    KB   *ast.KnowledgeBase
}
```

**Updated Loading Function:**
```go
func loadRulesFromSlice(rules []persistence.Rule, manifest *audit.AuditManifest) []RuleKB {
	var ruleKBs []RuleKB
	for _, r := range rules {
		log.Printf("Loading rule: %s", r.Name)
		kb := ast.NewKnowledgeLibrary()
		rb := builder.NewRuleBuilder(kb)
		err := rb.BuildRuleFromResource("FleetRules", "0.0.1", pkg.NewBytesResource([]byte(r.GRL)))
		if err != nil {
			log.Printf("Error building rule %s: %v", r.Name, err)
			continue
		}
		kbInstance, err := kb.NewKnowledgeBaseInstance("FleetRules", "0.0.1")
		if err != nil {
			log.Printf("Error creating KB instance for %s: %v", r.Name, err)
			continue
		}
		ruleKBs = append(ruleKBs, RuleKB{Name: r.Name, KB: kbInstance})
	}

	// Sort by manifest order (critical for DEFCON sequencing)
	sort.Slice(ruleKBs, func(i, j int) bool {
		metaI := manifest.GetRuleMeta(ruleKBs[i].Name)
		metaJ := manifest.GetRuleMeta(ruleKBs[j].Name)
		orderI := 9999 // fallback
		orderJ := 9999
		if metaI != nil {
			orderI = metaI.Order
		}
		if metaJ != nil {
			orderJ = metaJ.Order
		}
		return orderI < orderJ
	})

	return ruleKBs
}
```

### Backend (`backend/grule/worker.go`)

**Execution Loop (main change):**
```go
for _, rkb := range w.ruleKBs {
    err := w.engine.ExecuteWithContext(ctx, dc, rkb.KB)
    if err != nil {
        log.Printf("[Worker] Execution failed for %s: %v", rkb.Name, err)
        continue
    }

    // Post-execution snapshot capture (this is the core improvement)
    imei := audit.GetIMEI(dc.Get("IncomingPacket"))
    snapshot, err := audit.ExtractSnapshot(dc, imei, ctx.Value("originalPacket"))
    if err != nil {
        log.Printf("[Snapshot] Failed for %s: %v", rkb.Name, err)
        continue
    }

    meta := w.manifest.GetRuleMeta(rkb.Name)
    if meta == nil {
        log.Printf("[Worker] Warning: No metadata for rule %s, skipping audit capture", rkb.Name)
        continue
    }

    entry := &audit.AuditEntry{
        IMEI:         imei,
        RuleName:     rkb.Name,
        Salience:     0, // Fetch from meta.Salience if available, or rule entry if needed
        Description:  meta.Description,
        Level:        meta.Level,
        IsAlert:      meta.IsAlert,
        StepNumber:   meta.Order,
        StageReached: meta.Description,
        Snapshot:     snapshot,
        IsPost:       true, // new flag
    }
    audit.Capture(entry)
}
```

### Audit Layer Adjustments

- **AuditEntry:** Add `IsPost bool` to `AuditEntry` struct in `backend/audit/types.go`.
- **Database:** Add `is_post` BOOLEAN column to `rule_execution_state` table.
- **Capture:** Update `Capture()` to persist the `is_post` flag.

## 4. Benefits

- True post-state in snapshots (shows rule effects)
- No GRL changes or injection risks
- Maintains rule ordering from manifest
- Compatible with existing deduplication