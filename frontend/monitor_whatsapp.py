import time
import mysql.connector
import requests
import subprocess
import json
import os
import argparse
from datetime import datetime, timedelta
import logging
import urllib3
from typing import Dict, List, Tuple, Optional

# MessageBird imports
import messagebird
from messagebird.conversation_message import MESSAGE_TYPE_HSM

# Disable SSL warnings for self-signed certificates
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

# Configuration constants
DEFAULT_REQUEST_TIMEOUT = 30  # Increased default timeout
DEFAULT_CHECK_INTERVAL = 5
DEFAULT_DATA_THRESHOLD = 5
DEFAULT_RETRY_ATTEMPTS = 3
DEFAULT_RETRY_DELAY = 5

# Create logs directory if it doesn't exist
log_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'logs')
os.makedirs(log_dir, exist_ok=True)
log_file = os.path.join(log_dir, 'elasticsearch_monitor.log')

def setup_logging(debug_mode=False, debug_level=1):
    """Setup logging configuration based on debug mode."""
    if debug_mode and debug_level == 3:
        # For debug level 3, only show minimal info output
        log_level = logging.INFO
        log_format = '%(message)s'
    else:
        log_level = logging.DEBUG if debug_mode else logging.INFO
        log_format = '%(asctime)s - %(levelname)s - %(message)s'
    
    # Clear existing handlers
    for handler in logging.root.handlers[:]:
        logging.root.removeHandler(handler)
    
    # Configure logging
    logging.basicConfig(
        level=log_level,
        format=log_format,
        handlers=[
            logging.FileHandler(log_file),
            logging.StreamHandler()
        ]
    )
    
    # For debug level 3, suppress mysql and urllib3 debug logs
    if debug_mode and debug_level == 3:
        logging.getLogger('mysql.connector').setLevel(logging.WARNING)
        logging.getLogger('urllib3').setLevel(logging.WARNING)

class ElasticsearchMonitor:
    def __init__(self, debug_mode=False, debug_level=1, force_alerts=False):
        self.debug_mode = debug_mode
        self.debug_level = debug_level if debug_mode else 0
        self.force_alerts = force_alerts  # Add force flag
        self.db_config = {
            'host': 'localhost',
            'user': 'gpscontrol',
            'password': 'qazwsxedc',
            'database': 'automation_app'
        }
        
        # Track problems for WhatsApp alerts
        self.whatsapp_problems = {}
        
        # Create session with default settings
        self.session = requests.Session()
        self.session.verify = False  # For self-signed certificates
        
        # Setup logging with debug mode
        setup_logging(debug_mode, debug_level)
        
        if debug_mode:
            if debug_level == 3:
                pass  # No verbose logging messages for minimal output
            else:
                logging.info("Debug mode enabled - detailed logging active")
        
        self.ensure_monitor_config_table()
        self.ensure_whatsapp_config_table()
        self.ensure_whatsapp_phones_table()
    
    def __del__(self):
        """Cleanup when object is destroyed"""
        if hasattr(self, 'session'):
            try:
                self.session.close()
            except:
                pass
    
    def get_db_connection(self):
        """Create database connection."""
        return mysql.connector.connect(**self.db_config)
    
    def ensure_monitor_config_table(self):
        """Create monitor configuration table if it doesn't exist."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor()
            
            create_table_sql = """
            CREATE TABLE IF NOT EXISTS monitor_config (
                id INT AUTO_INCREMENT PRIMARY KEY,
                service_type VARCHAR(50) NOT NULL UNIQUE,
                check_interval_minutes INT DEFAULT 5,
                data_threshold_minutes INT DEFAULT 5,
                restart_enabled BOOLEAN DEFAULT TRUE,
                alert_enabled BOOLEAN DEFAULT FALSE,
                alert_webhook_url VARCHAR(500) DEFAULT NULL,
                request_timeout_seconds INT DEFAULT 30,
                retry_attempts INT DEFAULT 3,
                retry_delay_seconds INT DEFAULT 5,
                last_check TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                last_success TIMESTAMP NULL,
                last_failure TIMESTAMP NULL,
                failure_count INT DEFAULT 0,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
            )
            """
            cursor.execute(create_table_sql)
            
            # Insert default configurations for each service type
            services_with_elastic = self.get_services_with_elastic()
            for service_type in services_with_elastic:
                cursor.execute("""
                    INSERT IGNORE INTO monitor_config (service_type) 
                    VALUES (%s)
                """, (service_type,))
            
            db.commit()
            logging.info("Monitor configuration table ready")
            
        except Exception as e:
            logging.error(f"Error creating monitor config table: {e}")
        finally:
            cursor.close()
            db.close()
    
    def ensure_whatsapp_config_table(self):
        """Create WhatsApp configuration table if it doesn't exist."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor()
            
            create_table_sql = """
            CREATE TABLE IF NOT EXISTS whatsapp_config (
                id INT AUTO_INCREMENT PRIMARY KEY,
                enabled BOOLEAN DEFAULT FALSE,
                api_key VARCHAR(255) NOT NULL,
                channel_id VARCHAR(255) NOT NULL,
                namespace VARCHAR(255) NOT NULL,
                template_name VARCHAR(100) DEFAULT 'server_alerts',
                language_code VARCHAR(10) DEFAULT 'es_MX',
                alert_interval_minutes INT DEFAULT 30,
                min_failure_count INT DEFAULT 1,
                last_alert_sent TIMESTAMP NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
            )
            """
            cursor.execute(create_table_sql)
            db.commit()
            logging.debug("WhatsApp configuration table ready")
            
        except Exception as e:
            logging.error(f"Error creating WhatsApp config table: {e}")
        finally:
            cursor.close()
            db.close()
    
    def ensure_whatsapp_phones_table(self):
        """Create WhatsApp phones table if it doesn't exist."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor()
            
            create_table_sql = """
            CREATE TABLE IF NOT EXISTS whatsapp_phones (
                id INT AUTO_INCREMENT PRIMARY KEY,
                phone_number VARCHAR(20) NOT NULL,
                contact_name VARCHAR(100) DEFAULT NULL,
                enabled BOOLEAN DEFAULT TRUE,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                UNIQUE KEY unique_phone (phone_number)
            )
            """
            cursor.execute(create_table_sql)
            db.commit()
            logging.debug("WhatsApp phones table ready")
            
        except Exception as e:
            logging.error(f"Error creating WhatsApp phones table: {e}")
        finally:
            cursor.close()
            db.close()
    
    def get_whatsapp_config(self) -> Dict:
        """Get WhatsApp configuration."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor(dictionary=True)
            
            cursor.execute("SELECT * FROM whatsapp_config WHERE enabled = TRUE LIMIT 1")
            config = cursor.fetchone()
            
            return config or {}
            
        except Exception as e:
            logging.error(f"Error getting WhatsApp config: {e}")
            return {}
        finally:
            cursor.close()
            db.close()
    
    def get_whatsapp_phones(self) -> List[str]:
        """Get enabled phone numbers for alerts."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor()
            
            cursor.execute("SELECT phone_number FROM whatsapp_phones WHERE enabled = TRUE")
            phones = [row[0] for row in cursor.fetchall()]
            
            return phones
            
        except Exception as e:
            logging.error(f"Error getting WhatsApp phones: {e}")
            return []
        finally:
            cursor.close()
            db.close()
    
    def should_send_whatsapp_alert(self, whatsapp_config: Dict) -> bool:
        """Check if it's time to send WhatsApp alert based on global configuration."""
        try:
            if not whatsapp_config.get('enabled', False):
                logging.debug("⛔ WhatsApp alerts disabled in config")
                return False
            
            # Check if we have any problems to report
            if not self.whatsapp_problems:
                logging.debug("⛔ No problems found to report")
                return False
            
            alert_interval = whatsapp_config.get('alert_interval_minutes', 30)
            min_failures = whatsapp_config.get('min_failure_count', 1)
            
            logging.info(f"🔍 Alert conditions check:")
            logging.info(f"   📊 Problems found: {len(self.whatsapp_problems)} services")
            logging.info(f"   🔢 Min failure count required: {min_failures}")
            logging.info(f"   ⏰ Alert interval: {alert_interval} minutes")
            logging.info(f"   🚀 Force mode: {self.force_alerts}")
            
            # Check global timing - has enough time passed since last alert?
            if not self.force_alerts:  # Skip timing check if forcing
                last_alert = whatsapp_config.get('last_alert_sent')
                if last_alert:
                    if isinstance(last_alert, str):
                        last_alert = datetime.fromisoformat(last_alert.replace('Z', '+00:00'))
                    
                    time_since_alert = datetime.now() - last_alert.replace(tzinfo=None)
                    logging.info(f"   ⏱️ Time since last alert: {time_since_alert}")
                    
                    if time_since_alert < timedelta(minutes=alert_interval):
                        logging.info(f"⛔ WhatsApp alert skipped - only {time_since_alert} passed, need {alert_interval} minutes")
                        logging.info(f"💡 TIP: Use --force to bypass timing check")
                        return False
                else:
                    logging.info(f"   ⏱️ No previous alert found - OK to send")
            else:
                logging.info(f"   🚀 FORCE MODE - bypassing timing check")
            
            # Check if any service has enough failures
            services_ready = 0
            for service_type, problems in self.whatsapp_problems.items():
                service_config = self.get_monitor_config(service_type)
                failure_count = service_config.get('failure_count', 0)
                logging.info(f"   🔸 {service_type}: {failure_count} failures (need {min_failures})")
                
                # Check if enough failures have occurred for this service
                if failure_count >= min_failures:
                    services_ready += 1
            
            logging.info(f"   ✅ Services ready to alert: {services_ready}/{len(self.whatsapp_problems)}")
            
            if services_ready > 0:
                logging.info("✅ WhatsApp alert conditions MET")
                return True
            else:
                logging.info("⛔ WhatsApp alert conditions NOT met - no services have enough failures")
                return False
            
        except Exception as e:
            logging.error(f"Error checking if should send WhatsApp alert: {e}")
            return False
    
    def format_whatsapp_messages(self) -> List[tuple]:
        """Format problems into individual WhatsApp messages per DOCUMENT for roadapp template."""
        try:
            messages = []
            
            for service_type, data in self.whatsapp_problems.items():
                
                for problem in data['problems']:
                    doc_name = problem['doc_name']
                    error_msg = problem['error']
                    count = problem.get('count', 0)
                    
                    # Parameter 1: Company-Service name
                    param1 = f"🚨 Jonobridge-{service_type.upper()}"
                    
                    # Parameter 2: Document name + simplified error
                    simplified_error = self.simplify_error_message(error_msg, count)
                    param2 = f"⚠️ {doc_name} ({simplified_error})"
                    
                    messages.append((param1, param2))
                    
                    # Log each message
                    logging.info(f"📝 WhatsApp message for {doc_name}:")
                    logging.info(f"   📋 Service: {param1}")
                    logging.info(f"   📄 Problem: {param2}")
                    logging.info(f"   📏 Length: {len(param2)} characters")
            
            logging.info(f"📊 Total individual messages to send: {len(messages)}")
            return messages
            
        except Exception as e:
            logging.error(f"Error formatting WhatsApp messages: {e}")
            return [("Jonobridge-ERROR", "Error al generar mensajes de monitoreo automático")]
    
    def simplify_error_message(self, error_msg: str, count: int = 0) -> str:
        """Extract just the essential error information from complex error messages."""
        try:
            if not error_msg or not error_msg.strip():
                return "No data found" if count == 0 else "Unknown error"
            
            error_msg = error_msg.strip()
            
            # Extract HTTP status codes (404, 500, 302, etc.)
            import re
            http_status = re.search(r'HTTP (\d{3})', error_msg)
            if http_status:
                return f"HTTP {http_status.group(1)}"
            
            # Look for just status codes in the message
            status_code = re.search(r'"status":(\d{3})', error_msg)
            if status_code:
                return f"HTTP {status_code.group(1)}"
            
            # Extract simple error types
            if "timeout" in error_msg.lower():
                return "Timeout"
            elif "connection" in error_msg.lower():
                return "Connection error"
            elif "no data" in error_msg.lower():
                return "No data"
            elif "index_not_found" in error_msg.lower():
                return "Index not found"
            elif "unauthorized" in error_msg.lower():
                return "Unauthorized"
            elif "forbidden" in error_msg.lower():
                return "Forbidden"
            
            # If error message is short (< 50 chars), use as-is
            if len(error_msg) <= 50:
                return error_msg
            
            # For long messages, try to extract the first meaningful part
            # Split by common separators and take the first part
            for separator in [': {', ' - ', ': "', ': (']:
                if separator in error_msg:
                    first_part = error_msg.split(separator)[0].strip()
                    if len(first_part) <= 50 and first_part:
                        return first_part
            
            # Last resort: truncate to first 50 characters
            return error_msg[:47] + "..." if len(error_msg) > 50 else error_msg
            
        except Exception as e:
            logging.debug(f"Error simplifying message: {e}")
            return "Error occurred"
    
    def messagebird_send_waba(self, phone, templateName, jparam, language_code, api_key, channel_id, namespace):
        """Simplified MessageBird function exactly like working snippet."""
        error_state = False
        try:
            # Import fresh each time like in the working snippet
            import messagebird
            from messagebird.conversation_message import MESSAGE_TYPE_HSM
            
            # Create client exactly like working snippet
            client = messagebird.Client(api_key)
            
            # Use exact payload structure from working snippet
            msg = client.conversation_start({
                'channelId': channel_id,
                'to': phone,
                'type': MESSAGE_TYPE_HSM,
                "content": {
                    "hsm": {
                        "namespace": namespace,
                        "templateName": templateName,
                        "language": {
                            "policy": "deterministic",
                            "code": language_code
                        },
                        "params": jparam
                    }
                }
            })
            message = msg.status
            
            # Log response details
            response_id = getattr(msg, 'id', 'N/A')
            created_time = getattr(msg, 'createdDatetime', 'N/A')
            
            logging.info(f"✅ MessageBird Response for {phone} - ID: {response_id}, Status: {message}")
            logging.info(f"   📱 Phone: {phone}")
            logging.info(f"   📅 Created: {created_time}")
            
        except Exception as e:
            error_state = True
            message = "Error:" + str(e)
            logging.error(f"❌ Error sending to {phone}: {message}")

        return {'error': error_state, 'message': message}
    
    def send_whatsapp_alert(self, whatsapp_config: Dict, message_list: List[tuple]) -> bool:
        """Fixed WhatsApp alert sending - exactly like working snippet pattern."""
        try:
            api_key = whatsapp_config['api_key']
            channel_id = whatsapp_config['channel_id']
            namespace = whatsapp_config['namespace']
            template_name = whatsapp_config.get('template_name', 'roadapp')
            language_code = whatsapp_config.get('language_code', 'es_MX')
            
            logging.info(f"📱 Preparing WhatsApp alerts:")
            logging.info(f"   📋 Template: {template_name}")
            logging.info(f"   🌍 Language: {language_code}")
            logging.info(f"   📨 Messages to send: {len(message_list)}")
            
            phones = self.get_whatsapp_phones()
            if not phones:
                logging.warning("No phone numbers configured for WhatsApp alerts")
                return False
            
            logging.info(f"📞 Found {len(phones)} phone numbers:")
            for phone in phones:
                logging.info(f"   📱 {phone}")
            
            total_success = 0
            total_attempts = 0
            
            # Send each message to all phones - exactly like working snippet
            for msg_idx, (param1, param2) in enumerate(message_list, 1):
                logging.info(f"\n📨 MESSAGE {msg_idx}/{len(message_list)}: {param1}")
                logging.info(f"   📄 Content: {param2}")
                
                # Prepare parameters exactly like working snippet
                jparam = [
                    {'default': param1},  # {{1}} - Company-Service name
                    {'default': param2}   # {{2}} - Problem details
                ]
                
                message_success = 0
                
                # Send this ONE message to ALL phones (like working snippet)
                for phone in phones:
                    try:
                        total_attempts += 1
                        logging.info(f"   📞 Sending to {phone}...")
                        
                        # Use simple MessageBird function exactly like working snippet
                        resultado = self.messagebird_send_waba(
                            phone=phone,
                            templateName=template_name,
                            jparam=jparam,
                            language_code=language_code,
                            api_key=api_key,
                            channel_id=channel_id,
                            namespace=namespace
                        )
                        
                        if not resultado['error']:
                            total_success += 1
                            message_success += 1
                            logging.info(f"      ✅ SUCCESS - {phone}: {resultado['message']}")
                        else:
                            logging.error(f"      ❌ FAILED - {phone}: {resultado['message']}")
                    
                    except Exception as e:
                        logging.error(f"      ❌ FAILED - Exception sending to {phone}: {e}")
                
                logging.info(f"   📊 Message {msg_idx} summary: {message_success}/{len(phones)} successful")
                
                # Small delay between messages only
                if msg_idx < len(message_list):
                    logging.info(f"   ⏱️ Waiting 1 second before next message...")
                    time.sleep(1)
            
            # Overall Summary
            logging.info(f"\n📊 COMPLETE WhatsApp sending summary:")
            logging.info(f"   📨 Total messages: {len(message_list)}")
            logging.info(f"   📱 Total phones: {len(phones)}")
            logging.info(f"   🔢 Total attempts: {total_attempts}")
            logging.info(f"   ✅ Successful sends: {total_success}")
            logging.info(f"   ❌ Failed sends: {total_attempts - total_success}")
            logging.info(f"   📈 Success rate: {(total_success/total_attempts*100):.1f}%")
            
            return total_success > 0
            
        except Exception as e:
            logging.error(f"Error in send_whatsapp_alert: {e}")
            import traceback
            logging.error(f"Traceback: {traceback.format_exc()}")
            return False
    
    def update_whatsapp_alert_timestamps(self):
        """Update last WhatsApp alert timestamp globally."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor()
            
            # Update the global WhatsApp config with last alert time
            cursor.execute("""
                UPDATE whatsapp_config 
                SET last_alert_sent = CURRENT_TIMESTAMP,
                    updated_at = CURRENT_TIMESTAMP
                WHERE enabled = TRUE
            """)
            
            db.commit()
            
        except Exception as e:
            logging.error(f"Error updating WhatsApp alert timestamps: {e}")
        finally:
            cursor.close()
            db.close()
    
    def process_whatsapp_alerts(self):
        """Process and send WhatsApp alerts if conditions are met."""
        try:
            logging.info("🔍 Processing WhatsApp alerts...")
            
            whatsapp_config = self.get_whatsapp_config()
            
            if not whatsapp_config:
                logging.warning("⛔ No WhatsApp configuration found or WhatsApp alerts disabled")
                logging.info("💡 TIP: Check whatsapp_config table - enabled should be TRUE")
                return
            
            logging.info(f"✅ WhatsApp configuration loaded:")
            logging.info(f"   📋 Template: {whatsapp_config.get('template_name', 'N/A')}")
            logging.info(f"   🌍 Language: {whatsapp_config.get('language_code', 'N/A')}")
            logging.info(f"   ⏰ Alert interval: {whatsapp_config.get('alert_interval_minutes', 'N/A')} minutes")
            logging.info(f"   🔢 Min failure count: {whatsapp_config.get('min_failure_count', 'N/A')}")
            logging.info(f"   📛 Enabled: {whatsapp_config.get('enabled', False)}")
            
            # If no problems found during monitoring cycle, check for existing failures in database
            if not self.whatsapp_problems:
                logging.info("🔍 No problems found in current cycle, checking database for existing failures...")
                self.check_existing_failures()
            
            # Check if we have problems to report
            if not self.whatsapp_problems:
                logging.warning("⚠️ No problems found to report via WhatsApp")
                return
            
            if not self.should_send_whatsapp_alert(whatsapp_config):
                logging.info("⭕ WhatsApp alert conditions not met - skipping")
                return
            
            logging.info("✅ WhatsApp alert conditions met - proceeding with alert")
                
            logging.info(f"📊 Problems summary:")
            for service_type, data in self.whatsapp_problems.items():
                logging.info(f"   🔸 {service_type}: {len(data['problems'])} problems")
            
            message_list = self.format_whatsapp_messages()
            
            if not message_list:
                logging.warning("⛔ No messages generated for WhatsApp alert")
                return
            
            logging.info("📱 Sending individual WhatsApp alerts using roadapp template...")
            logging.info(f"📨 Total messages to send: {len(message_list)}")
            
            for i, (param1, param2) in enumerate(message_list, 1):
                logging.info(f"📝 Message {i}: {param1}")
                logging.info(f"📄 Preview: {param2[:50]}...")
            
            success = self.send_whatsapp_alert(whatsapp_config, message_list)
            
            if success:
                logging.info("🎉 WhatsApp alerts sent successfully!")
                self.update_whatsapp_alert_timestamps()
                logging.info("📅 WhatsApp alert timestamp updated")
            else:
                logging.error("💥 Failed to send WhatsApp alerts")
                
        except Exception as e:
            logging.error(f"💥 Error processing WhatsApp alerts: {e}")
            logging.error(f"Exception type: {type(e).__name__}")
            import traceback
            logging.error(f"Traceback: {traceback.format_exc()}")
    
    def check_existing_failures(self):
        """Check database for services with existing failures and create problems list."""
        try:
            logging.info("🔍 Checking database for services with existing failures...")
            
            db = self.get_db_connection()
            cursor = db.cursor(dictionary=True)
            
            # Get services with failures
            cursor.execute("""
                SELECT service_type, failure_count, last_failure 
                FROM monitor_config 
                WHERE failure_count > 0
            """)
            
            failed_services = cursor.fetchall()
            logging.info(f"📊 Found {len(failed_services)} services with failures in database")
            
            for service in failed_services:
                service_type = service['service_type']
                failure_count = service['failure_count']
                last_failure = service['last_failure']
                
                logging.info(f"   🔸 {service_type}: {failure_count} failures, last: {last_failure}")
                
                # Get documents for this service to create problem entries
                documents = self.get_elastic_documents(service_type)
                
                if documents:
                    problems = []
                    for doc in documents:
                        problems.append({
                            'doc_name': doc['elastic_doc_name'],
                            'error': f"Service has {failure_count} accumulated failures",
                            'count': 0,
                            'client_id': doc['client_id']
                        })
                    
                    self.whatsapp_problems[service_type] = {
                        'problems': problems,
                        'count': len(problems),
                        'timestamp': datetime.now()
                    }
                    
                    logging.info(f"   ✅ Added {len(problems)} problems for {service_type}")
            
            logging.info(f"📊 Total services with problems: {len(self.whatsapp_problems)}")
            
        except Exception as e:
            logging.error(f"Error checking existing failures: {e}")
        finally:
            cursor.close()
            db.close()
    
    def get_services_with_elastic(self) -> List[str]:
        """Get list of service types that have elastic_doc_name column."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor()
            
            cursor.execute("""
                SELECT DISTINCT table_name 
                FROM information_schema.columns 
                WHERE column_name = 'elastic_doc_name' 
                AND table_schema = 'automation_app'
            """)
            
            services = [row[0] for row in cursor.fetchall()]
            logging.debug(f"Services with elastic_doc_name column: {services}")
            return services
            
        except Exception as e:
            logging.error(f"Error getting services with elastic: {e}")
            return []
        finally:
            cursor.close()
            db.close()
    
    def get_monitor_config(self, service_type: str) -> Dict:
        """Get monitor configuration for a service type."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor(dictionary=True)
            
            cursor.execute("""
                SELECT * FROM monitor_config WHERE service_type = %s
            """, (service_type,))
            
            config = cursor.fetchone()
            if not config:
                # Return default config
                config = {
                    'check_interval_minutes': DEFAULT_CHECK_INTERVAL,
                    'data_threshold_minutes': DEFAULT_DATA_THRESHOLD,
                    'restart_enabled': True,
                    'alert_enabled': False,
                    'alert_webhook_url': None,
                    'request_timeout_seconds': DEFAULT_REQUEST_TIMEOUT,
                    'retry_attempts': DEFAULT_RETRY_ATTEMPTS,
                    'retry_delay_seconds': DEFAULT_RETRY_DELAY,
                    'failure_count': 0
                }
            
            return config
            
        except Exception as e:
            logging.error(f"Error getting monitor config for {service_type}: {e}")
            return {}
        finally:
            cursor.close()
            db.close()
    
    def get_elastic_documents(self, service_type: str) -> List[Dict]:
        """Get all elastic documents for a service type."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor()
            
            # Get required columns for the service
            required_columns = ['client_id', 'elastic_doc_name', 'elastic_url', 'elastic_user', 'elastic_password']
            existing_columns = []
            
            cursor.execute("""
                SELECT column_name 
                FROM information_schema.columns 
                WHERE table_name = %s 
                AND table_schema = 'automation_app'
                AND column_name IN ('client_id', 'elastic_doc_name', 'elastic_url', 'elastic_user', 'elastic_password')
            """, (service_type,))
            
            # Fix: Access the column_name correctly (tuple access)
            results = cursor.fetchall()
            existing_columns = [row[0] for row in results]  # row is a tuple: (column_name,)
            
            logging.debug(f"Existing columns for {service_type}: {existing_columns}")
            
            if 'elastic_doc_name' not in existing_columns:
                logging.debug(f"No elastic_doc_name column found in {service_type}")
                return []
            
            # Build query with only existing columns
            select_columns = ', '.join([col for col in required_columns if col in existing_columns])
            
            if self.debug_mode:
                logging.debug(f"Selecting columns: {select_columns} from {service_type}")
            
            # Create new cursor with dictionary mode for data query
            cursor.close()
            cursor = db.cursor(dictionary=True)
            
            query = f"""
                SELECT {select_columns}
                FROM {service_type} 
                WHERE elastic_doc_name IS NOT NULL 
                AND elastic_doc_name != ''
            """
            
            if self.debug_mode:
                logging.debug(f"Executing query: {query}")
            
            cursor.execute(query)
            documents = cursor.fetchall()
            
            if self.debug_mode:
                logging.debug(f"Found {len(documents)} documents in {service_type}")
                for i, doc in enumerate(documents):
                    logging.debug(f"  Document {i+1}: {doc.get('elastic_doc_name', 'Unknown')}")
            
            return documents
            
        except Exception as e:
            logging.error(f"Error getting elastic documents for {service_type}: {e}")
            import traceback
            logging.error(f"Traceback: {traceback.format_exc()}")
            return []
        finally:
            cursor.close()
            db.close()
    
    def check_elasticsearch_data(self, elastic_url: str, elastic_user: str, elastic_password: str, 
                                doc_name: str, threshold_minutes: int = 5, 
                                timeout_seconds: int = 30, retry_attempts: int = 3, 
                                retry_delay: int = 5) -> Tuple[bool, int, str]:
        """Check if data exists in Elasticsearch within the threshold time with retry logic."""
        
        # Normalize doc_name to lowercase and strip whitespace for Elasticsearch compatibility
        normalized_doc_name = doc_name.strip().lower()
        
        # Prepare the query
        query = {
            "query": {
                "range": {
                    "time": {
                        "gte": f"now-{threshold_minutes}m",
                        "lte": "now"
                    }
                }
            }
        }
        
        url = f"{elastic_url}/{normalized_doc_name}/_count"
        last_error = ""
        
        for attempt in range(retry_attempts):
            try:
                if self.debug_level == 2:
                    # Debug level 2: Show only Elasticsearch queries
                    logging.info(f"🔍 [{attempt + 1}/{retry_attempts}] Querying: {normalized_doc_name}")
                    if doc_name != normalized_doc_name:
                        logging.info(f"    🔧 Original: {repr(doc_name)} → Normalized: {repr(normalized_doc_name)}")
                    logging.info(f"    🌐 URL: {url}")
                    logging.info(f"    🕐 Query: last {threshold_minutes} minutes")
                elif self.debug_mode:
                    logging.debug(f"Attempt {attempt + 1}/{retry_attempts} for {normalized_doc_name} (original: {doc_name})")
                    logging.debug(f"Query: {json.dumps(query, indent=2)}")
                    logging.debug(f"URL: {url}")
                
                response = self.session.get(
                    url,
                    auth=(elastic_user, elastic_password),
                    headers={'Content-Type': 'application/json'},
                    json=query,
                    timeout=timeout_seconds
                )
                
                if self.debug_level == 2:
                    logging.info(f"    📊 Response: HTTP {response.status_code}")
                elif self.debug_mode:
                    logging.debug(f"Response status: {response.status_code}")
                    logging.debug(f"Response content: {response.text[:500]}...")
                
                if response.status_code == 200:
                    result = response.json()
                    count = result.get('count', 0)
                    if self.debug_level == 2:
                        logging.info(f"    ✅ Count: {count} records")
                    elif self.debug_level == 3 and count == 0:
                        # Debug level 3: Show only problems - minimal format
                        pass  # Will be handled in monitor_service_type
                    elif self.debug_mode:
                        logging.debug(f"Elasticsearch returned count: {count} for {normalized_doc_name} (original: {doc_name})")
                    return count > 0, count, ""
                else:
                    last_error = f"HTTP {response.status_code}: {response.text}"
                    if self.debug_level == 2:
                        logging.info(f"    ⛔ Error: HTTP {response.status_code}")
                    elif self.debug_level == 3:
                        # Debug level 3: Show only problems - minimal format
                        pass  # Will be handled in monitor_service_type
                    else:
                        logging.warning(f"Elasticsearch query failed for {normalized_doc_name} (original: {doc_name}) (attempt {attempt + 1}): {last_error}")
                    
                    if response.status_code in [404, 401, 403]:
                        # Don't retry for these errors
                        break
                        
            except requests.exceptions.Timeout:
                last_error = f"Request timeout after {timeout_seconds} seconds"
                if self.debug_level == 2:
                    logging.info(f"    ⏰ Timeout: {timeout_seconds}s")
                elif self.debug_level == 3:
                    # Debug level 3: Show only problems - minimal format
                    pass  # Will be handled in monitor_service_type
                else:
                    logging.warning(f"Timeout checking {normalized_doc_name} (original: {doc_name}) (attempt {attempt + 1}): {last_error}")
            except requests.exceptions.ConnectionError as e:
                last_error = f"Connection error: {str(e)}"
                if self.debug_level == 2:
                    logging.info(f"    🔌 Connection error")
                elif self.debug_level == 3:
                    # Debug level 3: Show only problems - minimal format
                    pass  # Will be handled in monitor_service_type
                else:
                    logging.warning(f"Connection error for {normalized_doc_name} (original: {doc_name}) (attempt {attempt + 1}): {last_error}")
            except Exception as e:
                last_error = f"Unexpected error: {str(e)}"
                if self.debug_level == 2:
                    logging.info(f"    💥 Error: {str(e)}")
                elif self.debug_level == 3:
                    # Debug level 3: Show only problems - minimal format
                    pass  # Will be handled in monitor_service_type
                else:
                    logging.error(f"Error checking Elasticsearch for {normalized_doc_name} (original: {doc_name}) (attempt {attempt + 1}): {last_error}")
            
            # Wait before retry (except on last attempt)
            if attempt < retry_attempts - 1:
                logging.debug(f"Waiting {retry_delay} seconds before retry...")
                time.sleep(retry_delay)
        
        return False, 0, last_error
    
    def get_client_namespace(self, client_id: int) -> Optional[str]:
        """Get client namespace by client_id."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor()
            
            cursor.execute("SELECT name FROM clients WHERE id = %s", (client_id,))
            result = cursor.fetchone()
            
            return result[0] if result else None
            
        except Exception as e:
            logging.error(f"Error getting client namespace for ID {client_id}: {e}")
            return None
        finally:
            cursor.close()
            db.close()
    
    def restart_service(self, namespace: str, service_type: str) -> bool:
        """Restart Kubernetes deployment for a service."""
        try:
            # First check if the deployment exists
            check_cmd = ['kubectl', 'get', 'deployment', service_type, '-n', namespace]
            check_result = subprocess.run(check_cmd, capture_output=True, text=True)
            
            if check_result.returncode != 0:
                logging.warning(f"Deployment {service_type} not found in namespace {namespace}")
                return False
            
            # Restart the deployment
            restart_cmd = ['kubectl', 'rollout', 'restart', f'deployment/{service_type}', '-n', namespace]
            result = subprocess.run(restart_cmd, capture_output=True, text=True)
            
            if result.returncode == 0:
                logging.info(f"Successfully restarted {service_type} in namespace {namespace}")
                return True
            else:
                logging.error(f"Failed to restart {service_type} in {namespace}: {result.stderr}")
                return False
                
        except Exception as e:
            logging.error(f"Error restarting service {service_type} in {namespace}: {e}")
            return False
    
    def send_alert(self, webhook_url: str, message: str) -> bool:
        """Send alert to webhook URL."""
        try:
            payload = {
                'text': message,
                'timestamp': datetime.now().isoformat()
            }
            
            response = self.session.post(
                webhook_url,
                json=payload,
                timeout=10
            )
            
            return response.status_code == 200
            
        except Exception as e:
            logging.error(f"Error sending alert: {e}")
            return False
    
    def update_monitor_status(self, service_type: str, success: bool, error_message: str = ""):
        """Update monitor status in database."""
        try:
            db = self.get_db_connection()
            cursor = db.cursor()
            
            if success:
                cursor.execute("""
                    UPDATE monitor_config 
                    SET last_check = CURRENT_TIMESTAMP,
                        last_success = CURRENT_TIMESTAMP,
                        failure_count = 0
                    WHERE service_type = %s
                """, (service_type,))
            else:
                cursor.execute("""
                    UPDATE monitor_config 
                    SET last_check = CURRENT_TIMESTAMP,
                        last_failure = CURRENT_TIMESTAMP,
                        failure_count = failure_count + 1
                    WHERE service_type = %s
                """, (service_type,))
            
            db.commit()
            
        except Exception as e:
            logging.error(f"Error updating monitor status for {service_type}: {e}")
        finally:
            cursor.close()
            db.close()
    
    def should_check_service(self, service_type: str, config: Dict) -> bool:
        """Check if it's time to monitor this service based on last check time."""
        try:
            # Force mode bypasses all timing checks
            if self.force_alerts:
                logging.debug(f"🚀 FORCE MODE - checking {service_type} regardless of timing")
                return True
                
            last_check = config.get('last_check')
            check_interval = config.get('check_interval_minutes', DEFAULT_CHECK_INTERVAL)
            
            if not last_check:
                logging.debug(f"No previous check found for {service_type}, will check now")
                return True
            
            # Convert to datetime if it's a string
            if isinstance(last_check, str):
                last_check = datetime.fromisoformat(last_check.replace('Z', '+00:00'))
            
            time_since_check = datetime.now() - last_check.replace(tzinfo=None)
            should_check = time_since_check >= timedelta(minutes=check_interval)
            
            if self.debug_mode:
                logging.debug(f"Service {service_type}: last_check={last_check}, interval={check_interval}min, "
                            f"time_since={time_since_check}, should_check={should_check}")
            
            return should_check
            
        except Exception as e:
            logging.error(f"Error checking if should monitor {service_type}: {e}")
            return True  # Default to checking if there's an error
    
    def monitor_service_type(self, service_type: str):
        """Monitor all documents for a specific service type."""
        config = self.get_monitor_config(service_type)
        documents = self.get_elastic_documents(service_type)
        
        if self.debug_level == 2:
            logging.info(f"📋 {service_type.upper()}: Found {len(documents)} documents")
        elif self.debug_level == 3:
            # Debug level 3: Only show service header if there will be problems to show
            pass  # We'll show the header later only if there are problems
        elif self.debug_level < 2:
            logging.info(f"🔍 Analyzing service type: {service_type}")
            logging.info(f"  📋 Found {len(documents)} elastic documents for {service_type}")
        
        if not documents:
            if self.debug_level < 2:
                logging.info(f"  ✅ No elastic documents found for {service_type} - nothing to monitor")
            self.update_monitor_status(service_type, True)
            return
        
        threshold_minutes = config.get('data_threshold_minutes', DEFAULT_DATA_THRESHOLD)
        restart_enabled = config.get('restart_enabled', True)
        alert_enabled = config.get('alert_enabled', False)
        alert_webhook_url = config.get('alert_webhook_url')
        timeout_seconds = config.get('request_timeout_seconds', DEFAULT_REQUEST_TIMEOUT)
        retry_attempts = config.get('retry_attempts', DEFAULT_RETRY_ATTEMPTS)
        retry_delay = config.get('retry_delay_seconds', DEFAULT_RETRY_DELAY)
        
        if self.debug_level < 2:
            logging.info(f"  ⚙️ Config - Threshold: {threshold_minutes}min, Timeout: {timeout_seconds}s, Restarts: {restart_enabled}")
        
        all_docs_ok = True
        error_messages = []
        docs_with_data = 0
        docs_without_data = 0
        service_header_shown = False  # Track if we've shown the service header for debug level 3
        service_problems = []  # Track problems for WhatsApp alerts
        
        for i, doc in enumerate(documents, 1):
            doc_name = doc['elastic_doc_name']
            client_id = doc['client_id']
            
            # Use default elastic credentials if not provided
            elastic_url = doc.get('elastic_url', 'https://opensearch.madd.com.mx:9200')
            elastic_user = doc.get('elastic_user', 'admin')
            elastic_password = doc.get('elastic_password', 'GPSc0ntr0l1')
            
            if self.debug_level < 3:
                logging.info(f"  📄 [{i}/{len(documents)}] Checking document: {doc_name} (client_id: {client_id})")
                logging.debug(f"       Elasticsearch URL: {elastic_url}")
            
            has_data, count, error_msg = self.check_elasticsearch_data(
                elastic_url, elastic_user, elastic_password, doc_name, 
                threshold_minutes, timeout_seconds, retry_attempts, retry_delay
            )
            
            if has_data:
                docs_with_data += 1
                if self.debug_level < 3:
                    logging.info(f"  ✅ {doc_name}: {count} records found in last {threshold_minutes} minutes")
            else:
                docs_without_data += 1
                all_docs_ok = False
                error_detail = f" ({error_msg})" if error_msg else ""
                
                # Add to problems list for WhatsApp alerts
                service_problems.append({
                    'doc_name': doc_name,
                    'error': error_msg if error_msg else "No data",
                    'count': count,
                    'client_id': client_id
                })
                
                # For debug level 3, show service header only when we have problems
                if self.debug_level == 3 and not service_header_shown:
                    logging.info(f"{service_type.upper()} service:")
                    service_header_shown = True
                
                if self.debug_level == 3:
                    # Simple format for debug level 3
                    if "HTTP 404" in error_msg:
                        logging.info(f"  {doc_name} (HTTP 404 - index not found)")
                    elif count == 0:
                        logging.info(f"  {doc_name} (0 records - NO DATA)")
                    elif "timeout" in error_msg.lower():
                        logging.info(f"  {doc_name} (timeout)")
                    elif "connection" in error_msg.lower():
                        logging.info(f"  {doc_name} (connection error)")
                    else:
                        logging.info(f"  {doc_name} ({error_msg})")
                else:
                    if self.debug_level < 3:
                        logging.info(f"  📄 [{i}/{len(documents)}] Checking document: {doc_name} (client_id: {client_id})")
                        logging.debug(f"       Elasticsearch URL: {elastic_url}")
                    
                    logging.warning(f"  ⛔ {doc_name}: No data in last {threshold_minutes} minutes{error_detail}")
                
                error_messages.append(f"{doc_name}: {error_msg}" if error_msg else f"{doc_name}: No data")
                
                # Get client namespace
                namespace = self.get_client_namespace(client_id)
                if not namespace:
                    if self.debug_level < 3:
                        logging.error(f"     🚨 Could not find namespace for client {client_id}")
                    continue
                
                if self.debug_level < 3:
                    logging.info(f"     🏷️ Client namespace: {namespace}")
                
                message = f"No data detected for {doc_name} (service: {service_type}, namespace: {namespace})"
                if error_msg:
                    message += f" - Error: {error_msg}"
                
                # Send alert if enabled
                if alert_enabled and alert_webhook_url:
                    if self.debug_level < 3:
                        logging.info(f"     📢 Sending alert to webhook")
                    alert_sent = self.send_alert(alert_webhook_url, message)
                    if alert_sent and self.debug_level < 3:
                        logging.info(f"     ✅ Alert sent successfully")
                    elif not alert_sent and self.debug_level < 3:
                        logging.error(f"     ⛔ Failed to send alert")
                else:
                    logging.debug(f"     📕 Alerts disabled for this service")
                
                # Restart service if enabled and it's not just a connection issue
                should_restart = restart_enabled and not ("timeout" in error_msg.lower() or "connection" in error_msg.lower())
                
                if should_restart:
                    if self.debug_level < 3:
                        logging.info(f"     🔄 Attempting to restart {service_type} in namespace {namespace}")
                    restart_success = self.restart_service(namespace, service_type)
                    
                    if restart_success:
                        if self.debug_level < 3:
                            logging.info(f"     ✅ Successfully triggered restart for {service_type} in {namespace}")
                    else:
                        if self.debug_level < 3:
                            logging.error(f"     ⛔ Failed to restart {service_type} in {namespace}")
                elif restart_enabled and self.debug_level < 3:
                    logging.info(f"     ⏸️ Restart skipped due to connection/timeout error")
                else:
                    logging.debug(f"     🚫 Restarts disabled for this service")
        
        # Store problems for WhatsApp alerts if any exist
        if service_problems:
            self.whatsapp_problems[service_type] = {
                'problems': service_problems,
                'count': docs_without_data,
                'timestamp': datetime.now()
            }
        
        # Summary for this service
        if self.debug_level < 3:
            # For debug level 3, don't show summary at all
            if not all_docs_ok:
                # Only show summary if there are problems for levels 1-2
                logging.info(f"  📊 Summary for {service_type}:")
                logging.info(f"     Documents with data: {docs_with_data}")
                logging.info(f"     Documents without data: {docs_without_data}")
                logging.info(f"     Overall status: {'✅ HEALTHY' if all_docs_ok else '⛔ ISSUES DETECTED'}")
        elif self.debug_level == 3 and service_header_shown:
            # For debug level 3, just add a blank line after problems
            logging.info("")
        
        # Update monitor status
        self.update_monitor_status(service_type, all_docs_ok, "; ".join(error_messages))
    
    def run_monitor_cycle(self):
        """Run one complete monitoring cycle."""
        if self.debug_level < 3:
            logging.info("=" * 60)
            logging.info("Starting monitoring cycle")
            logging.info("=" * 60)
        
        # Clear previous cycle's problems
        self.whatsapp_problems = {}
        
        services = self.get_services_with_elastic()
        if self.debug_level < 3:
            logging.info(f"Found {len(services)} services with elastic configuration: {services}")
        
        services_checked = 0
        services_skipped = 0
        
        for service_type in services:
            try:
                config = self.get_monitor_config(service_type)
                logging.debug(f"Config for {service_type}: {config}")
                
                # Check if it's time to monitor this service
                if self.should_check_service(service_type, config):
                    if self.debug_level < 3:
                        logging.info(f"🔍 Checking service: {service_type}")
                    self.monitor_service_type(service_type)
                    services_checked += 1
                else:
                    if self.debug_level < 3:
                        logging.info(f"⏳ Skipping {service_type} - not time to check yet")
                    services_skipped += 1
                
            except Exception as e:
                if self.debug_level < 3:
                    logging.error(f"Error monitoring service {service_type}: {e}")
                self.update_monitor_status(service_type, False, str(e))
        
        # Process WhatsApp alerts after all services have been checked
        if self.whatsapp_problems:
            if self.debug_level < 3:
                logging.info("📱 Processing WhatsApp alerts...")
            self.process_whatsapp_alerts()
        
        if self.debug_level < 3:
            logging.info("=" * 60)
            logging.info(f"Monitoring cycle completed - Checked: {services_checked}, Skipped: {services_skipped}")
            logging.info("=" * 60)
    
    def run(self):
        """Main monitoring loop."""
        logging.info("Starting Elasticsearch Monitor")
        
        while True:
            try:
                self.run_monitor_cycle()
                
                # Wait 1 minute before next cycle (services have individual intervals)
                logging.info("Waiting 1 minute before next cycle...")
                time.sleep(60)  # Check every minute, but services have their own intervals
                
            except KeyboardInterrupt:
                logging.info("Monitor stopped by user")
                break
            except Exception as e:
                logging.error(f"Unexpected error in monitor loop: {e}")
                time.sleep(60)  # Wait 1 minute before retrying

if __name__ == '__main__':
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Elasticsearch Service Monitor')
    parser.add_argument('-d', '--debug', type=int, nargs='?', const=1, default=0,
                        help='Enable debug mode (1=full debug, 2=elasticsearch only, 3=minimal problems only)')
    parser.add_argument('--once', action='store_true',
                        help='Run monitoring cycle once and exit (useful for testing)')
    parser.add_argument('--force', action='store_true',
                        help='Force WhatsApp alerts even if timing interval not met (useful for testing)')
    
    args = parser.parse_args()
    
    # Create monitor with debug mode if specified
    debug_mode = args.debug > 0
    debug_level = args.debug
    force_alerts = args.force
    monitor = ElasticsearchMonitor(debug_mode=debug_mode, debug_level=debug_level, force_alerts=force_alerts)
    
    if debug_mode:
        if debug_level == 2:
            logging.info("🔍 Debug Level 2: Showing only Elasticsearch queries")
        elif debug_level == 3:
            pass  # No verbose header for minimal output
        else:
            logging.info("Debug mode enabled - detailed logging active")
    
    if force_alerts:
        logging.info("🚀 FORCE MODE: WhatsApp alerts will bypass timing restrictions")
    
    if args.once:
        logging.info("Running single monitoring cycle...")
        monitor.run_monitor_cycle()
        logging.info("Single cycle completed. Exiting.")
    else:
        monitor.run()