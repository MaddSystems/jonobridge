# JonoBridge Diagnostic API

The Diagnostic API provides comprehensive endpoints for monitoring and managing Kubernetes pods and namespaces in the JonoBridge system.

## Overview

The API is built with Flask and includes interactive Swagger documentation. All endpoints are RESTful and return JSON responses.

## Accessing the API

### Via Web Interface

1. Click the **"API Docs"** link in the left navigation menu
2. This opens the Swagger UI with full interactive documentation
3. All endpoints can be tested directly from the Swagger interface

### Via Command Line

Access endpoints directly using curl or any HTTP client:

```bash
# List all namespaces
curl -X GET http://localhost:5000/api/v1/diagnostic/namespaces

# Diagnose a specific namespace
curl -X GET http://localhost:5000/api/v1/diagnostic/diagnose/scaniamx

# Restart a specific pod
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-pod/scaniamx/listener-66ccf94b56-cwvfw

# Restart a deployment
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-deployment/scaniamx/listener
```

## API Endpoints

### 1. List Namespaces

**Endpoint:** `GET /api/v1/diagnostic/namespaces`

**Description:** Lists all Kubernetes namespaces (excluding system namespaces) with pod status summary.

**Response Example:**
```json
{
  "success": true,
  "count": 2,
  "namespaces": [
    {
      "name": "scaniamx",
      "pod_count": 4,
      "running_pods": 4,
      "failed_pods": 0,
      "pending_pods": 0,
      "status": "healthy"
    },
    {
      "name": "clientb",
      "pod_count": 3,
      "running_pods": 2,
      "failed_pods": 1,
      "pending_pods": 0,
      "status": "critical"
    }
  ]
}
```

**Status Codes:**
- `200`: Success
- `500`: Server error

### 2. Diagnose Namespace

**Endpoint:** `GET /api/v1/diagnostic/diagnose/<namespace>`

**Description:** Provides detailed diagnosis of all pods in a specific namespace.

**Parameters:**
- `namespace` (path parameter): Kubernetes namespace name (e.g., `scaniamx`)

**Response Example:**
```json
{
  "success": true,
  "namespace": "scaniamx",
  "overall_status": "healthy",
  "pod_count": 4,
  "pods": [
    {
      "name": "listener-66ccf94b56-cwvfw",
      "status": "Running",
      "ready": "1/1",
      "restarts": 85,
      "age": "80d",
      "labels": {
        "app": "listener"
      },
      "containers": [
        {
          "name": "listener",
          "ready": true,
          "state": "Running",
          "restart_count": 85
        }
      ]
    },
    {
      "name": "meitrackprotocol-9c8f7554b-tp8mw",
      "status": "Running",
      "ready": "1/1",
      "restarts": 92,
      "age": "80d",
      "labels": {
        "app": "meitrackprotocol"
      },
      "containers": [
        {
          "name": "meitrackprotocol",
          "ready": true,
          "state": "Running",
          "restart_count": 92
        }
      ]
    }
  ]
}
```

**Status Information:**
- `healthy`: All pods are running
- `warning`: Some pods are pending or not fully ready
- `critical`: One or more pods are failed
- `empty`: No pods in namespace

**Status Codes:**
- `200`: Success
- `404`: Namespace not found
- `500`: Server error

### 3. Restart Pod

**Endpoint:** `POST /api/v1/diagnostic/restart-pod/<namespace>/<pod_name>`

**Description:** Restarts a specific pod by deleting it. Kubernetes will automatically recreate it based on the deployment configuration.

**Parameters:**
- `namespace` (path parameter): Kubernetes namespace name
- `pod_name` (path parameter): Name of the pod to restart

**Response Example:**
```json
{
  "success": true,
  "message": "Pod listener-66ccf94b56-cwvfw restart initiated",
  "namespace": "scaniamx",
  "pod_name": "listener-66ccf94b56-cwvfw",
  "note": "Pod is being terminated and will be recreated by Kubernetes"
}
```

**Status Codes:**
- `200`: Success
- `400`: Invalid parameters
- `404`: Pod not found
- `500`: Server error

### 4. Restart Deployment

**Endpoint:** `POST /api/v1/diagnostic/restart-deployment/<namespace>/<deployment_name>`

**Description:** Restarts all pods in a deployment using Kubernetes rolling restart. This ensures minimal downtime.

**Parameters:**
- `namespace` (path parameter): Kubernetes namespace name
- `deployment_name` (path parameter): Name of the deployment (e.g., `listener`, `meitrackprotocol`)

**Response Example:**
```json
{
  "success": true,
  "message": "Deployment listener restart initiated",
  "namespace": "scaniamx",
  "deployment_name": "listener",
  "note": "Deployment is performing a rolling restart of all pods"
}
```

**Status Codes:**
- `200`: Success
- `400`: Invalid parameters
- `500`: Server error

## Status Interpretation

### Namespace Status

| Status | Meaning | Action |
|--------|---------|--------|
| `healthy` | All pods are running and ready | No action needed |
| `warning` | Some pods are pending or initializing | Wait a moment, may resolve automatically |
| `critical` | One or more pods are failed | Investigate logs and consider restart |
| `empty` | No pods in namespace | Normal for empty namespaces |

### Pod Status

| Status | Meaning |
|--------|---------|
| `Running` | Pod is active and running |
| `Pending` | Pod is waiting to be scheduled |
| `Failed` | Pod terminated with error |
| `CrashLoopBackOff` | Pod keeps crashing and restarting |
| `Waiting` | Pod is waiting for resources |
| `Terminated` | Pod has stopped |

## Common Use Cases

### Check Overall System Health

```bash
curl -X GET http://localhost:5000/api/v1/diagnostic/namespaces
```

Identify namespaces with issues (status != "healthy") and drill down.

### Deep Dive into a Problem Namespace

```bash
curl -X GET http://localhost:5000/api/v1/diagnostic/diagnose/scaniamx
```

Check individual pod status, restart counts, and container states.

### Restart a Failing Pod

```bash
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-pod/scaniamx/listener-66ccf94b56-cwvfw
```

This removes the pod; Kubernetes recreates it automatically.

### Restart Entire Service

```bash
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-deployment/scaniamx/listener
```

Performs a rolling restart of all instances of the service (zero-downtime).

## Response Format

### Success Response

```json
{
  "success": true,
  "data": { /* endpoint-specific data */ }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Detailed error message"
}
```

## Authentication

All API endpoints require authentication through the same login system as the web interface. Include your session cookie when accessing the API.

### Using curl with Authentication

```bash
# First, login
curl -c cookies.txt -X POST http://localhost:5000/ \
  -d "username=admin&password=yourpassword"

# Then use the API with the saved cookies
curl -b cookies.txt -X GET http://localhost:5000/api/v1/diagnostic/namespaces
```

## Error Handling

### Common Errors

**404 Namespace Not Found**
```json
{
  "success": false,
  "error": "Namespace 'invalid-ns' not found"
}
```

**500 Connection Error**
```json
{
  "success": false,
  "error": "Error listing namespaces: Failed to establish connection to Kubernetes"
}
```

## Performance Considerations

- **Listing namespaces**: ~500ms (includes pod counts)
- **Diagnosing namespace**: ~300ms (includes detailed pod info)
- **Restarting pod**: ~100ms (triggers deletion, actual restart happens asynchronously)
- **Restarting deployment**: ~100ms (kubectl command execution)

## Integration Examples

### Python

```python
import requests
import json

# Get namespace diagnostics
url = "http://localhost:5000/api/v1/diagnostic/diagnose/scaniamx"
response = requests.get(url, cookies={'session': 'your_session_cookie'})
data = response.json()

if data['success']:
    print(f"Namespace: {data['namespace']}")
    print(f"Status: {data['overall_status']}")
    for pod in data['pods']:
        print(f"  - {pod['name']}: {pod['status']}")
else:
    print(f"Error: {data['error']}")
```

### JavaScript

```javascript
// Get namespace diagnostics
fetch('/api/v1/diagnostic/diagnose/scaniamx')
  .then(response => response.json())
  .then(data => {
    if (data.success) {
      console.log(`Namespace: ${data.namespace}`);
      console.log(`Status: ${data.overall_status}`);
      data.pods.forEach(pod => {
        console.log(`  - ${pod.name}: ${pod.status}`);
      });
    } else {
      console.error(`Error: ${data.error}`);
    }
  });

// Restart a pod
fetch('/api/v1/diagnostic/restart-pod/scaniamx/listener-66ccf94b56-cwvfw', {
  method: 'POST'
})
  .then(response => response.json())
  .then(data => console.log(data));
```

## API Documentation

Interactive API documentation is available at:

- **Swagger UI**: `/apidocs/` (recommended)
- **ReDoc**: `/redoc/`
- **Raw OpenAPI Spec**: `/apispec.json`

All endpoints support the Swagger/OpenAPI 2.0 specification and can be tested directly from the web interface.

## Troubleshooting

### Kubernetes Connection Issues

If you get connection errors:

1. Verify kubectl is configured: `kubectl cluster-info`
2. Check kubeconfig location: `echo $KUBECONFIG`
3. Verify credentials: `kubectl auth can-i get pods --all-namespaces`

### CORS Issues

If accessing the API from a different domain:

1. Check CORS headers in responses
2. API is designed for same-origin requests
3. For cross-origin, configure Flask CORS extension

### Authentication Issues

If receiving 403 or 401 errors:

1. Verify login session is valid
2. Check session cookie is being sent
3. Verify user has admin role for restart operations

## Best Practices

1. **Monitor Before Restart**: Always use diagnose endpoint before restarting
2. **Use Deployment Restart**: Prefer deployment restart over individual pod restart for zero-downtime
3. **Check Logs**: After restart, monitor logs to verify recovery
4. **Document Issues**: Keep track of frequent failures for root cause analysis
5. **Schedule Restarts**: Plan maintenance windows for mass restarts if needed