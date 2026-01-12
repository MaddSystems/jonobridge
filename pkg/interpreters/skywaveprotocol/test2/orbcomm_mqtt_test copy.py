#!/usr/bin/env python3
"""
ORBCOMM Data Retrieval and MQTT Publisher
Retrieves data from ORBCOMM provider for ST9101 device and sends via MQTT
"""

import requests
import json
import time
import xml.etree.ElementTree as ET
from datetime import datetime, timedelta
import paho.mqtt.client as mqtt
import argparse
import logging
import os
from typing import Optional, List, Dict

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class ORBCOMMClient:
    """Client for ORBCOMM API communication"""
    
    def __init__(self, access_id: str, password: str, from_id: str, base_url: str = "https://isatdatapro.skywave.com"):
        self.access_id = access_id
        self.password = password
        self.from_id = from_id
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        
    def authenticate(self) -> bool:
        """Set up session for ORBCOMM API"""
        try:
            # ORBCOMM uses query parameters for authentication
            self.session.headers.update({
                'User-Agent': 'ORBCOMM-Test-Client/1.0',
                'Accept': 'application/xml, text/xml, */*',
            })
            
            logger.info("Session configured for ORBCOMM API")
            return True
            
        except Exception as e:
            logger.error(f"Session setup failed: {e}")
            return False
    
    def get_messages(self, mobile_id: str = None, hours_back: int = 24) -> Optional[List[Dict]]:
        """
        Retrieve messages from ORBCOMM for a specific device
        
        Args:
            mobile_id: Device ID (SAT ID) - optional, can filter by from_id
            hours_back: How many hours back to retrieve messages
            
        Returns:
            List of parsed message dictionaries or None if failed
        """
        try:
            # Calculate time range (though the XML endpoint might not use time filters the same way)
            end_time = datetime.utcnow()
            start_time = end_time - timedelta(hours=hours_back)
            
            # ORBCOMM XML API parameters
            params = {
                'access_id': self.access_id,
                'password': self.password,
                'from_id': self.from_id,
            }
            
            # Add mobile_id if provided
            if mobile_id:
                params['mobile_id'] = mobile_id
            
            # Use the XML endpoint
            url = f"{self.base_url}/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml"
            
            logger.info(f"Retrieving messages from ORBCOMM XML API")
            logger.info(f"URL: {url}")
            logger.info(f"From ID: {self.from_id}")
            if mobile_id:
                logger.info(f"Mobile ID: {mobile_id}")
            
            response = self.session.get(url, params=params, timeout=30)
            logger.info(f"Response status: {response.status_code}")
            
            if response.status_code == 200:
                content = response.text
                logger.info(f"Received {len(content)} characters of XML data")
                
                # Parse XML
                try:
                    root = ET.fromstring(content)
                    
                    # Check for error
                    error_id_elem = root.find('.//ErrorID')
                    if error_id_elem is not None and error_id_elem.text != '0':
                        error_id = int(error_id_elem.text)
                        logger.warning(f"ORBCOMM API returned error ID: {error_id}")
                        
                        # Common ORBCOMM error codes
                        error_messages = {
                            21786: "No messages found for the specified criteria",
                            21787: "Invalid mobile ID",
                            21788: "Invalid time range",
                            21789: "Authentication failed",
                        }
                        
                        error_msg = error_messages.get(error_id, f"Unknown error code: {error_id}")
                        logger.warning(f"Error details: {error_msg}")
                        return []
                    
                    # Parse messages
                    messages = []
                    for message_elem in root.findall('.//ReturnMessage'):
                        message = self._parse_message(message_elem)
                        if message:
                            messages.append(message)
                    
                    logger.info(f"Parsed {len(messages)} messages")
                    return messages
                    
                except ET.ParseError as e:
                    logger.error(f"Failed to parse XML response: {e}")
                    return None
                
            else:
                logger.error(f"HTTP {response.status_code}: {response.text}")
                return None
                
        except Exception as e:
            logger.error(f"Error retrieving messages: {e}")
            return None
    
    def _parse_message(self, message_elem) -> Optional[Dict]:
        """Parse a single ReturnMessage element into a dictionary"""
        try:
            message = {}
            
            # Basic message info
            message['ID'] = message_elem.find('ID').text if message_elem.find('ID') is not None else None
            message['MessageUTC'] = message_elem.find('MessageUTC').text if message_elem.find('MessageUTC') is not None else None
            message['ReceiveUTC'] = message_elem.find('ReceiveUTC').text if message_elem.find('ReceiveUTC') is not None else None
            message['SIN'] = int(message_elem.find('SIN').text) if message_elem.find('SIN') is not None else None
            message['MobileID'] = message_elem.find('MobileID').text if message_elem.find('MobileID') is not None else None
            message['RegionName'] = message_elem.find('RegionName').text if message_elem.find('RegionName') is not None else None
            message['OTAMessageSize'] = int(message_elem.find('OTAMessageSize').text) if message_elem.find('OTAMessageSize') is not None else None
            
            # Parse payload
            payload_elem = message_elem.find('Payload')
            if payload_elem is not None:
                payload = {}
                payload['Name'] = payload_elem.get('Name')
                payload['SIN'] = int(payload_elem.get('SIN')) if payload_elem.get('SIN') else None
                payload['MIN'] = int(payload_elem.get('MIN')) if payload_elem.get('MIN') else None
                
                # Parse fields
                fields = {}
                for field_elem in payload_elem.findall('.//Field'):
                    field_name = field_elem.get('Name')
                    field_value = field_elem.get('Value')
                    
                    # Decode field values based on type
                    if field_name in ['Latitude', 'Longitude', 'PrevLatitude', 'PrevLongitude']:
                        # Convert from scaled integer to decimal degrees
                        if field_value:
                            fields[field_name] = float(field_value) / 10000.0
                        else:
                            fields[field_name] = None
                    elif field_name in ['Speed', 'Heading', 'InputVoltage', 'EventTime', 'PrevGpsTime', 'Version', 'GpsFixAge']:
                        # Convert to appropriate numeric type
                        if field_value:
                            try:
                                # Try integer first
                                fields[field_name] = int(field_value)
                            except ValueError:
                                # Try float
                                fields[field_name] = float(field_value)
                        else:
                            fields[field_name] = None
                    elif field_name in ['IsPowerSourceExternal']:
                        # Convert to boolean
                        fields[field_name] = field_value.lower() == 'true' if field_value else None
                    else:
                        # Keep as string
                        fields[field_name] = field_value
                
                payload['Fields'] = fields
                message['Payload'] = payload
            
            return message
            
        except Exception as e:
            logger.error(f"Error parsing message: {e}")
            return None
    


class MQTTPublisher:
    """MQTT client for publishing messages"""
    
    def __init__(self, broker_host: str, broker_port: int = 1883, client_id: str = "orbcomm_test"):
        self.broker_host = broker_host
        self.broker_port = broker_port
        self.client_id = client_id
        self.client = mqtt.Client(
            client_id=client_id,
            callback_api_version=mqtt.CallbackAPIVersion.VERSION2
        )
        self.connected = False
        
        # Set up callbacks
        self.client.on_connect = self._on_connect
        self.client.on_disconnect = self._on_disconnect
        self.client.on_publish = self._on_publish
        
    def _on_connect(self, client, userdata, flags, reason_code, properties):
        if reason_code == 0:
            self.connected = True
            logger.info("Connected to MQTT broker")
        else:
            logger.error(f"Failed to connect to MQTT broker: {reason_code}")
            
    def _on_disconnect(self, client, userdata, disconnect_flags, reason_code, properties):
        self.connected = False
        logger.info("Disconnected from MQTT broker")
        
    def _on_publish(self, client, userdata, mid, reason_code, properties):
        logger.debug(f"Message published: {mid}")
        
    def connect(self) -> bool:
        """Connect to MQTT broker"""
        try:
            logger.info(f"Connecting to MQTT broker at {self.broker_host}:{self.broker_port}")
            self.client.connect(self.broker_host, self.broker_port, 60)
            self.client.loop_start()
            
            # Wait for connection
            timeout = 10
            while not self.connected and timeout > 0:
                time.sleep(0.5)
                timeout -= 0.5
                
            return self.connected
            
        except Exception as e:
            logger.error(f"MQTT connection failed: {e}")
            return False
            
    def publish(self, topic: str, payload: str) -> bool:
        """Publish message to MQTT topic"""
        try:
            if not self.connected:
                logger.error("Not connected to MQTT broker")
                return False
                
            result = self.client.publish(topic, payload)
            
            if result.rc == mqtt.MQTT_ERR_SUCCESS:
                logger.info(f"Published to {topic}: {len(payload)} bytes")
                return True
            else:
                logger.error(f"Failed to publish to {topic}: {result.rc}")
                return False
                
        except Exception as e:
            logger.error(f"Publish error: {e}")
            return False
            
    def disconnect(self):
        """Disconnect from MQTT broker"""
        self.client.loop_stop()
        self.client.disconnect()

def format_for_skywave_protocol(messages: List[Dict], remote_addr: str = "orbcomm.test") -> str:
    """
    Format the parsed messages for the skywave protocol test
    Encodes JSON as hex string for safe MQTT transport
    """
    try:
        # Create the tracker data format expected by main.go
        tracker_data = {
            "payload": json.dumps(messages),  # Convert list of dicts to JSON string
            "remoteaddr": remote_addr
        }
        
        return json.dumps(tracker_data, indent=2)
        
    except Exception as e:
        logger.error(f"Error formatting data: {e}")
        return json.dumps({"error": str(e)})

def main():
    parser = argparse.ArgumentParser(description='ORBCOMM to MQTT Publisher')
    parser.add_argument('--mqtt-host', default='localhost', help='MQTT broker host')
    parser.add_argument('--mqtt-port', type=int, default=1883, help='MQTT broker port')
    parser.add_argument('--topic', default='tracker/from-tcp', help='MQTT topic to publish to')
    parser.add_argument('--interval', type=int, default=60, help='Polling interval in seconds')
    parser.add_argument('--hours-back', type=int, default=24, help='Hours of historical data to retrieve')
    parser.add_argument('--once', action='store_true', help='Run once and exit')
    parser.add_argument('--verbose', '-v', action='store_true', help='Enable verbose logging')
    
    args = parser.parse_args()
    
    if args.verbose:
        logging.getLogger().setLevel(logging.DEBUG)
    
    # ORBCOMM credentials from the working curl command
    access_id = "70001184"
    password = "JEUTPKKH"
    from_id = "13969586728"
    from_id = "3898451771"
    device_id = "02092247SKY6A70"  # SAT ID from datos.md
    
    logger.info("Starting ORBCOMM to MQTT publisher")
    logger.info(f"Access ID: {access_id}")
    logger.info(f"From ID: {from_id}")
    logger.info(f"Device ID: {device_id}")
    logger.info(f"MQTT Broker: {args.mqtt_host}:{args.mqtt_port}")
    logger.info(f"MQTT Topic: {args.topic}")
    
    # Initialize clients
    orbcomm = ORBCOMMClient(access_id, password, from_id)
    mqtt_pub = MQTTPublisher(args.mqtt_host, args.mqtt_port)
    
    # Set up session
    if not orbcomm.authenticate():
        logger.error("Failed to set up ORBCOMM session")
        return 1
    
    # Connect to MQTT
    if not mqtt_pub.connect():
        logger.error("Failed to connect to MQTT broker")
        return 1
    
    try:
        while True:
            logger.info("Retrieving messages from ORBCOMM...")
            
            # Get messages from ORBCOMM
            messages = orbcomm.get_messages(device_id, args.hours_back)
            
            if messages:
                logger.info(f"Retrieved {len(messages)} messages")
                
                # Format for skywave protocol
                formatted_data = format_for_skywave_protocol(messages)
                
                # Publish to MQTT
                if mqtt_pub.publish(args.topic, formatted_data):
                    logger.info("Successfully published to MQTT")
                    
                    # Also publish raw JSON for debugging
                    debug_topic = f"{args.topic}/debug"
                    mqtt_pub.publish(debug_topic, json.dumps(messages, indent=2))
                    logger.info(f"Published raw JSON to {debug_topic}")
                else:
                    logger.error("Failed to publish to MQTT")
            else:
                logger.warning("No messages retrieved from ORBCOMM")
            
            if args.once:
                break
                
            logger.info(f"Waiting {args.interval} seconds before next poll...")
            time.sleep(args.interval)
            
    except KeyboardInterrupt:
        logger.info("Stopping due to user interrupt")
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        return 1
    finally:
        mqtt_pub.disconnect()
        logger.info("Cleanup completed")
    
    return 0

if __name__ == "__main__":
    exit(main())
