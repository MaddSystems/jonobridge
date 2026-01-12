# Pagination Implementation Summary

## ✅ Changes Applied to main.go

### 1. Added XML Parsing Support
```go
import (
    "encoding/xml"  // NEW: Added for parsing pagination fields
    // ... other imports
)
```

### 2. Created SkyWave Response Structure
```go
// SkyWave XML response structure for pagination
type SkyWaveResponse struct {
    XMLName     xml.Name `xml:"GetReturnMessagesResult"`
    ErrorID     string   `xml:"ErrorID"`
    More        string   `xml:"More"`
    NextStartID string   `xml:"NextStartID"`
}
```

### 3. Rewrote processSkyWaveData() with Pagination Loop

**BEFORE (Single Request):**
```go
func processSkyWaveData(...) error {
    // Build URL
    // Make ONE request
    // Publish ONE response to MQTT
    // DONE - stops here even if More=true
}
```

**AFTER (Paginated Requests):**
```go
func processSkyWaveData(...) error {
    var nextStartID string
    pageNum := 0
    
    for {  // LOOP until More=false
        pageNum++
        
        // Build URL with start_id for pagination (if not first page)
        if nextStartID != "" {
            // Include &start_id=XXX for subsequent pages
        }
        
        // Make HTTP request
        // Parse XML to extract <More> and <NextStartID>
        // Publish this page to MQTT
        
        if More == "true" && NextStartID != "" {
            nextStartID = NextStartID
            continue  // Fetch next page
        } else {
            break  // No more pages, exit loop
        }
    }
}
```

## 🎯 Key Features Implemented

### ✅ Full Pagination Support
- Checks `<More>` field in XML response
- Extracts `<NextStartID>` for next page
- Loops until `More=false`
- Safety limit: max 100 pages to prevent infinite loops

### ✅ Dual Mode Compatibility
- **SkyWave Mode** (SKYWAVE_* env vars set): Uses pagination
- **HTTP Mode** (no SKYWAVE_* vars): Falls back to original processHttpData() - NO CHANGES

### ✅ Enhanced Logging
```
Page 1: Hex string length: 15234 characters
Page 1: Successfully published to MQTT topic 'skywave/xml'
More pages available, NextStartID: 20574290653
Page 2: Hex string length: 8456 characters
Page 2: Successfully published to MQTT topic 'skywave/xml'
No more pages. Total pages retrieved: 2
Pagination complete: 2 pages published to MQTT
```

### ✅ URL Parameter Handling
- Preserves `from_id` as mobile filter
- Adds `start_id` parameter for pagination
- Works with both HTTP_URL env var and hardcoded URL

## 📊 Expected Behavior Changes

### Before (Missing Events):
```
Request 1: 17:00-17:03
  Response: Page 1 (100 events) - Published ✅
           More=true, NextStartID=123456
           Pages 2+ (50 events) - NEVER FETCHED ❌
  
Result: Only 100 events sent to MQTT
```

### After (All Events):
```
Request 1: 17:00-17:03
  Response Page 1 (100 events) - Published ✅
           More=true, NextStartID=123456
  
Request 2: with start_id=123456
  Response Page 2 (50 events) - Published ✅
           More=false
  
Result: All 150 events sent to MQTT ✅
```

## 🧪 Testing Instructions

### 1. Build the Updated Binary
```bash
cd /home/ubuntu/jonobridge/pkg/inputs/httprequest
go build -o httprequest .
```

### 2. Run with Verbose Logging
```bash
export MQTT_BROKER_HOST="your-mqtt-host"
export SKYWAVE_ACCESS_ID="70001184"
export SKYWAVE_PASSWORD="JEUTPKKH"
export SKYWAVE_FROM_ID="13969586728"
export HTTP_URL="https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml"
export HTTP_POLLING_TIME="180"

./httprequest -v
```

### 3. Watch for Pagination Logs
Look for:
```
Using polling interval of 3m0s — UTC: ...
Page 1: Successfully published to MQTT topic 'skywave/xml'
More pages available, NextStartID: XXXXX
Fetching page 2 with start_id=XXXXX
Page 2: Successfully published to MQTT topic 'skywave/xml'
No more pages. Total pages retrieved: 2
Pagination complete: 2 pages published to MQTT
```

### 4. Verify All Event Types Are Now Appearing
Monitor your MQTT consumer/database and check for:
- ✅ StationaryIntervalCell (MIN=48)
- ✅ DigInp2Lo (MIN=53)
- ✅ DigInp2Hi (MIN=52)
- ✅ IdlingStart (MIN=21)
- ✅ IgnitionOn, MovingStart, etc.

## 🔒 Safety Features

### 1. Pagination Safety Limit
- Maximum 100 pages per polling cycle
- Prevents infinite loops if API misbehaves

### 2. Error Handling Per Page
- If page N fails, error is logged with page number
- Helps identify which page caused issues

### 3. Backward Compatibility
- HTTP mode (non-SkyWave) unchanged
- Existing deployments without SKYWAVE_* vars work as before

### 4. Graceful XML Parse Failures
- If XML parsing fails, page is still published
- Logs warning but continues operation

## 📈 Performance Impact

### Network:
- More HTTP requests if pagination occurs
- Each page is a separate GET request
- Typically 1-3 pages per polling cycle

### MQTT:
- More MQTT publishes (one per page)
- Each page published separately
- Downstream consumers receive multiple messages per cycle

### Memory:
- Minimal increase (one page at a time)
- No buffering of all pages in memory

## 🚀 Next Steps (Optional Enhancements)

1. **Add Event Counting**
   - Parse `<ReturnMessage>` count per page
   - Log: "Page 1: 100 events, Page 2: 50 events"

2. **Increase Time Window**
   - Consider setting `HTTP_POLLING_TIME=300` (5 minutes)
   - Provides safety overlap for delayed events

3. **Add Payload Type Stats**
   - Count events by Payload Name
   - Log distribution: "StationaryIntervalCell: 80, IdlingStart: 20"

4. **Merge Pages Before Publishing**
   - Combine all pages into single XML
   - Publish once instead of per-page
   - Reduces MQTT message count

## ⚠️ Important Notes

- **This change ONLY affects SkyWave mode** (when SKYWAVE_* env vars are set)
- **HTTP mode remains unchanged** (processHttpData() untouched)
- **Each page is published separately** to MQTT as hex-encoded XML
- **Downstream consumers** should expect multiple MQTT messages per polling cycle
- **Build succeeded** - binary is ready for deployment

## 📝 Files Modified

- ✅ `/home/ubuntu/jonobridge/pkg/inputs/httprequest/main.go` - Added pagination
- ✅ Build verified - compiles successfully
- ✅ Backward compatible - dual mode preserved
