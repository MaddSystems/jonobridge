import mysql.connector
import json
import sys
from datetime import datetime

def main():
    # IMEI fijo
    imei = "864352043111848"

    # Obtener el parámetro de posición (1-based, default 1)
    position = 1
    if len(sys.argv) > 1:
        try:
            position = int(sys.argv[1])
            if position < 1:
                print("Error: La posición debe ser un número entero positivo (1, 2, 3, ...)")
                return
        except ValueError:
            print("Error: El parámetro debe ser un número entero")
            return

    try:
        # Conectar a la base de datos
        conn = mysql.connector.connect(
            host='localhost',
            user='gpscontrol',
            password='qazwsxedc',
            database='grule'
        )

        cursor = conn.cursor(dictionary=True)

        # Query para obtener el registro en la posición especificada para el IMEI (ordenado por fecha ascendente)
        offset = position - 1
        query = """
        SELECT id, imei, rule_name, execution_time, context_snapshot
        FROM rule_execution_state
        WHERE imei = %s
        ORDER BY execution_time ASC, id ASC
        LIMIT 1 OFFSET %s
        """

        cursor.execute(query, (imei, offset))
        result = cursor.fetchone()

        if result:
            print(f"Registro #{position} más antiguo encontrado para IMEI {imei}:")
            print(f"ID: {result['id']}")
            print(f"IMEI: {result['imei']}")
            print(f"Rule Name: {result['rule_name']}")
            print(f"Execution Time: {result['execution_time']}")

            # Parsear y mostrar el context_snapshot
            if result['context_snapshot']:
                try:
                    snapshot = json.loads(result['context_snapshot'])
                    print("\nContext Snapshot (JSON detallado):")
                    print(json.dumps(snapshot, indent=2, ensure_ascii=False))
                except json.JSONDecodeError as e:
                    print(f"Error al parsear context_snapshot: {e}")
                    print(f"Raw context_snapshot: {result['context_snapshot']}")
            else:
                print("Context Snapshot: NULL")
        else:
            print(f"No se encontró el registro #{position} para IMEI {imei} (posiblemente no hay suficientes registros)")

        cursor.close()
        conn.close()

    except mysql.connector.Error as err:
        print(f"Error de MySQL: {err}")
    except Exception as e:
        print(f"Error general: {e}")

if __name__ == "__main__":
    main()