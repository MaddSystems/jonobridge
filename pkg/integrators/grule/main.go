package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonobridge/grule-engine/engine"
	"github.com/jonobridge/grule-engine/engine/audit"
)

var (
	portalEndpoint string
	mqttClient     mqtt.Client
)

func main() {
	// ---------------------------------------------------------
	// LOG DE VERIFICACIÓN DE DEPLOY (Simple y visible)
	// ---------------------------------------------------------
	log.Println("==== NEW DEPLOYMENT VERIFICATION LOG v1.0.11 ====")
	log.Println("==== CHECK: PacketWrapper has IsOfflineFor5Min ====")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("=== Iniciando Grule Engine Universal ===")

	engine.InitializeFixedBuffers(24 * time.Hour)

	// Configurar endpoint del portal
	portalEndpoint = os.Getenv("PORTAL_ENDPOINT")
	if portalEndpoint == "" {
		portalEndpoint = "/grule"
	} else if !strings.HasPrefix(portalEndpoint, "/") {
		portalEndpoint = "/" + portalEndpoint
	}

	auditEnabled := os.Getenv("GRULE_AUDIT_ENABLED")
	auditLevel := os.Getenv("GRULE_AUDIT_LEVEL")
	log.Printf("GRULE_AUDIT_ENABLED: %s, GRULE_AUDIT_LEVEL: %s", auditEnabled, auditLevel)
	log.Printf("PORTAL_ENDPOINT: %s", portalEndpoint)

	// Inicializar motor (MySQL + reglas)
	engine.Initialize()
	log.Println("✅ Engine inicializado correctamente")

	// Inicializar WorkerPool ANTES de suscribirse a MQTT
	engine.InitializeWorkerPool(nil)
	log.Println("✅ WorkerPool inicializado correctamente")

	// Inicializar auditoría si está habilitada
	if auditEnabled == "Y" {
		db := engine.GetDB()
		if db != nil {
			audit.InitDB(db)
			audit.CreateProgressAuditTable()
			log.Printf("✅ Audit layer initialized (level: %s)", auditLevel)
		} else {
			log.Println("⚠️  Could not initialize audit: database connection not available")
		}
	}

	// Configurar cliente MQTT
	broker := os.Getenv("MQTT_BROKER_HOST")
	if broker == "" {
		log.Println("❌ MQTT_BROKER_HOST no definido")
		return
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + broker + ":1883")
	opts.SetClientID("GRULE_ENGINE_UNIVERSAL_" + os.Getenv("HOSTNAME") + "_" + fmt.Sprint(time.Now().UnixNano()%1e6))
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(func(_ mqtt.Client, msg mqtt.Message) {
		// ProcessPacketMessage envía al worker pool
		engine.ProcessPacketMessage(string(msg.Payload()))
	})

	mqttClient = mqtt.NewClient(opts)
	log.Printf("📡 Conectando a MQTT broker: %s", broker)

	for {
		if token := mqttClient.Connect(); token.Wait() && token.Error() == nil {
			log.Println("✅ Grule Engine conectado a MQTT")
			break
		}
		log.Println("⚠️  Reintentando conexión MQTT en 5 segundos...")
		time.Sleep(5 * time.Second)
	}

	mqttClient.Subscribe("tracker/jonoprotocol", 1, nil)
	log.Println("✅ Grule Engine suscrito a tracker/jonoprotocol")

	// Configurar API REST
	setupAPI()

	// Iniciar servidor HTTP en puerto 8081
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🌐 API + Swagger running on port %s", port)
	log.Printf("📍 Endpoint: %s", portalEndpoint)
	log.Println("=== Sistema Universal Listo ===")

	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	select {} // Keep running
}

func setupAPI() {
	// Configurar CORS para permitir Flask (external-web) acceder a la API
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

	// API REST Endpoints - Sistema de gestión de reglas
	http.HandleFunc(portalEndpoint+"/api/rules", corsMiddleware(rulesHandler))
	http.HandleFunc(portalEndpoint+"/api/rules/", corsMiddleware(ruleByIDHandler))
	http.HandleFunc(portalEndpoint+"/api/rules/components", corsMiddleware(ruleComponentsHandler))
	http.HandleFunc(portalEndpoint+"/api/validate", corsMiddleware(validateHandler))
	http.HandleFunc(portalEndpoint+"/api/reload", corsMiddleware(reloadHandler))
	http.HandleFunc(portalEndpoint+"/api/health", corsMiddleware(healthHandler))

	// Audit API Endpoints - Sistema universal de auditoría
	http.HandleFunc(portalEndpoint+"/api/audit/summary", corsMiddleware(auditSummaryHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/details", corsMiddleware(auditDetailsHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/grid", corsMiddleware(auditGridHandler))

	// Progress Audit API Endpoints - Control de auditoría de progreso
	http.HandleFunc(portalEndpoint+"/api/audit/progress/enable", corsMiddleware(progressEnableHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/disable", corsMiddleware(progressDisableHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/clear", corsMiddleware(progressClearHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/status", corsMiddleware(progressStatusHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress", corsMiddleware(progressQueryHandler))

	// Progress Audit - Movie Frames API (jqGrid 2 niveles)
	http.HandleFunc(portalEndpoint+"/api/audit/progress/rules", corsMiddleware(progressRulesHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/summary", corsMiddleware(progressSummaryHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/timeline", corsMiddleware(progressTimelineHandler))
	http.HandleFunc(portalEndpoint+"/api/audit/progress/snapshot", corsMiddleware(snapshotHandler))

	// Swagger UI (documentación interactiva)
	http.HandleFunc(portalEndpoint+"/swagger", swaggerUIHandler)
	http.HandleFunc(portalEndpoint+"/swagger.json", swaggerJSONHandler)

	// Root - redirige a Swagger
	http.HandleFunc(portalEndpoint+"/", rootHandler)
	http.HandleFunc(portalEndpoint, rootHandler)
} // GET /api/rules - Listar reglas
// POST /api/rules - Crear regla
func rulesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		rules, err := engine.GetAllRules()
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"count":   len(rules),
			"rules":   rules,
		})

	case http.MethodPost:
		var rule engine.Rule
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		// Validar sintaxis antes de guardar
		if err := engine.ValidateRule(rule.GRL); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Syntax error: %v", err),
			})
			return
		}

		id, err := engine.CreateRule(rule)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"id":      id,
			"message": "Rule created successfully",
		})

	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// GET /api/rules/:id - Ver regla
// PUT /api/rules/:id - Actualizar regla
// DELETE /api/rules/:id - Eliminar regla
func ruleByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extraer ID de la URL
	path := strings.TrimPrefix(r.URL.Path, portalEndpoint+"/api/rules/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Invalid rule ID"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		rule, err := engine.GetRule(id)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"rule":    rule,
		})

	case http.MethodPut:
		var rule engine.Rule
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}
		rule.ID = id

		// Validar sintaxis
		if err := engine.ValidateRule(rule.GRL); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Syntax error: %v", err),
			})
			return
		}

		if err := engine.UpdateRule(rule); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Rule updated successfully",
		})

	case http.MethodDelete:
		if err := engine.DeleteRule(id); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Rule deleted successfully",
		})

	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// POST /api/validate - Validar sintaxis de regla
func validateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		GRL string `json:"grl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if err := engine.ValidateRule(payload.GRL); err != nil {
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

// POST /api/reload - Forzar recarga de reglas
func reloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	engine.ForceReload()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Rules reloaded successfully",
	})
}

// GET /api/health - Health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "grule-engine-universal",
		"version":   "2.0.0",
	})
}

// GET /api/audit/summary - Resumen de IMEIs con alertas (sin duplicados)
func auditSummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		log.Printf("Error obteniendo summary: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(summaries),
		"data":    summaries,
	})
}

// GET /api/audit/details - Detalles de alertas por IMEI
func auditDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		log.Printf("Error obteniendo details para IMEI %s: %v", imei, err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"imei":    imei,
		"count":   len(details),
		"data":    details,
	})
}

// GET /api/audit/grid - Endpoint para jqGrid (paginación, búsqueda, ordenamiento)
func auditGridHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Parámetros jqGrid
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	rows, _ := strconv.Atoi(r.URL.Query().Get("rows"))
	if rows < 1 {
		rows = 25
	}

	sidx := r.URL.Query().Get("sidx") // columna para ordenar
	if sidx == "" {
		sidx = "last_alert_date"
	}

	sord := r.URL.Query().Get("sord") // ASC o DESC
	if sord == "" {
		sord = "DESC"
	}

	searchText := r.URL.Query().Get("searchText")

	// Calcular offset
	offset := (page - 1) * rows

	// Obtener datos con paginación
	summaries, total, err := audit.GetIMEISummariesPaginated(rows, offset, sidx, sord, searchText)
	if err != nil {
		log.Printf("Error en auditGridHandler: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	// Calcular total de páginas
	totalPages := (total + rows - 1) / rows

	// Formato jqGrid
	response := map[string]interface{}{
		"page":    page,
		"total":   totalPages,
		"records": total,
		"rows":    summaries,
	}

	json.NewEncoder(w).Encode(response)
}

// Swagger UI Handler
func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Grule Engine API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '` + portalEndpoint + `/swagger.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// Swagger JSON
func swaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	swagger := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       "Grule Engine API",
			"description": "GPS Fleet Rules Manager - Universal Audit System",
			"version":     "2.0.0",
		},
		"host":     "jonobridge.madd.com.mx",
		"basePath": portalEndpoint + "/api",
		"schemes":  []string{"https"},
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
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"valid":   map[string]string{"type": "boolean"},
									"error":   map[string]string{"type": "string"},
									"message": map[string]string{"type": "string"},
								},
							},
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
			"/health": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":  "Health check",
					"produces": []string{"application/json"},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Service is healthy",
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"status":    map[string]string{"type": "string"},
									"timestamp": map[string]string{"type": "string"},
									"service":   map[string]string{"type": "string"},
									"version":   map[string]string{"type": "string"},
								},
							},
						},
					},
				},
			},
			"/audit/summary": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get audit summary by IMEI",
					"description": "Returns summary of alerts by IMEI (one row per IMEI)",
					"produces":    []string{"application/json"},
					"parameters": []interface{}{
						map[string]interface{}{
							"name":        "limit",
							"in":          "query",
							"type":        "integer",
							"default":     100,
							"description": "Maximum number of results",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "Success"},
					},
				},
			},
			"/audit/details": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get audit details by IMEI",
					"description": "Returns full alert history for a specific IMEI",
					"produces":    []string{"application/json"},
					"parameters": []interface{}{
						map[string]interface{}{
							"name":        "imei",
							"in":          "query",
							"type":        "string",
							"required":    true,
							"description": "IMEI to filter",
						},
						map[string]interface{}{
							"name":        "limit",
							"in":          "query",
							"type":        "integer",
							"default":     100,
							"description": "Maximum number of results",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "Success"},
					},
				},
			},
			"/audit/grid": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get audit grid data (jqGrid)",
					"description": "Returns paginated audit data for jqGrid with search and sort",
					"produces":    []string{"application/json"},
					"parameters": []interface{}{
						map[string]interface{}{
							"name":    "page",
							"in":      "query",
							"type":    "integer",
							"default": 1,
						},
						map[string]interface{}{
							"name":    "rows",
							"in":      "query",
							"type":    "integer",
							"default": 25,
						},
						map[string]interface{}{
							"name":    "sidx",
							"in":      "query",
							"type":    "string",
							"default": "last_alert_date",
						},
						map[string]interface{}{
							"name":    "sord",
							"in":      "query",
							"type":    "string",
							"default": "DESC",
						},
						map[string]interface{}{
							"name":        "searchText",
							"in":          "query",
							"type":        "string",
							"description": "Search filter",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "Success"},
					},
				},
			},
			"/audit/progress/enable": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "Enable progress audit",
					"description": "Activates progress audit tracking to monitor rule execution states",
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
					"description": "Deletes all stored progress audit records",
					"produces":    []string{"application/json"},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "Progress audit data cleared"},
					},
				},
			},
			"/audit/progress/status": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get progress audit status",
					"description": "Returns current state of progress audit (enabled/disabled)",
					"produces":    []string{"application/json"},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "Status retrieved"},
					},
				},
			},
			"/audit/progress": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Query progress audit by IMEI",
					"description": "Retrieves rule execution progress history for a specific IMEI",
					"produces":    []string{"application/json"},
					"parameters": []interface{}{
						map[string]interface{}{
							"name":        "imei",
							"in":          "query",
							"type":        "string",
							"required":    true,
							"description": "IMEI to query",
						},
						map[string]interface{}{
							"name":        "limit",
							"in":          "query",
							"type":        "integer",
							"default":     50,
							"description": "Maximum number of results",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{"description": "Progress data retrieved"},
					},
				},
			},
			"/audit/progress/rules": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "List available rules for progress audit",
					"description": "Returns list of rules with frame counts and last execution time (for rule selector in frontend)",
					"produces":    []string{"application/json"},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "List of rules retrieved",
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"success": map[string]string{"type": "boolean"},
									"count":   map[string]string{"type": "integer"},
									"rules": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"type": "object",
											"properties": map[string]interface{}{
												"rule_name":      map[string]string{"type": "string"},
												"total_frames":   map[string]string{"type": "integer"},
												"total_imeis":    map[string]string{"type": "integer"},
												"last_execution": map[string]string{"type": "string"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/audit/progress/summary": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get IMEIs summary with max step reached (jqGrid Level 1)",
					"description": "Returns list of IMEIs with their maximum step reached, total frames, and last execution time. Optionally filtered by rule name.",
					"produces":    []string{"application/json"},
					"parameters": []interface{}{
						map[string]interface{}{
							"name":        "rule_name",
							"in":          "query",
							"type":        "string",
							"required":    false,
							"description": "Filter by specific rule name",
						},
						map[string]interface{}{
							"name":        "limit",
							"in":          "query",
							"type":        "integer",
							"default":     100,
							"description": "Maximum number of results",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Summary retrieved",
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"success": map[string]string{"type": "boolean"},
									"count":   map[string]string{"type": "integer"},
									"data": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"type": "object",
											"properties": map[string]interface{}{
												"imei":            map[string]string{"type": "string"},
												"rule_name":       map[string]string{"type": "string"},
												"max_step":        map[string]string{"type": "integer"},
												"total_frames":    map[string]string{"type": "integer"},
												"last_frame_time": map[string]string{"type": "string"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/audit/progress/timeline": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get chronological timeline of frames for an IMEI (jqGrid Level 2)",
					"description": "Returns ordered list of execution frames showing buffer evolution, metrics, geofences, and flags. Shows 'movie frames' of rule execution progress.",
					"produces":    []string{"application/json"},
					"parameters": []interface{}{
						map[string]interface{}{
							"name":        "imei",
							"in":          "query",
							"type":        "string",
							"required":    true,
							"description": "IMEI to query timeline for",
						},
						map[string]interface{}{
							"name":        "rule_name",
							"in":          "query",
							"type":        "string",
							"required":    false,
							"description": "Filter by specific rule name",
						},
						map[string]interface{}{
							"name":        "limit",
							"in":          "query",
							"type":        "integer",
							"default":     500,
							"description": "Maximum number of frames",
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Timeline retrieved",
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"success": map[string]string{"type": "boolean"},
									"imei":    map[string]string{"type": "string"},
									"count":   map[string]string{"type": "integer"},
									"frames": map[string]interface{}{
										"type": "array",
										"items": map[string]interface{}{
											"type": "object",
											"properties": map[string]interface{}{
												"id":             map[string]string{"type": "integer"},
												"rule_name":      map[string]string{"type": "string"},
												"step_number":    map[string]string{"type": "integer"},
												"stage_reached":  map[string]string{"type": "string"},
												"stop_reason":    map[string]string{"type": "string"},
												"buffer_size":    map[string]string{"type": "integer"},
												"metrics_ready":  map[string]string{"type": "boolean"},
												"geofence_eval":  map[string]string{"type": "string"},
												"snapshot":       map[string]string{"type": "object", "description": "JSON with buffer_circular, jammer_metrics, geofence_checks, wrapper_flags, packet_current"},
												"execution_time": map[string]string{"type": "string"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"definitions": map[string]interface{}{
			"Rule": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id":       map[string]string{"type": "integer"},
					"name":     map[string]string{"type": "string"},
					"grl":      map[string]string{"type": "string"},
					"active":   map[string]string{"type": "boolean"},
					"priority": map[string]string{"type": "integer"},
				},
			},
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

// Root handler - redirige a Swagger UI
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, portalEndpoint+"/swagger", http.StatusMovedPermanently)
}

// ========================== PROGRESS AUDIT HANDLERS ==========================

// POST /api/audit/progress/enable - Activa auditoría de progreso
func progressEnableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	if err := audit.EnableProgressAudit(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Progress audit enabled",
		"status":  "enabled",
	})
}

// POST /api/audit/progress/disable - Desactiva auditoría de progreso
func progressDisableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	if err := audit.DisableProgressAudit(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Progress audit disabled",
		"status":  "disabled",
	})
}

// POST /api/audit/progress/clear - Limpia datos de auditoría de progreso
func progressClearHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	if err := audit.ClearProgressAudit(); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Progress audit data cleared",
	})
}

// GET /api/audit/progress/status - Obtiene estado actual de auditoría de progreso
func progressStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	enabled := audit.IsProgressAuditEnabled()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"enabled": enabled,
		"status": map[string]interface{}{
			"enabled":   enabled,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}

// GET /api/audit/progress?imei=xxx - Consulta progreso por IMEI
func progressQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	progressList, err := audit.GetProgressByIMEI(imei, limit)
	if err != nil {
		log.Printf("Error obteniendo progress para IMEI %s: %v", imei, err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"imei":    imei,
		"count":   len(progressList),
		"data":    progressList,
	})
}

// GET /api/audit/progress/rules - Lista de reglas disponibles para auditoría
func progressRulesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	rules, err := audit.GetAvailableRules()
	if err != nil {
		log.Printf("Error obteniendo reglas disponibles: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(rules),
		"rules":   rules,
	})
}

// GET /api/audit/progress/summary - jqGrid Nivel 1: Lista de IMEIs con max step (con paginación)
func progressSummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Parámetros de jqGrid para paginación y ordenamiento
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	rows, _ := strconv.Atoi(r.URL.Query().Get("rows"))
	if rows < 1 {
		rows = 25 // Default
	}
	sidx := r.URL.Query().Get("sidx")
	if sidx == "" {
		sidx = "last_frame_time"
	}
	sord := r.URL.Query().Get("sord")
	if sord == "" {
		sord = "DESC"
	}

	// Filtros personalizados
	ruleName := r.URL.Query().Get("rule_name")
	imeiSearch := r.URL.Query().Get("imei")

	offset := (page - 1) * rows

	// Llamar a la nueva función paginada
	summary, totalRecords, err := audit.GetProgressSummaryPaginated(rows, offset, sidx, sord, ruleName, imeiSearch)
	if err != nil {
		log.Printf("Error obteniendo summary paginado: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	totalPages := 0
	if totalRecords > 0 {
		totalPages = (totalRecords + rows - 1) / rows
	}

	// Formatear respuesta para jqGrid
	response := map[string]interface{}{
		"page":    page,
		"total":   totalPages,
		"records": totalRecords,
		"rows":    summary, // jqGrid espera los datos en un campo 'rows'
	}

	json.NewEncoder(w).Encode(response)
}

// GET /api/audit/progress/timeline?imei=xxx - jqGrid Nivel 2: Timeline de frames (con paginación)
func progressTimelineHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Parámetros de jqGrid para paginación y ordenamiento
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	rows, _ := strconv.Atoi(r.URL.Query().Get("rows"))
	if rows < 1 {
		rows = 20 // Default para subgrid
	}
	sidx := r.URL.Query().Get("sidx")
	if sidx == "" {
		sidx = "execution_time"
	}
	sord := r.URL.Query().Get("sord")
	if sord == "" {
		sord = "ASC" // Por defecto ascendente para timeline
	}

	// Filtros personalizados
	imei := r.URL.Query().Get("imei")
	if imei == "" {
		http.Error(w, `{"error":"IMEI parameter required"}`, http.StatusBadRequest)
		return
	}
	ruleName := r.URL.Query().Get("rule_name")

	offset := (page - 1) * rows

	// Llamar a la nueva función paginada del timeline
	timeline, totalRecords, err := audit.GetFrameTimelinePaginated(imei, rows, offset, sidx, sord, ruleName)
	if err != nil {
		log.Printf("Error obteniendo timeline paginado para IMEI %s: %v", imei, err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	totalPages := 0
	if totalRecords > 0 {
		totalPages = (totalRecords + rows - 1) / rows
	}

	// Formatear respuesta para jqGrid
	response := map[string]interface{}{
		"page":    page,
		"total":   totalPages,
		"records": totalRecords,
		"rows":    timeline, // jqGrid espera los datos en un campo 'rows'
	}

	json.NewEncoder(w).Encode(response)
}

// GET /api/audit/progress/snapshot?id=xxx - Obtiene el snapshot de un frame específico
func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"ID parameter is required"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID parameter"}`, http.StatusBadRequest)
		return
	}

	snapshot, err := audit.GetSnapshotByID(id)
	if err != nil {
		log.Printf("Error getting snapshot for ID %d: %v", id, err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"id":       id,
		"snapshot": snapshot,
	})
}

// GET /api/rules/components?rule_name=xxx - Obtiene los componentes internos de una regla (paquete)
func ruleComponentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	ruleName := r.URL.Query().Get("rule_name")
	if ruleName == "" {
		http.Error(w, `{"error":"rule_name parameter is required"}`, http.StatusBadRequest)
		return
	}

	components, err := engine.GetInternalRuleNames(ruleName)
	if err != nil {
		log.Printf("Error getting components for rule %s: %v", ruleName, err)
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"rule_name":  ruleName,
		"components": components,
	})
}
