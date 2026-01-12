package buffer

import (
	"sync"
	"time"
)

type BufferManager struct {
	buffers   map[string]*FixedCircularBuffer
	mutex     sync.RWMutex
	retention time.Duration
}

func NewBufferManager(retention time.Duration) *BufferManager {
	bm := &BufferManager{
		buffers:   make(map[string]*FixedCircularBuffer),
		retention: retention,
	}
	// Start cleanup routine
	go bm.startCleanupRoutine()
	return bm
}

func (bm *BufferManager) GetOrCreateBuffer(imei string) *FixedCircularBuffer {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if buffer, exists := bm.buffers[imei]; exists {
		buffer.lastUpdate = time.Now()
		return buffer
	}

	newBuffer := &FixedCircularBuffer{
		size:       0,
		lastUpdate: time.Now(),
	}
	bm.buffers[imei] = newBuffer
	return newBuffer
}

func (bm *BufferManager) startCleanupRoutine() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bm.cleanupInactiveBuffers()
		}
	}
}

func (bm *BufferManager) cleanupInactiveBuffers() {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	now := time.Now()
	for imei, buffer := range bm.buffers {
		if now.Sub(buffer.lastUpdate) > bm.retention {
			delete(bm.buffers, imei)
		}
	}
}
