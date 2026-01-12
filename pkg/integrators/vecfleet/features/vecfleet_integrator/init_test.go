package vecfleetintegrator_test

import (
	"testing"
	vecfleetintegrator "vecfleet/features/vecfleet_integrator"
)

func TestInitialize(t *testing.T) {
	jsonData := `{
		"IMEI": "866811062543247",
		"DataPackets": 1,
		"ListPackets": {
			"packet_1": {
				"Altitude": 0,
				"Datetime": "2024-11-01T16:35:47Z",
				"EventCode": {"Code": 35, "Name": "Track By Time Interval"},
				"Latitude": 19.52101,
				"Longitude": -99.211608,
				"Speed": 0,
				"Direction": 0,
				"GsmSignalStrength": 8,
				"Hdop": 0.0,
				"Mileage": 0,
				"NumberOfSatellites": 0,
				"PositioningStatus": "A",
				"RunTime": 250632
			}
		}
	}`

	err := vecfleetintegrator.Init(jsonData)
	if err != nil {
		t.Errorf("Expected nil, but got error: %v", err)
	} else {
		t.Logf("Test passed successfully with no errors.")
	}
}
