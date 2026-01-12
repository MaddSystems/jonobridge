from typing import Dict, Any, List
from ..base_service import BaseService

class HTTPRequestService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "httprequest"
    
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
            "http/get"
        ] 
    
    @property
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "http_url": "varchar(255)",
            "http_polling_time": "int",
            "skywave_access_id": "varchar(255)",
            "skywave_password": "varchar(255)",
            "skywave_from_id": "varchar(255)",
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "http_url": ["https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed/0Tu2D5pSArQ7u3M8ZTtafnh333Z47tIMc/message.xml,https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed/0O30YLmvuKNVn76RecWOeag8h01e2bHcW/message.xml,https://api.findmespot.com/spot-main-web/consumer/rest-api/2.0/public/feed/0ZHtnQqAsWBXh8MTH5vXBuiJgKRbImUIY/message.xml","https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml"],
            "http_polling_time": ["30", "180"],
            "skywave_access_id": ["70001184"],
            "skywave_password": ["JEUTPKKH"],
            "skywave_from_id": ["13969586728"],
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
                    "name": "httprequest",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "httprequest"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "httprequest"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "httprequest",
                                "image": "maddsystems/httprequest:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
                                    {
                                        "name": "HTTP_URL",
                                        "value": None
                                    },
                                    {
                                        "name": "HTTP_POLLING_TIME",
                                        "value": None
                                    },
                                    {
                                        "name": "SKYWAVE_ACCESS_ID",
                                        "value": None
                                    },
                                    {
                                        "name": "SKYWAVE_PASSWORD",
                                        "value": None
                                    },
                                    {
                                        "name": "SKYWAVE_FROM_ID",
                                        "value": None
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
                    "name": "httprequest-service",
                    "namespace": None
                },
                "spec": {
                    "type": "LoadBalancer",
                    "selector": {
                        "app": "httprequest"
                    },
                    "ports": [
                        {
                            "port": 80,
                            "targetPort": 80,
                            "protocol": "TCP"
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
            if env["name"] == "HTTP_URL":
                env["value"] = params["http_url"]
            elif env["name"] == "HTTP_POLLING_TIME":
                env["value"] = str(params["http_polling_time"])
            elif env["name"] == "SKYWAVE_ACCESS_ID":
                env["value"] = params.get("skywave_access_id", "")
            elif env["name"] == "SKYWAVE_PASSWORD":
                env["value"] = params.get("skywave_password", "")
            elif env["name"] == "SKYWAVE_FROM_ID":
                env["value"] = params.get("skywave_from_id", "")
        
        # Update service
        if "service" in templates:
            service = templates["service"]
            service["metadata"]["namespace"] = client_name
