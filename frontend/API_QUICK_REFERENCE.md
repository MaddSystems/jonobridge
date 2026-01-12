# API Quick Reference Card

## 🚀 Getting Started

```bash
# 1. Install dependency
pip install flasgger

# 2. Start application
python main.py

# 3. Open Swagger UI
http://localhost:5000/apidocs/
```

## 📋 Endpoints Quick Reference

| Endpoint | Method | Purpose | Example |
|----------|--------|---------|---------|
| `/api/v1/diagnostic/namespaces` | GET | List all namespaces | `curl http://localhost:5000/api/v1/diagnostic/namespaces` |
| `/api/v1/diagnostic/diagnose/{ns}` | GET | Diagnose namespace pods | `curl http://localhost:5000/api/v1/diagnostic/diagnose/scaniamx` |
| `/api/v1/diagnostic/restart-pod/{ns}/{pod}` | POST | Restart single pod | `curl -X POST http://localhost:5000/api/v1/diagnostic/restart-pod/scaniamx/listener-xxx` |
| `/api/v1/diagnostic/restart-deployment/{ns}/{dep}` | POST | Restart deployment | `curl -X POST http://localhost:5000/api/v1/diagnostic/restart-deployment/scaniamx/listener` |

## 📊 Status Codes Explained

| Status | Color | Meaning | Action |
|--------|-------|---------|--------|
| `healthy` | 🟢 Green | All pods running | None needed |
| `warning` | 🟡 Yellow | Some pods pending | Monitor/wait |
| `critical` | 🔴 Red | Pods failed | Investigate/restart |
| `empty` | ⚪ Gray | No pods | Normal |

## 🔍 Quick Diagnostics

### Check Everything
```bash
curl http://localhost:5000/api/v1/diagnostic/namespaces
```

**Look for:**
- Any namespace with `status != "healthy"`
- High `failed_pods` count
- `pending_pods` > 0

### Deep Dive Into Problem Namespace
```bash
curl http://localhost:5000/api/v1/diagnostic/diagnose/scaniamx
```

**Look for:**
- Pods with `status != "Running"`
- High `restarts` count (> 10)
- `ready: "0/1"` (not ready)
- Containers with `state: "Waiting"`

### Check Specific Pod
From diagnose output, find pod name and check:
- Is it `Running`?
- How many `restarts`?
- How old (`age`)?
- Container `state`?

## 🔧 Common Fixes

### Pod Is CrashLooping
```bash
# 1. Check logs
kubectl logs <pod-name> -n <namespace>

# 2. Restart pod
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-pod/<namespace>/<pod-name>

# 3. Verify recovery
curl http://localhost:5000/api/v1/diagnostic/diagnose/<namespace>
```

### Service is Down
```bash
# 1. Check deployment status
curl http://localhost:5000/api/v1/diagnostic/diagnose/<namespace>

# 2. Restart entire deployment (zero downtime)
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-deployment/<namespace>/<deployment-name>

# 3. Monitor restart
curl http://localhost:5000/api/v1/diagnostic/diagnose/<namespace>
```

### Multiple Pods Failing
```bash
# 1. Check all namespaces
curl http://localhost:5000/api/v1/diagnostic/namespaces

# 2. Identify problematic namespace
# 3. Diagnose
curl http://localhost:5000/api/v1/diagnostic/diagnose/<namespace>

# 4. Restart deployment (better than pod restart)
curl -X POST http://localhost:5000/api/v1/diagnostic/restart-deployment/<namespace>/<deployment>

# 5. Verify health
curl http://localhost:5000/api/v1/diagnostic/diagnose/<namespace>
```

## 📱 Using Swagger UI (Recommended)

### Step-by-Step

1. **Open** → http://localhost:5000/apidocs/
2. **Find** → Click on endpoint you want to test
3. **Expand** → Click the endpoint row
4. **Input** → Fill in any parameters needed
5. **Execute** → Click blue "Execute" button
6. **View** → See response below

### Example: Diagnose Namespace

1. Open http://localhost:5000/apidocs/
2. Find `/api/v1/diagnostic/diagnose/{namespace}`
3. Click to expand
4. Click "Try it out"
5. Enter `scaniamx` in namespace field
6. Click "Execute"
7. See full response with pod details

## 🔗 Response Structure

### Success Response
```json
{
  "success": true,
  "data": "endpoint-specific"
}
```

### Error Response
```json
{
  "success": false,
  "error": "error description"
}
```

## 📈 Monitoring Dashboard Integration

### JavaScript Example
```javascript
// Get namespace status
fetch('/api/v1/diagnostic/namespaces')
  .then(r => r.json())
  .then(data => {
    if (data.success) {
      data.namespaces.forEach(ns => {
        console.log(`${ns.name}: ${ns.status}`);
      });
    }
  });
```

### Python Example
```python
import requests

response = requests.get('http://localhost:5000/api/v1/diagnostic/namespaces')
if response.ok:
    data = response.json()
    for ns in data.get('namespaces', []):
        print(f"{ns['name']}: {ns['status']}")
```

## 🐛 Troubleshooting

| Problem | Solution |
|---------|----------|
| 404 Namespace not found | Check namespace name with `kubectl get ns` |
| Connection refused | Verify app running: `curl http://localhost:5000` |
| No pods found | Might be empty namespace - check with `kubectl get pods -n <ns>` |
| Restart not working | Verify pod name is exact match from diagnose output |
| Swagger UI not loading | Verify flasgger installed: `pip list \| grep flasgger` |

## ✅ Health Check Script

```bash
#!/bin/bash

echo "🔍 Checking JonoBridge API..."

# Check API is running
echo -n "API Health: "
curl -s http://localhost:5000/api/v1/diagnostic/namespaces | grep -q success && echo "✅ OK" || echo "❌ FAILED"

# Check namespaces
echo -n "Namespaces: "
COUNT=$(curl -s http://localhost:5000/api/v1/diagnostic/namespaces | grep -o '"count":[0-9]*' | cut -d: -f2)
echo "Found $COUNT"

# Check for critical status
echo -n "Health Status: "
curl -s http://localhost:5000/api/v1/diagnostic/namespaces | grep -q '"critical"' && echo "⚠️  Issues found" || echo "✅ Healthy"
```

## 🎯 Best Practices

1. ✅ **Always diagnose before restarting** - Know what's wrong first
2. ✅ **Use deployment restart** - Better than pod restart (zero downtime)
3. ✅ **Check after restart** - Verify recovery was successful
4. ✅ **Monitor restarts** - Track if pods keep failing
5. ✅ **Keep logs** - Review logs after failures for root cause
6. ✅ **Schedule restarts** - Plan maintenance windows for mass restarts
7. ✅ **Document issues** - Keep track of recurring problems

## 📚 Documentation Files

- **API.md** - Complete API documentation with all details
- **INSTALL_API.md** - Installation, setup, and testing guide
- **API_ARCHITECTURE.md** - System architecture and data flows
- **API_IMPLEMENTATION.md** - Implementation summary and details

## 🔐 Security Notes

- All endpoints require authentication (Flask login)
- All operations are logged
- Kubernetes permissions required for operations
- Pod restart operations have 30-second grace period
- No data is stored or cached

## 📞 Getting Help

1. Check documentation in `/frontend/API.md`
2. Test endpoint in Swagger UI at `/apidocs/`
3. Review logs for error details
4. Verify Kubernetes access: `kubectl cluster-info`
5. Check this quick reference card for common issues