package alerts

import (
	"fmt"
	"sync"

	"github.com/jonobridge/grule-backend/persistence"
)

type AlertsCapability struct {
	store persistence.StateStore

	// In-memory guard to prevent race conditions
	// Key format: "imei:alertID"
	sentGuard map[string]bool
	guardMu   sync.RWMutex
}

func NewAlertsCapability(store persistence.StateStore) *AlertsCapability {
	return &AlertsCapability{
		store:     store,
		sentGuard: make(map[string]bool),
	}
}

func (c *AlertsCapability) Name() string {
	return "alerts"
}

func (c *AlertsCapability) Version() string {
	return "1.0.0"
}

func (c *AlertsCapability) GetDataContextName() string {
	return "actions"
}

func (c *AlertsCapability) Initialize(imei string) error {
	return nil
}

func (c *AlertsCapability) GetSnapshot() map[string]interface{} {
	return map[string]interface{}{}
}

func (c *AlertsCapability) CastString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
