package usecases

import (
	"gpsinfo/features/gpsinfo_integrator/data"
	"fmt"

	"github.com/MaddSystems/jonobridge/common/utils"
)

func RunIntegrator(imei, dataStr string) error {
	// First get the vehicle plates for this IMEI from our JSON data

	platesStr, err := utils.GetPlates(imei)
	if err != nil {
		platesStr = "Desconocido" // Use a default value instead of empty string
	}
	//utils.VPrint("Plates: %s", platesStr)
	ecoStr, err := utils.GetEco(imei)
	if err != nil {
		ecoStr = "Desconocido" // Use a default value instead of empty string
	}

	vinStr, err := utils.GetVin(imei)
	if err != nil {
		vinStr = "Desconocido" // Use a default value instead of empty string
	}

	// Send the data via SOAP with error handling
	if err := data.ProcessAndSendGpsinfo(platesStr, ecoStr, vinStr, dataStr); err != nil {
		utils.VPrint("Error processing and sending data: %v", err)
		return fmt.Errorf("error processing and sending data: %w", err)
	}

	return nil
}
