package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/jonobridge/grule-backend/adapters"
	"github.com/jonobridge/grule-backend/api"
	"github.com/jonobridge/grule-backend/audit"
	"github.com/jonobridge/grule-backend/capabilities"
	"github.com/jonobridge/grule-backend/capabilities/alerts"
	"github.com/jonobridge/grule-backend/capabilities/buffer"
	"github.com/jonobridge/grule-backend/capabilities/geofence"
	"github.com/jonobridge/grule-backend/capabilities/metrics"
	"github.com/jonobridge/grule-backend/capabilities/timing"
	"github.com/jonobridge/grule-backend/grule"
	"github.com/jonobridge/grule-backend/persistence"
)

func main() {
	log.Println("====================================")
	log.Println("🚀 NEW BACKEND v3.0.2 - Declarative Audit System")
	log.Println("====================================")
	log.Println("Starting GRULE Backend (Jammer) on port 8081...")

	// Enable Alert Audit by default (saves alerts to MySQL)
	if os.Getenv("GRULE_AUDIT_ENABLED") == "" {
		os.Setenv("GRULE_AUDIT_ENABLED", "Y")
		log.Println("✅ Alert Audit enabled (GRULE_AUDIT_ENABLED=Y)")
	}
	if os.Getenv("GRULE_AUDIT_LEVEL") == "" {
		os.Setenv("GRULE_AUDIT_LEVEL", "ALL")
		log.Println("✅ Alert Audit level set to ALL")
	}

	// Configure portal endpoint (like the old backend)
	portalEndpoint := os.Getenv("PORTAL_ENDPOINT")
	if portalEndpoint == "" {
		portalEndpoint = "/grule"
	} else if !strings.HasPrefix(portalEndpoint, "/") {
		portalEndpoint = "/" + portalEndpoint
	}
	log.Printf("📍 Endpoint: %s", portalEndpoint)

	// 1. Persistence
	store := persistence.NewMySQLStateStore()

	// Initialize Audit DB (for Progress Audit)
	audit.InitDB(store.GetDB())
	audit.CreateProgressAuditTable()

	// Initialize Alert Audit DB (for tracking system-generated alerts)
	audit.InitAlertAuditDB(store.GetDB())

	// 2. Registry & Capabilities
	reg := capabilities.NewRegistry()

	bufCap := buffer.NewBufferCapability()
	reg.Register(bufCap)

	metCap := metrics.NewMetricsCapability(bufCap)
	reg.Register(metCap)

	geoCap := geofence.NewGeofenceCapability(store)
	reg.Register(geoCap)

	timCap := timing.NewTimingCapability()
	reg.Register(timCap)

	alrtCap := alerts.NewAlertsCapability(store)
	reg.Register(alrtCap)

	// 3. Grule Components
	ctxBuilder := grule.NewContextBuilder(reg)
	adapter := adapters.NewGPSTrackerAdapter()

	// 4. Load Rules
	rules, err := store.LoadActiveRules()
	if err != nil {
		log.Fatalf("Failed to load rules: %v", err)
	}
	kbs := loadRulesFromSlice(rules)

	// Build audit manifest from rules
	manifest := audit.NewAuditManifest()
	if err := manifest.LoadFromRules(rules); err != nil {
		log.Printf("Warning: Failed to load audit manifests: %v", err)
	}
	log.Printf("Loaded %d audit rule entries", manifest.Count())

	// 5. Worker
	worker := grule.NewWorker(ctxBuilder, adapter, kbs, manifest)

	// 6. MQTT Client
	go startMQTT(worker)

	// 7. API Server
	server := &api.Server{
		Store:           store,
		CapabilitiesDir: "./capabilities",
		PortalEndpoint:  portalEndpoint,
		Worker:          worker,
		ReloadFunc: func() ([]*ast.KnowledgeBase, *audit.AuditManifest, error) {
			log.Println("🔄 Reloading rules and manifests...")
			// 1. Load active rules from DB
			newRules, err := store.LoadActiveRules()
			if err != nil {
				return nil, nil, err
			}

			// 2. Build new KnowledgeBases
			newKBs := loadRulesFromSlice(newRules)

			// 3. Build new Manifest
			newManifest := audit.NewAuditManifest()
			if err := newManifest.LoadFromRules(newRules); err != nil {
				log.Printf("Warning: Failed to reload audit manifests: %v", err)
			}

			log.Printf("✅ Reloaded %d rules and %d manifest entries", len(newRules), newManifest.Count())
			return newKBs, newManifest, nil
		},
	}

	// CORS Middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

	http.HandleFunc(portalEndpoint+"/api/rules", corsMiddleware(server.RulesHandler))
	http.HandleFunc(portalEndpoint+"/api/rules/", corsMiddleware(server.RuleByIDHandler))
	http.HandleFunc(portalEndpoint+"/api/rules/components", corsMiddleware(server.RuleComponentsHandler))
	http.HandleFunc(portalEndpoint+"/api/validate", corsMiddleware(server.ValidateHandler))
	http.HandleFunc(portalEndpoint+"/api/reload", corsMiddleware(server.ReloadHandler))
	http.HandleFunc(portalEndpoint+"/api/schema/capabilities", corsMiddleware(server.SchemaHandler))

	// Alert Audit Routes (for tracking alerts generated by the system)
	http.HandleFunc(portalEndpoint+"/api/audit/grid", corsMiddleware(server.AuditGridHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/summary", corsMiddleware(server.AuditSummaryHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/details", corsMiddleware(server.AuditDetailsHandler))

	// Progress Audit Routes
	http.HandleFunc(portalEndpoint+"/api/audit/progress/enable", corsMiddleware(server.ProgressEnableHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/disable", corsMiddleware(server.ProgressDisableHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/clear", corsMiddleware(server.ProgressClearHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/status", corsMiddleware(server.ProgressStatusHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/summary", corsMiddleware(server.ProgressSummaryHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/timeline", corsMiddleware(server.ProgressTimelineHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/snapshot", corsMiddleware(server.SnapshotHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/query", corsMiddleware(server.ProgressQueryHandler))

	// Swagger UI Routes
	http.HandleFunc(portalEndpoint+"/swagger", swaggerUIHandler(portalEndpoint))
	http.HandleFunc(portalEndpoint+"/swagger.json", swaggerJSONHandler(portalEndpoint))

	log.Println("Listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

// swaggerUIHandler serves the Swagger UI
func swaggerUIHandler(portalEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Grule Backend API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css">
    <style>body { margin: 0; padding: 0; }</style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: '` + portalEndpoint + `/swagger.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}
}

// swaggerJSONHandler serves the Swagger JSON specification
func swaggerJSONHandler(portalEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers for Swagger UI
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		swagger := map[string]interface{}{
			"swagger": "2.0",
			"info": map[string]interface{}{
				"title":       "Grule Backend API v2.0",
				"description": "Declarative Audit System - GPS Fleet Rules Manager",
				"version":     "2.0.0",
			},
			"host":     "jonobridge.madd.com.mx",
			"basePath": portalEndpoint + "/api",
			"schemes":  []string{"https", "http"},
			"paths": map[string]interface{}{
				"/rules": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":  "List all rules",
						"produces": []string{"application/json"},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Success"},
						},
					},
					"post": map[string]interface{}{
						"summary":  "Create new rule",
						"consumes": []string{"application/json"},
						"produces": []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{
								"in":       "body",
								"name":     "body",
								"required": true,
								"schema":   map[string]string{"$ref": "#/definitions/RuleInput"},
							},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Rule created"},
						},
					},
				},
				"/rules/{id}": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":  "Get rule by ID",
						"produces": []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{
								"name":     "id",
								"in":       "path",
								"required": true,
								"type":     "integer",
							},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Success"},
						},
					},
					"put": map[string]interface{}{
						"summary":  "Update rule",
						"consumes": []string{"application/json"},
						"produces": []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "id", "in": "path", "required": true, "type": "integer"},
							map[string]interface{}{"in": "body", "name": "body", "required": true, "schema": map[string]string{"$ref": "#/definitions/RuleInput"}},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Rule updated"},
						},
					},
					"delete": map[string]interface{}{
						"summary":  "Delete rule",
						"produces": []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "id", "in": "path", "required": true, "type": "integer"},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Rule deleted"},
						},
					},
				},
				"/rules/components": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":  "Get rule components",
						"produces": []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "rule_name", "in": "query", "required": true, "type": "string"},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Components retrieved"},
						},
					},
				},
				"/validate": map[string]interface{}{
					"post": map[string]interface{}{
						"summary":     "Validate GRL syntax",
						"description": "Validates Grule Rule Language syntax without saving",
						"consumes":    []string{"application/json"},
						"produces":    []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{
								"in":       "body",
								"name":     "body",
								"required": true,
								"schema": map[string]interface{}{
									"type":     "object",
									"required": []string{"grl"},
									"properties": map[string]interface{}{
										"grl": map[string]string{"type": "string", "example": "rule SpeedAlert \"Test\" salience 100 { when Jono.Speed > 0 then actions.Log(\"test\"); }"},
									},
								},
							},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{
								"description": "Validation result",
							},
						},
					},
				},
				"/reload": map[string]interface{}{
					"post": map[string]interface{}{
						"summary":     "Force reload rules",
						"description": "Forces immediate reload of all rules from database",
						"produces":    []string{"application/json"},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Rules reloaded"},
						},
					},
				},
				"/audit/grid": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":     "Get alert audit grid (jqGrid paginated)",
						"description": "Returns paginated alert summaries for all IMEIs with search and sort support",
						"produces":    []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "page", "in": "query", "type": "integer", "default": 1},
							map[string]interface{}{"name": "rows", "in": "query", "type": "integer", "default": 25},
							map[string]interface{}{"name": "sidx", "in": "query", "type": "string", "default": "last_alert_date"},
							map[string]interface{}{"name": "sord", "in": "query", "type": "string", "default": "DESC"},
							map[string]interface{}{"name": "searchText", "in": "query", "type": "string", "description": "Search filter"},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Success"},
						},
					},
				},
				"/audit/summary": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":     "Get alert summaries by IMEI",
						"description": "Returns summary of alerts grouped by IMEI",
						"produces":    []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "limit", "in": "query", "type": "integer", "default": 100, "description": "Maximum number of results"},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Success"},
						},
					},
				},
				"/audit/details": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":     "Get alert details for specific IMEI",
						"description": "Returns detailed alert information for a given IMEI",
						"produces":    []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "imei", "in": "query", "required": true, "type": "string", "description": "IMEI to filter"},
							map[string]interface{}{"name": "limit", "in": "query", "type": "integer", "default": 100, "description": "Maximum number of results"},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Success"},
						},
					},
				},
				"/audit/progress/enable": map[string]interface{}{
					"post": map[string]interface{}{
						"summary":     "Enable progress audit",
						"description": "Activates progress audit tracking",
						"produces":    []string{"application/json"},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Progress audit enabled"},
						},
					},
				},
				"/audit/progress/disable": map[string]interface{}{
					"post": map[string]interface{}{
						"summary":     "Disable progress audit",
						"description": "Deactivates progress audit tracking",
						"produces":    []string{"application/json"},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Progress audit disabled"},
						},
					},
				},
				"/audit/progress/clear": map[string]interface{}{
					"post": map[string]interface{}{
						"summary":     "Clear progress audit data",
						"description": "Deletes all progress audit records",
						"produces":    []string{"application/json"},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Data cleared"},
						},
					},
				},
				"/audit/progress/status": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":     "Get progress audit status",
						"description": "Returns current state (enabled/disabled)",
						"produces":    []string{"application/json"},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Status retrieved"},
						},
					},
				},
				"/audit/progress/query": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":     "Query progress by IMEI",
						"description": "Retrieves execution progress for specific IMEI",
						"produces":    []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "imei", "in": "query", "required": true, "type": "string"},
							map[string]interface{}{"name": "limit", "in": "query", "type": "integer", "default": 50},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Progress data retrieved"},
						},
					},
				},
				"/audit/progress/summary": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":     "Get IMEI summary (jqGrid Level 1)",
						"description": "Returns IMEIs with max step, total frames",
						"produces":    []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "page", "in": "query", "type": "integer", "default": 1},
							map[string]interface{}{"name": "rows", "in": "query", "type": "integer", "default": 25},
							map[string]interface{}{"name": "sidx", "in": "query", "type": "string", "default": "last_frame_time"},
							map[string]interface{}{"name": "sord", "in": "query", "type": "string", "default": "DESC"},
							map[string]interface{}{"name": "rule_name", "in": "query", "type": "string"},
							map[string]interface{}{"name": "imei", "in": "query", "type": "string"},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Summary retrieved"},
						},
					},
				},
				"/audit/progress/timeline": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":     "Get frame timeline (jqGrid Level 2)",
						"description": "Returns chronological execution frames for IMEI",
						"produces":    []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "imei", "in": "query", "required": true, "type": "string"},
							map[string]interface{}{"name": "page", "in": "query", "type": "integer", "default": 1},
							map[string]interface{}{"name": "rows", "in": "query", "type": "integer", "default": 20},
							map[string]interface{}{"name": "sidx", "in": "query", "type": "string", "default": "execution_time"},
							map[string]interface{}{"name": "sord", "in": "query", "type": "string", "default": "ASC"},
							map[string]interface{}{"name": "rule_name", "in": "query", "type": "string"},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Timeline retrieved"},
						},
					},
				},
				"/audit/progress/snapshot": map[string]interface{}{
					"get": map[string]interface{}{
						"summary":  "Get snapshot by ID",
						"produces": []string{"application/json"},
						"parameters": []interface{}{
							map[string]interface{}{"name": "id", "in": "query", "required": true, "type": "integer"},
						},
						"responses": map[string]interface{}{
							"200": map[string]interface{}{"description": "Snapshot retrieved"},
						},
					},
				},
			},
			"definitions": map[string]interface{}{
				"RuleInput": map[string]interface{}{
					"type":     "object",
					"required": []string{"name", "grl"},
					"properties": map[string]interface{}{
						"name":     map[string]string{"type": "string", "example": "SpeedAlert"},
						"grl":      map[string]string{"type": "string", "example": "rule SpeedAlert {...}"},
						"priority": map[string]interface{}{"type": "integer", "default": 100},
						"active":   map[string]interface{}{"type": "boolean", "default": true},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(swagger)
	}
}

func loadRulesFromSlice(rules []persistence.Rule) []*ast.KnowledgeBase {
	var kbs []*ast.KnowledgeBase
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
		kbs = append(kbs, kbInstance)
	}
	return kbs
}

func startMQTT(worker *grule.Worker) {
	broker := os.Getenv("MQTT_BROKER_HOST")
	if broker == "" {
		broker = "mosquitto"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + broker + ":1883")
	opts.SetClientID("grule-backend-refactor-" + os.Getenv("HOSTNAME"))
	opts.SetAutoReconnect(true)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		go worker.Process(payload)
	})

	client := mqtt.NewClient(opts)
	for {
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			log.Printf("Connected to MQTT at %s", broker)
			client.Subscribe("tracker/jonoprotocol", 0, nil)
			break
		}
		time.Sleep(5 * time.Second)
		log.Printf("Retrying MQTT connection to %s...", broker)
	}
}
