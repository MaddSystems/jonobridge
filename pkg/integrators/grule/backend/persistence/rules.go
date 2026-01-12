package persistence

import (
	"log"
)

type Rule struct {

	ID            int64  `json:"id"`

	Name          string `json:"name"`

	Description   string `json:"description,omitempty"`

	GRL           string `json:"grl"`

	AuditManifest string `json:"audit_manifest,omitempty"` // NEW: YAML content

	Active        bool   `json:"active"`

	Priority      int    `json:"priority"`

	CreatedAt     string `json:"created_at,omitempty"`

	UpdatedAt     string `json:"updated_at,omitempty"`

}



func (s *MySQLStateStore) GetAllRules() ([]Rule, error) {

	query := `

		SELECT id, name, COALESCE(description, '') as description, 

		       grl_content, audit_manifest, active, priority, 

		       created_at, updated_at

		FROM fleet_rules 

		ORDER BY priority DESC, id DESC

	`



	rows, err := s.db.Query(query)

	if err != nil {

		return nil, err

	}

	defer rows.Close()



	var rules []Rule

	for rows.Next() {

		var r Rule

		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.GRL, &r.AuditManifest,

			&r.Active, &r.Priority, &r.CreatedAt, &r.UpdatedAt); err != nil {

			log.Printf("⚠️  Error escaneando regla: %v", err)

			continue

		}

		rules = append(rules, r)

	}



	return rules, nil

}



func (s *MySQLStateStore) GetRule(id int64) (*Rule, error) {

	var r Rule

	query := `SELECT id, name, grl_content, audit_manifest, active, priority FROM fleet_rules WHERE id = ?`

	err := s.db.QueryRow(query, id).Scan(&r.ID, &r.Name, &r.GRL, &r.AuditManifest, &r.Active, &r.Priority)

	if err != nil {

		return nil, err

	}

	return &r, nil

}



func (s *MySQLStateStore) CreateRule(rule Rule) (int64, error) {

	query := `

		INSERT INTO fleet_rules (name, description, grl_content, audit_manifest, priority, active)

		VALUES (?, ?, ?, ?, ?, ?)

	`

	priority := rule.Priority

	if priority == 0 {

		priority = 100

	}

	result, err := s.db.Exec(query, rule.Name, rule.Description, rule.GRL, rule.AuditManifest, priority, rule.Active)

	if err != nil {

		return 0, err

	}

	return result.LastInsertId()

}



func (s *MySQLStateStore) UpdateRule(rule Rule) error {

	query := `

		UPDATE fleet_rules 

		SET name = ?, description = ?, grl_content = ?, audit_manifest = ?,

		    priority = ?, active = ?, updated_at = NOW()

		WHERE id = ?

	`

	_, err := s.db.Exec(query, rule.Name, rule.Description, rule.GRL, rule.AuditManifest,

		rule.Priority, rule.Active, rule.ID)

	return err

}



func (s *MySQLStateStore) DeleteRule(id int64) error {

	_, err := s.db.Exec("DELETE FROM fleet_rules WHERE id = ?", id)

	return err

}



func (s *MySQLStateStore) LoadActiveRules() ([]Rule, error) {

	query := `

		SELECT id, name, grl_content, audit_manifest

		FROM fleet_rules 

		WHERE active = 1

		ORDER BY priority DESC

		LIMIT 500

	`

	rows, err := s.db.Query(query)

	if err != nil {

		return nil, err

	}

	defer rows.Close()



	var rules []Rule

	for rows.Next() {

		var r Rule

		if err := rows.Scan(&r.ID, &r.Name, &r.GRL, &r.AuditManifest); err != nil {

			continue

		}

		rules = append(rules, r)

	}

	return rules, nil

}