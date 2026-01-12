#!/usr/bin/env python3
"""
SkyWave API JSON Retriever

This program retrieves data from the SkyWave satellite API for the last specified minutes
and prints the response to the console.

Based on the knowhow documentation for SkyWave to MVT366 integration.
"""

import requests
import sys
from datetime import datetime, timezone, timedelta
import urllib.parse
import xml.etree.ElementTree as ET

# SkyWave API Configuration - VERIFY THESE VALUES
SKYWAVE_CONFIG = {
    'access_id': '70001184',  # Replace with your valid access_id
    'password': 'JEUTPKKH',   # Replace with your valid password
    'from_id': '13969586728'  # This appears to function as the mobile_id for filtering by terminal
}

# Number of minutes before current UTC time to retrieve records
MINUTES_BACK = 250

# API Endpoint
BASE_URL = 'https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml'

def build_api_url():
    """Build the complete API URL with authentication and time range parameters."""
    # Calculate UTC time range
    now_utc = datetime.now(timezone.utc)
    from_time = (now_utc - timedelta(minutes=MINUTES_BACK)).strftime("%Y-%m-%d %H:%M:%S")
    end_time = now_utc.strftime("%Y-%m-%d %H:%M:%S")

    params = {
        'access_id': SKYWAVE_CONFIG["access_id"],
        'password': SKYWAVE_CONFIG["password"],
        'from_id': SKYWAVE_CONFIG["from_id"],
        'start_utc': from_time,
        'end_utc': end_time  # Explicit end time to ensure range
    }

    # URL encode parameters to handle spaces and special characters
    url = BASE_URL + '?' + urllib.parse.urlencode(params, safe=':')  # Keep ':' unencoded for time
    return url

def parse_error_response(content):
    """Parse XML response to extract error details."""
    try:
        root = ET.fromstring(content)
        error_id = root.find('ErrorID')
        if error_id is not None:
            return error_id.text
        return "Unknown error"
    except ET.ParseError:
        return "Failed to parse XML response"

def make_api_request():
    """Make the API request and return the response."""
    url = build_api_url()

    print("=== SkyWave API Request ===")
    print(f"Timestamp (UTC): {datetime.now(timezone.utc).isoformat()}")
    print(f"URL: {url}")
    print(f"Access ID: {SKYWAVE_CONFIG['access_id']}")
    print(f"From ID (used as mobile filter): {SKYWAVE_CONFIG['from_id']}")
    print(f"Retrieving records from last {MINUTES_BACK} minutes")
    print("=" * 50)

    try:
        # Make the HTTP GET request
        response = requests.get(url, timeout=30)

        print("=== HTTP Response ===")
        print(f"Status Code: {response.status_code}")
        print(f"Content-Type: {response.headers.get('content-type', 'N/A')}")
        print(f"Content-Length: {len(response.content)} bytes")
        print(f"Headers: {response.headers}")
        print("=" * 50)

        print("=== Response Content ===")
        content = response.text
        print(content)
        print("=" * 50)

        if response.status_code == 200:
            # Check for error in response
            error_id = parse_error_response(content)
            if error_id != "0" and error_id is not None:
                print(f"❌ API Error: ErrorID {error_id}")
                if error_id == "513":
                    print("Error 513: Likely due to invalid parameter name or value. Verify 'from_id' is correct for your setup.")
                    print("Note: Documentation suggests 'mobile_id' for terminal filtering, but 'from_id' works in your curl example.")
                return None

            # Save to file for analysis
            timestamp = datetime.now(timezone.utc).strftime("%Y%m%d_%H%M%S")
            filename = f"skywave_response_{timestamp}.xml"

            with open(filename, 'w', encoding='utf-8') as f:
                f.write(content)

            print(f"Response saved to: {filename}")
            return content
        else:
            print(f"❌ API Request Failed with status code: {response.status_code}")
            print(f"Response: {content}")
            return None

    except requests.exceptions.RequestException as e:
        print(f"❌ Request Error: {e}")
        return None
    except Exception as e:
        print(f"❌ Unexpected Error: {e}")
        return None

def main():
    """Main function to run the SkyWave API retriever."""
    print("SkyWave Satellite API Data Retriever")
    print("=" * 50)
    print("⚠️ Please verify SKYWAVE_CONFIG values (access_id, password, from_id) are correct.")
    print("⚠️ 'from_id' is used here as the mobile/terminal filter based on your working curl example.")
    print("⚠️ If Error 513 persists, try replacing 'from_id' with 'mobile_id' or contact support for exact parameters.")
    print("=" * 50)

    # Make the API request
    response_data = make_api_request()

    if response_data:
        print("\n✅ API request completed successfully!")
        print(f"Retrieved {len(response_data)} characters of data")
    else:
        print("\n❌ API request failed! Check credentials, parameters, and time format.")
        sys.exit(1)

if __name__ == "__main__":
    main()