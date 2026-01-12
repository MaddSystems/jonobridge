package grule

import (
	"log"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/jonobridge/grule-backend/capabilities"
	"github.com/jonobridge/grule-backend/capabilities/alerts"
	"github.com/jonobridge/grule-backend/capabilities/buffer"
	"github.com/jonobridge/grule-backend/capabilities/geofence"
	"github.com/jonobridge/grule-backend/capabilities/metrics"
	"github.com/jonobridge/grule-backend/capabilities/timing"
)

type ContextBuilder struct {
	registry *capabilities.Registry
}

func NewContextBuilder(registry *capabilities.Registry) *ContextBuilder {
	return &ContextBuilder{registry: registry}
}

func (cb *ContextBuilder) Build(packet *IncomingPacket) (ast.IDataContext, error) {
	dc := ast.NewDataContext()

	// 1. Add IncomingPacket
	dc.Add("IncomingPacket", packet)

	// 2. Build State Wrapper
	stateWrapper := &StateWrapper{
		imei:        packet.IMEI,
		AlertStates: make(map[string]bool),
	}

	// Resolve capabilities
	if cap := cb.registry.Get("geofence"); cap != nil {
		stateWrapper.Geo = cap.(*geofence.GeofenceCapability)
		stateWrapper.Geo.UpdateLastPacket(packet.IMEI, packet.Latitude, packet.Longitude)
	}
	if cap := cb.registry.Get("buffer"); cap != nil {
		stateWrapper.Buf = cap.(*buffer.BufferCapability)
	}
	if cap := cb.registry.Get("metrics"); cap != nil {
		stateWrapper.Met = cap.(*metrics.MetricsCapability)
	}
	if cap := cb.registry.Get("timing"); cap != nil {
		stateWrapper.Tim = cap.(*timing.TimingCapability)
	}
	if cap := cb.registry.Get("alerts"); cap != nil {
		stateWrapper.Alrt = cap.(*alerts.AlertsCapability)
		log.Printf("✅ [ContextBuilder] AlertsCapability attached for IMEI %s", packet.IMEI)
		// Global state is checked directly in GRL now
	} else {
		log.Printf("❌ [ContextBuilder] AlertsCapability NOT FOUND for IMEI %s", packet.IMEI)
	}

	dc.Add("state", stateWrapper)

	// 3. Build Actions Wrapper
	actionsWrapper := &ActionsWrapper{
		imei:         packet.IMEI,
		stateWrapper: stateWrapper,
	}
	if stateWrapper.Alrt != nil {
		actionsWrapper.alrt = stateWrapper.Alrt
	}

	dc.Add("actions", actionsWrapper)

	return dc, nil
}

// StateWrapper mimics PersistentState
type StateWrapper struct {
	imei string

	Geo  *geofence.GeofenceCapability
	Buf  *buffer.BufferCapability
	Met  *metrics.MetricsCapability
	Tim  *timing.TimingCapability
	Alrt *alerts.AlertsCapability

	// Logic Variables (public for Rules)
	JammerAvgSpeed90min int64
	JammerAvgGsm5       int64

	// Universal Alert State (keyed by rule name)
	AlertStates map[string]bool
}

func (s *StateWrapper) UpdateMemoryBuffer(speed int64, gsmSignal int64, datetime time.Time, posStatus string, lat, lon float64) bool {
	if s.Buf == nil {
		return false
	}
	// Also update timing state
	if s.Tim != nil {
		s.Tim.UpdateState(s.imei, datetime, posStatus)
	}
	return s.Buf.AddToBuffer(s.imei, speed, gsmSignal, datetime, posStatus, lat, lon)
}

// IsAlertSentForRule checks if an alert has been sent for a specific rule (local state first, then global)
func (s *StateWrapper) IsAlertSentForRule(ruleName string) bool {
	// 1. Check local state (fast, intra-packet)
	if sent, exists := s.AlertStates[ruleName]; exists && sent {
		log.Printf("🛡️ [StateWrapper] BLOCKED (local): IMEI=%s, Rule=%s", s.imei, ruleName)
		return true
	}

	// 2. Check global guard (inter-packet)
	if s.Alrt != nil && s.Alrt.IsAlertSent(s.imei, ruleName) {
		// Cache true result locally for subsequent cycles
		s.AlertStates[ruleName] = true
		log.Printf("🛡️ [StateWrapper] BLOCKED (global): IMEI=%s, Rule=%s", s.imei, ruleName)
		return true
	}

	log.Printf("❌ [StateWrapper] ALLOWED: IMEI=%s, Rule=%s", s.imei, ruleName)
	return false
}

// MarkAlertSentForRule marks an alert as sent for a specific rule (updates both local and global)
func (s *StateWrapper) MarkAlertSentForRule(ruleName string) bool {
	log.Printf("🎯 [StateWrapper] MarkAlertSentForRule called: IMEI=%s, Rule=%s", s.imei, ruleName)

	if s.Alrt == nil {
		log.Printf("❌ [StateWrapper] Alrt is nil! Cannot mark alert sent")
		return false
	}

	// 1. Mark global (Atomic Check-and-Set)
	// Returns true if WE won the race and marked it.
	// Returns false if it was already marked.
	wonRace := s.Alrt.MarkAlertSent(s.imei, ruleName)

	// 2. Mark local (stop future cycles)
	// Even if we lost the race (wonRace=false), the alert IS sent, so we must block local cycles.
	s.AlertStates[ruleName] = true

	if wonRace {
		log.Printf("✅ [StateWrapper] Alert successfully marked (Winner): IMEI=%s, Rule=%s", s.imei, ruleName)
	} else {
		log.Printf("🛡️ [StateWrapper] Alert already marked (Lost Race): IMEI=%s, Rule=%s", s.imei, ruleName)
	}

	return wonRace
}

func (s *StateWrapper) IsOfflineFor(minutes int64) bool {
	if s.Tim == nil {
		return false
	}
	return s.Tim.IsOfflineFor(s.imei, minutes)
}

func (s *StateWrapper) GetAverageSpeed90Min(imei string) int64 {
	if s.Met == nil {
		return 0
	}
	return s.Met.GetAverageSpeed90Min(imei)
}

func (s *StateWrapper) GetAverageGSMLast5(imei string) int64 {
	if s.Met == nil {
		return 0
	}
	return s.Met.GetAverageGSMLast5(imei)
}

func (s *StateWrapper) IsInsideGroup(groupName string, lat, lon float64) bool {
	if s.Geo == nil {
		return false
	}
	return s.Geo.IsInsideGroup(groupName, lat, lon)
}

func (s *StateWrapper) MarkAlertSent(alertID string) bool {
	log.Printf("🎯 [StateWrapper] MarkAlertSent called: IMEI=%s, AlertID=%s, Alrt=%v", s.imei, alertID, s.Alrt != nil)
	if s.Alrt == nil {
		log.Printf("❌ [StateWrapper] Alrt is nil! Cannot mark alert sent")
		return false
	}
	return s.Alrt.MarkAlertSent(s.imei, alertID)
}

func (s *StateWrapper) IsAlertSent(alertID string) bool {
	if s.Alrt == nil {
		return false
	}
	return s.Alrt.IsAlertSent(s.imei, alertID)
}

// ActionsWrapper mimics ActionsHelper
type ActionsWrapper struct {
	imei         string
	alrt         *alerts.AlertsCapability
	stateWrapper *StateWrapper
}

func (a *ActionsWrapper) Log(message string) {
	if a.alrt != nil {
		a.alrt.Log(message)
	}
}

func (a *ActionsWrapper) SendTelegram(message string) {
	if a.alrt != nil {
		a.alrt.SendTelegram(message)
	}
}

func (a *ActionsWrapper) CastString(v interface{}) string {
	if a.alrt != nil {
		return a.alrt.CastString(v)
	}
	return ""
}
