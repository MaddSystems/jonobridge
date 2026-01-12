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

def send_defcon2_test():
    server = "jonobridge.madd.com.mx"
    port = 8056
    
    # Configuración del test
    imei = "864352045580768"
    packet_interval = 15 # Segundos entre paquetes
    
    # FASE 1: LLENAR BUFFER (10 paquetes VÁLIDOS)
    print("\n=== FASE 1: Llenando Buffer (10 paquetes VÁLIDOS) ===")
    num_valid_packets = 10
    
    lat_start = 19.000001
    lon_start = -99.000001
    increment = 0.000001

    current_time = datetime.now()
    
    # Alinear al próximo segundo exacto
    time.sleep(1 - (current_time.microsecond / 1000000.0)) 
    next_send = datetime.now()

    for i in range(num_valid_packets):
        # Esperar al tiempo de envío
        sleep_time = (next_send - datetime.now()).total_seconds()
        if sleep_time > 0:
            time.sleep(sleep_time)

        # Datos VÁLIDOS (Status A)
        lat = f"{lat_start + i * increment:.6f}"
        lon = f"{lon_start - i * increment:.6f}"
        utc = datetime.utcnow().strftime('%y%m%d%H%M%S')
        
        output = payload(
            imei=imei, 
            eventCode="35", 
            latitude=lat, 
            longitude=lon, 
            utc=utc, 
            status="A",  # A = Valid
            sats="9", 
            gsmStrenght="12", 
            speed="98", 
            direction="76", 
            accuracy="1", 
            altitude="2239", 
            mileage="0", 
            runtime="1348", 
            mcc="0", mnc="0", lac="0000", cellId="0000", 
            portStatus="0000", AD1="0000", AD2="0000", AD3="0000", battery="80", AD5="0000", eventInfo="00000000"
        )

        print(f"[{i+1}/{num_valid_packets}] Sending VALID (Status A) at {datetime.now()}: {output.strip()}")

        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((server, port))
        s.sendall(output.encode('utf-8'))
        s.close()

        next_send += timedelta(seconds=packet_interval)

    # FASE 2: PROVOCAR DEFCON 2 -> 3 (Paquetes INVÁLIDOS > 5 minutos)
    # Necesitamos al menos 5 minutos offline.
    # 5 minutos = 300 segundos.
    # Con intervalo de 15s -> 300/15 = 20 paquetes mínimo.
    # Enviaremos 25 para asegurar.
    
    print("\n=== FASE 2: Provocando DEFCON 2 -> 3 (25 paquetes INVÁLIDOS ~ 6 minutos) ===")
    num_invalid_packets = 25
    
    for i in range(num_invalid_packets):
        sleep_time = (next_send - datetime.now()).total_seconds()
        if sleep_time > 0:
            time.sleep(sleep_time)

        # Datos INVÁLIDOS (Status V)
        utc = datetime.utcnow().strftime('%y%m%d%H%M%S')
        
        output = payload(
            imei=imei, 
            eventCode="35", 
            latitude="0.000000", 
            longitude="0.000000", 
            utc=utc, 
            status="V",  # V = Invalid
            sats="0",    # Sin satélites
            gsmStrenght="12", # Buena señal GSM
            speed="0", 
            direction="0", 
            accuracy="0", 
            altitude="0", 
            mileage="0", 
            runtime="1348", 
            mcc="0", mnc="0", lac="0000", cellId="0000", 
            portStatus="0000", AD1="0000", AD2="0000", AD3="0000", battery="80", AD5="0000", eventInfo="00000000"
        )

        print(f"[{i+1}/{num_invalid_packets}] Sending INVALID (Status V) at {datetime.now()}: {output.strip()}")

        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((server, port))
        s.sendall(output.encode('utf-8'))
        s.close()

        next_send += timedelta(seconds=packet_interval)

if __name__ == "__main__":
    send_defcon2_test()
