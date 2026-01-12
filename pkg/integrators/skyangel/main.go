package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	skyangelintegrator "skyangel/features/skyangel_integrator"
	"sync"
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
	// Parse command-line flags
	flag.Parse()
	
	// Initialize the integrator explicitly
	skyangelintegrator.Initialize()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize MQTT client options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		log.Fatal("MQTT_BROKER_HOST environment variable not set")
	}
	brokerURL := fmt.Sprintf("tcp://%s:1883", mqttBrokerHost)
	opts.AddBroker(brokerURL)
	clientID := fmt.Sprintf("SKYANGEL_%s_%s_%d",
		"SKYANGEL",
		os.Getenv("HOSTNAME"),
		time.Now().UnixNano()%100000)
	opts.SetClientID(clientID)

	// Configure settings for multiple listeners
	opts.SetCleanSession(false) // Maintain persistent session
	opts.SetAutoReconnect(true) // Auto reconnect on connection loss
	opts.SetKeepAlive(60 * time.Second)
	opts.SetOrderMatters(true) // Maintain message order
	opts.SetResumeSubs(true)   // Resume stored subscriptions
	
	// Add connection timeouts
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetWriteTimeout(10 * time.Second)

	// WaitGroup to handle concurrent message processing
	var wg sync.WaitGroup
	var mqttClient mqtt.Client

	// Define MQTT message handler with timeout and goroutine
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		// Process messages concurrently to avoid blocking
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			// Create timeout context for message processing
			msgCtx, msgCancel := context.WithTimeout(ctx, 30*time.Second)
			defer msgCancel()

			// Channel to handle the result
			done := make(chan error, 1)
			
			// Run the integrator in a goroutine with timeout
			go func() {
				done <- skyangelintegrator.Init(string(msg.Payload()))
			}()

			// Wait for completion or timeout
			select {
			case err := <-done:
				if err != nil {
					utils.VPrint("Error processing message: %v", err)
				}
			case <-msgCtx.Done():
				utils.VPrint("Message processing timed out for topic %s", msg.Topic())
			}
		}()
	})

	// Attempt to connect to the MQTT broker with timeout and context
	connectionTimeout := time.After(5 * time.Minute) // Max 5 minutes to establish connection

	for {
		select {
		case <-ctx.Done():
			utils.VPrint("Shutting down before MQTT connection established")
			return
		case <-connectionTimeout:
			log.Fatal("Failed to connect to MQTT broker within timeout period")
		default:
			mqttClient = mqtt.NewClient(opts)
			if token := mqttClient.Connect(); token.WaitTimeout(30*time.Second) && token.Error() != nil {
				utils.VPrint("Error connecting to MQTT broker at %s: %v. Retrying in 5 seconds...", brokerURL, token.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			utils.VPrint("Successfully connected to the MQTT broker")
			goto connected
		}
	}

connected:
	// Subscribe to the topic
	topic := "tracker/jonoprotocol"
	if token := mqttClient.Subscribe(topic, 1, nil); token.WaitTimeout(10*time.Second) && token.Error() != nil {
		log.Fatalf("Error subscribing to topic %s: %v", topic, token.Error())
	}
	utils.VPrint("Subscribed to topic: %s", topic)

	// Wait for shutdown signal or context cancellation
	select {
	case sig := <-sigChan:
		utils.VPrint("Received signal %v, shutting down gracefully...", sig)
	case <-ctx.Done():
		utils.VPrint("Context cancelled, shutting down...")
	}

	// Graceful shutdown
	utils.VPrint("Shutting down...")
	
	// Unsubscribe and disconnect
	if token := mqttClient.Unsubscribe(topic); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		utils.VPrint("Error unsubscribing: %v", token.Error())
	}
	
	mqttClient.Disconnect(250) // Wait 250ms for disconnect
	
	// Wait for pending message processing to complete (with timeout)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		utils.VPrint("All messages processed")
	case <-time.After(10 * time.Second):
		utils.VPrint("Timeout waiting for message processing to complete")
	}
	
	utils.VPrint("Shutdown complete")
}