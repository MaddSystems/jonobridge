from typing import Dict, Any, List
from ..base_service import BaseService

class WebsqldumpService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "websqldump"
    
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
            "mysql_user": "varchar(255)",         # MYSQL_USE
            "mysql_password": "varchar(255)",     # MYSQL_PASS
            "mysql_database": "varchar(255)",     # MYSQL_DB
            "portal_endpoint": "varchar(255)",
            "web_user": "varchar(255)",           # WEB_USER
            "web_password": "varchar(255)",       # WEB_PASSWORD
            "websqldump_port": "int",
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "mysql_user": ["gpscontrol"],         # MYSQL_USE
            "mysql_password": ["qazwsxedc"],     # MYSQL_PASS
            "mysql_database": ["gruasaya"],     # MYSQL_DB
            "portal_endpoint": ["gruasaya"],
            "web_user": ["gruasaya"],            # WEB_USER
            "web_password": ["tSdyvgtdCQvhXxMx"], # WEB_PASSWORD
            "websqldump_port": ["10000-19999"],
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
                    "name": "websqldump",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "websqldump"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "websqldump"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "websqldump",
                                "image": "maddsystems/websqldump:1.0.0",
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
                                        "name": "MYSQL_PASSWORD",
                                        "value": None
                                    },
                                    {
                                        "name": "MYSQL_DATABASE",
                                        "value": None
                                    },
                                    {
                                        "name": "WEB_USER",
                                        "value": None
                                    },
                                    {
                                        "name": "WEB_PASSWORD",
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
                    "name": "websqldump",
                    "namespace": None
                },
                "spec": {
                    "type": "LoadBalancer",
                    "selector": {
                        "app": "websqldump"
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
            if env["name"] == "MYSQL_HOST":
                env["value"] = "host.minikube.internal"
            if env["name"] == "MYSQL_PORT":
                env["value"] = "3306"
            if env["name"] == "MYSQL_USER":
                env["value"] = params["mysql_user"]
            if env["name"] == "MYSQL_PASSWORD":   
                env["value"] = params["mysql_password"]
            if env["name"] == "MYSQL_DATABASE":
                env["value"] = params["mysql_database"]
            if env["name"] == "WEB_USER":
                env["value"] = params["web_user"]
            if env["name"] == "WEB_PASSWORD":
                env["value"] = params["web_password"]
            if env["name"] == "PORTAL_ENDPOINT":
                env["value"] = params["portal_endpoint"]
        
        # Update service
                # Update service namespace
        if "service" in templates:
            service = templates["service"]
            templates["service"]["metadata"]["namespace"] = client_name
            service["spec"]["ports"][0].update({
                "port": params["websqldump_port"],
                "protocol": "TCP",
                "targetPort": 8081
            })