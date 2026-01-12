package geofence

import (
	"log"
	"math"
)

func (c *GeofenceCapability) IsInsideCircle(lat, lon, clat, clon, radius float64) bool {
	// Haversine distance in meters
	const earthRadius = 6371000.0 // meters
	dLat := (clat - lat) * math.Pi / 180.0
	dLon := (clon - lon) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat*math.Pi/180.0)*math.Cos(clat*math.Pi/180.0)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c2 := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c2
	return distance <= radius
}

func (c *GeofenceCapability) IsInsideGroup(groupName string, latitude, longitude float64) bool {
	geofences, err := c.store.GetGeofencesByGroup(groupName)
	if err != nil {
		log.Printf("❌ Error querying geofences for group %s: %v", groupName, err)
		return false
	}

	for _, g := range geofences {
		// Check Circle
		if (g.ShapeType == "circle" || g.ShapeType == "Circle") && g.CenterLat.Valid && g.CenterLon.Valid && g.Radius.Valid {
			if c.IsInsideCircle(latitude, longitude, g.CenterLat.Float64, g.CenterLon.Float64, g.Radius.Float64) {
				return true
			}
		}

		// Check Polygon (Bounding Box approximation as per original)
		if (g.ShapeType == "polygon" || g.ShapeType == "Polygon") && g.BoundingBoxMinX.Valid && g.BoundingBoxMaxX.Valid && g.BoundingBoxMinY.Valid && g.BoundingBoxMaxY.Valid {
			if latitude >= g.BoundingBoxMinY.Float64 && latitude <= g.BoundingBoxMaxY.Float64 &&
				longitude >= g.BoundingBoxMinX.Float64 && longitude <= g.BoundingBoxMaxX.Float64 {
				return true
			}
		}
	}
	return false
}
