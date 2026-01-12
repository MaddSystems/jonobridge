package geofence

import (
	"sync"

	"github.com/jonobridge/grule-backend/persistence"
)

type coords struct {
	lat float64
	lon float64
}

type GeofenceCapability struct {
	store      persistence.StateStore
	lastCoords sync.Map // imei -> coords
}

func NewGeofenceCapability(store persistence.StateStore) *GeofenceCapability {
	return &GeofenceCapability{
		store: store,
	}
}

func (c *GeofenceCapability) Name() string {
	return "geofence"
}

func (c *GeofenceCapability) Version() string {
	return "1.0.0"
}

func (c *GeofenceCapability) GetDataContextName() string {
	return "geo"
}

func (c *GeofenceCapability) Initialize(imei string) error {
	return nil
}

func (c *GeofenceCapability) GetSnapshot() map[string]interface{} {
	return map[string]interface{}{}
}

// UpdateLastPacket updates the last known coordinates for an IMEI
func (c *GeofenceCapability) UpdateLastPacket(imei string, lat, lon float64) {
	c.lastCoords.Store(imei, coords{lat: lat, lon: lon})
}

// GetSnapshotData implements SnapshotProvider
func (c *GeofenceCapability) GetSnapshotData(imei string) map[string]interface{} {
	if c == nil {
		return nil
	}

	val, ok := c.lastCoords.Load(imei)
	if !ok {
		return nil
	}
	co := val.(coords)

	return map[string]interface{}{
		"geofence_checks": map[string]bool{
			"inside_taller":    c.IsInsideGroup("Taller", co.lat, co.lon),
			"inside_clientes":  c.IsInsideGroup("CLIENTES", co.lat, co.lon),
			"inside_resguardo": c.IsInsideGroup("Resguardo/Cedis/Puerto", co.lat, co.lon),
		},
	}
}
