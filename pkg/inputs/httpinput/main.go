package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/MaddSystems/jonobridge/common/utils"
)

var (
	portal_endpoint string
	mqttClient      mqtt.Client
)

// portalHandler reads and prints all query parameters from the request.
func portalHandler(w http.ResponseWriter, r *http.Request) {
	utils.VPrint("Received request to portal endpoint")

	// Get all query parameters
	queryParams := r.URL.Query()

	if len(queryParams) == 0 {
		utils.VPrint("No query parameters received.")
		fmt.Fprintln(w, "No query parameters received.")
		return
	}

	// Create a map to hold the query parameters for JSON conversion.
	paramsMap := make(map[string]string)
	utils.VPrint("Received query parameters:")
	for key, values := range queryParams {
		// Use the first value for each key.
		value := values[0]

		// Unescape the value to handle characters like %u00e1
		unescapedValue, err := url.QueryUnescape(value)
		if err != nil {
			// If there's an error, use the original value but log the error
			utils.VPrint("Error unescaping value for key '%s': %v", key, err)
			unescapedValue = value
		} else {
			// Handle non-standard %uXXXX escapes
			unescapedValue = unescapePercentU(unescapedValue)
		}

		paramsMap[key] = unescapedValue
		utils.VPrint("  %s: %s", key, unescapedValue)
	}

	// Convert the map to a JSON object.
	jsonBytes, err := json.Marshal(paramsMap)
	if err != nil {
		utils.VPrint("Error marshalling JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert the JSON object to a hex string.
	hexString := hex.EncodeToString(jsonBytes)
	utils.VPrint("Hex encoded JSON: %s", hexString)

	// Publish the hex string to MQTT if the client is available.
	if mqttClient != nil && mqttClient.IsConnected() {
		topic := "gpsgate"
		token := mqttClient.Publish(topic, 0, false, hexString)
		token.Wait()
		if token.Error() != nil {
			utils.VPrint("Error publishing to MQTT topic %s: %v", topic, token.Error())
		} else {
			utils.VPrint("Successfully published message to MQTT topic: %s", topic)
		}
	} else {
		utils.VPrint("MQTT client not connected, skipping publish.")
	}

	// Respond to the client
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<HttpExpression><Result>1</Result></HttpExpression>")
}

// unescapePercentU handles non-standard %uXXXX URL escapes.
func unescapePercentU(s string) string {
	for {
		idx := strings.Index(s, "%u")
		if idx == -1 || len(s) < idx+6 {
			break
		}

		// Get the 4 hex digits.
		hexCode := s[idx+2 : idx+6]

		// Convert hex to an integer.
		code, err := strconv.ParseInt(hexCode, 16, 32)
		if err != nil {
			// If conversion fails, just move on.
			s = s[idx+1:]
			continue
		}

		// Replace the %uXXXX sequence with the actual character.
		s = s[:idx] + string(rune(code)) + s[idx+6:]
	}
	return s
}

func main() {
	// Parse command-line flags to process the "-v" flag for verbose logging.
	flag.Parse()

	// Setup MQTT client
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Println("MQTT_BROKER_HOST environment variable not set. MQTT will be disabled.")
	} else {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("tcp://%s:1883", mqttBrokerHost))
		opts.SetClientID("httpinput-listener")
		mqttClient = mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("Error connecting to MQTT broker: %v", token.Error())
		} else {
			log.Println("Successfully connected to MQTT broker.")
		}
	}

	// Get the portal endpoint from environment variable or use default
	portal_endpoint = os.Getenv("PORTAL_ENDPOINT")
	if portal_endpoint == "" {
		portal_endpoint = "/test"
	}

	// Ensure the endpoint path starts with a "/"
	if !strings.HasPrefix(portal_endpoint, "/") {
		portal_endpoint = "/" + portal_endpoint
	}

	// Register the handler for the endpoint
	http.HandleFunc(portal_endpoint, portalHandler)

	// Provide instructions on how to use the endpoint
	utils.VPrint("Starting HTTP server on port 8081")
	utils.VPrint("Listening on endpoint: %s", portal_endpoint)
	utils.VPrint("Example usage:")
	utils.VPrint(`curl -v "https://jonobridge.madd.com.mx%s?cmd=_ExternalNotification&SIGNAL_SOS=true&RULE_NAME=SOS&EVENT_TIME=2010-10-03T15:02:56&POS_LATITUDE=59.28687&POS_LONGITUDE=18.08927"`, portal_endpoint)

	// Start the HTTP server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
