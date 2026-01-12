package main

import (
	activetrackintegrator "activetrack/features/activetrack_integrator"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MaddSystems/jonobridge/common/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TrackerData represents the JSON structure to send via HTTP
type TrackerData struct {
	TrackerData string `json:"trackerdata"`
}

func main() {
	// Parse command-line flags
	flag.Parse()

	// Configure and initialize Activetrack before anything else
	activetrackToken := os.Getenv("ACTIVETRACK_TOKEN")
	activetrackURL := os.Getenv("ACTIVETRACK_URL")
	fmt.Println("token length:", len(activetrackToken))
	fmt.Println("url:", activetrackURL)

	// Initialize the integrator explicitly
	activetrackintegrator.Initialize()

	// Initialize MQTT client options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Fatal("MQTT_BROKER_HOST environment variable not set")
	}
	brokerURL := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	opts.AddBroker(brokerURL)
	clientID := fmt.Sprintf("ACTIVETRACK_%s_%s_%d",
		"ACTIVETRACK",
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
		payload := msg.Payload()
		displayLen := 80
		if len(payload) < displayLen {
			displayLen = len(payload)
		}
		utils.VPrint("Received message on topic %s: %s", msg.Topic(), payload[:displayLen])

		// INSERT HERE ROUTINE TO DO SOMETHING WITH
		activetrackintegrator.Init(string(msg.Payload()))
		// That contains jono protocol
	})

	// Attempt to connect to the MQTT broker in a loop until successful

	for {
		mqttClient = mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			//utils.VPrint("Error connecting to MQTT broker at %s: %v. Retrying in 5 seconds...", brokerURL, token.Error())
			time.Sleep(5 * time.Second) // Wait before retrying
			continue                    // Retry the connection
		}
		utils.VPrint("Successfully connected to the MQTT broker")
		break // Exit the loop once connected
	}

	// Subscribe to the topic
	topic := "tracker/jonoprotocol"
	if token := mqttClient.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic %s: %v", topic, token.Error())
	}
	utils.VPrint("Subscribed to topic: %s", topic)

	// Start a timer to periodically send all pending data - but less frequently
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			utils.VPrint("Minute timer triggered, sending all pending data...")
			activetrackintegrator.SendAllPendingData()
		}
	}()

	// Keep the application running
	select {}
}
