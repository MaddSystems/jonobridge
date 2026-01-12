package usecases

import (
	"fmt"
	"unigis/features/unigis_integrator/data"

	"github.com/MaddSystems/jonobridge/common/utils"
)

func RunIntegrator(imei, dataStr string) error {
	utils.VPrint("IMEI: %s", imei)

	// Get the vehicle plates for this IMEI
	platesStr, err := utils.GetPlates(imei)
	if err != nil {
		return fmt.Errorf("failed to get plates: %w", err)
	}

	utils.VPrint("Plates: %s", platesStr)
	fmt.Println(platesStr)

	// Send the data via SOAP
	if err := data.SendSOAPRequest(platesStr, dataStr); err != nil {
		return fmt.Errorf("error al enviar solicitud SOAP: %w", err)
	}

	return nil
}
