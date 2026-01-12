#!/usr/bin/env python3
"""
SkyWave API JSON Retriever

This program retrieves data from the SkyWave satellite API
and prints the response to the console.

Based on the knowhow documentation for SkyWave to MVT366 integration.
"""

import requests
import json
import sys
from datetime import datetime

# SkyWave API Configuration
SKYWAVE_CONFIG = {
    'access_id': '70001184',
    'password': 'JEUTPKKH',
    'from_id': '13969586728'
}

# API Endpoint
BASE_URL = 'https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml'

def build_api_url():
    """Build the complete API URL with authentication parameters."""
    params = [
        f'access_id={SKYWAVE_CONFIG["access_id"]}',
        f'password={SKYWAVE_CONFIG["password"]}',
        f'from_id={SKYWAVE_CONFIG["from_id"]}'
    ]

    url = f'{BASE_URL}/?{"&".join(params)}'
    return url

def make_api_request():
    """Make the API request and return the response."""
    url = build_api_url()

    print("=== SkyWave API Request ===")
    print(f"Timestamp: {datetime.now().isoformat()}")
    print(f"URL: {url}")
    print(f"Access ID: {SKYWAVE_CONFIG['access_id']}")
    print(f"From ID: {SKYWAVE_CONFIG['from_id']}")
    print("=" * 50)

    try:
        # Make the HTTP GET request
        response = requests.get(url, timeout=30)

        print("=== HTTP Response ===")
        print(f"Status Code: {response.status_code}")
        print(f"Content-Type: {response.headers.get('content-type', 'N/A')}")
        print(f"Content-Length: {len(response.content)} bytes")
        print("=" * 50)

        if response.status_code == 200:
            print("=== Response Content ===")
            # Try to pretty-print if it's XML
            content = response.text
            print(content)

            # Save to file for analysis
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            filename = f"skywave_response_{timestamp}.xml"

            with open(filename, 'w', encoding='utf-8') as f:
                f.write(content)

            print("=" * 50)
            print(f"Response saved to: {filename}")
            return content
        else:
            print(f"❌ API Request Failed with status code: {response.status_code}")
            print(f"Response: {response.text}")
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

    # Make the API request
    response_data = make_api_request()

    if response_data:
        print("\n✅ API request completed successfully!")
        print(f"Retrieved {len(response_data)} characters of data")
    else:
        print("\n❌ API request failed!")
        sys.exit(1)

if __name__ == "__main__":
    main()
