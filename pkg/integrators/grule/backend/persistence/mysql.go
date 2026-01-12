package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLStateStore struct {
	db *sql.DB
}

func NewMySQLStateStore() *MySQLStateStore {
	return &MySQLStateStore{
		db: initMySQL(),
	}
}

func initMySQL() *sql.DB {
	dsn := getMySQLDSN()
	log.Printf("🔌 Connecting to MySQL with DSN: %s", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Error connecting to MySQL: %v", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Could not connect to MySQL: %v", err)
	}

	log.Println("✅ Connected to MySQL (Backend)")

	// Get current database
	var currentDB string
	db.QueryRow("SELECT DATABASE()").Scan(&currentDB)
	log.Printf("📁 Current database: %s", currentDB)

	// Create tables if they don't exist
	createFleetRulesTable(db)
	createVehicleRuleStateTable(db)

	return db
}

func createFleetRulesTable(db *sql.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS fleet_rules (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		grl_content TEXT NOT NULL,
		audit_manifest TEXT NOT NULL,
		active BOOLEAN DEFAULT false,
		priority INT DEFAULT 100,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		created_by VARCHAR(50),
		INDEX idx_active_priority (active, priority DESC),
		UNIQUE KEY unique_name (name)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	if _, err := db.Exec(schema); err != nil {
		log.Printf("❌ Error creating fleet_rules table: %v", err)
		return
	}
	log.Println("✅ fleet_rules table verified/created")
}

func verifyTables(db *sql.DB) {
	query := `SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE()`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("⚠️  Error checking tables: %v", err)
		return
	}
	defer rows.Close()

	var dbName string
	db.QueryRow("SELECT DATABASE()").Scan(&dbName)
	log.Printf("📁 Tables in database '%s'", dbName)
}

func createVehicleRuleStateTable(db *sql.DB) {
	log.Println("Creating vehicle_rule_state table...")
	schema := `
	CREATE TABLE IF NOT EXISTS vehicle_rule_state (
		imei VARCHAR(20) NOT NULL,
		key_name VARCHAR(100) NOT NULL,
		value_text TEXT NULL,
		value_int BIGINT NULL,
		value_float DOUBLE NULL,
		value_time DATETIME(6) NULL,
		updated_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
		PRIMARY KEY (imei, key_name),
		INDEX idx_updated (updated_at),
		INDEX idx_imei (imei)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	if _, err := db.Exec(schema); err != nil {
		log.Printf("❌ Error creating vehicle_rule_state table: %v", err)
		return
	}
	log.Println("📋 vehicle_rule_state table verified/created")
}

func getMySQLDSN() string {
	host := getEnvWithDefault("MYSQL_HOST", "127.0.0.1")
	port := getEnvWithDefault("MYSQL_PORT", "3306")
	user := getEnvWithDefault("MYSQL_USER", "gpscontrol")
	pass := getEnvWithDefault("MYSQL_PASS", "qazwsxedc")
	dbname := getEnvWithDefault("MYSQL_DB", "grule")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4", user, pass, host, port, dbname)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (s *MySQLStateStore) GetDB() *sql.DB {
	return s.db
}

func (s *MySQLStateStore) GetString(imei, key string) (string, error) {
	var value string
	query := "SELECT value_text FROM vehicle_rule_state WHERE imei = ? AND key_name = ?"
	err := s.db.QueryRow(query, imei, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (s *MySQLStateStore) SetString(imei, key, value string) error {
	query := `
		INSERT INTO vehicle_rule_state (imei, key_name, value_text, updated_at)
		VALUES (?, ?, ?, NOW(6))
		ON DUPLICATE KEY UPDATE value_text = ?, updated_at = NOW(6)
	`
	_, err := s.db.Exec(query, imei, key, value, value)
	return err
}

func (s *MySQLStateStore) GetInt64(imei, key string) (int64, error) {
	var value int64
	query := "SELECT value_int FROM vehicle_rule_state WHERE imei = ? AND key_name = ?"
	err := s.db.QueryRow(query, imei, key).Scan(&value)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return value, err
}

func (s *MySQLStateStore) SetInt64(imei, key string, value int64) error {
	query := `
		INSERT INTO vehicle_rule_state (imei, key_name, value_int, updated_at)
		VALUES (?, ?, ?, NOW(6))
		ON DUPLICATE KEY UPDATE value_int = ?, updated_at = NOW(6)
	`
	_, err := s.db.Exec(query, imei, key, value, value)
	return err
}

func (s *MySQLStateStore) GetTime(imei, key string) (time.Time, error) {
	var value time.Time
	query := "SELECT value_time FROM vehicle_rule_state WHERE imei = ? AND key_name = ?"
	err := s.db.QueryRow(query, imei, key).Scan(&value)
	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}
	return value, err
}

func (s *MySQLStateStore) SetTime(imei, key string, value time.Time) error {
	query := `
		INSERT INTO vehicle_rule_state (imei, key_name, value_time, updated_at)
		VALUES (?, ?, ?, NOW(6))
		ON DUPLICATE KEY UPDATE value_time = ?, updated_at = NOW(6)
	`
	_, err := s.db.Exec(query, imei, key, value, value)
	return err
}

func (s *MySQLStateStore) GetGeofencesByGroup(groupName string) ([]Geofence, error) {
	query := `
		SELECT s.id, s.name, s.shapeType, s.centerLat, s.centerLon, s.radius,
		       s.boundingBoxMinX, s.boundingBoxMaxX, s.boundingBoxMinY, s.boundingBoxMaxY
		FROM geofences.geofences s
		JOIN geofences.geofence_group_mapping m ON s.id = m.geofence_id
		JOIN geofences.geofence_groups gg ON m.group_id = gg.id
		WHERE gg.name = ?
		ORDER BY s.id
	`

	rows, err := s.db.Query(query, groupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var geofences []Geofence
	for rows.Next() {
		var g Geofence
		if err := rows.Scan(&g.ID, &g.Name, &g.ShapeType, &g.CenterLat, &g.CenterLon, &g.Radius,
			&g.BoundingBoxMinX, &g.BoundingBoxMaxX, &g.BoundingBoxMinY, &g.BoundingBoxMaxY); err != nil {
			log.Printf("⚠️  Error scanning geofence: %v", err)
			continue
		}
		geofences = append(geofences, g)
	}
	return geofences, nil
}
