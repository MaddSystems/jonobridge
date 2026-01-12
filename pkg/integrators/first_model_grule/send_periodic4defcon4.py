import socket
from socket import SOL_SOCKET, SO_REUSEADDR
import sys
from datetime import datetime, timedelta
import time

def crc(source):
    b = 0
    for i in source:
        b = b + ord(i)
    ret = hex(b % 256)
    ret = ret.upper()
    ret = ret.replace("0X", "")
    return ret

def charcounter(source):
    c = 0
    for i in source:
        c = c + 1
    return c

class identifier:
    idCounter = 64
    def __init__(self):
        if identifier.idCounter < 123:
            identifier.idCounter += 1
        else:
            identifier.idCounter = 65

def payload(imei, eventCode, latitude, longitude, utc, status, sats, gsmStrenght, speed, direction, accuracy, altitude, mileage, runtime, mcc, mnc, lac, cellId, portStatus, AD1, AD2, AD3, battery, AD5, eventInfo):
    # MVT380
    imei = imei.strip()
    newIdentifier = identifier()
    mydataidentifier = str(chr(newIdentifier.idCounter))

    first_output = "," + imei + ",AAA," + eventCode + "," + latitude + "," + longitude + "," + utc + "," + status + "," + str(sats) + "," + str(gsmStrenght) + "," + str(speed) + "," + str(direction) + "," + str(accuracy) + "," + str(altitude) + "," + str(mileage) + "," + str(runtime) + "," + mcc + "|" + mnc + "|" + lac + "|" + cellId + "," + portStatus + "," + AD1 + "|" + AD2 + "|" + AD3 + "|" + str(battery) + "|" + AD5 + "," + eventInfo + ",*"
    totalchar = charcounter(first_output) + 4
    header = "$$" + mydataidentifier + str(totalchar)
    preoutput = header + first_output
    output = preoutput + crc(preoutput) + chr(13) + chr(10)
    return output

def send_defcon4_test():
    server = "jonobridge.madd.com.mx"
    port = 8056
    
    # Configuración del test
    imei = "864352045580768"
    packet_interval = 15 # Segundos entre paquetes
    
    next_send = datetime.now()  # Initialize next_send
    
    # FASE 1: LLENAR EL BUFFER (CALENTAMIENTO)
    # Necesitamos al menos 10 paquetes válidos para activar BufferHas10
    print("\n=== FASE 1: LLENAR EL BUFFER (CALENTAMIENTO) ===")
    print("Enviando 11 paquetes VÁLIDOS para activar B10...")
    num_warmup_packets = 11
    
    # Ubicación inicial para calentamiento (incrementando)
    lat_start = 19.000001
    lon_start = -99.00001
    
    for i in range(num_warmup_packets):
        # Esperar al tiempo de envío (rápido para calentamiento)
        sleep_time = (next_send - datetime.now()).total_seconds()
        if sleep_time > 0:
            time.sleep(sleep_time)

        # Incrementar coordenadas para cada paquete
        lat = f"{lat_start + i * 0.000001:.6f}"
        lon = f"{lon_start - i * 0.00001:.5f}"
        
        # Datos VÁLIDOS (Status A)
        utc = datetime.utcnow().strftime('%y%m%d%H%M%S')
        
        print(f"Frame {i+1}: Coordinates - lat={lat}, lon={lon}")
        
        # Speed 50 km/h (> 20), GSM 20 (> 9)
        output = payload(
            imei=imei, 
            eventCode="35", 
            latitude=lat, 
            longitude=lon, 
            utc=utc, 
            status="A",  # A = Valid
            sats="12", 
            gsmStrenght="20", # Strong Signal
            speed="50",       # High Speed
            direction="180", 
            accuracy="1", 
            altitude="100", 
            mileage="0", 
            runtime="1000", 
            mcc="334", mnc="020", lac="1234", cellId="5678", 
            portStatus="0000", AD1="0000", AD2="0000", AD3="0000", battery="100", AD5="0000", eventInfo="00000000"
        )

        print(f"[{i+1}/{num_warmup_packets}] Sending WARMUP VALID (Status A) at {datetime.now()}: {output.strip()}")

        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((server, port))
        s.sendall(output.encode('utf-8'))
        s.close()

        next_send += timedelta(seconds=1)  # Enviar cada 1 segundo para calentamiento rápido

    # FASE 2: PROVOCAR DEFCON 4 (Paquetes INVÁLIDOS > 5 minutos)
    # Necesitamos al menos 5 minutos offline.
    # 5 minutos = 300 segundos.
    # Con intervalo de 15s -> 300/15 = 20 paquetes mínimo.
    # Enviaremos 24 para asegurar que pasamos el umbral de 6 minutos (360 segundos).
    
    print("\n=== FASE 2: Provocando DEFCON 4 (24 paquetes INVÁLIDOS ~ 6 minutos, GSM bajado a 10 para cambio de señal) ===")
    print("Iniciando simulación de Jammer...")
    num_invalid_packets = 24
    
    for i in range(num_invalid_packets):
        sleep_time = (next_send - datetime.now()).total_seconds()
        if sleep_time > 0:
            time.sleep(sleep_time)

        # Datos INVÁLIDOS (Status V)
        # La regla revisa "IsOfflineFor(5)". 
        # Esto cuenta el tiempo desde el ÚLTIMO paquete válido.
        # Al enviar paquetes inválidos, el contador aumenta.
        
        utc = datetime.utcnow().strftime('%y%m%d%H%M%S')
        #19.52073, -99.21152
        output = payload(
            imei=imei, 
            eventCode="35", 
            latitude="19.520730", 
            longitude="-99.211520", 
            utc=utc, 
            status="V",  # V = Invalid
            sats="0",    # Sin satélites
            gsmStrenght="10", # Baja señal GSM (cambiado de 15 a 10, < 15) para activar condición de jammer
            speed="100", # Speed 100 to maintain 'AvgSpeed >= 10' condition in buffer
            direction="0", 
            accuracy="0", 
            altitude="0", 
            mileage="0", 
            runtime="1348", 
            mcc="334", mnc="020", lac="1234", cellId="5678", 
            portStatus="0000", AD1="0000", AD2="0000", AD3="0000", battery="100", AD5="0000", eventInfo="00000000"
        )

        print(f"[{i+1}/{num_invalid_packets}] Sending INVALID (Status V) at {datetime.now()}: {output.strip()}")

        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((server, port))
        s.sendall(output.encode('utf-8'))
        s.close()

        next_send += timedelta(seconds=packet_interval)

if __name__ == "__main__":
    send_defcon4_test()
