package usecases

import (
	"fmt"
	"gruasaya/features/gruasaya_integrator/data"

	"github.com/MaddSystems/jonobridge/common/utils"
)

func RunIntegrator(imei, dataStr string) error {
	// Send the data via SOAP with error handling
	if err := data.ProcessAndSendGruasaya(dataStr); err != nil {
		utils.VPrint("Error processing and sending data: %v", err)
		return fmt.Errorf("error processing and sending data: %w", err)
	}

	return nil
}
