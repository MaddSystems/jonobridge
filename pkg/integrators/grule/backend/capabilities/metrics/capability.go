package metrics

import (
	"github.com/jonobridge/grule-backend/capabilities/buffer"
)

type MetricsCapability struct {
	bufferCap *buffer.BufferCapability
}

func NewMetricsCapability(bufferCap *buffer.BufferCapability) *MetricsCapability {
	return &MetricsCapability{
		bufferCap: bufferCap,
	}
}

func (c *MetricsCapability) Name() string {
	return "metrics"
}

func (c *MetricsCapability) Version() string {
	return "1.0.0"
}

func (c *MetricsCapability) GetDataContextName() string {
	return "metrics"
}

func (c *MetricsCapability) Initialize(imei string) error {
	return nil
}

func (c *MetricsCapability) GetSnapshot() map[string]interface{} {
	return map[string]interface{}{}
}

// GetSnapshotData implements SnapshotProvider
func (c *MetricsCapability) GetSnapshotData(imei string) map[string]interface{} {
	if c == nil {
		return nil
	}

	return map[string]interface{}{
		"jammer_metrics": map[string]interface{}{
			"avg_speed_90min": c.GetAverageSpeed90Min(imei),
			"avg_gsm_last5":   c.GetAverageGSMLast5(imei),
		},
	}
}
