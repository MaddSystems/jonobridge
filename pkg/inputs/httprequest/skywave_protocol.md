# SkyWave XML Protocol Documentation

## Overview

This document describes the XML structure returned by the SkyWave satellite API endpoint for retrieving satellite messages. The API returns position and telemetry data from satellite terminals in a structured XML format.

## API Endpoint

```
https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml
```

## Authentication Parameters

- `access_id`: Account identifier (e.g., `70001184`)
- `password`: API password (e.g., `JEUTPKKH`)
- `from_id`: Starting message ID for pagination (e.g., `13969586728`)

## XML Response Structure

### Root Element: `GetReturnMessagesResult`

The root element contains metadata about the API response and the collection of messages.

```xml
<GetReturnMessagesResult>
    <ErrorID>0</ErrorID>
    <More>true</More>
    <NextStartID>20368122913</NextStartID>
    <Messages>
        <!-- Array of ReturnMessage elements -->
    </Messages>
</GetReturnMessagesResult>
```

#### Root Element Fields

| Field | Type | Description |
|-------|------|-------------|
| `ErrorID` | Integer | Error code (0 = success) |
| `More` | Boolean | Indicates if more messages are available for pagination |
| `NextStartID` | Long | Next message ID for subsequent API calls |
| `Messages` | Array | Collection of satellite messages |

### Message Element: `ReturnMessage`

Each individual satellite message is contained within a `ReturnMessage` element.

```xml
<ReturnMessage>
    <ID>20368122913</ID>
    <MessageUTC>2025-09-05 02:19:48</MessageUTC>
    <ReceiveUTC>2025-09-05 02:19:47</ReceiveUTC>
    <SIN>126</SIN>
    <MobileID>02092247SKY6A70</MobileID>
    <Payload Name="StationaryIntervalCell" SIN="126" MIN="48">
        <Fields>
            <!-- Field elements -->
        </Fields>
    </Payload>
    <RegionName>CELLMTBP</RegionName>
    <OTAMessageSize>17</OTAMessageSize>
</ReturnMessage>
```

#### ReturnMessage Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `ID` | Long | Unique message identifier | `20368122913` |
| `MessageUTC` | DateTime | When message was sent by satellite | `2025-09-05 02:19:48` |
| `ReceiveUTC` | DateTime | When message was received by ground station | `2025-09-05 02:19:47` |
| `SIN` | Integer | Service Identification Number | `126` |
| `MobileID` | String | Satellite terminal identifier | `02092247SKY6A70` |
| `Payload` | Complex | Message payload with position data | See below |
| `RegionName` | String | Geographic region identifier | `CELLMTBP` |
| `OTAMessageSize` | Integer | Over-the-air message size in bytes | `17` |

### Payload Element

The `Payload` element contains the actual satellite message data with attributes and fields.

```xml
<Payload Name="StationaryIntervalCell" SIN="126" MIN="48">
    <Fields>
        <Field Name="Latitude" Value="1171261"/>
        <Field Name="Longitude" Value="-5952697"/>
        <Field Name="Speed" Value="0"/>
        <Field Name="Heading" Value="361"/>
        <Field Name="EventTime" Value="1757038722"/>
        <Field Name="GpsFixAge" Value="1023"/>
    </Fields>
</Payload>
```

#### Payload Attributes

| Attribute | Type | Description | Example |
|-----------|------|-------------|---------|
| `Name` | String | Payload type identifier | `StationaryIntervalCell` |
| `SIN` | Integer | Service Identification Number | `126` |
| `MIN` | Integer | Message Identification Number | `48` |

### Field Elements

Individual data fields within the payload are represented as `Field` elements.

```xml
<Field Name="Latitude" Value="1171261"/>
<Field Name="Longitude" Value="-5952697"/>
<Field Name="Speed" Value="0"/>
<Field Name="Heading" Value="361"/>
<Field Name="EventTime" Value="1757038722"/>
<Field Name="GpsFixAge" Value="1023"/>
```

#### Common Field Types

| Field Name | Type | Description | Example | Notes |
|------------|------|-------------|---------|-------|
| `Latitude` | Integer | Latitude coordinate (SkyWave format) | `1171261` | Requires conversion to decimal degrees |
| `Longitude` | Integer | Longitude coordinate (SkyWave format) | `-5952697` | Requires conversion to decimal degrees |
| `Speed` | Integer | Speed in knots | `0` | 0 = stationary |
| `Heading` | Integer | Direction in degrees | `361` | 361 = invalid/unknown |
| `EventTime` | Long | Unix timestamp | `1757038722` | Seconds since epoch |
| `GpsFixAge` | Integer | Age of GPS fix in seconds | `1023` | 1023 = invalid/old |

## Payload Types

### Position-Related Payload Types

| Payload Name | Description | SIN | MIN |
|--------------|-------------|-----|-----|
| `StationaryIntervalCell` | Stationary position report (cellular) | 126 | 48 |
| `StationaryIntervalSat` | Stationary position report (satellite) | 126 | 1 |
| `MovingIntervalCell` | Moving position report (cellular) | 126 | 49 |
| `MovingIntervalSat` | Moving position report (satellite) | 126 | 2 |
| `DistanceCell` | Distance-based report (cellular) | 126 | 50 |
| `MovingEnd` | End of movement report | 126 | 4 |
| `MovingStart` | Start of movement report | 126 | 3 |
| `IgnitionOn` | Ignition turned on | 126 | 5 |

### Other Payload Types

| Payload Name | Description | SIN | MIN |
|--------------|-------------|-----|-----|
| `PowerBackup` | Power backup event | 126 | 3 |

## Coordinate System Conversion

### SkyWave Coordinate Format

SkyWave provides coordinates as integers that need conversion to decimal degrees:

#### Latitude Conversion
- **7-digit format**: `1171261` → `11.71261°`
- **8-digit format**: `11712610` → `11.712610°`

#### Longitude Conversion
- **7-digit format**: `-5952697` → `-59.52697°`
- **8-digit format**: `-59526970` → `-59.526970°`

### Conversion Algorithm

```python
def convert_skywave_coordinate(coord_str):
    coord = int(coord_str)
    degrees = coord // 100000
    minutes = (coord % 100000) / 100000 * 60
    decimal_degrees = degrees + minutes
    return decimal_degrees
```

## Parsing Guidelines

### 1. Root Element Validation
- Always check `ErrorID` for API errors
- Use `More` and `NextStartID` for pagination
- Expect zero or more `ReturnMessage` elements

### 2. Message Processing
- Parse `MessageUTC` and `ReceiveUTC` as timestamps
- Extract `MobileID` for device identification
- Process payload based on `Name` attribute

### 3. Payload Field Extraction
- Iterate through all `Field` elements
- Convert coordinate values using the algorithm above
- Handle missing or invalid field values gracefully

### 4. Error Handling
- Check for missing required fields
- Validate coordinate ranges
- Handle timestamp parsing errors

## Example XML Response

```xml
<?xml version="1.0" encoding="utf-8"?>
<GetReturnMessagesResult xmlns="http://www.skywave.com">
    <ErrorID>0</ErrorID>
    <More>true</More>
    <NextStartID>20368122913</NextStartID>
    <Messages>
        <ReturnMessage>
            <ID>20368122913</ID>
            <MessageUTC>2025-09-05 02:19:48</MessageUTC>
            <ReceiveUTC>2025-09-05 02:19:47</ReceiveUTC>
            <SIN>126</SIN>
            <MobileID>02092247SKY6A70</MobileID>
            <Payload Name="StationaryIntervalCell" SIN="126" MIN="48">
                <Fields>
                    <Field Name="Latitude" Value="1171261"/>
                    <Field Name="Longitude" Value="-5952697"/>
                    <Field Name="Speed" Value="0"/>
                    <Field Name="Heading" Value="361"/>
                    <Field Name="EventTime" Value="1757038722"/>
                    <Field Name="GpsFixAge" Value="1023"/>
                </Fields>
            </Payload>
            <RegionName>CELLMTBP</RegionName>
            <OTAMessageSize>17</OTAMessageSize>
        </ReturnMessage>
    </Messages>
</GetReturnMessagesResult>
```

## Python Parsing Example

```python
import xml.etree.ElementTree as ET
from datetime import datetime

def parse_skywave_xml(xml_content):
    root = ET.fromstring(xml_content)

    # Check for errors
    error_id = int(root.find('ErrorID').text)
    if error_id != 0:
        raise ValueError(f"API Error: {error_id}")

    messages = []
    for message in root.findall('.//ReturnMessage'):
        msg_data = {
            'id': message.find('ID').text,
            'message_utc': message.find('MessageUTC').text,
            'receive_utc': message.find('ReceiveUTC').text,
            'mobile_id': message.find('MobileID').text,
            'sin': int(message.find('SIN').text),
            'region': message.find('RegionName').text,
            'ota_size': int(message.find('OTAMessageSize').text)
        }

        # Parse payload
        payload = message.find('Payload')
        if payload is not None:
            msg_data['payload_name'] = payload.get('Name')
            msg_data['payload_sin'] = payload.get('SIN')
            msg_data['payload_min'] = payload.get('MIN')

            # Parse fields
            fields = {}
            for field in payload.findall('.//Field'):
                name = field.get('Name')
                value = field.get('Value')
                fields[name] = value

            msg_data['fields'] = fields

        messages.append(msg_data)

    return {
        'error_id': error_id,
        'more': root.find('More').text.lower() == 'true',
        'next_start_id': root.find('NextStartID').text,
        'messages': messages
    }
```

## Go Parsing Example

```go
package main

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
)

// XML Structures
type GetReturnMessagesResult struct {
    XMLName      xml.Name       `xml:"GetReturnMessagesResult"`
    ErrorID      int            `xml:"ErrorID"`
    More         bool           `xml:"More"`
    NextStartID  string         `xml:"NextStartID"`
    Messages     []ReturnMessage `xml:"Messages>ReturnMessage"`
}

type ReturnMessage struct {
    ID            string  `xml:"ID"`
    MessageUTC    string  `xml:"MessageUTC"`
    ReceiveUTC    string  `xml:"ReceiveUTC"`
    SIN           int     `xml:"SIN"`
    MobileID      string  `xml:"MobileID"`
    Payload       Payload `xml:"Payload"`
    RegionName    string  `xml:"RegionName"`
    OTAMessageSize int    `xml:"OTAMessageSize"`
}

type Payload struct {
    Name   string `xml:"Name,attr"`
    SIN    int    `xml:"SIN,attr"`
    MIN    int    `xml:"MIN,attr"`
    Fields []Field `xml:"Fields>Field"`
}

type Field struct {
    Name  string `xml:"Name,attr"`
    Value string `xml:"Value,attr"`
}

func main() {
    // Make API request
    resp, err := http.Get("https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml/?access_id=70001184&password=JEUTPKKH&from_id=13969586728")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    // Parse XML
    var result GetReturnMessagesResult
    err = xml.Unmarshal(body, &result)
    if err != nil {
        panic(err)
    }

    // Process results
    fmt.Printf("Error ID: %d\n", result.ErrorID)
    fmt.Printf("More messages: %t\n", result.More)
    fmt.Printf("Next Start ID: %s\n", result.NextStartID)
    fmt.Printf("Message count: %d\n", len(result.Messages))

    // Process each message
    for _, msg := range result.Messages {
        fmt.Printf("Message ID: %s, Mobile ID: %s\n", msg.ID, msg.MobileID)
        // Process payload fields...
    }
}
```

## Notes

- All timestamps are in UTC
- Coordinate values require conversion from SkyWave format to decimal degrees
- The API supports pagination using the `from_id` parameter
- `GpsFixAge` value of 1023 indicates invalid or very old GPS data
- `Heading` value of 361 indicates invalid or unknown direction
- The XML namespace may be present: `xmlns="http://www.skywave.com"`

---

**Document Version**: 1.0
**Last Updated**: September 5, 2025
**Based on SkyWave API Response**: September 5, 2025</content>
<parameter name="filePath">/home/ubuntu/jonobridge/pkg/inputs/httprequest/skywave_protocol.md
