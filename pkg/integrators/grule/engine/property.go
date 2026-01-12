package engine

import (
	"context"
	"log"
	"sync"
)

// Property representa un contenedor de propiedades de evento con sincronización via channels
// Usado para comunicar cambios de estado entre el contexto principal y las goroutines
type Property struct {
	mu                                   sync.RWMutex
	DebugProcessed                       bool
	ResetProcessed                       bool
	PositionInvalidDetectedProcessed     bool
	PositionInvalidDetectedFailed        bool
	PositionInvalidDetected              bool
	MovingWithWeakSignalProcessed        bool
	MovingWithWeakSignalFailed           bool
	MovingWithWeakSignal                 bool
	OutsideAllSafeZonesProcessed         bool
	OutsideAllSafeZonesFailed            bool
	OutsideAllSafeZones                  bool
	JammerPatternFullyConfirmedProcessed bool
	JammerPatternFullyConfirmedFailed    bool
	JammerPatternFullyConfirmed          bool
	BufferUpdated                        bool
	BufferHas10                          bool
	MetricsReady                         bool
	CurrentlyInvalid                     bool
	EvaluationSkipped                    bool
	AlertFired                           bool
	ctx                                  context.Context
	cancel                               context.CancelFunc
	updatedChannels                      map[string]chan bool // Channels para cada propiedad
	gID                                  string               // Goroutine ID para debugging
}

// NewProperty crea una nueva instancia de Property con contexto y channels
func NewProperty(gID string, ctx context.Context) *Property {
	propCtx, cancel := context.WithCancel(ctx)
	prop := &Property{
		ctx:             propCtx,
		cancel:          cancel,
		updatedChannels: make(map[string]chan bool),
		gID:             gID,
	}

	// Inicializar channels para cada condición y flag de ejecución
	prop.updatedChannels["DebugProcessed"] = make(chan bool, 10)
	prop.updatedChannels["ResetProcessed"] = make(chan bool, 10)
	prop.updatedChannels["PositionInvalidDetectedProcessed"] = make(chan bool, 10)
	prop.updatedChannels["PositionInvalidDetectedFailed"] = make(chan bool, 10)
	prop.updatedChannels["PositionInvalidDetected"] = make(chan bool, 10)
	prop.updatedChannels["MovingWithWeakSignalProcessed"] = make(chan bool, 10)
	prop.updatedChannels["MovingWithWeakSignalFailed"] = make(chan bool, 10)
	prop.updatedChannels["MovingWithWeakSignal"] = make(chan bool, 10)
	prop.updatedChannels["OutsideAllSafeZonesProcessed"] = make(chan bool, 10)
	prop.updatedChannels["OutsideAllSafeZonesFailed"] = make(chan bool, 10)
	prop.updatedChannels["OutsideAllSafeZones"] = make(chan bool, 10)
	prop.updatedChannels["JammerPatternFullyConfirmedProcessed"] = make(chan bool, 10)
	prop.updatedChannels["JammerPatternFullyConfirmedFailed"] = make(chan bool, 10)
	prop.updatedChannels["JammerPatternFullyConfirmed"] = make(chan bool, 10)
	prop.updatedChannels["BufferUpdated"] = make(chan bool, 10)
	prop.updatedChannels["BufferHas10"] = make(chan bool, 10)
	prop.updatedChannels["MetricsReady"] = make(chan bool, 10)
	prop.updatedChannels["CurrentlyInvalid"] = make(chan bool, 10)
	prop.updatedChannels["EvaluationSkipped"] = make(chan bool, 10)
	prop.updatedChannels["AlertFired"] = make(chan bool, 10)

	return prop
}

// SetPositionInvalidDetected actualiza PositionInvalidDetected con sincronización via channel
func (p *Property) SetPositionInvalidDetected(value bool) {
	p.mu.Lock()
	oldValue := p.PositionInvalidDetected
	p.PositionInvalidDetected = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] PositionInvalidDetected cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("PositionInvalidDetected", value)
	}
}

// SetMovingWithWeakSignal actualiza MovingWithWeakSignal con sincronización via channel
func (p *Property) SetMovingWithWeakSignal(value bool) {
	p.mu.Lock()
	oldValue := p.MovingWithWeakSignal
	p.MovingWithWeakSignal = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] MovingWithWeakSignal cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("MovingWithWeakSignal", value)
	}
}

// SetOutsideAllSafeZones actualiza OutsideAllSafeZones con sincronización via channel
func (p *Property) SetOutsideAllSafeZones(value bool) {
	p.mu.Lock()
	oldValue := p.OutsideAllSafeZones
	p.OutsideAllSafeZones = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] OutsideAllSafeZones cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("OutsideAllSafeZones", value)
	}
}

// SetJammerPatternFullyConfirmed actualiza JammerPatternFullyConfirmed con sincronización via channel
func (p *Property) SetJammerPatternFullyConfirmed(value bool) {
	p.mu.Lock()
	oldValue := p.JammerPatternFullyConfirmed
	p.JammerPatternFullyConfirmed = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] JammerPatternFullyConfirmed cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("JammerPatternFullyConfirmed", value)
	}
}

// notifyUpdate envía el cambio por el channel correspondiente
func (p *Property) notifyUpdate(propName string, value bool) {
	ch, exists := p.updatedChannels[propName]
	if !exists {
		return
	}

	select {
	case ch <- value:
		log.Printf("[GID:%s] ✅ [PROPERTY] Notificación de %s enviada: %v", p.gID, propName, value)
	case <-p.ctx.Done():
		log.Printf("[GID:%s] ⚠️  [PROPERTY] Context cancelado, no se puede enviar notificación", p.gID)
		return
	default:
		log.Printf("[GID:%s] ⚠️  [PROPERTY] Canal lleno para %s, descartando notificación", p.gID, propName)
	}
}

// GetPositionInvalidDetected obtiene PositionInvalidDetected de forma thread-safe
func (p *Property) GetPositionInvalidDetected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.PositionInvalidDetected
}

// GetMovingWithWeakSignal obtiene MovingWithWeakSignal de forma thread-safe
func (p *Property) GetMovingWithWeakSignal() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.MovingWithWeakSignal
}

// GetOutsideAllSafeZones obtiene OutsideAllSafeZones de forma thread-safe
func (p *Property) GetOutsideAllSafeZones() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.OutsideAllSafeZones
}

// GetJammerPatternFullyConfirmed obtiene JammerPatternFullyConfirmed de forma thread-safe
func (p *Property) GetJammerPatternFullyConfirmed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.JammerPatternFullyConfirmed
}

// Close cancela el contexto y cierra todos los channels
func (p *Property) Close() {
	p.cancel()
	for _, ch := range p.updatedChannels {
		close(ch)
	}
}

// WatchCondition permite escuchar cambios en una condición específica
func (p *Property) WatchCondition(condName string) <-chan bool {
	ch, exists := p.updatedChannels[condName]
	if !exists {
		return nil
	}
	return ch
}

// ===== EXECUTION GUARD SETTERS/GETTERS =====

// SetDebugProcessed actualiza DebugProcessed
func (p *Property) SetDebugProcessed(value bool) {
	p.mu.Lock()
	oldValue := p.DebugProcessed
	p.DebugProcessed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("DebugProcessed", value)
	}
}

func (p *Property) GetDebugProcessed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.DebugProcessed
}

// SetResetProcessed actualiza ResetProcessed
func (p *Property) SetResetProcessed(value bool) {
	p.mu.Lock()
	oldValue := p.ResetProcessed
	p.ResetProcessed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("ResetProcessed", value)
	}
}

func (p *Property) GetResetProcessed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ResetProcessed
}

// ===== CONDITION 1 SETTERS/GETTERS =====

func (p *Property) SetPositionInvalidDetectedProcessed(value bool) {
	p.mu.Lock()
	oldValue := p.PositionInvalidDetectedProcessed
	p.PositionInvalidDetectedProcessed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("PositionInvalidDetectedProcessed", value)
	}
}

func (p *Property) GetPositionInvalidDetectedProcessed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.PositionInvalidDetectedProcessed
}

func (p *Property) SetPositionInvalidDetectedFailed(value bool) {
	p.mu.Lock()
	oldValue := p.PositionInvalidDetectedFailed
	p.PositionInvalidDetectedFailed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("PositionInvalidDetectedFailed", value)
	}
}

func (p *Property) GetPositionInvalidDetectedFailed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.PositionInvalidDetectedFailed
}

// ===== CONDITION 2 SETTERS/GETTERS =====

func (p *Property) SetMovingWithWeakSignalProcessed(value bool) {
	p.mu.Lock()
	oldValue := p.MovingWithWeakSignalProcessed
	p.MovingWithWeakSignalProcessed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("MovingWithWeakSignalProcessed", value)
	}
}

func (p *Property) GetMovingWithWeakSignalProcessed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.MovingWithWeakSignalProcessed
}

func (p *Property) SetMovingWithWeakSignalFailed(value bool) {
	p.mu.Lock()
	oldValue := p.MovingWithWeakSignalFailed
	p.MovingWithWeakSignalFailed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("MovingWithWeakSignalFailed", value)
	}
}

func (p *Property) GetMovingWithWeakSignalFailed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.MovingWithWeakSignalFailed
}

// ===== CONDITION 3 SETTERS/GETTERS =====

func (p *Property) SetOutsideAllSafeZonesProcessed(value bool) {
	p.mu.Lock()
	oldValue := p.OutsideAllSafeZonesProcessed
	p.OutsideAllSafeZonesProcessed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("OutsideAllSafeZonesProcessed", value)
	}
}

func (p *Property) GetOutsideAllSafeZonesProcessed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.OutsideAllSafeZonesProcessed
}

func (p *Property) SetOutsideAllSafeZonesFailed(value bool) {
	p.mu.Lock()
	oldValue := p.OutsideAllSafeZonesFailed
	p.OutsideAllSafeZonesFailed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("OutsideAllSafeZonesFailed", value)
	}
}

func (p *Property) GetOutsideAllSafeZonesFailed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.OutsideAllSafeZonesFailed
}

// ===== CONDITION 4 SETTERS/GETTERS =====

func (p *Property) SetJammerPatternFullyConfirmedProcessed(value bool) {
	p.mu.Lock()
	oldValue := p.JammerPatternFullyConfirmedProcessed
	p.JammerPatternFullyConfirmedProcessed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("JammerPatternFullyConfirmedProcessed", value)
	}
}

func (p *Property) GetJammerPatternFullyConfirmedProcessed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.JammerPatternFullyConfirmedProcessed
}

func (p *Property) SetJammerPatternFullyConfirmedFailed(value bool) {
	p.mu.Lock()
	oldValue := p.JammerPatternFullyConfirmedFailed
	p.JammerPatternFullyConfirmedFailed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("JammerPatternFullyConfirmedFailed", value)
	}
}

func (p *Property) GetJammerPatternFullyConfirmedFailed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.JammerPatternFullyConfirmedFailed
}

// SetBufferUpdated actualiza BufferUpdated con sincronización via channel
func (p *Property) SetBufferUpdated(value bool) {
	p.mu.Lock()
	oldValue := p.BufferUpdated
	p.BufferUpdated = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] BufferUpdated cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("BufferUpdated", value)
	}
}

func (p *Property) GetBufferUpdated() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.BufferUpdated
}

// SetBufferHas10 actualiza BufferHas10 con sincronización via channel
func (p *Property) SetBufferHas10(value bool) {
	p.mu.Lock()
	oldValue := p.BufferHas10
	p.BufferHas10 = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] BufferHas10 cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("BufferHas10", value)
	}
}

func (p *Property) GetBufferHas10() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.BufferHas10
}

// SetMetricsReady actualiza MetricsReady con sincronización via channel
func (p *Property) SetMetricsReady(value bool) {
	p.mu.Lock()
	oldValue := p.MetricsReady
	p.MetricsReady = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] MetricsReady cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("MetricsReady", value)
	}
}

func (p *Property) GetMetricsReady() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.MetricsReady
}

// SetCurrentlyInvalid actualiza CurrentlyInvalid con sincronización via channel
func (p *Property) SetCurrentlyInvalid(value bool) {
	p.mu.Lock()
	oldValue := p.CurrentlyInvalid
	p.CurrentlyInvalid = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] CurrentlyInvalid cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("CurrentlyInvalid", value)
	}
}

func (p *Property) GetCurrentlyInvalid() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.CurrentlyInvalid
}

// SetEvaluationSkipped actualiza EvaluationSkipped con sincronización via channel
func (p *Property) SetEvaluationSkipped(value bool) {
	p.mu.Lock()
	oldValue := p.EvaluationSkipped
	p.EvaluationSkipped = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] EvaluationSkipped cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("EvaluationSkipped", value)
	}
}

func (p *Property) GetEvaluationSkipped() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.EvaluationSkipped
}

// SetAlertFired actualiza AlertFired con sincronización via channel
func (p *Property) SetAlertFired(value bool) {
	p.mu.Lock()
	oldValue := p.AlertFired
	p.AlertFired = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] AlertFired cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("AlertFired", value)
	}
}

func (p *Property) GetAlertFired() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.AlertFired
}