package audit

import (
	"log"
	"os"
	"sync"
	"time"
)

// ExecutionCapture gestiona la captura de ejecuciones de reglas DURANTE la ejecución
// Thread-safe para uso concurrente desde múltiples workers
type ExecutionCapture struct {
	mu         sync.RWMutex
	imei       string
	executions []RuleExecution
	alertFired bool
	startTime  time.Time
}

// globalCaptures mantiene las capturas activas por IMEI
var (
	globalCaptures = make(map[string]*ExecutionCapture)
	capturesMutex  sync.RWMutex
)

// StartCapture inicia una nueva captura para un IMEI
func StartCapture(imei string) *ExecutionCapture {
	capturesMutex.Lock()
	defer capturesMutex.Unlock()

	capture := &ExecutionCapture{
		imei:       imei,
		executions: []RuleExecution{},
		alertFired: false,
		startTime:  time.Now(),
	}

	globalCaptures[imei] = capture
	return capture
}

// RecordExecution registra una ejecución de regla DURANTE su ejecución
// Esta función se llama desde ActionsHelper.RecordExecution()
func (ec *ExecutionCapture) RecordExecution(exec RuleExecution) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	exec.Timestamp = time.Now()
	ec.executions = append(ec.executions, exec)

	if exec.AlertFired {
		ec.alertFired = true
	}

	log.Printf("📝 [AUDIT] Capturado: %s - %s (Alert: %v)", ec.imei, exec.RuleName, exec.AlertFired)
}

// FinishCapture finaliza la captura y guarda en base de datos
func FinishCapture(imei string) error {
	capturesMutex.Lock()
	capture, exists := globalCaptures[imei]
	delete(globalCaptures, imei)
	capturesMutex.Unlock()

	if !exists || capture == nil {
		return nil // No hay nada que guardar
	}

	// Verificar si auditoría está habilitada
	if os.Getenv("GRULE_AUDIT_ENABLED") != "Y" {
		return nil
	}

	// Aplicar filtro de nivel
	level := os.Getenv("GRULE_AUDIT_LEVEL")
	if level == "NONE" {
		return nil
	}
	if level == "ERROR" && !capture.alertFired {
		return nil // Solo guardar si hay alerta
	}

	// Si hay ejecuciones, guardar
	if len(capture.executions) > 0 {
		return SaveExecutions(imei, capture.executions)
	}

	return nil
}

// GetCapture obtiene la captura activa para un IMEI
func GetCapture(imei string) *ExecutionCapture {
	capturesMutex.RLock()
	defer capturesMutex.RUnlock()
	return globalCaptures[imei]
}

// RecordProgress registra el progreso de ejecución (para auditoría de progreso)
func RecordProgress(progress ProgressAudit) {
	if !IsProgressAuditEnabled() {
		return // No guardar si está desactivado
	}

	if err := SaveProgressAudit(progress); err != nil {
		log.Printf("❌ Error guardando progress audit: %v", err)
	}
}
