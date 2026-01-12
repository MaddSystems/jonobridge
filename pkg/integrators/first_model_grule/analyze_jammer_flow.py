import mysql.connector
import json
from datetime import datetime

# Configuración de conexión
db_config = {
    'host': 'localhost',
    'user': 'gpscontrol',
    'password': 'qazwsxedc',
    'database': 'grule'
}

def analyze_jammer_flow(imei):
    conn = mysql.connector.connect(**db_config)
    cursor = conn.cursor(dictionary=True)

    # Traemos los últimos 40 frames para ver la evolución
    query = """
        SELECT id, stage_reached, stop_reason, context_snapshot, execution_time 
        FROM rule_execution_state 
        WHERE imei = %s 
        ORDER BY execution_time ASC 
        LIMIT 40
    """
    
    cursor.execute(query, (imei,))
    rows = cursor.fetchall()

    print(f"{'ID':<6} | {'TIME':<10} | {'STAGE':<20} | {'B10':<3} | {'OFF':<3} | {'MET':<3} | {'REASON'}")
    print("-" * 80)

    for row in rows:
        snap = json.loads(row['context_snapshot'])
        flags = snap.get('wrapper_flags', {})
        
        # Extraemos banderas críticas
        b10 = "✅" if flags.get('BufferHas10') else "❌"
        off = "✅" if flags.get('IsOfflineFor5Min') else "❌"
        met = "✅" if flags.get('MetricsReady') else "❌"
        
        time_str = row['execution_time'].strftime('%H:%M:%S')
        
        print(f"{row['id']:<6} | {time_str:<10} | {row['stage_reached']:<20} | {b10:<3} | {off:<3} | {met:<3} | {row['stop_reason']}")

    conn.close()

if __name__ == "__main__":
    # Cambia esto por el IMEI que estás probando
    analyze_jammer_flow('864352045580768')