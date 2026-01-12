package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var sqlLog string

type TrackerData struct {
	Payload    string `json:"payload"`
	RemoteAddr string `json:"remoteaddr"`
}
type TrackerAssign struct {
	Imei       string `json:"imei"`
	Protocol   string `json:"protocol"`
	RemoteAddr string `json:"remoteaddr"`
}
type Logger struct {
	Verbose bool
}

type ConnectionInfo struct {
	RemoteAddr string
	Protocol   string
}

func (l *Logger) Println(v ...interface{}) {
	if l.Verbose {
		log.Println(v...)
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if l.Verbose {
		log.Printf(format, v...)
	}
}

// Global connection map with mutex for thread safety
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

type SendCommandRequest struct {
	IMEI string `json:"imei" binding:"required"`
	Data string `json:"data" binding:"required"` // Hex encoded data
}

func setupRouter(logger *Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/trackerlist", func(c *gin.Context) {
			connMutex.Lock()
			trackers := make([]TrackerAssign, 0)
			for imei, connInfo := range imeiConnections {
				tracker := TrackerAssign{
					Imei:       imei,
					Protocol:   connInfo.Protocol,
					RemoteAddr: connInfo.RemoteAddr,
				}
				trackers = append(trackers, tracker)
			}
			connMutex.Unlock()
			c.JSON(http.StatusOK, trackers)
		})

		v1.POST("/sendcommand", func(c *gin.Context) {
			var req SendCommandRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
				return
			}

			// Check if IMEI exists
			connMutex.Lock()
			connInfo, exists := imeiConnections[req.IMEI]
			connMutex.Unlock()

			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"error": "Tracker not found"})
				return
			}

			// Log the received hex string
			logger.Printf("Received hex string for IMEI %s: %s", req.IMEI, req.Data)

			// Convert hex string to bytes
			data, err := hex.DecodeString(req.Data)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hex data"})
				return
			}

			// Log the decoded data
			logger.Printf("Decoded data: %s", string(data))
			logger.Printf("Decoded hex dump:\n%s", hex.Dump(data))

			// Send data to connection
			err = SendDataToConnection(connInfo.RemoteAddr, data, logger)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Command sent successfully"})
		})
	}

	return r
}

// SendDataToConnection sends data to a specific TCP connection if it exists
func SendDataToConnection(remoteAddr string, data []byte, logger *Logger) error {
	connMutex.Lock()
	conn, exists := activeConnections[remoteAddr]
	connMutex.Unlock()

	if !exists {
		return fmt.Errorf("no active connection for address: %s", remoteAddr)
	}

	// Log what we're about to send
	//logger.Printf("Sending to %s - Data as string: %s", remoteAddr, string(data))
	//logger.Printf("Sending to %s - Hex dump:\n%s", remoteAddr, hex.Dump(data))

	_, err := conn.Write(data)
	if err != nil {
		logger.Printf("Error writing to connection %s: %v\n", remoteAddr, err)
		return err
	}

	//logger.Printf("Successfully sent %d bytes to %s", len(data), remoteAddr)
	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func getMySQLDSN() string {
	host := getEnvWithDefault("MYSQL_HOST", "host.minikube.internal")
	port := getEnvWithDefault("MYSQL_PORT", "3306")
	user := getEnvWithDefault("MYSQL_USER", "gpscontrol")
	pass := getEnvWithDefault("MYSQL_PASS", "qazwsxedc")
	dbname := getEnvWithDefault("MYSQL_DB", "bridge")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)
}

func main() {
	verbose := flag.Bool("v", false, "Enable verbose logging")
	flag.Parse()
	logger := Logger{Verbose: *verbose}
	sqlLog = os.Getenv("SQL_LOG")
	if sqlLog == "" {
		sqlLog = "False"
	}

	if sqlLog != "False" {
		logger.Printf("Sql Log True")
		dsn := getMySQLDSN()
		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Error al cargar el driver mysql: %v", err)
			return
		}
		defer db.Close()
		// Setup logging

		// Create devices table if it doesn't exist
		logger.Println("Creating devices table if not exists...")

		createDevicesSQL := `CREATE TABLE IF NOT EXISTS rawdevices (
			imei varchar(255) NOT NULL,
			imei_count int(11) DEFAULT 0,
			PRIMARY KEY (imei)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1`

		_, err = db.Exec(createDevicesSQL)
		if err != nil {
			log.Printf("Error creating jonodevices table: %v", err)
			return
		}
	} else {
		logger.Printf("Sql Log False")
	}

	// Setup MQTT client options
	mqttBrokerHost := os.Getenv("MQTT_BROKER_HOST")
	opts := mqtt.NewClientOptions()
	opts.SetClientID("udp-tcp-listener")
	opts.AddBroker(fmt.Sprintf("tcp://%s:1883", mqttBrokerHost))

	// Create and start a client using the above ClientOptions
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Println("Error connecting to MQTT broker:", token.Error())
		return
	}

	// Subscribe to tracker/send topic
	token := client.Subscribe("tracker/send", 0, func(client mqtt.Client, msg mqtt.Message) {
		var data TrackerData
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			logger.Println("Error unmarshaling MQTT message:", err)
			return
		}

		// Convert payload from hex string to bytes
		rawBytes, err := hex.DecodeString(data.Payload)
		if err != nil {
			logger.Println("Error decoding hex payload:", err)
			return
		}

		// Send data to the connection
		if err := SendDataToConnection(data.RemoteAddr, rawBytes, &logger); err != nil {
			logger.Println("Error sending data to connection:", err)
			return
		}
	})

	if token.Wait() && token.Error() != nil {
		logger.Println("Error subscribing to MQTT topic:", token.Error())
		return
	}

	// Subscribe to tracker/send topic
	assign2remoteaddr := client.Subscribe("tracker/assign-imei2remoteaddr", 0, func(client mqtt.Client, msg mqtt.Message) {
		var data TrackerAssign
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			logger.Println("Error unmarshaling JSON:", err)
			return
		}

		// Convert hex string to bytes
		imei := data.Imei
		protocol := data.Protocol
		logger.Printf("Assigning imei: %s to address: %s from protocol %s\n", imei, data.RemoteAddr, protocol)
		connMutex.Lock()
		imeiConnections[imei] = ConnectionInfo{
			RemoteAddr: data.RemoteAddr,
			Protocol:   protocol,
		}
		connMutex.Unlock()
	})

	if assign2remoteaddr.Wait() && assign2remoteaddr.Error() != nil {
		logger.Println("Error subscribing to topic:", token.Error())
		return
	}

	// Use goroutines to handle both TCP and UDP listeners simultaneously
	go startTCPListener("0.0.0.0:1024", client, &logger)
	go startUDPListener("0.0.0.0:1024", client, &logger)

	// Start API server
	router := setupRouter(&logger)
	go func() {
		if err := router.Run(":8080"); err != nil {
			logger.Printf("Failed to start API server: %v", err)
			os.Exit(1)
		}
	}()

	// Keep the main goroutine alive
	select {}
}

func startTCPListener(address string, client mqtt.Client, logger *Logger) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Println("Error starting TCP listener:", err)
		return
	}
	defer listener.Close()

	logger.Println("TCP server listening on", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println("Error accepting TCP connection:", err)
			continue
		}
		go handleTCPConnection(conn, client, logger)
	}
}

func handleTCPDisconnection(conn net.Conn, client mqtt.Client, logger *Logger, remoteAddr string) {
	connMutex.Lock()
	delete(activeConnections, remoteAddr)
	// Clean up IMEI mapping
	for imei, addr := range imeiConnections {
		if addr.RemoteAddr == remoteAddr {
			delete(imeiConnections, imei)
			logger.Printf("Removed IMEI mapping for disconnected device: %s\n", imei)
		}
	}
	connMutex.Unlock()

	conn.Close()
	logger.Printf("Connection closed for %s\n", remoteAddr)
}

func handleTCPConnection(conn net.Conn, client mqtt.Client, logger *Logger) {
	remoteAddr := conn.RemoteAddr().String()

	// Register connection
	connMutex.Lock()
	activeConnections[remoteAddr] = conn
	connMutex.Unlock()

	defer func() {
		handleTCPDisconnection(conn, client, logger, remoteAddr)
	}()

	logger.Println("New TCP connection established from", remoteAddr)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			logger.Println("Error reading from TCP connection:", err)
			return
		}

		logger.Printf("TCP Received from %s:", remoteAddr)
		//logger.Printf("%s", hex.Dump(buffer[:min(32, n)]))
		logger.Println(hex.Dump(buffer[:n]))
		// Convert buffer to hex string and publish to MQTT
		hexString := hex.EncodeToString(buffer[:n])
		data := TrackerData{
			Payload:    hexString,
			RemoteAddr: remoteAddr,
		}
		byte_tracker_data_json, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error in byte_tracker_data_json creating JSON: %v", err)
			return
		}
		tracker_data_json := string(byte_tracker_data_json)
		client.Publish("tracker/from-tcp", 0, false, tracker_data_json)
		logger.Println("Successfully sent data MQTT: tracker/from-tcp")
	}
}

func startUDPListener(address string, client mqtt.Client, logger *Logger) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		logger.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Println("Error starting UDP listener:", err)
		return
	}
	defer conn.Close()

	logger.Println("UDP server listening on", address)

	for {
		// Clear the buffer before reading
		buffer := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			logger.Println("Error reading from UDP connection:", err)
			continue
		}
		// logger.Printf("UDP: Received %s from %s\n", string(buffer[:n]), remoteAddr)
		logger.Printf("UDP: Received from %s:\n%s", remoteAddr, hex.Dump(buffer[:min(32, n)]))
		// Convert buffer to hex string and publish to MQTT
		hexString := hex.EncodeToString(buffer[:n])
		client.Publish("tracker/from-udp", 0, false, hexString)
		logger.Printf("Successfully sent data MQTT: tracker/from-udp")
		if sqlLog != "False" {
			if db == nil {
				logger.Println("Database connection not initialized")
				continue
			}
			fields := strings.Split(string(buffer), ",")
			if len(fields) > 1 {
				imei := fields[1]
				// First try to update existing record
				updateSQL := `
				INSERT INTO rawdevices (imei, imei_count, payload) 
				VALUES (?, 1, ?) 
				ON DUPLICATE KEY UPDATE 
				imei_count = imei_count + 1, 
				payload = VALUES(payload)`

				// Execute the query with imei and hexString as parameters
				result, err2 := db.Exec(updateSQL, imei, hexString)
				if err2 != nil {
					logger.Printf("Error updating database: %v\n", err2)
					return
				}

				rowsAffected, _ := result.RowsAffected()
				logger.Printf("Rows affected: %d\n", rowsAffected)
			}
		}
	}
}
