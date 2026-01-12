package usecases

import (
	"fmt"
	"motumcloud/features/motumcloud_integrator/data"

	"github.com/MaddSystems/jonobridge/common/utils"
)

func RunIntegrator(imei, dataStr string) error {
	// First get the vehicle plates for this IMEI from our JSON data

	platesStr, err := utils.GetPlates(imei)
	if err != nil {
		platesStr = "Desconocido" // Use a default value instead of empty string
	}
	ecoStr, err := utils.GetEco(imei)
	if err != nil {
		ecoStr = "Desconocido"
	}
	vinStr, err := utils.GetVin(imei)
	if err != nil {
		vinStr = "Desconocido"
	}
	if err := data.ProcessAndSendMotumcloud(platesStr, ecoStr, vinStr, dataStr); err != nil {
		utils.VPrint("Error processing and sending data: %v", err)
		return fmt.Errorf("error processing and sending data: %w", err)
	}

	return nil
}
