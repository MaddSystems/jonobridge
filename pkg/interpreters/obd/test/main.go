package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Decode JT/T 808 response
func decodeTrackerResponse(data []byte) {
	if len(data) < 12 {
		log.Println("Invalid packet length")
		return
	}

	// Ensure start/end bytes are correct
	if data[0] != 0x7E || data[len(data)-1] != 0x7E {
		log.Println("Invalid packet format")
		return
	}

	// Extract header fields
	messageID := binary.BigEndian.Uint16(data[1:3])
	attribute := binary.BigEndian.Uint16(data[3:5])
	terminalID := fmt.Sprintf("%X", data[5:11])
	messageSerial := binary.BigEndian.Uint16(data[11:13])

	fmt.Println("Message ID:", messageID)
	fmt.Println("Attribute:", attribute)
	fmt.Println("Terminal ID:", terminalID)
	fmt.Println("Message Serial:", messageSerial)

	// Extract location data (0200 message)
	if messageID == 0x0200 {
		status := binary.BigEndian.Uint32(data[13:17])
		latitude := binary.BigEndian.Uint32(data[17:21])  // GPS latitude
		longitude := binary.BigEndian.Uint32(data[21:25]) // GPS longitude
		speed := binary.BigEndian.Uint16(data[25:27])
		direction := binary.BigEndian.Uint16(data[27:29])
		timestamp := fmt.Sprintf("%02X-%02X-%02X %02X:%02X:%02X",
			data[29], data[30], data[31], data[32], data[33], data[34])

		fmt.Println("Status:", status)
		fmt.Println("Latitude:", float64(latitude)/1000000)   // Convert to degrees
		fmt.Println("Longitude:", float64(longitude)/1000000) // Convert to degrees
		fmt.Println("Speed (km/h):", float64(speed)/10)
		fmt.Println("Direction:", direction)
		fmt.Println("Timestamp:", timestamp)
	}

	// TODO: Decode extra extensions if present
}

// Create general response message (8001)
func createGeneralResponse(messageID uint16, messageSerial uint16, terminalID []byte, result uint8) []byte {
	// Construct response according to JT/T 808 protocol
	responseData := make([]byte, 0, 16)

	// Start marker
	responseData = append(responseData, 0x7E)

	// Message ID (8001 - general response)
	responseData = append(responseData, 0x80, 0x01)

	// Message attributes (length to be filled later)
	responseData = append(responseData, 0x00, 0x05)

	// Terminal ID (copy from request)
	responseData = append(responseData, terminalID...)

	// Message serial
	serialBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(serialBytes, messageSerial)
	responseData = append(responseData, serialBytes...)

	// Response message serial
	responseSerialBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(responseSerialBytes, messageSerial)
	responseData = append(responseData, responseSerialBytes...)

	// Original message ID
	msgIDBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(msgIDBytes, messageID)
	responseData = append(responseData, msgIDBytes...)

	// Result (0: success)
	responseData = append(responseData, result)

	// Calculate checksum (XOR of all bytes except start and end markers)
	var checksum byte
	for i := 1; i < len(responseData); i++ {
		checksum ^= responseData[i]
	}
	responseData = append(responseData, checksum)

	// End marker
	responseData = append(responseData, 0x7E)

	return responseData
}

// Find complete messages in buffer (packets may be split or combined)
func findCompleteMessages(buffer []byte) ([][]byte, []byte) {
	var messages [][]byte
	remaining := buffer

	for {
		// Find start marker
		startIndex := -1
		for i, b := range remaining {
			if b == 0x7E {
				startIndex = i
				break
			}
		}

		if startIndex == -1 {
			// No start marker found, keep all data
			return messages, remaining
		}

		// Discard data before start marker
		if startIndex > 0 {
			remaining = remaining[startIndex:]
			startIndex = 0
		}

		// Find end marker
		endIndex := -1
		for i := startIndex + 1; i < len(remaining); i++ {
			if remaining[i] == 0x7E {
				endIndex = i
				break
			}
		}

		if endIndex == -1 {
			// No end marker found, wait for more data
			return messages, remaining
		}

		// Extract complete message
		message := remaining[startIndex : endIndex+1]
		messages = append(messages, message)

		// Update remaining buffer
		remaining = remaining[endIndex+1:]

		// If no more data, exit loop
		if len(remaining) == 0 {
			break
		}
	}

	return messages, remaining
}

// Print formatted hex dump with timestamp
func printHexDumpWithTimestamp(prefix string, data []byte) {
	// Get current time
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Convert data to hex dump format
	hexStr := hex.EncodeToString(data)

	// Format hex dump with spaces for better readability (2 bytes/characters per group)
	var formattedHex string
	for i := 0; i < len(hexStr); i += 2 {
		if i+2 <= len(hexStr) {
			formattedHex += hexStr[i:i+2] + " "
		} else {
			formattedHex += hexStr[i:] + " "
		}
	}

	// Print formatted output with ASCII representation
	log.Printf("%s [%s] %s\n", timestamp, prefix, formattedHex)

	// Add ASCII representation
	var sb strings.Builder
	sb.WriteString("ASCII: ")
	for _, b := range data {
		if b >= 32 && b <= 126 { // Printable ASCII characters
			sb.WriteByte(b)
		} else {
			sb.WriteByte('.')
		}
	}
	log.Printf("%s [%s] %s\n", timestamp, prefix, sb.String())
}

// Handle connection from tracker
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Log client information
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("New connection from %s", remoteAddr)

	var buffer []byte
	readBuffer := make([]byte, 4096)

	// Set read deadline
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	if err != nil {
		log.Printf("Error setting read deadline: %v", err)
		return
	}

	for {
		// Read data from connection
		n, err := conn.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by %s", remoteAddr)
			} else {
				log.Printf("Error reading from %s: %v", remoteAddr, err)
			}
			break
		}

		// Log raw input with timestamp
		printHexDumpWithTimestamp("RAW INPUT", readBuffer[:n])

		// Append new data to buffer
		buffer = append(buffer, readBuffer[:n]...)

		// Find complete messages
		messages, remaining := findCompleteMessages(buffer)
		buffer = remaining

		// Process complete messages
		for _, message := range messages {
			// Log complete message with timestamp
			printHexDumpWithTimestamp("COMPLETE MESSAGE", message)

			// Try to decode message
			decodeTrackerResponse(message)

			// Check if message has valid format
			if len(message) >= 13 && message[0] == 0x7E && message[len(message)-1] == 0x7E {
				// Extract message ID and serial
				messageID := binary.BigEndian.Uint16(message[1:3])
				messageSerial := binary.BigEndian.Uint16(message[11:13])
				terminalID := message[5:11]

				// Send response based on message type
				var response []byte

				switch messageID {
				case 0x0100: // Registration
					// Send registration response (8100)
					log.Printf("Registration request from %s", remoteAddr)
					// In real implementation, you would generate authentication code
					// For now, just send general response with success
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				case 0x0200: // Location report
					// Send general response (8001)
					log.Printf("Location report from %s", remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				case 0x0002: // Heartbeat
					log.Printf("Heartbeat from %s", remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				default:
					// For any other message type, send general response (8001)
					log.Printf("Unknown message type %04X from %s", messageID, remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				}

				// Send response
				if response != nil {
					_, err := conn.Write(response)
					if err != nil {
						log.Printf("Error sending response to %s: %v", remoteAddr, err)
						break
					}
					// Log outgoing response with timestamp
					printHexDumpWithTimestamp("RESPONSE", response)
				}
			}
		}

		// Reset read deadline
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		if err != nil {
			log.Printf("Error resetting read deadline: %v", err)
			break
		}
	}
}

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "0.0.0.0:8600")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	log.Println("Server started on 0.0.0.0:8600")

	// Track active connections
	var wg sync.WaitGroup

	for {
		// Accept connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Handle connection in goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConnection(conn)
		}()
	}
}
