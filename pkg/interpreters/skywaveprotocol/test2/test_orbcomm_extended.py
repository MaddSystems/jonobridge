#!/usr/bin/env python3
"""
ORBCOMM Data Testing Script - Extended Search
Tests different time ranges and Mobile ID variants to find device data
"""

import requests
import json
from datetime import datetime, timedelta
import logging

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Credentials
USERNAME = 'rodolfo@gpscontrol.com.mx'
PASSWORD = 'GPSc0ntr0l1*'
BASE_URL = 'https://isatdatapro.orbcomm.com'
ENDPOINT = '/GLGW/2/RestMessages.svc/JSON/get_return_messages'

def test_orbcomm_with_variants():
    """Test ORBCOMM API with different Mobile ID variants and time ranges"""
    
    # Different Mobile ID variants from datos.md
    mobile_id_variants = [
        '02092247SKY6A70',    # From datos.md SAT ID
        '02092234SKYB62F',    # From test_orbcomm.md 
        '353500723153153',    # IMEI from datos.md
        '353500723152270',    # IMEI from test_orbcomm.md
    ]
    
    # Different time ranges to test
    time_ranges = [
        {'hours': 1, 'name': '1 hour'},
        {'hours': 6, 'name': '6 hours'},
        {'hours': 24, 'name': '24 hours'},
        {'hours': 72, 'name': '3 days'},
        {'hours': 168, 'name': '1 week'},
        {'hours': 720, 'name': '30 days'},
    ]
    
    session = requests.Session()
    session.auth = (USERNAME, PASSWORD)
    session.headers.update({
        'User-Agent': 'ORBCOMM-Test-Client/1.0',
        'Accept': 'application/json',
    })
    
    logger.info("Testing ORBCOMM API with different Mobile IDs and time ranges")
    logger.info("=" * 80)
    
    successful_requests = []
    
    for mobile_id in mobile_id_variants:
        logger.info(f"\nTesting Mobile ID: {mobile_id}")
        logger.info("-" * 60)
        
        for time_range in time_ranges:
            end_time = datetime.utcnow()
            start_time = end_time - timedelta(hours=time_range['hours'])
            
            params = {
                'startUTC': start_time.strftime('%Y-%m-%dT%H:%M:%SZ'),
                'endUTC': end_time.strftime('%Y-%m-%dT%H:%M:%SZ'),
                'mobileID': mobile_id,
                'formatType': 'JSON'
            }
            
            try:
                url = f"{BASE_URL}{ENDPOINT}"
                response = session.get(url, params=params, timeout=30)
                
                if response.status_code == 200:
                    data = response.json()
                    error_id = data.get('ErrorID', 0)
                    messages = data.get('Messages')
                    
                    logger.info(f"  {time_range['name']:10} - Status: {response.status_code}, ErrorID: {error_id}")
                    
                    if error_id == 0 and messages:
                        logger.info(f"    ✅ SUCCESS! Found {len(messages)} messages")
                        successful_requests.append({
                            'mobile_id': mobile_id,
                            'time_range': time_range['name'],
                            'start_time': start_time.isoformat(),
                            'end_time': end_time.isoformat(),
                            'messages': messages,
                            'full_response': data
                        })
                        
                        # Log first message details if available
                        if messages and len(messages) > 0:
                            first_msg = messages[0]
                            logger.info(f"    First message preview: {str(first_msg)[:200]}...")
                            
                    elif error_id == 21786:
                        logger.info(f"    ⚠️  No messages in this time range")
                    else:
                        logger.info(f"    ❌ Error: {error_id}")
                        
                else:
                    logger.error(f"  {time_range['name']:10} - HTTP {response.status_code}")
                    
            except Exception as e:
                logger.error(f"  {time_range['name']:10} - Exception: {e}")
    
    # Summary
    logger.info("\n" + "=" * 80)
    logger.info("SUMMARY")
    logger.info("=" * 80)
    
    if successful_requests:
        logger.info(f"✅ Found {len(successful_requests)} successful request(s) with data:")
        for i, req in enumerate(successful_requests, 1):
            logger.info(f"\n{i}. Mobile ID: {req['mobile_id']}")
            logger.info(f"   Time Range: {req['time_range']}")
            logger.info(f"   Period: {req['start_time']} to {req['end_time']}")
            logger.info(f"   Messages: {len(req['messages'])}")
            
            # Show sample data
            if req['messages']:
                sample_msg = req['messages'][0]
                logger.info(f"   Sample message: {json.dumps(sample_msg, indent=4)}")
    else:
        logger.warning("❌ No messages found for any Mobile ID or time range")
        logger.info("\nPossible reasons:")
        logger.info("1. Device has not transmitted any data recently")
        logger.info("2. Mobile ID format is incorrect")
        logger.info("3. Device may be inactive or offline")
        logger.info("4. Account may not have access to this device")
        
        logger.info("\nNext steps:")
        logger.info("1. Check ORBCOMM portal for device status")
        logger.info("2. Verify device is powered and has satellite connectivity")
        logger.info("3. Contact ORBCOMM support to verify device configuration")
    
    return successful_requests

def test_additional_endpoints():
    """Test additional ORBCOMM endpoints that might work"""
    
    session = requests.Session()
    session.auth = (USERNAME, PASSWORD)
    session.headers.update({
        'User-Agent': 'ORBCOMM-Test-Client/1.0',
        'Accept': 'application/json, application/xml',
    })
    
    # Additional endpoints to try
    additional_endpoints = [
        '/GLGW/2/RestMessages.svc/JSON/get_forward_messages',
        '/GLGW/2/RestMessages.svc/JSON/get_forward_statuses',
        '/GLGW/2/RestMessages.svc/XML/get_return_messages',
        '/GLGW/2/RestMessages.svc/get_return_messages',
        '/GLGW/2/RestMessages.svc',  # Base service endpoint
    ]
    
    logger.info("\n" + "=" * 80)
    logger.info("TESTING ADDITIONAL ENDPOINTS")
    logger.info("=" * 80)
    
    for endpoint in additional_endpoints:
        url = f"{BASE_URL}{endpoint}"
        logger.info(f"\nTesting: {url}")
        
        try:
            response = session.get(url, timeout=30)
            logger.info(f"Status: {response.status_code}")
            logger.info(f"Content-Type: {response.headers.get('Content-Type', 'Unknown')}")
            
            if response.status_code == 200:
                content = response.text[:500]
                logger.info(f"Content preview: {content}...")
            else:
                logger.info(f"Error: {response.text[:200]}...")
                
        except Exception as e:
            logger.error(f"Exception: {e}")

if __name__ == "__main__":
    try:
        successful_requests = test_orbcomm_with_variants()
        test_additional_endpoints()
        
        if successful_requests:
            logger.info(f"\n🎉 Found working configuration(s)! Use Mobile ID and time range from successful requests above.")
        else:
            logger.info(f"\n📝 No data found, but API endpoint is working. Check device status or try different parameters.")
            
    except KeyboardInterrupt:
        logger.info("\nTest interrupted by user")
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        raise
