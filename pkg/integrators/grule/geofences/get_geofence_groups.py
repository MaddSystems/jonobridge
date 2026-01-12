import requests
from connect2server import get_server1_token
import time

def get_geofence_groups(token, appid):
    url = f"https://api4server1.gpscontrol.com.mx/applications/{appid}/geofencegroups"
    headers = {'Authorization': f'Bearer {token}'}
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        return {'data': response.json()}
    else:
        return {'error': f"HTTP {response.status_code}: {response.text}"}

def get_geofence_by_id(token, appid, geofence_id, timeout=10):
    url = f"https://api4server1.gpscontrol.com.mx/applications/{appid}/geofences/{geofence_id}"
    headers = {'Authorization': f'Bearer {token}'}
    try:
        response = requests.get(url, headers=headers, timeout=timeout)
        if response.status_code == 200:
            return {'data': response.json()}
        else:
            return {'error': f"HTTP {response.status_code}: {response.text}"}
    except requests.exceptions.Timeout:
        return {'error': f"Timeout getting geofence {geofence_id}"}
    except requests.exceptions.RequestException as e:
        return {'error': f"Error getting geofence {geofence_id}: {str(e)}"}

def get_details_for_groups(username, passcode, appid, group_names, max_geofences=None):
    token_result = get_server1_token(username, passcode, appid)
    if token_result['error']:
        return {'error': token_result['error']}
    token = token_result['token']
    groups_result = get_geofence_groups(token, appid)
    if 'error' in groups_result:
        return groups_result
    groups = groups_result['data']
    details = {}
    for group in groups:
        if group.get('name') in group_names:
            group_name = group['name']
            geofence_ids = group.get('geofenceIds', [])
            
            # Limit geofences if specified
            if max_geofences:
                geofence_ids = geofence_ids[:max_geofences]
                print(f"Processing {group_name}: {len(geofence_ids)} geofences (limited to {max_geofences})")
            else:
                print(f"Processing {group_name}: {len(geofence_ids)} geofences")
            
            geofences = []
            for i, gid in enumerate(geofence_ids, 1):
                print(f"  [{i}/{len(geofence_ids)}] Fetching geofence {gid}...", end='\r')
                gf_result = get_geofence_by_id(token, appid, gid)
                if 'error' in gf_result:
                    geofences.append({'id': gid, 'error': gf_result['error']})
                else:
                    geofences.append(gf_result['data'])
                time.sleep(0.1)  # Small delay to avoid overwhelming the server
            print(f"  [{len(geofence_ids)}/{len(geofence_ids)}] Complete!                              ")
            details[group_name] = geofences
    return {'data': details}

if __name__ == "__main__":
    username = "admindesarrollo"
    passcode = "GPSc0ntr0l00"
    appid = 112
    geofence_id = 154519
    # Get token once
    token_result = get_server1_token(username, passcode, appid)
    if token_result['error']:
        print("Error getting token:", token_result['error'])
        exit(1)
    token = token_result['token']
    # Test get groups
    print("Fetching geofence groups...")
    groups = get_geofence_groups(token, appid)
    print(f"Found {len(groups['data'])} groups\n")
    
    # Test get geofence by id
    print(f"Fetching single geofence {geofence_id}...")
    geofence = get_geofence_by_id(token, appid, geofence_id)
    print("Geofence by ID:", geofence, "\n")
    
    # Test get details for specific groups (limit to 5 geofences per group for testing)
    print("Fetching details for specific groups (limited to 5 geofences each for testing)...")
    group_names = ["Resguardo/CEDIS/Puerto", "Taller", "CLIENTES"]
    details = get_details_for_groups(username, passcode, appid, group_names, max_geofences=5)
    print("\nDetails for specific groups (sample):", details)
