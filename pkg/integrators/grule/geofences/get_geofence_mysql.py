import mysql.connector
import os
from connect2server import get_db_connection

def get_geofences_from_db():
    """Obtiene todas las geocercas de la base de datos."""
    conn = get_db_connection(suppress_message=True, database="geofences")
    if not conn:
        return []

    geofences = []
    try:
        cursor = conn.cursor(dictionary=True)
        query = "SELECT * FROM geofences ORDER BY name"
        cursor.execute(query)
        for row in cursor:
            geofences.append(row)
    except mysql.connector.Error as err:
        print(f"Error al leer geocercas: {err}")
    finally:
        cursor.close()
        conn.close()
    return geofences

def get_geofences_by_group(group_name):
    """Obtiene todas las geocercas de un grupo específico."""
    conn = get_db_connection(suppress_message=True, database="geofences")
    if not conn:
        return None
    
    try:
        cursor = conn.cursor(dictionary=True)
        
        # First, get the group ID
        query_group = "SELECT id FROM geofence_groups WHERE name = %s"
        cursor.execute(query_group, (group_name,))
        group_result = cursor.fetchone()
        
        if not group_result:
            print(f"❌ Grupo '{group_name}' no encontrado")
            return None
        
        group_id = group_result['id']
        
        # Get all geofences for this group
        query = """
            SELECT g.* FROM geofences g
            INNER JOIN geofence_group_mapping gm ON g.id = gm.geofence_id
            WHERE gm.group_id = %s
            ORDER BY g.name
        """
        cursor.execute(query, (group_id,))
        geofences = cursor.fetchall()
        
        return geofences
    
    except mysql.connector.Error as err:
        print(f"Error al consultar geocercas: {err}")
        return None
    finally:
        cursor.close()
        conn.close()

def print_geofences(geofences, title="Geocercas"):
    """Imprime las geocercas de forma legible."""
    if not geofences:
        print("No se encontraron geocercas")
        return
    
    print(f"\n{'='*80}")
    print(f"📍 {title}: {len(geofences)} geocercas encontradas")
    print(f"{'='*80}\n")
    
    for i, gf in enumerate(geofences, 1):
        print(f"{i}. {gf['name']}")
        print(f"   ID: {gf['id']}")
        print(f"   Tipo: {gf['shapeType']}")
        print(f"   Descripción: {gf['description'] or 'N/A'}")
        
        if gf['shapeType'] == 'circle':
            if gf['centerLat'] is not None:
                print(f"   Posición: ({gf['centerLat']}, {gf['centerLon']})")
            if gf['radius'] is not None:
                print(f"   Radio: {gf['radius']} metros")
        elif gf['shapeType'] == 'polygon':
            if gf['boundingBoxMinX'] is not None:
                print(f"   Área: ({gf['boundingBoxMinX']:.4f} a {gf['boundingBoxMaxX']:.4f}, "
                      f"{gf['boundingBoxMinY']:.4f} a {gf['boundingBoxMaxY']:.4f})")
        
        print()

if __name__ == "__main__":
    # Opción 1: Obtener todas las geocercas
    print("🔍 Obteniendo todas las geocercas...")
    all_geofences = get_geofences_from_db()
    print_geofences(all_geofences, "Todas las geocercas")
    
    # Opción 2: Obtener geocercas de un grupo específico
    print("\n" + "="*80)
    group_name = "CLIENTES"
    print(f"🔍 Buscando geocercas del grupo '{group_name}'...")
    group_geofences = get_geofences_by_group(group_name)
    
    if group_geofences is not None:
        print_geofences(group_geofences, f"Grupo {group_name}")
        print(f"{'='*80}")
        print(f"✅ Total en grupo: {len(group_geofences)} geocercas")
        print(f"{'='*80}")

