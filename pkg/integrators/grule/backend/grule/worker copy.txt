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

type Worker struct {
	contextBuilder *ContextBuilder
	adapter        TrackerAdapter
	knowledgeBases []*ast.KnowledgeBase
	manifest       *audit.AuditManifest
}

func NewWorker(cb *ContextBuilder, adapter TrackerAdapter, kbs []*ast.KnowledgeBase, manifest *audit.AuditManifest) *Worker {
	return &Worker{
		contextBuilder: cb,
		adapter:        adapter,
		knowledgeBases: kbs,
		manifest:       manifest,
	}
}

func (w *Worker) UpdateRules(kbs []*ast.KnowledgeBase, manifest *audit.AuditManifest) {
	w.knowledgeBases = kbs
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
			eng.Listeners = []engine.GruleEngineListener{
				audit.NewAuditListener(w.manifest),
			}
		} else {
			log.Printf("⚠️ [Worker] No audit manifest loaded, listener will not be attached")
		}

		for _, kb := range w.knowledgeBases {
			// Create context with dataContext and imei for listener
			ctx := context.WithValue(context.Background(), "dataContext", dataContext)
			ctx = context.WithValue(ctx, "imei", packet.IMEI)
			ctx = context.WithValue(ctx, "originalPacket", packet) // Pass original packet to bypass DataContext wrappers
			err = eng.ExecuteWithContext(ctx, dataContext, kb)
			if err != nil {
				log.Printf("❌ [Worker] Error executing rules for IMEI %s: %v", packet.IMEI, err)
			}
		}
	}
}
