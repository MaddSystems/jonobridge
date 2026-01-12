package avocadocloudintegrator

import (
	"avocadocloud/features/avocadocloud_integrator/data"
	"avocadocloud/features/avocadocloud_integrator/usecases"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

// Initialize sets up any required resources for the Avocadocloud integrator
func Initialize(db *sql.DB) {
	// Initialize Avocadocloud with database connection
	data.InitAvocadocloud(db)
	utils.VPrint("Avocadocloud integrator initialized")
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
