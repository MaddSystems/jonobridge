package actions

import (
	"fmt"
	"time"

	"github.com/jonobridge/grule-integrator/engine/audit"
)

// ActionsHelper expone funciones para llamar desde reglas GRL
// Incluye captura universal de auditoría con snapshot automático del contexto
type ActionsHelper struct {
	// Datos del contexto actual (se inyectan desde grule_worker.go)
	IMEI              string
	Speed             int64
	Latitude          float64
	Longitude         float64
	Altitude          int64
	GSMSignalStrength int64
	Satellites        int64
	PositioningStatus string
	EventCode         string
	Datetime          time.Time
}

// Audit registra la ejecución de una regla DURANTE su ejecución
// FUNCIÓN UNIVERSAL: Construye automáticamente el snapshot desde el contexto
// USO: actions.Audit("RuleName", "Description", salience, alertFired)
func (a *ActionsHelper) Audit(ruleName, description string, salience int64, alertFired bool) {
	capture := audit.GetCapture(a.IMEI)
	if capture == nil {
		return // No hay captura activa
	}

	// Construir snapshot automáticamente desde el contexto
	conditions := map[string]interface{}{
		"IMEI":              a.IMEI,
		"Speed":             a.Speed,
		"Latitude":          a.Latitude,
		"Longitude":         a.Longitude,
		"Altitude":          a.Altitude,
		"GSMSignalStrength": a.GSMSignalStrength,
		"Satellites":        a.Satellites,
		"PositioningStatus": a.PositioningStatus,
		"EventCode":         a.EventCode,
		"Datetime":          a.Datetime,
	}

	// Detectar acciones ejecutadas (se puede mejorar con un tracker)
	actions := []string{}
	if alertFired {
		actions = append(actions, "AlertTriggered")
	}

	exec := audit.RuleExecution{
		RuleName:    ruleName,
		Description: description,
		Salience:    salience,
		Status:      "PASSED",
		Timestamp:   time.Now(),
		DurationMs:  0,
		Conditions:  conditions,
		Actions:     actions,
		AlertFired:  alertFired,
	}

	capture.RecordExecution(exec)
}

// Wrappers para funciones existentes en alerts.go y commands.go
// Estas funciones permiten llamar las acciones desde el ActionsHelper

func (a *ActionsHelper) SendTelegram(message string) {
	SendTelegram(message)
}

func (a *ActionsHelper) SendEmail(subject, body string) {
	SendEmail(subject, body)
}

func (a *ActionsHelper) CutEngine(imei string) {
	CutEngine(imei)
}

func (a *ActionsHelper) RestoreEngine(imei string) {
	RestoreEngine(imei)
}

func (a *ActionsHelper) SendRawHex(imei, hexCommand string) {
	SendRawHex(imei, hexCommand)
}

func (a *ActionsHelper) Log(message string, args ...interface{}) {
	Log(message, args...)
}

// CastString convierte cualquier valor a string para uso en GRL
func (a *ActionsHelper) CastString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// Concat une dos strings (alternativa segura a + en algunas versiones de GRL)
func (a *ActionsHelper) Concat(s1, s2 string) string {
	return s1 + s2
}
