package audit

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

var alertDB *sql.DB

// InitAlertAuditDB initializes the database connection for alert audit
func InitAlertAuditDB(database *sql.DB) {
	alertDB = database
}

// GetIMEISummaries returns summary of alerts for all IMEIs
func GetIMEISummaries(limit int) ([]IMEISummary, error) {
	if alertDB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `
		SELECT imei, last_alert_date, total_alerts_24h, alert_types, 
		       last_rule_executed, last_alert_location
		FROM alert_summary
		ORDER BY last_alert_date DESC
		LIMIT ?
	`

	rows, err := alertDB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []IMEISummary
	for rows.Next() {
		var s IMEISummary
		var alertTypesJSON string

		err := rows.Scan(
			&s.IMEI,
			&s.LastAlertDate,
			&s.TotalAlerts24h,
			&alertTypesJSON,
			&s.LastRuleExecuted,
			&s.LastAlertLocation,
		)

		if err != nil {
			log.Printf("Error scanning summary: %v", err)
			continue
		}

		json.Unmarshal([]byte(alertTypesJSON), &s.AlertTypes)
		summaries = append(summaries, s)
	}

	return summaries, nil
}

// GetIMEISummariesPaginated returns paginated alert summary with search and sorting
func GetIMEISummariesPaginated(limit, offset int, sortBy, sortOrder, searchText string) ([]IMEISummary, int, error) {
	if alertDB == nil {
		return nil, 0, fmt.Errorf("database not initialized")
	}

	// Validate allowed sort columns
	allowedSortColumns := map[string]bool{
		"imei":               true,
		"last_alert_date":    true,
		"total_alerts_24h":   true,
		"last_rule_executed": true,
	}

	if !allowedSortColumns[sortBy] {
		sortBy = "last_alert_date"
	}

	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	// Build WHERE clause
	whereClause := ""
	args := []interface{}{}

	if searchText != "" {
		whereClause = "WHERE imei LIKE ? OR last_rule_executed LIKE ?"
		searchPattern := "%" + searchText + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM alert_summary %s", whereClause)
	var total int
	err := alertDB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated records
	query := fmt.Sprintf(`
		SELECT imei, last_alert_date, total_alerts_24h, alert_types, 
		       last_rule_executed, last_alert_location
		FROM alert_summary
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereClause, sortBy, sortOrder)

	args = append(args, limit, offset)

	rows, err := alertDB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var summaries []IMEISummary
	for rows.Next() {
		var s IMEISummary
		var alertTypesJSON string

		err := rows.Scan(
			&s.IMEI,
			&s.LastAlertDate,
			&s.TotalAlerts24h,
			&alertTypesJSON,
			&s.LastRuleExecuted,
			&s.LastAlertLocation,
		)

		if err != nil {
			log.Printf("Error scanning summary: %v", err)
			continue
		}

		json.Unmarshal([]byte(alertTypesJSON), &s.AlertTypes)
		summaries = append(summaries, s)
	}

	return summaries, total, nil
}

// GetAlertDetails returns detailed alert information for a specific IMEI
func GetAlertDetails(imei string, limit int) ([]AlertDetail, error) {
	if alertDB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `
		SELECT id, imei, alert_date, rule_name, rule_description, salience,
		       conditions_snapshot, actions_executed, telegram_sent,
		       latitude, longitude, speed
		FROM alert_details
		WHERE imei = ?
		ORDER BY alert_date DESC
		LIMIT ?
	`

	rows, err := alertDB.Query(query, imei, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []AlertDetail
	for rows.Next() {
		var d AlertDetail
		var conditionsJSON, actionsJSON string

		err := rows.Scan(
			&d.ID,
			&d.IMEI,
			&d.AlertDate,
			&d.RuleName,
			&d.RuleDescription,
			&d.Salience,
			&conditionsJSON,
			&actionsJSON,
			&d.TelegramSent,
			&d.Latitude,
			&d.Longitude,
			&d.Speed,
		)

		if err != nil {
			log.Printf("Error scanning detail: %v", err)
			continue
		}

		json.Unmarshal([]byte(conditionsJSON), &d.Conditions)
		json.Unmarshal([]byte(actionsJSON), &d.Actions)
		details = append(details, d)
	}

	return details, nil
}

// SaveAlertAudit saves alert information to alert_details and alert_summary tables
// This is called from Capture() when a rule with is_alert: true executes
func SaveAlertAudit(entry *AuditEntry) error {
	if alertDB == nil {
		return fmt.Errorf("alert database not initialized")
	}

	// Extract packet data from snapshot
	packet := extractPacketDataFromSnapshot(entry.Snapshot)

	// Prepare actions array
	actionsJSON, _ := json.Marshal([]string{"SendTelegram"}) // Default action for alerts

	// Prepare conditions snapshot from the rich context
	conditionsJSON, _ := json.Marshal(entry.Snapshot)

	// Insert into alert_details
	insertQuery := `
		INSERT INTO alert_details 
		(imei, alert_date, rule_name, rule_description, salience,
		 conditions_snapshot, actions_executed, telegram_sent,
		 latitude, longitude, speed)
		VALUES (?, NOW(6), ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := alertDB.Exec(
		insertQuery,
		entry.IMEI,
		entry.RuleName,
		entry.Description,
		entry.Salience,
		conditionsJSON,
		actionsJSON,
		true, // telegram_sent = true (alert was fired)
		packet.Latitude,
		packet.Longitude,
		packet.Speed,
	)

	if err != nil {
		log.Printf("❌ [AlertAudit] Error inserting alert_details: %v", err)
		return err
	}

	// Update or insert alert_summary
	if err := updateAlertSummary(entry.IMEI, entry.RuleName, packet); err != nil {
		log.Printf("⚠️ [AlertAudit] Error updating alert_summary: %v", err)
		// Don't fail if summary update fails
	}

	log.Printf("✅ [AlertAudit] Saved alert for IMEI %s, Rule %s", entry.IMEI, entry.RuleName)
	return nil
}

// updateAlertSummary updates or inserts the summary record for an IMEI
func updateAlertSummary(imei, ruleName string, packet PacketData) error {
	location := fmt.Sprintf("%.6f,%.6f", packet.Latitude, packet.Longitude)

	// Check if summary exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM alert_summary WHERE imei = ?)`
	err := alertDB.QueryRow(checkQuery, imei).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// Update existing record
		// Get current alert_types JSON
		var alertTypesJSON string
		getQuery := `SELECT alert_types FROM alert_summary WHERE imei = ?`
		err := alertDB.QueryRow(getQuery, imei).Scan(&alertTypesJSON)
		if err != nil {
			return err
		}

		// Parse existing alert types
		var alertTypes map[string]int
		if alertTypesJSON != "" && alertTypesJSON != "null" {
			json.Unmarshal([]byte(alertTypesJSON), &alertTypes)
		}
		if alertTypes == nil {
			alertTypes = make(map[string]int)
		}

		// Increment count for this rule type
		alertTypes[ruleName]++

		// Marshal back to JSON
		updatedTypesJSON, _ := json.Marshal(alertTypes)

		// Update record
		updateQuery := `
			UPDATE alert_summary 
			SET last_alert_date = NOW(6),
			    total_alerts_24h = total_alerts_24h + 1,
			    alert_types = ?,
			    last_rule_executed = ?,
			    last_alert_location = ?
			WHERE imei = ?
		`
		_, err = alertDB.Exec(updateQuery, updatedTypesJSON, ruleName, location, imei)
		return err
	} else {
		// Insert new record
		alertTypes := map[string]int{ruleName: 1}
		alertTypesJSON, _ := json.Marshal(alertTypes)

		insertQuery := `
			INSERT INTO alert_summary 
			(imei, last_alert_date, total_alerts_24h, alert_types, 
			 last_rule_executed, last_alert_location)
			VALUES (?, NOW(6), 1, ?, ?, ?)
		`
		_, err = alertDB.Exec(insertQuery, imei, alertTypesJSON, ruleName, location)
		return err
	}
}

// PacketData holds extracted packet information for alert logging
type PacketData struct {
	Latitude  float64
	Longitude float64
	Speed     int64
}

// extractPacketDataFromSnapshot extracts packet fields from the snapshot map
func extractPacketDataFromSnapshot(snapshot map[string]interface{}) PacketData {
	packet := PacketData{}

	// Try to get packet_current from snapshot
	if packetCurrent, ok := snapshot["packet_current"].(map[string]interface{}); ok {
		if lat, ok := packetCurrent["Latitude"].(float64); ok {
			packet.Latitude = lat
		}
		if lon, ok := packetCurrent["Longitude"].(float64); ok {
			packet.Longitude = lon
		}
		if speed, ok := packetCurrent["Speed"].(float64); ok {
			packet.Speed = int64(speed)
		} else if speed, ok := packetCurrent["Speed"].(int64); ok {
			packet.Speed = speed
		}
	}

	return packet
}
