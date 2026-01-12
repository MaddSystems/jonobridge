from typing import Dict, Any, List
from ..base_service import BaseService

class Send2MySql(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "send2mysql"  # This must match the database table name
    
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
            "mysql_user": "varchar(255)",         # MYSQL_USE
            "mysql_password": "varchar(255)",     # MYSQL_PASS
            "mysql_database": "varchar(255)",     # MYSQL_DB
            "mysql_update": "varchar(255)",       # MYSQL_UPDATE
            "plates_url": "varchar(255)",
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "mysql_user": ["gpscontrol"],         # MYSQL_USE
            "mysql_password": ["qazwsxedc"],     # MYSQL_PASS
            "mysql_database": ["bridge"],     # MYSQL_DB
            "mysql_update": ["N"],                # MYSQL_UPDATE
            "plates_url": ["https://pluto.dudewhereismy.com.mx/imei/search?appId=2911"],
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
                    "name": "send2mysql",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "send2mysql"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "send2mysql"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "send2mysql",
                                "image": "maddsystems/send2mysql:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
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
                                    {
                                        "name": "MYSQL_UPDATE",
                                        "value": None
                                    },
                                    {
                                        "name": "PLATES_URL",
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
                    "name": "send2mysql",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "send2mysql"
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
            if env["name"] == "MYSQL_UPDATE":
                env["value"] = params["mysql_update"]
            if env["name"] == "PLATES_URL":
                env["value"] = params["plates_url"]
        # Update service namespace
        if "service" in templates:
            templates["service"]["metadata"]["namespace"] = client_name
