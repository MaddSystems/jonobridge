from typing import Dict, Any, List
from ..base_service import BaseService

class UScentralService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "uscentral"  # This must match the database table name
    
    @property
    def service_type(self) -> str:
        """Return the name of the type of service, which matches the database table name."""
        return "integration"  # This must match the database table name
    
    @property
    def inputs(self) -> List[str]:
        """Return the inputs of the service."""
        return [
            "tracker/jonoprotocol",
        ]   

    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "uscentral_user": ["gpstrackerservice@tecnotrust.com.mx"],
            "uscentral_user_key": ["6Gn^J3c$LU"],
            "uscentral_url": ["https://us-central1-trust-11b7d.cloudfunctions.net"],
            "uscentral_token": ["XXXXXXX"],
            "elastic_doc_name": ["xxxxxx_uscentral"],
            "elastic_url": ["https://opensearch.madd.com.mx:9200"],
            "elastic_user": ["admin"],
            "elastic_password": ["GPSc0ntr0l1"],
            "replicas": ["1"]
        }

    @property
    def outputs(self) -> List[str]:
        """Return the outputs of the service."""
        return [] 

    @property
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "uscentral_user": "varchar(255)",
            "uscentral_user_key": "varchar(255)",
            "uscentral_url": "varchar(255)",
            "uscentral_token": "varchar(255)",
            "elastic_doc_name": "varchar(255)",
            "elastic_url": "varchar(255)",
            "elastic_user": "varchar(255)",
            "elastic_password": "varchar(255)",
            "replicas": "int",  
        }
    
    
    @property
    def kubernetes_templates(self) -> Dict[str, Any]:
        """Return the Kubernetes manifest templates for the service."""
        return {
            "deployment": {
                "apiVersion": "apps/v1",
                "kind": "Deployment",
                "metadata": {
                    "name": "uscentral",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "uscentral"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "uscentral"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "uscentral",
                                "image": "maddsystems/uscentral:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
                                    {
                                        "name": "USCENTRAL_USER",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "USCENTRAL_USER_KEY",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "USCENTRAL_URL",
                                        "value": None
                                    },
                                    {
                                        "name": "USCENTRAL_TOKEN",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "ELASTIC_DOC_NAME",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "ELASTIC_URL",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "ELASTIC_USER",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "ELASTIC_PASSWORD",
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
                    "name": "uscentral",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "uscentral"
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
        
        # Update environment variables
        container = deployment["spec"]["template"]["spec"]["containers"][0]
        for env in container["env"]:
            if env["name"] == "USCENTRAL_USER":
                env["value"] = params["uscentral_user"]
            if env["name"] == "USCENTRAL_USER_KEY":
                env["value"] = params["uscentral_user_key"]
            if env["name"] == "USCENTRAL_URL":
                env["value"] = params["uscentral_url"]
            if env["name"] == "USCENTRAL_TOKEN":
                env["value"] = params["uscentral_token"]
            if env["name"] == "ELASTIC_DOC_NAME":
                env["value"] = params["elastic_doc_name"]
            if env["name"] == "ELASTIC_URL":
                env["value"] = params["elastic_url"]
            if env["name"] == "ELASTIC_USER":
                env["value"] = params["elastic_user"]
            if env["name"] == "ELASTIC_PASSWORD":
                env["value"] = params["elastic_password"]   

        # Update service namespace
        if "service" in templates:
            templates["service"]["metadata"]["namespace"] = client_name
