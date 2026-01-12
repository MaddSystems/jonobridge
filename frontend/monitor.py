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
    def __init__(self, debug_mode=False, debug_level=1):
        self.debug_mode = debug_mode
        self.debug_level = debug_level if debug_mode else 0
        self.db_config = {
            'host': 'localhost',
            'user': 'gpscontrol',
            'password': 'qazwsxedc',
            'database': 'automation_app'
        }
        
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
                        logging.info(f"    📝 Original: {repr(doc_name)} → Normalized: {repr(normalized_doc_name)}")
                    logging.info(f"    🌐 URL: {url}")
                    logging.info(f"    🔍 Query: last {threshold_minutes} minutes")
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
                        logging.info(f"    ❌ Error: HTTP {response.status_code}")
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
            logging.info(f"  ⚙️  Config - Threshold: {threshold_minutes}min, Timeout: {timeout_seconds}s, Restarts: {restart_enabled}")
        
        all_docs_ok = True
        error_messages = []
        docs_with_data = 0
        docs_without_data = 0
        service_header_shown = False  # Track if we've shown the service header for debug level 3
        
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
                    
                    logging.warning(f"  ❌ {doc_name}: No data in last {threshold_minutes} minutes{error_detail}")
                
                error_messages.append(f"{doc_name}: {error_msg}" if error_msg else f"{doc_name}: No data")
                
                # Get client namespace
                namespace = self.get_client_namespace(client_id)
                if not namespace:
                    if self.debug_level < 3:
                        logging.error(f"     🚨 Could not find namespace for client {client_id}")
                    continue
                
                if self.debug_level < 3:
                    logging.info(f"     🏷️  Client namespace: {namespace}")
                
                message = f"No data detected for {doc_name} (service: {service_type}, namespace: {namespace})"
                if error_msg:
                    message += f" - Error: {error_msg}"
                
                # Send alert if enabled
                if alert_enabled and alert_webhook_url:
                    if self.debug_level < 3:
                        logging.info(f"     🔔 Sending alert to webhook")
                    alert_sent = self.send_alert(alert_webhook_url, message)
                    if alert_sent and self.debug_level < 3:
                        logging.info(f"     ✅ Alert sent successfully")
                    elif not alert_sent and self.debug_level < 3:
                        logging.error(f"     ❌ Failed to send alert")
                else:
                    logging.debug(f"     🔕 Alerts disabled for this service")
                
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
                            logging.error(f"     ❌ Failed to restart {service_type} in {namespace}")
                elif restart_enabled and self.debug_level < 3:
                    logging.info(f"     ⏸️  Restart skipped due to connection/timeout error")
                else:
                    logging.debug(f"     🚫 Restarts disabled for this service")
        
        # Summary for this service
        if self.debug_level < 3:
            # For debug level 3, don't show summary at all
            if not all_docs_ok:
                # Only show summary if there are problems for levels 1-2
                logging.info(f"  📊 Summary for {service_type}:")
                logging.info(f"     Documents with data: {docs_with_data}")
                logging.info(f"     Documents without data: {docs_without_data}")
                logging.info(f"     Overall status: {'✅ HEALTHY' if all_docs_ok else '❌ ISSUES DETECTED'}")
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
                        logging.info(f"📊 Checking service: {service_type}")
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
                        help='Enable debug mode (1=full debug, 2=elasticsearch only, default=1)')
    parser.add_argument('--once', action='store_true',
                        help='Run monitoring cycle once and exit (useful for testing)')
    
    args = parser.parse_args()
    
    # Create monitor with debug mode if specified
    debug_mode = args.debug > 0
    debug_level = args.debug
    monitor = ElasticsearchMonitor(debug_mode=debug_mode, debug_level=debug_level)
    
    if debug_mode:
        if debug_level == 2:
            logging.info("🔍 Debug Level 2: Showing only Elasticsearch queries")
        elif debug_level == 3:
            pass  # No verbose header for minimal output
        else:
            logging.info("Debug mode enabled - detailed logging active")
    
    if args.once:
        logging.info("Running single monitoring cycle...")
        monitor.run_monitor_cycle()
        logging.info("Single cycle completed. Exiting.")
    else:
        monitor.run()