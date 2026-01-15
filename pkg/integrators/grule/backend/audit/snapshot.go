package audit

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/jonobridge/grule-backend/capabilities"
)

// ExtractSnapshot builds rich context snapshot using SnapshotProvider pattern
// Each capability self-reports its data - no modification needed when adding new capabilities
// packetOverride allows passing the original packet pointer if DataContext holds a wrapped value
func ExtractSnapshot(dc ast.IDataContext, imei string, packetOverride interface{}) (map[string]interface{}, error) {
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
		log.Printf("📸 [Snapshot] Using packetOverride (Type: %T)", packetOverride)
	} else {
		packetObj = dc.Get("IncomingPacket")
		log.Printf("📸 [Snapshot] Using IncomingPacket from DataContext (Type: %T)", packetObj)
	}

	if packetObj != nil {
		// Explicit type check for debugging (using reflection to avoid circular dependency)
		typeStr := reflect.TypeOf(packetObj).String()
		log.Printf("[ExtractSnapshot] packetObj type: %s", typeStr)

		var realPacket interface{} = packetObj

		// Unwrap Grule node if needed
		// We check for an interface that has a GetValue() or Value() method
		// Grule internal nodes usually implement these.
		if gv, ok := packetObj.(interface{ GetValue() reflect.Value }); ok {
			val := gv.GetValue()
			if val.IsValid() && val.CanInterface() {
				realPacket = val.Interface()
				log.Printf("[ExtractSnapshot] Unwrapped via GetValue() to: %T", realPacket)
			}
		} else if gv, ok := packetObj.(interface{ Value() reflect.Value }); ok {
			val := gv.Value()
			if val.IsValid() && val.CanInterface() {
				realPacket = val.Interface()
				log.Printf("[ExtractSnapshot] Unwrapped via Value() to: %T", realPacket)
			}
		} else if gv, ok := packetObj.(interface{ GetValue() interface{} }); ok {
			realPacket = gv.GetValue()
			log.Printf("[ExtractSnapshot] Unwrapped via GetValue() interface to: %T", realPacket)
		}

		// Now extract manually if it matches our expected type name string
		realTypeStr := reflect.TypeOf(realPacket).String()
		if realTypeStr == "*grule.IncomingPacket" {
			v := reflect.ValueOf(realPacket)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			// Manual extraction map to bypass JSON marshaling issues
			extracted := map[string]interface{}{
				"IMEI":                    getFieldString(v, "IMEI"),
				"Speed":                   getFieldInt(v, "Speed"),
				"GSMSignalStrength":       getFieldInt(v, "GSMSignalStrength"),
				"Datetime":                fmt.Sprintf("%v", getFieldValue(v, "Datetime")),
				"PositioningStatus":       getFieldString(v, "PositioningStatus"),
				"Latitude":                getFieldFloat(v, "Latitude"),
				"Longitude":               getFieldFloat(v, "Longitude"),
				"BufferUpdated":           getFieldBool(v, "BufferUpdated"),
				"BufferHas10":             getFieldBool(v, "BufferHas10"),
				"IsOfflineFor5Min":        getFieldBool(v, "IsOfflineFor5Min"),
				"PositionInvalidDetected": getFieldBool(v, "PositionInvalidDetected"),	
				"MetricsReady":            getFieldBool(v, "MetricsReady"),
				"MovingWithWeakSignal":    getFieldBool(v, "MovingWithWeakSignal"),
				"OutsideAllSafeZones":     getFieldBool(v, "OutsideAllSafeZones"),
			}
			
			snapshot["packet_current"] = extracted
			log.Printf("[ExtractSnapshot] Manually extracted packet_current: BufferUpdated=%v, BufferHas10=%v", 
				extracted["BufferUpdated"], extracted["BufferHas10"])
		} else {
			// Fallback for other types
			extracted := safeExtract(realPacket)
			snapshot["packet_current"] = extracted
			if realTypeStr != typeStr {
				log.Printf("[ExtractSnapshot] Unwrapped fallback type: %s", realTypeStr)
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

// Helper functions for reflection extraction

func getFieldValue(v reflect.Value, name string) interface{} {
	f := v.FieldByName(name)
	if f.IsValid() && f.CanInterface() {
		return f.Interface()
	}
	return nil
}

func getFieldString(v reflect.Value, name string) string {
	f := v.FieldByName(name)
	if f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}
	return ""
}

func getFieldInt(v reflect.Value, name string) int {
	f := v.FieldByName(name)
	if f.IsValid() {
		if f.Kind() == reflect.Int || f.Kind() == reflect.Int64 {
			return int(f.Int())
		}
	}
	return 0
}

func getFieldFloat(v reflect.Value, name string) float64 {
	f := v.FieldByName(name)
	if f.IsValid() && f.Kind() == reflect.Float64 {
		return f.Float()
	}
	return 0
}

func getFieldBool(v reflect.Value, name string) bool {
	f := v.FieldByName(name)
	if f.IsValid() && f.Kind() == reflect.Bool {
		return f.Bool()
	}
	return false
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
		log.Printf("[safeExtract] Input is NIL")
		return nil
	}
	
	// Key debug: Log type and value
	log.Printf("[safeExtract] Input type: %T, value: %+v", v, v)

	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("[Snapshot] json.Marshal error: %v, type: %T, value: %+v", err, v, v)
		return fmt.Sprintf("%+v", v)
	}
	
	// Key debug: Log raw JSON
	if len(data) < 200 {
		log.Printf("[safeExtract] Raw JSON from marshal: %s", string(data))
	} else {
		log.Printf("[safeExtract] Raw JSON from marshal (truncated): %s...", string(data[:200]))
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
