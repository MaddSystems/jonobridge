package movupintegrator

import (
	"movup/features/movup_integrator/data"
	"movup/features/movup_integrator/usecases"
	"encoding/json"
	"fmt"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

// Initialize sets up any required resources for the Movup integrator
func Initialize() {
	// Initialize Movup
	data.InitMovup()
	utils.VPrint("Movup integrator initialized")
}

// Init processes incoming data for AltoTrack integration
func Init(jsonData string) error {
	// Add panic recovery to prevent crashes
	defer func() {
		if r := recover(); r != nil {
			utils.VPrint("Recovered from panic in Init: %v", r)
		}
	}()
	
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
