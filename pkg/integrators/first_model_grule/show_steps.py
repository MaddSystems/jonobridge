import mysql.connector
import json
import sys
from datetime import datetime

def main():
    # Obtener el parámetro de IMEI
    if len(sys.argv) < 2:
        print("Error: Debe proporcionar un IMEI como parámetro")
        print("Uso: python3 show_steps.py <IMEI>")
        return

    imei = sys.argv[1]

    try:
        # Conectar a la base de datos
        conn = mysql.connector.connect(
            host='localhost',
            user='gpscontrol',
            password='qazwsxedc',
            database='grule'
        )

        cursor = conn.cursor(dictionary=True)

        # Query para obtener todos los registros para el IMEI especificado, ordenados por tiempo
        query = """
        SELECT id, imei, rule_name, execution_time, context_snapshot, buffer_size, stage_reached, stop_reason
        FROM rule_execution_state
        WHERE imei = %s
        ORDER BY execution_time ASC, id ASC
        """

        cursor.execute(query, (imei,))
        results = cursor.fetchall()

        if results:
            print(f"Pasos encontrados para IMEI {imei}: {len(results)} registros")
            print("=" * 80)

            for i, result in enumerate(results, 1):
                print(f"\n--- Paso {i} ---")
                print(f"ID: {result['id']}")
                print(f"Execution Time: {result['execution_time']}")
                print(f"Buffer Size: {result['buffer_size']}")
                print(f"Stage Reached: {result['stage_reached']}")
                print(f"Stop Reason: {result['stop_reason']}")

                # Parsear y mostrar el context_snapshot
                if result['context_snapshot']:
                    try:
                        snapshot = json.loads(result['context_snapshot'])
                        buffer_size = len(snapshot.get('buffer_circular', []))
                        print(f"Buffer Circular: {buffer_size} posiciones")

                        # Mostrar resumen del buffer
                        if 'buffer_circular' in snapshot and snapshot['buffer_circular']:
                            print("Últimas posiciones en buffer:")
                            for j, pos in enumerate(snapshot['buffer_circular'][-3:], 1):  # Últimas 3
                                print(f"  {j}. Vel: {pos.get('speed', 'N/A')} km/h, "
                                      f"Lat: {pos.get('latitude', 'N/A'):.6f}, "
                                      f"Lon: {pos.get('longitude', 'N/A'):.6f}, "
                                      f"Time: {pos.get('datetime', 'N/A')}")

                        # Mostrar métricas de jammer si están listas
                        if snapshot.get('wrapper_flags', {}).get('MetricsReady', False):
                            metrics = snapshot.get('jammer_metrics', {})
                            print(f"Métricas: avg_gsm_last5={metrics.get('avg_gsm_last5', 'N/A')}, "
                                  f"jammer_positions={metrics.get('jammer_positions', 'N/A')}")

                        # Mostrar si se disparó alerta
                        if snapshot.get('wrapper_flags', {}).get('AlertFired', False):
                            print("⚠️ ALERTA DISPARADA")

                    except json.JSONDecodeError as e:
                        print(f"Error al parsear context_snapshot: {e}")
                else:
                    print("Context Snapshot: NULL")

                print("-" * 40)
        else:
            print(f"No se encontraron registros para el IMEI {imei}")

        cursor.close()
        conn.close()

    except mysql.connector.Error as err:
        print(f"Error de MySQL: {err}")
    except Exception as e:
        print(f"Error general: {e}")

if __name__ == "__main__":
    main()