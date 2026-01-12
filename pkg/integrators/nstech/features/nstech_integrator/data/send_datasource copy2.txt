package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/MaddSystems/jonobridge/common/models"
	"github.com/MaddSystems/jonobridge/common/utils"
)

// Struct for the NSTECH API OAuth token response
type NstechTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
}

// Position struct for NSTECH API
type NstechPosition struct {
	TechnologyID string  `json:"technology_id"`
	AccountID    string  `json:"account_id"`
	Date         string  `json:"date"`
	DeviceID     string  `json:"device_id"`
	PositionType string  `json:"position_type"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Ignition     string  `json:"ignition"`
	Speed        float64 `json:"speed"`
	Odometer     float64 `json:"odometer"`
}

// Event struct for NSTECH API
type NstechEvent struct {
	TechnologyID string  `json:"technology_id"`
	AccountID    string  `json:"account_id"`
	Date         string  `json:"date"`
	DeviceID     string  `json:"device_id"`
	EventType    string  `json:"event_type"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

var nstech_url string
var nstech_token_url string
var elastic_doc_name string
var nstech_technology_id string
var nstech_account_id string
var nstech_client_id string
var nstech_client_secret string

// OAuth token cache
var oauthToken string
var oauthTokenExpiry time.Time

// Position batching configuration
const positionBatchSize = 5

// Position batch buffer
var positionBatch []NstechPosition
var batchMutex sync.Mutex

// Event batching configuration
const eventBatchSize = 5

// Event batch buffer
var eventBatch []NstechEvent
var eventBatchMutex sync.Mutex

// mapJonoEventToNstech maps JonoProtocol event codes to NSTECH event types
func mapJonoEventToNstech(eventCode int) string {
	switch eventCode {
	case 1:
		return "PanicButton" // SOS Pressed
	case 2:
		return "BatteryConnect" // Input 2 Active (Ignition On/Door Open)
	case 10:
		return "BatteryDisconnect" // Input 2 Inactive (Ignition Off/Door Close)
	case 4:
		return "BatteryConnect" // Input 4 Active (Ignition On for MVT380)
	case 12:
		return "BatteryDisconnect" // Input 4 Inactive (Ignition Off for MVT380)
	case 17:
		return "BatteryDisconnect" // Low Battery
	case 18:
		return "BatteryDisconnect" // Low External Battery
	case 22:
		return "BatteryConnect" // External Battery On
	case 23:
		return "BatteryDisconnect" // External Battery Cut
	case 24:
		return "GPSDisconnected" // GPS Signal Lost
	case 25:
		return "GPSConnected" // GPS Signal Recovery
	case 28:
		return "GPSDisconnected" // GPS Antenna Cut
	case 29:
		return "BatteryConnect" // Device Reboot (Power On)
	case 65:
		return "PanicButton" // Press Input 1 (SOS) to Call
	case 66, 67, 68, 69:
		return "PanicButton" // Press Input 2-5 to Call (treated as panic)
	// For events without clear mapping, we'll use a generic approach
	// or skip them depending on requirements
	default:
		return "" // Empty string means skip this event
	}
}

// Initialize function to be called once at startup
func InitNstech() {
	elastic_doc_name = os.Getenv("ELASTIC_DOC_NAME")
	nstech_token_url = os.Getenv("NSTECH_TOKEN_URL")
	nstech_url = os.Getenv("NSTECH_URL")
	nstech_technology_id = os.Getenv("NSTECH_TECHNOLOGY_ID")
	nstech_account_id = os.Getenv("NSTECH_ACCOUNT_ID")
	nstech_client_id = os.Getenv("NSTECH_CLIENT_ID")
	nstech_client_secret = os.Getenv("NSTECH_CLIENT_SECRET")
}

// getOAuthToken is the only token retrieval function we now use

// getOAuthToken gets a token using the OAuth2 client credentials flow
func getOAuthToken() (string, error) {
	// Return cached token if still valid
	if oauthToken != "" && time.Now().Before(oauthTokenExpiry) {
		utils.VPrint("Using cached OAuth token (expires %v)", oauthTokenExpiry)
		return oauthToken, nil
	}

	if nstech_token_url == "" {
		return "", fmt.Errorf("NSTECH_TOKEN_URL not configured")
	}

	if nstech_client_id == "" || nstech_client_secret == "" {
		return "", fmt.Errorf("NSTECH_CLIENT_ID or NSTECH_CLIENT_SECRET not configured")
	}

	// Prepare form data
	formData := url.Values{}
	formData.Set("client_id", nstech_client_id)
	formData.Set("client_secret", nstech_client_secret)
	formData.Set("grant_type", "client_credentials")

	client := &http.Client{Timeout: 15 * time.Second}
	utils.VPrint("Requesting OAuth token from %s", nstech_token_url)

	req, err := http.NewRequest("POST", nstech_token_url, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating token request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error requesting token: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading token response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp NstechTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("error parsing token response: %v", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("received empty access token")
	}

	// Cache the token with a 10% safety margin
	expirySeconds := tokenResp.ExpiresIn
	if expirySeconds <= 0 {
		expirySeconds = 3600 // Default to 1 hour if not specified
	}

	safeExpirySeconds := int(float64(expirySeconds) * 0.9) // 90% of the actual expiry time
	oauthTokenExpiry = time.Now().Add(time.Duration(safeExpirySeconds) * time.Second)
	oauthToken = tokenResp.AccessToken

	utils.VPrint("OAuth token received, expires in %d seconds (using %d seconds)",
		expirySeconds, safeExpirySeconds)

	return oauthToken, nil
}

// sendPositionDataBatch sends multiple position data to the NSTECH API
func sendPositionDataBatch(positions []NstechPosition) error {
	if len(positions) == 0 {
		return nil
	}

	// Get OAuth token
	token, err := getOAuthToken()
	if err != nil {
		return fmt.Errorf("error obtaining OAuth token: %v", err)
	}

	if nstech_url == "" {
		return fmt.Errorf("NSTECH_URL not configured")
	}

	// Construct the correct production URL for positions, which includes the /integra path
	positionsURL := nstech_url + "/integra/v1/positions"

	// Create wrapper for positions array as required by the API
	type PositionsPayload struct {
		Positions []NstechPosition `json:"positions"`
	}
	payload := PositionsPayload{
		Positions: positions,
	}

	// Convert payload to JSON
	positionJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling position data: %v", err)
	}

	// Create HTTP request
	client := &http.Client{Timeout: 15 * time.Second}
	utils.VPrint("Sending batch of %d positions to %s", len(positions), positionsURL)

	req, err := http.NewRequest("POST", positionsURL, bytes.NewBuffer(positionJSON))
	if err != nil {
		return fmt.Errorf("error creating position request: %v", err)
	}

	// Add headers - matching the content type used in Python
	req.Header.Set("Content-Type", "application/*+json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending position data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading position response: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		utils.VPrint("Error sending position batch: StatusCode=%d, Response=%s",
			resp.StatusCode, string(body))

		// If 401, invalidate cached token to force refresh on next request
		if resp.StatusCode == 401 {
			utils.VPrint("Received 401, invalidating cached OAuth token")
			oauthToken = ""
			oauthTokenExpiry = time.Time{}
		}

		return fmt.Errorf("position endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	utils.VPrint("Position batch sent successfully: StatusCode=%d, Count=%d", resp.StatusCode, len(positions))

	// Log to Elastic for the first position in the batch (representative)
	if len(positions) > 0 {
		logData := utils.ElasticLogData{
			Client:     elastic_doc_name,
			IMEI:       positions[0].DeviceID,
			Payload:    string(positionJSON),
			Time:       time.Now().Format(time.RFC3339),
			StatusCode: resp.StatusCode,
			StatusText: resp.Status,
		}
		if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
		}
	}

	return nil
}

// addPositionToBatch adds a position to the batch buffer and sends when batch is full
func addPositionToBatch(position NstechPosition) error {
	batchMutex.Lock()
	defer batchMutex.Unlock()

	// Add position to batch
	positionBatch = append(positionBatch, position)
	utils.VPrint("Added position to batch, current batch size: %d/%d", len(positionBatch), positionBatchSize)

	// Send batch when it reaches the configured size
	if len(positionBatch) >= positionBatchSize {
		// Make a copy of the batch to send
		batchToSend := make([]NstechPosition, len(positionBatch))
		copy(batchToSend, positionBatch)

		// Clear the batch buffer
		positionBatch = positionBatch[:0]

		// Send the batch (unlock mutex before sending to avoid blocking)
		batchMutex.Unlock()
		err := sendPositionDataBatch(batchToSend)
		batchMutex.Lock()

		if err != nil {
			// Log error but don't crash - continue processing
			utils.VPrint("Error sending position batch (will retry later): %v", err)
			// Put the data back in the batch for retry
			positionBatch = append(batchToSend, positionBatch...)
		}
	}

	return nil
}

// flushPositionBatch sends any remaining positions in the batch buffer
// This function can be called during shutdown or when you want to force send pending positions
func flushPositionBatch() error {
	batchMutex.Lock()
	defer batchMutex.Unlock()

	if len(positionBatch) > 0 {
		// Make a copy of the batch to send
		batchToSend := make([]NstechPosition, len(positionBatch))
		copy(batchToSend, positionBatch)

		// Clear the batch buffer
		positionBatch = positionBatch[:0]

		// Send the batch - if it fails, log but don't crash
		if err := sendPositionDataBatch(batchToSend); err != nil {
			utils.VPrint("Error flushing position batch (data may be lost): %v", err)
			// Don't return error to avoid crashing
		}
	}

	return nil
}

// sendEventDataBatch sends multiple event data to the NSTECH API
func sendEventDataBatch(events []NstechEvent) error {
	if len(events) == 0 {
		return nil
	}

	// Get OAuth token
	token, err := getOAuthToken()
	if err != nil {
		return fmt.Errorf("error obtaining OAuth token: %v", err)
	}

	if nstech_url == "" {
		return fmt.Errorf("NSTECH_URL not configured")
	}

	// Construct the correct production URL for events, which includes the /integra path
	eventsURL := nstech_url + "/integra/v1/events"

	// Create wrapper for events array as required by the API
	type EventsPayload struct {
		Events []NstechEvent `json:"events"`
	}
	payload := EventsPayload{
		Events: events,
	}

	// Convert payload to JSON
	eventJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling event data: %v", err)
	}

	// Create HTTP request
	client := &http.Client{Timeout: 15 * time.Second}
	utils.VPrint("Sending batch of %d events to %s", len(events), eventsURL)

	req, err := http.NewRequest("POST", eventsURL, bytes.NewBuffer(eventJSON))
	if err != nil {
		return fmt.Errorf("error creating event request: %v", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending event data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading event response: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		utils.VPrint("Error sending event batch: StatusCode=%d, Response=%s",
			resp.StatusCode, string(body))

		// If 401, invalidate cached token to force refresh on next request
		if resp.StatusCode == 401 {
			utils.VPrint("Received 401, invalidating cached OAuth token")
			oauthToken = ""
			oauthTokenExpiry = time.Time{}
		}

		return fmt.Errorf("event endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	utils.VPrint("Event batch sent successfully: StatusCode=%d, Count=%d", resp.StatusCode, len(events))

	// Log to Elastic for the first event in the batch (representative)
	if len(events) > 0 {
		logData := utils.ElasticLogData{
			Client:     elastic_doc_name,
			IMEI:       events[0].DeviceID,
			Payload:    string(eventJSON),
			Time:       time.Now().Format(time.RFC3339),
			StatusCode: resp.StatusCode,
			StatusText: resp.Status,
		}
		if err := utils.SendToElastic(logData, elastic_doc_name); err != nil {
			utils.VPrint("Error sending to Elasticsearch: %v", err)
		}
	}

	return nil
}

// addEventToBatch adds an event to the batch buffer and sends when batch is full
func addEventToBatch(event NstechEvent) error {
	eventBatchMutex.Lock()
	defer eventBatchMutex.Unlock()

	// Add event to batch
	eventBatch = append(eventBatch, event)
	utils.VPrint("Added event to batch, current batch size: %d/%d", len(eventBatch), eventBatchSize)

	// Send batch when it reaches the configured size
	if len(eventBatch) >= eventBatchSize {
		// Make a copy of the batch to send
		batchToSend := make([]NstechEvent, len(eventBatch))
		copy(batchToSend, eventBatch)

		// Clear the batch buffer
		eventBatch = eventBatch[:0]

		// Send the batch (unlock mutex before sending to avoid blocking)
		eventBatchMutex.Unlock()
		err := sendEventDataBatch(batchToSend)
		eventBatchMutex.Lock()

		if err != nil {
			// Log error but don't crash - continue processing
			utils.VPrint("Error sending event batch (will retry later): %v", err)
			// Put the data back in the batch for retry
			eventBatch = append(batchToSend, eventBatch...)
		}
	}

	return nil
}

// flushEventBatch sends any remaining events in the batch buffer
func flushEventBatch() error {
	eventBatchMutex.Lock()
	defer eventBatchMutex.Unlock()

	if len(eventBatch) > 0 {
		// Make a copy of the batch to send
		batchToSend := make([]NstechEvent, len(eventBatch))
		copy(batchToSend, eventBatch)

		// Clear the batch buffer
		eventBatch = eventBatch[:0]

		// Send the batch - if it fails, log but don't crash
		if err := sendEventDataBatch(batchToSend); err != nil {
			utils.VPrint("Error flushing event batch (data may be lost): %v", err)
			// Don't return error to avoid crashing
		}
	}

	return nil
}

func ProcessAndSendNstech(plates, eco, vin, dataStr string) error {
	// Add defensive error handling to prevent container crashes
	defer func() {
		if r := recover(); r != nil {
			utils.VPrint("Recovered from panic in ProcessAndSendNstech: %v", r)
		}
	}()

	// Parse the incoming JSON data
	var data models.JonoModel
	err := json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		utils.VPrint("Error deserializando JSON: %v", err)
		// Don't return error to avoid crashing - just log and continue
		return nil
	}

	// Process all packets in the data
	for _, packet := range data.ListPackets {
		utils.VPrint("IMEI: %s", data.IMEI)
		utils.VPrint("Coordinates: %f,%f", packet.Latitude, packet.Longitude)
		imei := data.IMEI
		event_code := packet.EventCode.Code

		// Check if this is a position event (code 35)
		if event_code == 35 && nstech_url != "" && nstech_technology_id != "" && nstech_account_id != "" {
			utils.VPrint("Detected position event (code 35), sending to NSTECH positions API")

			// Create position object with the required fields
			position := NstechPosition{
				TechnologyID: nstech_technology_id,
				AccountID:    nstech_account_id,
				Date:         packet.Datetime.Format("2006-01-02T15:04:05.000Z"),
				DeviceID:     imei,
				PositionType: "GPRS",
				Latitude:     packet.Latitude,
				Longitude:    packet.Longitude,
				Ignition:     "On", // As specified in requirements, ignition is always "On" in jonoprotocol
				Speed:        float64(packet.Speed),
				Odometer:     float64(packet.Mileage),
			}

			utils.VPrint("Sending position: Lat=%f, Lon=%f, Speed=%d, Ignition=%s, DeviceID=%s",
				packet.Latitude, packet.Longitude, packet.Speed, "On", imei)

			// Add position to batch (will send automatically when batch is full)
			if err := addPositionToBatch(position); err != nil {
				utils.VPrint("Error adding position to batch: %v", err)
			}
		} else {
			// Check if this event can be mapped to NSTECH event type
			nstechEventType := mapJonoEventToNstech(event_code)

			if nstechEventType != "" && nstech_url != "" && nstech_technology_id != "" && nstech_account_id != "" {
				utils.VPrint("Detected mappable event (code %d -> %s), sending to NSTECH events API", event_code, nstechEventType)

				// Create event object with the required fields
				event := NstechEvent{
					TechnologyID: nstech_technology_id,
					AccountID:    nstech_account_id,
					Date:         packet.Datetime.Format("2006-01-02T15:04:05.000Z"),
					DeviceID:     imei,
					EventType:    nstechEventType,
					Latitude:     packet.Latitude,
					Longitude:    packet.Longitude,
				}

				utils.VPrint("Sending event: Type=%s, Lat=%f, Lon=%f, DeviceID=%s",
					nstechEventType, packet.Latitude, packet.Longitude, imei)

				// Add event to batch (will send automatically when batch is full)
				if err := addEventToBatch(event); err != nil {
					utils.VPrint("Error adding event to batch: %v", err)
				}
			} else {
				// For events that don't map to NSTECH event types, discard them
				if nstechEventType == "" {
					utils.VPrint("Event code %d not mapped to NSTECH event type, discarding event", event_code)
				} else {
					utils.VPrint("Event code %d mapped but NSTECH configuration missing, discarding event", event_code)
				}
				// Event is discarded - no further processing
			}
		}
	}

	// Flush any remaining positions and events in the batches after processing all packets
	// Don't let flush errors crash the container
	if err := flushPositionBatch(); err != nil {
		utils.VPrint("Error flushing remaining positions (non-fatal): %v", err)
	}

	if err := flushEventBatch(); err != nil {
		utils.VPrint("Error flushing remaining events (non-fatal): %v", err)
	}

	return nil
}
