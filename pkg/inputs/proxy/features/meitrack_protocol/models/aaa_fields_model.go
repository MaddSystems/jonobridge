package models

type AAAModel struct {
	GeneralModel
	EventCode          any
	Latitude           float64
	Longitude          float64
	Datetime           string
	PositioningStatus  string
	NumberOfSatellites int
	GsmSignalStrength  int
	Speed              int
	Direction          int
	Hdop               float64
	Altitude           float64
	Mileage            int
	RunTime            int
	BaseStationInfo    any
	IoPortStatus       string
	AnalogInputs       any
	AssistedEventInfo  string
	CustomizedData     string
	ProtocolVersion    int
	FuelPercentage     string
	TemperatureSensor  string
	MaxAcceleration    int
	MaxDesceleration   int
	Checksum           string
}
