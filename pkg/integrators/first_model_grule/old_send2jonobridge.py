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

    # batt=hex((int(battery*1024)/6)).replace("0","").replace("x","")
    # batt=batt.upper()
    # batt=batt.zfill(4)

    first_output = "," + imei + ",AAA," + eventCode + "," + latitude + "," + longitude + "," + utc + "," + status + "," + str(sats) + "," + str(gsmStrenght) + "," + str(speed) + "," + str(direction) + "," + str(accuracy) + "," + str(altitude) + "," + str(mileage) + "," + str(runtime) + "," + mcc + "|" + mnc + "|" + lac + "|" + cellId + "," + portStatus + "," + AD1 + "|" + AD2 + "|" + AD3 + "|" + str(battery) + "|" + AD5 + "," + eventInfo + ",*"
    totalchar = charcounter(first_output) + 4
    header = "$$" + mydataidentifier + str(totalchar)
    preoutput = header + first_output
    output = preoutput + crc(preoutput) + chr(13) + chr(10)
    return output

def send2platform(server, port, eventCode):
    imei = "864352045580768"
    eventCode = eventCode
    latitude = "19.611106"
    longitude = "-99.028335"
    utc = datetime.utcnow().strftime('%y%m%d%H%M%S')
    status = "A"
    sats = "9"
    gsmStrenght = "12"
    speed = "98"
    direction = "76"
    accuracy = "1"  # 1 Perfect, 2-3 Wonderful 4-6 good 7-8 medium 9-20 below average 21-50 poor
    altitude = "2239"
    mileage = "0"
    runtime = "1348"
    mcc = "0"
    mnc = "0"
    lac = "0000"
    cellId = "0000"
    portStatus = "0000"
    AD1 = "0000"
    AD2 = "0000"
    AD3 = "0000"
    battery = "80"
    AD5 = "0000"
    eventInfo = "00000000"
    output = payload(imei, eventCode, latitude, longitude, utc, status, sats, gsmStrenght, speed, direction, accuracy, altitude, mileage, runtime, mcc, mnc, lac, cellId, portStatus, AD1, AD2, AD3, battery, AD5, eventInfo)
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    print(output)
    s.connect((server, port))
    s.sendall(output.encode('utf-8'))
    s.close

url= "jonobridge.madd.com.mx"
port = 8055
event = "1"
send2platform(url, port, event )
