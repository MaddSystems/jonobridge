package activetrackintegrator

import (
	"activetrack/features/activetrack_integrator/data"
	"activetrack/features/activetrack_integrator/usecases"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

// Initialize sets up any required resources for the ActiveTrack integrator
func Initialize() {
	// Try multiple paths to be container-friendly
	possiblePaths := []string{
		"/data_spoof.json", // Root of container (preferred for container usage)
		filepath.Join(os.Getenv("HOME"), "jonobridge", "pkg", "integrators", "activetrack", "data_spoof.json"), // Development path
	}

	fileExists := false
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			utils.VPrint("Found existing IMEI spoof mapping file at: %s", path)
			fileExists = true
			break
		}
	}

	// If no file exists in any location, create it at the container root
	if !fileExists {
		utils.VPrint("IMEI spoof mapping file doesn't exist, creating at container root...")
		containerPath := "/data_spoof.json"
		if err := utils.FetchAndSaveImeiMappings(containerPath); err != nil {
			utils.VPrint("Failed to create IMEI spoof mapping file: %v", err)
		} else {
			utils.VPrint("Successfully created IMEI spoof mapping file at %s", containerPath)
		}
	}

	// Initialize Activetrac
	data.InitActivetrac()
	utils.VPrint("Activetrack integrator initialized")
}

// Init processes incoming data for ActiveTrack integration
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

// SendAllPendingData sends all pending data for all IMEIs
func SendAllPendingData() {
	data.SendAllPendingData()
}
