import requests
import json
import time
from connect2server import get_db_connection

def get_server1_token(username, passcode, appid):
    """Obtiene un token de autenticación de la plataforma GPS."""
    token = ""
    error = ""

    if passcode != "":
        url = "http://server1.gpscontrol.com.mx/comGPSGate/api/v.1/applications/" + str(appid) + "/tokens"
        r = requests.post(url, json={"username": username, "password": passcode})
        t = r.text[1:69]
        if r.status_code == 200:
            if t != "The user does not have neither _APIRead nor _APIReadWrite privileges. To get an API token, please assign required privileges.":
                j = r.json()
                token = j["token"]
            else:
                error = "Faltan permisos en la cuenta. _APIRead y _APIReadWrite"
        else:
            error = "Respuesta de la plataforma: " + r.text
    else:
        error = "El token no tiene keycode valido."
    return dict(token=token, error=error)

def get_geofence_by_id(token, appid, geofence_id, timeout=10):
    """Obtiene los detalles de una geocerca por ID."""
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

def get_geofence_groups(token, appid):
    """Obtiene los grupos de geocercas disponibles."""
    url = f"https://api4server1.gpscontrol.com.mx/applications/{appid}/geofencegroups"
    headers = {'Authorization': f'Bearer {token}'}
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        return {'data': response.json()}
    else:
        return {'error': f"HTTP {response.status_code}: {response.text}"}

if __name__ == "__main__":
    username = "admindesarrollo"
    passcode = "GPSc0ntr0l00"
    appid = 112
    
    # Get token
    print("Obteniendo token...")
    token_result = get_server1_token(username, passcode, appid)
    if token_result['error']:
        print(f"Error: {token_result['error']}")
        exit(1)
    token = token_result['token']
    print(f"✅ Token obtenido\n")
    
    # Get groups
    print("Obteniendo grupos...")
    groups_result = get_geofence_groups(token, appid)
    if 'error' in groups_result:
        print(f"Error: {groups_result['error']}")
        exit(1)
    
    groups = groups_result['data']
    print(f"✅ {len(groups)} grupos encontrados\n")
    
    # Find CLIENTES group and get first few geofences
    for group in groups:
        if group.get('name') == 'CLIENTES':
            geofence_ids = group.get('geofenceIds', [])[:5]  # Solo primeros 5
            print(f"📍 Grupo: {group.get('name')}")
            print(f"📊 Primeras {len(geofence_ids)} geocercas para inspeccionar:\n")
            
            for i, gid in enumerate(geofence_ids, 1):
                print(f"\n{'='*80}")
                print(f"Geocerca [{i}] ID: {gid}")
                print(f"{'='*80}")
                
                gf_result = get_geofence_by_id(token, appid, gid)
                
                if 'error' in gf_result:
                    print(f"❌ Error: {gf_result['error']}")
                else:
                    data = gf_result['data']
                    print(json.dumps(data, indent=2, ensure_ascii=False))
                
                time.sleep(0.5)
            break
