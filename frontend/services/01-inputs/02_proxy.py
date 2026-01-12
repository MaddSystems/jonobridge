from typing import Dict, Any, List
from ..base_service import BaseService

class ProxyService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "proxy"
    
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
            "tracker/from-tcp"
        ] 

    @property
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "platform_host": "varchar(255)",
            "port": "int",
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "platform_host": ["mdvr.trackermexico.com.mx:50005","server1.gpscontrol.com.mx:8500","ms03.trackermexico.com.mx:10003"],
            "port": ["8000-8999"],
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
                    "name": "proxy",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "proxy"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "proxy"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "proxy",
                                "image": "maddsystems/proxy:1.0.0",
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
                                        "name": "PLATFORM_HOST",
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
                    "name": "listener-service",
                    "namespace": None
                },
                "spec": {
                    "type": "LoadBalancer",
                    "selector": {
                        "app": "proxy"
                    },
                    "ports": [{
                        "name": "tcp",
                        "port": None,
                        "protocol": "TCP",
                        "targetPort": 1024,
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
        for env in container["env"]:
            if env["name"] == "PLATFORM_HOST":
                env["value"] = params["platform_host"]
        
        # Update service
        if "service" in templates:
            service = templates["service"]
            service["metadata"]["namespace"] = client_name
            service["spec"]["ports"][0].update({
                "port": params["port"],
            })
            