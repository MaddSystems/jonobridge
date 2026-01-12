package recursoconfiable_integrator_test

import (
	"recursoconfiable/features/recursoconfiable_integrator"
	"testing"
)

func TestInitialize(t *testing.T) {
	jsonData := `{
		"IMEI": "867869060389442",
		"DataPackets": 1,
		"ListPackets": {
			"packet_1": {
				"Altitude": 0,
				"Datetime": "2024-11-01T16:35:47Z",
				"EventCode": {"Code": 35, "Name": "Track By Time Interval"},
				"Latitude": 19.52101,
				"Longitude": -99.211608,
				"Speed": 0,
				"Extras": {
					"Direction": 0,
					"GsmSignalStrength": 8,
					"Hdop": 0,
					"Mileage": 0,
					"NumberOfSatellites": 0,
					"PositioningStatus": false,
					"RunTime": 250632,
					"SystemFlag": false
				}
			}
		}
	}`

	err := recursoconfiable_integrator.Initialize(jsonData)
	if err != nil {
		t.Errorf("Expected nil, but got error: %v", err)
	} else {
		t.Logf("Test passed successfully with no errors.")
	}
}
