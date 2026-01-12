from typing import Dict, Any, List
from ..base_service import BaseService

class ListenerService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "listener"
    
    @property
    def service_type(self) -> str:
        """Return the name of the type of service, which matches the database table name."""
        return "input"  # This must match the database table name

    @property
    def inputs(self) -> List[str]:
        """Return the inputs of the service."""
        return [
            "",
        ]   

    @property
    def outputs(self) -> List[str]:
        """Return the outputs of the service."""
        return [
            "tracker/from-udp",
            "tracker/from-tcp"
        ] 

    @property
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "port": "int",
            "api_port": "int",
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "port": ["8000-8999"],
            "api_port": ["9000-9999"],
            "replicas": ["1"]
        }

    @property
    def kubernetes_templates(self) -> Dict[str, Any]:
        """Return the Kubernetes manifest templates for the service."""
        return {
            "namespace": {
                "apiVersion": "v1",
                "kind": "Namespace",
                "metadata": {
                    "name": None
                }
            },
            "deployment": {
                "apiVersion": "apps/v1",
                "kind": "Deployment",
                "metadata": {
                    "name": "listener",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "listener"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "listener"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "listener",
                                "image": "maddsystems/listener:1.0.0",
                                "imagePullPolicy": "Always",
                                "ports": [{
                                    "containerPort": None
                                }],
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
                                    {
                                        "name": "SQL_LOG",
                                        "value": "False"
                                    }
                                ]
                            }]
                        }
                    }
                }
            },
            "service": {
                "apiVersion": "v1",
                "kind": "Service",
                "metadata": {
                    "name": "listener-service",
                    "namespace": None
                },
                "spec": {
                    "type": "LoadBalancer",
                    "selector": {
                        "app": "listener"
                    },
                    "ports": [{
                        "name": "tcp",
                        "port": None,
                        "protocol": "TCP",
                        "targetPort": 1024,
                    },
                    {
                        "name": "udp",
                        "port": None,
                        "protocol": "UDP",
                        "targetPort": 1024,
                    },
                    {
                        "name": "tcp-api",
                        "port": None,
                        "protocol": "TCP",
                        "targetPort": 8080,
                    }]
                }
            }
        }

    def _customize_manifests(self, templates: Dict[str, Any], client_name: str, params: Dict[str, Any]) -> None:
        """Customize the Kubernetes manifests with service-specific logic."""
        print(f"Customizing manifests for {self.service_name} with params:", params)
        
        # Set namespace
        templates["namespace"]["metadata"]["name"] = client_name
        
        # Update deployment
        deployment = templates["deployment"]
        deployment["metadata"]["namespace"] = client_name
        deployment["spec"]["replicas"] = params["replicas"]
        
        # Update container
        container = deployment["spec"]["template"]["spec"]["containers"][0]
        container["ports"][0]["containerPort"] = params["port"]
        
        # Update environment variables
        # for env in container["env"]:
        #     if env["name"] == "MEITRACK_HOST":
        #         env["value"] = params["meitrack_host"]
        new_port = "1"+str(params["port"])
        # Update service
        if "service" in templates:
            service = templates["service"]
            service["metadata"]["namespace"] = client_name
            service["spec"]["ports"][0].update({
                "port": params["port"],
                "protocol": "TCP",
                "targetPort": 1024
            })
            service["spec"]["ports"][1].update({
                "port": params["port"],
                "protocol": "UDP",
                "targetPort": 1024
            })
            service["spec"]["ports"][2].update({
                "port": params["api_port"],
                "protocol": "TCP",
                "targetPort": 8080
            })
