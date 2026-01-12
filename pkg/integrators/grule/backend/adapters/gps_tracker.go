package adapters

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/jonobridge/grule-backend/grule"
)

type GPSTrackerAdapter struct{}

func NewGPSTrackerAdapter() *GPSTrackerAdapter {
	return &GPSTrackerAdapter{}
}

func (a *GPSTrackerAdapter) Parse(payload string) ([]*grule.IncomingPacket, error) {
	var jono models.JonoModel
	if err := json.Unmarshal([]byte(payload), &jono); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	if len(jono.ListPackets) == 0 {
		return nil, fmt.Errorf("no packets in payload")
	}

	var packets []*grule.IncomingPacket
	for _, p := range jono.ListPackets {
		if p.Datetime.IsZero() {
			log.Printf("⚠️ Invalid timestamp for IMEI %s, skipping", jono.IMEI)
			continue
		}

		speedKmH := int64(float64(p.Speed) * 3.6)

		var gsm int64
		if p.GSMSignalStrength != nil {
			gsm = int64(*p.GSMSignalStrength)
		}

		packet := &grule.IncomingPacket{
			IMEI:              jono.IMEI,
			Speed:             speedKmH,
			GSMSignalStrength: gsm,
			Datetime:          p.Datetime,
			PositioningStatus: p.PositioningStatus,
			Latitude:          p.Latitude,
			Longitude:         p.Longitude,
			// Initialize flags to false
			BufferUpdated:           false,
			BufferHas10:             false,
			IsOfflineFor5Min:        false,
			PositionInvalidDetected: false,
			MetricsReady:            false,
			MovingWithWeakSignal:    false,
			OutsideAllSafeZones:     false,
		}
		packets = append(packets, packet)
	}

	return packets, nil
}
