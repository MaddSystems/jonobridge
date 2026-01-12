package metrics

func (c *MetricsCapability) GetAverageSpeed90Min(imei string) int64 {
	entries := c.bufferCap.GetEntriesInTimeWindow90Min(imei)
	if len(entries) == 0 {
		return 0
	}

	totalSpeed := 0.0
	for _, entry := range entries {
		totalSpeed += float64(entry.Speed)
	}

	return int64(totalSpeed / float64(len(entries)))
}

func (c *MetricsCapability) GetAverageGSMLast5(imei string) int64 {
	entries := c.bufferCap.GetEntriesInTimeWindow90Min(imei)
	if len(entries) < 5 {
		return 0
	}

	totalGSM := 0.0
	start := len(entries) - 5
	for i := start; i < len(entries); i++ {
		totalGSM += float64(entries[i].GSMSignalStrength)
	}

	return int64(totalGSM / 5.0)
}
