package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// JonoModel represents the JSON structure to be printed
type JonoModel struct {
	IMEI        string                `json:"IMEI"`
	Message     *string               `json:"Message"`
	DataPackets int                   `json:"DataPackets"`
	ListPackets map[string]DataPacket `json:"ListPackets"`
}

type DataPacket struct {
	Altitude                     int                         `json:"Altitude"`
	Datetime                     time.Time                   `json:"Datetime"`
	EventCode                    EventCode                   `json:"EventCode"`
	Latitude                     float64                     `json:"Latitude"`
	Longitude                    float64                     `json:"Longitude"`
	Speed                        int                         `json:"Speed"`
	RunTime                      int                         `json:"RunTime"`
	FuelPercentage               int                         `json:"FuelPercentage"`
	Direction                    int                         `json:"Direction"`
	HDOP                         float64                     `json:"HDOP"`
	Mileage                      int                         `json:"Mileage"`
	PositioningStatus            string                      `json:"PositioningStatus"`
	NumberOfSatellites           int                         `json:"NumberOfSatellites"`
	GSMSignalStrength            *int                        `json:"GSMSignalStrength"`
	AnalogInputs                 *AnalogInputs               `json:"AnalogInputs"`
	IoPortStatus                 *IoPortsStatus              `json:"IoPortStatus"`
	BaseStationInfo              *BaseStationInfo            `json:"BaseStationInfo"`
	OutputPortStatus             *OutputPortStatus           `json:"OutputPortStatus"`
	InputPortStatus              *InputPortStatus            `json:"InputPortStatus"`
	SystemFlag                   *SystemFlag                 `json:"SystemFlag"`
	TemperatureSensor            *TemperatureSensor          `json:"TemperatureSensor"`
	CameraStatus                 *CameraStatus               `json:"CameraStatus"`
	CurrentNetworkInfo           *CurrentNetworkInfo         `json:"CurrentNetworkInfo"`
	FatigueDrivingInformation    *FatigueDrivingInformation  `json:"FatigueDrivingInformation"`
	AdditionalAlertInfoADASDMS   *AdditionalAlertInfoADASDMS `json:"AdditionalAlertInfoADASDMS"`
	BluetoothBeaconA             *BluetoothBeacon            `json:"BluetoothBeaconA"`
	BluetoothBeaconB             *BluetoothBeacon            `json:"BluetoothBeaconB"`
	TemperatureAndHumiditySensor *TemperatureAndHumidity     `json:"TemperatureAndHumiditySensor"`
	DriverStatus                 *DriverStatus               `json:"DriverStatus"`
	ICCardStatus                 *ICCardStatus               `json:"ICCardStatus"`
	LocationTracking             *LocationTracking           `json:"LocationTracking"`
	TerminalUpgrade              *TerminalUpgrade            `json:"TerminalUpgrade"`
	BatteryStatus                *BatteryStatus              `json:"BatteryStatus"`
	TirePressure                 *TirePressure               `json:"TirePressure"`
	FuelConsumption              *FuelConsumption            `json:"FuelConsumption"`
	EngineStatus                 *EngineStatus               `json:"EngineStatus"`
	VehicleStatus                *VehicleStatus              `json:"VehicleStatus"`
	AlarmStatus                  *AlarmStatus                `json:"AlarmStatus"`
	AdditionalInfo               *AdditionalInfo             `json:"AdditionalInfo"`
}

type EventCode struct {
	Code int    `json:"Code"`
	Name string `json:"Name"`
}

type BaseStationInfo struct {
	MCC    *string `json:"MCC"`
	MNC    *string `json:"MNC"`
	LAC    *string `json:"LAC"`
	CellID *string `json:"CellID"`
}

type AnalogInputs struct {
	AD1  *string `json:"AD1"`
	AD2  *string `json:"AD2"`
	AD3  *string `json:"AD3"`
	AD4  *string `json:"AD4"`
	AD5  *string `json:"AD5"`
	AD6  *string `json:"AD6"`
	AD7  *string `json:"AD7"`
	AD8  *string `json:"AD8"`
	AD9  *string `json:"AD9"`
	AD10 *string `json:"AD10"`
}

type OutputPortStatus struct {
	Output1 *string `json:"Output1"`
	Output2 *string `json:"Output2"`
	Output3 *string `json:"Output3"`
	Output4 *string `json:"Output4"`
	Output5 *string `json:"Output5"`
	Output6 *string `json:"Output6"`
	Output7 *string `json:"Output7"`
	Output8 *string `json:"Output8"`
}

type InputPortStatus struct {
	Input1 *string `json:"Input1"`
	Input2 *string `json:"Input2"`
	Input3 *string `json:"Input3"`
	Input4 *string `json:"Input4"`
	Input5 *string `json:"Input5"`
	Input6 *string `json:"Input6"`
	Input7 *string `json:"Input7"`
	Input8 *string `json:"Input8"`
}

type SystemFlag struct {
	EEP2                *string `json:"EEP2"`
	ACC                 *string `json:"ACC"`
	AntiTheft           *string `json:"AntiTheft"`
	VibrationFlag       *string `json:"VibrationFlag"`
	MovingFlag          *string `json:"MovingFlag"`
	ExternalPowerSupply *string `json:"ExternalPowerSupply"`
	Charging            *string `json:"Charging"`
	SleepMode           *string `json:"SleepMode"`
	FMS                 *string `json:"FMS"`
	FMSFunction         *string `json:"FMSFunction"`
	SystemFlagExtras    *string `json:"SystemFlagExtras"`
}

type TemperatureSensor struct {
	SensorNumber *string `json:"SensorNumber"`
	Value        *string `json:"Value"`
}

type CameraStatus struct {
	CameraNumber *string `json:"CameraNumber"`
	Status       *string `json:"Status"`
}

type CurrentNetworkInfo struct {
	Version    *string `json:"Version"`
	Type       *string `json:"Type"`
	Descriptor *string `json:"Descriptor"`
}

type FatigueDrivingInformation struct {
	Version    *string `json:"Version"`
	Type       *string `json:"Type"`
	Descriptor *string `json:"Descriptor"`
}

type AdditionalAlertInfoADASDMS struct {
	AlarmProtocol *string `json:"AlarmProtocol"`
	AlarmType     *string `json:"AlarmType"`
	PhotoName     *string `json:"PhotoName"`
}

type BluetoothBeacon struct {
	Version        *string `json:"Version"`
	DeviceName     *string `json:"DeviceName"`
	MAC            *string `json:"MAC"`
	BatteryPower   *string `json:"BatteryPower"`
	SignalStrength *string `json:"SignalStrength"`
}

type TemperatureAndHumidity struct {
	DeviceName           *string `json:"DeviceName"`
	MAC                  *string `json:"MAC"`
	BatteryPower         *string `json:"BatteryPower"`
	Temperature          *string `json:"Temperature"`
	Humidity             *string `json:"Humidity"`
	AlertHighTemperature *string `json:"AlertHighTemperature"`
	AlertLowTemperature  *string `json:"AlertLowTemperature"`
	AlertHighHumidity    *string `json:"AlertHighHumidity"`
	AlertLowHumidity     *string `json:"AlertLowHumidity"`
}

type IoPortsStatus struct {
	Port1 int `json:"Port1"`
	Port2 int `json:"Port2"`
	Port3 int `json:"Port3"`
	Port4 int `json:"Port4"`
	Port5 int `json:"Port5"`
	Port6 int `json:"Port6"`
	Port7 int `json:"Port7"`
	Port8 int `json:"Port8"`
}

type DriverStatus struct {
	Status        *string `json:"Status"`       // 0x01: On duty, 0x02: Off duty
	ICCardStatus  *string `json:"ICCardStatus"` // 0x00: Success, 0x01: Authentication failed, 0x02: Card locked, etc.
	DriverName    *string `json:"DriverName"`
	CertCode      *string `json:"CertCode"`
	CertAuthority *string `json:"CertAuthority"`
	CertValidity  *string `json:"CertValidity"`
}

type ICCardStatus struct {
	Status        *string `json:"Status"`        // 0x01: Card inserted, 0x02: Card pulled out
	ReadingResult *string `json:"ReadingResult"` // 0x00: Success, 0x01: Authentication failed, etc.
}

type LocationTracking struct {
	TimeInterval    *int `json:"TimeInterval"`    // Time interval in seconds
	ValidTime       *int `json:"ValidTime"`       // Valid time in seconds
	LocationReports *int `json:"LocationReports"` // Number of location reports
}

type TerminalUpgrade struct {
	UpgradeType   *string `json:"UpgradeType"`   // 0x00: Terminal, 0x12: IC card reader, etc.
	UpgradeResult *string `json:"UpgradeResult"` // 0x00: Success, 0x01: Failure, etc.
}

type BatteryStatus struct {
	Voltage       *float64 `json:"Voltage"`       // Battery voltage
	ChargeStatus  *string  `json:"ChargeStatus"`  // Charging status
	BatteryHealth *string  `json:"BatteryHealth"` // Battery health status
}

type TirePressure struct {
	FrontLeft  *float64 `json:"FrontLeft"`  // Front left tire pressure
	FrontRight *float64 `json:"FrontRight"` // Front right tire pressure
	RearLeft   *float64 `json:"RearLeft"`   // Rear left tire pressure
	RearRight  *float64 `json:"RearRight"`  // Rear right tire pressure
}

type FuelConsumption struct {
	TotalFuelConsumed *float64 `json:"TotalFuelConsumed"` // Total fuel consumed in ml
	FuelRate          *float64 `json:"FuelRate"`          // Fuel rate in L/100km
}

type EngineStatus struct {
	EngineRPM        *int     `json:"EngineRPM"`        // Engine RPM
	CoolantTemp      *int     `json:"CoolantTemp"`      // Coolant temperature in °C
	OilPressure      *float64 `json:"OilPressure"`      // Oil pressure in kPa
	ThrottlePosition *float64 `json:"ThrottlePosition"` // Throttle position in %
}

type VehicleStatus struct {
	ACCStatus        *string `json:"ACCStatus"`        // ACC on/off
	PositionStatus   *string `json:"PositionStatus"`   // Positioned/Not positioned
	OilCircuitStatus *string `json:"OilCircuitStatus"` // Oil circuit status
	DoorStatus       *string `json:"DoorStatus"`       // Door status
}

type AlarmStatus struct {
	EmergencyAlarm *string `json:"EmergencyAlarm"` // Emergency alarm status
	OverSpeedAlarm *string `json:"OverSpeedAlarm"` // Over-speed alarm status
	FatigueDriving *string `json:"FatigueDriving"` // Fatigue driving alarm status
	GNSSFailure    *string `json:"GNSSFailure"`    // GNSS module failure
}

type AdditionalInfo struct {
	RotationSpeed   *int    `json:"RotationSpeed"`   // Rotation speed in RPM
	BaseStationData *string `json:"BaseStationData"` // Base station data
	WifiData        *string `json:"WifiData"`        // Wifi data
}

// Decode JT/T 808 response
func decodeTrackerResponse(data []byte) JonoModel {
	jonoModel := JonoModel{
		ListPackets: make(map[string]DataPacket),
	}

	if len(data) < 12 {
		log.Println("Invalid packet length")
		return jonoModel
	}

	// Ensure start/end bytes are correct
	if data[0] != 0x7E || data[len(data)-1] != 0x7E {
		log.Println("Invalid packet format")
		return jonoModel
	}

	// Extract header fields
	messageID := binary.BigEndian.Uint16(data[1:3])
	attribute := binary.BigEndian.Uint16(data[3:5])
	terminalID := fmt.Sprintf("%X", data[5:11])
	messageSerial := binary.BigEndian.Uint16(data[11:13])

	jonoModel.IMEI = terminalID

	// Check if the message is fragmented (packetized)
	isFragmented := (attribute & 0x2000) != 0
	if isFragmented {
		totalPackets := binary.BigEndian.Uint16(data[13:15])
		currentPacket := binary.BigEndian.Uint16(data[15:17])
		fmt.Printf("Fragmented message: Packet %d of %d\n", currentPacket, totalPackets)
	}

	// Extract location data (0200 message)
	if messageID == 0x0200 {
		status := binary.BigEndian.Uint32(data[13:17])
		latitude := binary.BigEndian.Uint32(data[17:21])  // GPS latitude
		longitude := binary.BigEndian.Uint32(data[21:25]) // GPS longitude
		speed := binary.BigEndian.Uint16(data[25:27])
		direction := binary.BigEndian.Uint16(data[27:29])
		timestamp := fmt.Sprintf("%02X-%02X-%02X %02X:%02X:%02X",
			data[29], data[30], data[31], data[32], data[33], data[34])

		dataPacket := DataPacket{
			Latitude:  float64(latitude) / 1000000,
			Longitude: float64(longitude) / 1000000,
			Speed:     int(speed) / 10,
			Direction: int(direction),
			Datetime:  time.Now(), // Replace with actual timestamp parsing
		}

		jonoModel.ListPackets[timestamp] = dataPacket

		// Decode additional information items
		if len(data) > 35 {
			additionalInfo := data[35:]
			for len(additionalInfo) > 2 {
				infoID := additionalInfo[0]
				infoLength := int(additionalInfo[1])
				if len(additionalInfo) < 2+infoLength {
					break
				}
				infoData := additionalInfo[2 : 2+infoLength]
				fmt.Printf("Additional Info ID: 0x%02X, Length: %d, Data: %X\n", infoID, infoLength, infoData)
				additionalInfo = additionalInfo[2+infoLength:]
			}
		}
	}

	return jonoModel
}

// Create general response message (8001)
func createGeneralResponse(messageID uint16, messageSerial uint16, terminalID []byte, result uint8) []byte {
	// Construct response according to JT/T 808 protocol
	responseData := make([]byte, 0, 16)

	// Start marker
	responseData = append(responseData, 0x7E)

	// Message ID (8001 - general response)
	responseData = append(responseData, 0x80, 0x01)

	// Message attributes (length to be filled later)
	responseData = append(responseData, 0x00, 0x05)

	// Terminal ID (copy from request)
	responseData = append(responseData, terminalID...)

	// Message serial
	serialBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(serialBytes, messageSerial)
	responseData = append(responseData, serialBytes...)

	// Response message serial
	responseSerialBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(responseSerialBytes, messageSerial)
	responseData = append(responseData, responseSerialBytes...)

	// Original message ID
	msgIDBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(msgIDBytes, messageID)
	responseData = append(responseData, msgIDBytes...)

	// Result (0: success)
	responseData = append(responseData, result)

	// Calculate checksum (XOR of all bytes except start and end markers)
	var checksum byte
	for i := 1; i < len(responseData); i++ {
		checksum ^= responseData[i]
	}
	responseData = append(responseData, checksum)

	// End marker
	responseData = append(responseData, 0x7E)

	return responseData
}

// Find complete messages in buffer (packets may be split or combined)
func findCompleteMessages(buffer []byte) ([][]byte, []byte) {
	var messages [][]byte
	remaining := buffer

	for {
		// Find start marker
		startIndex := -1
		for i, b := range remaining {
			if b == 0x7E {
				startIndex = i
				break
			}
		}

		if startIndex == -1 {
			// No start marker found, keep all data
			return messages, remaining
		}

		// Discard data before start marker
		if startIndex > 0 {
			remaining = remaining[startIndex:]
			startIndex = 0
		}

		// Find end marker
		endIndex := -1
		for i := startIndex + 1; i < len(remaining); i++ {
			if remaining[i] == 0x7E {
				endIndex = i
				break
			}
		}

		if endIndex == -1 {
			// No end marker found, wait for more data
			return messages, remaining
		}

		// Extract complete message
		message := remaining[startIndex : endIndex+1]
		messages = append(messages, message)

		// Update remaining buffer
		remaining = remaining[endIndex+1:]

		// If no more data, exit loop
		if len(remaining) == 0 {
			break
		}
	}

	return messages, remaining
}

// Handle connection from tracker
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Log client information
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("New connection from %s", remoteAddr)

	var buffer []byte
	readBuffer := make([]byte, 4096)

	// Set read deadline
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	if err != nil {
		log.Printf("Error setting read deadline: %v", err)
		return
	}

	for {
		// Read data from connection
		n, err := conn.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by %s", remoteAddr)
			} else {
				log.Printf("Error reading from %s: %v", remoteAddr, err)
			}
			break
		}

		// Append new data to buffer
		buffer = append(buffer, readBuffer[:n]...)

		// Find complete messages
		messages, remaining := findCompleteMessages(buffer)
		buffer = remaining

		// Process complete messages
		for _, message := range messages {
			log.Printf("Received from %s: %s", remoteAddr, hex.EncodeToString(message))

			// Decode message and populate JonoModel
			jonoModel := decodeTrackerResponse(message)

			// Print JonoModel as JSON
			jsonData, err := json.MarshalIndent(jonoModel, "", "  ")
			if err != nil {
				log.Printf("Error marshaling JonoModel to JSON: %v", err)
			} else {
				fmt.Println(string(jsonData))
			}

			// Check if message has valid format
			if len(message) >= 13 && message[0] == 0x7E && message[len(message)-1] == 0x7E {
				// Extract message ID and serial
				messageID := binary.BigEndian.Uint16(message[1:3])
				messageSerial := binary.BigEndian.Uint16(message[11:13])
				terminalID := message[5:11]

				// Send response based on message type
				var response []byte

				switch messageID {
				case 0x0100: // Registration
					log.Printf("Registration request from %s", remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				case 0x0200: // Location report
					log.Printf("Location report from %s", remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				case 0x0002: // Heartbeat
					log.Printf("Heartbeat from %s", remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				case 0x0102: // Terminal authentication
					log.Printf("Terminal authentication from %s", remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				case 0x8105: // Terminal control
					log.Printf("Terminal control from %s", remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				default:
					log.Printf("Unknown message type %04X from %s", messageID, remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				}

				// Send response
				if response != nil {
					_, err := conn.Write(response)
					if err != nil {
						log.Printf("Error sending response to %s: %v", remoteAddr, err)
						break
					}
					log.Printf("Sent response to %s: %s", remoteAddr, hex.EncodeToString(response))
				}
			}
		}

		// Reset read deadline
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		if err != nil {
			log.Printf("Error resetting read deadline: %v", err)
			break
		}
	}
}

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "0.0.0.0:8600")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	log.Println("Server started on 0.0.0.0:8600")

	// Track active connections
	var wg sync.WaitGroup

	for {
		// Accept connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Handle connection in goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConnection(conn)
		}()
	}

	// Wait for all connections to finish
	wg.Wait()
}
