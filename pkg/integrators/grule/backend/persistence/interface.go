package persistence

import (
	"database/sql"
	"time"
)

type Geofence struct {
	ID              int
	Name            string
	ShapeType       string
	CenterLat       sql.NullFloat64
	CenterLon       sql.NullFloat64
	Radius          sql.NullFloat64
	BoundingBoxMinX sql.NullFloat64
	BoundingBoxMaxX sql.NullFloat64
	BoundingBoxMinY sql.NullFloat64
	BoundingBoxMaxY sql.NullFloat64
}

type StateStore interface {
	GetString(imei, key string) (string, error)
	SetString(imei, key, value string) error
	GetInt64(imei, key string) (int64, error)
	SetInt64(imei, key string, value int64) error
	GetTime(imei, key string) (time.Time, error)
	SetTime(imei, key string, value time.Time) error
	GetGeofencesByGroup(groupName string) ([]Geofence, error)
}