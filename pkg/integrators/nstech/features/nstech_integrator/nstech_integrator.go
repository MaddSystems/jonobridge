package nstechintegrator

import (
	"encoding/json"
	"nstech/features/nstech_integrator/data"
	"nstech/features/nstech_integrator/usecases"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

// Initialize sets up any required resources for the Nstech integrator
func Initialize() {
	// Initialize Nstech
	data.InitNstech()
	utils.VPrint("Nstech integrator initialized")
}

// Init processes incoming data for AltoTrack integration
func Init(jsonData string) error {
	var data models.JonoModel

	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		utils.VPrint("Error unmarshalling JSON (non-fatal): %v", err)
		return nil // Don't propagate errors that could crash
	}

	if data.IMEI == "" {
		utils.VPrint("Warning: Empty IMEI in received data")
	}

	// Always return nil to prevent crashes
	usecases.RunIntegrator(data.IMEI, jsonData)
	return nil
}
