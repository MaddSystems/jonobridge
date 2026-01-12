# Diagnostic API Implementation Summary

## What Was Implemented

A comprehensive **Diagnostic API** has been added to JonoBridge that allows you to:

1. **List all Kubernetes namespaces** with pod health status
2. **Diagnose individual namespaces** with detailed pod information
3. **Restart individual pods** with automatic recreation
4. **Restart entire deployments** with rolling updates

All endpoints include:
- ✅ Full Swagger/OpenAPI documentation
- ✅ Interactive testing interface
- ✅ Error handling and validation
- ✅ Kubernetes integration
- ✅ Real-time pod status information

## Quick Start

### Installation

```bash
# Install the new dependency
pip install flasgger

# Or update all dependencies
pip install -r requirements.txt
```

### Running

```bash
python main.py
```

### Access Points

| Location | URL | Purpose |
|----------|-----|---------|
| **Web UI** | http://localhost:5000 | Main JonoBridge interface |
| **API Docs** | http://localhost:5000/apidocs/ | Interactive Swagger UI |
| **API Menu** | Click "API Docs" in sidebar | Link from web interface |

## API Endpoints

### 1. List Namespaces
```bash
GET /api/v1/diagnostic/namespaces
```
Returns all namespaces with pod count and health status.

### 2. Diagnose Namespace
```bash
GET /api/v1/diagnostic/diagnose/{namespace}
```
Returns detailed information about all pods in a namespace.

### 3. Restart Pod
```bash
POST /api/v1/diagnostic/restart-pod/{namespace}/{pod_name}
```
Restarts a specific pod by deleting it (Kubernetes recreates automatically).

### 4. Restart Deployment
```bash
POST /api/v1/diagnostic/restart-deployment/{namespace}/{deployment_name}
```
Performs a rolling restart of all pods in a deployment.

## Usage Example

### Using Curl

```bash
# List all namespaces
curl http://localhost:5000/api/v1/diagnostic/namespaces

# Check namespace health
curl http://localhost:5000/api/v1/diagnostic/diagnose/scaniamx

# Restart a pod
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-pod/scaniamx/listener-66ccf94b56-cwvfw

# Restart all instances of a service
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-deployment/scaniamx/listener
```

### Using Swagger UI

1. Open http://localhost:5000/apidocs/
2. Click on an endpoint
3. Click "Try it out"
4. Fill in parameters
5. Click "Execute"
6. View results

## Response Examples

### Healthy Namespace
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
    }
  ]
}
```

### Pod Diagnosis
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
      "containers": [
        {
          "name": "listener",
          "ready": true,
          "state": "Running",
          "restart_count": 85
        }
      ]
    }
  ]
}
```

## Key Features

### 1. Real-time Status Detection

The API automatically detects:
- ✅ Running pods
- ✅ Pending pods
- ✅ Failed pods
- ✅ CrashLoopBackOff state
- ✅ Container restart counts
- ✅ Pod age

### 2. Intelligent Status Classification

```
Status       Meaning                          Recommended Action
────────────────────────────────────────────────────────────────
healthy      All pods running                 None needed
warning      Some pods pending/not ready      Monitor and wait
critical     One or more pods failed          Investigate and restart
empty        No pods in namespace             Normal
```

### 3. Pod Restart Strategies

**Individual Pod Restart:**
- Direct deletion with 30-second grace period
- Kubernetes automatically recreates from ReplicaSet
- Minimal overhead but single pod impact

**Deployment Restart:**
- Uses Kubernetes rolling restart
- Gradually terminates old pods
- Gradually starts new pods
- **Zero downtime**
- **Recommended for production**

### 4. Comprehensive Information

Each endpoint provides:
- Pod names and current states
- Ready status (1/1, 0/1, etc.)
- Restart count
- Pod age
- Container details
- Pod labels

## Architecture

```
User/Dashboard
    ↓
Swagger UI (/apidocs)
    ↓
API Blueprint (/api/v1/diagnostic)
    ↓
Kubernetes Python Client
    ↓
Kubernetes API / kubectl
    ↓
Minikube/Cluster
```

## Files Modified/Created

### New Files
- `api/diagnostic_api.py` - Complete API implementation (320+ lines)
- `api/__init__.py` - Package marker
- `API.md` - Full API documentation
- `INSTALL_API.md` - Installation and testing guide
- `API_ARCHITECTURE.md` - Architecture diagrams and flows

### Modified Files
- `main.py` - Added Swagger initialization and blueprint registration
- `requirements.txt` - Added `flasgger` dependency
- `templates/base.html` - Added "API Docs" link in navigation

## Benefits

1. **Easy Diagnostics** - Know exactly what's wrong with your services
2. **Quick Restarts** - Restart pods/deployments without manual kubectl
3. **Rolling Updates** - Zero-downtime deployments
4. **Swagger Docs** - Interactive testing interface
5. **REST Standard** - Works with any HTTP client
6. **Built-in Logging** - All operations are logged

## Common Use Cases

### Monitor Overall Health
```bash
curl http://localhost:5000/api/v1/diagnostic/namespaces
```
Check all namespaces at a glance.

### Investigate Issue
```bash
curl http://localhost:5000/api/v1/diagnostic/diagnose/problematic-namespace
```
See which pods are failing and why.

### Fix Failing Service
```bash
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-deployment/ns/service-name
```
Restart service with zero downtime.

### Check After Restart
```bash
curl http://localhost:5000/api/v1/diagnostic/diagnose/ns
```
Verify all pods recovered properly.

## Testing Checklist

- [ ] Install flasgger: `pip install flasgger`
- [ ] Start application: `python main.py`
- [ ] Access Swagger: http://localhost:5000/apidocs/
- [ ] Test list namespaces endpoint
- [ ] Test diagnose endpoint with your namespace
- [ ] Verify pod information is correct
- [ ] Test pod restart (use test pod)
- [ ] Verify pod was recreated
- [ ] Test deployment restart
- [ ] Verify rolling restart worked

## Troubleshooting

### Flasgger not installed
```bash
pip install flasgger
```

### Kubernetes connection fails
```bash
# Verify kubectl works
kubectl cluster-info

# Check kubeconfig
echo $KUBECONFIG

# Test pod access
kubectl get pods -n scaniamx
```

### Wrong namespace in response
```bash
# Verify namespace exists
kubectl get namespaces

# Use exact namespace name
curl http://localhost:5000/api/v1/diagnostic/diagnose/scaniamx
```

### Pod restart not working
```bash
# Verify pod exists
kubectl get pods -n scaniamx

# Try with kubectl directly
kubectl delete pod <pod-name> -n scaniamx

# If that works, API should work
```

## Next Steps

1. **Deploy** - Install flasgger and restart application
2. **Test** - Use Swagger UI to test endpoints
3. **Integrate** - Add to your monitoring dashboard
4. **Monitor** - Set up automated health checks
5. **Extend** - Add more diagnostic endpoints as needed

## Additional Resources

- **Full Documentation**: See `API.md` for complete endpoint details
- **Installation Guide**: See `INSTALL_API.md` for setup and testing
- **Architecture**: See `API_ARCHITECTURE.md` for system design
- **Swagger UI**: http://localhost:5000/apidocs/ for interactive docs

## Support

For issues:
1. Check logs: `grep -i error logs/`
2. Verify Kubernetes access: `kubectl auth can-i get pods`
3. Test endpoint manually: `curl -v http://localhost:5000/api/v1/diagnostic/namespaces`
4. Review documentation in `API.md`