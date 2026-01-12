#!/usr/bin/env python3
"""
PoC script para DICSA API.
Genera token con las credenciales del manual y envía datos inventados a /data/.
"""
import sys
import requests
import datetime

TOKEN_URL = "https://viatix.com.mx/red_energy/API_GPS/token/"
DATA_URL = "https://viatix.com.mx/red_energy/API_GPS/data/"
USER = "dicuser"
PASSWORD = "dicpwd25"

# Valores inventados (puedes cambiarlos)
PAYLOAD = {
    "placa": "ABC123",
    "neconomico": "NECO-01",
    "latitud": "19.432608",
    "longitud": "-99.133209",
    "altitud": "2240",
    "direccion": "Ciudad de Mexico, Centro",
    "velocidad": "60",
    # formato de ejemplo: "YYYY-MM-DD HH:MM:SS"
    "fecha_localizacion": datetime.datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S"),
}


def get_token(session):
    params = {"user": USER, "password": PASSWORD}
    print(f"Solicitando token en: {TOKEN_URL} con params={params}")
    resp = session.get(TOKEN_URL, params=params, timeout=15)
    resp.raise_for_status()
    data = resp.json()
    token = data.get("token_access") or data.get("token")
    if not token:
        raise RuntimeError(f"Respuesta de token inválida: {data}")
    print(f"Token recibido: {token}")
    return token


def send_data(session, token):
    # En el manual indican enviar token tanto en params como en Header
    params = {"user": USER, "password": PASSWORD, "token": token}
    # añadir los campos inventados
    params.update(PAYLOAD)

    headers = {"Authorization": token}

    print(f"Enviando datos a: {DATA_URL} - params={params} - headers={headers}")
    # El manual muestra un POST, pero acepta parámetros en la URL; usamos POST con params
    resp = session.post(DATA_URL, params=params, headers=headers, timeout=15)
    # imprimir status y cuerpo
    print(f"HTTP {resp.status_code}")
    try:
        print(resp.text)
    except Exception:
        print("No se pudo leer el cuerpo de la respuesta")
    resp.raise_for_status()
    return resp


def main(dry_run=False):
    s = requests.Session()
    try:
        token = get_token(s)
    except Exception as e:
        print(f"Error obteniendo token: {e}")
        return 2

    if dry_run:
        print("DRY RUN - no se enviarán datos. Parámetros preparados:")
        print({**{"user": USER, "password": PASSWORD, "token": token}, **PAYLOAD})
        return 0

    try:
        send_data(s, token)
    except Exception as e:
        print(f"Error enviando datos: {e}")
        return 3

    print("Proceso completado.")
    return 0


if __name__ == '__main__':
    dry = False
    if len(sys.argv) > 1 and sys.argv[1] in ("-n", "--dry-run"):
        dry = True
    sys.exit(main(dry_run=dry))
