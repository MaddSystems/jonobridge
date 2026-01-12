from typing import Dict, Any, List
from ..base_service import BaseService

class Service(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "httpinput"
    
    @property
    def service_type(self) -> str:
        """Return the name of the type of service, which matches the database table name."""
        return "input"  # This must match the database table name

    @property
    def inputs(self) -> List[str]:
        """Return the inputs of the service."""
        return []   

    @property
    def outputs(self) -> List[str]:
        """Return the outputs of the service."""
        return [
            "httpinput/get"
        ] 
    
    @property
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "portal_endpoint": "varchar(255)",
            "webrelay_port": "int",   
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "portal_endpoint": ["test"],
            "webrelay_port": ["10000-19999"],
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
                    "name": "httpinput",
                    "namespace": None
                },
                "spec": {
                    "revisionHistoryLimit": 2,  # Keep only 2 old ReplicaSets
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "httpinput"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "httpinput"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "httpinput",
                                "image": "maddsystems/httpinput:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
                                    {
                                        "name": "PORTAL_ENDPOINT",
                                        "value": None
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
                    "name": "httpinput",
                    "namespace": None
                },
                "spec": {
                    "type": "LoadBalancer",
                    "selector": {
                        "app": "httpinput"
                    },
                    "ports": [{
                        "name": "tcp-portal",
                        "port": None,
                        "protocol": "TCP",
                        "targetPort": 8081,
                    }
                ]
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

        
        # Update environment variables
        for env in container["env"]:
            if env["name"] == "PORTAL_ENDPOINT":
                env["value"] = params["portal_endpoint"]
        
        # Update service
                # Update service namespace
        if "service" in templates:
            service = templates["service"]
            templates["service"]["metadata"]["namespace"] = client_name
            service["spec"]["ports"][0].update({
                "port": params["webrelay_port"],
                "protocol": "TCP",
                "targetPort": 8081
            })