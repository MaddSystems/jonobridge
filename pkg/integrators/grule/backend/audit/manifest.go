package audit

import (
	"fmt"
	"log"
	"github.com/jonobridge/grule-backend/persistence"
	"gopkg.in/yaml.v3"
)

type AuditManifest struct {
	rules map[string]*RuleMeta
}

type RuleMeta struct {
	Enabled     bool     // false = skip auditing this rule
	Description string
	Level       string   // debug, info, warning, critical
	IsAlert     bool
	Order       int
	Snapshot    []string // Fields to capture: ["packet", "state"]. Empty = all
}

func NewAuditManifest() *AuditManifest {
	return &AuditManifest{
		rules: make(map[string]*RuleMeta),
	}
}

// ParseYAML parses a single YAML manifest string
func (m *AuditManifest) ParseYAML(yamlContent string) error {
	var manifest struct {
		Stages []struct {
			Rule  string `yaml:"rule"`
			Order int    `yaml:"order"`
			Audit struct {
				Enabled     *bool    `yaml:"enabled"`
				Description string   `yaml:"description"`
				Level       string   `yaml:"level"`
				IsAlert     bool     `yaml:"is_alert"`
				Snapshot    []string `yaml:"snapshot"`
			} `yaml:"audit"`
		} `yaml:"stages"`
	}

	if err := yaml.Unmarshal([]byte(yamlContent), &manifest); err != nil {
		return fmt.Errorf("YAML parse error: %w", err)
	}

	for _, stage := range manifest.Stages {
		enabled := true
		if stage.Audit.Enabled != nil {
			enabled = *stage.Audit.Enabled
		}

		m.rules[stage.Rule] = &RuleMeta{
			Enabled:     enabled,
			Description: stage.Audit.Description,
			Level:       stage.Audit.Level,
			IsAlert:     stage.Audit.IsAlert,
			Order:       stage.Order,
			Snapshot:    stage.Audit.Snapshot,
		}
	}
	return nil
}

// LoadFromRules loads all audit manifests from fleet_rules records
func (m *AuditManifest) LoadFromRules(rules []persistence.Rule) error {
	for _, rule := range rules {
		if rule.AuditManifest == "" {
			log.Printf("[Manifest] Rule '%s' has no audit manifest, skipping", rule.Name)
			continue
		}

		if err := m.ParseYAML(rule.AuditManifest); err != nil {
			log.Printf("[Manifest] Warning: Invalid manifest for rule '%s': %v", rule.Name, err)
			continue // Don't fail startup, just skip this rule
		}

		log.Printf("[Manifest] Loaded manifest for rule '%s'", rule.Name)
	}
	return nil
}

func (m *AuditManifest) GetRuleMeta(ruleName string) *RuleMeta {
	if meta, ok := m.rules[ruleName]; ok {
		return meta
	}
	return nil
}

func (m *AuditManifest) Count() int {
	return len(m.rules)
}
