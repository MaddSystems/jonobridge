package meitrack_integrator

import (
	"fmt"
	"forwarder/features/meitrack_integrator/usecases"
)

func Initialize(data string) ([]string, error) {
	result, err := usecases.BuildOutput(data)
	if err != nil {
		return nil, fmt.Errorf("error meitrakc integrator")
	}
	return result, nil
}
