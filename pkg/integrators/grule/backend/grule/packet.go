package grule

import "time"

type IncomingPacket struct {
	IMEI              string
	Speed             int64
	GSMSignalStrength int64
	Datetime          time.Time
	PositioningStatus string
	Latitude          float64
	Longitude         float64

	// Logic flags
	BufferUpdated           bool
	BufferHas10             bool
	IsOfflineFor5Min        bool
	PositionInvalidDetected bool
	MetricsReady            bool
	MovingWithWeakSignal    bool
	OutsideAllSafeZones     bool
}
