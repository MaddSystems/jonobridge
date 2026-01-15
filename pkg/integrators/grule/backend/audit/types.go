package audit

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// RuleExecution representa UNA ejecución de regla capturada DURANTE la ejecución
type RuleExecution struct {
	RuleName    string                 `json:"rule_name"`
	Description string                 `json:"description"`
	Salience    int64                  `json:"salience"`
	Status      string                 `json:"status"` // "PASSED", "FAILED", "SKIPPED"
	Timestamp   time.Time              `json:"timestamp"`
	DurationMs  float64                `json:"duration_ms"`
	Conditions  map[string]interface{} `json:"conditions"`
	Actions     []string               `json:"actions"`
	AlertFired  bool                   `json:"alert_fired"`
}

// IMEISummary representa el resumen de alertas de un IMEI
type IMEISummary struct {
	IMEI              string         `json:"imei"`
	LastAlertDate     time.Time      `json:"last_alert_date"`
	TotalAlerts24h    int            `json:"total_alerts_24h"`
	AlertTypes        map[string]int `json:"alert_types"`
	LastRuleExecuted  string         `json:"last_rule_executed"`
	LastAlertLocation string         `json:"last_alert_location"`
}

// AlertDetail representa el detalle completo de una alerta específica
type AlertDetail struct {
	ID              int64                  `json:"id"`
	IMEI            string                 `json:"imei"`
	AlertDate       time.Time              `json:"alert_date"`
	RuleName        string                 `json:"rule_name"`
	RuleDescription string                 `json:"rule_description"`
	Salience        int64                  `json:"salience"`
	Conditions      map[string]interface{} `json:"conditions"`
	Actions         []string               `json:"actions"`
	TelegramSent    bool                   `json:"telegram_sent"`
	Latitude        float64                `json:"latitude"`
	Longitude       float64                `json:"longitude"`
	Speed           int64                  `json:"speed"`
}

// JSONMap es un wrapper para JSON que implementa sql.Scanner y driver.Valuer
type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// ProgressAudit representa el estado de ejecución de una regla

type ProgressAudit struct {

	ID                 int64     `json:"id"`

	IMEI               string    `json:"imei"`

	RuleID             int64     `json:"rule_id"`

	RuleName           string    `json:"rule_name"`

	ComponentsExecuted []string  `json:"components_executed"`

	ComponentDetails   JSONMap   `json:"component_details"`

	StepNumber         int       `json:"step_number"`

	StageReached       string    `json:"stage_reached"`

	Level              string    `json:"level"`

	IsPost             bool      `json:"is_post"`

	StopReason         string    `json:"stop_reason"`

	BufferSize         int       `json:"buffer_size"`

	MetricsReady       bool      `json:"metrics_ready"`

	GeofenceEval       string    `json:"geofence_eval"`

	ContextSnapshot    JSONMap   `json:"context_snapshot"`

	ExecutionTime      time.Time `json:"execution_time"`

}


