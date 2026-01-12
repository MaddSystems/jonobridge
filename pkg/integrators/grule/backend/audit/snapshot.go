package audit

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/jonobridge/grule-backend/capabilities"
)

// extractSnapshot builds rich context snapshot using SnapshotProvider pattern
// Each capability self-reports its data - no modification needed when adding new capabilities
// packetOverride allows passing the original packet pointer if DataContext holds a wrapped value
func extractSnapshot(dc ast.IDataContext, imei string, packetOverride interface{}) (map[string]interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[AuditListener] Panic during snapshot: %v", r)
		}
	}()

	snapshot := make(map[string]interface{})

	// 1. Extract packet_current
	var packetObj interface{}
	if packetOverride != nil {
		packetObj = packetOverride
	} else {
		packetObj = dc.Get("IncomingPacket")
	}

	if packetObj != nil {
		extracted := safeExtract(packetObj)
		snapshot["packet_current"] = extracted

		// Debug: Log what we extracted
		if m, ok := extracted.(map[string]interface{}); ok {
			if len(m) == 0 {
				log.Printf("[Snapshot] WARNING: packet_current is empty map. packetObj type: %T", packetObj)
			}
		}
	} else {
		log.Printf("[Snapshot] WARNING: IncomingPacket is nil in DataContext")
	}

	// 2. Collect all SnapshotProviders from DataContext
	providers := collectSnapshotProviders(dc)

	// 3. Each provider contributes its own data (Open/Closed Principle)
	for _, provider := range providers {
		if data := provider.GetSnapshotData(imei); data != nil {
			for key, value := range data {
				snapshot[key] = value
			}
		}
	}

	return snapshot, nil
}

// collectSnapshotProviders gathers all capabilities that implement SnapshotProvider
func collectSnapshotProviders(dc ast.IDataContext) []capabilities.SnapshotProvider {
	var providers []capabilities.SnapshotProvider

	// Get StateWrapper which holds all capabilities
	stateObj := dc.Get("state")
	if stateObj == nil {
		return providers
	}

	v := reflect.ValueOf(stateObj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return providers
	}

	// Iterate over exported fields to find SnapshotProviders
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		// Skip unexported fields to avoid panic
		if !field.CanInterface() {
			continue
		}
		if p, ok := field.Interface().(capabilities.SnapshotProvider); ok {
			providers = append(providers, p)
		}
	}

	return providers
}

// safeExtract handles conversion of any object to map[string]interface{} using JSON
func safeExtract(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("[Snapshot] json.Marshal error: %v, type: %T", err, v)
		return fmt.Sprintf("%+v", v)
	}
	// Debug: log the raw JSON
	if len(data) < 50 {
		log.Printf("[Snapshot] safeExtract raw JSON: %s, type: %T", string(data), v)
	}
	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		log.Printf("[Snapshot] json.Unmarshal error: %v", err)
		return fmt.Sprintf("%+v", v)
	}
	return res
}

// getIMEI extracts IMEI field from object using reflection
func getIMEI(obj interface{}) string {
	if obj == nil {
		return ""
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}
	f := v.FieldByName("IMEI")
	if f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}
	return ""
}

// getFromContext is a helper for the listener (though we use reflection now,
// the plan might still want this for other things)
func getFromContext[T any](dc ast.IDataContext, name string) T {
	var zero T
	obj := dc.Get(name)
	if obj == nil {
		return zero
	}
	typed, ok := obj.(T)
	if !ok {
		return zero
	}
	return typed
}
