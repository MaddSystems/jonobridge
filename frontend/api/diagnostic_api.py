"""
Diagnostic API for Kubernetes pod monitoring and management.
Provides endpoints to list namespaces, diagnose pod status, and restart pods.
"""

from flask import Blueprint, jsonify, request
from kubernetes import client as k8s_client, config
import subprocess
import logging

# Create blueprint
diagnostic_bp = Blueprint('diagnostic_api', __name__, url_prefix='/api/v1/diagnostic')

# Setup logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


@diagnostic_bp.route('/namespaces', methods=['GET'])
def list_namespaces():
    """
    List all Kubernetes namespaces with pod counts and status.
    
    ---
    tags:
      - Diagnostic
    summary: List all Kubernetes namespaces
    description: Returns a list of all namespaces with pod counts and health status
    responses:
      200:
        description: List of namespaces retrieved successfully
        schema:
          type: object
          properties:
            success:
              type: boolean
              example: true
            namespaces:
              type: array
              items:
                type: object
                properties:
                  name:
                    type: string
                    example: scaniamx
                  pod_count:
                    type: integer
                    example: 4
                  running_pods:
                    type: integer
                    example: 4
                  failed_pods:
                    type: integer
                    example: 0
                  status:
                    type: string
                    enum: [healthy, warning, critical]
                    example: healthy
      500:
        description: Error retrieving namespaces
    """
    try:
        try:
            config.load_kube_config()
        except Exception as kube_error:
            logger.error(f"Failed to load kube config: {str(kube_error)}")
            # Try to load from service account
            config.load_incluster_config()
        
        v1 = k8s_client.CoreV1Api()
        
        namespaces = []
        
        # Get all namespaces
        ns_list = v1.list_namespace()
        
        for ns in ns_list.items:
            ns_name = ns.metadata.name
            
            # Skip system namespaces
            if ns_name.startswith('kube-') or ns_name == 'default':
                continue
            
            try:
                # Get pods in this namespace
                pods = v1.list_namespaced_pod(ns_name)
                
                total_pods = len(pods.items)
                running_pods = sum(1 for p in pods.items if p.status.phase == 'Running')
                failed_pods = sum(1 for p in pods.items if p.status.phase in ['Failed', 'CrashLoopBackOff'])
                pending_pods = sum(1 for p in pods.items if p.status.phase == 'Pending')
                
                # Determine status
                if total_pods == 0:
                    status = "empty"
                elif failed_pods > 0:
                    status = "critical"
                elif pending_pods > 0:
                    status = "warning"
                elif running_pods == total_pods:
                    status = "healthy"
                else:
                    status = "warning"
                
                namespaces.append({
                    'name': ns_name,
                    'pod_count': total_pods,
                    'running_pods': running_pods,
                    'failed_pods': failed_pods,
                    'pending_pods': pending_pods,
                    'status': status
                })
            except k8s_client.rest.ApiException as e:
                logger.warning(f"Error listing pods in namespace {ns_name}: {e}")
        
        return jsonify({
            'success': True,
            'count': len(namespaces),
            'namespaces': sorted(namespaces, key=lambda x: x['name'])
        }), 200
        
    except Exception as e:
        logger.error(f"Error listing namespaces: {str(e)}")
        return jsonify({
            'success': False,
            'error': f"Error listing namespaces: {str(e)}"
        }), 500


@diagnostic_bp.route('/diagnose/<namespace>', methods=['GET'])
def diagnose_namespace(namespace):
    """
    Diagnose the health status of all pods in a namespace.
    
    ---
    tags:
      - Diagnostic
    summary: Diagnose namespace pod status
    description: Returns detailed health information for all pods in a namespace
    parameters:
      - name: namespace
        in: path
        type: string
        required: true
        description: The Kubernetes namespace name
        example: scaniamx
    responses:
      200:
        description: Namespace diagnosis completed successfully
        schema:
          type: object
          properties:
            success:
              type: boolean
              example: true
            namespace:
              type: string
              example: scaniamx
            overall_status:
              type: string
              enum: [healthy, warning, critical]
              example: healthy
            pods:
              type: array
              items:
                type: object
                properties:
                  name:
                    type: string
                    example: listener-66ccf94b56-cwvfw
                  status:
                    type: string
                    enum: [Running, Pending, Failed, CrashLoopBackOff]
                    example: Running
                  ready:
                    type: string
                    example: "1/1"
                  restarts:
                    type: integer
                    example: 85
                  age:
                    type: string
                    example: 80d
                  containers:
                    type: array
                    items:
                      type: object
                      properties:
                        name:
                          type: string
                        ready:
                          type: boolean
                        state:
                          type: string
      404:
        description: Namespace not found
      500:
        description: Error diagnosing namespace
    """
    try:
        try:
            config.load_kube_config()
        except Exception as kube_error:
            logger.error(f"Failed to load kube config: {str(kube_error)}")
            # Try to load from service account
            config.load_incluster_config()
        
        v1 = k8s_client.CoreV1Api()
        
        try:
            # Check if namespace exists
            v1.read_namespace(namespace)
        except k8s_client.rest.ApiException as e:
            if e.status == 404:
                return jsonify({
                    'success': False,
                    'error': f"Namespace '{namespace}' not found"
                }), 404
            raise
        
        # Get all pods in namespace
        pods = v1.list_namespaced_pod(namespace)
        
        pod_details = []
        has_issues = False
        
        for pod in pods.items:
            pod_name = pod.metadata.name
            pod_status = pod.status.phase
            
            # Check if pod has issues
            if pod_status not in ['Running']:
                has_issues = True
            
            # Get container info
            containers = []
            if pod.status.container_statuses:
                for container in pod.status.container_statuses:
                    containers.append({
                        'name': container.name,
                        'ready': container.ready,
                        'state': _get_container_state(container),
                        'restart_count': container.restart_count
                    })
            
            # Calculate age
            age_str = _calculate_age(pod.metadata.creation_timestamp)
            
            # Get restart count
            restarts = sum(c.restart_count for c in pod.status.container_statuses) if pod.status.container_statuses else 0
            
            # Get ready status
            ready_count = sum(1 for c in pod.status.container_statuses if c.ready) if pod.status.container_statuses else 0
            total_containers = len(pod.spec.containers)
            ready_str = f"{ready_count}/{total_containers}"
            
            pod_details.append({
                'name': pod_name,
                'status': pod_status,
                'ready': ready_str,
                'restarts': restarts,
                'age': age_str,
                'containers': containers,
                'labels': dict(pod.metadata.labels) if pod.metadata.labels else {}
            })
        
        # Determine overall status
        if not pod_details:
            overall_status = "empty"
        elif has_issues:
            overall_status = "critical"
        else:
            overall_status = "healthy"
        
        return jsonify({
            'success': True,
            'namespace': namespace,
            'overall_status': overall_status,
            'pod_count': len(pod_details),
            'pods': pod_details
        }), 200
        
    except Exception as e:
        logger.error(f"Error diagnosing namespace {namespace}: {str(e)}")
        return jsonify({
            'success': False,
            'error': f"Error diagnosing namespace: {str(e)}"
        }), 500


@diagnostic_bp.route('/restart-pod/<namespace>/<pod_name>', methods=['POST'])
def restart_pod(namespace, pod_name):
    """
    Restart a specific pod by deleting it (Kubernetes will recreate it).
    
    ---
    tags:
      - Diagnostic
    summary: Restart a pod
    description: Restarts a pod by deleting it. Kubernetes will automatically recreate it based on deployment configuration.
    parameters:
      - name: namespace
        in: path
        type: string
        required: true
        description: The Kubernetes namespace name
        example: scaniamx
      - name: pod_name
        in: path
        type: string
        required: true
        description: The pod name to restart
        example: listener-66ccf94b56-cwvfw
    responses:
      200:
        description: Pod restart initiated successfully
        schema:
          type: object
          properties:
            success:
              type: boolean
              example: true
            message:
              type: string
              example: Pod listener-66ccf94b56-cwvfw restart initiated
            namespace:
              type: string
              example: scaniamx
            pod_name:
              type: string
              example: listener-66ccf94b56-cwvfw
      400:
        description: Invalid namespace or pod name
      404:
        description: Pod not found
      500:
        description: Error restarting pod
    """
    try:
        if not namespace or not pod_name:
            return jsonify({
                'success': False,
                'error': "Namespace and pod name are required"
            }), 400
        
        config.load_kube_config()
        v1 = k8s_client.CoreV1Api()
        
        try:
            # Check if pod exists
            v1.read_namespaced_pod(pod_name, namespace)
        except k8s_client.rest.ApiException as e:
            if e.status == 404:
                return jsonify({
                    'success': False,
                    'error': f"Pod '{pod_name}' not found in namespace '{namespace}'"
                }), 404
            raise
        
        # Delete the pod (Kubernetes will recreate it)
        try:
            v1.delete_namespaced_pod(
                name=pod_name,
                namespace=namespace,
                body=k8s_client.V1DeleteOptions(grace_period_seconds=30)
            )
            
            logger.info(f"Pod {pod_name} in namespace {namespace} marked for deletion")
            
            return jsonify({
                'success': True,
                'message': f"Pod {pod_name} restart initiated",
                'namespace': namespace,
                'pod_name': pod_name,
                'note': 'Pod is being terminated and will be recreated by Kubernetes'
            }), 200
            
        except k8s_client.rest.ApiException as e:
            logger.error(f"Error deleting pod {pod_name}: {e}")
            return jsonify({
                'success': False,
                'error': f"Error restarting pod: {str(e)}"
            }), 500
        
    except Exception as e:
        logger.error(f"Error in restart_pod: {str(e)}")
        return jsonify({
            'success': False,
            'error': f"Error restarting pod: {str(e)}"
        }), 500


@diagnostic_bp.route('/restart-deployment/<namespace>/<deployment_name>', methods=['POST'])
def restart_deployment(namespace, deployment_name):
    """
    Restart all pods in a deployment by rolling restart.
    
    ---
    tags:
      - Diagnostic
    summary: Restart a deployment
    description: Restarts all pods in a deployment using kubectl rollout restart
    parameters:
      - name: namespace
        in: path
        type: string
        required: true
        description: The Kubernetes namespace name
        example: scaniamx
      - name: deployment_name
        in: path
        type: string
        required: true
        description: The deployment name to restart
        example: listener
    responses:
      200:
        description: Deployment restart initiated successfully
        schema:
          type: object
          properties:
            success:
              type: boolean
              example: true
            message:
              type: string
              example: Deployment listener restart initiated
            namespace:
              type: string
              example: scaniamx
            deployment_name:
              type: string
              example: listener
      400:
        description: Invalid namespace or deployment name
      500:
        description: Error restarting deployment
    """
    try:
        if not namespace or not deployment_name:
            return jsonify({
                'success': False,
                'error': "Namespace and deployment name are required"
            }), 400
        
        # Use kubectl rollout restart command
        cmd = ['kubectl', 'rollout', 'restart', f'deployment/{deployment_name}', '-n', namespace]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode == 0:
            logger.info(f"Deployment {deployment_name} in namespace {namespace} restart initiated")
            return jsonify({
                'success': True,
                'message': f"Deployment {deployment_name} restart initiated",
                'namespace': namespace,
                'deployment_name': deployment_name,
                'note': 'Deployment is performing a rolling restart of all pods'
            }), 200
        else:
            error_msg = result.stderr or result.stdout
            logger.error(f"Error restarting deployment {deployment_name}: {error_msg}")
            return jsonify({
                'success': False,
                'error': f"Error restarting deployment: {error_msg}"
            }), 500
        
    except Exception as e:
        logger.error(f"Error in restart_deployment: {str(e)}")
        return jsonify({
            'success': False,
            'error': f"Error restarting deployment: {str(e)}"
        }), 500


def _get_container_state(container_status):
    """Extract the current state of a container."""
    if container_status.state.running:
        return "Running"
    elif container_status.state.waiting:
        return f"Waiting: {container_status.state.waiting.reason}"
    elif container_status.state.terminated:
        return f"Terminated: {container_status.state.terminated.reason}"
    else:
        return "Unknown"


def _calculate_age(creation_timestamp):
    """Calculate human-readable age from creation timestamp."""
    from datetime import datetime, timezone, timedelta
    
    now = datetime.now(timezone.utc)
    created = creation_timestamp.replace(tzinfo=timezone.utc)
    age = now - created
    
    days = age.days
    hours, remainder = divmod(age.seconds, 3600)
    
    if days > 0:
        return f"{days}d"
    elif hours > 0:
        return f"{hours}h"
    else:
        minutes = remainder // 60
        return f"{minutes}m"


@diagnostic_bp.route('/health/<namespace>', methods=['GET'])
def health_check(namespace):
    """
    Simple health check for a namespace - returns OK or NOT OK.
    
    ---
    tags:
      - Diagnostic
    summary: Quick health check for namespace
    description: Simple endpoint that returns OK if all pods are running, NOT OK if any pod is not running
    parameters:
      - name: namespace
        in: path
        type: string
        required: true
        description: The Kubernetes namespace name
        example: benjaminluna
    responses:
      200:
        description: Health check completed
        schema:
          type: object
          properties:
            status:
              type: string
              enum: [OK, NOT_OK]
              example: OK
            namespace:
              type: string
              example: benjaminluna
            healthy_pods:
              type: integer
              example: 4
            total_pods:
              type: integer
              example: 4
            unhealthy_pods:
              type: array
              items:
                type: string
              example: []
      404:
        description: Namespace not found
      500:
        description: Error checking health
    """
    try:
        try:
            config.load_kube_config()
        except Exception as kube_error:
            logger.error(f"Failed to load kube config: {str(kube_error)}")
            config.load_incluster_config()
        
        v1 = k8s_client.CoreV1Api()
        
        try:
            v1.read_namespace(namespace)
        except k8s_client.rest.ApiException as e:
            if e.status == 404:
                return jsonify({
                    'status': 'NOT_OK',
                    'error': f"Namespace '{namespace}' not found"
                }), 404
            raise
        
        # Get all pods in namespace
        pods = v1.list_namespaced_pod(namespace)
        
        total_pods = len(pods.items)
        healthy_pods = 0
        unhealthy_pods = []
        
        for pod in pods.items:
            pod_name = pod.metadata.name
            if pod.status.phase == 'Running':
                # Check if all containers are ready
                if pod.status.container_statuses:
                    all_ready = all(c.ready for c in pod.status.container_statuses)
                    if all_ready:
                        healthy_pods += 1
                    else:
                        unhealthy_pods.append(pod_name)
                else:
                    healthy_pods += 1
            else:
                unhealthy_pods.append(pod_name)
        
        status = 'OK' if len(unhealthy_pods) == 0 else 'NOT_OK'
        
        return jsonify({
            'status': status,
            'namespace': namespace,
            'healthy_pods': healthy_pods,
            'total_pods': total_pods,
            'unhealthy_pods': unhealthy_pods
        }), 200
        
    except Exception as e:
        logger.error(f"Error checking namespace health {namespace}: {str(e)}")
        return jsonify({
            'status': 'NOT_OK',
            'error': f"Error checking health: {str(e)}"
        }), 500


@diagnostic_bp.route('/restart-namespace/<namespace>', methods=['POST'])
def restart_namespace(namespace):
    """
    Restart all deployments in a namespace.
    
    ---
    tags:
      - Diagnostic
    summary: Restart entire namespace
    description: Restarts all deployments in a namespace with rolling updates (zero downtime)
    parameters:
      - name: namespace
        in: path
        type: string
        required: true
        description: The Kubernetes namespace name
        example: benjaminluna
    responses:
      200:
        description: Namespace restart initiated
        schema:
          type: object
          properties:
            success:
              type: boolean
              example: true
            message:
              type: string
              example: All deployments in namespace restarted
            namespace:
              type: string
              example: benjaminluna
            deployments_restarted:
              type: integer
              example: 4
      404:
        description: Namespace not found
      500:
        description: Error restarting namespace
    """
    try:
        try:
            config.load_kube_config()
        except Exception as kube_error:
            logger.error(f"Failed to load kube config: {str(kube_error)}")
            config.load_incluster_config()
        
        v1 = k8s_client.CoreV1Api()
        apps_v1 = k8s_client.AppsV1Api()
        
        try:
            v1.read_namespace(namespace)
        except k8s_client.rest.ApiException as e:
            if e.status == 404:
                return jsonify({
                    'success': False,
                    'error': f"Namespace '{namespace}' not found"
                }), 404
            raise
        
        # Get all deployments in namespace
        deployments = apps_v1.list_namespaced_deployment(namespace)
        
        restarted_count = 0
        failed_deployments = []
        
        for deployment in deployments.items:
            deployment_name = deployment.metadata.name
            
            try:
                # Trigger rolling restart using kubectl patch
                cmd = [
                    'kubectl', 'rollout', 'restart', 
                    f'deployment/{deployment_name}',
                    '-n', namespace
                ]
                
                result = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
                
                if result.returncode == 0:
                    logger.info(f"Deployment {deployment_name} in namespace {namespace} restarted")
                    restarted_count += 1
                else:
                    error_msg = result.stderr or result.stdout
                    logger.error(f"Error restarting deployment {deployment_name}: {error_msg}")
                    failed_deployments.append(deployment_name)
                    
            except Exception as e:
                logger.error(f"Exception restarting deployment {deployment_name}: {str(e)}")
                failed_deployments.append(deployment_name)
        
        return jsonify({
            'success': True,
            'message': f"All deployments in namespace {namespace} restarted",
            'namespace': namespace,
            'deployments_restarted': restarted_count,
            'deployments_failed': failed_deployments,
            'note': 'Deployments are performing rolling restarts of all pods'
        }), 200
        
    except Exception as e:
        logger.error(f"Error in restart_namespace: {str(e)}")
        return jsonify({
            'success': False,
            'error': f"Error restarting namespace: {str(e)}"
        }), 500
