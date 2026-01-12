package audit

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

var db *sql.DB

func InitDB(dbConn *sql.DB) {
	db = dbConn
	createOptimizedTables()
}

func createOptimizedTables() {
	if db == nil {
		log.Println("⚠️  Audit database not initialized")
		return
	}

	schema1 := `
	CREATE TABLE IF NOT EXISTS alert_summary (
		imei VARCHAR(20) PRIMARY KEY,
		last_alert_date DATETIME(6) NOT NULL,
		total_alerts_24h INT DEFAULT 0,
		alert_types JSON,
		last_rule_executed VARCHAR(100),
		last_alert_location VARCHAR(100),
		updated_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
		INDEX idx_last_alert (last_alert_date),
		INDEX idx_total_alerts (total_alerts_24h)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	if _, err := db.Exec(schema1); err != nil {
		log.Printf("⚠️  Error creating alert_summary: %v", err)
	}

	schema2 := `
	CREATE TABLE IF NOT EXISTS alert_details (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		imei VARCHAR(20) NOT NULL,
		alert_date DATETIME(6) NOT NULL,
		rule_name VARCHAR(100) NOT NULL,
		rule_description VARCHAR(255),
		salience INT,
		conditions_snapshot JSON,
		actions_executed JSON,
		telegram_sent BOOLEAN DEFAULT false,
		latitude DECIMAL(10, 6),
		longitude DECIMAL(10, 6),
		speed INT,
		created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6),
		INDEX idx_imei_date (imei, alert_date),
		INDEX idx_rule_name (rule_name),
		INDEX idx_alert_date (alert_date)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	if _, err := db.Exec(schema2); err != nil {
		log.Printf("⚠️  Error creating alert_details: %v", err)
	}
}

func SaveExecutions(imei string, executions []RuleExecution) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	alertExecutions := []RuleExecution{}
	for _, exec := range executions {
		if exec.AlertFired {
			alertExecutions = append(alertExecutions, exec)
		}
	}

	if len(alertExecutions) == 0 {
		return nil
	}

	for _, exec := range alertExecutions {
		conditionsJSON, _ := json.Marshal(exec.Conditions)
		actionsJSON, _ := json.Marshal(exec.Actions)

		telegramSent := false
		for _, action := range exec.Actions {
			if action == "SendTelegram" {
				telegramSent = true
				break
			}
		}

		lat, _ := exec.Conditions["Latitude"].(float64)
		lon, _ := exec.Conditions["Longitude"].(float64)
		speed, _ := exec.Conditions["Speed"].(int)

		query := `
			INSERT INTO alert_details 
			(imei, alert_date, rule_name, rule_description, salience, 
			 conditions_snapshot, actions_executed, telegram_sent, 
			 latitude, longitude, speed)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		_, err := db.Exec(query,
			imei,
			exec.Timestamp,
			exec.RuleName,
			exec.Description,
			exec.Salience,
			string(conditionsJSON),
			string(actionsJSON),
			telegramSent,
			lat,
			lon,
			speed,
		)

		if err != nil {
			log.Printf("❌ Error insertando alert_detail: %v", err)
			continue
		}
	}

	return updateSummary(imei, alertExecutions)
}

func updateSummary(imei string, alertExecutions []RuleExecution) error {
	if len(alertExecutions) == 0 {
		return nil
	}

	lastAlert := alertExecutions[len(alertExecutions)-1]

	alertTypes := make(map[string]int)
	query := `
		SELECT rule_name, COUNT(*) as count
		FROM alert_details
		WHERE imei = ? AND alert_date >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
		GROUP BY rule_name
	`

	rows, err := db.Query(query, imei)
	if err != nil {
		log.Printf("❌ Error contando alertas: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var ruleName string
			var count int
			if err := rows.Scan(&ruleName, &count); err == nil {
				alertTypes[ruleName] = count
			}
		}
	}

	alertTypesJSON, _ := json.Marshal(alertTypes)
	totalAlerts := 0
	for _, count := range alertTypes {
		totalAlerts += count
	}

	lat, _ := lastAlert.Conditions["Latitude"].(float64)
	lon, _ := lastAlert.Conditions["Longitude"].(float64)
	location := fmt.Sprintf("%.6f,%.6f", lat, lon)

	upsertQuery := `
		INSERT INTO alert_summary 
		(imei, last_alert_date, total_alerts_24h, alert_types, last_rule_executed, last_alert_location)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			last_alert_date = VALUES(last_alert_date),
			total_alerts_24h = VALUES(total_alerts_24h),
			alert_types = VALUES(alert_types),
			last_rule_executed = VALUES(last_rule_executed),
			last_alert_location = VALUES(last_alert_location)
	`

	_, err = db.Exec(upsertQuery,
		imei,
		lastAlert.Timestamp,
		totalAlerts,
		string(alertTypesJSON),
		lastAlert.RuleName,
		location,
	)

	return err
}

// Progress Audit DB Functions

var progressAuditEnabled bool = false

func CreateProgressAuditTable() {
	if db == nil {
		log.Println("⚠️  Progress audit database not initialized")
		return
	}

	schema := `
	CREATE TABLE IF NOT EXISTS rule_execution_state (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		imei VARCHAR(20) NOT NULL,
		rule_id BIGINT NOT NULL,
		rule_name VARCHAR(100) NOT NULL,
		components_executed JSON,
		component_details JSON,
		step_number INT DEFAULT 0,
		stage_reached VARCHAR(50) NOT NULL,
		level VARCHAR(20) DEFAULT 'info',
		stop_reason VARCHAR(100) NOT NULL,
		buffer_size INT DEFAULT 0,
		metrics_ready BOOLEAN DEFAULT false,
		geofence_eval VARCHAR(50) DEFAULT 'not_evaluated',
		context_snapshot JSON,
		execution_time DATETIME(6) NOT NULL,
		INDEX idx_imei_time (imei, execution_time),
		INDEX idx_imei_step (imei, step_number),
		INDEX idx_rule_stage (rule_name, stage_reached),
		INDEX idx_execution_time (execution_time),
		INDEX idx_rule_id (rule_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	db.Exec(schema)
}

func EnableProgressAudit() error {
	CreateProgressAuditTable()
	progressAuditEnabled = true
	log.Println("✅ Progress Audit ENABLED")
	return nil
}

func DisableProgressAudit() error {
	progressAuditEnabled = false
	return nil
}

func IsProgressAuditEnabled() bool {
	return progressAuditEnabled
}

func SaveProgressAudit(progress ProgressAudit) error {
	if !progressAuditEnabled || db == nil {
		return nil
	}

	contextSnapshotJSON, _ := json.Marshal(progress.ContextSnapshot)
	componentsExecutedJSON, _ := json.Marshal(progress.ComponentsExecuted)
	componentDetailsJSON, _ := json.Marshal(progress.ComponentDetails)

	query := `
		INSERT INTO rule_execution_state 
		(imei, rule_id, rule_name, components_executed, component_details, step_number, stage_reached, level, stop_reason, buffer_size, 
		 metrics_ready, geofence_eval, context_snapshot, execution_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query,
		progress.IMEI,
		progress.RuleID,
		progress.RuleName,
		string(componentsExecutedJSON),
		string(componentDetailsJSON),
		progress.StepNumber,
		progress.StageReached,
		progress.Level,
		progress.StopReason,
		progress.BufferSize,
		progress.MetricsReady,
		progress.GeofenceEval,
		string(contextSnapshotJSON),
		progress.ExecutionTime,
	)
	return err
}

func SaveProgressFrame(imei string, ruleID int64, ruleName string, componentsExecuted []string, componentDetails map[string]interface{},
	stepNumber int, stageReached, level, stopReason string, bufferSize int, metricsReady bool, geofenceEval string,
	contextSnapshot map[string]interface{}) error {

	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	snapshotJSON, _ := json.Marshal(contextSnapshot)
	componentsJSON, _ := json.Marshal(componentsExecuted)
	detailsJSON, _ := json.Marshal(componentDetails)

	query := `
		INSERT INTO rule_execution_state 
		(imei, rule_id, rule_name, components_executed, component_details, step_number, stage_reached, level, stop_reason, 
		 buffer_size, metrics_ready, geofence_eval, context_snapshot, execution_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(6))
	`

	_, err := db.Exec(query, imei, ruleID, ruleName, string(componentsJSON), string(detailsJSON), stepNumber, stageReached, level, stopReason,
		bufferSize, metricsReady, geofenceEval, string(snapshotJSON))

	return err
}

func GetAvailableRules() ([]map[string]interface{}, error) {
	if db == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	query := `
		SELECT DISTINCT rule_name, COUNT(*) as total_frames, COUNT(DISTINCT imei) as total_imeis, MAX(execution_time) as last_execution
		FROM rule_execution_state GROUP BY rule_name ORDER BY last_execution DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var ruleName string
		var totalFrames, totalImeis int
		var lastExecution string
		rows.Scan(&ruleName, &totalFrames, &totalImeis, &lastExecution)
		result = append(result, map[string]interface{}{
			"rule_name":      ruleName,
			"total_frames":   totalFrames,
			"total_imeis":    totalImeis,
			"last_execution": lastExecution,
		})
	}
	return result, nil
}

func GetProgressSummaryPaginated(limit, offset int, sortBy, sortOrder, ruleName, imeiSearch string) ([]map[string]interface{}, int, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("db not initialized")
	}

	// Map frontend column names to actual SQL columns/aggregates
	columnMap := map[string]string{
		"last_frame_time": "MAX(execution_time)",
		"max_step":        "MAX(step_number)",
		"total_frames":    "COUNT(*)",
		"imei":            "imei",
		"rule_name":       "rule_name",
	}

	// Get the actual SQL expression for sorting
	sortExpr, ok := columnMap[sortBy]
	if !ok {
		sortExpr = "MAX(execution_time)" // default
	}

	whereClause := ""
	args := []interface{}{}
	if ruleName != "" {
		whereClause = "WHERE rule_name = ?"
		args = append(args, ruleName)
	}
	if imeiSearch != "" {
		if whereClause == "" {
			whereClause = "WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "imei LIKE ?"
		args = append(args, "%"+imeiSearch+"%")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (SELECT 1 FROM rule_execution_state %s GROUP BY imei, rule_name) AS sub", whereClause)
	var total int
	db.QueryRow(countQuery, args...).Scan(&total)

	query := fmt.Sprintf(`
		SELECT imei, rule_name, MAX(step_number), COUNT(*), MAX(execution_time)
		FROM rule_execution_state %s GROUP BY imei, rule_name ORDER BY %s %s LIMIT ? OFFSET ?
	`, whereClause, sortExpr, sortOrder)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var imei, rName, lastTime string
		var maxStep, totalFrames int
		rows.Scan(&imei, &rName, &maxStep, &totalFrames, &lastTime)
		result = append(result, map[string]interface{}{
			"imei": imei, "rule_name": rName, "max_step": maxStep, "total_frames": totalFrames, "last_frame_time": lastTime,
		})
	}
	return result, total, nil
}

func GetFrameTimelinePaginated(imei string, limit, offset int, sortBy, sortOrder, ruleName string) ([]map[string]interface{}, int, error) {
	if db == nil {
		return nil, 0, fmt.Errorf("db not initialized")
	}

	whereClause := "WHERE imei = ?"
	args := []interface{}{imei}
	if ruleName != "" {
		whereClause += " AND rule_name = ?"
		args = append(args, ruleName)
	}

	var total int
	db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM rule_execution_state %s", whereClause), args...).Scan(&total)

	query := fmt.Sprintf(`
		SELECT id, rule_id, rule_name, components_executed, component_details, step_number, stage_reached, level, stop_reason,
		buffer_size, metrics_ready, geofence_eval, context_snapshot, execution_time
		FROM rule_execution_state %s ORDER BY %s %s LIMIT ? OFFSET ?
	`, whereClause, sortBy, sortOrder)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result := []map[string]interface{}{}
	for rows.Next() {
		var id, step, bufSize int
		var rId int64
		var rName, stage, level, stop, geo string
		var metrics bool
		var snapJSON, compJSON, detJSON []byte
		var execTime sql.NullTime

		rows.Scan(&id, &rId, &rName, &compJSON, &detJSON, &step, &stage, &level, &stop, &bufSize, &metrics, &geo, &snapJSON, &execTime)

		var snap, det map[string]interface{}
		var comp []string
		json.Unmarshal(snapJSON, &snap)
		json.Unmarshal(compJSON, &comp)
		json.Unmarshal(detJSON, &det)

		frame := map[string]interface{}{
			"id": id, "rule_id": rId, "rule_name": rName, "components_executed": comp, "component_details": det,
			"step_number": step, "stage_reached": stage, "level": level, "stop_reason": stop, "buffer_size": bufSize,
			"metrics_ready": metrics, "geofence_eval": geo, "snapshot": snap,
		}
		if execTime.Valid {
			frame["execution_time"] = execTime.Time
		}
		result = append(result, frame)
	}
	return result, total, nil
}

func GetSnapshotByID(id int) (map[string]interface{}, error) {
	if db == nil {
		return nil, fmt.Errorf("db not initialized")
	}
	var jsonBytes []byte
	err := db.QueryRow("SELECT context_snapshot FROM rule_execution_state WHERE id = ?", id).Scan(&jsonBytes)
	if err != nil {
		return nil, err
	}
	var snap map[string]interface{}
	json.Unmarshal(jsonBytes, &snap)
	return snap, nil
}
