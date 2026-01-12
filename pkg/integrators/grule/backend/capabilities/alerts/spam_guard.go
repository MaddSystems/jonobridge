package alerts

import (
	"log"
)

func (c *AlertsCapability) MarkAlertSent(imei, alertID string) bool {
	guardKey := imei + ":" + alertID

	c.guardMu.Lock()
	defer c.guardMu.Unlock()

	if c.sentGuard[guardKey] {
		log.Printf("🛡️ [AlertGuard] Race detected: IMEI=%s, AlertID=%s already sent", imei, alertID)
		return false
	}

	c.sentGuard[guardKey] = true
	log.Printf("✅ [AlertGuard] Marked alert sent: IMEI=%s, AlertID=%s", imei, alertID)
	return true
}

func (c *AlertsCapability) IsAlertSent(imei, alertID string) bool {
	guardKey := imei + ":" + alertID

	// Check in-memory guard (fast, prevents duplicate alerts)
	c.guardMu.RLock()
	sent := c.sentGuard[guardKey]
	c.guardMu.RUnlock()

	if sent {
		log.Printf("🛡️ [AlertGuard] BLOCKED: IMEI=%s, AlertID=%s (already sent)", imei, alertID)
		return true
	}

	log.Printf("❌ [AlertGuard] ALLOWED: IMEI=%s, AlertID=%s", imei, alertID)
	return false
}
