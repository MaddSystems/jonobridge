from flask import Flask, request, render_template, redirect, jsonify, url_for, flash, session, send_file
import mysql.connector
import json
import subprocess
import yaml
from services import ServiceFactory
import os
from kubernetes import client as k8s_client, config
import traceback
import sys
import time
from datetime import datetime, timedelta
import requests
from functools import wraps
import binascii
import tempfile
import zipfile
import shutil
import mysql.connector
import time
import os
import sys
import traceback
import json
import requests
import subprocess

# Import Swagger/Flasgger for API documentation
from flasgger import Swagger

# Import BSJ packet generator
from bsj_packet_generator import generate_bsj_command

# Import diagnostic API blueprint
from api.diagnostic_api import diagnostic_bp

app = Flask(__name__, static_url_path='/static', static_folder='static')
app.secret_key = 'your_secret_key_here'

# Add CORS headers and disable range requests
@app.after_request
def add_cors_headers(response):
    response.headers['Access-Control-Allow-Origin'] = '*'
    response.headers['Access-Control-Allow-Methods'] = 'GET, POST, PUT, DELETE, OPTIONS'
    response.headers['Access-Control-Allow-Headers'] = 'Content-Type, Authorization'
    response.headers['Accept-Ranges'] = 'none'
    return response

@app.route('/api/v1/diagnostic/<path:path>', methods=['OPTIONS'])
def handle_options(path):
    """Handle CORS preflight requests"""
    return '', 200

# Initialize Swagger/Flasgger for API documentation
swagger = Swagger(app, 
    template={
        "swagger": "2.0",
        "info": {
            "title": "JonoBridge Diagnostic API",
            "description": "API for monitoring and managing Kubernetes pods and namespaces",
            "contact": {
                "name": "MADD Systems"
            },
            "version": "1.0.0"
        },
        "basePath": "/api/v1",
        "schemes": ["https", "http"]
    }
)

# Register diagnostic API blueprint
app.register_blueprint(diagnostic_bp)

# User credentials (in a real application, these should be stored securely in a database)
USERS = {
    'admin': 'password123',
    'jorge': 'qazwsxedc',
    'soporte': 'GPSc0ntr0l1'
}

def get_minikube_ip():
    """Get and verify Minikube IP address."""
    try:
        print("Environment variables:")
        print(os.environ)

        # Debugging Minikube Command
        print("Running minikube ip command...")
        result = subprocess.run(
            ['minikube', 'ip'],
            capture_output=True,
            text=True,
            check=True
        )
        minikube_ip = result.stdout.strip()
        print(f"Minikube IP: {minikube_ip}")

        if not minikube_ip:
            raise subprocess.CalledProcessError(1, 'minikube ip', 'Empty IP returned')

        subprocess.run(['ping', '-c', '1', minikube_ip], check=True, capture_output=True)
        print(f"Minikube IP {minikube_ip} is accessible")
        return minikube_ip
    except subprocess.CalledProcessError as e:
        print(f"Command failed: {e.cmd}")
        print(f"Return code: {e.returncode}")
        print(f"Output: {e.output}")
        print(f"Stderr: {e.stderr}")
        sys.exit(1)

# Get Minikube IP at startup
MINIKUBE_IP = get_minikube_ip()

def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if 'username' not in session:
            return redirect(url_for('index'))
        return f(*args, **kwargs)
    return decorated_function

@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        username = request.form.get('username')
        password = request.form.get('password')
        
        if username in USERS and USERS[username] == password:
            session['username'] = username
            return redirect(url_for('clients'))
        else:
            flash('Invalid username or password', 'danger')
            return redirect(url_for('index'))
    
    if 'username' in session:
        return redirect(url_for('clients'))
        
    return render_template('index.html')

@app.route('/clients')
@login_required
def clients():
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients')
        clients = cursor.fetchall()
        cursor.close()
        db.close()
        return render_template('clients.html', clients=clients)
    except Exception as e:
        flash(f'Error: {str(e)}', 'error')
        return render_template('clients.html', clients=[])

@app.route('/logout')
def logout():
    session.clear()
    return redirect(url_for('index'))

MOSQUITTO_TEMPLATE = {
    'configmap': {
        'apiVersion': 'v1',
        'kind': 'ConfigMap',
        'metadata': {
            'name': 'mosquitto-config',
            'namespace': None
        },
        'data': {
            'mosquitto.conf': '''
listener 1883
allow_anonymous true
max_queued_messages 1000
max_inflight_messages 100
persistence true
persistence_location /mosquitto/data/
autosave_interval 60
queue_qos0_messages true
'''
        }
    },
    'deployment': {
        'apiVersion': 'apps/v1',
        'kind': 'Deployment',
        'metadata': {
            'name': 'mosquitto',
            'namespace': None
        },
        'spec': {
            'replicas': 1,
            'selector': {
                'matchLabels': {
                    'app': 'mosquitto'
                }
            },
            'template': {
                'metadata': {
                    'labels': {
                        'app': 'mosquitto'
                    }
                },
                'spec': {
                    'containers': [{
                        'name': 'mosquitto',
                        'image': 'eclipse-mosquitto:latest',
                        'ports': [{
                            'containerPort': 1883
                        }],
                        'volumeMounts': [{
                            'name': 'config',
                            'mountPath': '/mosquitto/config'
                        }]
                    }],
                    'volumes': [{
                        'name': 'config',
                        'configMap': {
                            'name': 'mosquitto-config'
                        }
                    }]
                }
            }
        }
    },
    'service': {
        'apiVersion': 'v1',
        'kind': 'Service',
        'metadata': {
            'name': 'mosquitto',
            'namespace': None
        },
        'spec': {
            'type': 'NodePort',
            'selector': {
            'app': 'mosquitto'
            },
            'ports': [
            {
                'name': 'mqtt',
                'port': 1883,
                'targetPort': 1883,
            },
            {
                'name': 'mqtt-websocket',
                'port': 9001,
                'targetPort': 9001,
            }
            ]
        }
    }
}

def check_and_create_missing_tables():
    """Check for missing tables and create them if needed."""
    try:
        db = get_db_connection()
        cursor = db.cursor()

        # Get required tables from schema
        required_tables = get_required_tables()

        # Get existing tables
        cursor.execute("SHOW TABLES")
        existing_tables = {table[0] for table in cursor.fetchall()}

        # Find and create missing tables
        missing_tables = []
        for table_name, create_statement in required_tables.items():
            if table_name not in existing_tables:
                missing_tables.append(table_name)
                print(f"Creating missing table: {table_name}")
                cursor.execute(create_statement)

        db.commit()
        return missing_tables

    except Exception as e:
        print(f"Error checking/creating tables: {str(e)}")
        raise
    finally:
        cursor.close()
        db.close()

def get_db_connection():
    try:
        connection = mysql.connector.connect(
            host=DB_CONFIG['host'],
            user=DB_CONFIG['user'],
            password=DB_CONFIG['password'],
            database=DB_NAME
        )
        return connection
    except mysql.connector.Error as err:
        print(f"Error connecting to database: {err}")
        raise

def get_service_data(client_id: int, action: str) -> dict:
    """Get service data for a client."""
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        services_data = {}
        
        # Get all available services from the factory
        available_services = service_factory.get_all_services()
        
        # For each service, check if client has configuration
        for service_name, service_instance in available_services.items():
            table_name = service_instance.service_name
            print(f"Checking {table_name} configuration for client {client_id}")
            
            try:
                cursor.execute(f'SELECT * FROM {table_name} WHERE client_id = %s', (client_id,))
                service_data = cursor.fetchone()
                
                if service_data:
                    # Remove id and client_id from service data
                    service_params = {k: v for k, v in service_data.items() 
                                  if k not in ['id', 'client_id', 'created_at']}
                    services_data[service_name] = service_params
                    print(f"Found {service_name} configuration:", service_params)
            except Exception as e:
                print(f"Error checking {service_name}: {e}")
        
        cursor.close()
        db.close()

        if not services_data:
            print(f"No services found for client {client_id}")
            return {}

        print(f"Retrieved services for client {client_id}:", services_data)
        return services_data
        
    except Exception as e:
        print(f"Error getting service data: {e}")
        return {}

def generate_kubernetes_manifests(namespace: str, service_name: str, service_data: dict) -> dict:
    """Generate Kubernetes manifests for a service."""
    try:
        service = service_factory.get_all_services().get(service_name)
        if not service:
            print(f"Service {service_name} not found")
            return {}

        return service.get_kubernetes_manifests(namespace, service_data)
    except Exception as e:
        print(f"Error generating manifests for {service_name}: {e}")
        return {}

def configure_nginx_stream(namespace, client_id):
    """Configure nginx stream for a namespace with improved error handling and validation."""
    print(f"\n=== Starting nginx configuration for client: {namespace} ===")
    results = []
    
    try:
        # Verify nginx service is running
        print("1. Checking nginx service status...")
        nginx_status = subprocess.run(['sudo', 'systemctl', 'is-active', 'nginx'], 
                                    capture_output=True, text=True)
        print(f"Nginx status: {nginx_status.stdout.strip()}")
        if nginx_status.returncode != 0:
            error_msg = "Nginx service is not running"
            results.append({'command': 'nginx-status', 'success': False, 'output': error_msg})
            return False, results

        # Create streams.d directory with proper permissions
        streams_dir = '/etc/nginx/streams.d'
        print(f"2. Setting up streams directory: {streams_dir}")
        
        try:
            # Check if directory exists
            dir_exists = os.path.exists(streams_dir)
            print(f"Streams directory exists: {dir_exists}")
            
            if not dir_exists:
                print("Creating streams directory...")
                subprocess.run(['sudo', 'mkdir', '-p', streams_dir], check=True)
                print("Setting directory permissions...")
                subprocess.run(['sudo', 'chown', 'root:root', streams_dir], check=True)
                subprocess.run(['sudo', 'chmod', '755', streams_dir], check=True)
            
            results.append({'command': 'create-streams-dir', 'success': True, 
                          'output': f"Streams directory ready: {streams_dir}"})
            
        except Exception as e:
            error_msg = f"Failed to create/setup streams directory: {str(e)}"
            print(error_msg)
            results.append({'command': 'create-streams-dir', 'success': False, 'output': error_msg})
            return False, results

        # Load kubernetes configuration
        print("3. Loading kubernetes configuration...")
        config.load_kube_config()
        v1 = k8s_client.CoreV1Api()
        
        # Get services and find bridge service
        print("4. Finding proxy service...")
        services = v1.list_namespaced_service(namespace=namespace)
        proxy_service = None

        for svc in services.items:
            if svc.metadata.name.endswith('listener-service'):
                proxy_service = svc
                break
        
        if not proxy_service:
            error_msg = f"No proxy service found for {namespace}"
            print(error_msg)
            results.append({'command': 'find-proxy-service', 'success': False, 'output': error_msg})
            return True, results
        
        print(f"Found proxy service: {proxy_service.metadata.name}")
        print("SVC:", svc.metadata.name)
        print("client_id:", client_id)  
        print("namespace:", namespace)

        # Set timeout to 1 minute from now
        timeout = datetime.now() + timedelta(minutes=1)
        nodeport = None
        nodeport_api = None

        while datetime.now() < timeout:
            # Get main TCP port
            nodeport, port = get_nodeport_and_tcp(svc.metadata.name, namespace, 'tcp')
            
            # Get TCP-API port if available
            nodeport_api_result = get_nodeport_and_tcp(svc.metadata.name, namespace, 'tcp-api')
            port_api, nodeport_api = nodeport_api_result if nodeport_api_result else (None, None)
            
            if isinstance(nodeport, int):
                print(f"Successfully got nodeport: {nodeport}")
                break
            print("Waiting for NodePort to be available...")
            time.sleep(5)

        if not isinstance(nodeport, int):
            print("Timeout reached while waiting for NodePort")
            return False, results

        print(f"Main Nodeport:{nodeport}, Main port:{port}")
        if nodeport_api:
            print(f"API Nodeport:{nodeport_api}, API port:{port_api}")
        
        # Get port mappings
        print("5. Getting port mappings...")
        port_info = next((port for port in proxy_service.spec.ports), None)
        if not port_info or not port_info.node_port:
            error_msg = f"Could not find port mappings for service {proxy_service.metadata.name}"
            print(error_msg)
            results.append({'command': 'get-ports', 'success': False, 'output': error_msg})
            return False, results
        
        print(f"Port mappings - Port: {port_info.port}, NodePort: {port_info.node_port}")

        # Verify Minikube IP
        print("6. Verifying Minikube IP...")
        results.append({'command': 'verify-minikube', 'success': True, 
                      'output': f"Minikube IP {MINIKUBE_IP} is accessible"})

        # Create or update nginx configuration
        print("7. Creating/updating nginx configuration...")
        nginx_file_path = f'/etc/nginx/streams.d/port-{port_info.port}.conf'
        temp_path = f'/tmp/port-{port_info.port}.conf'

        existing_servers = set()  # Use a set to prevent duplicates
        if os.path.exists(nginx_file_path):
            # Read existing configuration to extract current servers
            with open(nginx_file_path, 'r') as f:
                for line in f:
                    line = line.strip()
                    if 'server' in line and MINIKUBE_IP in line:
                        existing_servers.add(line)

        # Add new server if not already present
        new_server = f"    server {MINIKUBE_IP}:{port_info.node_port};"
        existing_servers.add(new_server)

        # Convert set back to sorted list for consistent output
        server_list = sorted(list(existing_servers))

        # Create updated nginx configuration
        nginx_config = f"""upstream port{port_info.port}_tcp {{
{chr(10).join(server_list)}
}}
upstream port{port_info.port}_udp {{
{chr(10).join(server_list)}
}}
server {{
    listen {port_info.port};
    proxy_pass port{port_info.port}_tcp;
    proxy_timeout 10m;
    proxy_connect_timeout 1m;
}}
server {{
    listen {port_info.port} udp;
    proxy_pass port{port_info.port}_udp;
    proxy_timeout 10m;
    proxy_connect_timeout 1m;
}}
"""

        # Add API port configuration if available
        if nodeport_api and port_api:
            api_server = f"    server {MINIKUBE_IP}:{nodeport_api};"
            nginx_config += f"""
upstream port{port_api}_tcp {{
    {api_server}
}}
server {{
    listen {port_api};
    proxy_pass port{port_api}_tcp;
    proxy_timeout 10m;
    proxy_connect_timeout 1m;
}}
"""
        print("Configuration content:")
        print(nginx_config)
        
        print(f"8. Writing configuration to temporary file: {temp_path}")
        with open(temp_path, 'w') as f:
            f.write(nginx_config)
        
        print(f"9. Moving configuration to final location: {nginx_file_path}")
        try:
            subprocess.run(['sudo', 'mv', temp_path, nginx_file_path], check=True)
            subprocess.run(['sudo', 'chown', 'root:root', nginx_file_path], check=True)
            subprocess.run(['sudo', 'chmod', '644', nginx_file_path], check=True)
            
            if not os.path.exists(nginx_file_path):
                raise Exception(f"File was not created at {nginx_file_path}")
                
            print(f"Configuration file created successfully")
            results.append({'command': 'write-config', 'success': True, 
                          'output': f"Created configuration file {nginx_file_path}"})
        except Exception as e:
            error_msg = f"Failed to create configuration file: {str(e)}"
            print(error_msg)
            results.append({'command': 'write-config', 'success': False, 'output': error_msg})
            return False, results

        # Test configuration
        print("10. Testing nginx configuration...")
        test_result = subprocess.run(['sudo', 'nginx', '-t'], capture_output=True, text=True)
        if test_result.returncode != 0:
            error_msg = f"Nginx configuration test failed: {test_result.stderr}"
            print(error_msg)
            results.append({'command': 'test-config', 'success': False, 'output': error_msg})
            subprocess.run(['sudo', 'rm', nginx_file_path], check=False)
            return False, results
        
        print("Configuration test passed")
        results.append({'command': 'test-config', 'success': True, 'output': "Configuration test passed"})
        # Setting api port
        
        # Reload nginx
        print("11. Reloading nginx...")
        reload_result = subprocess.run(['sudo', 'systemctl', 'reload', 'nginx'], capture_output=True, text=True)
        if reload_result.returncode != 0:
            error_msg = f"Nginx reload failed: {reload_result.stderr}"
            print(error_msg)
            results.append({'command': 'reload-nginx', 'success': False, 'output': error_msg})
            subprocess.run(['sudo', 'rm', nginx_file_path], check=False)
            return False, results

        print("Nginx reloaded successfully")
        results.append({'command': 'reload-nginx', 'success': True, 'output': "Nginx reloaded successfully"})
        
        # Update port and api_port in the clients table
        print("12. Updating client table with port and API port...")
        try:
            db = get_db_connection()
            cursor = db.cursor()
            # Check if api_port is None and set default value to 9 if it is
            if port_api is None:
                port_api = 0
                print(f"API port is null, setting default value: {port_api}")
            
            update_query = "UPDATE clients SET port = %s, api_port = %s WHERE id = %s"
            cursor.execute(update_query, (port_info.port, port_api, client_id))
            db.commit()
            cursor.close()
            db.close()
            print(f"Client table updated with port={port_info.port}, api_port={port_api}")
            results.append({'command': 'update-client-table', 'success': True, 
                           'output': f"Client ports updated: port={port_info.port}, api_port={port_api}"})
        except Exception as e:
            error_msg = f"Failed to update client table: {str(e)}"
            print(error_msg)
            results.append({'command': 'update-client-table', 'success': False, 'output': error_msg})
            # Continue execution even if update fails
        
        print("\n=== Nginx configuration completed successfully ===")
        return True, results

    except Exception as e:
        error_msg = f"Error configuring nginx stream: {str(e)}"
        print(error_msg)
        traceback.print_exc()
        results.append({'command': 'configure-nginx', 'success': False, 'output': error_msg})
        return False, results

def configure_nginx_endpoint(namespace, client_id):
    """Configure nginx endpoint for portal service in a namespace."""
    print(f"\n=== Starting nginx endpoint configuration for client: {namespace} ===")
    results = []
    
    try:
        # Verify nginx service is running
        print("1. Checking nginx service status...")
        nginx_status = subprocess.run(['sudo', 'systemctl', 'is-active', 'nginx'], 
                                    capture_output=True, text=True)
        print(f"Nginx status: {nginx_status.stdout.strip()}")
        if nginx_status.returncode != 0:
            error_msg = "Nginx service is not running"
            results.append({'command': 'nginx-status', 'success': False, 'output': error_msg})
            return False, results

        # Create endpoints directory with proper permissions
        endpoints_dir = '/etc/nginx/endpoints'
        print(f"2. Setting up endpoints directory: {endpoints_dir}")
        
        try:
            # Check if directory exists
            main_dir_exists = os.path.exists(endpoints_dir)
            print(f"Endpoints directory exists: {main_dir_exists}")
            
            if not main_dir_exists:
                print("Creating endpoints directory...")
                subprocess.run(['sudo', 'mkdir', '-p', endpoints_dir], check=True)
                subprocess.run(['sudo', 'chown', 'root:root', endpoints_dir], check=True)
                subprocess.run(['sudo', 'chmod', '755', endpoints_dir], check=True)
            
            results.append({'command': 'create-endpoints-dir', 'success': True, 
                          'output': f"Endpoints directory ready: {endpoints_dir}"})
            
        except Exception as e:
            error_msg = f"Failed to create/setup endpoints directory: {str(e)}"
            print(error_msg)
            results.append({'command': 'create-endpoints-dir', 'success': False, 'output': error_msg})
            return False, results

        # Load kubernetes configuration
        print("3. Loading kubernetes configuration...")
        config.load_kube_config()
        v1 = k8s_client.CoreV1Api()
        
        # Get services and find portal service
        print("4. Finding portal service...")
        services = v1.list_namespaced_service(namespace=namespace)
        portal_service = None

        for svc in services.items:
            if svc.metadata.name == 'portal' or svc.metadata.name == 'webrelay' or svc.metadata.name == 'httpinput' or svc.metadata.name == 'websqldump' or svc.metadata.name == 'grule':
                portal_service = svc
                break
        
        if not portal_service:
            print("No portal service found, skipping endpoint configuration")
            results.append({'command': 'fsourind-portal-service', 'success': True, 'output': "No portal service found"})
            return True, results
        
        print(f"Found portal service: {portal_service.metadata.name}")

        # Set timeout to 1 minute from now
        timeout = datetime.now() + timedelta(minutes=1)
        nodeport = None

        while datetime.now() < timeout:
            # Find the NodePort
            for port in portal_service.spec.ports:
                if port.node_port:
                    nodeport = port.node_port
                    service_port = port.port
                    break
            
            if nodeport:
                print(f"Found portal NodePort: {nodeport}")
                break
                
            print("Waiting for NodePort to be available...")
            time.sleep(5)

        if not nodeport:
            print("Timeout reached while waiting for NodePort")
            return False, results

        # Get port mappings and endpoint from database
        print("5. Getting portal configuration...")
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        service_name = svc.metadata.name
        if service_name == 'portal' or service_name == 'webrelay' or service_name == 'httpinput' or service_name == 'websqldump' or service_name== 'grule':
            try:
                cursor.execute(f"SELECT portal_endpoint FROM {service_name} WHERE client_id = %s", (client_id,))
                portal_config = cursor.fetchone()
                
                print(f"Retrieved portal config: {portal_config}")
                
                if not portal_config or not portal_config.get('portal_endpoint'):
                    error_msg = "Portal endpoint not configured in database"
                    print(error_msg)
                    results.append({'command': 'get-endpoint', 'success': False, 'output': error_msg})
                    cursor.close()
                    db.close()
                    return False, results
                
                # Extract portal_endpoint from configuration
                portal_endpoint = portal_config.get('portal_endpoint')
                print(f"Using portal endpoint: {portal_endpoint}")
                
                # Ensure we have a valid portal endpoint before continuing
                if not portal_endpoint:
                    error_msg = "Portal endpoint is empty even though it exists"
                    print(error_msg)
                    results.append({'command': 'get-endpoint', 'success': False, 'output': error_msg})
                    cursor.close()
                    db.close()
                    return False, results
                    
            except Exception as e:
                error_msg = f"Failed to get portal configuration: {str(e)}"
                print(error_msg)
                results.append({'command': 'get-endpoint', 'success': False, 'output': error_msg})
                cursor.close()
                db.close()
                return False, results

        # Verify Minikube IP
        print("6. Verifying Minikube IP...")
        results.append({'command': 'verify-minikube', 'success': True, 
                      'output': f"Minikube IP {MINIKUBE_IP} is accessible"})

        # Create nginx configuration
        print("7. Creating nginx endpoint configuration...")
        # Store file directly in the endpoints directory without subdirectory
        nginx_file_path = f'{endpoints_dir}/{namespace}.conf'
        temp_path = f'/tmp/{namespace}_endpoint.conf'

        # Create the endpoint configuration
        nginx_config = f"""location /{portal_endpoint} {{
    proxy_pass http://{MINIKUBE_IP}:{nodeport}/{portal_endpoint};
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}}
"""

        print("Configuration content:")
        print(nginx_config)
        
        print(f"8. Writing configuration to temporary file: {temp_path}")
        with open(temp_path, 'w') as f:
            f.write(nginx_config)
        
        print(f"9. Moving configuration to final location: {nginx_file_path}")
        try:
            subprocess.run(['sudo', 'mv', temp_path, nginx_file_path], check=True)
            subprocess.run(['sudo', 'chown', 'root:root', nginx_file_path], check=True)
            subprocess.run(['sudo', 'chmod', '644', nginx_file_path], check=True)
            
            if not os.path.exists(nginx_file_path):
                raise Exception(f"File was not created at {nginx_file_path}")
                
            print(f"Configuration file created successfully")
            results.append({'command': 'write-config', 'success': True, 
                          'output': f"Created configuration file {nginx_file_path}"})
        except Exception as e:
            error_msg = f"Failed to create configuration file: {str(e)}"
            print(error_msg)
            results.append({'command': 'write-config', 'success': False, 'output': error_msg})
            return False, results

        # Test configuration
        print("10. Testing nginx configuration...")
        test_result = subprocess.run(['sudo', 'nginx', '-t'], capture_output=True, text=True)
        if test_result.returncode != 0:
            error_msg = f"Nginx configuration test failed: {test_result.stderr}"
            print(error_msg)
            results.append({'command': 'test-config', 'success': False, 'output': error_msg})
            subprocess.run(['sudo', 'rm', nginx_file_path], check=False)
            return False, results
        
        print("Configuration test passed")
        results.append({'command': 'test-config', 'success': True, 'output': "Configuration test passed"})

        # Reload nginx
        print("11. Reloading nginx...")
        reload_result = subprocess.run(['sudo', 'systemctl', 'reload', 'nginx'], capture_output=True, text=True)
        if reload_result.returncode != 0:
            error_msg = f"Nginx reload failed: {reload_result.stderr}"
            print(error_msg)
            results.append({'command': 'reload-nginx', 'success': False, 'output': error_msg})
            subprocess.run(['sudo', 'rm', nginx_file_path], check=False)
            return False, results

        print("Nginx reloaded successfully")
        results.append({'command': 'reload-nginx', 'success': True, 'output': "Nginx reloaded successfully"})
        
        # Update port and api_port in the clients table
        print("12. Updating client table with port and API port...")
        try:
            db = get_db_connection()
            cursor = db.cursor()
            # Check if api_port is None and set default value to 9 if it is
            if port_api is None:
                port_api = 0
                print(f"API port is null, setting default value: {port_api}")
            
            update_query = "UPDATE clients SET port = %s, api_port = %s WHERE id = %s"
            cursor.execute(update_query, (port_info.port, port_api, client_id))
            db.commit()
            cursor.close()
            db.close()
            print(f"Client table updated with port={port_info.port}, api_port={port_api}")
            results.append({'command': 'update-client-table', 'success': True, 
                           'output': f"Client ports updated: port={port_info.port}, api_port={port_api}"})
        except Exception as e:
            error_msg = f"Failed to update client table: {str(e)}"
            print(error_msg)
            results.append({'command': 'update-client-table', 'success': False, 'output': error_msg})
            # Continue execution even if update fails
        
        # Update web_port in the clients table
        print("12. Updating client table with web port...")
        try:
            # Update query to set web_port
            update_query = "UPDATE clients SET web_port = %s WHERE id = %s"
            cursor.execute(update_query, (nodeport, client_id))
            db.commit()
            print(f"Client table updated with web_port={nodeport}")
            results.append({'command': 'update-web-port', 'success': True, 
                          'output': f"Client web port updated: web_port={nodeport}"})
        except Exception as e:
            error_msg = f"Failed to update client web port: {str(e)}"
            print(error_msg)
            results.append({'command': 'update-web-port', 'success': False, 'output': error_msg})
            # Continue execution even if update fails
            
        # Close database connection that was opened earlier
        cursor.close()
        db.close()
        
        print("\n=== Nginx endpoint configuration completed successfully ===")
        return True, results

    except Exception as e:
        error_msg = f"Error configuring nginx endpoint: {str(e)}"
        print(error_msg)
        traceback.print_exc()
        results.append({'command': 'configure-nginx-endpoint', 'success': False, 'output': error_msg})
        return False, results

def cleanup_nginx_config(namespace: str) -> None:
    """Clean up nginx configuration for a namespace."""
    print(f"\n=== Cleaning up nginx configuration for {namespace} ===")
    try:
        # Load kubernetes configuration to get the port information
        config.load_kube_config()
        v1 = k8s_client.CoreV1Api()
        
        # Get services to find the port
        services = v1.list_namespaced_service(namespace=namespace)
        port = None
        
        for svc in services.items:
            if svc.metadata.name.endswith('listener-service'):
                port_info = next((port for port in svc.spec.ports), None)
                if port_info:
                    port = port_info.port
                break
        
        if not port:
            print(f"No port found for namespace {namespace}, skipping stream cleanup")
        else:
            # Clean up in streams.d directory
            nginx_streams_path = f"/etc/nginx/streams.d/port-{port}.conf"
            if os.path.exists(nginx_streams_path):
                print(f"Removing nginx stream configuration: {nginx_streams_path}")
                subprocess.run(['sudo', 'rm', nginx_streams_path], check=True)
                print(f"Successfully cleaned up nginx stream configuration for {namespace}")
                
        # Clean up portal endpoint configuration - now directly in endpoints directory
        nginx_endpoints_path = f"/etc/nginx/endpoints/{namespace}.conf"
        
        if os.path.exists(nginx_endpoints_path):
            print(f"Removing nginx endpoint configuration: {nginx_endpoints_path}")
            subprocess.run(['sudo', 'rm', nginx_endpoints_path], check=True)
            
        # Test and reload nginx configuration
        print("Testing nginx configuration...")
        test_result = subprocess.run(['sudo', 'nginx', '-t'], capture_output=True, text=True)
        if test_result.returncode == 0:
            print("Reloading nginx...")
            subprocess.run(['sudo', 'systemctl', 'reload', 'nginx'], check=True)
            print(f"Successfully cleaned up all nginx configurations for {namespace}")
        else:
            print(f"Nginx configuration test failed after cleanup: {test_result.stderr}")
            
    except Exception as e:
        print(f"Error cleaning up nginx config: {e}")
        raise

@app.route('/status/<int:client_id>')
@login_required
def status_page(client_id):
    db = get_db_connection()
    cursor = db.cursor(dictionary=True)
    cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
    client = cursor.fetchone()
    cursor.close()
    db.close()
    
    if not client:
        flash('Client not found', 'error')
        return redirect(url_for('index'))
        
    return render_template('status.html', client=client)

@app.route('/setup/<int:client_id>')
@login_required
def setup_page(client_id):
    db = get_db_connection()
    cursor = db.cursor(dictionary=True)
    cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
    client = cursor.fetchone()

    if not client:
        flash('Client not found', 'error')
        return redirect(url_for('index'))

    available_services = service_factory.get_all_services()
    setup_inputs = []
    setup_interpreters = []
    setup_integrations = []
    client_nodes = {}
    structure = {}
    for service_name, service_instance in available_services.items():
        service_info = {
            "service_name": service_instance.service_name,
            "parameters": service_instance.parameters,
            "parameters_helpers": service_instance.parameters_helpers,
            "service_inputs": service_instance.inputs,
            "services_outputs": service_instance.outputs
        }
        
        if service_instance.service_type == "input":
            setup_inputs.append(service_info)
            structure
        elif service_instance.service_type == "interpreter":
            setup_interpreters.append(service_info)
        elif service_instance.service_type == "integration":
            setup_integrations.append(service_info)
        try:
            query=f'SELECT * FROM {service_name} WHERE client_id = {client_id}'
            print("query:   ", query)
            cursor.execute(query)
            service_data = cursor.fetchone()
            if service_data:
                # Remove id and client_id from service data
                cleaned_data = {k: v for k, v in service_data.items() 
                                if k not in ['id', 'client_id', 'created_at']}
                client_nodes[service_name] = {
                    'parameters': cleaned_data,
                    'service_name': service_name
                }
                print(f"Loaded node data for {service_name}:", client_nodes[service_name])
        except Exception as e:
            print(f"Error fetching {service_name} config: {str(e)}")
            continue

    # Get existing connections
    cursor.execute('''
        SELECT source_service, target_service, source_id, target_id 
        FROM service_connections 
        WHERE client_id = %s
    ''', (client_id,))
    connections = cursor.fetchall()
    
    # Convert connections to serializable format
    serializable_connections = []
    for conn in connections:
        serializable_connections.append({
            'source': {
                'service': conn['source_service']
            },
            'target': {
                'service': conn['target_service']
            }
        })

    # Prepare service lists for template
    setup_services = []
    setup_interpreters = []
    setup_integrations = []

    for service_name, service_instance in available_services.items():
        service_info = {
            "service_name": service_instance.service_name,
            "parameters": service_instance.parameters if hasattr(service_instance, 'parameters') else {},
            "parameters_helpers": service_instance.parameters_helpers if hasattr(service_instance, 'parameters_helpers') else {},
            "inputs": service_instance.inputs if hasattr(service_instance, 'inputs') else [],
            "outputs": service_instance.outputs if hasattr(service_instance, 'outputs') else []
        }
        
        if service_instance.service_type == "input":
            setup_services.append(service_info)
        elif service_instance.service_type == "interpreter":
            setup_interpreters.append(service_info)
        elif service_instance.service_type == "integration":
            setup_integrations.append(service_info)

    # Convert client_nodes to serializable format
    serializable_client_nodes = {}
    for service_name, node_data in client_nodes.items():
        serializable_client_nodes[service_name] = {
            'parameters': node_data.get('parameters', {}),
            'service_name': node_data.get('service_name', '')
        }
        
    print("setup_services:", setup_services)
    print("setup_interpreters:", setup_interpreters)
    print("setup_integrations:", setup_integrations)
    print("")
    print("client_nodes:", client_nodes)
    print("connections:", serializable_connections)
    
    cursor.close()
    db.close()
        
    # Build the service hierarchy
    service_hierarchy = build_service_hierarchy(available_services)
    print("\nService Hierarchy:")
    print(json.dumps(service_hierarchy, indent=2))
        
    return render_template('setup.html', 
                         client=client, 
                         setup_services=setup_services,
                         setup_interpreters=setup_interpreters,
                         setup_integrations=setup_integrations,
                         client_nodes=serializable_client_nodes,
                         connections=serializable_connections,
                         service_hierarchy=service_hierarchy)

def build_service_hierarchy(available_services):
    # Initialize the structure
    hierarchy = {
        "Input": []
    }
    
    # First, identify all input services (services with no inputs or only have outputs)
    input_services = {}
    interpreter_services = {}
    integration_services = {}
    
    for service_name, service_instance in available_services.items():
        service_info = {
            "service_name": service_instance.service_name,
            "inputs": getattr(service_instance, 'inputs', []),
            "outputs": getattr(service_instance, 'outputs', []),
            "service_type": service_instance.service_type
        }
        
        if service_instance.service_type == "input":
            input_services[service_name] = service_info
        elif service_instance.service_type == "interpreter":
            interpreter_services[service_name] = service_info
        elif service_instance.service_type == "integration":
            integration_services[service_name] = service_info

    # Build the hierarchy for each input service
    for input_name, input_service in input_services.items():
        input_node = {
            "name": input_name,
            "Interpreter": []
        }
        
        # Find interpreters that can accept this input service's outputs
        for interpreter_name, interpreter_service in interpreter_services.items():
            # Check if any output from input service matches any input of interpreter
            if any(output in interpreter_service['inputs'] 
                  for output in input_service['outputs']):
                
                interpreter_node = {
                    "name": interpreter_name,
                    "integration": []
                }
                
                # Find integrations that can accept this interpreter's outputs
                for integration_name, integration_service in integration_services.items():
                    if any(output in integration_service['inputs'] 
                          for output in interpreter_service['outputs']):
                        interpreter_node["integration"].append(integration_name)
                
                input_node["Interpreter"].append(interpreter_node)
        
        hierarchy["Input"].append(input_node)
    
    return hierarchy

@app.route('/check_status/<int:client_id>')
@login_required
def check_status(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT name FROM clients WHERE id = %s', (client_id,))
        client_data = cursor.fetchone()
        print(f"Client data from DB: {client_data}")
        cursor.close()
        db.close()
        
        if not client_data:
            return jsonify({'success': False, 'error': 'Client not found'})
            
        namespace = client_data['name']
        print(f"Using namespace: {namespace}")
        
        # Load kubernetes configuration
        print("Loading kubernetes config...")
        config.load_kube_config()
        print("Creating CoreV1Api client...")
        api = k8s_client.CoreV1Api()
        print("CoreV1Api client created successfully")
        
        # Get pods
        pods = []
        print("Fetching pods...")
        pod_list = api.list_namespaced_pod(namespace=namespace)
        print(f"Found {len(pod_list.items)} pods")
        for pod in pod_list.items:
            pod_info = {
                'name': pod.metadata.name,
                'status': pod.status.phase,
                'ready': all(container.ready for container in pod.status.container_statuses) if pod.status.container_statuses else False,
                'restarts': sum(container.restart_count for container in pod.status.container_statuses) if pod.status.container_statuses else 0,
                'age': pod.metadata.creation_timestamp.isoformat()
            }
            print(f"Pod info: {pod_info}")
            pods.append(pod_info)
        
        # Get services
        services = []
        print("Fetching services...")
        svc_list = api.list_namespaced_service(namespace=namespace)
        print(f"Found {len(svc_list.items)} services")
        for svc in svc_list.items:
            ports = []
            for port in svc.spec.ports:
                port_info = {
                    'name': port.name,
                    'port': port.port,
                }
                if hasattr(port, 'node_port'):
                    port_info['node_port'] = port.node_port
                ports.append(port_info)
            
            svc_info = {
                'name': svc.metadata.name,
                'type': svc.spec.type,
                'cluster_ip': svc.spec.cluster_ip,
                'ports': ports,
                'age': svc.metadata.creation_timestamp.isoformat()
            }
            print(f"Service info: {svc_info}")
            services.append(svc_info)
        
        return jsonify({
            'success': True,
            'pods': pods,
            'services': services
        })
        
    except Exception as e:
        print(f"Error checking status: {str(e)}")
        return jsonify({
            'success': False,
            'error': str(e)
        })

@app.route('/stop/<int:client_id>')
@login_required
def stop_page(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        cursor.close()
        db.close()
        
        if not client:
            flash('Client not found', 'error')
            return redirect(url_for('index'))
            
        namespace = client['name']
        print(f"Stopping services for namespace: {namespace}")
        
        return render_template('stop.html', client=client)
        
    except Exception as e:
        print(f"Error in stop_page: {str(e)}")
        flash(f'Error: {str(e)}', 'error')
        return redirect(url_for('index'))

def reset_client_ports(client_id):
    """Reset client port and api_port to 0 in the database.
    
    Args:
        client_id: The ID of the client to update
        
    Returns:
        bool: True if successful, False otherwise
    """
    try:
        db = get_db_connection()
        cursor = db.cursor()
        update_query = "UPDATE clients SET port = %s, api_port = %s, web_port = %s WHERE id = %s"
        cursor.execute(update_query, (0, 0, 0, client_id))
        db.commit()
        cursor.close()
        db.close()
        print(f"Reset port, api_port, and web_port to 0 for client {client_id}")
        return True
    except Exception as e:
        print(f"Error resetting port values: {str(e)}")
        return False

@app.route('/download_logs/<int:client_id>')
@login_required
def download_logs(client_id):
    """Generate and download logs as a zip file before stopping services."""
    from flask import send_file
    import tempfile
    import os
    import zipfile
    from datetime import datetime
    
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        cursor.close()
        db.close()

        if not client:
            return jsonify({'error': 'Client not found'}), 404

        namespace = client['name']
        timestamp = datetime.now().strftime("%Y-%m-%d:%H:%M")
        zip_filename = f"{namespace}-{timestamp}.zip"
        
        print(f"Generating logs for namespace: {namespace}")
        
        # Create temporary directory
        temp_dir = tempfile.mkdtemp()
        logs_dir = os.path.join(temp_dir, f"{namespace}-{timestamp}")
        os.makedirs(logs_dir)
        
        try:
            # Load kubernetes configuration
            config.load_kube_config()
            
            # Get all pods in the namespace
            print(f"Getting pods from namespace: {namespace}")
            result = subprocess.run(
                ['kubectl', 'get', 'pods', '-n', namespace, '-o', 'name'],
                capture_output=True,
                text=True
            )
            
            pods_found = False
            
            if result.returncode == 0 and result.stdout.strip():
                pod_names = result.stdout.strip().split('\n')
                
                # Download logs for each pod
                for pod_full_name in pod_names:
                    if pod_full_name.startswith('pod/'):
                        pod_name = pod_full_name[4:]  # Remove 'pod/' prefix
                        log_filename = f"{pod_name}.log"
                        log_filepath = os.path.join(logs_dir, log_filename)
                        
                        print(f"Downloading logs for pod: {pod_name}")
                        pods_found = True
                        
                        try:
                            log_result = subprocess.run(
                                ['kubectl', 'logs', '-n', namespace, pod_name, '--tail=1000'],
                                capture_output=True,
                                text=True,
                                timeout=30  # Add timeout to prevent hanging
                            )
                            
                            with open(log_filepath, 'w') as log_file:
                                if log_result.returncode == 0:
                                    log_content = log_result.stdout if log_result.stdout else "No logs available for this pod"
                                    log_file.write(log_content)
                                else:
                                    log_file.write(f"Error retrieving logs: {log_result.stderr}")
                                    
                        except subprocess.TimeoutExpired:
                            print(f"Timeout getting logs for pod {pod_name}")
                            with open(log_filepath, 'w') as log_file:
                                log_file.write(f"Timeout occurred while retrieving logs for pod {pod_name}")
                        except subprocess.CalledProcessError as e:
                            print(f"Error getting logs for pod {pod_name}: {e}")
                            with open(log_filepath, 'w') as log_file:
                                log_file.write(f"Error retrieving logs: {e}")
            
            if not pods_found:
                # No pods found or namespace doesn't exist, create an info file
                info_filepath = os.path.join(logs_dir, "namespace_info.txt")
                with open(info_filepath, 'w') as info_file:
                    if result.returncode != 0:
                        info_file.write(f"Namespace '{namespace}' not found or error accessing it.\n")
                        info_file.write(f"kubectl error: {result.stderr}\n")
                    else:
                        info_file.write(f"No pods found in namespace '{namespace}' at {timestamp}\n")
                    info_file.write(f"Timestamp: {timestamp}\n")
                    info_file.write(f"Client ID: {client_id}\n")
            
            # Always create a summary file
            summary_filepath = os.path.join(logs_dir, "summary.txt")
            with open(summary_filepath, 'w') as summary_file:
                summary_file.write(f"Log Archive Summary\n")
                summary_file.write(f"==================\n")
                summary_file.write(f"Client: {namespace}\n")
                summary_file.write(f"Client ID: {client_id}\n")
                summary_file.write(f"Generated: {timestamp}\n")
                summary_file.write(f"Pods found: {'Yes' if pods_found else 'No'}\n")
                if pods_found:
                    summary_file.write(f"Number of pods: {len([f for f in os.listdir(logs_dir) if f.endswith('.log')])}\n")
            
            # Create zip file
            zip_filepath = os.path.join(temp_dir, zip_filename)
            with zipfile.ZipFile(zip_filepath, 'w', zipfile.ZIP_DEFLATED) as zipf:
                for root, dirs, files in os.walk(logs_dir):
                    for file in files:
                        file_path = os.path.join(root, file)
                        arc_name = os.path.relpath(file_path, temp_dir)
                        zipf.write(file_path, arc_name)
            
            print(f"Created zip file: {zip_filepath}")
            
            # Return the zip file for download
            return send_file(
                zip_filepath,
                as_attachment=True,
                download_name=zip_filename,
                mimetype='application/zip'
            )
            
        except Exception as e:
            print(f"Error creating log archive: {e}")
            # Create an error zip with error information
            error_filepath = os.path.join(logs_dir, "error.txt")
            with open(error_filepath, 'w') as error_file:
                error_file.write(f"Error occurred while generating logs:\n")
                error_file.write(f"Error: {str(e)}\n")
                error_file.write(f"Timestamp: {timestamp}\n")
                error_file.write(f"Client: {namespace}\n")
            
            zip_filepath = os.path.join(temp_dir, zip_filename)
            with zipfile.ZipFile(zip_filepath, 'w', zipfile.ZIP_DEFLATED) as zipf:
                zipf.write(error_filepath, f"{namespace}-{timestamp}/error.txt")
            
            return send_file(
                zip_filepath,
                as_attachment=True,
                download_name=zip_filename,
                mimetype='application/zip'
            )
        finally:
            # Clean up temporary files (handled by Flask after sending file)
            pass
                
    except Exception as e:
        print(f"Error in download_logs: {str(e)}")
        return jsonify({'error': str(e)}), 500

@app.route('/stop_services/<int:client_id>', methods=['GET'])
@login_required
def stop_services(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        cursor.close()
        db.close()

        if not client:
            return jsonify({'error': 'Client not found'}), 404

        namespace = client['name']
        print(f"Processing service stop for namespace: {namespace}")

        # Clean up nginx configuration first
        print("Cleaning up nginx configuration...")
        cleanup_nginx_config(namespace)

        # Stop all services first
        print("Stopping Kubernetes services...")
        try:
            # Load kubernetes configuration
            config.load_kube_config()
            v1 = k8s_client.CoreV1Api()
            
            # Delete the namespace which will delete all resources in it
            print(f"Deleting namespace: {namespace}")
            try:
                v1.delete_namespace(name=namespace)
                print(f"Successfully deleted namespace: {namespace}")
                
                # Reset port and api_port to 0 after successfully stopping services
                reset_client_ports(client_id)
                
                return jsonify({'message': 'Services stopping, ports reset'}), 200
            except k8s_client.rest.ApiException as e:
                if e.status != 404:  # Ignore if namespace doesn't exist
                    print(f"Error deleting namespace: {str(e)}")
                    return jsonify({'error': f'Error deleting namespace: {str(e)}'}), 500
                
                # Reset port and api_port to 0 even if namespace was already deleted
                reset_client_ports(client_id)
                
                return jsonify({'message': 'Namespace already deleted, ports reset'}), 200
            
        except Exception as e:
            print(f"Error stopping Kubernetes services: {str(e)}")
            return jsonify({'error': f'Error stopping services: {str(e)}'}), 500
        
    except Exception as e:
        print(f"Error in stop_services: {str(e)}")
        return jsonify({'error': str(e)}), 500

@app.route('/delete/<int:client_id>')
@login_required
def delete_page(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        cursor.close()
        db.close()
        
        if not client:
            flash('Client not found', 'error')
            return redirect(url_for('clients'))
            
        namespace = client['name']
        print(f"Delete page for namespace: {namespace}")
        
        return render_template('delete.html', client=client)
        
    except Exception as e:
        print(f"Error in delete_page: {str(e)}")
        flash(f'Error: {str(e)}', 'error')
        return redirect(url_for('clients'))

@app.route('/delete_client/<int:client_id>', methods=['POST'])
@login_required
def delete_client(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        
        if not client:
            flash('Client not found', 'error')
            return redirect(url_for('clients'))
            
        namespace = client['name']
        results = []
        print(f"Processing client deletion for namespace: {namespace}")
        
        # Clean up nginx configuration first (must be done before deleting namespace)
        print("Cleaning up nginx configuration...")
        try:
            cleanup_nginx_config(namespace)
            results.append(('Nginx configuration cleanup', 'success', 'Configuration removed successfully'))
        except Exception as e:
            print(f"Error cleaning up nginx config: {str(e)}")
            results.append(('Nginx configuration cleanup', 'error', str(e)))
        
        # Reset client ports to 0 first
        reset_result = reset_client_ports(client_id)
        results.append((
            'Port reset', 
            'success' if reset_result else 'error',
            'Client ports reset to 0' if reset_result else 'Failed to reset ports'
        ))
        
        # Delete namespace after nginx cleanup
        cmd = ['kubectl', 'delete', 'namespace', namespace]
        try:
            print(f"Executing command: {' '.join(cmd)}")
            output = subprocess.check_output(cmd, stderr=subprocess.STDOUT)
            results.append(('Kubernetes namespace deletion', 'success', output.decode()))
        except subprocess.CalledProcessError as e:
            if b'not found' not in e.output:
                print(f"Error deleting namespace: {e.output.decode()}")
                results.append(('Kubernetes namespace deletion', 'error', e.output.decode()))
        
        # Delete from database
        print("Deleting client from database...")
        try:
            cursor.execute('DELETE FROM clients WHERE id = %s', (client_id,))
            db.commit()
            results.append(('Database deletion', 'success', 'Client removed from database'))
        except Exception as e:
            print(f"Error deleting from database: {str(e)}")
            results.append(('Database deletion', 'error', str(e)))
            
        cursor.close()
        db.close()
        
        # Check results and set appropriate flash message
        has_error = any(result[1] == 'error' for result in results)
        if has_error:
            flash('Some operations failed during client deletion. Check the logs for details.', 'warning')
        else:
            flash('Client deleted successfully', 'success')
            
        return redirect(url_for('clients'))
        
    except Exception as e:
        print(f"Error in delete_client: {str(e)}")
        flash(f'Error: {str(e)}', 'error')
        return redirect(url_for('clients'))

@app.route('/deploy/<int:client_id>')
@login_required
def deploy(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        cursor.close()
        db.close()

        if not client:
            flash('Client not found', 'error')
            return redirect(url_for('index'))

        services = []
        all_manifests = {'services': {}}
        debug_info = []
        namespace = client['name'].lower().replace(' ', '')

        # Get all services for this client
        service_data = get_service_data(client_id, 'deploy')
        debug_info.append(f"Service data: {service_data}")

        # Add namespace manifest
        all_manifests['namespace'] = {
            'apiVersion': 'v1',
            'kind': 'Namespace',
            'metadata': {
                'name': namespace
            }
        }

        # Add mosquitto manifests
        mosquitto_config = dict(MOSQUITTO_TEMPLATE['configmap'])
        mosquitto_deployment = dict(MOSQUITTO_TEMPLATE['deployment'])
        mosquitto_service = dict(MOSQUITTO_TEMPLATE['service'])
        
        # Update namespace in mosquitto manifests
        mosquitto_config['metadata']['namespace'] = namespace
        mosquitto_deployment['metadata']['namespace'] = namespace
        mosquitto_service['metadata']['namespace'] = namespace
        
        all_manifests['mosquitto'] = {
            'config': mosquitto_config,
            'deployment': mosquitto_deployment,
            'service': mosquitto_service
        }

        # Process each service
        available_services = service_factory.get_all_services()
        for service_name, service_instance in available_services.items():
            debug_info.append(f"Processing service: {service_name}")
            
            if service_name in service_data:
                service_params = service_data[service_name]
                debug_info.append(f"Found data for {service_name}: {service_params}")
                
                try:
                    # Generate manifests for this service
                    service_manifests = service_instance.get_kubernetes_manifests(namespace, service_params)
                    debug_info.append(f"Generated manifests for {service_name}: {service_manifests.keys() if service_manifests else 'None'}")
                    
                    if service_manifests:
                        all_manifests['services'][service_name] = service_manifests
                        services.append(service_name)
                except Exception as e:
                    debug_info.append(f"Error generating manifests for {service_name}: {str(e)}")

        # Create yaml directory if it doesn't exist
        os.makedirs('yaml', exist_ok=True)

        # Write namespace manifest
        with open(f'yaml/{namespace}-namespace.yaml', 'w') as f:
            yaml.dump(all_manifests['namespace'], f)

        # Write mosquitto manifests
        with open(f'yaml/{namespace}-mosquitto-config.yaml', 'w') as f:
            yaml.dump(all_manifests['mosquitto']['config'], f)
        with open(f'yaml/{namespace}-mosquitto-deployment.yaml', 'w') as f:
            yaml.dump(all_manifests['mosquitto']['deployment'], f)
        with open(f'yaml/{namespace}-mosquitto-service.yaml', 'w') as f:
            yaml.dump(all_manifests['mosquitto']['service'], f)

        # Write service manifests
        for service_name, manifests in all_manifests['services'].items():
            if 'deployment' in manifests:
                with open(f'yaml/{namespace}-{service_name}-deployment.yaml', 'w') as f:
                    yaml.dump(manifests['deployment'], f)
            if 'service' in manifests:
                with open(f'yaml/{namespace}-{service_name}-service.yaml', 'w') as f:
                    yaml.dump(manifests['service'], f)

        return render_template('deploy.html',
                           client=client,
                           manifests=all_manifests,
                           services=services,
                           debug_info=debug_info)

    except Exception as e:
        print(f"Error in deploy: {str(e)}")
        traceback.print_exc()
        flash(f'Error: {str(e)}', 'error')
        return redirect(url_for('clients'))

@app.route('/deploy_services/<int:client_id>', methods=['POST'])
@login_required
def deploy_services(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        
        if not client:
            return jsonify({
                'success': False,
                'message': 'Client not found',
            }), 404

        namespace = client['name'].lower().replace(' ', '')
        print(f"Deploying services for namespace: {namespace}")

        # Check if services are already deployed
        try:
            print("=== Checking for existing pods ===")
            # Check for existing pods
            pod_cmd = ['kubectl', 'get', 'pods', '--namespace', namespace]
            print(f"Running command: {' '.join(pod_cmd)}")
            result = subprocess.run(pod_cmd, capture_output=True, text=True)
            print(f"Pod check stdout: '{result.stdout}'")
            print(f"Pod check stderr: '{result.stderr}'")
            print(f"Pod check return code: {result.returncode}")
            
            if result.returncode == 0 and "No resources found" not in result.stdout and len(result.stdout.strip()) > 0:
                print("Found existing pods, returning error response")
                return jsonify({
                    'success': False,
                    'message': 'Services are already deployed for this client. Please stop the services first if you want to redeploy.'
                }), 400

            print("=== Checking for existing services ===")
            # Check for existing services
            svc_cmd = ['kubectl', 'get', 'services', '--namespace', namespace]
            print(f"Running command: {' '.join(svc_cmd)}")
            result = subprocess.run(svc_cmd, capture_output=True, text=True)
            print(f"Service check stdout: '{result.stdout}'")
            print(f"Service check stderr: '{result.stderr}'")
            print(f"Service check return code: {result.returncode}")
            
            if result.returncode == 0 and "No resources found" not in result.stdout and len(result.stdout.strip()) > 0:
                print("Found existing services, returning error response")
                return jsonify({
                    'success': False,
                    'message': 'Services are already deployed for this client. Please stop the services first if you want to redeploy.'
                }), 400

            print("=== No existing resources found, proceeding with deployment ===")

        except subprocess.CalledProcessError as e:
            print(f"Error checking deployment status: {e}")
            print(f"Error output: {e.stderr}")
            pass

        # Load kubernetes configuration
        config.load_kube_config()
        v1 = k8s_client.CoreV1Api()
        apps_v1 = k8s_client.AppsV1Api()

        # Create namespace if it doesn't exist
        try:
            v1.create_namespace(body={"apiVersion": "v1", "kind": "Namespace", "metadata": {"name": namespace}})
            print(f"Created namespace: {namespace}")
        except k8s_client.rest.ApiException as e:
            if e.status != 409:  # Ignore if namespace already exists
                raise

        # Deploy Mosquitto first
        print("Deploying Mosquitto...")
        try:
            # Create ConfigMap
            MOSQUITTO_TEMPLATE['configmap']['metadata']['namespace'] = namespace
            v1.create_namespaced_config_map(namespace=namespace, body=MOSQUITTO_TEMPLATE['configmap'])
            
            # Create Deployment
            MOSQUITTO_TEMPLATE['deployment']['metadata']['namespace'] = namespace
            apps_v1.create_namespaced_deployment(namespace=namespace, body=MOSQUITTO_TEMPLATE['deployment'])
            
            # Create Service
            MOSQUITTO_TEMPLATE['service']['metadata']['namespace'] = namespace
            v1.create_namespaced_service(namespace=namespace, body=MOSQUITTO_TEMPLATE['service'])
            
            print("Mosquitto deployed successfully")
        except k8s_client.rest.ApiException as e:
            if e.status != 409:  # Ignore if resources already exist
                print(f"Error deploying Mosquitto: {e}")
                raise

        # Get all available services
        available_services = service_factory.get_all_services()
        print(f"\nDeploying services for client {client_id}...")
        print(f"Available services: {list(available_services.keys())}")

        # Deploy each service
        for service_name, service_instance in available_services.items():
            print(f"\nChecking service: {service_name}")
            try:
                # Get service data from database
                cursor.execute(f'SELECT * FROM {service_name} WHERE client_id = %s', (client_id,))
                service_data = cursor.fetchone()
                
                if service_data:
                    print(f"Found configuration for {service_name}:", service_data)
                    
                    # Remove id and client_id from service data
                    service_params = {k: v for k, v in service_data.items() 
                                  if k not in ['id', 'client_id', 'created_at']}
                    
                    try:
                        # Generate manifests
                        manifests = service_instance.get_kubernetes_manifests(namespace, service_params)
                        print(f"Generated manifests for {service_name}")
                        
                        # Deploy service
                        if 'deployment' in manifests:
                            try:
                                apps_v1.create_namespaced_deployment(
                                    namespace=namespace,
                                    body=manifests['deployment']
                                )
                                print(f"{service_name} deployment created")
                            except k8s_client.rest.ApiException as e:
                                if e.status != 409:  # Ignore if deployment already exists
                                    raise

                        if 'service' in manifests:
                            try:
                                v1.create_namespaced_service(
                                    namespace=namespace,
                                    body=manifests['service']
                                )
                                print(f"{service_name} service created")
                            except k8s_client.rest.ApiException as e:
                                if e.status != 409:  # Ignore if service already exists
                                    raise
                                    
                    except Exception as e:
                        print(f"Error deploying {service_name}: {str(e)}")
                        raise
                else:
                    print(f"No configuration found for {service_name}")
            except Exception as e:
                print(f"Error processing {service_name}: {str(e)}")
                raise
        # Configure Nginx after services are deployed
        print("\nConfiguring Nginx streams...")
        try:
            nginx_success, nginx_results = configure_nginx_stream(namespace,client_id)
            if not nginx_success:
                error_msg = "\n".join([r['output'] for r in nginx_results if not r['success']])
                raise Exception(f"Failed to configure Nginx streams: {error_msg}")
            print("Nginx streams configuration completed successfully")
            
            # Configure Nginx endpoint for portal service if it exists
            print("\nConfiguring Nginx endpoint for portal...")
            endpoint_success, endpoint_results = configure_nginx_endpoint(namespace, client_id)
            # We don't raise an exception if portal endpoint config fails, just log it
            if not endpoint_success:
                print("Portal endpoint configuration did not complete successfully. This may be normal if no portal service is deployed.")
                print("\n".join([r['output'] for r in endpoint_results if not r['success']]))
            else:
                print("Nginx portal endpoint configuration completed successfully")
                
        except Exception as e:
            print(f"Error configuring Nginx: {str(e)}")
            raise
        return jsonify({
            'success': True,
            'message': f'Successfully deployed services for client {client["name"]}',
            'redirect': url_for('clients')
        })

    except Exception as e:
        print(f"Error in deploy_services: {str(e)}")
        traceback.print_exc()
        return jsonify({
            'success': False,
            'message': f'Error deploying services: {str(e)}',
            'redirect': url_for('clients')
        }), 400
    finally:
        cursor.close()
        db.close()

@app.route('/client/<int:client_id>')
@login_required
def client_details(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)

        # Get client info
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        if not client:
            flash('Client not found', 'error')
            return redirect(url_for('index'))

        # Get available services and their configurations
        available_services = service_factory.get_all_services()
        client_services = {}

        # Get existing service configurations
        for service_name, service_instance in available_services.items():
            try:
                cursor.execute(f'SELECT * FROM {service_name} WHERE client_id = %s', (client_id,))
                service_data = cursor.fetchone()
                if service_data:
                    # Remove internal fields and structure the data properly
                    cleaned_data = {k: v for k, v in service_data.items() 
                                  if k not in ['id', 'client_id', 'created_at']}
                    client_services[service_name] = cleaned_data
            except Exception as e:
                print(f"Error fetching {service_name} config: {str(e)}")
                continue

        # Convert service instances to dict with parameters
        available_services = {name: {"parameters": instance.parameters} 
                            for name, instance in available_services.items()}

        return render_template('client.html',
                           client=client,
                           available_services=available_services,
                           client_services=client_services)

    except Exception as e:
        print(f"Error in client_details: {str(e)}")
        flash(f'Error: {str(e)}', 'error')
        return redirect(url_for('clients'))
    finally:
        cursor.close()
        db.close()

def generate_schema_from_services():
    """Generate database schema from service classes."""
    schema = []
    
    # Get all services
    available_services = service_factory.get_all_services()
    
    # Generate tables for each service
    for service_name, service_instance in available_services.items():
        # Get parameters schema from service
        params = service_instance.parameters
        if not params:
            continue
            
        # Build CREATE TABLE statement
        columns = [
            "id INT AUTO_INCREMENT PRIMARY KEY",
            "client_id INT NOT NULL",
            *[f"{name} {dtype}" for name, dtype in params.items()],
            "created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
            "FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE"
        ]
        
        create_table = f"""
        CREATE TABLE IF NOT EXISTS {service_name} (
            {', '.join(columns)}
        )"""
        schema.append(create_table)
        
    return schema

@app.route('/create_database', methods=['POST'])
@login_required
def create_database():
    connection = None
    cursor = None
    try:
        # Connect to MySQL without specifying a database
        connection = mysql.connector.connect(
            host=DB_CONFIG['host'],
            user=DB_CONFIG['user'],
            password=DB_CONFIG['password']
        )
        cursor = connection.cursor()

        # Create the database if it doesn't exist
        cursor.execute(f"CREATE DATABASE IF NOT EXISTS {DB_NAME}")
        print(f"Database {DB_NAME} created successfully")

        # Switch to the new database
        cursor.execute(f"USE {DB_NAME}")

        # Generate and execute schema
        schema_statements = generate_schema_from_services()
        for statement in schema_statements:
            if statement.strip():
                try:
                    print(f"Executing: {statement}")
                    cursor.execute(statement)
                except mysql.connector.Error as e:
                    if e.errno == 1050:  # Table already exists
                        print(f"Table already exists, continuing...")
                        continue
                    raise
        
        connection.commit()
        print("Schema created successfully")

        flash('Database created successfully', 'success')
        return redirect(url_for('admin'))

    except Exception as e:
        print(f"Error creating database: {str(e)}")
        flash(f'Error creating database: {str(e)}', 'error')
        return redirect(url_for('index'))
    finally:
        if cursor:
            cursor.close()
        if connection:
            connection.close()

@app.route('/delete_database', methods=['POST'])
@login_required
def delete_database():
    try:
        # Connect to MySQL without specifying a database
        connection = mysql.connector.connect(
            host=DB_CONFIG['host'],
            user=DB_CONFIG['user'],
            password=DB_CONFIG['password']
        )
        cursor = connection.cursor()

        # Drop the database if it exists
        cursor.execute(f"DROP DATABASE IF EXISTS {DB_NAME}")
        print(f"Database {DB_NAME} dropped successfully")

        # Recreate the database
        cursor.execute(f"CREATE DATABASE {DB_NAME}")
        print(f"Database {DB_NAME} created successfully")

        # Switch to the new database
        cursor.execute(f"USE {DB_NAME}")

        # Read and execute schema.sql
        schema_file = os.path.join(os.path.dirname(__file__), 'schema.sql')
        with open(schema_file, 'r') as f:
            schema = f.read()
            # Split into individual statements
            statements = schema.split(';')
            for statement in statements:
                if statement.strip():
                    cursor.execute(statement)
        
        connection.commit()
        print("Schema recreated successfully")

        flash('Database reset successfully', 'success')
        return redirect(url_for('admin'))

    except Exception as e:
        print(f"Error resetting database: {str(e)}")
        flash(f'Error resetting database: {str(e)}', 'error')
        return redirect(url_for('admin'))
    finally:
        if cursor:
            cursor.close()
        if connection:
            connection.close()

def get_required_tables():
    """Get the list of required tables and their create statements from schema.sql."""
    schema_file = os.path.join(os.path.dirname(__file__), 'schema.sql')
    print(f"Reading schema from: {schema_file}")
    
    with open(schema_file, 'r') as f:
        schema = f.read()
    print(f"Schema content length: {len(schema)}")
    
    # Extract CREATE TABLE statements and table names
    tables = {}
    statements = schema.split(';')
    print(f"Found {len(statements)} statements")
    
    for statement in statements:
        statement = statement.strip()
        print(f"\nProcessing statement: {statement[:100]}...")  # Print first 100 chars
        
        if statement.upper().startswith('CREATE TABLE'):
            print("Found CREATE TABLE statement")
            # Extract table name - handle both with and without backticks
            try:
                if '`' in statement:
                    # Extract name between backticks
                    table_name = statement.split('`')[1]
                else:
                    # Extract name between 'EXISTS' and '('
                    table_name = statement.split('EXISTS')[1].split('(')[0].strip()
                print(f"Found table in schema: {table_name}")
                tables[table_name] = statement + ';'  # Add back the semicolon
            except Exception as e:
                print(f"Error extracting table name from statement: {statement}")
                print(f"Error: {str(e)}")
                continue
    
    print(f"Final tables found: {list(tables.keys())}")
    return tables

@app.route('/check_database', methods=['POST'])
@login_required
def check_database():
    try:
        # Try to connect to the database
        db = get_db_connection()
        cursor = db.cursor()

        # Check if all required tables exist
        missing_tables = check_and_create_missing_tables()

        if missing_tables:
            flash(f'Created missing tables: {", ".join(missing_tables)}', 'success')
        else:
            flash('All required tables exist', 'success')

        return redirect(url_for('admin'))

    except mysql.connector.Error as e:
        if e.errno == 1049:  # Database doesn't exist
            flash('Database does not exist. Please create it first.', 'error')
        else:
            flash(f'Database error: {str(e)}', 'error')
        return redirect(url_for('admin'))
    except Exception as e:
        flash(f'Error checking database: {str(e)}', 'error')
        return redirect(url_for('admin'))

@app.route('/add_missing_tables', methods=['POST'])
@login_required
def add_missing_tables():
    try:
        # Try to connect to the database
        db = get_db_connection()
        cursor = db.cursor()

        # Get existing tables
        cursor.execute("SHOW TABLES")
        existing_tables = {table[0] for table in cursor.fetchall()}
        print(f"Existing tables: {existing_tables}")

        # Get tables from schema.sql
        schema_tables = get_required_tables()
        print(f"Schema tables found: {list(schema_tables.keys())}")
        
        # Generate schema statements from services
        service_schema = generate_schema_from_services()
        
        # Combine both sources
        missing_tables = []
        
        # Create tables from schema.sql
        for table_name, statement in schema_tables.items():
            if table_name not in existing_tables:
                print(f"Creating missing table from schema.sql: {table_name}")
                print(f"Using statement: {statement}")
                try:
                    cursor.execute(statement)
                    missing_tables.append(table_name)
                except Exception as e:
                    print(f"Error creating table {table_name}: {str(e)}")
                    print(f"Statement that failed: {statement}")
                    raise
        
        # Create tables from service schema
        for statement in service_schema:
            if statement.strip():
                # Extract table name from CREATE TABLE statement
                try:
                    table_name = statement.split('CREATE TABLE IF NOT EXISTS')[1].split('(')[0].strip()
                    if table_name not in existing_tables:
                        print(f"Creating missing table from service schema: {table_name}")
                        print(f"Using statement: {statement}")
                        cursor.execute(statement)
                        missing_tables.append(table_name)
                except Exception as e:
                    print(f"Error with service table: {str(e)}")
                    print(f"Statement that failed: {statement}")
                    raise
        
        # Ensure service_connections table exists
        if 'service_connections' not in existing_tables:
            create_service_connections_table()
            missing_tables.append('service_connections')
        
        db.commit()

        if missing_tables:
            message = f'Created missing tables: {", ".join(missing_tables)}'
            print(message)
            flash(message, 'success')
        else:
            message = 'All required tables exist'
            print(message)
            flash(message, 'success')

        return redirect(url_for('admin'))

    except mysql.connector.Error as e:
        error_message = f'Database error: {str(e)}'
        print(error_message)
        flash(error_message, 'error')
        return redirect(url_for('index'))
    except Exception as e:
        error_message = f'Error checking/creating tables: {str(e)}'
        print(error_message)
        flash(error_message, 'error')
        return redirect(url_for('index'))
    finally:
        if cursor:
            cursor.close()
        if db:
            db.close()

@app.route('/add_client', methods=['POST'])
@login_required
def add_client():
    try:
        # Get client details from form
        client_name = request.form.get('client_name')
        port = request.form.get('port')
        api_port = request.form.get('api_port')
        web_port = request.form.get('web_port')
        
        if not client_name:
            flash('Client name is required', 'error')
            return redirect(url_for('index'))

        # Clean client name (remove spaces and convert to lowercase)
        client_name = client_name.strip().lower()

        # Connect to database
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)

        # Check if client already exists
        cursor.execute('SELECT * FROM clients WHERE name = %s', (client_name,))
        if cursor.fetchone():
            flash(f'Client {client_name} already exists', 'error')
            return redirect(url_for('index'))
            
        # Set default values if not provided
        if port is None:
            port = 0
        if api_port is None:
            api_port = 0
        if web_port is None:
            web_port = 0
            
        # Add new client with port, API port and web port
        cursor.execute('INSERT INTO clients (name, port, api_port, web_port) VALUES (%s, %s, %s, %s)', 
                      (client_name, port, api_port, web_port))
        client_id = cursor.lastrowid
        db.commit()

        print(f"Added new client: {client_name} (ID: {client_id})")
        flash(f'Client {client_name} added successfully', 'success')
        
        return redirect(url_for('clients'))

    except mysql.connector.Error as e:
        error_message = f'Database error: {str(e)}'
        print(error_message)
        flash(error_message, 'error')
        return redirect(url_for('clients'))
    except Exception as e:
        error_message = f'Error adding client: {str(e)}'
        print(error_message)
        flash(error_message, 'error')
        return redirect(url_for('clients'))
    finally:
        if cursor:
            cursor.close()
        if db:
            db.close()

@app.route('/create_automation', methods=['POST'])
@login_required
def create_automation():
    db = None
    cursor = None
    try:
        data = request.get_json()
        if not data:
            return jsonify({'message': 'No data provided'}), 400

        client_id = data.get('client_id')
        print("Client ID: ", client_id)
        services = data.get('services', [])
        print("Services: ", services)
        removed_services = data.get('removed_services', [])
        connections = data.get('connections', [])
        print("Removed Services: ", removed_services)
        print("Connections:", connections)
        
        if not client_id:
            return jsonify({'message': 'Client ID is required'}), 400

        db = get_db_connection()
        cursor = db.cursor(dictionary=True, buffered=True)

        # Verify client exists
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        if not client:
            return jsonify({'message': 'Client not found'}), 404

        # Remove services that were deleted
        for service_name in removed_services:
            try:
                cursor.execute(f'DELETE FROM {service_name} WHERE client_id = %s', (client_id,))
                print(f"Removed service {service_name} for client {client_id}")
            except Exception as e:
                print(f"Error removing service {service_name}: {str(e)}")

        # Store service IDs for connections
        service_ids = {}

        # Update or insert new services

        for service in services:
            service_name = service.get('name')
            print(f"Processing service: {service_name}")
            if not service_name:
                continue

            service_params = {k: v for k, v in service.items() if k != 'name'}
            print(f"Processing service {service_name} with params:", service_params)
            
            try:
                # Check if service exists
                cursor.execute(f'SELECT id FROM {service_name} WHERE client_id = %s', (client_id,))
                existing = cursor.fetchone()

                if existing:
                    # Update existing service
                    set_clause = ', '.join(f'{k} = %s' for k in service_params.keys())
                    values = list(service_params.values()) + [client_id]
                    update_query = f'UPDATE {service_name} SET {set_clause} WHERE client_id = %s'
                    cursor.execute(update_query, values)
                    service_ids[service_name] = existing['id']
                    print(f"Updated existing service {service_name}")
                else:
                    # Insert new service
                    columns = ['client_id'] + list(service_params.keys())
                    placeholders = ', '.join(['%s'] * (len(service_params) + 1))
                    values = [client_id] + list(service_params.values())
                    insert_query = f'INSERT INTO {service_name} ({", ".join(columns)}) VALUES ({placeholders})'
                    cursor.execute(insert_query, values)
                    service_ids[service_name] = cursor.lastrowid
                    print(f"Inserted new service {service_name}")

                print(f"Updated service {service_name} for client {client_id}")
            except Exception as e:
                print(f"Error updating service {service_name}: {str(e)}")
                if db:
                    db.rollback()
                return jsonify({'message': f'Error updating service {service_name}: {str(e)}'}), 500

        # Save new connections
        if connections:
            connection_values = []
            for conn in connections:
                source = conn.get('source', {})
                target = conn.get('target', {})
                source_service = source.get('service')
                target_service = target.get('service')
                print(f"Processing connection from {source_service} to {target_service}")
                
                if source_service and target_service:
                    # Get the actual service IDs from the database
                    try:
                        # Get source service ID
                        cursor.execute(f'SELECT id FROM {source_service} WHERE client_id = %s', (client_id,))
                        source_result = cursor.fetchone()
                        source_id = source_result['id'] if source_result else None
                        
                        # Get target service ID
                        cursor.execute(f'SELECT id FROM {target_service} WHERE client_id = %s', (client_id,))
                        target_result = cursor.fetchone()
                        target_id = target_result['id'] if target_result else None
                        
                        if source_id is not None and target_id is not None:
                            connection_values.append((
                                client_id,
                                source_service,
                                source_id,
                                target_service,
                                target_id
                            ))
                            print(f"Added connection: {source_service}({source_id}) -> {target_service}({target_id})")
                    except Exception as e:
                        print(f"Error getting service IDs: {str(e)}")
                        continue

            if connection_values:
                cursor.executemany(
                    'INSERT INTO service_connections (client_id, source_service, source_id, target_service, target_id) VALUES (%s, %s, %s, %s, %s)',
                    connection_values
                )
                print(f"Saved {len(connection_values)} connections")
        db.commit()
        return jsonify({'message': 'Services and connections updated successfully'})

    except Exception as e:
        if db:
            db.rollback()
        print(f"Error in create_automation: {str(e)}")
        return jsonify({'message': f'Error: {str(e)}'}), 500
    finally:
        if cursor:
            cursor.close()
        if db:
            db.close()

def get_service_nodeport(service_name: str, namespace: str) -> int:
    """
    Get the NodePort of a Kubernetes service by executing kubectl describe service command.
    
    Args:
        service_name (str): Name of the service
        namespace (str): Kubernetes namespace
        
    Returns:
        int: The NodePort number, or None if not found
    """
    try:
        # Execute kubectl describe service command
        cmd = ['kubectl', 'describe', 'service', service_name, '-n', namespace]
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            print(f"Error getting service description: {result.stderr}")
            return None
            
        # Parse the output to find NodePort
        for line in result.stdout.split('\n'):
            if 'NodePort:' in line and 'TCP' in line:
                # Extract the numeric value from something like "NodePort: tcp  30099/TCP"
                parts = line.split()
                for part in parts:
                    if '/' in part:
                        nodeport = part.split('/')[0]
                        return int(nodeport)
        
        print(f"NodePort not found in service description")
        return None
        
    except Exception as e:
        print(f"Error getting NodePort: {str(e)}")
        return None

def get_nodeport_and_tcp(service_name: str, namespace: str, port_type) -> int:
    """
    Get the TCP-api NodePort (typically port 8080) of a Kubernetes service.
    
    Args:
        service_name (str): Name of the service
        namespace (str): Kubernetes namespace
        
    Returns:
        int: The TCP-api NodePort number, or None if not found
    """

    try:
        cmd = ['kubectl', 'describe', 'service', service_name, '-n', namespace]
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            print(f"Error getting service description: {result.stderr}")
            return None

        nodeport = None    
        port = None
        print(f"port type:{port_type}")
        lines = result.stdout.split('\n')
        print(lines)
        for i, line in enumerate(lines):
            # Check if line starts with Port: and contains the exact port_type
            if line.strip().startswith('Port:'):
                parts = line.strip().split()
                if parts[1] == port_type:  # Check exact match of port_type
                    port = parts[2].split('/')[0]
            
            # Check if line starts with NodePort: and contains the exact port_type
            if line.strip().startswith('NodePort:'):
                parts = line.strip().split()
                if parts[1] == port_type:  # Check exact match of port_type
                    nodeport = parts[2].split('/')[0]

        if port and nodeport:
            return int(port), int(nodeport)
        else:
            # If we didn't find a tcp-api port, return None
            return None
        
    except Exception as e:
        print(f"Error getting TCP-api NodePort: {str(e)}")
        return None
    
@app.route('/get_pod_names/<int:client_id>')
@login_required
def get_pod_names(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT name FROM clients WHERE id = %s', (client_id,))
        client_data = cursor.fetchone()
        cursor.close()
        db.close()
        
        if not client_data:
            return jsonify({'success': False, 'error': 'Client not found'})
            
        namespace = client_data['name']
        
        # Get pod names using kubectl
        cmd = ['kubectl', 'get', 'pods', '-n', namespace, '-o', 'custom-columns=APP:metadata.labels.app', '--no-headers']
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            return jsonify({'success': False, 'error': f'Failed to get pod names: {result.stderr}'})
            
        pod_names = [name.strip() for name in result.stdout.split('\n') if name.strip()]
        return jsonify({'success': True, 'pod_names': pod_names})
        
    except Exception as e:
        return jsonify({'success': False, 'error': str(e)})

def create_service_connections_table():
    """Create the service_connections table if it doesn't exist."""
    db = get_db_connection()
    cursor = db.cursor()
    
    try:
        # Check if table exists
        cursor.execute("SHOW TABLES LIKE 'service_connections'")
        if not cursor.fetchone():
            print("Creating service_connections table...")
            
            # Create the table
            create_table_sql = """
            CREATE TABLE IF NOT EXISTS `service_connections` (
                `id` INT AUTO_INCREMENT PRIMARY KEY,
                `client_id` INT NOT NULL,
                `source_service` VARCHAR(255) NOT NULL,
                `source_id` INT NOT NULL,
                `target_service` VARCHAR(255) NOT NULL,
                `target_id` INT NOT NULL,
                `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (`client_id`) REFERENCES `clients`(`id`) ON DELETE CASCADE,
                INDEX `idx_client_source` (`client_id`, `source_service`, `source_id`),
                INDEX `idx_client_target` (`client_id`, `target_service`, `target_id`)
            )"""
            
            cursor.execute(create_table_sql)
            db.commit()
            print("service_connections table created successfully")
        else:
            print("service_connections table already exists")
            
    except Exception as e:
        print(f"Error creating service_connections table: {str(e)}")
        raise
    finally:
        cursor.close()
        db.close()

@app.route('/get_pod_logs/<int:client_id>/<string:app_name>')
@login_required
def get_pod_logs(client_id, app_name):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT name FROM clients WHERE id = %s', (client_id,))
        client_data = cursor.fetchone()
        cursor.close()
        db.close()
        
        if not client_data:
            return jsonify({'success': False, 'error': 'Client not found'})
            
        namespace = client_data['name']
        
        # Get logs using kubectl
        cmd = ['kubectl', 'logs', '-l', f'app={app_name}', '-n', namespace]
        #cmd = ['kubectl', 'logs', '-f', '-l', f'app={app_name}', '-n', namespace]
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode != 0:
            return jsonify({'success': False, 'error': f'Failed to get logs: {result.stderr}'})
            
        return jsonify({'success': True, 'logs': result.stdout})
        
    except Exception as e:
        return jsonify({'success': False, 'error': str(e)})

@app.route('/check_namespace_status/<int:client_id>', methods=['GET'])
@login_required
def check_namespace_status(client_id):
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients WHERE id = %s', (client_id,))
        client = cursor.fetchone()
        cursor.close()
        db.close()

        if not client:
            return jsonify({'error': 'Client not found'}), 404

        namespace = client['name']
        
        # Load kubernetes configuration
        config.load_kube_config()
        v1 = k8s_client.CoreV1Api()
        
        services_stopped = True
        pods_stopped = True
        
        try:
            # Check for services
            services = v1.list_namespaced_service(namespace)
            if services.items:
                services_stopped = False
        except k8s_client.rest.ApiException as e:
            if e.status != 404:  # 404 means namespace doesn't exist, which is what we want
                return jsonify({'error': f'Error checking services: {str(e)}'}), 500
        
        try:
            # Check for pods
            pods = v1.list_namespaced_pod(namespace)
            if pods.items:
                pods_stopped = False
        except k8s_client.rest.ApiException as e:
            if e.status != 404:  # 404 means namespace doesn't exist, which is what we want
                return jsonify({'error': f'Error checking pods: {str(e)}'}), 500
        
        return jsonify({
            'services_stopped': services_stopped,
            'pods_stopped': pods_stopped
        }), 200
        
    except Exception as e:
        print(f"Error checking namespace status: {str(e)}")
        return jsonify({'error': str(e)}), 500

@app.route('/trackers')
@login_required
def trackers():
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients')
        clients = cursor.fetchall()
        cursor.close()
        db.close()

        # Get all available services
        try:
            db = get_db_connection()
            cursor = db.cursor(dictionary=True)
            cursor.execute('SELECT * FROM clients')
            clients = cursor.fetchall()
            cursor.close()
            db.close()
            trackers = []
            all_trackers = []
            for client in clients:
                cmd = ['kubectl', 'get', 'services', '--namespace', client["name"]]
                result = subprocess.run(cmd, capture_output=True, text=True)
                print("looking for services in namespace", client["name"])
                if "No resources found" in result.stderr:
                    print(f"No services found in namespace {client['name']}, skipping...")
                    continue
                else:
                    print("result.stdout:", result.stdout)
                    print("result.stderr:", result.stderr)
                    nodeport, port = get_nodeport_and_tcp('listener-service', f'{client["name"]}','tcp-api')
                    print(f"node:{nodeport}, port:{port}")
                    if nodeport:  # Only process if we found a tcp-api nodeport
                        trackerlist_url = "http://"+MINIKUBE_IP+":"+str(nodeport)+"/api/v1/trackerlist"
                        print(f"Fetching from URL: {trackerlist_url}")
                        try:
                            response = requests.get(trackerlist_url, timeout=5)  # Add timeout
                            print(f"Response status for {client['name']}: {response.status_code}")
                            if response.ok:
                                client_trackers = response.json()
                                print(f"Got trackers for {client['name']}: {client_trackers}")
                                for tracker in client_trackers:
                                    all_trackers.append({
                                        'client': client['name'],
                                        'imei': tracker['imei'],
                                        'protocol': tracker['protocol'],
                                        'port': port
                                    })
                            else:
                                print(f"Bad response for {client['name']}: {response.text}")
                        except requests.exceptions.Timeout:
                            print(f"Timeout fetching trackers for {client['name']}")
                        except requests.exceptions.ConnectionError:
                            print(f"Connection error fetching trackers for {client['name']}")
                        except Exception as e:
                            print(f"Error fetching trackers for {client['name']}: {str(e)}")
                print("\nAll Trackers:",all_trackers)

            return render_template('trackers.html', trackers=all_trackers)
        except Exception as e:
            print(f"Error getting nodeports: {str(e)}")
            return render_template('trackers.html', trackers=[])

    except Exception as e:
        flash(f'Error: {str(e)}', 'error')
        return render_template('index.html', clients=[], services={})

@app.route('/api/sendcommand', methods=['POST'])
@login_required
def send_command():
    try:
        data = request.get_json()
        client = data.get('client')
        imei = data.get('imei')
        payload = data.get('payload')
        
        if not all([client, imei, payload]):
            return jsonify({'error': 'Missing required fields'}), 400

        # Get nodeport for the client
        port, nodeport = get_nodeport_and_tcp('listener-service', client,'tcp-api')
        if not nodeport:
            return jsonify({'error': f'Could not find nodeport for client {client}'}), 404

        # Send command to the actual endpoint
        url = f"http://{MINIKUBE_IP}:{nodeport}/api/v1/sendcommand"
        command_data = {
            "imei": imei,
            "data": payload  # Use payload directly, it's already hex
        }
        
        response = requests.post(
            url,
            headers={"Content-Type": "application/json"},
            json=command_data
        )
        
        # Return the response from the service
        return jsonify(response.json()), response.status_code
        
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/api/trackers/data', methods=['POST'])
@login_required
def get_trackers_data():
    try:
        # Get DataTables parameters
        draw = request.json.get('draw', 1)
        start = request.json.get('start', 0)
        length = request.json.get('length', 10)
        search = request.json.get('search', {}).get('value', '')
        order_column = request.json.get('order', [{}])[0].get('column', 0)
        order_dir = request.json.get('order', [{}])[0].get('dir', 'asc')

        # Get all trackers first
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        cursor.execute('SELECT * FROM clients')
        clients = cursor.fetchall()
        cursor.close()
        db.close()

        all_trackers = []
        for client in clients:
            has_services, nodeport, port = get_namespace_services(client["name"])
            if not has_services:
                continue
                
            if nodeport:
                client_trackers = fetch_client_trackers(client["name"], nodeport, port)
                all_trackers.extend(client_trackers)
                
        print("\nAll Trackers:", all_trackers)
        # Filter based on search term
        if search:
            search = search.lower()
            filtered_data = [
                t for t in all_trackers
                if search in t['client'].lower() 
                or search in t['imei'].lower() 
                or search in t['protocol'].lower()
                or (t['port'] is not None and search in str(t['port']))
            ]
        else:
            filtered_data = all_trackers

        # Sort data
        column_names = ['client', 'imei', 'protocol']
        if 0 <= order_column < len(column_names):
            key = column_names[order_column]
            filtered_data.sort(
                key=lambda x: x[key],
                reverse=(order_dir == 'desc')
            )

        # Paginate
        total_records = len(all_trackers)
        total_filtered = len(filtered_data)
        paginated_data = filtered_data[start:start + length]

        return jsonify({
            'draw': draw,
            'recordsTotal': total_records,
            'recordsFiltered': total_filtered,
            'data': paginated_data
        })

    except Exception as e:
        print(f"Error in get_trackers_data: {str(e)}")
        return jsonify({
            'draw': 1,
            'recordsTotal': 0,
            'recordsFiltered': 0,
            'data': [],
            'error': str(e)
        }), 500

def get_namespace_services(namespace: str):
    """
    Get services for a given namespace.
    Returns tuple of (success, nodeport, port) where success is a boolean indicating if services were found
    """
    try:
        cmd = ['kubectl', 'get', 'services', '--namespace', namespace]
        result = subprocess.run(cmd, capture_output=True, text=True)
        print("looking for services in namespace", namespace)
        
        if result.returncode != 0 or "No resources found" in result.stderr:
            print(f"No services found in namespace {namespace}, skipping...")
            return False, None, None
        
        # Get nodeport and port information
        result = get_nodeport_and_tcp('listener-service', namespace, 'tcp-api')
        if result is None:
            print(f"No nodeport/port found for {namespace}")
            return False, None, None
            
        port, nodeport = result
        print(f"nodeport:{nodeport} port:{port}")
        return True, nodeport, port
        
    except Exception as e:
        print(f"Error getting namespace services: {str(e)}")
        return False, None, None

def fetch_client_trackers(client_name: str, nodeport: int, port: int):
    """
    Fetch trackers for a specific client.
    Returns list of tracker dictionaries.
    """
    try:
        url = f"http://{MINIKUBE_IP}:{nodeport}/api/v1/trackerlist"
        response = requests.get(url, timeout=5)
        if response.ok:
            client_trackers = response.json()
            return [{
                'client': client_name,
                'imei': tracker['imei'],
                'protocol': tracker['protocol'],
                'port': port
            } for tracker in client_trackers]
    except Exception as e:
        print(f"Error fetching trackers for {client_name}: {str(e)}")
    return []

@app.route('/admin')
@login_required
def admin():
    return render_template('admin.html')

@app.route('/api-docs')
@login_required
def api_docs():
    """Redirect to Swagger API documentation."""
    return redirect('/apidocs/')

@app.route('/admin/monitor', methods=['GET', 'POST'], endpoint='admin_monitor')
@login_required
def admin_monitor():
    import mysql.connector
    db = mysql.connector.connect(
        host=DB_CONFIG['host'],
        user=DB_CONFIG['user'],
        password=DB_CONFIG['password'],
        database=DB_NAME
    )
    cursor = db.cursor(dictionary=True)
    message = None

    # Add phone functionality
    if request.method == 'POST':
        if 'add_phone' in request.form:
            phone_number = request.form.get('new_phone_number')
            contact_name = request.form.get('new_contact_name')
            enabled = 1 if request.form.get('new_phone_enabled') == 'on' else 0
            cursor.execute(
                "INSERT INTO whatsapp_phones (phone_number, contact_name, enabled) VALUES (%s, %s, %s)",
                (phone_number, contact_name, enabled)
            )
            db.commit()
            message = 'Phone number added successfully.'
        elif 'delete_phone' in request.form:
            phone_id = request.form.get('delete_phone_id')
            cursor.execute("DELETE FROM whatsapp_phones WHERE id=%s", (phone_id,))
            db.commit()
            message = 'Phone number deleted successfully.'
        else:
            # Update monitor config
            monitor_fields = [
                'check_interval_minutes', 'data_threshold_minutes', 'restart_enabled',
                'alert_enabled', 'alert_webhook_url', 'request_timeout_seconds', 'retry_attempts',
                'retry_delay_seconds'
            ]
            cursor.execute("SELECT * FROM monitor_config LIMIT 1")
            monitor_config = cursor.fetchone()
            if monitor_config:
                for field in monitor_fields:
                    value = request.form.get(field)
                    if value is not None:
                        cursor.execute(f"UPDATE monitor_config SET {field}=%s WHERE id=%s", (value, monitor_config['id']))

            # WhatsApp config (global)
            whatsapp_fields = [
                'enabled', 'api_key', 'channel_id', 'namespace', 'template_name', 'language_code',
                'alert_interval_minutes', 'min_failure_count'
            ]
            for field in whatsapp_fields:
                value = request.form.get(field)
                if value is not None:
                    cursor.execute(f"UPDATE whatsapp_config SET {field}=%s WHERE id=1", (value,))

            # WhatsApp phones (global)
            phone_id = request.form.get('phone_id')
            phone_number = request.form.get('phone_number')
            contact_name = request.form.get('contact_name')
            enabled = request.form.get('phone_enabled')
            if phone_id and phone_number:
                cursor.execute(
                    "UPDATE whatsapp_phones SET phone_number=%s, contact_name=%s, enabled=%s WHERE id=%s",
                    (phone_number, contact_name, enabled, phone_id)
                )
            db.commit()
            message = 'Monitor parameters updated successfully.'

    # Fetch current values
    cursor.execute("SELECT * FROM monitor_config LIMIT 1")
    monitor_config = cursor.fetchone()
    cursor.execute("SELECT * FROM whatsapp_config LIMIT 1")
    whatsapp_config = cursor.fetchone()
    cursor.execute("SELECT * FROM whatsapp_phones")
    whatsapp_phones = cursor.fetchall()
    db.close()

    return render_template(
        'admin_monitor.html',
        monitor_config=monitor_config,
        whatsapp_config=whatsapp_config,
        whatsapp_phones=whatsapp_phones,
        message=message
    )

# Database configuration
DB_CONFIG = {
    'host': 'localhost',
    'user': 'gpscontrol',
    'password': 'qazwsxedc'
}

DB_NAME = 'automation_app'

# Initialize ServiceFactory
service_factory = ServiceFactory()
available_services = service_factory.get_all_services()
print("\nAvailable services in factory:", list(available_services.keys()))
for name, service in available_services.items():
    print(f"Service: {name}, Class: {service.__class__.__name__}")

# GT06 Protocol Handling Functions
def get_next_gt06_serial(imei):
    """
    Get the next serial number for a GT06 device by IMEI.
    Increments the counter each time and handles wraparound at 0xFFFF.
    
    Args:
        imei (str): The IMEI of the device
        
    Returns:
        str: The next serial number as a 4-character hex string (e.g., "05F2")
    """
    try:
        db = get_db_connection()
        cursor = db.cursor(dictionary=True)
        
        # Try to get the current serial number
        cursor.execute('SELECT serial_number FROM gt06_serials WHERE imei = %s', (imei,))
        result = cursor.fetchone()
        
        if result:
            # Increment the serial number
            current_serial = result['serial_number']
            new_serial = (current_serial + 1) % 65536  # Wrap around at 0xFFFF
            
            # Update the database
            cursor.execute(
                'UPDATE gt06_serials SET serial_number = %s WHERE imei = %s',
                (new_serial, imei)
            )
        else:
            # First time for this IMEI, start at 0
            new_serial = 0
            cursor.execute(
                'INSERT INTO gt06_serials (imei, serial_number) VALUES (%s, %s)',
                (imei, new_serial)
            )
        
        db.commit()
        
        # Convert to 4-character hex string (e.g., "05F2")
        return f"{new_serial:04X}"
        
    except Exception as e:
        print(f"Error getting next GT06 serial: {str(e)}")
        # Return a default value in case of error
        return "0000"
    finally:
        if cursor:
            cursor.close()
        if db:
            db.close()

def calculate_crc16(data):
    """
    Calculate CRC-16 (CRC-ITU) for the given data.
    Uses the lookup table from the GT06 protocol documentation appendix.
    
    Args:
        data (bytes): The data to calculate CRC for
        
    Returns:
        bytes: The 2-byte CRC in big-endian format
    """
    # CRC-ITU lookup table from the documentation
    crctab16 = [
        0x0000, 0x1189, 0x2312, 0x329B, 0x4624, 0x57AD, 0x6536, 0x74BF,
        0x8C48, 0x9DC1, 0xAF5A, 0xBED3, 0xCA6C, 0xDBE5, 0xE97E, 0xF8F7,
        0x1081, 0x0108, 0x3393, 0x221A, 0x56A5, 0x472C, 0x75B7, 0x643E,
        0x9CC9, 0x8D40, 0xBFDB, 0xAE52, 0xDAED, 0xCB64, 0xF9FF, 0xE876,
        0x2102, 0x308B, 0x0210, 0x1399, 0x6726, 0x76AF, 0x4434, 0x55BD,
        0xAD4A, 0xBCC3, 0x8E58, 0x9FD1, 0xEB6E, 0xFAE7, 0xC87C, 0xD9F5,
        0x3183, 0x200A, 0x1291, 0x0318, 0x77A7, 0x662E, 0x54B5, 0x453C,
        0xBDCB, 0xAC42, 0x9ED9, 0x8F50, 0xFBEF, 0xEA66, 0xD8FD, 0xC974,
        0x4204, 0x538D, 0x6116, 0x709F, 0x0420, 0x15A9, 0x2732, 0x36BB,
        0xCE4C, 0xDFC5, 0xED5E, 0xFCD7, 0x8868, 0x99E1, 0xAB7A, 0xBAF3,
        0x5285, 0x430C, 0x7197, 0x601E, 0x14A1, 0x0528, 0x37B3, 0x263A,
        0xDECD, 0xCF44, 0xFDDF, 0xEC56, 0x90A9, 0x8120, 0xB3BB, 0xA232,
        0x5AC5, 0x4B4C, 0x79D7, 0x685E, 0x1CE1, 0x0D68, 0x3FF3, 0x2E7A,
        0xE70E, 0xF687, 0xC41C, 0xD595, 0xA12A, 0xB0A3, 0x8238, 0x93B1,
        0x6B46, 0x7ACF, 0x4854, 0x59DD, 0x2D62, 0x3CEB, 0x0E70, 0x1FF9,
        0xF78F, 0xE606, 0xD49D, 0xC514, 0xB1AB, 0xA022, 0x92B9, 0x8330,
        0x7BC7, 0x6A4E, 0x58D5, 0x495C, 0x3DE3, 0x2C6A, 0x1EF1, 0x0F78
    ]

    # Initialize with 0xFFFF
    fcs = 0xFFFF

    # Calculate CRC for each byte in the input data
    for byte in data:
        fcs = (fcs >> 8) ^ crctab16[(fcs ^ byte) & 0xFF]

    # Perform bitwise NOT and ensure it's a 16-bit value
    crc_result = ~fcs & 0xFFFF

    # Return CRC as 2 bytes, big-endian
    return crc_result.to_bytes(2, byteorder='big')

def generate_gt06_command(command_str, serial_no_hex, language_code=2):
    """
    Generates a GT06 server-to-terminal command packet (Protocol 0x80).
    Now returns both hex and ascii formats.

    Args:
        command_str: The command content string (e.g., "WHERE#").
        serial_no_hex: The 2-byte information serial number as a hex string
                       (e.g., "05F2").
        language_code: The language code (1 for Chinese, 2 for English).

    Returns:
        A dictionary with both 'hex' and 'ascii' formats of the command packet.
    """
    try:
        # --- Fixed Values ---
        start_bit = b'\x78\x78'
        protocol_no = b'\x80'
        # Example Server Flag Bit (can be customized if needed)
        server_flag = b'\x00\x00\x00\x01'
        stop_bit = b'\x0d\x0a'

        # --- Process Inputs ---
        command_content = command_str.encode('ascii')
        if language_code == 1:
            language_bytes = b'\x00\x01' # Chinese
        elif language_code == 2:
            language_bytes = b'\x00\x02' # English
        else:
            print("Error: Invalid language code. Use 1 for Chinese or 2 for English.")
            return None

        # Ensure serial number is 2 bytes (4 hex chars)
        if len(serial_no_hex) != 4:
             print(f"Error: Serial number hex '{serial_no_hex}' must be 4 characters (2 bytes).")
             return None
        serial_no_bytes = binascii.unhexlify(serial_no_hex)
        if len(serial_no_bytes) != 2:
             print(f"Error: Serial number '{serial_no_hex}' must be 2 bytes.")
             return None

        # --- Calculate Lengths ---
        # Length of Command = Server Flag Bit length (4) + Command Content length
        len_command_val = 4 + len(command_content)
        if len_command_val > 255:
            print("Error: Command content too long.")
            return None
        len_command_byte = len_command_val.to_bytes(1, byteorder='big')

        # Information Content = Length of Command byte + Server Flag + Command Content + Language
        information_content = len_command_byte + server_flag + command_content + language_bytes

        # Packet Length = len(Proto No + Info Content + Serial No + Error Check)
        core_data_len = 1 + len(information_content) + 2 # Proto No(1) + Info Content + Serial No(2)
        packet_len_val = core_data_len + 2 # Add 2 for CRC bytes
        if packet_len_val > 255:
             print("Error: Calculated Packet Length exceeds 1 byte.")
             return None
        packet_len_byte = packet_len_val.to_bytes(1, byteorder='big')

        # --- Calculate CRC ---
        # Data for CRC: Packet Length byte + Protocol No + Information Content + Serial No bytes
        data_for_crc = packet_len_byte + protocol_no + information_content + serial_no_bytes
        # Use the function based on the manual's table
        error_check = calculate_crc16(data_for_crc)

        # --- Assemble Final Packet ---
        full_packet = (
            start_bit +
            packet_len_byte +
            protocol_no +
            information_content +
            serial_no_bytes +
            error_check +
            stop_bit
        )

        # Return both hex and ASCII formats
        hex_format = binascii.hexlify(full_packet).decode('ascii')
        
        # For ASCII format, we'll keep the raw bytes to be sent directly
        ascii_format = full_packet
        
        return {
            'hex': hex_format,
            'ascii': ascii_format
        }

    except binascii.Error:
        print(f"Error: Invalid hexadecimal string for serial number: '{serial_no_hex}'")
        return None
    except Exception as e:
        print(f"An error occurred: {e}")
        return None
    
@app.route('/api/generate_gt06_packet', methods=['POST'])
def generate_gt06_packet_api():
    """API endpoint to generate GT06 command packets based on the command and IMEI."""
    try:
        data = request.get_json()
        if not data:
            return jsonify({"error": "No data provided"}), 400
        
        imei = data.get('imei')
        command = data.get('command')
        
        if not imei or not command:
            return jsonify({"error": "IMEI and command are required"}), 400
        
        # Get the next serial number for this IMEI
        serial_no_hex = get_next_gt06_serial(imei)
        
        # Default language is English (2)
        language_code = 2
        
        # Generate the packet
        packet_result = generate_gt06_command(command, serial_no_hex, language_code)
        
        if not packet_result:
            return jsonify({"error": "Failed to generate packet. Check command format."}), 400
        
        # Return both hex and ASCII formats
        return jsonify({
            "packet": packet_result['ascii'].decode('latin-1'),  # Use latin-1 to preserve binary bytes as characters
            "hex": packet_result['hex'],
            "serial": serial_no_hex,
            "message": "Packet generated successfully"
        })
        
    except Exception as e:
        print(f"Error generating GT06 packet: {str(e)}")
        return jsonify({"error": f"Internal server error: {str(e)}"}), 500

@app.route('/api/generate_bsj_packet', methods=['POST'])
def generate_bsj_packet_api():
    """API endpoint to generate BSJ command packets with auto-derived phone and serial."""
    try:
        data = request.get_json()
        if not data:
            return jsonify({"error": "No data provided"}), 400
        
        command = data.get('command')
        imei = data.get('imei')
        
        if not command or not imei:
            return jsonify({"error": "Command and IMEI are required"}), 400
        
        # Derive phone number from last 12 digits of IMEI
        if len(imei) < 12:
            # Pad with zeros if IMEI is too short
            phone = imei.zfill(12)
        else:
            # Take last 12 digits
            phone = imei[-12:]
        
        # Get next serial number using same function as GT06
        serial_hex = get_next_gt06_serial(imei)
        serial = int(serial_hex, 16)  # Convert hex to decimal
        
        print(f"Using derived phone: {phone} and serial: {serial} for IMEI: {imei}")
        
        # Generate the packet
        packet_result = generate_bsj_command(command, phone, serial)
        
        if not packet_result:
            return jsonify({"error": "Failed to generate packet. Check command format."}), 400
        
        # Return the hex format
        return jsonify({
            "hex": packet_result,
            "message": "BSJ packet generated successfully"
        })
        
    except Exception as e:
        print(f"Error generating BSJ packet: {str(e)}")
        return jsonify({"error": f"Internal server error: {str(e)}"}), 500

# svc=listener-service
if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0')