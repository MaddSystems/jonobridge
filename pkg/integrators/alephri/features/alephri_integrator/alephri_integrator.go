package alephriintegrator

import (
	"alephri/features/alephri_integrator/data"
	"alephri/features/alephri_integrator/usecases"
	"encoding/json"
	"fmt"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

// Initialize sets up any required resources for the Alephri integrator
func Initialize() {
	// Initialize Alephri
	data.InitAlephri()
	utils.VPrint("Alephri integrator initialized")
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
