#!/usr/bin/env python3
"""
Send GPS data JSON to the cortedecorriente endpoint as a GET request with query parameters.

This script takes the provided JSON, decodes any %u00XX sequences (like in POS_ADDRESS),
and sends all fields as query parameters to the endpoint.
"""

import json
import re
import sys
from urllib.parse import quote, urlencode

import requests

# Function to decode %u00XX sequences to Latin-1 characters
def decode_percent_u_sequences(text: str) -> str:
    if not isinstance(text, str):
        return text
    # Replace %u00e1 -> %e1 etc.
    text = re.sub(r'%u00([0-9A-Fa-f]{2})', r'%\1', text)
    # Now decode as Latin-1 bytes
    try:
        # Use unquote_to_bytes to handle %XX
        from urllib.parse import unquote_to_bytes
        b = unquote_to_bytes(text)
        return b.decode('latin-1')
    except Exception:
        # If decode fails, return the replaced string
        return text

# The JSON data as a string
json_data = '''{"APPLICATION_ID":"112","APPLICATION_NAME":"Monitoreo","APP_ID":"112","ASSIGNED_VEHICLE_DESC":"N/S: KMFWB3WR2KU008614\\nPLACAS: U77BAF","ASSIGNED_VEHICLE_ID":"13073","ASSIGNED_VEHICLE_NAME":"U77BAF","ASSIGNED_VEHICLE_USERNAME":"Jhinojosa","CellID":"19224622","DEVICE_IMEI":"864606043913275","DEVICE_PHONE":" 526566639430","EVENT_DURATION":"425","EVENT_TIME":"2025-09-30T18:49:08","ErrorLevel":"1.1","Event start time":"2025-09-30T18:42:00","EventCode":"32","GEOFENCE_DESCRIPTION":"CLIENTE","GEOFENCE_ID":"322067","GEOFENCE_NAME":"WA CARTER","GEOFENCE_TAG_DESCRIPTION":"Puerto","GEOFENCE_TAG_ID":"64488","GEOFENCE_TAG_NAME":"Resguardo/CEDIS/Puerto","GF_ACTION_OUTSIDE":"WA CARTER","GPRMC":"$GPRMC,184905,A,1930.5781,N,09909.4335,W,003.2,261.0,300925,000.0,E*64\\r\\n","Immobilizer":"False","Input3":"False","LAC":"40198","MCC":"334","MNC":"50","Output4":"False","POS_ADDRESS":"J%u00fapiter 80, Norte-Bas%u00edlica de Guadalupe, 07700 Ciudad de M%u00e9xico, Federal District","POS_HEADING":"261","POS_LATITUDE":"19.50964","POS_LONGITUDE":"-99.15722","POS_TIME":"2025-09-30T18:49:05","POS_TIME_LAST_VALID":"2025-09-30T18:49:05","Presencia":"0","RULE_NAME":"CorteDeCorrienteExterna","SIGNAL_BATTERY_VOLTAGE":"0.5875","SIGNAL_GSMSIGNALLEVEL":"6","SIGNAL_IGNITION":"True","SIGNAL_LOW_BATTERY":"False","SIGNAL_POWERON":"False","SIGNAL_SATELLITECOUNT":"7","SIGNAL_SOS":"False","SIGNAL_SPEED":"1.66666666666667","SIGNAL_VOLTAGE":"0.09","USER_DESCRIPTION":"N/S: KMFWB3WR2KU008614\\nPLACAS: U77BAF","USER_NAME":"U77BAF","USER_USERNAME":"Jhinojosa","cmd":"_ExternalNotification","id":"316458499","stage":"start"}'''

def main():
    try:
        # Parse the JSON
        data = json.loads(json_data)
    except json.JSONDecodeError as e:
        print(f"Error parsing JSON: {e}", file=sys.stderr)
        sys.exit(1)

    # Prepare query parameters
    params = {}
    for key, value in data.items():
        if isinstance(value, str):
            # Decode %u00XX sequences
            decoded_value = decode_percent_u_sequences(value)
            params[key] = decoded_value
        else:
            # Convert to string for other types
            params[key] = str(value)

    # Base URL
    base_url = "https://jonobridge.madd.com.mx/cortedecorriente"

    # Send GET request
    try:
        response = requests.get(base_url, params=params, timeout=30)
        print("Final GET URL:", response.request.url)
        print("HTTP Status:", response.status_code)
        print("Response Headers:", dict(response.headers))
        print("Response Body:")
        print(response.text)
    except requests.RequestException as e:
        print(f"Error sending request: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()