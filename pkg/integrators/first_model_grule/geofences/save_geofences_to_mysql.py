import mysql.connector
import os
import json
import requests
import time
from connect2server import get_db_connection, create_database_and_tables

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

def save_geofence(name, description=None, geo_type="circle", coordinates=None, radius=None):
    """Guarda o actualiza una geocerca en la base de datos. Retorna el ID."""
    conn = get_db_connection(suppress_message=True)
    if not conn:
        return None

    try:
        cursor = conn.cursor()
        # Check if geofence already exists
        cursor.execute("SELECT id FROM geofences WHERE name = %s", (name,))
        result = cursor.fetchone()

        if result:
            geofence_id = result[0]
            # Update existing geofence
            if geo_type == "circle" and coordinates:
                centerLat = coordinates.get("latitude")
                centerLon = coordinates.get("longitude")
                query = """
                    UPDATE geofences
                    SET description = %s, shapeType = %s, centerLat = %s, centerLon = %s, radius = %s
                    WHERE id = %s
                """
                cursor.execute(query, (description, geo_type, centerLat, centerLon, radius, geofence_id))
            elif geo_type == "polygon" and coordinates:
                # For polygons, calculate bounding box
                lats = [c["latitude"] for c in coordinates]
                lons = [c["longitude"] for c in coordinates]
                query = """
                    UPDATE geofences
                    SET description = %s, shapeType = %s, boundingBoxMinX = %s, boundingBoxMaxX = %s, 
                        boundingBoxMinY = %s, boundingBoxMaxY = %s
                    WHERE id = %s
                """
                cursor.execute(query, (description, geo_type, min(lons), max(lons), min(lats), max(lats), geofence_id))
            print(f"Geocerca '{name}' actualizada.")
        else:
            # Get next ID
            cursor.execute("SELECT MAX(id) FROM geofences")
            max_id_result = cursor.fetchone()
            geofence_id = (max_id_result[0] or 0) + 1
            
            # Insert new geofence
            if geo_type == "circle" and coordinates:
                centerLat = coordinates.get("latitude")
                centerLon = coordinates.get("longitude")
                query = """
                    INSERT INTO geofences (id, name, description, shapeType, centerLat, centerLon, radius)
                    VALUES (%s, %s, %s, %s, %s, %s, %s)
                """
                cursor.execute(query, (geofence_id, name, description, geo_type, centerLat, centerLon, radius))
            elif geo_type == "polygon" and coordinates:
                # For polygons, calculate bounding box
                lats = [c["latitude"] for c in coordinates]
                lons = [c["longitude"] for c in coordinates]
                query = """
                    INSERT INTO geofences (id, name, description, shapeType, boundingBoxMinX, boundingBoxMaxX, 
                                          boundingBoxMinY, boundingBoxMaxY)
                    VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                """
                cursor.execute(query, (geofence_id, name, description, geo_type, min(lons), max(lons), min(lats), max(lats)))
            print(f"Geocerca '{name}' insertada.")
        
        conn.commit()
        return geofence_id
    except mysql.connector.Error as err:
        print(f"Error al guardar geocerca: {err}")
        conn.rollback()
        return None
    finally:
        conn.close()

def save_geofence_group(group_id, group_name):
    """Guarda o actualiza un grupo de geocercas en la base de datos."""
    conn = get_db_connection(suppress_message=True)
    if not conn:
        return False

    try:
        cursor = conn.cursor()
        # Check if group already exists
        cursor.execute("SELECT id FROM geofence_groups WHERE id = %s", (group_id,))
        result = cursor.fetchone()

        if not result:
            # Insert new group
            query = "INSERT INTO geofence_groups (id, name) VALUES (%s, %s)"
            cursor.execute(query, (group_id, group_name))
        
        conn.commit()
        return True
    except mysql.connector.Error as err:
        print(f"Error al guardar grupo: {err}")
        conn.rollback()
        return False
    finally:
        conn.close()

def save_geofence_group_mapping(group_id, geofence_id):
    """Guarda la relación entre un grupo y una geocerca."""
    conn = get_db_connection(suppress_message=True)
    if not conn:
        return False

    try:
        cursor = conn.cursor()
        # Check if mapping already exists
        cursor.execute("SELECT 1 FROM geofence_group_mapping WHERE group_id = %s AND geofence_id = %s", 
                      (group_id, geofence_id))
        result = cursor.fetchone()

        if not result:
            # Insert new mapping
            query = "INSERT INTO geofence_group_mapping (group_id, geofence_id) VALUES (%s, %s)"
            cursor.execute(query, (group_id, geofence_id))
            conn.commit()
        
        return True
    except mysql.connector.Error as err:
        print(f"Error al guardar mapeo: {err}")
        conn.rollback()
        return False
    finally:
        conn.close()

def save_geofences_from_groups(username, passcode, appid, group_names, max_geofences=None, is_test=False):
    """Descarga geofences de grupos especificados y los guarda en la base de datos."""
    # Apply test limit if is_test is True
    if is_test:
        max_geofences = 5
        print("🧪 MODO TESTING - Limitado a 5 geocercas por grupo\n")
    
    # Get token
    token_result = get_server1_token(username, passcode, appid)
    if token_result['error']:
        print(f"Error al obtener token: {token_result['error']}")
        return False
    
    token = token_result['token']
    
    # Get groups
    groups_result = get_geofence_groups(token, appid)
    if 'error' in groups_result:
        print(f"Error al obtener grupos: {groups_result['error']}")
        return False
    
    groups = groups_result['data']
    total_saved = 0
    total_errors = 0
    
    for group in groups:
        if group.get('name') in group_names:
            group_name = group['name']
            group_id = group.get('id')
            geofence_ids = group.get('geofenceIds', [])
            
            # Save group first
            if group_id:
                save_geofence_group(group_id, group_name)
            
            # Limit geofences if specified
            if max_geofences:
                geofence_ids = geofence_ids[:max_geofences]
                print(f"\n📍 Procesando grupo '{group_name}': {len(geofence_ids)} geocercas (limitadas a {max_geofences})")
            else:
                print(f"\n📍 Procesando grupo '{group_name}': {len(geofence_ids)} geocercas")
            
            for i, gid in enumerate(geofence_ids, 1):
                print(f"  [{i}/{len(geofence_ids)}] Descargando geocerca {gid}...", end='\r')
                gf_result = get_geofence_by_id(token, appid, gid)
                
                if 'error' in gf_result:
                    print(f"  [{i}/{len(geofence_ids)}] ❌ Error en geocerca {gid}: {gf_result['error']}")
                    total_errors += 1
                else:
                    geofence_data = gf_result['data']
                    name = geofence_data.get('name', f'Geofence_{gid}')
                    description = geofence_data.get('description', group_name)
                    shape_type = geofence_data.get('shapeType', '').strip()
                    geofence_id = None
                    
                    # Check for Circle shape
                    if shape_type == 'Circle':
                        circle_shape = geofence_data.get('circleShape', {})
                        center = circle_shape.get('center', {})
                        lat = center.get('latitude')
                        lon = center.get('longitude')
                        radius = circle_shape.get('radius')
                        
                        if lat is not None and lon is not None:
                            coordinates = {"latitude": lat, "longitude": lon}
                            geofence_id = save_geofence(name, description, "circle", coordinates, radius)
                            if geofence_id:
                                total_saved += 1
                                # Save mapping if group_id exists
                                if group_id:
                                    save_geofence_group_mapping(group_id, geofence_id)
                            else:
                                total_errors += 1
                        else:
                            print(f"  [{i}/{len(geofence_ids)}] ⚠️  Circle sin coordenadas: {name}")
                            total_errors += 1
                    
                    # Check for Polygon shape
                    elif shape_type == 'Polygon':
                        polygon_shape = geofence_data.get('polygonShape', {})
                        vertices = polygon_shape.get('vertices', [])
                        
                        if vertices and len(vertices) > 0:
                            coordinates = [{"latitude": v.get('latitude'), "longitude": v.get('longitude')} for v in vertices]
                            geofence_id = save_geofence(name, description, "polygon", coordinates)
                            if geofence_id:
                                total_saved += 1
                                # Save mapping if group_id exists
                                if group_id:
                                    save_geofence_group_mapping(group_id, geofence_id)
                            else:
                                total_errors += 1
                        else:
                            print(f"  [{i}/{len(geofence_ids)}] ⚠️  Polygon sin vertices: {name}")
                            total_errors += 1
                    
                    # Check for Route shape
                    elif shape_type == 'Route':
                        route_shape = geofence_data.get('routeShape', {})
                        vertices = route_shape.get('vertices', [])
                        radius = route_shape.get('radius')
                        
                        if vertices and len(vertices) > 0:
                            # Treat route as polygon with the given vertices
                            coordinates = [{"latitude": v.get('latitude'), "longitude": v.get('longitude')} for v in vertices]
                            geofence_id = save_geofence(name, description, "polygon", coordinates)
                            if geofence_id:
                                total_saved += 1
                                # Save mapping if group_id exists
                                if group_id:
                                    save_geofence_group_mapping(group_id, geofence_id)
                            else:
                                total_errors += 1
                        else:
                            print(f"  [{i}/{len(geofence_ids)}] ⚠️  Route sin vertices: {name}")
                            total_errors += 1
                    
                    else:
                        print(f"  [{i}/{len(geofence_ids)}] ⚠️  Tipo no soportado: {shape_type}")
                        total_errors += 1
                
                time.sleep(0.1)  # Pequeño delay para no sobrecargar el servidor
            
            print(f"  [{len(geofence_ids)}/{len(geofence_ids)}] ✅ Grupo completado!                              ")
    
    print(f"\n{'='*60}")
    print(f"📊 RESUMEN: {total_saved} geocercas guardadas, {total_errors} errores")
    print(f"{'='*60}")
    
    return True

if __name__ == "__main__":
    # Crear base de datos y tablas si no existen
    print("📋 Inicializando base de datos...\n")
    create_database_and_tables()
    
    # Configuración de credenciales
    username = "admindesarrollo"
    passcode = "GPSc0ntr0l00"
    appid = 112
    
    # Grupos a descargar
    group_names = ["Resguardo/CEDIS/Puerto", "Taller", "CLIENTES"]
    
    # OPCIONES DE LIMITACIÓN:
    # is_test = True  → 5 geocercas por grupo (testing)
    # is_test = False + max_geofences = None → Todas las geocercas
    # is_test = False + max_geofences = 10 → Solo 10 geocercas por grupo
    
    is_test = False
    max_geofences = None  # Cambiar a un número para limitar (ej: 10, 50, 100)
    
    # Ejecutar descarga y guardado
    print("🚀 Iniciando descarga y guardado de geocercas...\n")
    save_geofences_from_groups(username, passcode, appid, group_names, max_geofences=max_geofences, is_test=is_test)
