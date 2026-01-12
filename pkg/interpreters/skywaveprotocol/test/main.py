import requests
import xml.etree.ElementTree as ET
import socket
import time
import os
from datetime import datetime

# Hardcoded credentials
ACCESS_ID = 70001184
PASSWORD = "JEUTPKKH"

headers_list = ["imei", "commandtype", "eventcode", "latitude", "longitude", "datetime", "positionstatus", "numsatellites", "gsmsignal", "speed", "direction", "hdop", "altitude", "mileage", "runtime", "MCC", "MNC", "LAC", "CI", "iostatus", "AD1", "AD2", "AD3", "batteryanalog", "externalpoweranalog", "geofence", "customizeddata", "protocolversion", "fuelpercentage", "temperaturesensor", "maxccelerationvalue", "maxdecelerationvalue"]

class PayloadBridge:
    def __init__(self):
        self.id = 0
        self.message_utc = ""
        self.receive_utc = ""
        self.sin = 0
        self.mobile_id = ""
        self.region_name = ""
        self.ota_message_size = ""
        self.type = ""
        self.min = ""
        self.latitude = ""
        self.longitude = ""
        self.speed = ""
        self.heading = ""
        self.event_time = ""
        self.gps_fix_age = ""

class GetReturnMessagesResult:
    def __init__(self):
        self.error_id = 0
        self.more = False
        self.next_start_utc = ""
        self.next_start_id = 0
        self.messages = []

    def parse_xml(self, data):
        root = ET.fromstring(data)
        error_id_elem = root.find('ErrorID')
        if error_id_elem is not None:
            self.error_id = int(error_id_elem.text)
        more_elem = root.find('More')
        if more_elem is not None:
            self.more = more_elem.text.lower() == 'true'
        next_start_id_elem = root.find('NextStartID')
        if next_start_id_elem is not None:
            self.next_start_id = int(next_start_id_elem.text)
        messages_elem = root.find('Messages')
        if messages_elem is not None:
            for msg in messages_elem:
                pb = PayloadBridge()
                id_elem = msg.find('ID')
                if id_elem is not None:
                    pb.id = int(id_elem.text)
                message_utc_elem = msg.find('MessageUTC')
                if message_utc_elem is not None:
                    pb.message_utc = message_utc_elem.text
                receive_utc_elem = msg.find('ReceiveUTC')
                if receive_utc_elem is not None:
                    pb.receive_utc = receive_utc_elem.text
                sin_elem = msg.find('SIN')
                if sin_elem is not None:
                    pb.sin = int(sin_elem.text)
                mobile_id_elem = msg.find('MobileID')
                if mobile_id_elem is not None:
                    pb.mobile_id = mobile_id_elem.text
                region_name_elem = msg.find('RegionName')
                if region_name_elem is not None:
                    pb.region_name = region_name_elem.text
                ota_message_size_elem = msg.find('OTAMessageSize')
                if ota_message_size_elem is not None:
                    pb.ota_message_size = ota_message_size_elem.text
                payload = msg.find('Payload')
                if payload is not None:
                    pb.type = payload.get('Name', '')
                    pb.min = payload.get('MIN', '')
                    fields = payload.find('Fields')
                    if fields is not None:
                        for field in fields:
                            name = field.get('Name')
                            value = field.get('Value')
                            if name == 'Latitude':
                                pb.latitude = value
                            elif name == 'Longitude':
                                pb.longitude = value
                            elif name == 'Speed':
                                pb.speed = value
                            elif name == 'Heading':
                                pb.heading = value
                            elif name == 'EventTime':
                                pb.event_time = value
                self.messages.append(pb)

    def returned_messages_bridge(self):
        return self.messages

class SkywaveDoc:
    def __init__(self, access_id, password, from_id):
        self.access_id = access_id
        self.password = password
        self.from_id = from_id

    def get_doc(self):
        url = f"https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml/?access_id={self.access_id}&password={self.password}&from_id={self.from_id}"
        print(f"\n\n Datos de los documentos: {self.access_id} {self.password} {self.from_id}")
        resp = requests.get(url)
        if resp.status_code == 200:
            return resp.text
        else:
            raise Exception(f"Response with other status {resp.status_code}")

class MVT366:
    def __init__(self, **kwargs):
        self.dataidentifier = kwargs.get('dataidentifier', 'H')
        self.datalength = kwargs.get('datalength', 0)
        self.imei = kwargs.get('imei', '')
        self.commandtype = kwargs.get('commandtype', 'AAA')
        self.eventcode = kwargs.get('eventcode', 35)
        self.latitude = kwargs.get('latitude', 0.0)
        self.longitude = kwargs.get('longitude', 0.0)
        self.datetime = kwargs.get('datetime', datetime.now())
        self.positionstatus = kwargs.get('positionstatus', True)
        self.numberofsatellites = kwargs.get('numberofsatellites', 0)
        self.gsmsignal = kwargs.get('gsmsignal', 0)
        self.speed = kwargs.get('speed', 0.0)
        self.direction = kwargs.get('direction', 0)
        self.hdop = kwargs.get('hdop', 0.0)
        self.altitude = kwargs.get('altitude', 21.232345)
        self.millage = kwargs.get('mileage', 0.0)
        self.runtime = kwargs.get('runtime', 0)
        self.mcc = kwargs.get('mcc', 0)
        self.mnc = kwargs.get('mnc', 0)
        self.lac = kwargs.get('lac', 0)
        self.ci = kwargs.get('ci', 0)
        self.ioportstatus = kwargs.get('ioportstatus', '')
        self.ad1 = kwargs.get('ad1', 0.0)
        self.ad2 = kwargs.get('ad2', 0.0)
        self.ad3 = kwargs.get('ad3', 0.0)
        self.batteryanalog = kwargs.get('batteryanalog', 0.0)
        self.externalpoweranalog = kwargs.get('externalpoweranalog', 0.0)
        self.geofence = kwargs.get('geofence', '')
        self.customizeddata = kwargs.get('customizeddata', '')
        self.protocolversion = kwargs.get('protocolversion', 3)
        self.fuelpercentage = kwargs.get('fuelpercentage', 0.0)
        self.temperaturesensor = kwargs.get('temperaturesensor', '')
        self.maxacceleration = kwargs.get('maxacceleration', 0)
        self.maxdeceleration = kwargs.get('maxdeceleration', 0)

    def to_map_alternative(self):
        mp = {}
        for header in headers_list:
            if header == "imei":
                mp[header] = self.imei
            elif header == "commandtype":
                mp[header] = self.commandtype
            elif header == "eventcode":
                mp[header] = self.eventcode
            elif header == "latitude":
                mp[header] = self.latitude
            elif header == "longitude":
                mp[header] = self.longitude
            elif header == "datetime":
                mp[header] = self.datetime
            elif header == "positionstatus":
                mp[header] = self.positionstatus
            elif header == "numsatellites":
                mp[header] = self.numberofsatellites
            elif header == "gsmsignal":
                mp[header] = self.gsmsignal
            elif header == "speed":
                mp[header] = self.speed
            elif header == "direction":
                mp[header] = self.direction
            elif header == "hdop":
                mp[header] = self.hdop
            elif header == "altitude":
                mp[header] = self.altitude
            elif header == "mileage":
                mp[header] = self.millage
            elif header == "runtime":
                mp[header] = self.runtime
            elif header == "LAC":
                mp[header] = self.lac
            elif header == "CI":
                mp[header] = self.ci
            elif header == "MCC":
                mp[header] = self.mcc
            elif header == "MNC":
                mp[header] = self.mnc
            elif header == "iostatus":
                mp[header] = self.ioportstatus
            elif header == "geofence":
                mp[header] = self.geofence
            elif header == "customizeddata":
                mp[header] = self.customizeddata
            elif header == "protocolversion":
                mp[header] = self.protocolversion
            elif header == "AD1":
                mp[header] = self.ad1
            elif header == "AD2":
                mp[header] = self.ad2
            elif header == "AD3":
                mp[header] = self.ad3
            elif header == "batteryanalog":
                mp[header] = self.batteryanalog
            elif header == "externalpoweranalog":
                mp[header] = self.externalpoweranalog
            elif header == "fuelpercentage":
                mp[header] = self.fuelpercentage
            elif header == "temperaturesensor":
                mp[header] = self.temperaturesensor
            elif header == "maxccelerationvalue":
                mp[header] = self.maxacceleration
            elif header == "maxdecelerationvalue":
                mp[header] = self.maxdeceleration
            else:
                raise ValueError(f"Header not recognized {header}")
        return mp

    def to_mvt366_message(self):
        mp = self.to_map_alternative()
        last = ""
        initbody = ""
        analoginput = ""
        iobody = ""
        postbody = ""
        basestation = ""
        for header in headers_list:
            value = mp[header]
            encoded = encode_value(header, value)
            if header in ["CI", "LAC", "MCC", "MNC"]:
                if header == "CI":
                    basestation += encoded
                else:
                    basestation += encoded + "|"
                continue
            if header in ["AD1", "AD2", "AD3", "batteryanalog", "externalpoweranalog"]:
                if header == "externalpoweranalog":
                    analoginput += encoded
                else:
                    analoginput += encoded + "|"
                continue
            if header in ["imei", "commandtype", "eventcode", "latitude", "longitude", "datetime", "positionstatus", "numsatellites", "gsmsignal", "speed", "direction", "hdop", "altitude", "mileage", "runtime"]:
                initbody += encoded + ","
                continue
            if header == "iostatus":
                iobody += encoded + ","
            if header in ["geofence", "customizeddata", "protocolversion", "fuelpercentage", "temperaturesensor", "maxccelerationvalue", "maxdecelerationvalue"]:
                postbody += encoded + ","
                continue
            if header == "maxdecelerationvalue":
                last += encoded
                continue
        body = f"{initbody}{analoginput}{iobody}{postbody}{basestation}{last}"
        self.datalength = len(body) + 5
        if self.dataidentifier == "" or self.dataidentifier == "z":
            self.dataidentifier = "H"
        else:
            # Handle the case where dataidentifier is a letter
            try:
                dataid = int(self.dataidentifier, 10)
                dataid += 1
                self.dataidentifier = chr(dataid)
            except ValueError:
                # If it's a letter, get its ASCII code, increment, and convert back to char
                dataid = ord(self.dataidentifier)
                dataid += 1
                self.dataidentifier = chr(dataid)
        header_str = f"$${self.dataidentifier}{self.datalength}"
        checksum = len(body + header_str) + 2
        lastpart = f"{checksum:X}\r\n"
        return f"{header_str},{body}*{lastpart}"

def encode_value(header, value):
    if header == "imei":
        return str(value)
    elif header == "commandtype":
        return str(value)
    elif header == "eventcode":
        return str(int(value))
    elif header == "latitude":
        return f"{float(value):.06f}"
    elif header == "longitude":
        return f"{float(value):.06f}"
    elif header == "datetime":
        y = value.year
        ys = f"{y:02d}"
        year = ys[2:]
        dt_str = f"{year}{value.month:02d}{value.day:02d}{value.hour:02d}{value.minute:02d}{value.second:02d}"
        return dt_str
    elif header == "positionstatus":
        return "A" if value else "V"
    elif header == "numsatellites":
        return str(int(value))
    elif header == "gsmsignal":
        return str(int(value))
    elif header == "speed":
        return f"{float(value):f}"
    elif header == "direction":
        return str(int(value))
    elif header == "hdop":
        return f"{float(value):f}"
    elif header == "altitude":
        return f"{float(value):f}"
    elif header == "mileage":
        return f"{float(value):f}"
    elif header == "runtime":
        return str(int(value))
    elif header == "LAC":
        return f"{int(value):X}"
    elif header == "CI":
        return f"{int(value):X}"
    elif header == "MCC":
        return str(int(value))
    elif header == "MNC":
        return str(int(value))
    elif header == "iostatus":
        return str(value)
    elif header == "AD1":
        val = (float(value) * 1024) / 6
        vals = str(int(val))
        return f"{int(vals):04X}"
    elif header == "AD2":
        val = (float(value) * 1024) / 6
        vals = str(int(val))
        return f"{int(vals):04X}"
    elif header == "AD3":
        val = (float(value) * 1024) / 6
        vals = str(int(val))
        return f"{int(vals):04X}"
    elif header == "batteryanalog":
        val = (float(value) * 1024) / 6
        vals = str(int(val))
        return f"{int(vals):04X}"
    elif header == "externalpoweranalog":
        val = (float(value) * 1024) / 6
        vals = str(int(val))
        return f"{int(vals):04X}"
    elif header == "geofence":
        return str(value)
    elif header == "customizeddata":
        return str(value)
    elif header == "protocolversion":
        return str(int(value))
    elif header == "fuelpercentage":
        vali = int(float(value))
        vald = f"{float(value):.2f}"
        return f"{vali:X}{vald}"
    elif header == "temperaturesensor":
        return str(value)
    elif header == "maxccelerationvalue":
        return str(int(value))
    elif header == "maxdecelerationvalue":
        return str(int(value))
    else:
        raise ValueError(f"Header not recognized {header}")

def from_bridge_payload(pb):
    if len(pb.latitude) >= 7 and len(pb.longitude) >= 7:
        # Parse lat
        if len(pb.latitude) == 7:
            lat_degrees = pb.latitude[:4]
            lat_decimal = pb.latitude[4:]
        else:
            lat_degrees = pb.latitude[:5]
            lat_decimal = pb.latitude[5:]
        lat_degrees_float = float(lat_degrees)
        lat_degrees_res = lat_degrees_float / 60
        lat_part_one = f"{lat_degrees_res:.2f}"
        lat_detwodec = float(lat_part_one)
        lat_decimal_float = float(lat_decimal)
        lat_decimal_res = lat_decimal_float / 60000
        if '-' in lat_degrees:
            lat_decimal_res *= -1
        lat = lat_detwodec + lat_decimal_res
        lat -= 0.003333
        # Parse lon
        if len(pb.longitude) == 7:
            lon_degrees = pb.longitude[:4]
            lon_decimal = pb.longitude[4:]
        else:
            lon_degrees = pb.longitude[:5]
            lon_decimal = pb.longitude[5:]
        lon_degrees_float = float(lon_degrees)
        lon_degrees_res = lon_degrees_float / 60
        lon_part_one = f"{lon_degrees_res:.2f}"
        lon_detwodec = float(lon_part_one)
        lon_decimal_float = float(lon_decimal)
        lon_decimal_res = lon_decimal_float / 60000
        if '-' in lon_degrees:
            lon_decimal_res *= -1
        lon = lon_detwodec + lon_decimal_res
        # Date
        date_partial = pb.receive_utc.replace(' ', 'T') + 'Z'
        dt = datetime.strptime(date_partial, '%Y-%m-%dT%H:%M:%SZ')
        # Speed
        speed = float(pb.speed)
        # Heading
        direction = int(pb.heading)
        return MVT366(
            imei=pb.mobile_id,
            latitude=lat,
            longitude=lon,
            speed=speed,
            direction=direction,
            datetime=dt,
            positionstatus=True,
            protocolversion=3,
            eventcode=35,
            commandtype="AAA",
            altitude=21.232345
        )
    else:
        raise ValueError("Latitude or longitude too short")

def read_since(doc):
    print("ReadSince")
    while True:
        try:
            d = doc.get_doc()
            sky = GetReturnMessagesResult()
            sky.parse_xml(d)
            messages = sky.returned_messages_bridge()
            for message in messages:
                try:
                    conn = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
                    t366 = from_bridge_payload(message)
                    mes = t366.to_mvt366_message()
                    conn.sendto(mes.encode(), ("13.89.38.9", 1805))
                    conn.close()
                except Exception as e:
                    print(e)
                    print("NO HAY CONEXION CON skywave.FromBridgePayload or send")
                    continue
        except Exception as e:
            print(e)
            continue
        if sky.more:
            doc = SkywaveDoc(doc.access_id, doc.password, sky.next_start_id)
            print("Next doc", doc.from_id)
        else:
            print("End document", doc.from_id)
            return doc

if __name__ == "__main__":
    fromid = os.getenv("FROMIDSKYWAVE")
    if not fromid:
        print("osvariable FROMIDSKYWAVE doesn't exists or empty")
        exit(1)
    fromid_uint = int(fromid)
    doc = SkywaveDoc(ACCESS_ID, PASSWORD, fromid_uint)
    lastdoc = read_since(doc)
    while True:
        lastdoc2 = read_since(lastdoc)
        time.sleep(180)  # 3 minutes
        lastdoc = read_since(lastdoc2)
