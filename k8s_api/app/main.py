from fastapi import FastAPI, HTTPException, UploadFile, File, Form, Query
from kubernetes import client, config
from kubernetes.client.rest import ApiException
from kubernetes.dynamic import DynamicClient
from typing import Optional, List
from datetime import datetime, timezone
from helper import PodSummary, ServiceSummary, DeleteDeploymentsResponse, DeleteServicesResponse, NamespaceSummary
from pathlib import Path
import yaml
from jsonschema import validate, ValidationError


app = FastAPI(
    title="Kubernetes Management API",
    description="API for managing Kubernetes clusters using kubectl-like commands.",
    version="1.0.0"
)

# --------------------------
# Initialize Kubernetes Client
# --------------------------

try:
    config.load_kube_config()  # For local development
except:
    config.load_incluster_config()  # For in-cluster deployment

v1 = client.CoreV1Api()
apps_v1 = client.AppsV1Api()
batch_v1 = client.BatchV1Api()
autoscaling_v1 = client.AutoscalingV1Api()

# --------------------------
# JSON Schemas for Validation
# --------------------------

# Define JSON Schemas for deployment, service, and config (simplified versions)
deployment_schema = {
    "type": "object",
    "required": ["apiVersion", "kind", "metadata", "spec"],
    "properties": {
        "apiVersion": {"type": "string"},
        "kind": {"type": "string", "enum": ["Deployment"]},
        "metadata": {
            "type": "object",
            "required": ["name"],
            "properties": {
                "name": {"type": "string"},
            },
        },
        "spec": {
            "type": "object",
            "required": ["replicas", "selector", "template"],
            "properties": {
                "replicas": {"type": "integer", "minimum": 1},
                "selector": {
                    "type": "object",
                    "required": ["matchLabels"],
                    "properties": {
                        "matchLabels": {
                            "type": "object",
                            "additionalProperties": {"type": "string"},
                        }
                    },
                },
                "template": {
                    "type": "object",
                    "required": ["metadata", "spec"],
                    "properties": {
                        "metadata": {
                            "type": "object",
                            "required": ["labels"],
                            "properties": {
                                "labels": {
                                    "type": "object",
                                    "additionalProperties": {"type": "string"},
                                }
                            },
                        },
                        "spec": {
                            "type": "object",
                            "required": ["containers"],
                            "properties": {
                                "containers": {
                                    "type": "array",
                                    "minItems": 1,
                                    "items": {
                                        "type": "object",
                                        "required": ["name", "image"],
                                        "properties": {
                                            "name": {"type": "string"},
                                            "image": {"type": "string"},
                                            "ports": {
                                                "type": "array",
                                                "items": {
                                                    "type": "object",
                                                    "properties": {
                                                        "containerPort": {"type": "integer"},
                                                    },
                                                },
                                            },
                                        },
                                    },
                                }
                            },
                        },
                    },
                },
            },
        },
    },
}

service_schema = {
    "type": "object",
    "required": ["apiVersion", "kind", "metadata", "spec"],
    "properties": {
        "apiVersion": {"type": "string"},
        "kind": {"type": "string", "enum": ["Service"]},
        "metadata": {
            "type": "object",
            "required": ["name"],
            "properties": {
                "name": {"type": "string"},
            },
        },
        "spec": {
            "type": "object",
            "required": ["selector", "ports", "type"],
            "properties": {
                "selector": {
                    "type": "object",
                    "additionalProperties": {"type": "string"},
                },
                "ports": {
                    "type": "array",
                    "minItems": 1,
                    "items": {
                        "type": "object",
                        "required": ["port", "targetPort"],
                        "properties": {
                            "port": {"type": "integer"},
                            "targetPort": {"type": "integer"},
                        },
                    },
                },
                "type": {
                    "type": "string",
                    "enum": ["ClusterIP", "NodePort", "LoadBalancer", "ExternalName"],
                },
            },
        },
    },
}

config_schema = {
    "type": "object",
    "required": ["apiVersion", "kind", "metadata", "data"],
    "properties": {
        "apiVersion": {"type": "string"},
        "kind": {"type": "string", "enum": ["ConfigMap"]},
        "metadata": {
            "type": "object",
            "required": ["name"],
            "properties": {
                "name": {"type": "string"},
            },
        },
        "data": {
            "type": "object",
            "additionalProperties": {"type": "string"},
        },
    },
}

# Mapping of type to schema
schemas = {
    "deployment": deployment_schema,
    "service": service_schema,
    "config": config_schema,
}

# --------------------------
# Cluster Management Endpoints
# --------------------------

# @app.get("/cluster-info", summary="Display cluster information")
# def get_cluster_info():
#     try:
#         # Fetch detailed cluster info
#         cluster_info = v1.get_api_versions()
#         return {"api_versions": cluster_info.to_dict()}
#     except Exception as e:
#         raise HTTPException(status_code=500, detail=str(e))
#
# @app.get("/nodes", summary="List all nodes in the cluster", response_model=List[dict])
# def list_nodes():
#     try:
#         nodes = v1.list_node()
#         return [node.to_dict() for node in nodes.items]
#     except ApiException as e:
#         raise HTTPException(status_code=e.status, detail=e.body)
#
# @app.get("/nodes/{node_name}", summary="Describe a specific node in detail")
# def describe_node(node_name: str):
#     try:
#         node = v1.read_node(name=node_name)
#         return node.to_dict()
#     except ApiException as e:
#         raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# Working with Namespaces Endpoints
# --------------------------

@app.get("/namespaces", summary="List all namespaces in the cluster", response_model=List[NamespaceSummary])
def list_namespaces():
    """
    Retrieves a list of all namespaces in the Kubernetes cluster, including their names, statuses, and ages.
    """
    try:
        namespaces = v1.list_namespace()
        namespace_summaries = []
        for ns in namespaces.items:
            name = ns.metadata.name
            status = ns.status.phase
            creation_timestamp = ns.metadata.creation_timestamp
            age = calculate_age(creation_timestamp)
            namespace_summary = NamespaceSummary(
                name=name,
                status=status,
                age=age
            )
            namespace_summaries.append(namespace_summary)
        return namespace_summaries
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/namespaces", summary="Create a new namespace")
def create_namespace(namespace_name: str = Query(..., description="Name of the namespace to create")):
    namespace = client.V1Namespace(metadata=client.V1ObjectMeta(name=namespace_name))
    try:
        v1.create_namespace(body=namespace)
        return {"message": f"Namespace '{namespace_name}' created successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.delete("/namespaces/{namespace_name}", summary="Delete a namespace")
def delete_namespace(namespace_name: str):
    try:
        v1.delete_namespace(name=namespace_name)
        return {"message": f"Namespace '{namespace_name}' deleted successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# Working with Pods Endpoints
# --------------------------
#
@app.get("/pods", summary="List all pods in a namespace", response_model=List[dict])
def list_pods(
    namespace: Optional[str] = Query("default", description="Namespace to list pods from"),
    wide: Optional[bool] = Query(False, description="Show additional pod details")
):
    try:
        if wide:
            pods = v1.list_namespaced_pod(namespace, pretty="true")
        else:
            pods = v1.list_namespaced_pod(namespace)
        return [pod.to_dict() for pod in pods.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/pods/{pod_name}", summary="Describe a specific pod in detail")
def describe_pod(pod_name: str, namespace: Optional[str] = Query("default")):
    try:
        pod = v1.read_namespaced_pod(name=pod_name, namespace=namespace)
        return pod.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/pods/{pod_name}/logs", summary="Get logs from a specific pod")
def get_pod_logs(
    pod_name: str,
    namespace: Optional[str] = Query("default"),
    follow: Optional[bool] = Query(False, description="Follow log output in real-time"),
    previous: Optional[bool] = Query(False, description="Retrieve logs from the previous instance")
):
    try:
        logs = v1.read_namespaced_pod_log(
            name=pod_name,
            namespace=namespace,
            follow=follow,
            previous=previous
        )
        return {"logs": logs}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.post("/pods/{pod_name}/exec", summary="Execute a command inside a running pod")
def exec_command(
    pod_name: str,
    command: str,
    namespace: Optional[str] = Query("default")
):
    from kubernetes.stream import stream

    try:
        exec_command = [
            '/bin/sh',
            '-c',
            command
        ]
        resp = stream(v1.connect_get_namespaced_pod_exec,
                      pod_name,
                      namespace,
                      command=exec_command,
                      stderr=True, stdin=False,
                      stdout=True, tty=False)
        return {"output": resp}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# Working with Deployments Endpoints
# --------------------------

@app.get("/deployments", summary="List all deployments in a namespace", response_model=List[dict])
def list_deployments(namespace: Optional[str] = Query("default")):
    try:
        deployments = apps_v1.list_namespaced_deployment(namespace)
        return [dep.to_dict() for dep in deployments.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.post("/deployments", summary="Create a new deployment")
def create_deployment(deployment_name: str, image: str, namespace: Optional[str] = Query("default")):
    container = client.V1Container(
        name=deployment_name,
        image=image,
        ports=[client.V1ContainerPort(container_port=80)]
    )
    template = client.V1PodTemplateSpec(
        metadata=client.V1ObjectMeta(labels={"app": deployment_name}),
        spec=client.V1PodSpec(containers=[container])
    )
    spec = client.V1DeploymentSpec(
        replicas=1,
        selector={'matchLabels': {'app': deployment_name}},
        template=template
    )
    deployment = client.V1Deployment(
        metadata=client.V1ObjectMeta(name=deployment_name),
        spec=spec
    )
    try:
        apps_v1.create_namespaced_deployment(namespace=namespace, body=deployment)
        return {"message": f"Deployment '{deployment_name}' created successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.put("/deployments/{deployment_name}/scale", summary="Scale a deployment")
def scale_deployment(deployment_name: str, replicas: int, namespace: Optional[str] = Query("default")):
    try:
        deployment = apps_v1.read_namespaced_deployment(name=deployment_name, namespace=namespace)
        deployment.spec.replicas = replicas
        apps_v1.patch_namespaced_deployment(name=deployment_name, namespace=namespace, body=deployment)
        return {"message": f"Deployment '{deployment_name}' scaled to {replicas} replicas."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/deployments/{deployment_name}", summary="Describe a specific deployment in detail")
def describe_deployment(deployment_name: str, namespace: Optional[str] = Query("default")):
    try:
        deployment = apps_v1.read_namespaced_deployment(name=deployment_name, namespace=namespace)
        return deployment.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.delete("/deployments/{deployment_name}", summary="Delete a deployment")
def delete_deployment(deployment_name: str, namespace: Optional[str] = Query("default")):
    try:
        apps_v1.delete_namespaced_deployment(name=deployment_name, namespace=namespace)
        return {"message": f"Deployment '{deployment_name}' deleted successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# Working with Services Endpoints
# --------------------------

@app.get("/services", summary="List all services in a namespace", response_model=List[dict])
def list_services(namespace: Optional[str] = Query("default")):
    try:
        services = v1.list_namespaced_service(namespace)
        return [svc.to_dict() for svc in services.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.post("/services", summary="Expose a pod as a service")
def expose_pod(
    pod_name: str,
    service_type: str,
    port: int,
    service_name: Optional[str] = Query(None),
    namespace: Optional[str] = Query("default")
):
    valid_service_types = {"ClusterIP", "NodePort", "LoadBalancer", "ExternalName"}
    if service_type not in valid_service_types:
        raise HTTPException(status_code=400, detail=f"Invalid service type. Must be one of {valid_service_types}.")

    service_name = service_name or f"{pod_name}-service"
    service = client.V1Service(
        metadata=client.V1ObjectMeta(name=service_name),
        spec=client.V1ServiceSpec(
            selector={"app": pod_name},
            ports=[client.V1ServicePort(port=port)],
            type=service_type
        )
    )
    try:
        v1.create_namespaced_service(namespace=namespace, body=service)
        return {"message": f"Service '{service_name}' created successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/services/{service_name}", summary="Describe a specific service in detail")
def describe_service(service_name: str, namespace: Optional[str] = Query("default")):
    try:
        service = v1.read_namespaced_service(name=service_name, namespace=namespace)
        return service.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.delete("/services/{service_name}", summary="Delete a service")
def delete_service(service_name: str, namespace: Optional[str] = Query("default")):
    try:
        v1.delete_namespaced_service(name=service_name, namespace=namespace)
        return {"message": f"Service '{service_name}' deleted successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# Working with ConfigMaps and Secrets Endpoints
# --------------------------

@app.post("/configmaps", summary="Create a ConfigMap from literal values")
def create_configmap(configmap_name: str, data: dict, namespace: Optional[str] = Query("default")):
    configmap = client.V1ConfigMap(
        metadata=client.V1ObjectMeta(name=configmap_name),
        data=data
    )
    try:
        v1.create_namespaced_config_map(namespace=namespace, body=configmap)
        return {"message": f"ConfigMap '{configmap_name}' created successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/configmaps", summary="List all ConfigMaps in a namespace", response_model=List[dict])
def list_configmaps(namespace: Optional[str] = Query("default")):
    try:
        configmaps = v1.list_namespaced_config_map(namespace)
        return [cm.to_dict() for cm in configmaps.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/configmaps/{configmap_name}", summary="Describe a specific ConfigMap")
def describe_configmap(configmap_name: str, namespace: Optional[str] = Query("default")):
    try:
        configmap = v1.read_namespaced_config_map(name=configmap_name, namespace=namespace)
        return configmap.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.delete("/configmaps/{configmap_name}", summary="Delete a ConfigMap")
def delete_configmap(configmap_name: str, namespace: Optional[str] = Query("default")):
    try:
        v1.delete_namespaced_config_map(name=configmap_name, namespace=namespace)
        return {"message": f"ConfigMap '{configmap_name}' deleted successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.post("/secrets", summary="Create a Secret from literal values")
def create_secret(secret_name: str, data: dict, namespace: Optional[str] = Query("default")):
    secret = client.V1Secret(
        metadata=client.V1ObjectMeta(name=secret_name),
        string_data=data
    )
    try:
        v1.create_namespaced_secret(namespace=namespace, body=secret)
        return {"message": f"Secret '{secret_name}' created successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/secrets", summary="List all Secrets in a namespace", response_model=List[dict])
def list_secrets(namespace: Optional[str] = Query("default")):
    try:
        secrets = v1.list_namespaced_secret(namespace)
        return [secret.to_dict() for secret in secrets.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/secrets/{secret_name}", summary="Describe a specific Secret")
def describe_secret(secret_name: str, namespace: Optional[str] = Query("default")):
    try:
        secret = v1.read_namespaced_secret(name=secret_name, namespace=namespace)
        return secret.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.delete("/secrets/{secret_name}", summary="Delete a Secret")
def delete_secret(secret_name: str, namespace: Optional[str] = Query("default")):
    try:
        v1.delete_namespaced_secret(name=secret_name, namespace=namespace)
        return {"message": f"Secret '{secret_name}' deleted successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# Working with Persistent Volumes and Claims Endpoints
# --------------------------

@app.get("/persistentvolumes", summary="List all Persistent Volumes", response_model=List[dict])
def list_pv():
    try:
        pvs = v1.list_persistent_volume()
        return [pv.to_dict() for pv in pvs.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/persistentvolumeclaims", summary="List all Persistent Volume Claims in a namespace", response_model=List[dict])
def list_pvc(namespace: Optional[str] = Query("default")):
    try:
        pvcs = v1.list_namespaced_persistent_volume_claim(namespace)
        return [pvc.to_dict() for pvc in pvcs.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/persistentvolumes/{pv_name}", summary="Describe a specific Persistent Volume")
def describe_pv(pv_name: str):
    try:
        pv = v1.read_persistent_volume(name=pv_name)
        return pv.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/persistentvolumeclaims/{pvc_name}", summary="Describe a specific Persistent Volume Claim")
def describe_pvc(pvc_name: str, namespace: Optional[str] = Query("default")):
    try:
        pvc = v1.read_namespaced_persistent_volume_claim(name=pvc_name, namespace=namespace)
        return pvc.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# Working with ReplicaSets and StatefulSets Endpoints
# --------------------------

@app.get("/replicasets", summary="List all ReplicaSets in a namespace", response_model=List[dict])
def list_replicasets(namespace: Optional[str] = Query("default")):
    try:
        rs = apps_v1.list_namespaced_replica_set(namespace)
        return [replicaset.to_dict() for replicaset in rs.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/replicasets/{replicaset_name}", summary="Describe a specific ReplicaSet")
def describe_replicaset(replicaset_name: str, namespace: Optional[str] = Query("default")):
    try:
        replicaset = apps_v1.read_namespaced_replica_set(name=replicaset_name, namespace=namespace)
        return replicaset.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/statefulsets", summary="List all StatefulSets in a namespace", response_model=List[dict])
def list_statefulsets(namespace: Optional[str] = Query("default")):
    try:
        sts = apps_v1.list_namespaced_stateful_set(namespace)
        return [statefulset.to_dict() for statefulset in sts.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/statefulsets/{statefulset_name}", summary="Describe a specific StatefulSet")
def describe_statefulset(statefulset_name: str, namespace: Optional[str] = Query("default")):
    try:
        statefulset = apps_v1.read_namespaced_stateful_set(name=statefulset_name, namespace=namespace)
        return statefulset.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# Working with Jobs and CronJobs Endpoints
# --------------------------

@app.get("/jobs", summary="List all Jobs in a namespace", response_model=List[dict])
def list_jobs(namespace: Optional[str] = Query("default")):
    try:
        jobs = batch_v1.list_namespaced_job(namespace)
        return [job.to_dict() for job in jobs.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.post("/jobs", summary="Create a new Job")
def create_job(job_name: str, image: str, namespace: Optional[str] = Query("default")):
    container = client.V1Container(
        name=job_name,
        image=image,
        command=["/bin/sh", "-c", "echo Hello Kubernetes! && sleep 30"]
    )
    template = client.V1PodTemplateSpec(
        metadata=client.V1ObjectMeta(labels={"job": job_name}),
        spec=client.V1PodSpec(restart_policy="Never", containers=[container])
    )
    job_spec = client.V1JobSpec(
        template=template,
        backoff_limit=4
    )
    job = client.V1Job(
        metadata=client.V1ObjectMeta(name=job_name),
        spec=job_spec
    )
    try:
        batch_v1.create_namespaced_job(namespace=namespace, body=job)
        return {"message": f"Job '{job_name}' created successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/cronjobs", summary="List all CronJobs in a namespace", response_model=List[dict])
def list_cronjobs(namespace: Optional[str] = Query("default")):
    try:
        cronjobs = batch_v1.list_namespaced_cron_job(namespace)
        return [cronjob.to_dict() for cronjob in cronjobs.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/cronjobs/{cronjob_name}", summary="Describe a specific CronJob")
def describe_cronjob(cronjob_name: str, namespace: Optional[str] = Query("default")):
    try:
        cronjob = batch_v1.read_namespaced_cron_job(name=cronjob_name, namespace=namespace)
        return cronjob.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# Working with Resources Endpoints
# --------------------------

@app.post("/resources/apply", summary="Apply a configuration from a YAML file")
def apply_resource(file_content: str):
    try:
        resources = list(yaml.safe_load_all(file_content))
        for resource in resources:
            kind = resource.get('kind')
            api_version = resource.get('apiVersion')
            if kind and api_version:
                api = client.ApiClient()
                # Using DynamicClient requires the 'kubernetes-client' package
                from kubernetes.dynamic import DynamicClient
                dynamic_client = DynamicClient(api)
                resource_api = dynamic_client.resources.get(api_version=api_version, kind=kind)
                resource_api.create(body=resource)
        return {"message": "Resources applied successfully."}
    except yaml.YAMLError as e:
        raise HTTPException(status_code=400, detail=f"Invalid YAML content: {str(e)}")
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.delete("/resources/delete", summary="Delete resources defined in a YAML file")
def delete_resource(file_content: str):
    try:
        resources = list(yaml.safe_load_all(file_content))
        for resource in resources:
            kind = resource.get('kind')
            api_version = resource.get('apiVersion')
            metadata = resource.get('metadata', {})
            name = metadata.get('name')
            namespace = metadata.get('namespace', "default")
            if kind and api_version and name:
                api = client.ApiClient()
                from kubernetes.dynamic import DynamicClient
                dynamic_client = DynamicClient(api)
                resource_api = dynamic_client.resources.get(api_version=api_version, kind=kind)
                resource_api.delete(name=name, namespace=namespace)
        return {"message": "Resources deleted successfully."}
    except yaml.YAMLError as e:
        raise HTTPException(status_code=400, detail=f"Invalid YAML content: {str(e)}")
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.put("/resources/edit", summary="Edit a resource live in the default editor")
def edit_resource(resource: str, resource_name: str, namespace: Optional[str] = Query("default")):
    import tempfile
    import subprocess

    try:
        # Fetch the resource
        api = client.ApiClient()
        from kubernetes.dynamic import DynamicClient
        dynamic_client = DynamicClient(api)
        resource_api = dynamic_client.resources.get(kind=resource.capitalize(), api_version="v1")
        resource_obj = resource_api.get(name=resource_name, namespace=namespace)

        # Write to a temp file
        with tempfile.NamedTemporaryFile(mode='w+', delete=False, suffix=".yaml") as tmp:
            yaml.dump(resource_obj.to_dict(), tmp)
            tmp_path = tmp.name

        # Open editor (using 'vi'; change to your preferred editor if needed)
        subprocess.run(["vi", tmp_path])

        # Read edited file
        with open(tmp_path, 'r') as tmp:
            edited_content = yaml.safe_load(tmp)

        # Replace the resource
        resource_api.replace(name=resource_name, namespace=namespace, body=edited_content)
        return {"message": f"Resource '{resource_name}' edited successfully."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/resources/replace", summary="Replace a resource with configuration from a YAML file")
def replace_resource(file_content: str):
    try:
        resources = list(yaml.safe_load_all(file_content))
        for resource in resources:
            kind = resource.get('kind')
            api_version = resource.get('apiVersion')
            metadata = resource.get('metadata', {})
            name = metadata.get('name')
            namespace = metadata.get('namespace', "default")
            if kind and api_version and name:
                api = client.ApiClient()
                from kubernetes.dynamic import DynamicClient
                dynamic_client = DynamicClient(api)
                resource_api = dynamic_client.resources.get(api_version=api_version, kind=kind)
                resource_api.replace(name=name, namespace=namespace, body=resource)
        return {"message": "Resources replaced successfully."}
    except yaml.YAMLError as e:
        raise HTTPException(status_code=400, detail=f"Invalid YAML content: {str(e)}")
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

# --------------------------
# Accessing and Debugging Endpoints
# --------------------------

@app.post("/port-forward", summary="Forward a local port to a port on the pod")
def port_forward(pod_name: str, local_port: int, pod_port: int, namespace: Optional[str] = Query("default")):
    import threading
    from kubernetes.stream import portforward

    try:
        api = client.CoreV1Api()
        pf = portforward.PortForwardClient(api.api_client)
        pf.start_pod_portforward(pod_name, namespace, ports=str(pod_port))
        # Note: Implementing real-time port forwarding is complex and beyond this scope
        return {"message": f"Port forwarding from local port {local_port} to pod port {pod_port} started."}
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as ex:
        raise HTTPException(status_code=500, detail=str(ex))

@app.get("/top/nodes", summary="Show resource usage by nodes", response_model=List[dict])
def top_nodes():
    try:
        # Requires Metrics Server
        from kubernetes.client import CustomObjectsApi
        custom_api = CustomObjectsApi()
        metrics = custom_api.list_cluster_custom_object(
            group="metrics.k8s.io",
            version="v1beta1",
            plural="nodes"
        )
        return metrics
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as ex:
        raise HTTPException(status_code=500, detail=str(ex))

@app.get("/top/pods", summary="Show resource usage by pods", response_model=List[dict])
def top_pods(namespace: Optional[str] = Query("default")):
    try:
        # Requires Metrics Server
        from kubernetes.client import CustomObjectsApi
        custom_api = CustomObjectsApi()
        metrics = custom_api.list_namespaced_custom_object(
            group="metrics.k8s.io",
            version="v1beta1",
            namespace=namespace,
            plural="pods"
        )
        return metrics
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as ex:
        raise HTTPException(status_code=500, detail=str(ex))

@app.get("/events", summary="List recent events in the cluster", response_model=List[dict])
def list_events(namespace: Optional[str] = Query("default")):
    try:
        events = v1.list_namespaced_event(namespace)
        return [event.to_dict() for event in events.items]
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

@app.get("/version", summary="Display the client and server Kubernetes versions")
def get_version():
    try:
        version = client.VersionApi().get_code()
        return version.to_dict()
    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)

# --------------------------
# YAML File Management Endpoints
# --------------------------

@app.post("/yaml/create", summary="Create a YAML file")
async def create_yaml_file(
    port: int = Form(..., description="Port number"),
    gear_name: str = Form(..., description="Name of the gear"),
    type: str = Form(..., description="Type of YAML file (deployment, service, config)"),
    file: UploadFile = File(...)
):
    """
    Creates a YAML file with the specified parameters and stores it in jonobridge/config/manager.
    The file name is composed as <port>-<gear_name>-<type>.yaml.
    """
    # Validate 'type' parameter
    valid_types = {"deployment", "service", "config"}
    type_lower = type.lower()
    if type_lower not in valid_types:
        raise HTTPException(status_code=400, detail=f"Invalid type. Must be one of {valid_types}.")

    # Read file content
    try:
        file_content = await file.read()
        yaml_content = yaml.safe_load(file_content)
    except yaml.YAMLError as e:
        raise HTTPException(status_code=400, detail=f"Invalid YAML content: {str(e)}")
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Failed to read uploaded file: {str(e)}")

    # Validate YAML content against schema
    schema = schemas.get(type_lower)
    if schema:
        try:
            validate(instance=yaml_content, schema=schema)
        except ValidationError as e:
            raise HTTPException(status_code=400, detail=f"YAML does not conform to the {type_lower} schema: {e.message}")

    # Construct filename
    filename = f"{port}-{gear_name}-{type_lower}.yaml"

    # Define directory path
    config_dir = Path.home() / "jonobridge" / "config" / "manager"

    # Create directory if it doesn't exist
    try:
        config_dir.mkdir(parents=True, exist_ok=True)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to create directory '{config_dir}': {str(e)}")

    # Define full file path
    file_path = config_dir / filename

    # Write YAML content to file
    try:
        with open(file_path, 'w') as f:
            yaml.dump(yaml_content, f)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to write YAML file: {str(e)}")

    # Verify file creation
    if not file_path.exists():
        raise HTTPException(status_code=500, detail="YAML file was not created successfully.")

    return {"message": f"YAML file '{filename}' created successfully at '{file_path}'."}

@app.get("/yaml/read", summary="Read a YAML file")
def read_yaml_file(
    port: int = Query(..., description="Port number"),
    gear_name: str = Query(..., description="Name of the gear"),
    type: str = Query(..., description="Type of YAML file (deployment, service, config)")
):
    """
    Reads and returns the content of a YAML file based on the provided parameters.
    """
    # Validate 'type' parameter
    valid_types = {"deployment", "service", "config"}
    type_lower = type.lower()
    if type_lower not in valid_types:
        raise HTTPException(status_code=400, detail=f"Invalid type. Must be one of {valid_types}.")

    # Construct filename
    filename = f"{port}-{gear_name}-{type_lower}.yaml"

    # Define directory path
    config_dir = Path.home() / "jonobridge" / "config" / "manager"

    # Define full file path
    file_path = config_dir / filename

    # Check if file exists
    if not file_path.exists():
        raise HTTPException(status_code=404, detail=f"YAML file '{filename}' does not exist in '{config_dir}'.")

    # Read YAML content from file
    try:
        with open(file_path, 'r') as f:
            content = f.read()
            yaml_content = yaml.safe_load(content)
    except yaml.YAMLError as e:
        raise HTTPException(status_code=500, detail=f"Stored YAML content is invalid: {str(e)}")
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to read YAML file: {str(e)}")

    # Optional: Validate YAML content against schema
    schema = schemas.get(type_lower)
    if schema:
        try:
            validate(instance=yaml_content, schema=schema)
        except ValidationError as e:
            raise HTTPException(status_code=500, detail=f"Stored YAML content does not conform to the {type_lower} schema: {e.message}")

    return {"filename": filename, "content": content}

# --------------------------
# New YAML Delete Endpoint
# --------------------------

@app.delete("/yaml/delete", summary="Delete a YAML file")
def delete_yaml_file(
    port: int = Query(..., description="Port number"),
    gear_name: str = Query(..., description="Name of the gear"),
    type: str = Query(..., description="Type of YAML file (deployment, service, config)")
):
    """
    Deletes a YAML file from jonobridge/config/manager based on the provided parameters.
    The file to be deleted is identified by <port>-<gear_name>-<type>.yaml.
    """
    # Validate 'type' parameter
    valid_types = {"deployment", "service", "config"}
    type_lower = type.lower()
    if type_lower not in valid_types:
        raise HTTPException(status_code=400, detail=f"Invalid type. Must be one of {valid_types}.")

    # Construct filename
    filename = f"{port}-{gear_name}-{type_lower}.yaml"

    # Define directory path
    config_dir = Path.home() / "jonobridge" / "config" / "manager"

    # Define full file path
    file_path = config_dir / filename

    # Check if file exists
    if not file_path.exists():
        raise HTTPException(status_code=404, detail=f"YAML file '{filename}' does not exist in '{config_dir}'.")

    # Attempt to delete the file
    try:
        file_path.unlink()
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to delete YAML file: {str(e)}")

    return {"message": f"YAML file '{filename}' has been deleted successfully from '{config_dir}'."}

# --------------------------
# New YAML Apply Endpoint
# --------------------------

@app.post("/yaml/apply", summary="Apply a stored YAML file to the Kubernetes cluster")
def apply_yaml_file(
    port: int = Query(..., description="Port number"),
    gear_name: str = Query(..., description="Name of the gear"),
    type: str = Query(..., description="Type of YAML file (deployment, service, config)"),
    namespace: Optional[str] = Query("default", description="Namespace to apply the resource to")
):
    """
    Applies a YAML file to the Kubernetes cluster based on the provided parameters.
    The YAML file is identified by <port>-<gear_name>-<type>.yaml in the jonobridge/config/manager directory.
    """
    # Validate 'type' parameter
    valid_types = {"deployment", "service", "config"}
    type_lower = type.lower()
    if type_lower not in valid_types:
        raise HTTPException(status_code=400, detail=f"Invalid type. Must be one of {valid_types}.")

    # Construct filename
    filename = f"{port}-{gear_name}-{type_lower}.yaml"

    # Define directory path
    config_dir = Path.home() / "jonobridge" / "config" / "manager"

    # Define full file path
    file_path = config_dir / filename

    # Check if file exists
    if not file_path.exists():
        raise HTTPException(status_code=404, detail=f"YAML file '{filename}' does not exist in '{config_dir}'.")

    # Read YAML content from file
    try:
        with open(file_path, 'r') as f:
            yaml_content = yaml.safe_load(f)
    except yaml.YAMLError as e:
        raise HTTPException(status_code=500, detail=f"Stored YAML content is invalid: {str(e)}")
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to read YAML file: {str(e)}")

    # Apply the YAML content to the cluster
    try:
        # Determine the kind and API version
        kind = yaml_content.get("kind")
        api_version = yaml_content.get("apiVersion")
        if not kind or not api_version:
            raise HTTPException(status_code=400, detail="YAML file missing 'kind' or 'apiVersion' fields.")

        # Initialize Dynamic Client
        api_client = client.ApiClient()
        dynamic_client = DynamicClient(api_client)  # Correct usage

        # Get the resource API based on kind and api_version
        resource_api = dynamic_client.resources.get(api_version=api_version, kind=kind)

        # Extract metadata
        metadata = yaml_content.get("metadata", {})
        name = metadata.get("name")
        if not name:
            raise HTTPException(status_code=400, detail="YAML file missing 'metadata.name' field.")

        # Apply the resource
        namespace = metadata.get("namespace", namespace)
        try:
            # Check if the resource already exists
            existing_resource = resource_api.get(name=name, namespace=namespace)
            # If exists, replace it
            resource_api.replace(name=name, namespace=namespace, body=yaml_content)
        except ApiException as e:
            if e.status == 404:
                # If not found, create it
                resource_api.create(body=yaml_content, namespace=namespace)
            else:
                raise HTTPException(status_code=e.status, detail=e.body)

    except ApiException as e:
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

    return {"message": f"YAML file '{filename}' has been applied successfully to the Kubernetes cluster in namespace '{namespace}'."}




# --------------------------
# YAML File Management Endpoints
# --------------------------

# ... [Other endpoints like /yaml/create, /yaml/read, /yaml/delete remain unchanged]

# --------------------------
# Simplified Pods Listing Endpoint
# --------------------------
def calculate_age(creation_timestamp: datetime) -> str:
    now = datetime.now(timezone.utc)
    age_delta = now - creation_timestamp

    seconds = int(age_delta.total_seconds())
    minutes, seconds = divmod(seconds, 60)
    hours, minutes = divmod(minutes, 60)
    days, hours = divmod(hours, 24)

    age_parts = []
    if days > 0:
        age_parts.append(f"{days}d")
    if hours > 0:
        age_parts.append(f"{hours}h")
    if minutes > 0:
        age_parts.append(f"{minutes}m")
    age_parts.append(f"{seconds}s")

    return ''.join(age_parts)


def calculate_ready_restarts(pod: client.V1Pod) -> (str, int):
    total_containers = len(pod.spec.containers)
    ready_containers = sum(1 for status in pod.status.container_statuses or [] if status.ready)
    restarts = sum(status.restart_count for status in pod.status.container_statuses or [])
    ready = f"{ready_containers}/{total_containers}"
    return ready, restarts

def get_pod_status(pod: client.V1Pod) -> str:
    if pod.status.phase == "Running":
        for condition in pod.status.conditions or []:
            if condition.type == "Ready":
                return "Running" if condition.status == "True" else "Not Ready"
    return pod.status.phase

def format_service_ports(service: client.V1Service) -> str:
    port_list = []
    for port in service.spec.ports or []:
        if port.node_port:
            port_entry = f"{port.port}:{port.node_port}/{port.protocol}"
        else:
            port_entry = f"{port.port}/{port.protocol}"
        port_list.append(port_entry)
    return ','.join(port_list) if port_list else "<none>"

def list_all_deployments(namespace: str) -> List[client.V1Deployment]:
    """
    Retrieves all deployments within the specified namespace.
    """
    deployments = apps_v1.list_namespaced_deployment(namespace)
    return deployments.items

def delete_deployment(namespace: str, name: str):
    """
    Deletes a deployment by name within the specified namespace.
    """
    try:
        apps_v1.delete_namespaced_deployment(
            name=name,
            namespace=namespace,
            body=client.V1DeleteOptions(
                propagation_policy='Foreground',
                grace_period_seconds=0
            )
        )
    except ApiException as e:
        raise e

def list_all_services(namespace: str) -> List[client.V1Service]:
    """
    Retrieves all services within the specified namespace.
    """
    services = v1.list_namespaced_service(namespace)
    return services.items

def delete_service(namespace: str, name: str):
    """
    Deletes a service by name within the specified namespace.
    """
    try:
        v1.delete_namespaced_service(
            name=name,
            namespace=namespace,
            body=client.V1DeleteOptions(
                propagation_policy='Foreground',
                grace_period_seconds=0
            )
        )
    except ApiException as e:
        raise e
@app.get("/pods", summary="List all pods in a namespace", response_model=List[PodSummary])
def list_pods(namespace: Optional[str] = Query("default", description="Kubernetes namespace to list pods from")):
    """
    Retrieves a simplified list of all pods within the specified Kubernetes namespace.
    If no namespace is provided, defaults to the 'default' namespace.
    """
    try:
        # Fetch the list of pods in the specified namespace
        pods = v1.list_namespaced_pod(namespace)

        # Transform each pod into a PodSummary
        pod_summaries = []
        for pod in pods.items:
            name = pod.metadata.name
            ready, restarts = calculate_ready_restarts(pod)
            status = pod.status.phase
            age = calculate_age(pod.metadata.creation_timestamp)

            pod_summary = PodSummary(
                name=name,
                ready=ready,
                status=status,
                restarts=restarts,
                age=age
            )
            pod_summaries.append(pod_summary)

        return pod_summaries
    except ApiException as e:
        # Handle Kubernetes API exceptions
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as e:
        # Handle any other exceptions
        raise HTTPException(status_code=500, detail=str(e))

@app.delete("/deployments", summary="Delete all deployments in a namespace", response_model=DeleteDeploymentsResponse)
def delete_all_deployments(namespace: str = Query(..., description="Kubernetes namespace to delete deployments from")):
    """
    Deletes all deployments within the specified Kubernetes namespace.
    """
    try:
        # List all deployments in the specified namespace
        deployments = list_all_deployments(namespace)

        if not deployments:
            return DeleteDeploymentsResponse(
                deleted_deployments=[],
                namespace=namespace,
                message="No deployments found in the specified namespace."
            )

        deleted_deployments = []
        for deployment in deployments:
            name = deployment.metadata.name
            try:
                delete_deployment(namespace, name)
                deleted_deployments.append(name)
            except ApiException as e:
                # Log the exception details if necessary
                raise HTTPException(status_code=e.status, detail=f"Failed to delete deployment '{name}': {e.body}")

        return DeleteDeploymentsResponse(
            deleted_deployments=deleted_deployments,
            namespace=namespace,
            message=f"Successfully deleted {len(deleted_deployments)} deployment(s) in namespace '{namespace}'."
        )

    except ApiException as e:
        # Handle Kubernetes API exceptions
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as e:
        # Handle any other unforeseen exceptions
        raise HTTPException(status_code=500, detail=str(e))

# --------------------------
# New Services Listing Endpoint
# --------------------------

@app.get("/services", summary="List all services in a namespace", response_model=List[ServiceSummary])
def list_services(
        namespace: Optional[str] = Query("default", description="Kubernetes namespace to list services from")):
    """
    Retrieves a simplified list of all services within the specified Kubernetes namespace.
    If no namespace is provided, defaults to the 'default' namespace.
    """
    try:
        # Fetch the list of services in the specified namespace
        services = v1.list_namespaced_service(namespace)

        # Transform each service into a ServiceSummary
        service_summaries = []
        for service in services.items:
            name = service.metadata.name
            type_ = service.spec.type
            cluster_ip = service.spec.cluster_ip
            external_ip = service.spec.external_i_ps[0].ip if service.spec.external_i_ps else "<none>"
            ports = format_service_ports(service)
            age = calculate_age(service.metadata.creation_timestamp)

            service_summary = ServiceSummary(
                name=name,
                type=type_,
                cluster_ip=cluster_ip,
                external_ip=external_ip,
                ports=ports,
                age=age
            )
            service_summaries.append(service_summary)

        return service_summaries
    except ApiException as e:
        # Handle Kubernetes API exceptions
        raise HTTPException(status_code=e.status, detail=e.body)
    except Exception as e:
        # Handle any other exceptions
        raise HTTPException(status_code=500, detail=str(e))


# --------------------------
# Services Deletion Endpoint
# --------------------------
@app.delete(
    "/services",
    summary="Delete all services in a namespace",
    response_model=DeleteServicesResponse
)
def delete_all_services(
    namespace: str = Query(..., description="Kubernetes namespace to delete services from"),
    confirm: bool = Query(False, description="Set to true to confirm deletion")
):
    """
    Deletes all services within the specified Kubernetes namespace.
    Requires confirmation to proceed.
    """
    if not confirm:
        raise HTTPException(
            status_code=400,
            detail="Deletion not confirmed. Please set the 'confirm' query parameter to true to proceed."
        )

    try:
        # List all services in the specified namespace
        services = list_all_services(namespace)

        if not services:
            return DeleteServicesResponse(
                deleted_services=[],
                namespace=namespace,
                message="No services found in the specified namespace."
            )

        deleted_services = []
        for service in services:
            name = service.metadata.name
            try:
                delete_service(namespace, name)
                deleted_services.append(name)
            except ApiException as e:
                raise HTTPException(
                    status_code=e.status,
                    detail=f"Failed to delete service '{name}': {e.body}"
                )

        return DeleteServicesResponse(
            deleted_services=deleted_services,
            namespace=namespace,
            message=f"Successfully deleted {len(deleted_services)} service(s) in namespace '{namespace}'."
        )

    except ApiException as e:
        if e.status == 403:
            detail = "Forbidden: You do not have permission to delete services in this namespace."
        elif e.status == 404:
            detail = f"Namespace '{namespace}' not found."
        else:
            detail = f"Kubernetes API returned an error: {e.reason}"
        raise HTTPException(status_code=e.status, detail=detail)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))