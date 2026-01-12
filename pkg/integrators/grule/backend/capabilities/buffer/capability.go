package buffer

import (
	"time"
)

type BufferCapability struct {
	manager *BufferManager
}

func NewBufferCapability() *BufferCapability {
	return &BufferCapability{
		manager: NewBufferManager(24 * time.Hour),
	}
}

func (c *BufferCapability) Name() string {
	return "buffer"
}

func (c *BufferCapability) Version() string {
	return "1.0.0"
}

func (c *BufferCapability) GetDataContextName() string {
	return "buffer"
}

func (c *BufferCapability) Initialize(imei string) error {
	return nil
}

func (c *BufferCapability) GetSnapshot() map[string]interface{} {
	return map[string]interface{}{}
}

// GetSnapshotData implements SnapshotProvider
func (c *BufferCapability) GetSnapshotData(imei string) map[string]interface{} {
	if c == nil {
		return nil
	}

	entries := c.GetEntriesInTimeWindow90Min(imei)
	var bufferData []map[string]interface{}
	for _, e := range entries {
		bufferData = append(bufferData, map[string]interface{}{
			"imei":               e.IMEI,
			"datetime":           e.Datetime.Format(time.RFC3339),
			"speed":              e.Speed,
			"gsm_signal":         e.GSMSignalStrength,
			"positioning_status": e.PositioningStatus,
			"latitude":           e.Latitude,
			"longitude":          e.Longitude,
			"is_valid":           e.IsValid,
		})
	}

	return map[string]interface{}{
		"buffer_circular": bufferData,
	}
}

func (c *BufferCapability) AddToBuffer(imei string, speed int64, gsmSignal int64, datetime time.Time, posStatus string, lat, lon float64) bool {
	buffer := c.manager.GetOrCreateBuffer(imei)
	entry := BufferEntry{
		IMEI:              imei,
		Datetime:          datetime,
		Speed:             speed,
		GSMSignalStrength: gsmSignal,
		PositioningStatus: posStatus,
		Latitude:          lat,
		Longitude:         lon,
	}
	return buffer.AddEntry(entry)
}

func (c *BufferCapability) GetEntriesInTimeWindow90Min(imei string) []BufferEntry {
	buffer := c.manager.GetOrCreateBuffer(imei)
	return buffer.GetEntriesInTimeWindow90Min()
}
