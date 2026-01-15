package grule

import (
	"context"
	"log"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/jonobridge/grule-backend/audit"
)

type TrackerAdapter interface {
	Parse(payload string) ([]*IncomingPacket, error)
}

type RuleKB struct {
	RuleName string
	KB       *ast.KnowledgeBase
	Order    int
}

type Worker struct {
	contextBuilder *ContextBuilder
	adapter        TrackerAdapter
	ruleKBs        []RuleKB
	manifest       *audit.AuditManifest
}

func NewWorker(cb *ContextBuilder, adapter TrackerAdapter, kbs []RuleKB, manifest *audit.AuditManifest) *Worker {
	return &Worker{
		contextBuilder: cb,
		adapter:        adapter,
		ruleKBs:        kbs,
		manifest:       manifest,
	}
}

func (w *Worker) UpdateRules(kbs []RuleKB, manifest *audit.AuditManifest) {
	w.ruleKBs = kbs
	w.manifest = manifest
}

func (w *Worker) Process(payload string) {
	log.Printf("📥 [Worker] Processing new payload: %s", payload)
	packets, err := w.adapter.Parse(payload)
	if err != nil {
		log.Printf("❌ [Worker] Error parsing packets: %v", err)
		return
	}

	log.Printf("📦 [Worker] Parsed %d packets", len(packets))

	for _, packet := range packets {
		log.Printf("🛰️ [Worker] Running rules for IMEI %s (Speed: %d)", packet.IMEI, packet.Speed)
		dataContext, err := w.contextBuilder.Build(packet)
		if err != nil {
			log.Printf("❌ [Worker] Error building context for IMEI %s: %v", packet.IMEI, err)
			continue
		}

		eng := engine.NewGruleEngine()

		// Wire listener for declarative audit
		if w.manifest != nil {
			listener := audit.NewAuditListener(w.manifest)
			listener.SetPacket(packet)
			eng.Listeners = []engine.GruleEngineListener{
				listener,
			}
		} else {
			log.Printf("⚠️ [Worker] No audit manifest loaded, listener will not be attached")
		}

		for _, rkb := range w.ruleKBs {
			// Create context with dataContext and imei for listener
			ctx := context.WithValue(context.Background(), "dataContext", dataContext)
			ctx = context.WithValue(ctx, "imei", packet.IMEI)
			ctx = context.WithValue(ctx, "originalPacket", packet) // Pass original packet to bypass DataContext wrappers
			
			log.Printf("🔍 [Worker] Executing rule '%s' (Order: %d)", rkb.RuleName, rkb.Order)
			err = eng.ExecuteWithContext(ctx, dataContext, rkb.KB)
			if err != nil {
				log.Printf("❌ [Worker] Error executing rule '%s' for IMEI %s: %v", rkb.RuleName, packet.IMEI, err)
			} else {
				log.Printf("✅ [Worker] Execution finished for '%s'", rkb.RuleName)
			}

			// Capture POST-execution snapshot
			if w.manifest != nil {
				meta := w.manifest.GetRuleMeta(rkb.RuleName)
				if meta != nil {
					if meta.Enabled {
						log.Printf("[Worker] Capturing post-snapshot for '%s' (using dc state, no override)", rkb.RuleName)
						snapshot, err := audit.ExtractSnapshot(dataContext, packet.IMEI, nil)
						if err != nil {
							log.Printf("❌ [Worker] Error capturing post-snapshot for rule '%s': %v", rkb.RuleName, err)
						} else {
							log.Printf("📤 [Worker] Sending post-capture for '%s'", rkb.RuleName)
							audit.Capture(&audit.AuditEntry{
								IMEI:         packet.IMEI,
								RuleName:     rkb.RuleName,
								Description:  meta.Description,
								Level:        meta.Level,
								IsAlert:      meta.IsAlert,
								StepNumber:   meta.Order,
								StageReached: meta.Description,
								Snapshot:     snapshot,
								IsPost:       true,
							})
						}
					}
				} else {
					log.Printf("⚠️ [Worker] No manifest meta for rule '%s'", rkb.RuleName)
				}
			} else {
				log.Printf("⚠️ [Worker] Manifest is nil, skipping post-capture")
			}
		}
	}
}
