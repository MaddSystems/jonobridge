import socket
from socket import SOL_SOCKET, SO_REUSEADDR
import sys
import os
from datetime import datetime, timedelta
import time
import random

# Add parent directory to import geofences module
sys.path.append(os.path.join(os.path.dirname(__file__), '..', 'geofences'))

try:
    from get_geofence_mysql import get_geofences_by_group
    GEOFENCE_CHECK_AVAILABLE = True
except ImportError:
    GEOFENCE_CHECK_AVAILABLE = False
    print("⚠️  Geofence validation unavailable (get_geofence_mysql not found)")

def crc(source):
    b = 0
    for i in source:
        b = b + ord(i)
    ret = hex(b % 256)
    ret = ret.upper()
    ret = ret.replace("0X", "")
    return ret

def charcounter(source):
    c = 0
    for i in source:
        c = c + 1
    return c

class identifier:
    idCounter = 64
    def __init__(self):
        if identifier.idCounter < 123:
            identifier.idCounter += 1
        else:
            identifier.idCounter = 65

def payload(imei, eventCode, latitude, longitude, utc, status, sats, gsmStrenght, speed, direction, accuracy, altitude, mileage, runtime, mcc, mnc, lac, cellId, portStatus, AD1, AD2, AD3, battery, AD5, eventInfo):
    # MVT380
    imei = imei.strip()
    newIdentifier = identifier()
    mydataidentifier = str(chr(newIdentifier.idCounter))

    first_output = "," + imei + ",AAA," + eventCode + "," + latitude + "," + longitude + "," + utc + "," + status + "," + str(sats) + "," + str(gsmStrenght) + "," + str(speed) + "," + str(direction) + "," + str(accuracy) + "," + str(altitude) + "," + str(mileage) + "," + str(runtime) + "," + mcc + "|" + mnc + "|" + lac + "|" + cellId + "," + portStatus + "," + AD1 + "|" + AD2 + "|" + AD3 + "|" + str(battery) + "|" + AD5 + "," + eventInfo + ",*"
    totalchar = charcounter(first_output) + 4
    header = "$$" + mydataidentifier + str(totalchar)
    preoutput = header + first_output
    output = preoutput + crc(preoutput) + chr(13) + chr(10)
    return output

# Configuration for 4 IMEIs with different targets
IMEI_CONFIG = [
    {
        "imei": "864352045580761",
        "name": "IMEI_1 (Defcon 4)",
        "target": 4,
        "gsm": 10,
        "lat": "19.600000",  # OUTSIDE all safe zones → Will trigger DEFCON 4
        "lon": "-99.300000",
        "stop_at_step": 24 # Full simulation
    },
    {
        "imei": "864352045580762",
        "name": "IMEI_2 (Defcon 4)",
        "target": 4,
        "gsm": 10,
        "lat": "19.350000",  # OUTSIDE all safe zones → Will trigger DEFCON 4
        "lon": "-99.150000",
        "stop_at_step": 24 # Full simulation
    },
    {
        "imei": "864352045580763",
        "name": "IMEI_3 (Pass Defcon 2, Fail Defcon 3 - Safe Zone)",
        "target": 2,
        "gsm": 10,
        "lat": "19.456200", # INSIDE TALLER SAN JUAN ARAGON → Passes DEFCON 2, stops at 3 (safe zone)
        "lon": "-99.092900",
        "stop_at_step": 24 # Full simulation
    },
    {
        "imei": "864352045580764",
        "name": "IMEI_4 (Fail Defcon 2 - Strong Signal)",
        "target": 1,
        "gsm": 25,  # STRONG signal → Will FAIL DEFCON 2 (GSM must be < 15)
        "lat": "19.520730",  # Location doesn't matter - fails before geofence check
        "lon": "-99.211520",
        "stop_at_step": 24 # Full simulation to show it never progresses past DEFCON 1
    }
]

SERVER = "jonobridge.madd.com.mx"
PORT = 8056
PACKET_INTERVAL = 15
WARMUP_PACKETS = 11
INVALID_PACKETS = 24

def is_inside_bounding_box(lat, lon, geofence):
    """Check if point is inside geofence bounding box (polygon approximation)"""
    if geofence['shapeType'].lower() == 'polygon':
        if (geofence['boundingBoxMinX'] is not None and 
            geofence['boundingBoxMaxX'] is not None and
            geofence['boundingBoxMinY'] is not None and
            geofence['boundingBoxMaxY'] is not None):
            
            return (lat >= geofence['boundingBoxMinY'] and 
                    lat <= geofence['boundingBoxMaxY'] and
                    lon >= geofence['boundingBoxMinX'] and 
                    lon <= geofence['boundingBoxMaxX'])
    return False

def is_inside_circle(lat, lon, geofence):
    """Check if point is inside circular geofence"""
    if geofence['shapeType'].lower() == 'circle':
        if (geofence['centerLat'] is not None and 
            geofence['centerLon'] is not None and
            geofence['radius'] is not None):
            
            # Haversine distance approximation
            import math
            earth_radius = 6371000.0  # meters
            dLat = (geofence['centerLat'] - lat) * math.pi / 180.0
            dLon = (geofence['centerLon'] - lon) * math.pi / 180.0
            a = (math.sin(dLat/2) * math.sin(dLat/2) +
                 math.cos(lat * math.pi / 180.0) * 
                 math.cos(geofence['centerLat'] * math.pi / 180.0) *
                 math.sin(dLon/2) * math.sin(dLon/2))
            c = 2 * math.atan2(math.sqrt(a), math.sqrt(1-a))
            distance = earth_radius * c
            return distance <= geofence['radius']
    return False

def check_if_inside_safe_zones(lat, lon):
    """Check if coordinate is inside any safe zone groups (CLIENTES, Taller, Resguardo)"""
    if not GEOFENCE_CHECK_AVAILABLE:
        return False, []
    
    groups = ['CLIENTES', 'Taller', 'Resguardo/Cedis/Puerto']
    inside_groups = []
    
    for group_name in groups:
        try:
            geofences = get_geofences_by_group(group_name)
            if not geofences:
                continue
            
            for gf in geofences:
                if is_inside_bounding_box(lat, lon, gf) or is_inside_circle(lat, lon, gf):
                    inside_groups.append(f"{group_name} ({gf['name']})")
                    break  # Found one in this group, no need to check more
        except Exception as e:
            continue
    
    return len(inside_groups) > 0, inside_groups

def print_test_prediction():
    """Print comprehensive prediction of test results based on rule logic"""
    print(f"\n{'='*80}")
    print("🎯 DEFCON TEST - PREDICTION & RULE ANALYSIS")
    print(f"{'='*80}\n")
    
    print("📋 DEFCON PROGRESSION RULES (from jammer_wargames.grl):")
    print("   DEFCON 0: Buffer update (always executes)")
    print("   DEFCON 1: PositioningStatus == 'V' (invalid GPS)")
    print("   DEFCON 2: Buffer has 10 + AvgSpeed >= 10 + AvgGSM < 15")
    print("   DEFCON 3: Offline >5min + NOT inside Taller + NOT inside CLIENTES")
    print("   DEFCON 4: OutsideAllSafeZones == true + Alert not sent yet\n")
    
    print("📊 TEST CONFIGURATION:")
    print(f"   Warmup: {WARMUP_PACKETS} valid packets (Status A, GSM=20, Speed=50)")
    print(f"   Simulation: {INVALID_PACKETS} invalid packets (Status V, GSM=10, Speed=100)")
    print(f"   Packet interval: {PACKET_INTERVAL}s")
    print(f"   Time to offline (5 min): {5*60/PACKET_INTERVAL:.0f} packets\n")
    
    defcon4_expected = []
    defcon2_expected = []
    
    for cfg in IMEI_CONFIG:
        lat = float(cfg["lat"])
        lon = float(cfg["lon"])
        inside_safe_zone, zones = check_if_inside_safe_zones(lat, lon)
        
        print(f"{'='*80}")
        print(f"📍 {cfg['name']} - {cfg['imei']}")
        print(f"{'='*80}")
        print(f"   Coordinates: ({lat}, {lon})")
def print_test_prediction():
    """Print detailed prediction of expected behavior for all IMEIs"""
    if not GEOFENCE_CHECK_AVAILABLE:
        print("⚠️  Geofence check unavailable - cannot predict DEFCON 3/4 behavior\n")
        return
    
    print(f"\n{'='*80}")
    print("🔮 DEFCON TEST PREDICTION - RULE ANALYSIS")
    print(f"{'='*80}\n")
    
    defcon4_expected = []
    defcon2_expected = []
    defcon1_expected = []
    
    for cfg in IMEI_CONFIG:
        lat = float(cfg["lat"])
        lon = float(cfg["lon"])
        gsm = int(cfg["gsm"])
        inside_safe_zone, zones = check_if_inside_safe_zones(lat, lon)
        
        print(f"{'='*80}")
        print(f"📍 {cfg['name']} - {cfg['imei']}")
        print(f"{'='*80}")
        print(f"   Coordinates: ({lat}, {lon})")
        print(f"   Test GSM: {gsm} ({'< 15 ✓' if gsm < 15 else '>= 15 ✗ - TOO STRONG'})")
        print(f"   Test Speed: 100 km/h (>= 10 ✓)")
        print(f"   Stops at step: {cfg['stop_at_step']}/{INVALID_PACKETS}\n")
        
        # Check geofences
        if inside_safe_zone:
            print(f"   🛡️  INSIDE SAFE ZONE: {', '.join(zones)}")
        else:
            print(f"   ⚠️  OUTSIDE ALL SAFE ZONES")
        
        print(f"\n   📈 DEFCON PROGRESSION PREDICTION:")
        print(f"   ✅ DEFCON 0: Will execute (buffer update)")
        print(f"   ✅ DEFCON 1: Will pass (Status V in Phase 2)")
        
        # Check DEFCON 2 requirements
        if gsm < 15:
            print(f"   ✅ DEFCON 2: Will pass (10 packets + Speed >= 10 + GSM < 15)")
            
            # Calculate if offline for 5 min
            packets_needed_for_offline = int(5 * 60 / PACKET_INTERVAL)  # 20 packets
            will_be_offline = cfg['stop_at_step'] >= packets_needed_for_offline
            
            if inside_safe_zone:
                print(f"   ❌ DEFCON 3: Will FAIL (inside safe zone: {', '.join(zones)})")
                print(f"   ❌ DEFCON 4: Will NOT trigger\n")
                print(f"   🎯 EXPECTED RESULT: Stops at DEFCON 2")
                print(f"      Reason: Vehicle is inside safe zone, rule blocks progression at DEFCON 3")
                defcon2_expected.append(cfg['imei'])
            else:
                if will_be_offline:
                    print(f"   ✅ DEFCON 3: Will pass (offline >5min + outside safe zones)")
                    print(f"   ✅ DEFCON 4: WILL TRIGGER ALERT 🚨\n")
                    print(f"   🎯 EXPECTED RESULT: DEFCON 4 - JAMMER ALERT FIRED")
                    print(f"      Reason: Outside all safe zones + offline >5min + weak GSM + moving")
                    defcon4_expected.append(cfg['imei'])
                else:
                    print(f"   ⏳ DEFCON 3: Will NOT pass yet (needs {packets_needed_for_offline - cfg['stop_at_step']} more packets for 5min offline)")
                    print(f"   ❌ DEFCON 4: Will NOT trigger (test stops too early)\n")
                    print(f"   🎯 EXPECTED RESULT: Stops at DEFCON 2")
                    print(f"      Reason: Test configured to stop at step {cfg['stop_at_step']} (before 5min offline)")
                    defcon2_expected.append(cfg['imei'])
        else:
            print(f"   ❌ DEFCON 2: Will FAIL (GSM {gsm} >= 15 - signal too strong, not a jammer)")
            print(f"   ❌ DEFCON 3: Will NOT reach (failed at DEFCON 2)")
            print(f"   ❌ DEFCON 4: Will NOT trigger\n")
            print(f"   🎯 EXPECTED RESULT: Stops at DEFCON 1")
            print(f"      Reason: GSM signal too strong ({gsm} >= 15), doesn't match jammer pattern")
            defcon1_expected.append(cfg['imei'])
        
        print(f"{'='*80}\n")
    
    # Summary
    print(f"{'='*80}")
    print("📊 TEST SUMMARY - EXPECTED RESULTS")
    print(f"{'='*80}")
    print(f"   🚨 IMEIs expected to trigger DEFCON 4: {len(defcon4_expected)}/4")
    for imei in defcon4_expected:
        cfg = next(c for c in IMEI_CONFIG if c['imei'] == imei)
        print(f"      ✅ {cfg['name']} - {imei}")
    
    print(f"\n   🛡️  IMEIs expected to stop at DEFCON 2: {len(defcon2_expected)}/4")
    for imei in defcon2_expected:
        cfg = next(c for c in IMEI_CONFIG if c['imei'] == imei)
        print(f"      ⏸️  {cfg['name']} - {imei}")
    
    print(f"\n   ⚠️  IMEIs expected to stop at DEFCON 1: {len(defcon1_expected)}/4")
    for imei in defcon1_expected:
        cfg = next(c for c in IMEI_CONFIG if c['imei'] == imei)
        print(f"      🔴 {cfg['name']} - {imei}")
    
    if len(defcon4_expected) >= 2:
        print(f"\n   ✅ TEST VALID: {len(defcon4_expected)} IMEIs will reach DEFCON 4 (>= 2 required)")
    else:
        print(f"\n   ❌ TEST INVALID: Only {len(defcon4_expected)} IMEIs will reach DEFCON 4 (need >= 2)")
    
    print(f"{'='*80}\n")

def send_packet(imei, output, label):
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(5)
        s.connect((SERVER, PORT))
        s.sendall(output.encode('utf-8'))
        s.close()
        print(f"[{datetime.now().strftime('%H:%M:%S')}] {label} - IMEI: {imei} - Sent OK")
    except Exception as e:
        print(f"[{datetime.now().strftime('%H:%M:%S')}] {label} - IMEI: {imei} - ERROR: {e}")

def run_test():
    # Print prediction BEFORE starting test
    print_test_prediction()
    
    print(f"\n{'='*60}")
    print(f"🚀 STARTING MULTI-IMEI TEST AT {datetime.now()}")
    print(f"Targets: 2x Defcon 4, 1x Defcon 2, 1x Defcon 3")
    print(f"{'='*60}\n")

    # Initialize next_send for each IMEI with a slight stagger (2 seconds apart)
    base_time = datetime.now()
    for i, cfg in enumerate(IMEI_CONFIG):
        cfg["next_send"] = base_time + timedelta(seconds=i * 2)
    
    # PHASE 1: WARMUP (11 VALID PACKETS)
    print("\n--- PHASE 1: WARMUP (11 Valid Packets per IMEI) ---")
    for step in range(WARMUP_PACKETS):
        for cfg in IMEI_CONFIG:
            # Sleep until it's time to send for this IMEI
            wait = (cfg["next_send"] - datetime.now()).total_seconds()
            if wait > 0:
                time.sleep(wait)
            
            utc = datetime.utcnow().strftime('%y%m%d%H%M%S')
            # Warmup is always valid Status A, strong signal, speed 50
            out = payload(
                imei=cfg["imei"], 
                eventCode="35", 
                latitude=cfg["lat"], 
                longitude=cfg["lon"], 
                utc=utc, 
                status="A", 
                sats="12", 
                gsmStrenght="20", 
                speed="50",
                direction="180", accuracy="1", altitude="100", mileage="0", runtime="1000",
                mcc="334", mnc="020", lac="1234", cellId="5678", 
                portStatus="0000", AD1="0000", AD2="0000", AD3="0000", battery="100", AD5="0000", eventInfo="00000000"
            )
            
            send_packet(cfg["imei"], out, f"WARMUP {step+1}/{WARMUP_PACKETS}")
            cfg["next_send"] += timedelta(seconds=1) # Fast warmup

    # PHASE 2: SIMULATION (24 INVALID PACKETS)
    print("\n--- PHASE 2: SIMULATION (24 Invalid Packets with Custom Signals) ---")
    # Reset next_send for normal interval
    base_time = datetime.now()
    for i, cfg in enumerate(IMEI_CONFIG):
        cfg["next_send"] = base_time + timedelta(seconds=i * 2)

    for step in range(INVALID_PACKETS):
        for cfg in IMEI_CONFIG:
            # Skip if we reached the stop step for this IMEI (to stay at a target Defcon)
            if step >= cfg["stop_at_step"]:
                continue

            wait = (cfg["next_send"] - datetime.now()).total_seconds()
            if wait > 0:
                time.sleep(wait)
            
            utc = datetime.utcnow().strftime('%y%m%d%H%M%S')
            
            # Defcon 2 and 4 need AvgSpeed >= 10
            speed = "100"
            
            out = payload(
                imei=cfg["imei"], 
                eventCode="35", 
                latitude=cfg["lat"], 
                longitude=cfg["lon"], 
                utc=utc, 
                status="V", # INVALID
                sats="0", 
                gsmStrenght=str(cfg["gsm"]), 
                speed=speed,
                direction="0", accuracy="0", altitude="0", mileage="0", runtime="1348",
                mcc="334", mnc="020", lac="1234", cellId="5678", 
                portStatus="0000", AD1="0000", AD2="0000", AD3="0000", battery="100", AD5="0000", eventInfo="00000000"
            )
            
            send_packet(cfg["imei"], out, f"SIMUL {step+1}/{INVALID_PACKETS} (T:{cfg['target']})")
            cfg["next_send"] += timedelta(seconds=PACKET_INTERVAL)

if __name__ == "__main__":
    run_test()
