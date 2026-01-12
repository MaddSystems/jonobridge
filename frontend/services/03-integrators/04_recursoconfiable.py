from typing import Dict, Any, List
from ..base_service import BaseService

class RecursoConfiableService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "recursoconfiable"  # This must match the database table name
    
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
    def outputs(self) -> List[str]:
        """Return the outputs of the service."""
        return [] 

    @property
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "user_name": "varchar(255)",
            "password": "varchar(255)",
            "url": "varchar(255)",
            "customer_id": "varchar(255)",
            "customer_name": "varchar(255)",
            "plates_url": "varchar(255)",
            "elastic_doc_name": "varchar(255)",
            "elastic_url": "varchar(255)",
            "elastic_user": "varchar(255)",
            "elastic_password": "varchar(255)",
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "customer_id": ["41013"],
            "customer_name": ["Juan Pablo Aguilar"],
            "user_name": ["ws_avl_controlGPS"],
            "password": ["NUSU#294AVJm$2"],
            "url": ["http://gps.rcontrol.com.mx/Tracking/wcf/RCService.svc?wsdl"],
            "plates_url": ["https://pluto.dudewhereismy.com.mx/imei/search?appId=2911"],
            "elastic_doc_name": ["logisticknc_rc"],
            "elastic_url": ["https://opensearch.madd.com.mx:9200"],
            "elastic_user": ["admin"],
            "elastic_password": ["GPSc0ntr0l1"],
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
                    "name": "recursoconfiable",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "recursoconfiable"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "recursoconfiable"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "recursoconfiable",
                                "image": "maddsystems/recursoconfiable:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
                                    {
                                        "name": "USER",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "PASSWORD",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "URL",
                                        "value": None
                                    },
                                    {
                                        "name": "CUSTOMER_NAME",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "CUSTOMER_ID",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "PLATES_URL",
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
                                    ,
                                    {
                                        "name": "CLIENT_ID",
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
                    "name": "recursoconfiable",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "recursoconfiable"
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
            if env["name"] == "USER":
                env["value"] = params["user_name"]
            if env["name"] == "PASSWORD":
                env["value"] = params["password"]
            if env["name"] == "URL":
                env["value"] = params["url"]
            if env["name"] == "CUSTOMER_ID":
                env["value"] = params["customer_id"]
            if env["name"] == "CUSTOMER_NAME":
                env["value"] = params["customer_name"]
            if env["name"] == "PLATES_URL":
                env["value"] = params["plates_url"]
            if env["name"] == "ELASTIC_DOC_NAME":
                env["value"] = params["elastic_doc_name"]
            if env["name"] == "ELASTIC_URL":
                env["value"] = params["elastic_url"]
            if env["name"] == "ELASTIC_USER":
                env["value"] = params["elastic_user"]
            if env["name"] == "ELASTIC_PASSWORD":
                env["value"] = params["elastic_password"]
            if env["name"] == "CLIENT_ID":
                env["value"] = client_name

        # Update service namespace
        if "service" in templates:
            templates["service"]["metadata"]["namespace"] = client_name
