# Project Jonobridge

The project is disigned with microservices within the Kubernetes Operator Architecture
The system has a predefined microservices: listener,forwarder, xpot, meitrackprotocol, pinoprotocol, proxy, recursoconfiable ruptelaprotocol, etc. Integrating these additional microservices into the existing Kubernetes operator architecture would likely involve a similar approach to the one detailed for the services.

## Overview of the Kubernetes Operator Architecture

The sources describe a Kubernetes operator, designed using the Kubebuilder framework, for deploying and managing two interconnected applications: listener, proxy, forwarder, xpot, meitrackprotocol, pinoprotocol, proxy, recursoconfiable, ruptelaprotocol, etc. along with a Mosquitto MQTT broker. Let's examine the key architectural patterns at play:

1. Custom Resource Definitions (CRDs) and the Operator Pattern

At the heart of this architecture lies the concept of CRDs, which extend Kubernetes's functionality by enabling the definition of new resource types.

The "listener", "forwarder", etc... are applications are managed through their corresponding CRDs: listener, forwarder, etc.

The operator itself embodies the Operator pattern, acting as a dedicated controller for these custom resources.

It continuously monitors the state of Jonobridge and Forwarder CR instances, reconciling their desired state (defined in the CRD specs) with the actual state within the cluster.

2. Reconciliation Loop and Event-Driven Architecture

The operator's core functionality revolves around the reconciliation loop, an event-driven process.

Whenever a deployment,services CR is created, updated, or deleted, the operator's Reconcile function is triggered.

This function then fetches the CR's specification, compares it to the current deployment status, and takes actions to ensure alignment.

This event-driven nature enables the operator to autonomously manage the lifecycle of the applications.

3. Microservice Architecture and Containerization

The "listener", "forwarder", etc.. applications are designed as separate microservices.

Each microservice is packaged as a Docker image, allowing for independent development, deployment, and scaling.

Docker provides containerization, ensuring consistent execution environments across different platforms.

This approach promotes modularity, maintainability, and flexibility within the overall system.

4. Message Queueing with MQTT and Mosquitto

Communication between services leverages the MQTT protocol (using the Mosquitto broker).

By ex:
"listener" connects to the MQTT broker, exposes REST endpoints, and registers IMEI-based clients.

"forwarder" subscribes to MQTT topics and forwards the data received to a specified target URL.

This asynchronous messaging pattern facilitates decoupled and scalable data exchange between the components.

5. Deployment on Kubernetes with Load Balancing

The entire system is deployed on a Kubernetes cluster, which provides orchestration, scaling, and self-healing capabilities.

Deployments define how each application is deployed and scaled, while Services expose them within and outside the cluster.

Nginx is used as a load balancer to distribute traffic to the exposed services.

NodePort is utilized to expose services externally on specific ports.
The sources provide a practical guide to building and deploying this operator, but they do not explicitly use the terms "Operator pattern" or "Microservice architecture." These terms are commonly used to describe the patterns observed in this system.

# Requirements to run Jonobrige

This code provides a basic implementation of the app using Kubebuilder, a popular tool for creating Kubernetes Operators with Go. It includes steps for defining Custom Resource Definitions (CRDs), setting up the Operator, and creating a Dockerfile. Below are the hardware and software requirements needed to successfully run Jonobridge.

## Hardware requirements
- Ubuntu server (headless) version 22.04
- Hard Disk 80 GB 
- Memory 8Gb
- CPUs : 2 (at least)

## Software requirements
### Go Lang Installation:

```
cd /tmp
wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz
```

Edit bash
```
vim ~/.bashrc 
```
Put inside:
```
export PATH=$PATH:/usr/local/go/bin 
export GOPATH="/home/user/go" 
export GO111MODULE="on" 
export GOPROXY="https://proxy.golang.org,direct"
```
Logout and login again

### Docker installation
**Step 1: Install Docker on Ubuntu 22.04**

Update the apt package index and install dependencies:
```
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg lsb-release -y
```

**Add Docker’s official GPG key:**

```
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg 
```

**Set up the repository:**

```
echo \
"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
$(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

**Install Docker Engine:**
```
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
```

**Start and enable Docker:**
```
sudo systemctl enable docker
sudo systemctl start docker
```

**Verify Docker installation**

```
docker --version
```

**Execute this commands in the console**

```
sudo groupadd docker
sudo usermod -aG docker $USER
```

Logout and login again

Test

```
docker run hello-world
```

### Follow the following guide to Install Kubectl

*https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/*


### Kubebuilder installation guide: (to to the link to install)
```
cd /tmp
wget https://github.com/kubernetes-sigs/kubebuilder/releases/download/v4.3.0/kubebuilder_linux_amd64
sudo mv kubebuilder_linux_amd64 /usr/local/bin/kubebuilder
sudo chmod +x /usr/local/bin/kubebuilder
```

### Kubernetes installation cluster (Minikube set up)
***Minikube Install***

```
# First Install VirtualBox
sudo apt update
sudo apt install -y virtualbox virtualbox-ext-pack

# Download and install Minikube
wget https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo cp minikube-linux-amd64 /usr/local/bin/minikube
sudo chmod +x /usr/local/bin/minikube
```

### Start Minikube**:*
```
minikube start --memory=8192 --cpus=4 --driver=docker --disk-size=50g --extra-config=kubelet.max-pods=250
```

### MySQL installation

**1\. Update Your System**

Ensure your package list is up to date by running the following commands:
```
sudo apt update
sudo apt upgrade -y 
```

**2\. Install MySQL Server** **Install MySQL Server using the following command:**

```
sudo apt install mysql-server -y 
```

**3\. Secure the MySQL Installation**

Run the security script to configure your MySQL installation:
```
sudo mysql_secure_installation 
```

During the process, you will:  
- Set a root password (if not set during installation).  
- Remove anonymous users.  
- Disable root login remotely.  
- Remove test databases.  
- Reload privilege tables.

**4\. Start and Enable MySQL Service**

To ensure the MySQL server is running and starts automatically at boot, run:
```
sudo systemctl start mysql
sudo systemctl enable mysql
```

**5\. Verify the MySQL Installation**

Check the status of the MySQL service to ensure it's running correctly:
```
sudo systemctl status mysql
```

**6\. Set MYSQL**
```
mysql -u root -p
CREATE DATABASE automation_app ;
CREATE DATABASE bridge;
CREATE USER 'gpscontrol'@'localhost' IDENTIFIED BY 'qazwsxedc';
CREATE USER IF NOT EXISTS 'gpscontrol'@'mini' IDENTIFIED BY 'qazwsxedc';
GRANT ALL PRIVILEGES ON * . * TO 'gpscontrol'@'localhost';
GRANT ALL PRIVILEGES ON bridge.* TO 'gpscontrol'@'mini' IDENTIFIED BY 'qazwsxedc';

-- Create the user (if it doesn't exist already)
CREATE USER IF NOT EXISTS 'gpscontrol'@'mini' IDENTIFIED BY 'qazwsxedc';

-- Grant privileges to the user
GRANT ALL PRIVILEGES ON bridge.* TO 'gpscontrol'@'mini';

-- Apply the changes

FLUSH PRIVILEGES;
FLUSH PRIVILEGES;
CREATE USER 'gpscontrol'@'mini' IDENTIFIED BY 'qazwsxedc';
ALTER USER 'gpscontrol'@'mini' IDENTIFIED BY 'qazwsxedc';
GRANT ALL PRIVILEGES ON bridge.* TO 'gpscontrol'@'mini';
FLUSH PRIVILEGES;
```

### DOCKER LOGIN (perform this step to gain access to the account)
```
docker login -u maddsystems -p GPSc0ntr0l18
```
### NGINX Installation

#### How to Install Nginx with UDP Support on Ubuntu

This guide explains the steps to install Nginx with UDP support on an Ubuntu system.

---

### Step 1: Update Your System

First, update your package list and upgrade installed packages:

```bash
sudo apt update && sudo apt upgrade -y
```

---

### Step 2: Install Nginx

Install Nginx from the official repositories:

```bash
sudo apt install nginx -y
```

---

### Step 3: Verify Nginx Stream Module Support

Check if Nginx is built with the `stream` module (required for UDP support):

```bash
nginx -V 2>&1 | grep -- '--with-stream'
```

- If the output includes `--with-stream`, your Nginx supports UDP.
- If not, you will need to compile Nginx from source with `--with-stream` (see Step 6).

---

### Step 4: Configure Nginx for TCP/UDP Proxying

1. Open the Nginx configuration file for editing:
   ```bash
   sudo mkdir /etc/nginx/streams.d
   sudo nano /etc/nginx/nginx.conf
   ```

2. Add the following `stream` block to configure a basic TCP/UDP proxy before the http block:
   ```nginx
    # Stream block for TCP/UDP configurations
    stream {
        # Include all configuration files from the streams.d directory
        include /etc/nginx/streams.d/*.conf;
    }
   ```

3. Save and exit the file.

---

### Step 5: Restart Nginx

Restart Nginx to apply the changes:

```bash
sudo systemctl restart nginx
```

---

### Step 6: Compile Nginx with Stream Module 

Follow these steps to build it from source:

1. Install the required build dependencies:

   ```bash
   sudo apt install build-essential libpcre3 libpcre3-dev zlib1g zlib1g-dev libssl-dev -y
   ```

2. Download the latest Nginx source code from the [official website](https://nginx.org/en/download.html). For example:

   ```bash
   cd /tmp
   wget https://nginx.org/download/nginx-1.24.0.tar.gz
   tar -zxvf nginx-1.24.0.tar.gz
   cd nginx-1.24.0
   ```

3. Configure the build with the `--with-stream` option:

   ```bash
      ./configure --with-stream \
            --prefix=/usr/share/nginx \
            --sbin-path=/usr/sbin/nginx \
            --modules-path=/usr/lib/nginx/modules \
            --conf-path=/etc/nginx/nginx.conf \
            --error-log-path=/var/log/nginx/error.log \
            --http-log-path=/var/log/nginx/access.log \
            --pid-path=/run/nginx.pid \
            --lock-path=/var/lock/nginx.lock \
            --user=nginx \
            --group=nginx \
            --with-stream_ssl_module \
            --with-http_ssl_module \
            --with-http_v2_module
   ```

4. Build and install Nginx:

   ```bash
   make
   sudo make install
   ```

5. 
```
sudo adduser --system --no-create-home --shell /bin/false --group --disabled-login nginx
sudo mkdir -p /var/log/nginx
sudo chown -R nginx:nginx /var/log/nginx
```

6.Create a systemd service file for Nginx. Create /etc/systemd/system/nginx.service with the following content:
```
[Unit]
Description=nginx - high performance web server
Documentation=https://nginx.org/en/docs/
After=network-online.target remote-fs.target nss-lookup.target
Wants=network-online.target

[Service]
Type=forking
PIDFile=/run/nginx.pid
ExecStartPre=/usr/sbin/nginx -t -q -g 'daemon on; master_process on;'
ExecStart=/usr/sbin/nginx -g 'daemon on; master_process on;'
ExecReload=/usr/sbin/nginx -g 'daemon on; master_process on;' -s reload
ExecStop=-/sbin/start-stop-daemon --quiet --stop --retry QUIT/5 --pidfile /run/nginx.pid
TimeoutStopSec=5
KillMode=mixed
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

7. Start and enable Nginx:
```
sudo systemctl daemon-reload
sudo systemctl start nginx
sudo systemctl enable nginx
sudo mkdir -p /etc/nginx/sites-available
sudo mkdir -p /etc/nginx/sites-enabled
```

8.Verify the installation:
```
nginx -V
```

9. sudo vim /etc/nginx/nginx.conf
```
user www-data;
worker_processes auto;
pid /run/nginx.pid;

# Include dynamically loaded modules
include /etc/nginx/modules-enabled/*.conf;

events {
    worker_connections 768;
    # multi_accept on;
}

# Stream block for TCP/UDP configurations
stream {
    # Include all configuration files from the streams.d directory
    include /etc/nginx/streams.d/*.conf;
}

http {
    ##
    # Basic Settings
    ##

    sendfile on;
    tcp_nopush on;
    types_hash_max_size 2048;
    proxy_read_timeout 60s;
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    send_timeout 60s;  
    # server_tokens off;

    # server_names_hash_bucket_size 64;
    # server_name_in_redirect off;

    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    ##
    # SSL Settings
    ##

    ssl_protocols TLSv1 TLSv1.1 TLSv1.2 TLSv1.3; # Dropping SSLv3, ref: POODLE
    ssl_prefer_server_ciphers on;

    ##
    # Logging Settings
    ##

    access_log off; # Disable access logging
    error_log /dev/null crit; # Log critical errors only and discard logs

    ##
    # Gzip Settings
    ##

    gzip on;

    # gzip_vary on;
    # gzip_proxied any;
    # gzip_comp_level 6;
    # gzip_buffers 16 8k;
    # gzip_http_version 1.1;
    # gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

    ##
    # Virtual Host Configs
    ##

    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
```
10. sudo nano /etc/nginx/sites-available/jonobridge.madd.com.mx
```
server {
    listen 443 ssl;
    server_name jonobridge.madd.com.mx;

    # SSL configuration
    ssl_certificate /etc/letsencrypt/live/jonobridge.madd.com.mx/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/jonobridge.madd.com.mx/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    # ACME Challenge location for Let's Encrypt
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    # Default service location
    location / {
        proxy_pass http://localhost:5000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name jonobridge.madd.com.mx;

    # ACME Challenge location for Let's Encrypt
    location /.well-known/acme-challenge/ {
        root /var/www/html;
        try_files $uri =404;
    }

    # Redirect all other HTTP traffic to HTTPS
    location / {
        return 301 https://$host$request_uri;
    }
}
```

11. Enable site
```
sudo ln -s /etc/nginx/sites-available/jonobridge.madd.com.mx /etc/nginx/sites-enabled/
```



### Step 7: Test Your TCP/UDP Configuration

1. Check the syntax of your configuration file:

   ```bash
   sudo nginx -t
   sudo systemctl reload nginx
   ```

2. Restart Nginx to apply your changes:

   ```bash
   sudo systemctl restart nginx
   ```

3. Use a UDP client to test if your proxy is working as expected.

---

### Notes

- Monitor Nginx logs for troubleshooting:
  ```bash
  sudo tail -f /var/log/nginx/error.log
  ```


## Final Step
 
***1.Cloning the Repository***  
*Now that SSH authentication is set up, clone your Bitbucket repository.*
```
git clone git@bitbucket.org:maddisontest/jonobridge.git
```

**2.Expected Output:**
```
Cloning into 'jonobridge'...
remote: Counting objects: 5, done.
remote: Total 5 (delta 0), reused 0 (delta 0), pack-reused 0
Unpacking objects: 100% (5/5), done.
```

**3.Navigate to the Repository Directory**
```
cd jonobridge/frontend
```

**Execute this:**
```
python3 -m venv venv
source venv/bin/activate
python main.py 
```

### Clean Docker space from time to time if much space allocated ( also have to restart minikube)
Clean Space:
```
docker system prune -a
```

# Writing the Supervisorctl installation instructions into a Markdown file
installation_guide = """
# Install `Supervisorctl` on Ubuntu

`supervisorctl` is a command-line tool to interact with Supervisor, a process control system that allows you to manage processes on Linux systems. Here’s how you can install and configure it on Ubuntu.

## 1. Update the Package List
First, make sure your system’s package list is up-to-date:
```bash
sudo apt update
```

## 2. Install Supervisor
Install the Supervisor package, which includes `supervisorctl`:
```bash
sudo apt install supervisor
```

## 3. Start and Enable Supervisor Service
Once installed, start the Supervisor service and enable it to start on boot:
```bash
sudo systemctl start supervisor
sudo systemctl enable supervisor
```

## 4. Verify Installation
Check the installed version of Supervisor:
```bash
supervisord --version
```

Verify that `supervisorctl` is working:
```bash
sudo supervisorctl status
```

## 5. Configure Supervisor
1. The main configuration file for Supervisor is located at:
   ```plaintext
   /etc/supervisor/supervisord.conf
   ```

2. Add program configurations by creating `.conf` files in the directory:
   ```plaintext
   /etc/supervisor/conf.d/
   ```

   Create file jonobridge.conf with the following content:
   ```ini
   [program:jonobridge]
   command=/home/ubuntu/jonobridge/frontend/venv/bin/python3 /home/ubuntu/jonobridge/frontend/main.py
   directory=/home/ubuntu/jonobridge/frontend
   user=ubuntu
   autostart=true
   autorestart=true
   stderr_logfile=/var/log/jonobridge.err.log
   stdout_logfile=/var/log/jonobridge.out.log
   environment=PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",HOME="/home/ubuntu",USER="ubuntu",MINIKUBE_HOME="/home/ubuntu/.minikube"
   ```

3. Reload Supervisor to apply changes:
   ```bash
   sudo supervisorctl reread
   sudo supervisorctl update
   ```

## 6. Using `supervisorctl`
- **Check the status of all programs:**  
  ```bash
  sudo supervisorctl status
  ```

- **Start/Stop/Restart a specific program:**  
  ```bash
  sudo supervisorctl start <program-name>
  sudo supervisorctl stop <program-name>
  sudo supervisorctl restart <program-name>
  ```

### Create cerificates:
```
sudo apt update
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d jonobridge.madd.com.mx

Create a new certificate:
sudo certbot certonly --webroot -w /usr/local/nginx/html -d jonobridge.madd.com.mx
```

### Renew:
```
sudo certbot renew --quiet --webroot -w /usr/local/nginx/html && sudo /usr/local/nginx/sbin/nginx -s reload
```

1. If  missing /etc/letsencrypt/options-ssl-nginx.conf file with the recommended SSL options. Here's what we need to do:
```
sudo bash -c 'cat > /etc/letsencrypt/options-ssl-nginx.conf' << 'EOL'
ssl_session_cache shared:le_nginx_SSL:10m;
ssl_session_timeout 1440m;
ssl_session_tickets off;

ssl_protocols TLSv1.2 TLSv1.3;
ssl_prefer_server_ciphers off;

ssl_ciphers "ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384";
```

2. Create the SSL parameters file if it's missing:
```
sudo openssl dhparam -out /etc/letsencrypt/ssl-dhparams.pem 2048
```

3. After creating these files, try testing Nginx again:
```
sudo nginx -t
sudo systemctl reload nginx
```


4. Add directory For certificatechallenge
```
sudo mkdir -p /var/www/html/.well-known/acme-challenge
sudo chown -R www-data:www-data /var/www/html
sudo chmod -R 755 /var/www/html
```

5. Renew certificate

Edit crontab
```
sudo crontab -e
```
add the following line
```
0 0 1 * * certbot renew --quiet --webroot -w /var/www/html --deploy-hook "systemctl reload nginx"
```   

## NGINX files setup for long traffic:


### Fixing "Failed to Allocate Directory Watch: Too Many Open Files" Error

This guide helps you resolve the issue caused by hitting the limit of open file descriptors.

---

### **1. Check Current Limits**
Run the following command to check the current limits for open files:

```bash
ulimit -n
```

For the system-wide limit:

```bash
cat /proc/sys/fs/file-max
```

---

### **2. Increase File Descriptor Limits**
To fix the issue, increase the limit for open files.

#### **Temporary Fix (For Current Session)**
To temporarily increase the limit for the current shell session:
```bash
ulimit -n 65536
```

#### **Permanent Fix (For All Sessions)**

##### a. Update System Configuration
Edit the file `/etc/security/limits.conf`:

```bash
sudo nano /etc/security/limits.conf
```

Add the following lines at the end:

```
* soft nofile 65536
* hard nofile 65536
```

#### b. Configure PAM (Pluggable Authentication Module)
Ensure the PAM limits module is enabled. Edit `/etc/pam.d/common-session` and `/etc/pam.d/common-session-noninteractive`:

```bash
sudo nano /etc/pam.d/common-session
sudo nano /etc/pam.d/common-session-noninteractive
```

Add this line if it’s not already there:

```
session required pam_limits.so
```

#### c. Update Systemd Configuration (For Services Like Nginx)
Edit the Nginx systemd service file:

```bash
sudo systemctl edit nginx
```

Add the following lines under `[Service]`:

```
[Service]
LimitNOFILE=65536
```

Save and reload the systemd daemon:

```bash
sudo systemctl daemon-reexec
```

---

### **3. Verify Changes**
After making the changes, restart your system to apply them globally. Then verify:

```bash
ulimit -n
```

For Nginx-specific settings:

```bash
sudo systemctl show nginx | grep LimitNOFILE
```
---

### **4. Troubleshoot Nginx Further**
If the problem persists, ensure that Nginx isn't keeping too many connections or files open:

- Check Nginx logs for issues: `/var/log/nginx/error.log`
- Check active connections:

```bash
netstat -anp | grep nginx
```

---

If Still issues:

## Immediate fix
sudo sh -c 'echo 8192 > /proc/sys/fs/inotify/max_user_instances'
sudo sh -c 'echo 524288 > /proc/sys/fs/inotify/max_user_watches'

## Persistent fix in /etc/sysctl.d/
echo "fs.inotify.max_user_instances = 8192" | sudo tee -a /etc/sysctl.d/10-user-watches.conf
echo "fs.inotify.max_user_watches = 524288" | sudo tee -a /etc/sysctl.d/10-user-watches.conf

## Apply changes
sudo sysctl -p /etc/sysctl.d/10-user-watches.conf

## Try nginx again
sudo systemctl start nginx




# Configure Permanently Limit Journal Size

Edit the persistent journal config:
```
sudo vim /etc/systemd/journald.conf
```
Uncomment and modify these lines:
```
SystemMaxUse=500M
SystemKeepFree=100M
SystemMaxFileSize=100M
SystemMaxFiles=10
```
Then restart journald:
```
sudo systemctl restart systemd-journald
```

You can free space immediately with:
```
sudo journalctl --vacuum-time=7d
```

Size in elastic search
```
GET _cat/indices?v&h=index,store.size,docs.count
```


### Example Bash Script (deletes indices older than 6 months)

sudo vim /usr/local/bin/cleanup-es.sh

```
#!/bin/bash

# Settings
ES_HOST="http://localhost:9200"
DATE_CUTOFF=$(date -d "6 months ago" +"%Y.%m.%d")

# Optional: define a prefix like "logstash-" or leave it empty
INDEX_PREFIX="logstash-"

# Get all indices matching the prefix
INDICES=$(curl -s "$ES_HOST/_cat/indices?h=index" | grep "^$INDEX_PREFIX")

for index in $INDICES; do
  index_date=$(echo $index | grep -oP '\d{4}\.\d{2}\.\d{2}')
  if [[ "$index_date" < "$DATE_CUTOFF" ]]; then
    echo "Deleting index: $index"
    curl -X DELETE "$ES_HOST/$index"
  fi
done
```

Save this as /usr/local/bin/cleanup-es.sh, and make it executable
```
sudo chmod +x /usr/local/bin/cleanup-es.sh
```

📆 Add to Cron (run daily at 2am)
```
sudo crontab -e
```

Add this line:
```
0 2 * * * /usr/local/bin/cleanup-es.sh >> /var/log/cleanup-es.log 2>&1
```