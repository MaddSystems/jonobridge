package main

import (
	"altotrack/features/altotrack_integrator"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaddSystems/jonobridge/common/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TrackerData represents the JSON structure to send via HTTP
type TrackerData struct {
	TrackerData string `json:"trackerdata"`
}

func main() {
	// Recover from panic and log the error
	defer func() {
		if r := recover(); r != nil {
			utils.VPrint("Recovered from panic in main: %v", r)
			os.Exit(1) // Exit with error code
		}
	}()

	// Parse command-line flags
	flag.Parse()

	// Configure and initialize altotrack before anything else
	altotrackintegrator.Initialize()

	// Initialize MQTT client options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Fatal("MQTT_BROKER_HOST environment variable not set")
	}
	brokerURL := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	opts.AddBroker(brokerURL)
	clientID := fmt.Sprintf("ALTOTRACK_%s_%s_%d",
		"ALTOTRACK",
		os.Getenv("HOSTNAME"),
		time.Now().UnixNano()%100000)
	opts.SetClientID(clientID)

	// Configure settings for multiple listeners
	opts.SetCleanSession(false) // Maintain persistent session
	opts.SetAutoReconnect(true) // Auto reconnect on connection loss
	opts.SetKeepAlive(60 * time.Second)
	opts.SetOrderMatters(true) // Maintain message order
	opts.SetResumeSubs(true)   // Resume stored subscriptions
	opts.SetConnectTimeout(10 * time.Second) // Timeout for connection

	// Worker pool to limit concurrent processing
	workerPool := make(chan struct{}, 10) // Limit to 10 concurrent workers

	// Define MQTT message handler with worker pool
	var mqttClient mqtt.Client
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		// Recover from panics in the handler
		defer func() {
			if r := recover(); r != nil {
				utils.VPrint("Recovered from panic in MQTT handler: %v", r)
			}
		}()

		// Attempt to acquire a worker slot
		select {
		case workerPool <- struct{}{}:
			// Run processing in a separate goroutine
			go func() {
				defer func() { <-workerPool }() // Release worker slot
				utils.VPrint("Processing message on topic %s", msg.Topic())
				if err := altotrackintegrator.Init(string(msg.Payload())); err != nil {
					utils.VPrint("Error processing message: %v", err)
				}
			}()
		default:
			utils.VPrint("Worker pool full, dropping message on topic %s", msg.Topic())
		}
	})

	// Set connection lost handler
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		utils.VPrint("Connection to MQTT broker lost: %v", err)
	})

	// Attempt to connect to the MQTT broker with retries
	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		mqttClient = mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.WaitTimeout(10*time.Second) && token.Error() != nil {
			utils.VPrint("Error connecting to MQTT broker at %s: %v. Retrying (%d/%d)...", brokerURL, token.Error(), i+1, maxRetries)
			time.Sleep(5 * time.Second)
			continue
		}
		utils.VPrint("Successfully connected to the MQTT broker")
		break
	}
	if !mqttClient.IsConnected() {
		log.Fatal("Failed to connect to MQTT broker after maximum retries")
	}

	// Subscribe to the topic
	topic := "tracker/jonoprotocol"
	if token := mqttClient.Subscribe(topic, 1, nil); token.WaitTimeout(10*time.Second) && token.Error() != nil {
		log.Fatalf("Error subscribing to topic %s: %v", topic, token.Error())
	}
	utils.VPrint("Subscribed to topic: %s", topic)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	utils.VPrint("Received shutdown signal, disconnecting MQTT client")
	mqttClient.Disconnect(250)
	utils.VPrint("Application terminated gracefully")
}