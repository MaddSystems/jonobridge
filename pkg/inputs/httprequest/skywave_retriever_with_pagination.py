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
MINUTES_BACK = 5

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
    """Make the API request and follow pagination to retrieve all pages.

    The SkyWave XML response may include <More>true</More> and a <NextStartID> value
    when there are additional pages. This function will loop until <More> is false
    and will combine all pages into a single saved XML file.
    """
    # We'll preserve the configured mobile filter in 'mobile_filter' and use
    # a separate pagination parameter 'start_id' (if provided by the API via NextStartID).
    mobile_filter = SKYWAVE_CONFIG.get('from_id')

    # Time window fixed for the whole run: now -> now - MINUTES_BACK
    now_utc = datetime.now(timezone.utc)
    from_time = (now_utc - timedelta(minutes=MINUTES_BACK)).strftime("%Y-%m-%d %H:%M:%S")
    end_time = now_utc.strftime("%Y-%m-%d %H:%M:%S")

    print("=== SkyWave API Request (paginated) ===")
    print(f"Timestamp (UTC): {now_utc.isoformat()}")
    print(f"Retrieving records from last {MINUTES_BACK} minutes: {from_time} -> {end_time}")
    print(f"Access ID: {SKYWAVE_CONFIG['access_id']}")
    print(f"Mobile filter (from_id): {mobile_filter}")
    print("=" * 50)

    all_pages = []
    next_start_id = None
    page = 0

    while True:
        page += 1
        # Build params for this page
        params = {
            'access_id': SKYWAVE_CONFIG["access_id"],
            'password': SKYWAVE_CONFIG["password"],
            # include mobile filter as provided
            'from_id': mobile_filter,
            'start_utc': from_time,
            'end_utc': end_time,
        }

        # If we have a pagination token from previous response, include it.
        if next_start_id:
            # Many SkyWave examples use 'start_id' or reuse 'from_id' as pagination param.
            # We'll send it as 'start_id' by default; if your API expects a different name,
            # change this to the documented parameter (e.g., 'from_id' replacement).
            params['start_id'] = next_start_id

        url = BASE_URL + '?' + urllib.parse.urlencode(params, safe=':')
        print(f"Requesting page {page}: {url}")

        try:
            resp = requests.get(url, timeout=30)
        except requests.exceptions.RequestException as e:
            print(f"❌ Request Error on page {page}: {e}")
            return None

        content = resp.text

        print(f"Page {page} status: {resp.status_code}, {len(content)} chars")

        if resp.status_code != 200:
            print(f"❌ API Request Failed on page {page} with status code: {resp.status_code}")
            print(content)
            return None

        # Save this page's content
        all_pages.append(content)

        # Parse pagination fields
        try:
            root = ET.fromstring(content)
            more_el = root.find('More')
            next_id_el = root.find('NextStartID')
            more = (more_el is not None and more_el.text and more_el.text.strip().lower() == 'true')
            next_start_id = next_id_el.text.strip() if next_id_el is not None and next_id_el.text else None
        except ET.ParseError:
            print("⚠️  Failed to parse XML for pagination; stopping after current page.")
            more = False
            next_start_id = None

        # If no more pages, break
        if not more:
            print(f"No more pages after page {page}.")
            break

        # Safety: avoid infinite loops
        if page > 100:
            print("⚠️  Reached page limit (100); stopping to avoid infinite loop")
            break

    # Combine pages into single file (simple concatenation)
    timestamp = datetime.now(timezone.utc).strftime("%Y%m%d_%H%M%S")
    filename = f"skywave_response_{timestamp}.xml"
    with open(filename, 'w', encoding='utf-8') as f:
        for p in all_pages:
            f.write(p)
            f.write('\n')

    print(f"Saved combined response ({len(all_pages)} pages) to: {filename}")
    return "\n".join(all_pages)

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