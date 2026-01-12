package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	localAddress  = flag.String("l", "0.0.0.0:1024", "Local address")
	remoteAddress = flag.String("r", os.Getenv("PLATFORM_HOST"), "Remote address")
	verbose       = flag.Bool("v", false, "Enable verbose logging")
	mqttClient    mqtt.Client
)

// Helper function to print verbose logs if enabled
func vPrint(format string, v ...interface{}) {
	if *verbose {
		log.Printf(format, v...)
	}
}

type TrackerData struct {
	Payload    string `json:"payload"`
	RemoteAddr string `json:"remoteaddr"`
}

type TrackerAssign struct {
	Imei       string `json:"imei"`
	Protocol   string `json:"protocol"`
	RemoteAddr string `json:"remoteaddr"`
}

type ConnectionInfo struct {
	RemoteAddr string
	Protocol   string
}

var (
	activeConnections = make(map[string]net.Conn)
	imeiConnections   = make(map[string]ConnectionInfo) // maps IMEI to ConnectionInfo
	connMutex         sync.Mutex
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SendDataToConnection(remoteAddr string, data []byte) error {
	connMutex.Lock()
	conn, exists := activeConnections[remoteAddr]
	connMutex.Unlock()
	fmt.Println("activeConnections: ", activeConnections)
	if !exists {
		return fmt.Errorf("no active connection for address: %s", remoteAddr)
	}

	_, err := conn.Write(data)
	if err != nil {
		vPrint("Error writing to connection %s: %v", remoteAddr, err)
		return err
	}
	return nil
}

func init() {
	// MQTT connection options
	opts := mqtt.NewClientOptions()
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	if mqttBrokerHost == "" {
		mqttBrokerHost = "localhost" // default value
	}

	opts.AddBroker(fmt.Sprintf("tcp://%s:1883", mqttBrokerHost))
	opts.SetClientID("proxy_service")
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		if token := client.Subscribe("tracker/assign", 0, handleAssignment); token.Wait() && token.Error() != nil {
			log.Printf("Failed to subscribe to tracker/assign: %v", token.Error())
		}
	})

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Failed to connect to MQTT broker: %v", token.Error())
	}
}

func handleAssignment(client mqtt.Client, msg mqtt.Message) {
	var assign TrackerAssign
	if err := json.Unmarshal(msg.Payload(), &assign); err != nil {
		log.Printf("Error unmarshaling assignment message: %v", err)
		return
	}

	if assign.Imei == "" || assign.RemoteAddr == "" {
		log.Printf("Invalid assignment message: missing IMEI or RemoteAddr")
		return
	}

	connMutex.Lock()
	defer connMutex.Unlock()

	// Verify the connection exists
	if _, exists := activeConnections[assign.RemoteAddr]; !exists {
		log.Printf("Cannot assign IMEI %s: no active connection for %s", assign.Imei, assign.RemoteAddr)
		return
	}

	// Update the IMEI mapping
	imeiConnections[assign.Imei] = ConnectionInfo{
		RemoteAddr: assign.RemoteAddr,
		Protocol:   assign.Protocol,
	}
	vPrint("Assigned IMEI %s to connection %s with protocol %s", assign.Imei, assign.RemoteAddr, assign.Protocol)
}

func main() {
	flag.Parse()

	fmt.Printf("Listening: %v\nProxying %v\n", *localAddress, *remoteAddress)

	// Subscribe to tracker/send topic
	token := mqttClient.Subscribe("tracker/send", 0, func(client mqtt.Client, msg mqtt.Message) {
		var data TrackerData
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			vPrint("Error unmarshaling MQTT message: %v", err)
			return
		}

		// Convert payload from hex string to bytes
		rawBytes, err := hex.DecodeString(data.Payload)
		if err != nil {
			vPrint("Error decoding hex payload: %v", err)
			return
		}

		// Send data to the connection
		if err := SendDataToConnection(data.RemoteAddr, rawBytes); err != nil {
			vPrint("Error sending data to connection: %v", err)
			return
		}
	})

	if token.Wait() && token.Error() != nil {
		vPrint("Error subscribing to MQTT topic: %v", token.Error())
		return
	}

	// Subscribe to tracker/send topic
	assign2remoteaddr := mqttClient.Subscribe("tracker/assign-imei2remoteaddr", 0, func(client mqtt.Client, msg mqtt.Message) {
		var data TrackerAssign
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			vPrint("Error unmarshaling JSON: %v", err)
			return
		}

		// Convert hex string to bytes
		imei := data.Imei
		protocol := data.Protocol
		vPrint("Assigning imei: %s to address: %s from protocol %s", imei, data.RemoteAddr, protocol)
		connMutex.Lock()
		imeiConnections[imei] = ConnectionInfo{
			RemoteAddr: data.RemoteAddr,
			Protocol:   protocol,
		}
		connMutex.Unlock()
	})

	if assign2remoteaddr.Wait() && assign2remoteaddr.Error() != nil {
		vPrint("Error subscribing to topic: %v", assign2remoteaddr.Error())
		return
	}

	// Start the HTTP server
	go func() {
		r := setupRouter()
		if err := r.Run(":8080"); err != nil {
			vPrint("Error starting HTTP server: %v", err)
		}
	}()

	addr, err := net.ResolveTCPAddr("tcp", *localAddress)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go proxyConnection(conn)
	}
}

// Deregister a client connection
func deregisterClient(remoteAddr string) {
	connMutex.Lock()
	defer connMutex.Unlock()

	// First, delete from activeConnections
	if conn, exists := activeConnections[remoteAddr]; exists {
		conn.Close()
		delete(activeConnections, remoteAddr)
	}

	// Find and remove any IMEI that maps to this remoteAddr
	var imeisToDelete []string
	for imei, connInfo := range imeiConnections {
		if connInfo.RemoteAddr == remoteAddr {
			imeisToDelete = append(imeisToDelete, imei)
		}
	}

	// Delete the found IMEIs
	for _, imei := range imeisToDelete {
		delete(imeiConnections, imei)
		vPrint("Deregistered IMEI %s due to connection close from %s", imei, remoteAddr)
	}
}

func proxyConnection(conn *net.TCPConn) {
	remoteAddr := conn.RemoteAddr().String()
	defer conn.Close()

	vPrint("New connection from: %s", remoteAddr)

	// Just add to activeConnections, but don't register IMEI yet
	connMutex.Lock()
	activeConnections[remoteAddr] = conn
	connMutex.Unlock()

	rAddr, err := net.ResolveTCPAddr("tcp", *remoteAddress)
	if err != nil {
		vPrint("Failed to resolve remote address: %v", err)
		deregisterClient(remoteAddr)
		return
	}

	rConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		vPrint("Failed to connect to remote server: %v", err)
		deregisterClient(remoteAddr)
		return
	}
	defer rConn.Close()

	// Create channels to signal connection closure
	clientClosed := make(chan struct{})
	serverClosed := make(chan struct{})

	// Forward data from client to remote server
	go func() {
		defer close(clientClosed)
		for {
			buffer := make([]byte, 1024*1024)
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					vPrint("Error reading from client: %v", err)
				}
				return
			}

			if n > 0 {
				// Forward data to the remote server
				_, err = rConn.Write(buffer[:n])
				if err != nil {
					vPrint("Error writing to remote server: %v", err)
					return
				}
				vPrint("From tracker to platform:\n%s", hex.Dump(buffer[:n]))
				//vPrint("From tracker to platform:\n%s", hex.Dump(buffer[:min(32, n)]))
				//vPrint("From tracker to platform:\n%s", hex.Dump(buffer[:n]))
				// Prepare and publish MQTT message
				hexString := hex.EncodeToString(buffer[:n])
				trackerData := TrackerData{
					Payload:    hexString,
					RemoteAddr: remoteAddr,
				}
				byte_tracker_data_json, err := json.Marshal(trackerData)
				if err != nil {
					log.Printf("Error in byte_tracker_data_json creating JSON: %v", err)
					return
				}

				// Check if MQTT client is available before publishing
				if mqttClient != nil && mqttClient.IsConnected() {
					tracker_data_json := string(byte_tracker_data_json)
					if token := mqttClient.Publish("tracker/from-tcp", 0, false, tracker_data_json); token.Wait() && token.Error() != nil {
						vPrint("Error publishing to MQTT: %v", token.Error())
					}
				} else {
					vPrint("MQTT client not available or not connected")
				}
			}
		}
	}()

	// Forward data from remote server to client
	go func() {
		defer close(serverClosed)
		buffer := make([]byte, 1024*1024)
		for {
			n, err := rConn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					vPrint("Error reading from remote server: %v", err)
				}
				return
			}

			if n > 0 {
				_, err = conn.Write(buffer[:n])
				if err != nil {
					vPrint("Error writing to client: %v", err)
					return
				}
				vPrint("From platform to tracker:\n%s", hex.Dump(buffer[:min(32, n)]))
			}
		}
	}()

	// Wait for either connection to close
	select {
	case <-clientClosed:
		vPrint("Client connection closed: %s", remoteAddr)
	case <-serverClosed:
		vPrint("Server connection closed for client: %s", remoteAddr)
	}

	// Clean up the connection
	deregisterClient(remoteAddr)
}
