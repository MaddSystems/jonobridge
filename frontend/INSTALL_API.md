# API Implementation Guide

This guide explains how to install, run, and test the new Diagnostic API.

## Installation

### 1. Install Dependencies

```bash
cd ~/jonobridge/frontend

# Install flasgger for Swagger documentation
pip install flasgger
```

Or update all dependencies:

```bash
pip install -r requirements.txt
```

### 2. Verify Installation

```bash
python -c "from flasgger import Swagger; print('Flasgger installed successfully')"
```

## Running the Application

### Start with Default Flask

```bash
python main.py
```

The application will start at `http://0.0.0.0:5000`

### Access Points

- **Web Interface**: `http://localhost:5000`
- **API Docs (Swagger)**: `http://localhost:5000/apidocs/`
- **API Docs (ReDoc)**: `http://localhost:5000/redoc/`
- **API Link in UI**: Click "API Docs" in the left navigation menu

## Testing the API

### Test 1: List All Namespaces

```bash
curl -X GET http://localhost:5000/api/v1/diagnostic/namespaces
```

Expected response (200 OK):
```json
{
  "success": true,
  "count": 1,
  "namespaces": [
    {
      "name": "scaniamx",
      "pod_count": 4,
      "running_pods": 4,
      "failed_pods": 0,
      "pending_pods": 0,
      "status": "healthy"
    }
  ]
}
```

### Test 2: Diagnose a Namespace

```bash
curl -X GET http://localhost:5000/api/v1/diagnostic/diagnose/scaniamx
```

Expected response (200 OK):
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
      "containers": [...]
    }
  ]
}
```

### Test 3: Restart a Pod

```bash
# Get pod name from diagnose response first, then:
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-pod/scaniamx/listener-66ccf94b56-cwvfw
```

Expected response (200 OK):
```json
{
  "success": true,
  "message": "Pod listener-66ccf94b56-cwvfw restart initiated",
  "namespace": "scaniamx",
  "pod_name": "listener-66ccf94b56-cwvfw",
  "note": "Pod is being terminated and will be recreated by Kubernetes"
}
```

### Test 4: Restart a Deployment

```bash
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-deployment/scaniamx/listener
```

Expected response (200 OK):
```json
{
  "success": true,
  "message": "Deployment listener restart initiated",
  "namespace": "scaniamx",
  "deployment_name": "listener",
  "note": "Deployment is performing a rolling restart of all pods"
}
```

## Using Swagger UI

### Access Swagger

1. Start the application
2. Open browser to `http://localhost:5000/apidocs/`
3. You should see the interactive Swagger documentation

### Test Endpoints in Swagger

1. Find an endpoint in the list
2. Click "Try it out"
3. Fill in required parameters
4. Click "Execute"
5. View the response below

### Example: Testing in Swagger

1. Go to `/api/v1/diagnostic/diagnose/{namespace}`
2. Click "Try it out"
3. Enter `scaniamx` in the namespace field
4. Click "Execute"
5. See the full response with all pod details

## Troubleshooting

### Issue: "Cannot connect to Kubernetes"

**Error Response:**
```json
{
  "success": false,
  "error": "Error listing namespaces: Failed to establish connection"
}
```

**Solution:**
```bash
# Verify kubectl is working
kubectl cluster-info

# Check kubeconfig
echo $KUBECONFIG

# Verify API can reach Kubernetes
kubectl get pods -n scaniamx
```

### Issue: "Namespace not found"

**Error Response:**
```json
{
  "success": false,
  "error": "Namespace 'invalid-ns' not found"
}
```

**Solution:**
```bash
# List available namespaces
kubectl get namespaces

# Use correct namespace name
curl -X GET http://localhost:5000/api/v1/diagnostic/diagnose/scaniamx
```

### Issue: Pod restart doesn't work

**Possible Causes:**

1. Pod doesn't exist - check pod name in diagnose endpoint
2. Missing Kubernetes permissions - verify kubectl can delete pods
3. Deployment protection - some deployments might be protected

**Solution:**
```bash
# Test with kubectl directly
kubectl delete pod listener-66ccf94b56-cwvfw -n scaniamx

# If that works, API should work too
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-pod/scaniamx/listener-66ccf94b56-cwvfw
```

## File Structure

```
frontend/
├── main.py                    # Main Flask application (UPDATED)
├── requirements.txt           # Dependencies (UPDATED - added flasgger)
├── api/
│   ├── __init__.py           # API package marker
│   └── diagnostic_api.py     # Diagnostic API blueprint (NEW)
├── templates/
│   └── base.html             # Base template (UPDATED - added API Docs link)
└── API.md                     # Full API documentation (NEW)
```

## Key Changes Made

### 1. `main.py`
- Added Flasgger import
- Initialized Swagger with metadata
- Registered diagnostic blueprint
- Added `/api-docs` route for UI redirect

### 2. `requirements.txt`
- Added `flasgger` dependency

### 3. `templates/base.html`
- Added "API Docs" link in left navigation menu
- Uses Font Awesome icon for API

### 4. New Files
- `api/diagnostic_api.py` - Complete API implementation
- `api/__init__.py` - Package marker
- `API.md` - Full API documentation

## Deployment Notes

### With Gunicorn (Production)

```bash
gunicorn -w 4 -b 0.0.0.0:5000 main:app
```

The API will be available at:
- `http://server:5000/api/v1/diagnostic/...`
- `http://server:5000/apidocs/`

### With Docker

Add to Dockerfile:
```dockerfile
RUN pip install flasgger
```

## Performance Tips

1. **Caching**: Consider implementing Redis caching for namespace lists
2. **Async**: Use async tasks for restart operations
3. **Rate Limiting**: Consider adding rate limiting for production
4. **Monitoring**: Log all API calls for audit trail

## Security Considerations

1. **Authentication**: Already handled by Flask login_required decorator
2. **RBAC**: Consider implementing role-based access control
3. **Audit Logging**: All restart operations are logged
4. **API Keys**: Could implement optional API key authentication

## Next Steps

1. Deploy the updated code
2. Install flasgger: `pip install flasgger`
3. Restart the application
4. Visit `/apidocs/` to view API documentation
5. Test endpoints using Swagger UI
6. Integrate into monitoring dashboard