from typing import Dict, Any, List
from ..base_service import BaseService

class AvocadocloudService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "avocadocloud"  # This must match the database table name
    
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
            "plates_url": "varchar(255)",
            "elastic_doc_name": "varchar(255)",
            "avocado_user": "varchar(255)",        
            "avocado_password": "varchar(255)",     
            "avocado_user_adm": "varchar(255)",     
            "avocado_url": "varchar(255)",     
            "mysql_user": "varchar(255)",         # MYSQL_USER
            "mysql_password": "varchar(255)",     # MYSQL_PASSWORD
            "mysql_database": "varchar(255)",     # MYSQL_DB
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "plates_url": ["https://pluto.dudewhereismy.com.mx/imei/search?appId=2834"],
            "elastic_doc_name": ["dafnegarrido_avocadocloud"],
            "avocado_user": ["phoenix"],        
            "avocado_password": ["Ph03niX-2018_"],     
            "avocado_user_adm": ["2135"],     
            "avocado_url": ["https://cerberusenlinea.com/WEB_SERVICE_PHOENIX_CLOUD_PRD-1.0/PH_PHOENIX_CLOUD_PRD_v01?wsdl/recibirEventosGPS"], 
            "mysql_user": ["gpscontrol"],         # MYSQL_USER
            "mysql_password": ["qazwsxedc"],      # MYSQL_PASSWORD
            "mysql_database": ["bridge"],         # MYSQL_DB
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
                    "name": "avocadocloud",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "avocadocloud"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "avocadocloud"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "avocadocloud",
                                "image": "maddsystems/avocadocloud:1.0.0",
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
                                        "name": "ELASTIC_DOC_NAME",
                                        "value": None
                                    },
                                    {
                                        "name": "AVOCADO_USER",
                                        "value": None
                                    },
                                    {
                                        "name": "AVOCADO_PASSWORD",
                                        "value": None
                                    },
                                    {
                                        "name": "AVOCADO_USER_ADM",
                                        "value": None
                                    },
                                    {
                                        "name": "AVOCADO_URL",
                                        "value": None
                                    },
                                    {
                                        "name": "MYSQL_HOST",
                                        "value": None
                                    },
                                    {
                                        "name": "MYSQL_PORT",
                                        "value": None
                                    },
                                    {
                                        "name": "MYSQL_USER",
                                        "value": None
                                    },
                                    {
                                        "name": "MYSQL_PASS",
                                        "value": None
                                    },
                                    {
                                        "name": "MYSQL_DB",
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
                    "name": "avocadocloud",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "avocadocloud"
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
            if env["name"] == "ELASTIC_DOC_NAME":
                env["value"] = params["elastic_doc_name"]
            if env["name"] == "AVOCADO_USER":
                env["value"] = params["avocado_user"]
            if env["name"] == "AVOCADO_PASSWORD":
                env["value"] = params["avocado_password"]
            if env["name"] == "AVOCADO_USER_ADM":
                env["value"] = params["avocado_user_adm"]
            if env["name"] == "AVOCADO_URL":
                env["value"] = params["avocado_url"]
            if env["name"] == "MYSQL_HOST":
                env["value"] = "host.minikube.internal"
            if env["name"] == "MYSQL_PORT":
                env["value"] = "3306"
            if env["name"] == "MYSQL_USER":
                env["value"] = params["mysql_user"]
            if env["name"] == "MYSQL_PASS":
                env["value"] = params["mysql_password"]
            if env["name"] == "MYSQL_DB":
                env["value"] = params["mysql_database"]
        # Update service namespace
        if "service" in templates:
            templates["service"]["metadata"]["namespace"] = client_name
