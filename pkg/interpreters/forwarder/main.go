package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/sync/semaphore"
)

var (
	verbose = flag.Bool("v", false, "Enable verbose logging") // Verbose flag
	// Limit concurrent connections to avoid overwhelming the target server
	connectionSemaphore = semaphore.NewWeighted(10) // Allow max 10 concurrent connections
	messageCounter      int64                       // Counter for processed messages
)

// Helper function to print verbose logs if enabled
func vPrint(format string, v ...interface{}) {
	if *verbose {
		log.Printf(format, v...)
	}
}

// TrackerData represents the JSON structure to send via HTTP
type TrackerData struct {
	TrackerData string `json:"trackerdata"`
}

// TrackerPayload represents the JSON structure from the tracker/from-tcp topic
type TrackerPayload struct {
	Payload string `json:"payload"`
}

// Function to post data to the HTTP server and publish to MQTT
func postDataToServer(message_send string, serverAddress string) {
	// Acquire semaphore to limit concurrent connections
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := connectionSemaphore.Acquire(ctx, 1); err != nil {
		vPrint("Failed to acquire connection semaphore: %v", err)
		return
	}
	defer connectionSemaphore.Release(1)

	// Attempt to establish a TCP connection to the server
	conn, err := net.DialTimeout("tcp", serverAddress, 3*time.Second)
	if err != nil {
		vPrint("Failed to connect to %s: %v", serverAddress, err)
		return // Exit the function if connection fails
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			vPrint("Error closing connection: %v", err)
		} else {
			vPrint("Connection closed.")
		}
	}()
	vPrint("Successfully connected to %s", serverAddress)

	// Set a shorter deadline for the write operation
	deadline := time.Now().Add(5 * time.Second)
	err = conn.SetDeadline(deadline)
	if err != nil {
		vPrint("Failed to set deadline: %v", err)
		return
	}

	// Send the message
	bytesWritten, err := fmt.Fprintf(conn, "%s\n", message_send)
	if err != nil {
		vPrint("Failed to send message: %v", err)
		return
	}
	vPrint("Sent %d bytes to %s", bytesWritten, serverAddress)

	// Set a shorter timeout for reading response
	err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		vPrint("Failed to set read deadline: %v", err)
		return
	}

	// Optionally, read a response from the server
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		vPrint("No response received or failed to read response: %v", err)
	} else {
		vPrint("Received response: %s", response)
	}

}

func main() {
	// Parse command-line flags
	flag.Parse()
	serverAddress := os.Getenv("FORWARDER_HOST")
	vPrint("Postdata to server:")
	if serverAddress == "" {
		serverAddress = "server1.gpscontrol.com.mx:8500"
		vPrint("FORWARDER_HOST not set. Using default server address: %s", serverAddress)
	} else {
		vPrint("Using FORWARDER_HOST: %s", serverAddress)
	}
	// MQTT connection options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Fatal("MQTT_BROKER_HOST environment variable not set")
	}
	brokerURL := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	opts.AddBroker(brokerURL)
	clientID := fmt.Sprintf("forwarder_%s_%s_%d",
		"forwarder",
		os.Getenv("HOSTNAME"),
		time.Now().UnixNano()%100000)
	opts.SetClientID(clientID)

	// Configure settings for multiple listeners
	opts.SetCleanSession(false) // Maintain persistent session
	opts.SetAutoReconnect(true) // Auto reconnect on connection loss
	opts.SetKeepAlive(60 * time.Second)
	opts.SetOrderMatters(true) // Maintain message order
	opts.SetResumeSubs(true)   // Resume stored subscriptions

	// Define MQTT message handler
	var mqttClient mqtt.Client
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		// Process each message in a separate goroutine to avoid blocking
		go func() {
			// Increment message counter (for statistics only)
			atomic.AddInt64(&messageCounter, 1)
			
			// Post the received data to the server
			trackerData := string(msg.Payload())
			var messageToForward string
			messageToForward = ""

			// Check if the message is from the tracker/from-tcp topic (JSON format)

			var trackerPayload TrackerPayload
			if err := json.Unmarshal(msg.Payload(), &trackerPayload); err == nil {
				// Successfully parsed JSON, use the payload field
				trackerData = trackerPayload.Payload
				//vPrint("Extracted payload from JSON: %s", trackerData)
			} else {
				vPrint("Failed to parse JSON from tracker/from-tcp: %v", err)
			}

			// Try to decode as hex, if it fails, use the original message
			bytes, err := hex.DecodeString(trackerData)
			if err != nil {
				// Not a hex string, use the original message
				messageToForward = trackerData
			} else {
				// Successfully decoded hex string
				messageToForward = string(bytes)
			}
			tracker_bytes := []byte(messageToForward)
			vPrint("Received message on topic %s:\n%v", msg.Topic(), hex.Dump(tracker_bytes[:min(32, len(tracker_bytes))]))
			postDataToServer(messageToForward, serverAddress)
		}()
	})

	// Attempt to connect to the MQTT broker in a loop until successful
	for {
		mqttClient = mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			vPrint("Error connecting to MQTT broker at %s: %v. Retrying in 5 seconds...", brokerURL, token.Error())
			time.Sleep(5 * time.Second) // Wait before retrying
			continue                    // Retry the connection
		}
		vPrint("Successfully connected to the MQTT broker")
		break // Exit the loop once connected
	}

	// Subscribe to both UDP and TCP topics
	topics := []string{"tracker/from-udp", "tracker/from-tcp"}
	for _, topic := range topics {
		vPrint("Subscribe topic: %s", topic)
		if token := mqttClient.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
			log.Fatalf("Error subscribing to topic %s: %v", topic, token.Error())
		}
		vPrint("Subscribed to topic: %s", topic)
	}

	// Start a goroutine to report message processing statistics
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				count := atomic.LoadInt64(&messageCounter)
				vPrint("Messages processed so far: %d", count)
			}
		}
	}()

	// Keep the application running
	select {}
}
