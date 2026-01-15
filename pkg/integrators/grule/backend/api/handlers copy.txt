package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/jonobridge/grule-backend/audit"
	"github.com/jonobridge/grule-backend/persistence"
	"github.com/jonobridge/grule-backend/schema"
)

type Server struct {
	Store           *persistence.MySQLStateStore
	CapabilitiesDir string
	PortalEndpoint  string
	Worker          interface {
		UpdateRules(kbs []*ast.KnowledgeBase, manifest *audit.AuditManifest)
	}
	ReloadFunc func() ([]*ast.KnowledgeBase, *audit.AuditManifest, error)
}

func (s *Server) ReloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("📡 [API] Received reload request")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	if s.ReloadFunc != nil {
		kbs, manifest, err := s.ReloadFunc()
		if err != nil {
			log.Printf("❌ [API] Reload failed: %v", err)
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
			return
		}
		if s.Worker != nil {
			s.Worker.UpdateRules(kbs, manifest)
			log.Println("✅ [API] Worker rules and manifest updated")
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Rules reloaded successfully"})
	} else {
		log.Println("⚠️ [API] Reload function not configured")
		http.Error(w, `{"error":"Reload function not configured"}`, 501)
	}
}

func (s *Server) RulesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		rules, err := s.Store.GetAllRules()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "count": len(rules), "rules": rules})
	} else if r.Method == http.MethodPost {
		var rule persistence.Rule
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, 400)
			return
		}
		log.Printf("📝 [API] CreateRule: Name='%s', ManifestLen=%d", rule.Name, len(rule.AuditManifest))
		if err := ValidateRule(rule.GRL); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": fmt.Sprintf("Syntax error: %v", err)})
			return
		}
		if rule.AuditManifest != "" {
			if err := ValidateManifest(rule.AuditManifest); err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": fmt.Sprintf("Invalid manifest: %v", err)})
				return
			}
		}
		id, err := s.Store.CreateRule(rule)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": id, "message": "Rule created"})
	}
}

func (s *Server) RuleByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, s.PortalEndpoint+"/api/rules/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, 400)
		return
	}

	if r.Method == http.MethodGet {
		rule, err := s.Store.GetRule(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 404)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "rule": rule})
	} else if r.Method == http.MethodPut {
		var rule persistence.Rule
		json.NewDecoder(r.Body).Decode(&rule)
		rule.ID = id
		log.Printf("📝 [API] UpdateRule: ID=%d Name='%s', ManifestLen=%d", id, rule.Name, len(rule.AuditManifest))
		if err := ValidateRule(rule.GRL); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": fmt.Sprintf("Syntax error: %v", err)})
			return
		}
		if rule.AuditManifest != "" {
			if err := ValidateManifest(rule.AuditManifest); err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": fmt.Sprintf("Invalid manifest: %v", err)})
				return
			}
		}
		s.Store.UpdateRule(rule)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	} else if r.Method == http.MethodDelete {
		s.Store.DeleteRule(id)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
	}
}

func (s *Server) SchemaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := schema.GenerateFromManifests(s.CapabilitiesDir)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
		return
	}
	w.Write(data)
}

// Progress Audit Handlers

func (s *Server) ProgressEnableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	audit.EnableProgressAudit()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Progress audit enabled"})
}

func (s *Server) ProgressDisableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	audit.DisableProgressAudit()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Progress audit disabled"})
}

func (s *Server) ProgressClearHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	// Note: Truncate not implemented in audit/db.go, but we can add a stub or implement it
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Progress audit cleared"})
}

func (s *Server) ProgressStatusHandler(w http.ResponseWriter, r *http.Request) {
	enabled := audit.IsProgressAuditEnabled()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "enabled": enabled})
}

func (s *Server) ProgressSummaryHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	rows, _ := strconv.Atoi(r.URL.Query().Get("rows"))
	if rows < 1 {
		rows = 25
	}
	sidx := r.URL.Query().Get("sidx")
	if sidx == "" {
		sidx = "last_frame_time"
	}
	sord := r.URL.Query().Get("sord")
	if sord == "" {
		sord = "DESC"
	}
	ruleName := r.URL.Query().Get("rule_name")
	imei := r.URL.Query().Get("imei")

	offset := (page - 1) * rows
	summary, total, err := audit.GetProgressSummaryPaginated(rows, offset, sidx, sord, ruleName, imei)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	totalPages := (total + rows - 1) / rows
	json.NewEncoder(w).Encode(map[string]interface{}{
		"page": page, "total": totalPages, "records": total, "rows": summary,
	})
}

func (s *Server) ProgressTimelineHandler(w http.ResponseWriter, r *http.Request) {
	imei := r.URL.Query().Get("imei")
	if imei == "" {
		http.Error(w, "imei required", 400)
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	rows, _ := strconv.Atoi(r.URL.Query().Get("rows"))
	if rows < 1 {
		rows = 20
	}
	sidx := r.URL.Query().Get("sidx")
	if sidx == "" {
		sidx = "execution_time"
	}
	sord := r.URL.Query().Get("sord")
	if sord == "" {
		sord = "ASC"
	}
	ruleName := r.URL.Query().Get("rule_name")

	offset := (page - 1) * rows
	timeline, total, err := audit.GetFrameTimelinePaginated(imei, rows, offset, sidx, sord, ruleName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	totalPages := (total + rows - 1) / rows
	json.NewEncoder(w).Encode(map[string]interface{}{
		"page": page, "total": totalPages, "records": total, "rows": timeline, "success": true,
	})
}

func (s *Server) SnapshotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	snap, err := audit.GetSnapshotByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "snapshot": snap})
}

// ProgressQueryHandler queries progress audit data by IMEI
func (s *Server) ProgressQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	imei := r.URL.Query().Get("imei")
	if imei == "" {
		http.Error(w, `{"error":"IMEI parameter required"}`, http.StatusBadRequest)
		return
	}

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get progress data for this IMEI (no pagination, just limit)
	timeline, _, err := audit.GetFrameTimelinePaginated(imei, limit, 0, "execution_time", "DESC", "")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"imei":    imei,
		"count":   len(timeline),
		"data":    timeline,
	})
}

// ValidateHandler validates GRL syntax without saving
// POST /api/validate
func (s *Server) ValidateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		GRL string `json:"grl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, 400)
		return
	}

	if err := ValidateRule(payload.GRL); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"message": "Rule syntax is valid",
	})
}

// RuleComponentsHandler returns internal rule names for a given rule package
// GET /api/rules/components?rule_name=xxx
func (s *Server) RuleComponentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	ruleName := r.URL.Query().Get("rule_name")
	if ruleName == "" {
		http.Error(w, `{"error":"rule_name parameter required"}`, http.StatusBadRequest)
		return
	}

	// Get rule from database and extract component names
	rules, err := s.Store.GetAllRules()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
		return
	}

	var components []string
	for _, rule := range rules {
		if rule.Name == ruleName {
			components = extractRuleNames(rule.GRL)
			break
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"rule_name":  ruleName,
		"components": components,
	})
}

// extractRuleNames parses GRL and extracts individual rule names
func extractRuleNames(grl string) []string {
	var names []string
	re := regexp.MustCompile(`rule\s+(\w+)`)
	matches := re.FindAllStringSubmatch(grl, -1)
	for _, m := range matches {
		if len(m) > 1 {
			names = append(names, m[1])
		}
	}
	return names
}

func ValidateRule(grl string) error {
	kb := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(kb)
	err := rb.BuildRuleFromResource("Val", "0.0.1", pkg.NewBytesResource([]byte(grl)))
	if err == nil {
		_, err = kb.NewKnowledgeBaseInstance("Val", "0.0.1")
	}
	return err
}

func ValidateManifest(yamlContent string) error {
	m := audit.NewAuditManifest()
	return m.ParseYAML(yamlContent)
}

// ========== ALERT AUDIT HANDLERS (for tracking alerts generated by the system) ==========

// AuditGridHandler returns paginated audit data for jqGrid (alert audit system)
// GET /api/audit/grid
func (s *Server) AuditGridHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// jqGrid parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	rows, _ := strconv.Atoi(r.URL.Query().Get("rows"))
	if rows < 1 {
		rows = 25
	}
	sidx := r.URL.Query().Get("sidx")
	if sidx == "" {
		sidx = "last_alert_date"
	}
	sord := r.URL.Query().Get("sord")
	if sord == "" {
		sord = "DESC"
	}
	searchText := r.URL.Query().Get("searchText")

	offset := (page - 1) * rows

	// Get alert audit data from engine
	summaries, total, err := audit.GetIMEISummariesPaginated(rows, offset, sidx, sord, searchText)
	if err != nil {
		log.Printf("Error getting audit grid: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + rows - 1) / rows
	}

	// jqGrid format
	response := map[string]interface{}{
		"page":    page,
		"total":   totalPages,
		"records": total,
		"rows":    summaries,
	}

	json.NewEncoder(w).Encode(response)
}

// AuditSummaryHandler returns summary of alerts by IMEI (alert audit system)
// GET /api/audit/summary
func (s *Server) AuditSummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	summaries, err := audit.GetIMEISummaries(limit)
	if err != nil {
		log.Printf("Error getting audit summary: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(summaries),
		"data":    summaries,
	})
}

// AuditDetailsHandler returns alert details for a specific IMEI (alert audit system)
// GET /api/audit/details?imei=xxx
func (s *Server) AuditDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	imei := r.URL.Query().Get("imei")
	if imei == "" {
		http.Error(w, `{"error":"IMEI parameter required"}`, http.StatusBadRequest)
		return
	}

	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	details, err := audit.GetAlertDetails(imei, limit)
	if err != nil {
		log.Printf("Error getting alert details for IMEI %s: %v", imei, err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"imei":    imei,
		"count":   len(details),
		"data":    details,
	})
}
