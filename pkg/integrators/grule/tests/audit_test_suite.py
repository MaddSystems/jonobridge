import time
import requests
import json
import sys
import paho.mqtt.client as mqtt
from datetime import datetime, timedelta

# Configuration
MQTT_BROKER = "localhost"
MQTT_PORT = 1883
MQTT_TOPIC = "tracker/jonoprotocol"
API_URL = "http://localhost:8081/api/audit/progress"
IMEI = "TEST_AUDIT_001"

def crc(source):
    b = 0
    for i in source:
        b = b + ord(i)
    ret = hex(b % 256)
    ret = ret.upper()
    ret = ret.replace("0X", "")
    return ret

def charcounter(source):
    return len(source)

class identifier:
    idCounter = 64
    def __init__(self):
        if identifier.idCounter < 123:
            identifier.idCounter += 1
        else:
            identifier.idCounter = 65

def payload(imei, eventCode, latitude, longitude, utc, status, sats, gsmStrenght, speed, direction, accuracy, altitude, mileage, runtime, mcc, mnc, lac, cellId, portStatus, AD1, AD2, AD3, battery, AD5, eventInfo):
    imei = imei.strip()
    newIdentifier = identifier()
    mydataidentifier = str(chr(newIdentifier.idCounter))

    first_output = "," + imei + ",AAA," + eventCode + "," + latitude + "," + longitude + "," + utc + "," + status + "," + str(sats) + "," + str(gsmStrenght) + "," + str(speed) + "," + str(direction) + "," + str(accuracy) + "," + str(altitude) + "," + str(mileage) + "," + str(runtime) + "," + mcc + "|" + mnc + "|" + lac + "|" + cellId + "," + portStatus + "," + AD1 + "|" + AD2 + "|" + AD3 + "|" + str(battery) + "|" + AD5 + "," + eventInfo + ",*"
    totalchar = charcounter(first_output) + 4
    header = "$$" + mydataidentifier + str(totalchar)
    preoutput = header + first_output
    output = preoutput + crc(preoutput) + chr(13) + chr(10)
    return output

def send_packet(packet_data):
    try:
        client = mqtt.Client(callback_api_version=mqtt.CallbackAPIVersion.VERSION2)
        client.connect(MQTT_BROKER, MQTT_PORT, 60)
        
        # Wrap in expected JSON format if backend expects JonoModel or raw string?
        # backend/grule/worker.go uses adapter.Parse(payload).
        # backend/adapters/gps_tracker.go likely expects the raw string or JSON.
        # Looking at backend/main.go: engine.ProcessPacketMessage(string(msg.Payload()))
        # engine/grule_worker.go: json.Unmarshal([]byte(payload), &jono)
        
        # So we need to wrap it in a JSON structure compatible with models.JonoModel
        jono_payload = {
            "IMEI": IMEI,
            "ListPackets": {
                "0": {
                    "IMEI": IMEI,
                    "Speed": int(packet_data.split(",")[9]), # Extract speed roughly
                    "Latitude": float(packet_data.split(",")[4]),
                    "Longitude": float(packet_data.split(",")[5]),
                    "Datetime": datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%SZ'), # ISO 8601
                    "PositioningStatus": packet_data.split(",")[7],
                    "GSMSignalStrength": int(packet_data.split(",")[8]),
                    "Altitude": 0,
                    "EventCode": {"Code": 35, "Name": "Location Update"}, # Dummy event code
                    "Direction": 0,
                    "Mileage": 0,
                    "NumberOfSatellites": int(packet_data.split(",")[8]) # Sats position
                }
            }
        }
        
        # Actually, let's use the raw string if the adapter supports it, but the backend worker 
        # explicitly Unmarshals JSON.
        # "backend/grule/worker.go": packets, err := w.adapter.Parse(payload)
        # If adapter is GPSTrackerAdapter, let's check it.
        
        # Assuming backend expects JSON wrapping the packet data or pre-parsed data.
        # Let's try sending the JSON format that engine/grule_worker.go expects.
        
        client.publish(MQTT_TOPIC, json.dumps(jono_payload))
        print(f"📤 Sent MQTT packet for {IMEI}")
        client.disconnect()
        return True
    except Exception as e:
        print(f"❌ Error sending MQTT packet: {e}")
        return False

def clear_audit():
    try:
        requests.post("http://localhost:8081/api/audit/progress/clear")
        requests.post("http://localhost:8081/api/audit/progress/enable")
        print("🧹 Audit log cleared and enabled")
    except Exception as e:
        print(f"⚠️ Failed to clear/enable audit: {e}")

def check_audit_log(step_description):
    print(f"🔍 Verifying audit for: {step_description}...")
    time.sleep(2) # Wait for processing
    
    try:
        resp = requests.get(f"{API_URL}/timeline", params={"imei": IMEI, "limit": 1})
        data = resp.json()
        
        if not data.get("success"):
            print("❌ API Error")
            return False
            
        frames = data.get("frames", [])
        if not frames:
            print("❌ No frames found")
            return False
            
        latest = frames[0]
        snapshot = latest.get("snapshot", {})
        
        print(f"   ✅ Stage: {latest.get('stage_reached')}")
        print(f"   ✅ Step: {latest.get('step_number')}")
        
        # Check rich data
        if "buffer_circular" in snapshot:
            print(f"   ✅ Buffer Size: {len(snapshot['buffer_circular'])}")
        else:
            print("   ❌ Missing buffer_circular in snapshot")
            
        if "jammer_metrics" in snapshot:
             metrics = snapshot['jammer_metrics']
             print(f"   ✅ Metrics: SpeedAvg={metrics.get('avg_speed_90min')}, GSMAvg={metrics.get('avg_gsm_last5')}")
        else:
             print("   ❌ Missing jammer_metrics in snapshot")

        return True
    except Exception as e:
        print(f"❌ Verification failed: {e}")
        return False

def run_test():
    print(f"🚀 Starting Declarative Audit Integration Test for IMEI: {IMEI}")
    clear_audit()
    
    # 1. Send Valid Packet (Normal Operation)
    print("\n--- Step 1: Sending Valid Packet ---")
    utc = datetime.utcnow().strftime('%y%m%d%H%M%S')
    # Using dummy payload function just to extract values for JSON
    pkt = payload(IMEI, "35", "19.4326", "-99.1332", utc, "A", "10", "25", "60", "0", "0", "0", "0", "0", "334", "020", "1", "1", "0000", "0", "0", "0", "100", "0", "0")
    send_packet(pkt)
    check_audit_log("Valid Packet Processing")
    
    # 2. Send Invalid Packet (Potential Jammer)
    print("\n--- Step 2: Sending Invalid Packet (Jammer Sim) ---")
    utc = datetime.utcnow().strftime('%y%m%d%H%M%S')
    pkt = payload(IMEI, "35", "19.4326", "-99.1332", utc, "V", "0", "5", "0", "0", "0", "0", "0", "0", "334", "020", "1", "1", "0000", "0", "0", "0", "100", "0", "0")
    send_packet(pkt)
    check_audit_log("Invalid Packet Processing")

if __name__ == "__main__":
    run_test()
