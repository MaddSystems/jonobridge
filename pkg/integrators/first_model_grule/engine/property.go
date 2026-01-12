package engine

import (
	"context"
	"log"
	"sync"
)

// Property representa un contenedor de propiedades de evento con sincronización via channels
// Usado para comunicar cambios de estado entre el contexto principal y las goroutines
type Property struct {
	mu              sync.RWMutex
	DebugProcessed  bool
	ResetProcessed  bool
	Cond1Processed  bool
	Cond1Failed     bool
	Cond1Passed     bool
	Cond2Processed  bool
	Cond2Failed     bool
	Cond2Passed     bool
	Cond3Processed  bool
	Cond3Failed     bool
	Cond3Passed     bool
	Cond4Processed  bool
	Cond4Failed     bool
	Cond4Passed     bool
	Cond5Processed  bool
	Cond5Failed     bool
	Cond5Passed     bool
	BufferUpdated   bool
	BufferHas10     bool
	MetricsReady    bool
	CurrentlyInvalid bool
	EvaluationSkipped bool
	AlertFired      bool
	ctx             context.Context
	cancel          context.CancelFunc
	updatedChannels map[string]chan bool // Channels para cada propiedad
	gID             string               // Goroutine ID para debugging
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
	prop.updatedChannels["Cond1Processed"] = make(chan bool, 10)
	prop.updatedChannels["Cond1Failed"] = make(chan bool, 10)
	prop.updatedChannels["Cond1Passed"] = make(chan bool, 10)
	prop.updatedChannels["Cond2Processed"] = make(chan bool, 10)
	prop.updatedChannels["Cond2Failed"] = make(chan bool, 10)
	prop.updatedChannels["Cond2Passed"] = make(chan bool, 10)
	prop.updatedChannels["Cond3Processed"] = make(chan bool, 10)
	prop.updatedChannels["Cond3Failed"] = make(chan bool, 10)
	prop.updatedChannels["Cond3Passed"] = make(chan bool, 10)
	prop.updatedChannels["Cond4Processed"] = make(chan bool, 10)
	prop.updatedChannels["Cond4Failed"] = make(chan bool, 10)
	prop.updatedChannels["Cond4Passed"] = make(chan bool, 10)
	prop.updatedChannels["Cond5Processed"] = make(chan bool, 10)
	prop.updatedChannels["Cond5Failed"] = make(chan bool, 10)
	prop.updatedChannels["Cond5Passed"] = make(chan bool, 10)
	prop.updatedChannels["BufferUpdated"] = make(chan bool, 10)
	prop.updatedChannels["BufferHas10"] = make(chan bool, 10)
	prop.updatedChannels["MetricsReady"] = make(chan bool, 10)
	prop.updatedChannels["CurrentlyInvalid"] = make(chan bool, 10)
	prop.updatedChannels["EvaluationSkipped"] = make(chan bool, 10)
	prop.updatedChannels["AlertFired"] = make(chan bool, 10)

	return prop
}

// SetCond1Passed actualiza Cond1Passed con sincronización via channel
func (p *Property) SetCond1Passed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond1Passed
	p.Cond1Passed = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] Cond1Passed cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("Cond1Passed", value)
	}
}

// SetCond2Passed actualiza Cond2Passed con sincronización via channel
func (p *Property) SetCond2Passed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond2Passed
	p.Cond2Passed = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] Cond2Passed cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("Cond2Passed", value)
	}
}

// SetCond3Passed actualiza Cond3Passed con sincronización via channel
func (p *Property) SetCond3Passed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond3Passed
	p.Cond3Passed = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] Cond3Passed cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("Cond3Passed", value)
	}
}

// SetCond4Passed actualiza Cond4Passed con sincronización via channel
func (p *Property) SetCond4Passed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond4Passed
	p.Cond4Passed = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] Cond4Passed cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("Cond4Passed", value)
	}
}

// SetCond5Passed actualiza Cond5Passed con sincronización via channel
func (p *Property) SetCond5Passed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond5Passed
	p.Cond5Passed = value
	p.mu.Unlock()

	if oldValue != value {
		log.Printf("[GID:%s] 📢 [PROPERTY DEBUG] Cond5Passed cambió: %v → %v", p.gID, oldValue, value)
		p.notifyUpdate("Cond5Passed", value)
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

// GetCond1Passed obtiene Cond1Passed de forma thread-safe
func (p *Property) GetCond1Passed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond1Passed
}

// GetCond2Passed obtiene Cond2Passed de forma thread-safe
func (p *Property) GetCond2Passed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond2Passed
}

// GetCond3Passed obtiene Cond3Passed de forma thread-safe
func (p *Property) GetCond3Passed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond3Passed
}

// GetCond4Passed obtiene Cond4Passed de forma thread-safe
func (p *Property) GetCond4Passed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond4Passed
}

// GetCond5Passed obtiene Cond5Passed de forma thread-safe
func (p *Property) GetCond5Passed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond5Passed
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

func (p *Property) SetCond1Processed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond1Processed
	p.Cond1Processed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond1Processed", value)
	}
}

func (p *Property) GetCond1Processed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond1Processed
}

func (p *Property) SetCond1Failed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond1Failed
	p.Cond1Failed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond1Failed", value)
	}
}

func (p *Property) GetCond1Failed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond1Failed
}

// ===== CONDITION 2 SETTERS/GETTERS =====

func (p *Property) SetCond2Processed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond2Processed
	p.Cond2Processed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond2Processed", value)
	}
}

func (p *Property) GetCond2Processed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond2Processed
}

func (p *Property) SetCond2Failed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond2Failed
	p.Cond2Failed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond2Failed", value)
	}
}

func (p *Property) GetCond2Failed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond2Failed
}

// ===== CONDITION 3 SETTERS/GETTERS =====

func (p *Property) SetCond3Processed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond3Processed
	p.Cond3Processed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond3Processed", value)
	}
}

func (p *Property) GetCond3Processed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond3Processed
}

func (p *Property) SetCond3Failed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond3Failed
	p.Cond3Failed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond3Failed", value)
	}
}

func (p *Property) GetCond3Failed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond3Failed
}

// ===== CONDITION 4 SETTERS/GETTERS =====

func (p *Property) SetCond4Processed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond4Processed
	p.Cond4Processed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond4Processed", value)
	}
}

func (p *Property) GetCond4Processed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond4Processed
}

func (p *Property) SetCond4Failed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond4Failed
	p.Cond4Failed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond4Failed", value)
	}
}

func (p *Property) GetCond4Failed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond4Failed
}

// ===== CONDITION 5 SETTERS/GETTERS =====

func (p *Property) SetCond5Processed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond5Processed
	p.Cond5Processed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond5Processed", value)
	}
}

func (p *Property) GetCond5Processed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond5Processed
}

func (p *Property) SetCond5Failed(value bool) {
	p.mu.Lock()
	oldValue := p.Cond5Failed
	p.Cond5Failed = value
	p.mu.Unlock()
	if oldValue != value {
		p.notifyUpdate("Cond5Failed", value)
	}
}

func (p *Property) GetCond5Failed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Cond5Failed
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
