package send2mysqlintegrator

import (
	"database/sql"
	"send2mysql/features/send2mysql_integrator/data"
	"send2mysql/features/send2mysql_integrator/usecases"
	"encoding/json"
	"fmt"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

// Initialize sets up any required resources for the Send2mysql integrator
func Initialize(db *sql.DB) {
	// Initialize Send2mysql data package with the database connection
	data.SetDB(db)
	// Call InitSend2mysql for any other non-DB initializations in the data package (e.g., httpTimeout).
	data.InitSend2mysql() 
	utils.VPrint("Send2mysql integrator initialized with DB connection and other settings")
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
