from typing import Dict, Any, List
from ..base_service import BaseService

class LobosoftwareService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "lobosoftware"  # This must match the database table name
    
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
            "plates_url": ["https://pluto.dudewhereismy.com.mx/imei/search?appId=3059"],
            "lobosoftware_user": ["GPSCONTROL"],
            "lobosoftware_user_key": ["GPSC0NTR0L.2024"],
            "lobosoftware_url": ["https://pruebascampanaws.lobos.com.mx/campana/gps/datagps"],
            "lobosoftware_token_url": ["https://pruebascampanaws.lobos.com.mx/campana/authorization/usergpsauthorization"],
            "elastic_doc_name": ["lobosoftware"],
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
            "plates_url": "varchar(255)",
            "lobosoftware_user": "varchar(255)",
            "lobosoftware_user_key": "varchar(255)",
            "lobosoftware_url": "varchar(255)",
            "lobosoftware_token_url": "varchar(255)",
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
                    "name": "lobosoftware",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "lobosoftware"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "lobosoftware"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "lobosoftware",
                                "image": "maddsystems/lobosoftware:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
                                    {
                                        "name": "PLATES_URL",
                                        "value": None
                                    },
                                    {
                                        "name": "LOBOSOFTWARE_USER",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "LOBOSOFTWARE_USER_KEY",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "LOBOSOFTWARE_URL",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "LOBOSOFTWARE_TOKEN_URL",
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
                    "name": "lobosoftware",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "lobosoftware"
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
            if env["name"] == "PLATES_URL":
                env["value"] = params["plates_url"]
            if env["name"] == "LOBOSOFTWARE_USER":
                env["value"] = params["lobosoftware_user"]
            if env["name"] == "LOBOSOFTWARE_USER_KEY":
                env["value"] = params["lobosoftware_user_key"]
            if env["name"] == "LOBOSOFTWARE_URL":
                env["value"] = params["lobosoftware_url"]
            if env["name"] == "LOBOSOFTWARE_TOKEN_URL":
                env["value"] = params["lobosoftware_token_url"]
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
