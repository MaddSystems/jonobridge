import requests # Módulo para hacer peticiones HTTP
import json # Módulo para manejar datos JSON
import hashlib # Probablemente necesario para generar las firmas (fuera de las fuentes)
import time # Para generar timestamp
import base64 # Para base64 encoding

# 1. Definición de Credenciales y URL (Extraídas de las fuentes [1, 2, 5, 6])
API_ACCOUNT = '292508153084379141'
PRIVATE_KEY = 'a0a1047cce70493c9d5d29704f05d0d9'
URL_WEBSERVICE = 'https://demoopenapi.jtjms-mx.com/webopenplatformapi/transport/gps/trajectoryUpload?uuid=9b2aa6c36a3f4c94aab0ade672c82786'

# --- FUNCIONES DE SEGURIDAD ---
# Ya no necesarias, usando digest directo

# 2. Definición de Parámetros de Negocio (Payload)
# ESTOS SON PLACEHOLDERS. DEBEN SER REEMPLAZADOS POR LOS CAMPOS OBLIGATORIOS ('Y')
# DEL SERVICIO 'trajectoryUpload' [2, 6].

business_parameters = [
    {
        "plateNumber": "ECO 102",
        "longitude": "-99.182505",
        "latitude": "19.546546",
        "address": "Mexico City, Mexico",  # Dirección aproximada basada en coordenadas
        "getDataTime": "2025-10-20 21:44:41",  # Fecha del reporte GPS
        "uploadDataTime": time.strftime("%Y-%m-%d %H:%M:%S"),  # Hora actual de subida
        "speed": "0",
        "direction": "0"  # Dirección no especificada, usar 0
    }
]

# 3. Serialización y Generación de Firmas
# El cuerpo del mensaje debe ser enviado como una cadena de texto JSON.
bizContent_json = json.dumps(business_parameters)

# Generar el Digest
business_digest = base64.b64encode(hashlib.md5((bizContent_json + PRIVATE_KEY).encode('utf-8')).digest()).decode('utf-8')

# 4. Construcción de Headers
timestamp_str = str(int(time.time() * 1000))
headers = {
    "apiAccount": API_ACCOUNT,
    "digest": business_digest,
    "timestamp": timestamp_str,
    "timezone": "GMT-6"  # Zona horaria
}

# 5. Envío de la Petición
try:
    print(f"Enviando solicitud a: {URL_WEBSERVICE}")

    # La solicitud POST lleva bizContent en el cuerpo como form data
    response = requests.post(
        URL_WEBSERVICE,
        data={'bizContent': bizContent_json},
        headers=headers
    )

    # 6. Manejo de la Respuesta
    if response.status_code == 200:
        print("Prueba de comunicación exitosa (código 200).")
        # Recuerde que se necesitan TRES pruebas exitosas [1, 5].
        print("Respuesta del servicio:")
        print(response.json())
    else:
        print(f"Error en la petición: Código HTTP {response.status_code}")
        print("Respuesta (posibles códigos de error):")
        # Los códigos de error se pueden consultar en el apartado TEST TOOLS [3, 7].
        print(response.text)

except requests.exceptions.RequestException as e:
    print(f"Error de conexión: {e}")

