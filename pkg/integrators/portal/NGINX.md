# Restructure setup:

1. First, modify your main Nginx configuration to include files from a specific directory:

```
server {
    server_name jonobridge.madd.com.mx;
    
    # Location for the default service
    location / {
        proxy_pass http://localhost:5000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Include all endpoint configurations from endpoints directory
    include /etc/nginx/endpoints/*.conf;
    
    # Enable SSL for the server
    listen 443 ssl;
    ssl_certificate /etc/letsencrypt/live/jonobridge.madd.com.mx/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/jonobridge.madd.com.mx/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;
}

server {
    if ($host = jonobridge.madd.com.mx) {
        return 301 https://$host$request_uri;
    }
    server_name jonobridge.madd.com.mx;
    return 404;
}
```

2. Create a directory for your endpoint configurations:
```
mkdir -p /etc/nginx/endpoints
```

3. Create individual configuration files for each endpoint in this directory:
```
For example, create /etc/nginx/endpoints/semov.conf:
```

### Contents:

```
location /api/v1.0/lobosoftware {
    proxy_pass http://192.168.49.2:31871/api/v1.0/lobosoftware;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

4. After adding or removing configuration files, reload Nginx to apply changes:

```
sudo nginx -t     # Test the configuration
sudo nginx -s reload  # Reload if the test passes
```