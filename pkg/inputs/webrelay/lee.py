"""
Consultar las posiciones en server1

Maria Lacayo
16/05/2025

"""

# Librerias
from datetime import datetime
import requests


# Obtener token del server1
def get_server1_token(ggs_user, ggs_password, appid):
    token = ""
    error = ""

    url = "http://server1.gpscontrol.com.mx/comGPSGate/api/v.1/applications/" + str(appid) + "/tokens"
    r = requests.post(url, json={"username": ggs_user, "password": ggs_password})
    t = r.text[1:69]
    if r.status_code == 200:
        if t != "The user does not have neither _APIRead nor _APIReadWrite privileges. To get an API token, please assign required privileges.":
            j = r.json()
            token = j["token"]
        else:
            error = "Faltan permisos en la cuenta. _APIRead y _APIReadWrite"
    else:
        error = "Respuesta de la plataforma: " + r.text
    return dict(token=token, error=error)

# Limpiar la fecha
def extraer_fecha_hora(datetime_str):
    dt = datetime.strptime(datetime_str.rstrip('Z'), "%Y-%m-%dT%H:%M:%S")

    date_str = dt.strftime("%m-%d-%Y")
    time_str = dt.strftime("%H:%M:%S")

    return date_str, time_str, dt

# Buscar el tracktrip más reciente
def encontrar_mas_reciente(data):
    mas_reciente = max(
        data,
        key=lambda x: datetime.strptime(x["endTrackPoint"]["utc"], "%Y-%m-%dT%H:%M:%SZ")
    )
    return mas_reciente

def main():
    print("Entrando ando")
    # Variables a configurar
    app_id = 424
    ggs_user = "admindesarrollo"
    ggs_password = "GPSc0ntr0l00"
    plates = "115"

    print(f"Obteniendo token para app_id: {app_id}, usuario: {ggs_user}")
    res = get_server1_token(ggs_user,ggs_password,app_id)
    token = res['token']
    error = res['error']

    if error == "":
        print(f"Token obtenido: {token[:10]}...") # Print first 10 chars of token for security
        url = f"https://server1.gpscontrol.com.mx/comGpsGate/api/v.1/applications/{app_id}/users?FromIndex=0&PageSize=1000"
        print(f"Consultando usuarios: {url}")
        response = requests.get(url, headers={'Authorization': token})
        users = response.json()
        print(f"Se encontraron {len(users)} usuarios")
        
        user = None
        for u in users:
            if u['name'] == plates:
                user = u
                print(f"Usuario encontrado: {u['name']}")
                break

        if user:
            print(f"Procesando dispositivos para usuario: {user['name']}")
            for device in user['devices']:
                imei = device['imei']
                plates = plates
                altitude = user['trackPoint']['position']['altitude']
                latitude = user['trackPoint']['position']['latitude']
                longitude = user['trackPoint']['position']['longitude']
                speed = user['trackPoint']['velocity']['groundSpeed']*3.6
                heading = user['trackPoint']['velocity']['heading']
                date, time, dt = extraer_fecha_hora(user['deviceActivity'])
                ignition_status = False
                stoping_time = ""
                stopping_date = ""
                
                print(f"IMEI: {imei}, Posición: {latitude}, {longitude}, Velocidad: {speed}")

                # Consultar si esta prendido o apagado el motor
                url = f"https://server1.gpscontrol.com.mx/comGpsGate/api/v.1/applications/{app_id}/users/{user['id']}/status"
                response = requests.get(url, headers={'Authorization': token})
                status = response.json()
                for variable in status['variables']:
                    if variable['name'] == "Ignition":
                        ignition_status = bool(variable['value'])
                        print(f"Estado de ignición: {ignition_status}")
                        break

                # Consultar el último viaje de la unidad
                dt = dt.strftime("%Y-%m-%d")
                url = f"https://server1.gpscontrol.com.mx/comGpsGate/api/v.1/applications/{app_id}/users/{user['id']}/tripinfos?Date={dt}"
                print(f"Consultando viajes: {url}")
                response = requests.get(url, headers={'Authorization': token})
                tripinfos = response.json()
                
                if not tripinfos:
                    print("No se encontraron viajes para la fecha indicada")
                    data = {
                        "imei": imei,
                        "plate": plates,
                        "altitude": altitude,     
                        "latitude": latitude,
                        "longitude": longitude,
                        "speed": speed,
                        "heading": heading,
                        "date": date,
                        "time": time,
                        "moving": False,
                        "ignitionStatus": ignition_status,
                        "stoppingDate": date,
                        "stopingTime": "00:00"
                    }
                    print("Datos obtenidos:", data)
                    return data
                    
                print(f"Se encontraron {len(tripinfos)} viajes")
                tripinfo = encontrar_mas_reciente(tripinfos)

                if tripinfo["totalDistance"] == 0:
                    moving = False
                    stopping_date = date

                    start_time = datetime.strptime(tripinfo["startTrackPoint"]["utc"], "%Y-%m-%dT%H:%M:%SZ")
                    end_time = datetime.strptime(tripinfo["endTrackPoint"]["utc"], "%Y-%m-%dT%H:%M:%SZ")

                    total_segundos = int((end_time - start_time).total_seconds())
                    horas, minutos = divmod(total_segundos // 60, 60)
                    stoping_time = f"{horas:02}:{minutos:02}"
                else:
                    moving = True

                data = {
                    "imei": imei,
                    "plate": plates,
                    "altitude": altitude,     
                    "latitude": latitude,
                    "longitude": longitude,
                    "speed": speed,
                    "heading": heading,
                    "date": date,
                    "time": time,
                    "moving": moving,
                    "ignitionStatus": ignition_status,
                    "stoppingDate": stopping_date,
                    "stopingTime": stoping_time
                }
                
                print("Datos obtenidos:", data)
                return data
        else:
            error_msg = "Las placas solicitadas no existen"
            print(f"Error: {error_msg}")
            return dict(error=error_msg)
    else:
        print(f"Error obteniendo token: {error}")
        return dict(error=f"{error}")

if __name__ == "__main__":
    result = main()
    print("\n===== RESULTADO FINAL =====")
    print(result)