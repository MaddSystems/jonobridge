#!/usr/bin/env python3
"""
Test program to validate NSTECH OAuth authentication.
This script attempts to get an OAuth token using the provided credentials.
"""
import os
import sys
import json
import requests
from datetime import datetime

def main():
    """Main function to test OAuth token acquisition"""
    # Configuration - using the credentials from your email
    token_url = "https://auth.nstech.com.br/realms/zeus/protocol/openid-connect/token"
    client_id = "52466691-482f-48a0-adfc-a68e776eb966"
    client_secret = "m6qUGrJ7dEVYTeAPLeV3BRVlEkveZCF8"
    
    # Print configuration (hide full secret)
    print(f"[{datetime.now()}] Configuration:")
    print(f"Token URL: {token_url}")
    print(f"Client ID: {client_id}")
    print(f"Client Secret: {client_secret[:4]}...{client_secret[-4:]}")
    
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
    
    print(f"\n[{datetime.now()}] Requesting OAuth token...")
    
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
                return 0
            else:
                print(f"\n❌ ERROR: No access token in response")
                return 1
                
        except ValueError:
            print(f"\n❌ ERROR: Invalid JSON response: {resp.text}")
            return 1
            
    except requests.exceptions.RequestException as e:
        print(f"\n❌ ERROR: Request failed: {e}")
        return 1

if __name__ == "__main__":
    sys.exit(main())
