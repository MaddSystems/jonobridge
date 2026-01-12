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

	log.Printf("Decoding message - Type: 0x%04X, Attribute: 0x%04X, Terminal ID: %s, Serial: %d",
		messageID, attribute, terminalID, messageSerial)

	// Check if the message is fragmented (packetized)
	isFragmented := (attribute & 0x2000) != 0
	bodyLength := attribute & 0x03FF // Bits 0-9 represent the body length

	log.Printf("Message body length: %d bytes, Fragmented: %v", bodyLength, isFragmented)

	var bodyStartIndex int = 13

	if isFragmented {
		totalPackets := binary.BigEndian.Uint16(data[13:15])
		currentPacket := binary.BigEndian.Uint16(data[15:17])
		bodyStartIndex = 17
		log.Printf("Fragmented message: Packet %d of %d", currentPacket, totalPackets)
	}

	// Extract location data (0200 message)
	if messageID == 0x0200 {
		if len(data) < bodyStartIndex+20 {
			log.Println("Location message too short")
			return
		}

		// Calculate body end index (excluding checksum and end marker)
		bodyEndIndex := len(data) - 2

		// Extract location information
		alarmFlag := binary.BigEndian.Uint32(data[bodyStartIndex : bodyStartIndex+4])
		status := binary.BigEndian.Uint32(data[bodyStartIndex+4 : bodyStartIndex+8])
		latitude := binary.BigEndian.Uint32(data[bodyStartIndex+8 : bodyStartIndex+12])
		longitude := binary.BigEndian.Uint32(data[bodyStartIndex+12 : bodyStartIndex+16])
		speed := binary.BigEndian.Uint16(data[bodyStartIndex+16 : bodyStartIndex+18])
		direction := binary.BigEndian.Uint16(data[bodyStartIndex+18 : bodyStartIndex+20])

		// Check if there's enough data for timestamp
		if bodyStartIndex+26 <= bodyEndIndex {
			timestamp := fmt.Sprintf("20%02X-%02X-%02X %02X:%02X:%02X",
				data[bodyStartIndex+20], data[bodyStartIndex+21], data[bodyStartIndex+22],
				data[bodyStartIndex+23], data[bodyStartIndex+24], data[bodyStartIndex+25])

			log.Println("======= LOCATION DATA =======")
			log.Printf("Alarm Flag: 0x%08X", alarmFlag)
			log.Printf("Status: 0x%08X", status)
			log.Printf("Latitude: %f (raw: %d)", float64(latitude)/1000000, latitude)
			log.Printf("Longitude: %f (raw: %d)", float64(longitude)/1000000, longitude)
			log.Printf("Speed: %f km/h (raw: %d)", float64(speed)/10, speed)
			log.Printf("Direction: %d degrees", direction)
			log.Printf("Timestamp: %s", timestamp)
			log.Println("============================")
		} else {
			log.Println("Timestamp data not available in location message")
		}

		// Decode additional information items
		if bodyStartIndex+26 < bodyEndIndex {
			additionalInfoStartIndex := bodyStartIndex + 26
			additionalInfo := data[additionalInfoStartIndex:bodyEndIndex]

			log.Println("Additional information items:")
			for len(additionalInfo) > 2 {
				infoID := additionalInfo[0]
				infoLength := int(additionalInfo[1])

				if 2+infoLength > len(additionalInfo) {
					log.Printf("Warning: Additional info item 0x%02X has length %d but only %d bytes remain",
						infoID, infoLength, len(additionalInfo)-2)
					break
				}

				infoData := additionalInfo[2 : 2+infoLength]
				log.Printf("Additional Info ID: 0x%02X, Length: %d, Data: %X", infoID, infoLength, infoData)

				// Special handling for known additional info types
				switch infoID {
				case 0xEA: // Basic data flow
					if len(infoData) >= 3 {
						subID := binary.BigEndian.Uint16(infoData[0:2])
						subLength := infoData[2]
						if len(infoData) >= 3+int(subLength) {
							subData := infoData[3 : 3+int(subLength)]
							log.Printf("Basic data flow - SubID: 0x%04X, Length: %d, Data: %X",
								subID, subLength, subData)
						}
					}
				case 0xEB: // Car extended data flow
					if len(infoData) >= 3 {
						subID := binary.BigEndian.Uint16(infoData[0:2])
						subLength := infoData[2]
						if len(infoData) >= 3+int(subLength) {
							subData := infoData[3 : 3+int(subLength)]
							log.Printf("Car extended data flow - SubID: 0x%04X, Length: %d, Data: %X",
								subID, subLength, subData)
						}
					}
				}

				additionalInfo = additionalInfo[2+infoLength:]
			}
		}
	} else if messageID == 0x0102 { // Terminal authentication
		if bodyStartIndex < len(data)-2 { // Ensure there's data after header
			authCode := data[bodyStartIndex : len(data)-2] // Exclude checksum and end marker
			log.Printf("Terminal authentication - Auth Code: %s", string(authCode))
		}
	} else if messageID == 0x0002 { // Heartbeat
		log.Println("Heartbeat message received")
	} else {
		log.Printf("Message type 0x%04X not specifically handled", messageID)
	}
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

// Create authentication response message (8100)
func createAuthResponse(messageSerial uint16, terminalID []byte, result uint8) []byte {
	// Construct response according to JT/T 808 protocol
	responseData := make([]byte, 0, 20)

	// Start marker
	responseData = append(responseData, 0x7E)

	// Message ID (8100 - authentication response)
	responseData = append(responseData, 0x81, 0x00)

	// Message attributes - will update length later
	responseData = append(responseData, 0x00, 0x05)

	// Terminal ID (copy from request)
	responseData = append(responseData, terminalID...)

	// Message serial
	serialBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(serialBytes, messageSerial)
	responseData = append(responseData, serialBytes...)

	// Serial number of response
	responseData = append(responseData, serialBytes...)

	// Result (0: success)
	responseData = append(responseData, result)

	// Authentication code (only if success)
	authCode := "SUCCESS123"
	if result == 0 {
		responseData = append(responseData, []byte(authCode)...)
	}

	// Update message length in attributes
	bodyLength := uint16(len(responseData) - 3) // Exclude start marker and message ID bytes
	attrBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(attrBytes, bodyLength)
	responseData[3] = attrBytes[0]
	responseData[4] = attrBytes[1]

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
					log.Printf("Registration request from %s", remoteAddr)
					// Send registration response (8100)
					response = createAuthResponse(messageSerial, terminalID, 0)
				case 0x0102: // Authentication
					log.Printf("Authentication request from %s", remoteAddr)
					// Send general response (8001) for authentication
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				case 0x0200: // Location report
					log.Printf("Location report from %s", remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				case 0x0002: // Heartbeat
					log.Printf("Heartbeat from %s", remoteAddr)
					response = createGeneralResponse(messageID, messageSerial, terminalID, 0)
				default:
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
	// Configure logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Println("JT/T 808 Server starting...")

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

	// Wait for all connections to finish
	wg.Wait()
}
