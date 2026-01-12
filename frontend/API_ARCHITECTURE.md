# Diagnostic API Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     JonoBridge Frontend                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────────┐          ┌──────────────────┐             │
│  │   Web UI         │          │   API UI         │             │
│  │   (HTML/JS)      │          │   (Swagger)      │             │
│  │                  │          │                  │             │
│  │ - Clients        │          │ - Interactive    │             │
│  │ - Setup          │          │ - Documentation  │             │
│  │ - Deploy         │          │ - Try it out     │             │
│  │ - Status         │          │ - Examples       │             │
│  └────────┬─────────┘          └────────┬─────────┘             │
│           │                             │                        │
│           └─────────────┬───────────────┘                        │
│                         │                                        │
│              ┌──────────▼────────────┐                          │
│              │   main.py            │                          │
│              │  Flask Application   │                          │
│              ├──────────────────────┤                          │
│              │ Routes:              │                          │
│              │ - /clients           │                          │
│              │ - /status            │                          │
│              │ - /deploy            │                          │
│              │ - /api-docs          │                          │
│              └──────────┬───────────┘                          │
│                         │                                        │
│              ┌──────────▼────────────────────┐                 │
│              │   Diagnostic API Blueprint    │                 │
│              │   (api/diagnostic_api.py)     │                 │
│              ├──────────────────────────────┤                 │
│              │ Routes:                      │                 │
│              │ - /api/v1/diagnostic/...     │                 │
│              ├──────────────────────────────┤                 │
│              │ Endpoints:                   │                 │
│              │ • namespaces                 │                 │
│              │ • diagnose/<namespace>       │                 │
│              │ • restart-pod/..             │                 │
│              │ • restart-deployment/..      │                 │
│              └──────────┬───────────────────┘                 │
│                         │                                        │
│              ┌──────────▼──────────┐                           │
│              │   Kubernetes        │                           │
│              │   Python Client     │                           │
│              └──────────┬──────────┘                           │
│                         │                                        │
└─────────────────────────┼────────────────────────────────────────┘
                          │
          ┌───────────────▼───────────────┐
          │                               │
      ┌───▼───┐                      ┌────▼────┐
      │kubectl│                      │ Minikube │
      │ API   │                      │/Cluster  │
      └───┬───┘                      └────┬─────┘
          │                               │
          └───────────────┬───────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
    ┌───▼────┐        ┌───▼────┐      ┌───▼────┐
    │Pods    │        │Services│      │Deployments
    │Running │        │        │      │
    │CrashLoop
    │Failed  │        │        │      │
    └────────┘        └────────┘      └────────┘
```

## API Call Flow

```
User Action
    │
    ▼
┌─────────────────────────┐
│ Click "API Docs" Button │
│ or Visit /apidocs/      │
└────────────┬────────────┘
             │
             ▼
    ┌────────────────────┐
    │  Swagger UI Loads  │
    │  - Lists endpoints │
    │  - Shows examples  │
    └────────┬───────────┘
             │
             ▼
    ┌────────────────────────────────┐
    │ User Selects Endpoint          │
    │ Example: /diagnose/<namespace> │
    └────────┬───────────────────────┘
             │
             ▼
    ┌────────────────────┐
    │ Fills Parameters   │
    │ - namespace: scaniamx
    └────────┬───────────┘
             │
             ▼
    ┌────────────────────┐
    │ Clicks "Execute"   │
    └────────┬───────────┘
             │
             ▼
    ┌────────────────────────────┐
    │ Request Sent to Backend    │
    │ GET /api/v1/diagnostic/... │
    └────────┬───────────────────┘
             │
             ▼
    ┌────────────────────────────┐
    │ API Blueprint Receives     │
    │ - Validates parameters     │
    │ - Checks authentication    │
    └────────┬───────────────────┘
             │
             ▼
    ┌──────────────────────────────┐
    │ Kubernetes Client Connects   │
    │ - Load kubeconfig            │
    │ - Query API                  │
    └────────┬─────────────────────┘
             │
             ▼
    ┌──────────────────────────────┐
    │ Process Response             │
    │ - Parse JSON                 │
    │ - Format data                │
    │ - Determine status           │
    └────────┬─────────────────────┘
             │
             ▼
    ┌──────────────────────────────┐
    │ Return JSON Response         │
    │ {                            │
    │   success: true,             │
    │   namespace: "scaniamx",     │
    │   pods: [...],               │
    │   overall_status: "healthy"  │
    │ }                            │
    └────────┬─────────────────────┘
             │
             ▼
    ┌──────────────────────────────┐
    │ Swagger UI Displays Result   │
    │ - Pretty printed JSON        │
    │ - Status code 200            │
    │ - Response time              │
    └──────────────────────────────┘
```

## Endpoint Call Examples

### 1. List Namespaces Flow

```
GET /api/v1/diagnostic/namespaces
           │
           ▼
   Load Kubernetes Config
           │
           ▼
   List All Namespaces
           │
           ▼
   For Each Namespace:
   ├─ Get Pod Count
   ├─ Count Running Pods
   ├─ Count Failed Pods
   ├─ Determine Status
   └─ Add to Response
           │
           ▼
   Return [
     { name: "scaniamx", pods: 4, status: "healthy" },
     { name: "clientb", pods: 3, status: "critical" }
   ]
```

### 2. Diagnose Namespace Flow

```
GET /api/v1/diagnostic/diagnose/scaniamx
           │
           ▼
   Verify Namespace Exists
           │
           ▼
   List All Pods in Namespace
           │
           ▼
   For Each Pod:
   ├─ Get Pod Name
   ├─ Get Pod Status
   ├─ Get Ready Status
   ├─ Count Restarts
   ├─ Calculate Age
   ├─ Get Container Info
   │  ├─ Name
   │  ├─ Ready State
   │  ├─ Current State
   │  └─ Restart Count
   └─ Add to Pods Array
           │
           ▼
   Determine Overall Status:
   ├─ All Running? → "healthy"
   ├─ Some Pending? → "warning"
   └─ Has Failed? → "critical"
           │
           ▼
   Return {
     namespace: "scaniamx",
     overall_status: "healthy",
     pods: [
       { name: "listener-...", status: "Running", ... },
       { name: "meitrack-...", status: "Running", ... },
       ...
     ]
   }
```

### 3. Restart Pod Flow

```
POST /api/v1/diagnostic/restart-pod/scaniamx/listener-66...
           │
           ▼
   Validate Parameters
   ├─ Namespace not empty?
   └─ Pod name not empty?
           │
           ▼
   Load Kubernetes Config
           │
           ▼
   Verify Pod Exists
   └─ Try to Read Pod Metadata
           │
           ▼
   Delete Pod with Grace Period (30s)
           │
           ▼
   Log Success
           │
           ▼
   Return {
     success: true,
     message: "Pod restart initiated",
     pod_name: "listener-66...",
     note: "Pod will be recreated by Kubernetes"
   }
           │
           ▼
   Kubernetes Detects Missing Pod
           │
           ▼
   ReplicaSet Creates New Pod
           │
           ▼
   Pod Transitions: Pending → Running
```

### 4. Restart Deployment Flow

```
POST /api/v1/diagnostic/restart-deployment/scaniamx/listener
           │
           ▼
   Validate Parameters
           │
           ▼
   Execute kubectl Command:
   kubectl rollout restart deployment/listener -n scaniamx
           │
           ▼
   kubectl:
   ├─ Gets Deployment
   ├─ Updates Pod Template Spec
   │  (forces new rollout)
   └─ Triggers Rolling Restart
           │
           ▼
   Kubernetes:
   ├─ Creates New ReplicaSet
   ├─ Gradually Starts New Pods
   ├─ Gradually Terminates Old Pods
   └─ Maintains Service Availability
           │
           ▼
   Return {
     success: true,
     message: "Deployment restart initiated",
     deployment_name: "listener",
     note: "Rolling restart in progress"
   }
```

## Status Status Determination Logic

```
Check Pod Statuses
    │
    ├─ No pods? → "empty"
    │
    ├─ Any Failed/CrashLoopBackOff? → "critical"
    │
    ├─ Any Pending? → "warning"
    │
    ├─ All Running? → "healthy"
    │
    └─ Mixed state? → "warning"
```

## Security Flow

```
Request to API Endpoint
         │
         ▼
Check Authentication
├─ Session Cookie Exists?
├─ Session Is Valid?
└─ User Logged In?
         │
         ▼ (If authenticated)
Load Kubernetes Config
         │
         ▼
Execute Request
         │
         ▼
Return Response
```

## Error Handling Flow

```
Request
  │
  ├─ Invalid Parameters?
  │  └─ Return 400 Bad Request
  │
  ├─ Not Authenticated?
  │  └─ Return 403 Forbidden
  │
  ├─ Namespace Not Found?
  │  └─ Return 404 Not Found
  │
  ├─ Kubernetes Error?
  │  ├─ Log Error
  │  └─ Return 500 Server Error
  │
  └─ Success?
     └─ Return 200 OK with Data
```

## Data Flow from Kubernetes

```
Kubernetes API
    │
    ├─ Namespace Objects
    │  └─ metadata.name
    │
    ├─ Pod Objects
    │  ├─ metadata.name
    │  ├─ metadata.creationTimestamp
    │  ├─ metadata.labels
    │  ├─ status.phase (Running, Pending, Failed)
    │  ├─ status.containerStatuses[]
    │  │  ├─ name
    │  │  ├─ ready (boolean)
    │  │  ├─ restartCount
    │  │  └─ state (running, waiting, terminated)
    │  └─ spec.containers[] (count)
    │
    └─ Deployment Objects
       └─ metadata.name
```

## Response Data Structure

```json
{
  "success": boolean,
  "error": "error message (if success=false)",
  
  // For /namespaces endpoint
  "count": number,
  "namespaces": [
    {
      "name": string,
      "pod_count": number,
      "running_pods": number,
      "failed_pods": number,
      "pending_pods": number,
      "status": "healthy|warning|critical|empty"
    }
  ],
  
  // For /diagnose endpoint
  "namespace": string,
  "overall_status": string,
  "pod_count": number,
  "pods": [
    {
      "name": string,
      "status": string,
      "ready": string,
      "restarts": number,
      "age": string,
      "labels": object,
      "containers": [
        {
          "name": string,
          "ready": boolean,
          "state": string,
          "restart_count": number
        }
      ]
    }
  ],
  
  // For /restart-pod and /restart-deployment endpoints
  "message": string,
  "namespace": string,
  "pod_name": string,  // or deployment_name
  "note": string
}
```