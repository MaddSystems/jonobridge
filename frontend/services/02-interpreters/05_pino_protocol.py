from typing import Dict, Any, List
from ..base_service import BaseService

class PinoProtocolService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "pinoprotocol"  # This must match the database table name
        
    @property
    def service_type(self) -> str:
        """Return the name of the type of service, which matches the database table name."""
        return "interpreter"  # This must match the database table name

    @property
    def inputs(self) -> List[str]:
        """Return the inputs of the service."""
        return [
            "tracker/from-tcp",
        ]   

    @property
    def outputs(self) -> List[str]:
        """Return the outputs of the service."""
        return [
            "tracker/jonoprotocol"
        ] 
    
    @property
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "replicas": "int"
        }

    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
             "replicas": ["1"]
        }

    @property
    def kubernetes_templates(self) -> Dict[str, Any]:
        """Return the Kubernetes manifest templates for the service."""
        return {
            "deployment": {
                "apiVersion": "apps/v1",
                "kind": "Deployment",
                "metadata": {
                    "name": "pinoprotocol",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "pinoprotocol"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "pinoprotocol"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "pinoprotocol",
                                "image": "maddsystems/pinoprotocol:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
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
                    "name": "pinoprotocol",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "pinoprotocol"
                    },
                    "ports": [{
                        "port": 1883,
                        "targetPort": 1883,
                        "protocol": "TCP"
                    }]
                }
            }
        }
    
    def _customize_manifests(self, templates: Dict[str, Any], client_name: str, params: Dict[str, Any]) -> None:
        """Customize the Kubernetes manifests with service-specific logic."""
        print(f"Customizing manifests for {self.service_name} with params:", params)
        
        # Update deployment
        deployment = templates["deployment"]
        deployment["metadata"]["namespace"] = client_name
        deployment["spec"]["replicas"] = params["replicas"]

        # Update service namespace
        if "service" in templates:
            templates["service"]["metadata"]["namespace"] = client_name
