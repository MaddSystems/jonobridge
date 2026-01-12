# ORBCOMM Test Configuration

This directory contains tools for testing ST9101 device support with the skywave protocol.

## Files

- `orbcomm_mqtt_test.py` - Python script to retrieve data from ORBCOMM and publish to MQTT
- `setup_orbcomm_test.sh` - Setup script for dependencies
- `test_orbcomm.md` - This documentation

## Device Information

- **Model**: ST9101 (ORBCOMM)
- **SAT ID**: 02092234SKYB62F
- **IMEI**: 353500723152270
- **Account**: rodolfo@gpscontrol.com.mx

## Quick Start

1. **Setup environment:**
   ```bash
   # Install dependencies
   pip3 install -r requirements.txt
   
   # Or use the setup script
   chmod +x setup_orbcomm_test.sh
   ./setup_orbcomm_test.sh
   ```

2. **Test ORBCOMM data retrieval (one-time):**
   ```bash
   python3 orbcomm_mqtt_test.py --once --verbose
   ```

3. **Test with sample data only (no API calls):**
   ```bash
   python3 orbcomm_mqtt_test.py --once --verbose
   ```

3. **Run skywave protocol processor:**
   ```bash
   # In another terminal
   ./skywaveprotocol -v
   ```

4. **Monitor MQTT messages:**
   ```bash
   # In another terminal
   mosquitto_sub -h localhost -t 'tracker/from-tcp' -v
   mosquitto_sub -h localhost -t 'tracker/jonoprotocol' -v
   ```

## Testing Process

1. **Data Flow:**
   ```
   ORBCOMM API → Python Script → MQTT (tracker/from-tcp) → Skywave Protocol → MQTT (tracker/jonoprotocol)
   ```

2. **Expected Data Format:**
   - ORBCOMM returns XML in `GetReturnMessagesResult` format
   - Python script formats it as JSON for MQTT
   - Skywave protocol parses XML and converts to Jono format

3. **API Testing:**
   - The script relies on successful API connection to retrieve data
   - No fallback data is available if the API is inaccessible

4. **Troubleshooting:**
   - Use `--verbose` flag to see detailed logs
   - Check `tracker/from-tcp/debug` topic for raw XML data
   - Verify ORBCOMM API credentials and endpoints

## API Endpoints Tried

The script attempts multiple ORBCOMM API endpoints:
- `/GlobalMessages/GetReturnMessages`
- `/api/messages/return`
- `/services/GetReturnMessages`

## Command Line Options

```bash
python3 orbcomm_mqtt_test.py [OPTIONS]

Options:
  --mqtt-host HOST     MQTT broker host (default: localhost)
  --mqtt-port PORT     MQTT broker port (default: 1883)
  --topic TOPIC        MQTT topic (default: tracker/from-tcp)
  --interval SECONDS   Polling interval (default: 60)
  --hours-back HOURS   Historical data hours (default: 24)
  --once               Run once and exit
  --verbose, -v        Enable verbose logging
```

## Expected Output

If ST9101 is supported, you should see:
1. XML data retrieved from ORBCOMM
2. Data published to MQTT
3. Skywave protocol parsing the data
4. Jono protocol output with device location/status

If not supported, you'll see parsing errors in the skywave protocol logs.
