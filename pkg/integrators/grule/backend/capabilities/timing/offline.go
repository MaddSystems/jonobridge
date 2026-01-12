package timing

import (
	"log"
)

func (c *TimingCapability) IsOfflineFor(imei string, minutes int64) bool {
	s := c.GetOrCreateState(imei)
	
	if s.LastValidPosition.IsZero() {
		return false
	}

	elapsed := s.CurrentPacketTime.Sub(s.LastValidPosition).Minutes()
	isOffline := elapsed >= float64(minutes)

	if isOffline {
		log.Printf("⏰ [%s] Offline for %.1f minutes (threshold: %d min)", imei, elapsed, minutes)
	}

	return isOffline
}
