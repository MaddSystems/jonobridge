#!/usr/bin/env python3
"""
Check if coordinates are inside geofences for DEFCON testing.
Helps validate that test coordinates will trigger the expected DEFCON levels.
"""
import sys
import os

# Add parent directory to path to import geofences module
sys.path.append(os.path.join(os.path.dirname(__file__), '..', 'geofences'))

from get_geofence_mysql import get_geofences_by_group

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

def check_coordinate(lat, lon, label="Point"):
    """Check if coordinate is inside any safe zone groups"""
    print(f"\n{'='*80}")
    print(f"🔍 Checking: {label}")
    print(f"   Coordinates: ({lat}, {lon})")
    print(f"{'='*80}")
    
    # Check all critical groups
    groups = ['CLIENTES', 'Taller', 'Resguardo/Cedis/Puerto']
    
    inside_any = False
    
    for group_name in groups:
        geofences = get_geofences_by_group(group_name)
        if not geofences:
            print(f"⚠️  Group '{group_name}' not found or empty")
            continue
        
        inside_group = False
        matched_geofences = []
        
        for gf in geofences:
            if is_inside_bounding_box(lat, lon, gf) or is_inside_circle(lat, lon, gf):
                inside_group = True
                inside_any = True
                matched_geofences.append(gf['name'])
        
        if inside_group:
            print(f"   ✅ INSIDE '{group_name}': {', '.join(matched_geofences)}")
        else:
            print(f"   ❌ OUTSIDE '{group_name}'")
    
    print(f"\n{'='*80}")
    if inside_any:
        print(f"❌ {label}: INSIDE SAFE ZONE → Will NOT trigger DEFCON 4")
        print(f"   Expected: DEFCON 2 (if moving with weak signal)")
    else:
        print(f"✅ {label}: OUTSIDE ALL SAFE ZONES → CAN trigger DEFCON 4")
        print(f"   Expected: DEFCON 4 (if offline >5min + moving + weak GSM)")
    print(f"{'='*80}\n")
    
    return not inside_any  # True if outside all zones (can trigger DEFCON 4)

def main():
    print("\n" + "="*80)
    print("🧪 DEFCON TEST COORDINATE VALIDATOR")
    print("="*80)
    
    # Import from send_multiple.py to get actual test coordinates
    sys.path.insert(0, os.path.dirname(__file__))
    from send_multiple import IMEI_CONFIG
    
    # Convert to test format
    test_coords = []
    for cfg in IMEI_CONFIG:
        test_coords.append({
            "imei": cfg["imei"],
            "name": cfg["name"],
            "lat": float(cfg["lat"]),
            "lon": float(cfg["lon"]),
        })
    
    results = []
    for coord in test_coords:
        can_trigger_defcon4 = check_coordinate(
            coord["lat"], 
            coord["lon"], 
            f"{coord['name']} - {coord['imei']}"
        )
        results.append({
            "imei": coord["imei"],
            "name": coord["name"],
            "can_trigger_defcon4": can_trigger_defcon4
        })
    
    # Summary
    print("\n" + "="*80)
    print("📊 SUMMARY")
    print("="*80)
    
    defcon4_count = sum(1 for r in results if r["can_trigger_defcon4"])
    
    for r in results:
        status = "✅ CAN" if r["can_trigger_defcon4"] else "❌ CANNOT"
        print(f"   {status} trigger DEFCON 4: {r['name']}")
    
    print(f"\n   Total IMEIs that CAN trigger DEFCON 4: {defcon4_count}/4")
    
    if defcon4_count >= 2:
        print(f"   ✅ TEST VALID: At least 2 IMEIs can reach DEFCON 4")
    else:
        print(f"   ❌ TEST INVALID: Need at least 2 IMEIs outside safe zones")
        print(f"\n   💡 SUGGESTED FIX:")
        print(f"      Use coordinates OUTSIDE all geofences for DEFCON 4 test")
        print(f"      Example: 19.600000, -99.300000 (random location far from geofences)")
    
    print("="*80 + "\n")

if __name__ == "__main__":
    main()
