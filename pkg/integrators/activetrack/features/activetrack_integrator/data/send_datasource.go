package data

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

// Uses a global lastGlobalSend variable to track when any data was last sent
// Sends data when either:
// The buffer for a specific IMEI reaches 50 entries
// It's been more than 1 minute since any data was sent
// Adds a new SendAllPendingData() function that could be called periodically to flush all buffers
// Keeps the same buffer management logic, removing only the entries that were sent

var payloadBuffer = make(map[string][]string) // IMEI -> slice of JSON payloads
var lastSent = make(map[string]time.Time)     // IMEI -> last sent time
var bufferMutex sync.Mutex                    // Mutex for thread-safe access
var lastGlobalSend = time.Now()               // Last time any payload was sent
var minSendInterval = 2 * time.Second         // Minimum time between API requests
var lastAPIRequest = time.Now()               // Track last API request time
var apiRateMutex sync.Mutex                   // Mutex for API rate limiting

// Add these variables at the package level
var activetrackToken string
var activetrackURL string

// Initialize function to be called once at startup
func InitActivetrac() {
	activetrackToken = os.Getenv("ACTIVETRACK_TOKEN")
	if activetrackToken == "" {
		activetrackToken = "87a0239529b1515b2a7a5a173699e3310a5404e1d8b2c8400a6139e4" // Default fallback
	}

	activetrackURL = os.Getenv("ACTIVETRACK_URL")
	if activetrackURL == "" {
		activetrackURL = "https://pegasus248.peginstances.com/receivers/json" // Default fallback
	}

	utils.VPrint("Initialized Activetrack with URL: %s and token length: %d", activetrackURL, len(activetrackToken))
}

func ProcessAndSendActiveTrack(plates, dataStr string) error {
	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		fmt.Println("Error deserializando JSON:", err)
		return fmt.Errorf("error deserializando JSON: %v", err)
	}

	imei := data.IMEI
	packetCount := 0

	// Lock the mutex once before processing all packets
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	// Process all packets in the data
	for _, packet := range data.ListPackets {
		packetCount++

		// Extract and process values from active_track_process
		eventcode := fmt.Sprintf("%d", packet.EventCode.Code)
		utils.VPrint("EventCode: %s", eventcode)
		valid := packet.PositioningStatus

		// Initialize battery values
		Battery_power := "0.0"

		// Safe handling of analog inputs
		if packet.AnalogInputs != nil {
			if packet.AnalogInputs.AD5 != nil {
				if bp, err := strconv.ParseFloat(*packet.AnalogInputs.AD5, 64); err == nil {
					Battery_power = fmt.Sprintf("%.2f", bp)
				}
			}
		}

		// Fix for timestamp handling
		Fecha_trama := packet.Datetime.Format(time.RFC3339)
		var timestamp int64

		// Properly parse the time
		if t, err := time.Parse(time.RFC3339, Fecha_trama); err == nil {
			timestamp = t.Unix()
		} else {
			// If parsing fails, use current time as fallback
			timestamp = time.Now().Unix()
			utils.VPrint("Error parsing time %s: %v, using current time instead", Fecha_trama, err)
		}

		Fecha_trama = strconv.FormatInt(timestamp, 10)

		ignition := "0"
		if eventcode == "2" {
			ignition = "1"
		} else if eventcode == "10" {
			ignition = "0"
		}
		// Fallback for empty status (treat as valid)
		if valid == "" {
			valid = "A" // Set to valid if empty
		}

		// Set boolean validity (true if "A", false otherwise)
		valid_boolean := true

		// Construct JSON payload
		currentTime := time.Now().UTC()
		strtimestamp := currentTime.Unix()
		power := ignition == "1"
		enum := 24
		label := "trckpnt"
		if eventcode == "1" {
			enum = 1
			label = "panic"
		}

		// Get values from the packet
		direction := packet.Direction
		satellites := packet.NumberOfSatellites

		hdop := packet.HDOP
		mileage := packet.Mileage
		var runtime int = 0
		if packet.RunTime > 0 {
			runtime = packet.RunTime
		}

		payload := fmt.Sprintf(
			`{"timestamp":%s,"device.id":"%s","position.latitude":%f,"position.longitude":%f,"server.timestamp":%d,"device.name":"%s","position.direction":%d,"position.speed":%d,"position.altitude":%d,"position.hdop":%f,"position.satellites":%d,"position.valid":%t,"event.enum":%d,"io.ignition":"%s","io.power":%t,"io.input1":false,"device.battery.level":%s,"metric.odometer":%d,"metric.hourmeter":%d,"protocol.id":"rt.platform","event.label":"%s","adas.snapshots.count":1,"adas.snapshot.timestamp":%d}`,
			Fecha_trama,
			imei,
			packet.Latitude,
			packet.Longitude,
			strtimestamp,
			plates,
			direction,
			packet.Speed,
			packet.Altitude,
			hdop,
			satellites,
			valid_boolean,
			enum,
			ignition,
			power,
			Battery_power,
			mileage,
			runtime,
			label,
			strtimestamp,
		)

		// For debugging
		utils.VPrint("Mapped packet %d fields: Direction=%d, Speed=%d, Battery=%.2f, Satellites=%d, HDOP=%.2f",
			packetCount, direction, packet.Speed, parseFloat(Battery_power, 0.0), satellites, hdop)

		// Add payload to buffer
		payloadBuffer[imei] = append(payloadBuffer[imei], payload)
	}

	// If no packets were processed, log a warning
	if packetCount == 0 {
		utils.VPrint("Warning: No packets found in data for IMEI %s", imei)
		return nil
	}

	utils.VPrint("Processed %d packets for IMEI %s", packetCount, imei)

	// Check conditions for sending:
	// 1. Buffer for this IMEI has at least 50 entries, OR
	// 2. It's been more than 1 minute since any IMEI's data was sent
	now := time.Now()
	size := len(payloadBuffer[imei])

	// Match original code conditions
	send := size >= 50 || now.Sub(lastGlobalSend) > time.Minute

	utils.VPrint("send flag: %t, imei: %s, size: %d, time since last send: %v",
		send, imei, size, now.Sub(lastGlobalSend))

	if send {
		utils.VPrint("------------- Sending for IMEI %s -------------", imei)

		// Prepare payload array (up to 50 entries like original code)
		count := size
		if count > 50 {
			count = 50
		}
		payloadArray := "[" + strings.Join(payloadBuffer[imei][:count], ",") + "]"

		// Send to server
		err := sendToActivetrac(imei, payloadArray)
		if err != nil {
			log.Printf("Error sending to Activetrac: %v", err)
			return err
		}

		// Remove sent entries and update timestamps
		if count < len(payloadBuffer[imei]) {
			payloadBuffer[imei] = payloadBuffer[imei][count:]
		} else {
			payloadBuffer[imei] = nil
		}

		lastSent[imei] = now
		lastGlobalSend = now

		log.Printf("Sent %d entries for IMEI %s", count, imei)
	}

	return nil
}

// Additional function to send all pending data for all IMEIs
// This mimics the sendAll function from the original code
func SendAllPendingData() {
	// Get a snapshot of the buffer to avoid holding the lock during API calls
	bufferMutex.Lock()
	pendingBuffers := make(map[string][]string)
	for imei, buffer := range payloadBuffer {
		if len(buffer) > 0 {
			pendingBuffers[imei] = make([]string, len(buffer))
			copy(pendingBuffers[imei], buffer)
		}
	}
	bufferMutex.Unlock()

	now := time.Now()
	imeiCount := len(pendingBuffers)
	successCount := 0

	utils.VPrint("========= SendAllPendingData called with %d IMEIs =========", imeiCount)

	// Process each IMEI's buffer
	for imei, buffer := range pendingBuffers {
		// Don't process too many IMEIs at once - cap at 5 to avoid overwhelming the API
		if successCount >= 5 {
			utils.VPrint("Processed 5 IMEIs, will process others in next cycle")
			break
		}

		count := len(buffer)
		if count > 50 {
			count = 50
		}

		utils.VPrint("Sending %d entries for IMEI %s", count, imei)
		payloadArray := "[" + strings.Join(buffer[:count], ",") + "]"

		// Make API call without holding lock
		err := sendToActivetrac(imei, payloadArray)

		if err != nil {
			if strings.Contains(err.Error(), "rate limit") {
				utils.VPrint("Rate limit encountered, stopping batch processing")
				break
			}
			utils.VPrint("Error sending to Activetrac: %v", err)
			continue
		}

		// Update the original buffer now that we've sent successfully
		bufferMutex.Lock()
		if originalBuffer, exists := payloadBuffer[imei]; exists {
			if count < len(originalBuffer) {
				payloadBuffer[imei] = originalBuffer[count:]
			} else {
				payloadBuffer[imei] = nil
			}
			lastSent[imei] = now
		}
		bufferMutex.Unlock()

		successCount++
	}

	// Update last global send time
	bufferMutex.Lock()
	lastGlobalSend = now
	bufferMutex.Unlock()

	utils.VPrint("Processed %d of %d IMEIs at %v", successCount, imeiCount, now)
}

func sendToActivetrac(imei, payload string) error {
	// Implement rate limiting without using the buffer mutex
	apiRateMutex.Lock()
	timeSinceLastRequest := time.Since(lastAPIRequest)
	if timeSinceLastRequest < minSendInterval {
		// Calculate wait time
		waitTime := minSendInterval - timeSinceLastRequest
		apiRateMutex.Unlock() // Release lock before sleeping
		utils.VPrint("Rate limiting: waiting %v before next request", waitTime)
		time.Sleep(waitTime)
		apiRateMutex.Lock() // Re-acquire lock after sleep
	}
	// Update last request time and release lock
	lastAPIRequest = time.Now()
	apiRateMutex.Unlock()

	// Use the cached token and URL values
	method := "POST"
	payloadReader := strings.NewReader(payload)
	utils.VPrint("payload: %s", payload)
	client := &http.Client{
		Timeout: 10 * time.Second, // Add timeout to prevent hanging requests
	}
	req, err := http.NewRequest(method, activetrackURL, payloadReader)
	if err != nil {
		log.Printf("Error en NewRequest: %v, URL: %v", err, activetrackURL)
		return err
	}

	req.Header.Add("Authenticate", activetrackToken)
	req.Header.Add("Content-Type", "application/json")

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error en client.Do: %v", err)
		return err
	}
	defer res.Body.Close()

	// Process response status code
	if res.StatusCode == 429 {
		// Special handling for rate limit exceeded
		utils.VPrint("Rate limit exceeded (429). Will retry with increased backoff.")
		apiRateMutex.Lock()
		minSendInterval = minSendInterval * 2 // Double the wait time
		if minSendInterval > 30*time.Second {
			minSendInterval = 30 * time.Second // Cap at 30 seconds
		}
		apiRateMutex.Unlock()
		return fmt.Errorf("rate limit exceeded (429)")
	} else if res.StatusCode != 200 {
		log.Printf("Error al enviar la información. StatusCode: %v", res.StatusCode)
		utils.VPrint("HTTP error: %d", res.StatusCode)
	} else {
		// Successful request, we can reduce the interval slightly
		apiRateMutex.Lock()
		if minSendInterval > 2*time.Second {
			minSendInterval = minSendInterval / 2
			if minSendInterval < 2*time.Second {
				minSendInterval = 2 * time.Second
			}
		}
		apiRateMutex.Unlock()
	}

	// Log success
	utils.VPrint("Activetrack Response: %d", res.StatusCode)

	// Handle elastic logging in a separate goroutine to avoid blocking
	elastic_doc_name := os.Getenv("ELASTIC_DOC_NAME")
	utils.VPrint("elastic_doc_name: %s", elastic_doc_name)

	go func(imeiCopy, payloadCopy string, statusCode int, statusText string, docName string) {
		logData := utils.ElasticLogData{
			Client:     docName,
			IMEI:       imeiCopy,
			Payload:    payload,
			Time:       time.Now().Format(time.RFC3339),
			StatusCode: statusCode,
			StatusText: statusText,
		}

		if err := utils.SendToElastic(logData, docName); err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
		}
	}(imei, payload, res.StatusCode, res.Status, elastic_doc_name)

	return nil
}

// Helper function to parse float with fallback
func parseFloat(s string, fallback float64) float64 {
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}
	return fallback
}
