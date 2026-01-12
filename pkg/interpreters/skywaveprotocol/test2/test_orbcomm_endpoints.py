#!/usr/bin/env python3
"""
ORBCOMM Endpoint Testing Script
Tests various ORBCOMM API endpoints to find the working one for ST9101 device
"""

import requests
import json
import time
from datetime import datetime, timedelta
import base64
import logging
from urllib.parse import urljoin

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Device and credentials from datos.md
CREDENTIALS = {
    'username': 'rodolfo@gpscontrol.com.mx',
    'password': 'GPSc0ntr0l1*'
}

DEVICE_INFO = {
    'sat_id': '02092247SKY6A70',
    'imei': '353500723153153',
    'model': 'ST9101'
}

def test_endpoint(base_url, endpoint, method='GET', auth_type='basic', params=None, data=None, headers=None):
    """
    Test a specific endpoint with given parameters
    """
    full_url = urljoin(base_url, endpoint)
    logger.info(f"Testing {method} {full_url}")
    
    # Prepare headers
    test_headers = {
        'User-Agent': 'ORBCOMM-Test-Client/1.0',
        'Accept': 'application/json, application/xml, text/xml, */*',
    }
    
    if headers:
        test_headers.update(headers)
    
    # Prepare authentication
    auth = None
    if auth_type == 'basic':
        auth = (CREDENTIALS['username'], CREDENTIALS['password'])
    elif auth_type == 'header':
        # Some APIs use Authorization header
        credentials = f"{CREDENTIALS['username']}:{CREDENTIALS['password']}"
        encoded_credentials = base64.b64encode(credentials.encode()).decode()
        test_headers['Authorization'] = f'Basic {encoded_credentials}'
    
    try:
        session = requests.Session()
        session.headers.update(test_headers)
        
        if method.upper() == 'GET':
            response = session.get(full_url, auth=auth, params=params, timeout=30)
        elif method.upper() == 'POST':
            if 'Content-Type' not in test_headers:
                test_headers['Content-Type'] = 'application/json'
            response = session.post(full_url, auth=auth, params=params, json=data, timeout=30)
        else:
            logger.error(f"Unsupported method: {method}")
            return None
        
        logger.info(f"Response Status: {response.status_code}")
        logger.info(f"Response Headers: {dict(response.headers)}")
        
        if response.status_code == 200:
            content = response.text
            logger.info(f"SUCCESS! Received {len(content)} characters")
            logger.info(f"Response preview: {content[:500]}...")
            
            # Try to parse as JSON
            try:
                json_data = response.json()
                logger.info("Response is valid JSON")
                return {
                    'success': True,
                    'status_code': response.status_code,
                    'content_type': response.headers.get('content-type', ''),
                    'data': json_data,
                    'raw_content': content
                }
            except json.JSONDecodeError:
                logger.info("Response is not JSON, might be XML")
                return {
                    'success': True,
                    'status_code': response.status_code,
                    'content_type': response.headers.get('content-type', ''),
                    'data': None,
                    'raw_content': content
                }
        else:
            logger.warning(f"HTTP {response.status_code}: {response.text}")
            return {
                'success': False,
                'status_code': response.status_code,
                'error': response.text
            }
            
    except requests.exceptions.RequestException as e:
        logger.error(f"Request failed: {e}")
        return {
            'success': False,
            'error': str(e)
        }

def test_orbcomm_endpoints():
    """
    Test various ORBCOMM endpoints with different configurations
    """
    # Time parameters
    end_time = datetime.utcnow()
    start_time = end_time - timedelta(hours=24)
    
    # Base parameters for requests
    base_params = {
        'startUTC': start_time.strftime('%Y-%m-%dT%H:%M:%SZ'),
        'endUTC': end_time.strftime('%Y-%m-%dT%H:%M:%SZ'),
        'mobileID': DEVICE_INFO['sat_id'],
        'formatType': 'JSON'
    }
    
    # Test configurations
    test_configs = [
        # Main ORBCOMM IDP endpoints
        {
            'name': 'ORBCOMM IDP REST JSON (Option 1)',
            'base_url': 'https://isatdatapro.orbcomm.com',
            'endpoint': '/GLGW/2/RestMessages.svc/get_return_messages.json',
            'method': 'GET',
            'auth_type': 'basic',
            'params': base_params
        },
        {
            'name': 'ORBCOMM IDP REST JSON (Option 2)',
            'base_url': 'https://isatdatapro.orbcomm.com',
            'endpoint': '/GLGW/2/RestMessages.svc/JSON/get_return_messages',
            'method': 'GET',
            'auth_type': 'basic',
            'params': base_params
        },
        {
            'name': 'ORBCOMM IDP REST JSON (Custom)',
            'base_url': 'https://isatdatapro.orbcomm.com',
            'endpoint': '/GLGW/2/RestMessages.svc/JSON',
            'method': 'GET',
            'auth_type': 'basic',
            'params': {**base_params, 'no_json': 'true'}
        },
        # Try with POST method
        {
            'name': 'ORBCOMM IDP REST JSON POST',
            'base_url': 'https://isatdatapro.orbcomm.com',
            'endpoint': '/GLGW/2/RestMessages.svc/get_return_messages.json',
            'method': 'POST',
            'auth_type': 'basic',
            'data': base_params
        },
        # Try XML endpoints
        {
            'name': 'ORBCOMM IDP REST XML',
            'base_url': 'https://isatdatapro.orbcomm.com',
            'endpoint': '/GLGW/2/RestMessages.svc/get_return_messages',
            'method': 'GET',
            'auth_type': 'basic',
            'params': {**base_params, 'formatType': 'XML'}
        },
        # Alternative base URLs
        {
            'name': 'Alternative API endpoint 1',
            'base_url': 'https://api.orbcomm.com',
            'endpoint': '/GlobalMessages/GetReturnMessages',
            'method': 'GET',
            'auth_type': 'basic',
            'params': base_params
        },
        {
            'name': 'Alternative API endpoint 2',
            'base_url': 'https://portal.orbcomm.com',
            'endpoint': '/api/v1/messages',
            'method': 'GET',
            'auth_type': 'basic',
            'params': base_params
        },
        # Try with different authentication
        {
            'name': 'ORBCOMM IDP with Header Auth',
            'base_url': 'https://isatdatapro.orbcomm.com',
            'endpoint': '/GLGW/2/RestMessages.svc/get_return_messages.json',
            'method': 'GET',
            'auth_type': 'header',
            'params': base_params
        },
        # Test just the base service endpoint
        {
            'name': 'Base Service Endpoint',
            'base_url': 'https://isatdatapro.orbcomm.com',
            'endpoint': '/GLGW/2/RestMessages.svc',
            'method': 'GET',
            'auth_type': 'basic',
            'params': {}
        }
    ]
    
    successful_configs = []
    
    logger.info("="*80)
    logger.info("ORBCOMM ENDPOINT TESTING")
    logger.info("="*80)
    logger.info(f"Device: {DEVICE_INFO['model']} ({DEVICE_INFO['sat_id']})")
    logger.info(f"Time range: {start_time.isoformat()} to {end_time.isoformat()}")
    logger.info(f"Account: {CREDENTIALS['username']}")
    logger.info("="*80)
    
    for i, config in enumerate(test_configs, 1):
        logger.info(f"\n[{i}/{len(test_configs)}] Testing: {config['name']}")
        logger.info("-" * 60)
        
        result = test_endpoint(
            base_url=config['base_url'],
            endpoint=config['endpoint'],
            method=config['method'],
            auth_type=config['auth_type'],
            params=config.get('params'),
            data=config.get('data'),
            headers=config.get('headers')
        )
        
        if result and result.get('success'):
            successful_configs.append({
                'config': config,
                'result': result
            })
            logger.info("✅ SUCCESS - This endpoint works!")
        else:
            logger.info("❌ FAILED")
        
        # Small delay between requests
        time.sleep(1)
    
    # Summary
    logger.info("\n" + "="*80)
    logger.info("SUMMARY")
    logger.info("="*80)
    
    if successful_configs:
        logger.info(f"✅ Found {len(successful_configs)} working endpoint(s):")
        for i, success in enumerate(successful_configs, 1):
            config = success['config']
            result = success['result']
            logger.info(f"\n{i}. {config['name']}")
            logger.info(f"   URL: {config['base_url']}{config['endpoint']}")
            logger.info(f"   Method: {config['method']}")
            logger.info(f"   Auth: {config['auth_type']}")
            logger.info(f"   Status: {result['status_code']}")
            logger.info(f"   Content-Type: {result['content_type']}")
            
            # Show data sample if available
            if result.get('data'):
                logger.info(f"   Data sample: {str(result['data'])[:200]}...")
            elif result.get('raw_content'):
                logger.info(f"   Content sample: {result['raw_content'][:200]}...")
    else:
        logger.error("❌ No working endpoints found!")
        logger.info("\nTroubleshooting suggestions:")
        logger.info("1. Verify credentials are correct")
        logger.info("2. Check if the device is active and has recent data")
        logger.info("3. Contact ORBCOMM support for correct API endpoints")
        logger.info("4. Try accessing the portal manually: https://partner-support.orbcomm.com")
    
    return successful_configs

def generate_curl_commands(successful_configs):
    """
    Generate curl commands for successful endpoints
    """
    if not successful_configs:
        return
    
    logger.info("\n" + "="*80)
    logger.info("CURL COMMANDS FOR WORKING ENDPOINTS")
    logger.info("="*80)
    
    for i, success in enumerate(successful_configs, 1):
        config = success['config']
        
        base_url = config['base_url']
        endpoint = config['endpoint']
        method = config['method']
        params = config.get('params', {})
        
        # Build curl command
        curl_cmd = f"curl -X {method.upper()}"
        
        # Add authentication
        if config['auth_type'] == 'basic':
            curl_cmd += f" -u '{CREDENTIALS['username']}:{CREDENTIALS['password']}'"
        elif config['auth_type'] == 'header':
            credentials = f"{CREDENTIALS['username']}:{CREDENTIALS['password']}"
            encoded_credentials = base64.b64encode(credentials.encode()).decode()
            curl_cmd += f" -H 'Authorization: Basic {encoded_credentials}'"
        
        # Add headers
        curl_cmd += " -H 'User-Agent: ORBCOMM-Test-Client/1.0'"
        curl_cmd += " -H 'Accept: application/json, application/xml, text/xml, */*'"
        
        # Build URL with parameters
        full_url = f"{base_url}{endpoint}"
        if params and method.upper() == 'GET':
            param_string = "&".join([f"{k}={v}" for k, v in params.items()])
            full_url += f"?{param_string}"
        
        curl_cmd += f" '{full_url}'"
        
        # Add data for POST
        if method.upper() == 'POST' and config.get('data'):
            curl_cmd += f" -d '{json.dumps(config['data'])}'"
            curl_cmd += " -H 'Content-Type: application/json'"
        
        logger.info(f"\n{i}. {config['name']}:")
        logger.info(curl_cmd)

if __name__ == "__main__":
    try:
        successful_configs = test_orbcomm_endpoints()
        generate_curl_commands(successful_configs)
        
    except KeyboardInterrupt:
        logger.info("\nTest interrupted by user")
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        raise
