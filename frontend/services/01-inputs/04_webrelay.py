from typing import Dict, Any, List
from ..base_service import BaseService

class WebrelayService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "webrelay"
    
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
            ""
        ] 
    
    @property
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "ggs_user": "varchar(255)",
            "ggs_password": "varchar(255)",
            "app_id": "int",
            "webrelay_token": "varchar(255)",
            "portal_endpoint": "varchar(255)",
            "webrelay_port": "int",   
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "ggs_user": ["admindesarrollo"],
            "ggs_password": ["GPSc0ntr0l00"],
            "app_id": ["424"],
            "webrelay_token": ["d655eea7616e05b35dc7b22dd83b6ebc"],
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
                    "name": "webrelay",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "webrelay"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "webrelay"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "webrelay",
                                "image": "maddsystems/webrelay:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
                                    {
                                        "name": "GGS_USER",
                                        "value": None
                                    },
                                    {
                                        "name": "GGS_PASSWORD",
                                        "value": None
                                    },
                                    {
                                        "name": "APP_ID",
                                        "value": None
                                    },
                                    {
                                        "name": "WEBRELAY_TOKEN",
                                        "value": None
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
                    "name": "webrelay",
                    "namespace": None
                },
                "spec": {
                    "type": "LoadBalancer",
                    "selector": {
                        "app": "webrelay"
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
            if env["name"] == "GGS_USER":
                env["value"] = params["ggs_user"]
            if env["name"] == "GGS_PASSWORD":   
                env["value"] = params["ggs_password"]
            if env["name"] == "APP_ID":
                env["value"] = str(params["app_id"])
            if env["name"] == "WEBRELAY_TOKEN":
                env["value"] = params["webrelay_token"]
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