import requests

# --- CONFIGURACIÓN BASADA EN EL DOCUMENTO FUENTE ---

# 1. URL del Ambiente de Pruebas (WSDL endpoint) [5]
URL_PRUEBAS = 'https://compassvmqa.centralus.cloudapp.azure.com/locations/locationReceiver.wsdl'

# 2. Credenciales (Reemplazar con las credenciales reales de Landstar Metro) [3, 5]
# !!! INSERTE SUS CREDENCIALES AQUÍ !!!
USUARIO_PRUEBA = "SU_USUARIO_AQUI"
CONTRASENA_PRUEBA = "SU_CONTRASENA_AQUI"

# Headers necesarios para las solicitudes SOAP
HEADERS = {'Content-Type': 'text/xml;charset=UTF-8',
           'SOAPAction': '' # Para este tipo de servicios, a veces es vacío o puede requerir el nombre del método
          }

# 3. Datos de Evento de Geolocalización (Usando el "Request Minimal" de ejemplo) [6]
# Estos son datos de ejemplo para enviar la posición.
ASSET_PLACA_EJEMPLO = "55BB1T" # Placa del vehículo [7]
LATITUD_EJEMPLO = "19.680356" # Latitud [7]
LONGITUD_EJEMPLO = "-99.201784" # Longitud [7]
FECHA_EJEMPLO = "2020-07-15T10:13:00.000" # Formato YYYY-MM-DTHH:MM:SS [7]
CODIGO_EVENTO = "0" # 0: Sin Evento (solo ubicación) [8]
DIRECCION_EJEMPLO = "FRACCIONAMIENTO DIAMANTE, 39893, ACAPULCO DE JUAREZ, GUERRERO." # Dirección [7]
VELOCIDAD_EJEMPLO = "60" # Velocidad [9]

# --- PASO 1: OBTENER EL TOKEN DE AUTENTICACIÓN (GetUserToken) ---

def get_user_token(user, password, url, headers):
    """
    Realiza la autenticación usando el método GetUserToken [3].
    Devuelve el token de sesión o None si falla.
    """
    # Estructura del Request para GetUserToken, basado en [10]
    soap_request = f"""
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:loc="http://xmlns.landstar.mx/compass/location-receiver/">
   <soapenv:Header/>
   <soapenv:Body>
      <loc:GetUserTokenRequest>
         <loc:userId>{user}</loc:userId>
         <loc:password>{password}</loc:password>
      </loc:GetUserTokenRequest>
   </soapenv:Body>
</soapenv:Envelope>
    """
    print("1. Intentando obtener el Token de Autenticación...")
    
    try:
        response = requests.post(url, data=soap_request, headers=headers)
        
        if response.status_code == 200:
            print("   Respuesta HTTP exitosa (Status 200).")
            # En una implementación real, se usaría un parser XML (como ElementTree)
            # para extraer el token de manera segura desde la respuesta [11].
            
            # Buscando el token directamente en el texto de la respuesta (simple demo)
            if "<xs:token" in response.text:
                start_index = response.text.find("<xs:token>") + len("<xs:token>")
                end_index = response.text.find("</xs:token>")
                token = response.text[start_index:end_index]
                print(f"   Token obtenido exitosamente. (Inicio: {token[:10]}...)")
                return token
            else:
                print("   Error: No se encontró el tag <xs:token> en la respuesta.")
                print("   Respuesta completa (ERROR DE AUTENTICACIÓN):", response.text) # Puede ser un error como el 2004 [4, 12]
                return None
        else:
            print(f"   Error HTTP en autenticación: {response.status_code}")
            return None
            
    except requests.exceptions.RequestException as e:
        print(f"   Error de conexión durante la autenticación: {e}")
        return None

# --- PASO 2: ENVIAR DATOS DE GEOLOCALIZACIÓN (GPSAssetTracking) ---

def send_gps_tracking(token, event_data, url, headers):
    """
    Envía los datos de geolocalización usando el método GPSAssetTracking [4, 6].
    """
    print("\n2. Intentando enviar los datos de Geolocalización (GPSAssetTracking)...")

    # Estructura del Request para GPSAssetTracking, basado en el ejemplo Minimal [6]
    # Se utiliza un token dinámico.
    soap_request = f"""
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:loc="http://xmlns.landstar.mx/compass/location-receiver/">
   <soapenv:Header/>
   <soapenv:Body>
      <loc:GPSAssetTrackingRequest>
         <loc:token>{token}</loc:token>
         <loc:events>
            <loc:Event>
               <loc:code>{event_data['code']}</loc:code>
               <loc:date>{event_data['date']}</loc:date>
               <loc:latitude>{event_data['latitude']}</loc:latitude>
               <loc:longitude>{event_data['longitude']}</loc:longitude>
               <loc:asset>{event_data['asset']}</loc:asset>
               <loc:direction>{event_data['direction']}</loc:direction>
               <loc:speed>{event_data['speed']}</loc:speed>
            </loc:Event>
         </loc:events>
      </loc:GPSAssetTrackingRequest>
   </soapenv:Body>
</soapenv:Envelope>
    """
    
    try:
        response = requests.post(url, data=soap_request, headers=headers)
        
        if response.status_code == 200:
            print("   Respuesta HTTP exitosa (Status 200).")
            # Verificando si se obtuvo un ID de trabajo (idJob), como en [13] y [14]
            if "<ns2:idJob>" in response.text:
                print("   Transacción de geolocalización exitosa (Se recibió un idJob).")
            else:
                print("   Respuesta recibida (Verificar errores):", response.text)
        else:
            print(f"   Error HTTP al enviar tracking: {response.status_code}")
            
    except requests.exceptions.RequestException as e:
        print(f"   Error de conexión durante el tracking: {e}")


# --- EJECUCIÓN DEL FLUJO DE PRUEBA ---

def main():
    if USUARIO_PRUEBA == "SU_USUARIO_AQUI" or CONTRASENA_PRUEBA == "SU_CONTRASENA_AQUI":
        print("ERROR: Debe reemplazar 'SU_USUARIO_AQUI' y 'SU_CONTRASENA_AQUI' en el script con las credenciales proporcionadas por Landstar Metro antes de ejecutar la prueba.")
        return

    # Datos mínimos requeridos para el GPS Asset Tracking [6]
    event_data = {
        'code': CODIGO_EVENTO,
        'date': FECHA_EJEMPLO,
        'latitude': LATITUD_EJEMPLO,
        'longitude': LONGITUD_EJEMPLO,
        'asset': ASSET_PLACA_EJEMPLO,
        'direction': DIRECCION_EJEMPLO,
        'speed': VELOCIDAD_EJEMPLO
    }
    
    # PASO 1: Obtener el token
    auth_token = get_user_token(USUARIO_PRUEBA, CONTRASENA_PRUEBA, URL_PRUEBAS, HEADERS)
    
    if auth_token:
        # PASO 2: Enviar los datos de geolocalización usando el token
        send_gps_tracking(auth_token, event_data, URL_PRUEBAS, HEADERS)
    else:
        print("\nPrueba fallida. No se pudo obtener el token. Revise las credenciales o la URL del servicio.")

if __name__ == "__main__":
    # Importante: El contacto en Landstar Metro le proporcionará un usuario y contraseña para el ambiente de pruebas [5].
    main()
