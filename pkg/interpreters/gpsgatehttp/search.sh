#!/bin/bash

# SSH Routine to Monitor and Store GPS Gate HTTP Pod Logs in Namespace 'gpsgatecortedecorriente'
# This script finds the current pod for app 'gpsgatehttp', streams logs with timestamps to a file,
# and exits when the pod is killed, then reports the log file.

NAMESPACE="gpsgatecortedecorriente"
APP_LABEL="app=gpsgatehttp"

echo "Starting pod log monitoring for namespace: $NAMESPACE, app: $APP_LABEL"
echo "Current time: $(date)"

# Find the current pod name (assumes one pod; adjust if multiple)
POD_NAME=$(kubectl get pods -n $NAMESPACE -l $APP_LABEL -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -z "$POD_NAME" ]; then
    echo "Error: No pod found for app '$APP_LABEL' in namespace '$NAMESPACE'. Check if the pod exists."
    exit 1
fi

LOG_FILE="log_${POD_NAME}.txt"
echo "Found pod: $POD_NAME"
echo "Streaming logs with timestamps to file: $LOG_FILE"
echo "Press Ctrl+C to stop early (logs will still be saved up to that point)"

# Stream logs with timestamps to file continuously
# --timestamps adds log timestamps; -f follows the logs until pod is killed/deleted
kubectl logs -n $NAMESPACE -f --timestamps $POD_NAME > "$LOG_FILE"

# After logs stop, extract error/fatal/panic/exit lines for diagnosis
ERROR_LOG_FILE="errors_${POD_NAME}.txt"
grep -iE "error|fatal|panic|exit" "$LOG_FILE" > "$ERROR_LOG_FILE"

# After logs stop (pod killed/deleted), print the file name and check pod status

echo "Pod logs have stopped (likely pod killed or deleted)."
echo "This is the log file: $LOG_FILE"
echo "Error summary file: $ERROR_LOG_FILE"
echo "Final time: $(date)"

# Smoking gun: get full pod status and container state
kubectl get pod $POD_NAME -n $NAMESPACE -o yaml > pod_${POD_NAME}_full.yaml
kubectl describe pod $POD_NAME -n $NAMESPACE > pod_${POD_NAME}_describe.txt

echo "Full pod YAML: pod_${POD_NAME}_full.yaml"
echo "Pod describe output: pod_${POD_NAME}_describe.txt"

# Print pod status.reason and status.message if present
echo "Pod status.reason and status.message:"
kubectl get pod $POD_NAME -n $NAMESPACE -o json | jq '.status.reason, .status.message'

# Get ReplicaSet and Deployment names
RS_NAME=$(kubectl get pod $POD_NAME -n $NAMESPACE -o jsonpath='{.metadata.ownerReferences[0].name}')
DEPLOY_NAME=$(kubectl get rs $RS_NAME -n $NAMESPACE -o jsonpath='{.metadata.ownerReferences[0].name}')

# Show last 10 events for ReplicaSet and Deployment
echo "Last 10 events for ReplicaSet $RS_NAME:"
kubectl get events -n $NAMESPACE --field-selector involvedObject.name=$RS_NAME | tail -10

echo "Last 10 events for Deployment $DEPLOY_NAME:"
kubectl get events -n $NAMESPACE --field-selector involvedObject.name=$DEPLOY_NAME | tail -10

# Get last 30 events for the namespace for context
echo "Recent cluster events (last 30):"
kubectl get events -n $NAMESPACE --sort-by=.metadata.creationTimestamp | tail -30 > events_${POD_NAME}.txt
cat events_${POD_NAME}.txt

echo "Check pod_${POD_NAME}_describe.txt, events_${POD_NAME}.txt, and above ReplicaSet/Deployment events for the exact reason Kubernetes killed the pod (look for OOMKilled, node pressure, controller actions, etc.)"