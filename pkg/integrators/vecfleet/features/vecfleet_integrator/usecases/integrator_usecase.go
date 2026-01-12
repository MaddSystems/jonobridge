package usecases

import (
	"fmt"
	"vecfleet/features/vecfleet_integrator/data"
)

func RunIntegrator(imei string, dataStr string) error {
	plates, err := GetPlates(imei)
	if err != nil {
		return err
	}

	//HTTP SEND DATA
	err = data.SendRequest(plates, dataStr)
	if err != nil {
		return fmt.Errorf("error al enviar solicitud: %w", err)
	}

	return nil
}
