# SkyWave Protocol to MVT366 Integration Guide

This document provides a comprehensive guide to understanding how the system retrieves data from SkyWave satellite devices and transforms it into the MVT366 format used by the tracking platform.

## Table of Contents

1. [System Overview](#system-overview)
2. [Authentication and Configuration](#authentication-and-configuration)
3. [Data Flow Architecture](#data-flow-architecture)
4. [SkyWave Message Structure](#skywave-message-structure)
5. [Coordinate Transformation](#coordinate-transformation)
6. [Field Mapping](#field-mapping)
7. [MVT366 Message Format](#mvt366-message-format)
8. [Network Communication](#network-communication)
9. [Main Program Flow](#main-program-flow)
10. [Troubleshooting](#troubleshooting)

## System Overview

The system integrates SkyWave satellite terminals with a tracking platform that uses the MVT366 protocol. It consists of three major components:

1. **SkyWave XML API Client**: Retrieves satellite messages via REST API from SkyWave's servers
2. **Message Parser**: Decodes XML messages and extracts position data from specific payload types
3. **MVT366 Formatter**: Converts the parsed data into MVT366 format and transmits via UDP

## Authentication and Configuration

### Required Credentials

The system requires three key parameters for API authentication:

- **Access ID**: `70001184` - SkyWave account identifier
- **Password**: `JEUTPKKH` - API access password  
- **From ID**: Dynamic value read from environment variable `FROMIDSKYWAVE`

### Environment Setup

The system requires the `FROMIDSKYWAVE` environment variable to be set:

```bash
export FROMIDSKYWAVE=13969586728
```

This value changes dynamically as messages are processed and represents the starting point for message retrieval.

### API Endpoint

The system uses the SkyWave API endpoint:

```
https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml/
```

## Data Flow Architecture

1. **Initialization**: Program reads `FROMIDSKYWAVE` environment variable
2. **API Polling**: Continuous polling of SkyWave API for new messages
3. **XML Parsing**: Parse returned XML to extract message structures
4. **Message Filtering**: Filter messages by payload type (only position-related messages)
5. **Coordinate Conversion**: Transform SkyWave coordinate format to decimal degrees
6. **MVT366 Formatting**: Convert data to MVT366 protocol format
7. **UDP Transmission**: Send formatted messages to tracking server at `13.89.38.9:1805`
8. **Pagination**: Update `From_id` for next batch of messages

## SkyWave Message Structure

### XML Response Structure

SkyWave returns messages in XML format with the following hierarchy:

```xml
<GetReturnMessagesResult>
    <ErrorID>0</ErrorID>
    <More>true</More>
    <NextStartID>123456789</NextStartID>
    <Messages>
        <ReturnMessage>
            <ID>987654321</ID>
            <MessageUTC>2025-09-03 18:02:25</MessageUTC>
            <ReceiveUTC>2025-09-03 18:02:30</ReceiveUTC>
            <MobileID>02092247SKY6A70</MobileID>
            <Payload Name="StationaryIntervalSat" SIN="126" MIN="1">
                <Fields>
                    <Field Name="Latitude" Value="1952100"/>
                    <Field Name="Longitude" Value="-9921171"/>
                    <Field Name="Speed" Value="0"/>
                    <Field Name="Heading" Value="361"/>
                    <Field Name="EventTime" Value="230903180225"/>
                </Fields>
            </Payload>
        </ReturnMessage>
    </Messages>
</GetReturnMessagesResult>
```

### Supported Message Types

The system only processes specific payload types that contain position data:

- `DistanceCell`
- `StationaryIntervalSat` 
- `MovingIntervalSat`
- `MovingEnd`
- `MovingStart`
- `IgnitionOn`
- `StationaryIntervalCell`

All other message types are ignored.

### Key Fields Extracted

From each supported message, the following fields are extracted:
- **Latitude**: String representation of latitude (needs conversion)
- **Longitude**: String representation of longitude (needs conversion)  
- **Speed**: Speed value as string
- **Heading**: Direction/heading as string
- **EventTime**: Timestamp as string

## Coordinate Transformation

### SkyWave Coordinate Format

SkyWave provides coordinates as strings that represent degrees and minutes in a compressed format. The conversion algorithm is complex:

#### For 7-character coordinates:
- Characters 1-4: Degrees portion
- Characters 5-7: Minutes portion

#### For 8-character coordinates:
- Characters 1-5: Degrees portion  
- Characters 6-8: Minutes portion

### Conversion Algorithm

```go
// Extract degrees and decimal portions
if len(latitude) == 7 {
    latdegrees = latitude[:4]    // First 4 chars
    latdecimal = latitude[4:]    // Remaining 3 chars
} else {
    latdegrees = latitude[:5]    // First 5 chars  
    latdecimal = latitude[5:]    // Remaining chars
}

// Convert degrees portion
latdegreesfloat := parseFloat(latdegrees) / 60
latdetwodec := round(latdegreesfloat, 2)

// Convert decimal portion
latdecimalres := parseFloat(latdecimal) / 60000

// Apply sign if negative
if latdegrees contains "-" {
    latdecimalres *= -1
}

// Final coordinate
finalLat := latdetwodec + latdecimalres - 0.003333
```

**Note**: There's a hardcoded offset of `-0.003333` applied to latitude for calibration purposes.

## Field Mapping

### Device Identification
- **SkyWave MobileID** → **MVT366 IMEI**
  - Example: `02092247SKY6A70` → `02092247SKY6A70`

### Position Data
- **SkyWave Latitude** → **MVT366 Latitude** (after complex conversion)
  - Example: `1952100` → `19.521003`
- **SkyWave Longitude** → **MVT366 Longitude** (after complex conversion)  
  - Example: `-9921171` → `-99.211715`

### Motion Data
- **SkyWave Speed** → **MVT366 Speed** (direct conversion)
- **SkyWave Heading** → **MVT366 Direction** (direct conversion)

### Timestamps  
- **SkyWave ReceiveUTC** → **MVT366 Datetime**
  - Converted from `2025-09-03 18:02:25` to `250903180225` format

### Fixed Values
- **CommandType**: Always `"AAA"`
- **EventCode**: Always `35`
- **Altitude**: Fixed at `21.232345`
- **PositionStatus**: Always `true` (A in MVT366)
- **ProtocolVersion**: Always `3`

## MVT366 Message Format

The final MVT366 message format:
```
$$H[length],[IMEI],[command],[eventcode],[latitude],[longitude],[datetime],[status],[satellites],[signal],[speed],[direction],[hdop],[altitude],[mileage],[runtime],[basestation],[io],[analog],[geofence],[custom],[version],[fuel],[temp],[accel],[decel]*[checksum]
```

Example output:
```
$$H166,02092247SKY6A70,AAA,35,19.521003,-99.211715,250903180225,A,0,0,0.000000,361,0.000000,21.232345,0.000000,0,0030|0030|0030|0030|0030,,,3,,,0,0*A5
```

## Network Communication

### UDP Transmission
- **Server**: `13.89.38.9:1805`
- **Protocol**: UDP
- **Connection**: New connection created for each message
- **Data**: Raw MVT366 message string as bytes

### Connection Handling
```go
conn, err := net.Dial("udp", "13.89.38.9:1805")
// Send message
conn.Write([]byte(mvt366Message))  
conn.Close()
```

## Main Program Flow

### 1. Initialization
```go
func main() {
    fromid := os.Getenv("FROMIDSKYWAVE")
    doc := skywave.SkywaveDoc{
        From_id: fromiduint, 
        Access_id: 70001184, 
        Password: "JEUTPKKH"
    }
    // Start continuous processing
}
```

### 2. Continuous Processing Loop
```go
func ReadSince(doc skywave.SkywaveDoc) skywave.SkywaveDoc {
    for {
        // 1. Get XML data from API
        xmlData, err := doc.GetDoc()
        
        // 2. Parse XML response  
        sky := skywave.GetReturnMessagesResult{}
        sky.ParseXML(xmlData)
        
        // 3. Extract position messages
        messages, err := sky.ReturnedMessagesBridge()
        
        // 4. Process each message
        for _, message := range messages {
            // Convert to MVT366
            mvt366Data, err := skywave.FromBridgePayload(message)
            mvt366Message, err := mvt366Data.ToMVT366Message()
            
            // Send via UDP
            conn, err := net.Dial("udp", "13.89.38.9:1805")
            conn.Write([]byte(mvt366Message))
            conn.Close()
        }
        
        // 5. Handle pagination
        if sky.More {
            doc.From_id = sky.NextStartID
        } else {
            return doc  // No more messages
        }
    }
}
```

### 3. Timing and Intervals
- **Continuous processing**: No delays between message batches within a document
- **Document completion**: 3-minute sleep between document cycles
- **Pagination**: Immediate processing of next batch when `More=true`

## Troubleshooting

### Common Issues

1. **Missing Environment Variable**:
   ```
   osvariable FROMIDSKYWAVE doesn't exists or empty
   ```
   **Solution**: Set the environment variable before running

2. **API Connection Issues**:
   ```
   NO HAY CONEXION CON SERVER1
   ```
   **Solution**: Check network connectivity to SkyWave API and UDP server

3. **Coordinate Conversion Errors**:
   ```
   no soy len de 7 [length1] [length2]
   ```
   **Solution**: Verify coordinate format from SkyWave API

4. **UDP Transmission Issues**:
   Check connectivity to `13.89.38.9:1805`

### Debug Information

The system outputs debug information:
- `ReadSince` - When starting message retrieval
- `Next doc [from_id]` - When moving to next message batch  
- `End document [from_id]` - When no more messages available

### Testing

Basic functionality test:
```bash
export FROMIDSKYWAVE=13969586728
go run main.go
```

## API Reference

### SkyWave API Parameters
- **access_id**: Account identifier (`70001184`)
- **password**: API password (`JEUTPKKH`)  
- **from_id**: Starting message ID (dynamic)

### Key Data Structures

- **SkywaveDoc**: API client configuration
- **GetReturnMessagesResult**: XML response parser
- **PayloadBridge**: Simplified message structure
- **MVT366**: Final protocol structure

---

**Created**: September 4, 2025  
**Last Updated**: September 4, 2025  
**Analysis Based On**: main.go, doc.go, skywave.go, mvt366.go
