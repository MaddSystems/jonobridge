package timing

import (
	"sync"
	"time"
)

type TimingState struct {
	LastValidPosition time.Time
	CurrentPacketTime time.Time
}

type TimingCapability struct {
	states map[string]*TimingState
	mutex  sync.RWMutex
}

func NewTimingCapability() *TimingCapability {
	return &TimingCapability{
		states: make(map[string]*TimingState),
	}
}

func (c *TimingCapability) Name() string {
	return "timing"
}

func (c *TimingCapability) Version() string {
	return "1.0.0"
}

func (c *TimingCapability) GetDataContextName() string {
	return "timing"
}

func (c *TimingCapability) Initialize(imei string) error {
	return nil
}

func (c *TimingCapability) GetSnapshot() map[string]interface{} {
	return map[string]interface{}{}
}

// GetSnapshotData implements SnapshotProvider
func (c *TimingCapability) GetSnapshotData(imei string) map[string]interface{} {
	if c == nil {
		return nil
	}

	c.mutex.RLock()
	state, ok := c.states[imei]
	c.mutex.RUnlock()

	if !ok {
		return nil
	}

	return map[string]interface{}{
		"timing_state": map[string]interface{}{
			"last_valid_position": state.LastValidPosition.Format(time.RFC3339),
			"current_packet_time": state.CurrentPacketTime.Format(time.RFC3339),
			"offline_minutes":     time.Since(state.LastValidPosition).Minutes(),
		},
	}
}

func (c *TimingCapability) GetOrCreateState(imei string) *TimingState {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if s, ok := c.states[imei]; ok {
		return s
	}
	s := &TimingState{}
	c.states[imei] = s
	return s
}

func (c *TimingCapability) UpdateState(imei string, datetime time.Time, posStatus string) {
	s := c.GetOrCreateState(imei)
	s.CurrentPacketTime = datetime
	if posStatus == "A" || posStatus == "true" {
		s.LastValidPosition = datetime
	}
}
