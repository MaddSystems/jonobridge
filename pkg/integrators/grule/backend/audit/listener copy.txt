package audit

import (
	"context"
	"log"
	"sync"

	"github.com/hyperjumptech/grule-rule-engine/ast"
)

type AuditListener struct {
	manifest    *AuditManifest
	enabled     bool            // Global kill switch
	loggedOnce  map[string]bool // Track logged unknown rules
	loggedMu    sync.Mutex
	dataContext ast.IDataContext // Store data context for use in ExecuteRuleEntry
	dcMu        sync.RWMutex
}

func NewAuditListener(manifest *AuditManifest) *AuditListener {
	return &AuditListener{
		manifest:   manifest,
		enabled:    true,
		loggedOnce: make(map[string]bool),
	}
}

// SetEnabled - Global kill switch for performance
func (l *AuditListener) SetEnabled(enabled bool) {
	l.enabled = enabled
}

// EvaluateRuleEntry is called when a rule is being evaluated (required by GruleEngineListener)
func (l *AuditListener) EvaluateRuleEntry(ctx context.Context, cycle uint64, entry *ast.RuleEntry, candidate bool) {
	// Not used for audit capture, but required by interface
}

// BeginCycle is called at the start of each evaluation cycle (required by GruleEngineListener)
func (l *AuditListener) BeginCycle(ctx context.Context, cycle uint64) {
	// Not used for audit capture, but required by interface
}

// ExecuteRuleEntry is called when a rule is being executed (required by GruleEngineListener)
func (l *AuditListener) ExecuteRuleEntry(ctx context.Context, cycle uint64, entry *ast.RuleEntry) {
	// Use EXISTING global toggle from audit.IsProgressAuditEnabled()
	// This ensures frontend controls (Activar/Desactivar) still work
	if !IsProgressAuditEnabled() {
		return
	}

	// Get data context from context value
	dcValue := ctx.Value("dataContext")
	if dcValue == nil {
		log.Printf("[AuditListener] No dataContext in context for rule '%s'", entry.RuleName)
		return
	}

	dc, ok := dcValue.(ast.IDataContext)
	if !ok {
		log.Printf("[AuditListener] Invalid dataContext type for rule '%s'", entry.RuleName)
		return
	}

	meta := l.manifest.GetRuleMeta(entry.RuleName)

	// Debug: Log all rule executions
	log.Printf("[AuditListener] Rule executed: '%s', Meta found: %v", entry.RuleName, meta != nil)

	// Unknown rule - log once, skip audit
	if meta == nil {
		l.logOnce(entry.RuleName)
		return
	}

	// Explicitly disabled
	if !meta.Enabled {
		log.Printf("[AuditListener] Rule '%s' disabled in manifest", entry.RuleName)
		return
	}

	// Debug: Log if it's an alert
	if meta.IsAlert {
		log.Printf("[AuditListener] ALERT RULE EXECUTED: '%s', IsAlert=true", entry.RuleName)
	}

	// Extract IMEI from packet for database storage
	packetObj := dc.Get("IncomingPacket")
	imei := getIMEI(packetObj)
	if imei == "" {
		// Try to get IMEI from context if packet extraction fails
		if contextIMEI, ok := ctx.Value("imei").(string); ok {
			imei = contextIMEI
		}
	}

	// Get original packet from context (to bypass possible DataContext wrappers/proxies)
	originalPacket := ctx.Value("originalPacket")

	// Extract snapshot with nil-safety
	snapshot, err := extractSnapshot(dc, imei, originalPacket)
	if err != nil {
		log.Printf("[AuditListener] Snapshot error for '%s': %v", entry.RuleName, err)
		snapshot = map[string]interface{}{"error": err.Error()}
	}

	// Synchronous capture for testing reliability (can be made async later if needed)
	Capture(&AuditEntry{
		IMEI:         imei,
		RuleName:     entry.RuleName,
		Salience:     int(entry.Salience),
		Description:  meta.Description,
		Level:        meta.Level,
		IsAlert:      meta.IsAlert,
		StepNumber:   meta.Order,       // From manifest
		StageReached: meta.Description, // Use description as stage name
		Snapshot:     snapshot,
	})
}

// logOnce logs unknown rule warning only once
func (l *AuditListener) logOnce(ruleName string) {
	l.loggedMu.Lock()
	defer l.loggedMu.Unlock()
	if !l.loggedOnce[ruleName] {
		log.Printf("[AuditListener] Rule '%s' not in manifest, skipping audit", ruleName)
		l.loggedOnce[ruleName] = true
	}
}
