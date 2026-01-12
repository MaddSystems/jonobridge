package audit

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// RuleExecution representa UNA ejecución de regla capturada DURANTE la ejecución
// No post-mortem, no hardcoded - captura valores REALES
type RuleExecution struct {
	RuleName    string                 `json:"rule_name"`
	Description string                 `json:"description"`
	Salience    int64                  `json:"salience"`
	Status      string                 `json:"status"` // "PASSED", "FAILED", "SKIPPED"
	Timestamp   time.Time              `json:"timestamp"`
	DurationMs  float64                `json:"duration_ms"`
	Conditions  map[string]interface{} `json:"conditions"`  // Valores REALES evaluados
	Actions     []string               `json:"actions"`     // Acciones REALES ejecutadas
	AlertFired  bool                   `json:"alert_fired"` // Si esta regla disparó alerta
}

// IMEISummary representa el resumen de alertas de un IMEI
type IMEISummary struct {
	IMEI              string         `json:"imei"`
	LastAlertDate     time.Time      `json:"last_alert_date"`
	TotalAlerts24h    int            `json:"total_alerts_24h"`
	AlertTypes        map[string]int `json:"alert_types"` // {"SpeedAlert": 5, "JammerAlert": 2}
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

// Scan implementa sql.Scanner interface
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

// Value implementa driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// StringSlice es un wrapper para []string que implementa sql.Scanner y driver.Valuer
type StringSlice []string

// Scan implementa sql.Scanner interface
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

// Value implementa driver.Valuer interface
func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// ProgressAudit representa el estado de ejecución de una regla (en progreso o completada)
type ProgressAudit struct {
	ID                int64     `json:"id"`
	IMEI              string    `json:"imei"`
	RuleID            int64     `json:"rule_id"`            // ID de la regla en fleet_rules
	RuleName          string    `json:"rule_name"`          // Nombre principal de la regla (ej: "Jammer Real - Detección Avanzada")
	ComponentsExecuted []string `json:"components_executed"` // ["BufferUpdate", "CalculateMetrics", "SkipIfValid"]
	ComponentDetails  JSONMap   `json:"component_details"`  // Detalles de cada componente ejecutado
	StageReached      string    `json:"stage_reached"`      // "buffer_update", "metrics_calc", "expression_eval", "alert_fired"
	StopReason        string    `json:"stop_reason"`        // "waiting_for_10", "position_valid", "inside_geofence", "alert_triggered"
	BufferSize        int       `json:"buffer_size"`        // 7/10
	MetricsReady      bool      `json:"metrics_ready"`      // ¿Ya tiene métricas calculadas?
	GeofenceEval      string    `json:"geofence_eval"`      // "not_evaluated", "inside_taller", "outside_all"
	ContextSnapshot   JSONMap   `json:"context_snapshot"`   // Snapshot COMPLETO (buffer+metrics+geofences+flags+packet)
	ExecutionTime     time.Time `json:"execution_time"`
}
