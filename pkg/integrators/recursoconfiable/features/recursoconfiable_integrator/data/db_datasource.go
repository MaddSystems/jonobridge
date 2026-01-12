package data

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}
	return db, nil
}

func GetPlates(db *sql.DB, imei string) (string, error) {

	var plates sql.NullString
	err := db.QueryRow("SELECT plates FROM devices WHERE imei=?", imei).Scan(&plates)

	if err == sql.ErrNoRows {
		return "", nil // No existe el dispositivo
	} else if err != nil {
		return "", fmt.Errorf("error en la consulta: %w", err)
	}
	return plates.String, nil
}

func InsertDevice(db *sql.DB, imei, date, eventCode, latitude, longitude, altitude, speed, angle, protocol string) error {
	cmd := "INSERT INTO devices (imei, creation_date, plates, event_code, latitude, longitude, altitude, speed, angle, protocol) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(cmd, imei, date, "", eventCode, latitude, longitude, altitude, speed, angle, protocol)
	if err != nil {
		return fmt.Errorf("error al insertar: %s", err)
	}
	return nil
}
