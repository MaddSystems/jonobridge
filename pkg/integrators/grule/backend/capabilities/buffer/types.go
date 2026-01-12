package buffer

import (
	"time"
)

type BufferEntry struct {
	IMEI              string    `json:"imei"`
	Datetime          time.Time `json:"datetime"`
	Speed             int64     `json:"speed"`              // YA EN KM/H
	GSMSignalStrength int64     `json:"gsm_signal"`         // 0-31
	PositioningStatus string    `json:"positioning_status"` // "true"/"A"/"false"/"V"
	IsValid           bool      `json:"is_valid"`           // Computed
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
}
