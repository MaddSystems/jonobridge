from typing import Dict, Any, List
from ..base_service import BaseService

class NstechService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "nstech"  # This must match the database table name
    
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
            "plates_url": ["https://pluto.dudewhereismy.com.mx/imei/search?appId=2357"],
            "nstech_client_id": ["52466691-482f-48a0-adfc-a68e776eb966"],
            "nstech_client_secret": ["m6qUGrJ7dEVYTeAPLeV3BRVlEkveZCF8"],
            "nstech_technology_id": ["52466691-482f-48a0-adfc-a68e776eb966"],
            "nstech_account_id": ["52a4b1da-8e17-49c5-b490-d98ff1b390e0"],
            "nstech_url": ["https://zeus.nstech.com.br/api"],
            "nstech_token_url": ["https://auth.nstech.com.br/realms/zeus/protocol/openid-connect/token"],
            "elastic_doc_name": ["nstech"],
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
            "nstech_client_id": "varchar(255)",
            "nstech_client_secret": "varchar(255)",
            "nstech_technology_id": "varchar(255)",
            "nstech_account_id": "varchar(255)",
            "nstech_url": "varchar(255)",
            "nstech_token_url": "varchar(255)",
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
                    "name": "nstech",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "nstech"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "nstech"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "nstech",
                                "image": "maddsystems/nstech:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
                                    {
                                        "name": "plates_url",
                                        "value": None
                                    },
                                    {
                                        "name": "NSTECH_CLIENT_ID",
                                        "value": None
                                    },
                                    {
                                        "name": "NSTECH_CLIENT_SECRET",
                                        "value": None
                                    },
                                    {
                                        "name": "NSTECH_TECHNOLOGY_ID",
                                        "value": None
                                    },
                                    {
                                        "name": "NSTECH_ACCOUNT_ID",
                                        "value": None
                                    },
                                    {
                                        "name": "NSTECH_URL",
                                        "value": None
                                    },
                                    {
                                        "name": "NSTECH_TOKEN_URL",
                                        "value": None
                                    },
                                    {
                                        "name": "ELASTIC_DOC_NAME",
                                        "value": None
                                    },
                                    {
                                        "name": "ELASTIC_URL",
                                        "value": None
                                    },
                                    {
                                        "name": "ELASTIC_USER",
                                        "value": None
                                    },
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
                    "name": "nstech",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "nstech"
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
            if env["name"] == "plates_url":
                env["value"] = params["plates_url"]
            if env["name"] == "NSTECH_CLIENT_ID":
                env["value"] = params["nstech_client_id"]
            if env["name"] == "NSTECH_CLIENT_SECRET":
                env["value"] = params["nstech_client_secret"]
            if env["name"] == "NSTECH_TECHNOLOGY_ID":
                env["value"] = params["nstech_technology_id"]
            if env["name"] == "NSTECH_ACCOUNT_ID":
                env["value"] = params["nstech_account_id"]
            if env["name"] == "NSTECH_URL":
                env["value"] = params["nstech_url"]
            if env["name"] == "NSTECH_TOKEN_URL":
                env["value"] = params["nstech_token_url"]
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
