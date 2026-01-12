package buffer

import (
	"log"
	"sync"
	"time"
)

type FixedCircularBuffer struct {
	entries    [10]BufferEntry
	size       int
	mutex      sync.RWMutex
	lastUpdate time.Time
}

func (b *FixedCircularBuffer) AddEntry(entry BufferEntry) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	entry.IsValid = (entry.PositioningStatus == "true" || entry.PositioningStatus == "A")

	if b.size > 0 {
		var lastEntryTime time.Time
		if b.size < 10 {
			lastEntryTime = b.entries[b.size-1].Datetime
		} else {
			lastEntryTime = b.entries[9].Datetime
		}

		if !entry.Datetime.After(lastEntryTime) {
			log.Printf("⚠️ [BUFFER SKIP] Trama descartada por fecha antigua o duplicada. IMEI: %s. BufferLast: %v, New: %v", 
				entry.IMEI, lastEntryTime, entry.Datetime)
			return b.size >= 10
		}
	}

	if b.size < 10 {
		b.entries[b.size] = entry
		b.size++
		return b.size >= 10
	} else {
		for i := 0; i < 9; i++ {
			b.entries[i] = b.entries[i+1]
		}
		b.entries[9] = entry
		return true
	}
}

func (b *FixedCircularBuffer) GetEntriesInTimeWindow90Min() []BufferEntry {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.size != 10 {
		return []BufferEntry{}
	}

	now := time.Now()
	var filtered []BufferEntry

	for i := 0; i < 10; i++ {
		entry := b.entries[i]
		diferencia := now.Sub(entry.Datetime).Minutes()

		if diferencia < 90 {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
