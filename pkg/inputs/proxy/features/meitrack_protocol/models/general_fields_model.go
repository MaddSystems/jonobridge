package models

import (
	"fmt"
	"strings"
)

type CommandType string

const (
	CommandAAA CommandType = "AAA"
	CommandCCE CommandType = "CCE"
)

type StartSignal string

const (
	StartSignalToServer StartSignal = "$$"
	StartSignalToDevice StartSignal = "@@"
)

type GeneralModel struct {
	StartSignal StartSignal
	Identifier  string
	DataLength  string
	IMEI        string
	CommandType CommandType
	Rest        string
	Message     string
}

func ParseGeneralFields(log string) (GeneralModel, error) {
	fmt.Println(log)
	if len(log) < 25 {
		return GeneralModel{}, fmt.Errorf("error data too short")
	}
	parts := strings.Split(log, ",")
	startSignal := StartSignal(parts[0][:2])
	if startSignal != StartSignalToDevice && startSignal != StartSignalToServer {
		return GeneralModel{}, fmt.Errorf("Invalid start signal")
	}

	return GeneralModel{
		StartSignal: startSignal,
		Identifier:  parts[0][2:3],
		DataLength:  parts[0][3:],
		IMEI:        parts[1],
		CommandType: CommandType(parts[2]),
		Rest:        strings.Join(parts[3:], ","),
		Message:     log,
	}, nil
}
