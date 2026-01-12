#!/usr/bin/env python3
"""
Test program to validate NSTECH OAuth authentication and test an events API call.
This script attempts to get an OAuth token using the provided credentials and then
sends a sample event.
"""
import os
import sys
import json
import requests
from datetime import datetime, timezone

def get_oauth_token(token_url, client_id, client_secret):
    """Get OAuth token using client credentials flow"""
    print(f"[{datetime.now()}] Requesting OAuth token...")
    
    # Prepare the form data for OAuth request
    data = {
        'client_id': client_id,
        'client_secret': client_secret,
        'grant_type': 'client_credentials'
    }
    
    # Prepare headers
    headers = {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Accept': 'application/json'
    }
    
    try:
        # Make the POST request to get the token
        resp = requests.post(token_url, data=data, headers=headers, timeout=15)
        
        # Print response status and headers
        print(f"HTTP Status: {resp.status_code}")
        print(f"Response Headers: {json.dumps(dict(resp.headers), indent=2)}")
        
        # Try to parse the JSON response
        try:
            json_resp = resp.json()
            print(f"\nResponse Body: {json.dumps(json_resp, indent=2)}")
            
            # Check if we got an access token
            if 'access_token' in json_resp:
                token = json_resp['access_token']
                expires_in = json_resp.get('expires_in', 'unknown')
                print(f"\n✅ SUCCESS: Received access token!")
                print(f"Token Type: {json_resp.get('token_type', 'bearer')}")
                print(f"Expires In: {expires_in} seconds")
                print(f"Token: {token[:10]}...{token[-10:]} (length: {len(token)})")
                
                # Print the Authorization header format for API calls
                print("\nUse this header for API calls:")
                print(f"Authorization: Bearer {token}")
                return token
            else:
                print(f"\n❌ ERROR: No access token in response: {json_resp}")
                return None
                
        except ValueError:
            print(f"\n❌ ERROR: Invalid JSON response: {resp.text}")
            return None
            
    except requests.exceptions.RequestException as e:
        print(f"\n❌ ERROR: Request failed: {e}")
        return None

def send_event(api_url, token, technology_id, account_id, device_id):
    """Send a test event to the NSTECH API"""
    print(f"\n[{datetime.now()}] Sending event data...")
    
    # Use the events API URL
    events_url = api_url
    
    # Create an event JSON payload based on the exact structure from the example
    event_data = {
        "events": [
            {
                "technology_id": technology_id,            # string UUID
                "account_id": account_id,                 # string UUID
                "date": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z",  # exact format
                "device_id": device_id,                   # string
                "event_type": "PanicButton",              # string - event type from NSTECH
                "latitude": -23.584492,                   # number (float)
                "longitude": -46.828401,                  # number (float)
            }
        ]
    }
    
    # Prepare headers
    headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': f"Bearer {token}"
    }
    
    print(f"Sending to URL: {events_url}")
    print(f"Headers: {json.dumps(headers, indent=2, default=str)}")
    print(f"Payload:")
    # Print payload with exact formatting to ensure JSON structure is correct
    payload_json = json.dumps(event_data, indent=2)
    print(payload_json)
    
    try:
        # Make the POST request
        resp = requests.post(events_url, json=event_data, headers=headers, timeout=15)
        
        # Print request details for debugging
        print(f"Request URL: {resp.request.url}")
        print(f"Request Method: {resp.request.method}")
        print(f"Request Headers: {resp.request.headers}")
        print(f"Request Body: {resp.request.body.decode('utf-8') if hasattr(resp.request, 'body') and resp.request.body else None}")
        
        # Print response status and body
        print(f"HTTP Status: {resp.status_code}")
        
        try:
            json_resp = resp.json() if resp.text else {}
            print(f"Response Body: {json.dumps(json_resp, indent=2)}")
            
            if 200 <= resp.status_code < 300:
                print("✅ SUCCESS: Event data sent successfully!")
                return True
            else:
                print(f"❌ ERROR: Failed to send event data: {resp.status_code}")
                return False
                
        except ValueError:
            print(f"Response Body: {resp.text}")
            return False
            
    except requests.exceptions.RequestException as e:
        print(f"❌ ERROR: Request failed: {e}")
        return False

def main():
    """Main function to test OAuth token acquisition and event sending"""
    # Configuration - using the credentials from your email
    token_url = "https://auth.nstech.com.br/realms/zeus/protocol/openid-connect/token"
    # Events API URL
    api_url = "https://stg.nstech.com.br/zeus/api/integra/v1/events"
    # IDs from the provided information
    technology_id = "52466691-482f-48a0-adfc-a68e776eb966"  # TechnologyId
    account_id = "52a4b1da-8e17-49c5-b490-d98ff1b390e0"     # Account_id
    client_id = "52466691-482f-48a0-adfc-a68e776eb966"
    client_secret = "m6qUGrJ7dEVYTeAPLeV3BRVlEkveZCF8"
    device_id = "867869061346979" # Actual device ID
    
    # Print configuration (hide full secret)
    print(f"[{datetime.now()}] Configuration:")
    print(f"Token URL: {token_url}")
    print(f"API URL: {api_url}")
    print(f"Technology ID: {technology_id}")
    print(f"Account ID: {account_id}")
    print(f"Client ID: {client_id}")
    print(f"Client Secret: {client_secret}")
    print(f"Device ID: {device_id}")
    
    # Step 1: Get the OAuth token
    token = get_oauth_token(token_url, client_id, client_secret)
    if not token:
        return 1
        
    # Step 2: Send a test event using the standard JSON API
    print("\n--- Testing NSTECH Events API ---")
    success = send_event(api_url, token, technology_id, account_id, device_id)
    
    # Return success status
    return 0 if success else 1

if __name__ == "__main__":
    sys.exit(main())
