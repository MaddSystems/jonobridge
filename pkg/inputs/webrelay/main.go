package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/utils"
)

var portal_endpoint string
var webrelay_token string // Token for Bearer authentication

// Constants for server1.gpscontrol.com.mx
const (
	DefaultAppID       = 424
	DefaultGGSUser     = "admindesarrollo"
	DefaultGGSPassword = "GPSc0ntr0l00"
)

// Configuration struct for server1 access
type Server1Config struct {
	AppID    int
	Username string
	Password string
}

// Vehicle data response struct
type VehicleData struct {
	IMEI           string  `json:"imei"`
	Plate          string  `json:"plate"`
	Altitude       float64 `json:"altitude"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Speed          float64 `json:"speed"`
	Heading        float64 `json:"heading"`
	Date           string  `json:"date"`
	Time           string  `json:"time"`
	Moving         bool    `json:"moving"`
	IgnitionStatus bool    `json:"ignitionStatus"`
	StoppingDate   string  `json:"stoppingDate"`
	StopingTime    string  `json:"stopingTime"`
	Error          string  `json:"error,omitempty"`
}

// Server1 API response types
type TokenResponse struct {
	Token string `json:"token"`
}

type Position struct {
	Altitude  float64 `json:"altitude"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Velocity struct {
	GroundSpeed float64 `json:"groundSpeed"`
	Heading     float64 `json:"heading"`
}

type TrackPoint struct {
	UTC      string   `json:"utc"`
	Position Position `json:"position"`
	Velocity Velocity `json:"velocity,omitempty"`
}

type TripInfo struct {
	StartTrackPoint TrackPoint `json:"startTrackPoint"`
	EndTrackPoint   TrackPoint `json:"endTrackPoint"`
	TotalDistance   float64    `json:"totalDistance"`
}

type Variable struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type Status struct {
	Variables []Variable `json:"variables"`
}

type Device struct {
	IMEI string `json:"imei"`
}

type User struct {
	ID             json.Number `json:"id"` // Change from string to json.Number
	Name           string      `json:"name"`
	Devices        []Device    `json:"devices"`
	TrackPoint     TrackPoint  `json:"trackPoint"`
	DeviceActivity string      `json:"deviceActivity"`
}

// authMiddleware implements Bearer token authentication
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If token is not set, skip auth
		if webrelay_token == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization header required"))
			return
		}

		// Check that it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authorization header must start with Bearer"))
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		if token != webrelay_token {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token"))
			return
		}

		// Token is valid, proceed
		next.ServeHTTP(w, r)
	}
}

// getServer1Token gets an authentication token from server1
func getServer1Token(config Server1Config) (string, error) {
	utils.VPrint("Getting server1 token for app_id: %d, user: %s", config.AppID, config.Username)
	url := fmt.Sprintf("http://server1.gpscontrol.com.mx/comGPSGate/api/v.1/applications/%d/tokens", config.AppID)

	payload := map[string]string{
		"username": config.Username,
		"password": config.Password,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error creating JSON payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		utils.VPrint("Error making token request: %v", err)
		return "", fmt.Errorf("error making request to server1: %v", err)
	}
	defer resp.Body.Close()

	utils.VPrint("Token request status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.VPrint("Error response body: %s", string(body))
		return "", fmt.Errorf("server1 returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.VPrint("Error reading token response body: %v", err)
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Fix the slice bounds error by checking length
	bodyStr := string(body)
	displayLen := len(bodyStr)
	if displayLen > 100 {
		displayLen = 100
	}
	utils.VPrint("Token response body: %s", bodyStr[:displayLen]) // Print first 100 chars for security

	// Check for error message in response
	if strings.Contains(string(body), "The user does not have neither _APIRead nor _APIReadWrite privileges") {
		return "", fmt.Errorf("faltan permisos en la cuenta. _APIRead y _APIReadWrite")
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("error parsing token response: %v", err)
	}

	utils.VPrint("Successfully parsed token")
	return tokenResp.Token, nil
}

// extractDateTime formats date and time from a datetime string
func extractDateTime(datetimeStr string) (date, timeStr string, dt time.Time, err error) {
	// Remove 'Z' if present
	datetimeStr = strings.TrimSuffix(datetimeStr, "Z")

	dt, err = time.Parse("2006-01-02T15:04:05", datetimeStr)
	if err != nil {
		return "", "", time.Time{}, err
	}

	date = dt.Format("01-02-2006")
	timeStr = dt.Format("15:04:05")

	return date, timeStr, dt, nil
}

// findMostRecent finds the most recent trip info
func findMostRecent(tripInfos []TripInfo) TripInfo {
	var mostRecent TripInfo
	var mostRecentTime time.Time

	for i, trip := range tripInfos {
		endTime, err := time.Parse("2006-01-02T15:04:05Z", trip.EndTrackPoint.UTC)
		if err != nil {
			continue
		}

		if i == 0 || endTime.After(mostRecentTime) {
			mostRecent = trip
			mostRecentTime = endTime
		}
	}

	return mostRecent
}

// getVehicleData gets vehicle data for a specific plate
func getVehicleData(plates string, config Server1Config) (VehicleData, error) {
	utils.VPrint("Getting vehicle data for plates: %s", plates)
	var data VehicleData

	// Get token
	token, err := getServer1Token(config)
	if err != nil {
		utils.VPrint("Failed to get token: %v", err)
		return data, fmt.Errorf("error getting token: %v", err)
	}

	utils.VPrint("Token obtained for app_id: %d (token starts with: %s...)", config.AppID, token[:10])

	// Get users
	usersURL := fmt.Sprintf("https://server1.gpscontrol.com.mx/comGpsGate/api/v.1/applications/%d/users?FromIndex=0&PageSize=1000", config.AppID)
	utils.VPrint("Requesting users from URL: %s", usersURL)

	req, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		utils.VPrint("Error creating request: %v", err)
		return data, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", token)
	utils.VPrint("Added Authorization header")

	client := &http.Client{Timeout: 30 * time.Second}
	utils.VPrint("Sending request to get users...")
	resp, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error fetching users: %v", err)
		return data, fmt.Errorf("error fetching users: %v", err)
	}
	defer resp.Body.Close()

	utils.VPrint("Users request status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		utils.VPrint("Error in users response: %s", string(bodyBytes))
		return data, fmt.Errorf("error getting users: status code %d", resp.StatusCode)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		utils.VPrint("Error reading users response: %v", readErr)
		return data, fmt.Errorf("error reading users response: %v", readErr)
	}

	// Fix the slice bounds error here too
	bodyStr := string(body)
	displayLen := min(200, len(bodyStr))
	utils.VPrint("Received users JSON (first %d chars): %s", displayLen, bodyStr[:displayLen])
	utils.VPrint("Parsing users JSON...")

	// Debug - check if we have array
	trimmed := strings.TrimSpace(string(body))
	if !strings.HasPrefix(trimmed, "[") {
		utils.VPrint("WARNING: Users response doesn't start with array bracket '['")
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		utils.VPrint("ERROR parsing users JSON: %v", err)
		utils.VPrint("Attempting to debug JSON structure...")

		// Try parsing as generic interface to see the structure
		var rawData interface{}
		if jsonErr := json.Unmarshal(body, &rawData); jsonErr == nil {
			utils.VPrint("JSON parsed as interface{}, type: %T", rawData)

			// If it's an array, check first element
			if arr, ok := rawData.([]interface{}); ok && len(arr) > 0 {
				first := arr[0]
				utils.VPrint("First element type: %T", first)

				// If first element is a map, check for 'id' field
				if obj, ok := first.(map[string]interface{}); ok {
					if id, exists := obj["id"]; exists {
						utils.VPrint("'id' field type: %T, value: %v", id, id)
					} else {
						utils.VPrint("'id' field not found in first element")
					}
				}
			}
		} else {
			utils.VPrint("Failed even generic JSON parsing: %v", jsonErr)
		}

		return data, fmt.Errorf("error parsing users: %v", err)
	}

	utils.VPrint("Found %d users", len(users))

	// Find user with matching plates
	var user *User
	for i := range users {
		if users[i].Name == plates {
			user = &users[i]
			break
		}
	}

	if user == nil {
		return data, fmt.Errorf("las placas solicitadas no existen")
	}

	utils.VPrint("Found user with plates: %s", user.Name)

	// Process devices
	if len(user.Devices) == 0 {
		return data, fmt.Errorf("no devices found for the requested plates")
	}

	// Use first device
	imei := user.Devices[0].IMEI
	altitude := user.TrackPoint.Position.Altitude
	latitude := user.TrackPoint.Position.Latitude
	longitude := user.TrackPoint.Position.Longitude
	speed := user.TrackPoint.Velocity.GroundSpeed * 3.6 // Convert to km/h
	heading := user.TrackPoint.Velocity.Heading

	// Extract date/time
	date, timeStr, dt, err := extractDateTime(user.DeviceActivity)
	if err != nil {
		return data, fmt.Errorf("error parsing device time: %v", err)
	}

	// Initialize default values
	ignitionStatus := false
	stoppingTime := ""
	stoppingDate := ""

	// Get ignition status
	statusURL := fmt.Sprintf("https://server1.gpscontrol.com.mx/comGpsGate/api/v.1/applications/%d/users/%s/status",
		config.AppID, user.ID.String()) // Use .String() method to convert json.Number to string
	req, _ = http.NewRequest("GET", statusURL, nil)
	req.Header.Add("Authorization", token)

	resp, err = client.Do(req)
	if err != nil {
		return data, fmt.Errorf("error fetching status: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		var status Status
		body, _ = ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &status); err == nil {
			for _, variable := range status.Variables {
				if variable.Name == "Ignition" {
					switch v := variable.Value.(type) {
					case bool:
						ignitionStatus = v
					case float64:
						ignitionStatus = v != 0
					case string:
						ignitionStatus = v == "true" || v == "1"
					}
					break
				}
			}
		}
	}
	resp.Body.Close()

	// Get trip info
	tripDate := dt.Format("2006-01-02")
	tripURL := fmt.Sprintf("https://server1.gpscontrol.com.mx/comGpsGate/api/v.1/applications/%d/users/%s/tripinfos?Date=%s",
		config.AppID, user.ID.String(), tripDate)

	req, _ = http.NewRequest("GET", tripURL, nil)
	req.Header.Add("Authorization", token)

	resp, err = client.Do(req)
	moving := true
	stoppingTime = ""  // Use assignment, not declaration
	stoppingDate = ""  // Use assignment, not declaration

	if err == nil && resp.StatusCode == http.StatusOK {
		var tripInfos []TripInfo
		body, _ = ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &tripInfos); err == nil {
			if len(tripInfos) > 0 {
				// Trip infos found, process them
				tripInfo := findMostRecent(tripInfos)

				if tripInfo.TotalDistance == 0 {
					moving = false
					stoppingDate = date

					// Calculate stopping time
					startTime, sErr := time.Parse("2006-01-02T15:04:05Z", tripInfo.StartTrackPoint.UTC)
					endTime, eErr := time.Parse("2006-01-02T15:04:05Z", tripInfo.EndTrackPoint.UTC)

					if sErr == nil && eErr == nil {
						totalSeconds := int(endTime.Sub(startTime).Seconds())
						hours := totalSeconds / 3600
						minutes := (totalSeconds % 3600) / 60
						stoppingTime = fmt.Sprintf("%02d:%02d", hours, minutes)
					}
				}
			} else {
				// No trips found, vehicle is not moving
				moving = false
				stoppingDate = date
				stoppingTime = "00:00"
			}
		} else {
			// Error parsing trip infos, assume not moving
			moving = false
			stoppingDate = date
			stoppingTime = "00:00"
		}
	} else {
		// Error getting trip info, assume not moving
		moving = false
		stoppingDate = date
		stoppingTime = "00:00"
	}
	resp.Body.Close()

	// Populate result
	data = VehicleData{
		IMEI:           imei,
		Plate:          plates,
		Altitude:       altitude,
		Latitude:       latitude,
		Longitude:      longitude,
		Speed:          speed,
		Heading:        heading,
		Date:           date,
		Time:           timeStr,
		Moving:         moving,
		IgnitionStatus: ignitionStatus,
		StoppingDate:   stoppingDate,
		StopingTime:    stoppingTime,
	}

	return data, nil
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// vehicleDataHandler implements the handler that gets vehicle data
func vehicleDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.VPrint("Received request to vehicle data endpoint")

		// Get plates from the request parameters - allow empty for "all vehicles"
		plates := r.URL.Query().Get("plates")
		utils.VPrint("Fetching data for plates: '%s'", plates)

		// Get app_id from environment or use default
		appIDStr := os.Getenv("APP_ID")
		appID := DefaultAppID
		if appIDStr != "" {
			if id, err := strconv.Atoi(appIDStr); err == nil {
				appID = id
			}
		}

		// Get username and password from environment or use defaults
		username := os.Getenv("GGS_USER")
		if username == "" {
			username = DefaultGGSUser
		}

		password := os.Getenv("GGS_PASSWORD")
		if password == "" {
			password = DefaultGGSPassword
		}

		// Configure server1 access
		config := Server1Config{
			AppID:    appID,
			Username: username,
			Password: password,
		}

		// Check if we should return all vehicles or specific one
		if plates == "" {
			// Return all vehicles
			utils.VPrint("Calling getAllVehicleData")
			dataArray, err := getAllVehicleData(config)
			if err != nil {
				utils.VPrint("Error getting all vehicle data: %v", err)
				errorResponse := VehicleData{Error: err.Error()}
				jsonResponse, _ := json.Marshal(errorResponse)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResponse)
				return
			}

			utils.VPrint("Successfully got data for %d vehicles", len(dataArray))
			w.Header().Set("Content-Type", "application/json")
			jsonResponse, err := json.Marshal(dataArray)
			if err != nil {
				utils.VPrint("Error marshaling JSON: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			w.Write(jsonResponse)
			utils.VPrint("Successfully returned data for all vehicles")
		} else {
			// Return specific vehicle
			utils.VPrint("Calling getVehicleData with plates: %s", plates)
			data, err := getVehicleData(plates, config)
			if err != nil {
				utils.VPrint("Error getting vehicle data: %v", err)
				data = VehicleData{
					Error: err.Error(),
					Plate: plates,
				}
			} else {
				utils.VPrint("Successfully got data for plates %s: %+v", plates, data)
			}

			// Return JSON response
			w.Header().Set("Content-Type", "application/json")
			jsonResponse, err := json.Marshal(data)
			if err != nil {
				utils.VPrint("Error marshaling JSON: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Write(jsonResponse)
			utils.VPrint("Successfully returned vehicle data for plates: %s", plates)
		}
	}
}

// getAllVehicleData gets data for all vehicles
func getAllVehicleData(config Server1Config) ([]VehicleData, error) {
	utils.VPrint("Getting data for all vehicles")
	var allData []VehicleData

	// Get token
	token, err := getServer1Token(config)
	if err != nil {
		utils.VPrint("Failed to get token: %v", err)
		return allData, fmt.Errorf("error getting token: %v", err)
	}

	utils.VPrint("Token obtained for app_id: %d", config.AppID)

	// Get users
	usersURL := fmt.Sprintf("https://server1.gpscontrol.com.mx/comGpsGate/api/v.1/applications/%d/users?FromIndex=0&PageSize=1000", config.AppID)
	utils.VPrint("Requesting users from URL: %s", usersURL)

	req, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		utils.VPrint("Error creating request: %v", err)
		return allData, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", token)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		utils.VPrint("Error fetching users: %v", err)
		return allData, fmt.Errorf("error fetching users: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		utils.VPrint("Error in users response: %s", string(bodyBytes))
		return allData, fmt.Errorf("error getting users: status code %d", resp.StatusCode)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		utils.VPrint("Error reading users response: %v", readErr)
		return allData, fmt.Errorf("error reading users response: %v", readErr)
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		utils.VPrint("ERROR parsing users JSON: %v", err)
		return allData, fmt.Errorf("error parsing users: %v", err)
	}

	utils.VPrint("Found %d users, processing all vehicles", len(users))

	// Process all users with devices
	for _, user := range users {
		if len(user.Devices) == 0 {
			continue // Skip users without devices
		}

		for _, device := range user.Devices {
			if device.IMEI == "" {
				continue // Skip devices without IMEI
			}

			utils.VPrint("Processing device IMEI: %s for user: %s", device.IMEI, user.Name)

			// Extract basic vehicle data
			imei := device.IMEI
			plates := user.Name
			altitude := user.TrackPoint.Position.Altitude
			latitude := user.TrackPoint.Position.Latitude
			longitude := user.TrackPoint.Position.Longitude
			speed := user.TrackPoint.Velocity.GroundSpeed * 3.6
			heading := user.TrackPoint.Velocity.Heading

			// Extract date/time
			date, timeStr, dt, err := extractDateTime(user.DeviceActivity)
			if err != nil {
				utils.VPrint("Error parsing device time for %s: %v", user.Name, err)
				continue // Skip this device
			}

			// Get ignition status
			ignitionStatus := false
			statusURL := fmt.Sprintf("https://server1.gpscontrol.com.mx/comGpsGate/api/v.1/applications/%d/users/%s/status",
				config.AppID, user.ID.String())
			req, _ := http.NewRequest("GET", statusURL, nil)
			req.Header.Add("Authorization", token)

			statusResp, err := client.Do(req)
			if err == nil && statusResp.StatusCode == http.StatusOK {
				var status Status
				statusBody, _ := ioutil.ReadAll(statusResp.Body)
				if err := json.Unmarshal(statusBody, &status); err == nil {
					for _, variable := range status.Variables {
						if variable.Name == "Ignition" {
							switch v := variable.Value.(type) {
							case bool:
								ignitionStatus = v
							case float64:
								ignitionStatus = v != 0
							case string:
								ignitionStatus = v == "true" || v == "1"
							}
							break
						}
					}
				}
			}
			if statusResp != nil {
				statusResp.Body.Close()
			}

			// Get trip info
			moving := true
			stoppingTime := ""
			stoppingDate := ""
			tripDate := dt.Format("2006-01-02")
			tripURL := fmt.Sprintf("https://server1.gpscontrol.com.mx/comGpsGate/api/v.1/applications/%d/users/%s/tripinfos?Date=%s",
				config.AppID, user.ID.String(), tripDate)

			req, _ = http.NewRequest("GET", tripURL, nil)
			req.Header.Add("Authorization", token)

			tripResp, err := client.Do(req)
			if err == nil && tripResp.StatusCode == http.StatusOK {
				var tripInfos []TripInfo
				tripBody, _ := ioutil.ReadAll(tripResp.Body)

				if err := json.Unmarshal(tripBody, &tripInfos); err == nil {
					if len(tripInfos) > 0 {
						// Trip infos found, process them
						tripInfo := findMostRecent(tripInfos)

						if tripInfo.TotalDistance == 0 {
							moving = false
							stoppingDate = date

							// Calculate stopping time
							startTime, sErr := time.Parse("2006-01-02T15:04:05Z", tripInfo.StartTrackPoint.UTC)
							endTime, eErr := time.Parse("2006-01-02T15:04:05Z", tripInfo.EndTrackPoint.UTC)

							if sErr == nil && eErr == nil {
								totalSeconds := int(endTime.Sub(startTime).Seconds())
								hours := totalSeconds / 3600
								minutes := (totalSeconds % 3600) / 60
								stoppingTime = fmt.Sprintf("%02d:%02d", hours, minutes)
							}
						}
					} else {
						// No trips found, vehicle is not moving
						moving = false
						stoppingDate = date
						stoppingTime = "00:00"
					}
				} else {
					// Error parsing trip infos, assume not moving
					moving = false
					stoppingDate = date
					stoppingTime = "00:00"
				}
			} else {
				// Error getting trip info, assume not moving
				moving = false
				stoppingDate = date
				stoppingTime = "00:00"
			}
			if tripResp != nil {
				tripResp.Body.Close()
			}

			// Create vehicle data entry
			vehicleData := VehicleData{
				IMEI:           imei,
				Plate:          plates,
				Altitude:       altitude,
				Latitude:       latitude,
				Longitude:      longitude,
				Speed:          speed,
				Heading:        heading,
				Date:           date,
				Time:           timeStr,
				Moving:         moving,
				IgnitionStatus: ignitionStatus,
				StoppingDate:   stoppingDate,
				StopingTime:    stoppingTime,
			}

			allData = append(allData, vehicleData)
		}
	}

	utils.VPrint("Processed %d vehicles total", len(allData))
	return allData, nil
}

func main() {
	// Parse command-line flags
	flag.Parse()
	portal_endpoint = os.Getenv("PORTAL_ENDPOINT")
	if portal_endpoint == "" {
		portal_endpoint = "/test"
	} else if !strings.HasPrefix(portal_endpoint, "/") {
		// Ensure portal_endpoint always starts with a slash
		portal_endpoint = "/" + portal_endpoint
	}

	// Read environment variables
	webrelay_token = os.Getenv("WEBRELAY_TOKEN")

	// Log authentication status
	if webrelay_token != "" {
		utils.VPrint("Bearer token authentication enabled")
	} else {
		utils.VPrint("Bearer token authentication disabled - WEBRELAY_TOKEN not set")
	}

	// Start HTTP server
	http.HandleFunc(portal_endpoint, authMiddleware(vehicleDataHandler()))
	utils.VPrint("Try accessing at: https://jonobridge.madd.com.mx%s?plates=105", portal_endpoint)
	utils.VPrint("For all vehicles: https://jonobridge.madd.com.mx%s", portal_endpoint)
	utils.VPrint("With authorization: curl -H \"Authorization: Bearer d655eea7616e05b35dc7b22dd83b6ebc\" \"https://jonobridge.madd.com.mx%s?plates=105\"", portal_endpoint)

	// Add a debug endpoint
	http.HandleFunc("/debug/endpoints", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Registered endpoints:\n")
		fmt.Fprintf(w, "- %s?plates=XXX (Vehicle data endpoint)\n", portal_endpoint)
		fmt.Fprintf(w, "- /debug/endpoints (this debugging endpoint)\n")
	})

	// Run HTTP server
	utils.VPrint("Starting HTTP server on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
