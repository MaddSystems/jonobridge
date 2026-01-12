package gruasayaintegrator

import (
	"gruasaya/features/gruasaya_integrator/data"
	"gruasaya/features/gruasaya_integrator/usecases"
	"encoding/json"
	"fmt"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

// Initialize sets up any required resources for the Gruasaya integrator
func Initialize() {
	// Initialize Gruasaya
	data.InitGruasaya()
	utils.VPrint("Gruasaya integrator initialized")
}

// Init processes incoming data for AltoTrack integration
func Init(jsonData string) error {
	var data models.JonoModel

	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		utils.VPrint("Error unmarshalling JSON: %v", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	if data.IMEI == "" {
		utils.VPrint("Warning: Empty IMEI in received data")
	}

	return usecases.RunIntegrator(data.IMEI, jsonData)
}
