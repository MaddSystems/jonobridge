package meitrack_integrator

import (
	"fmt"
	"meitrack/features/meitrack_integrator/usecases"
)

func Initialize(data, meitrack_mock_imei, meitrack_mock_value string) ([]string, bool, error) {
	result, isADASEvent, err := usecases.BuildOutput(data, meitrack_mock_imei, meitrack_mock_value)
	if err != nil {
		return nil, false, fmt.Errorf("error meitrack integrator: %w", err)
	}
	return result, isADASEvent, nil
}
