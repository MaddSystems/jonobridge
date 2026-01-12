import time
import requests
import json
import paho.mqtt.client as mqtt
from datetime import datetime

# Configuration
MQTT_BROKER = "localhost"
MQTT_PORT = 1883
MQTT_TOPIC = "tracker/jonoprotocol"
API_BASE_URL = "http://localhost:8081/api"
IMEI = "TEST_AUDIT_999"

def create_test_rule():
    print("📝 Creating test rule and manifest...")
    rule_data = {
        "name": "AuditIntegrationTest",
        "grl": 'rule AuditIntegrationTest { when IncomingPacket.Speed >= 0 then actions.Log("Audit test fired"); Retract("AuditIntegrationTest"); }',
        "audit_manifest": """
stages:
  - rule: AuditIntegrationTest
    order: 1
    audit:
      enabled: true
      description: "Integration Test Step"
      level: info
      is_alert: false
""",
        "active": True,
        "priority": 100
    }
    
    try:
        # First try to delete if exists (optional, depends on API)
        resp = requests.post(f"{API_BASE_URL}/rules", json=rule_data)
        print(f"   Response: {resp.status_code} - {resp.text}")
        
        # Force reload rules
        resp = requests.post(f"{API_BASE_URL}/reload")
        print("   Rules reloaded")
        return True
    except Exception as e:
        print(f"❌ Failed to create rule: {e}")
        return False

def send_packet(status="A", speed=60):
    try:
        client = mqtt.Client(callback_api_version=mqtt.CallbackAPIVersion.VERSION2)
        client.connect(MQTT_BROKER, MQTT_PORT, 60)
        
        jono_payload = {
            "IMEI": IMEI,
            "ListPackets": {
                "0": {
                    "IMEI": IMEI,
                    "Speed": int(speed / 3.6), # FIX: Must be int
                    "Latitude": 19.4326,
                    "Longitude": -99.1332,
                    "Datetime": datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%SZ'),
                    "PositioningStatus": status,
                    "GSMSignalStrength": 25,
                    "Altitude": 100,
                    "EventCode": {"Code": 35, "Name": "TestEvent"}, # FIX: Code is int
                    "Direction": 0,
                    "Mileage": 0,
                    "NumberOfSatellites": 12
                }
            }
        }
        
        client.publish(MQTT_TOPIC, json.dumps(jono_payload))
        print(f"📤 Sent MQTT packet for {IMEI} (Status={status}, Speed={speed})")
        client.disconnect()
        return True
    except Exception as e:
        print(f"❌ Error sending MQTT packet: {e}")
        return False

def check_audit_log():
    print(f"🔍 Verifying audit for IMEI {IMEI}...")
    time.sleep(3) # Wait for processing
    
    try:
        resp = requests.get(f"{API_BASE_URL}/audit/progress/timeline", params={"imei": IMEI, "limit": 5})
        data = resp.json()
        
        if not data.get("success"):
            print(f"❌ API Error: {data}")
            return False
            
        frames = data.get("rows", [])
        if not frames:
            print("❌ No audit frames found in DB")
            return False
            
        print(f"✅ Found {len(frames)} audit frames")
        for i, frame in enumerate(frames):
            print(f"   Frame {i}: Rule={frame.get('rule_name')}, Stage={frame.get('stage_reached')}, Step={frame.get('step_number')}, Level={frame.get('level')}")
            
            # Verify new fields
            if 'level' not in frame:
                print(f"      ❌ Missing 'level' in frame {i}")
                return False

            snapshot = frame.get("snapshot", {})
            if "packet_current" in snapshot:
                print(f"      ✅ Snapshot has packet data: Speed={snapshot['packet_current'].get('Speed')}")
            if "buffer_circular" in snapshot:
                print(f"      ✅ Snapshot has buffer data ({len(snapshot['buffer_circular'])} entries)")
                
        return True
    except Exception as e:
        print(f"❌ Verification failed: {e}")
        return False

def test_dynamic_frontend():
    print("\n🧪 Testing Dynamic Frontend Data...")
    rule_name = "DynamicTestRule"
    rule_data = {
        "name": rule_name,
        "grl": f'rule {rule_name} {{ when IncomingPacket.GSMSignalStrength > 0 then actions.Log("Dynamic test fired"); Retract("{rule_name}"); }}',
        "audit_manifest": """
stages:
  - rule: DynamicTestRule
    order: 7
    audit:
      enabled: true
      description: "Custom Dynamic Stage"
      level: critical
      is_alert: true
""",
        "active": True,
        "priority": 100
    }
    
    try:
        # Create rule
        requests.post(f"{API_BASE_URL}/rules", json=rule_data)
        requests.post(f"{API_BASE_URL}/reload")
        
        # Send packet
        send_packet(speed=88)
        time.sleep(2)
        
        # Verify - get latest frames (reverse sort to get newest first)
        resp = requests.get(f"{API_BASE_URL}/audit/progress/timeline", params={"imei": IMEI, "limit": 10})
        data = resp.json()
        frames = data.get("rows", [])
        
        if not frames:
            print("   ❌ No frames for dynamic test")
            return False
        
        # Find DynamicTestRule frame (should be newest)
        dynamic_frame = None
        for frame in frames:
            if frame.get('rule_name') == 'DynamicTestRule':
                dynamic_frame = frame
                break
        
        if not dynamic_frame:
            print(f"   ❌ DynamicTestRule frame not found. Got {len(frames)} frames:")
            for f in frames:
                print(f"      - {f.get('rule_name')}: {f.get('stage_reached')}")
            return False
            
        frame = dynamic_frame
        print(f"   ✅ Frame: Stage='{frame.get('stage_reached')}', Level='{frame.get('level')}', Step={frame.get('step_number')}")
        
        passed = (frame.get('stage_reached') == "Custom Dynamic Stage" and 
                  frame.get('level') == "critical" and 
                  frame.get('step_number') == 7)
        
        if passed:
            print("   ✅ Dynamic data verification PASSED")
        else:
            print("   ❌ Dynamic data verification FAILED")
            
        return passed
    except Exception as e:
        print(f"   ❌ Dynamic test failed with error: {e}")
        return False

def run_test():
    print(f"🚀 Starting Phase 3 & 5 Integration Test")
    
    # 1. Setup Rule
    if not create_test_rule():
        return
        
    # 2. Enable Audit
    try:
        requests.post(f"{API_BASE_URL}/audit/progress/enable")
        print("✅ Audit enabled")
    except:
        print("⚠️ Could not enable audit via API")

    # 3. Send Packet
    send_packet(status="A", speed=75)
    
    # 4. Verify Phase 3
    p3_passed = check_audit_log()
    
    # 5. Verify Phase 5
    p5_passed = test_dynamic_frontend()
    
    if p3_passed and p5_passed:
        print("\n✨ ALL INTEGRATION TESTS PASSED! ✨")
    else:
        print("\n❌ SOME INTEGRATION TESTS FAILED ❌")

if __name__ == "__main__":
    run_test()
