from typing import Dict, Any, List
from ..base_service import BaseService

class GruleService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "grule"  # This must match the database table name
    
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
            "mysql_user": "varchar(255)",         # MYSQL_USER
            "mysql_password": "varchar(255)",     # MYSQL_PASS
            "mysql_database": "varchar(255)",     # MYSQL_DB
            "mysql_update": "varchar(255)",       # MYSQL_UPDATE
            "portal_endpoint": "varchar(255)",    # PORTAL_ENDPOINT
            "telegram_bot_token": "varchar(255)", # TELEGRAM_BOT_TOKEN
            "telegram_chat_id": "varchar(255)",   # TELEGRAM_CHAT_ID
            "grule_web_port": "int",             # Service port
            "grule_audit_enabled": "varchar(1)", # GRULE_AUDIT_ENABLED
            "grule_audit_level": "varchar(255)", # GRULE_AUDIT_LEVEL
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "mysql_user": ["gpscontrol"],         # MYSQL_USER
            "mysql_password": ["qazwsxedc"],     # MYSQL_PASS
            "mysql_database": ["grule"],         # MYSQL_DB
            "mysql_update": ["N"],               # MYSQL_UPDATE
            "portal_endpoint": ["grule"],        # PORTAL_ENDPOINT
            "telegram_bot_token": ["1437635839:AAEVHvYYzzBaoya42zA1_1X9X1RQjlpdlUo"], # TELEGRAM_BOT_TOKEN
            "telegram_chat_id": ["-1002135388607"], # TELEGRAM_CHAT_ID
            "grule_web_port": ["10003"],        # Service port
            "grule_audit_enabled": ["Y","N"],       # GRULE_AUDIT_ENABLED
            "grule_audit_level": ["NONE","ERROR","ALL"],      # GRULE_AUDIT_LEVEL
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
                    "name": "grule",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "grule"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "grule"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "grule",
                                "image": "maddsystems/grule:1.0.0",
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
                                        "name": "PORTAL_ENDPOINT",
                                        "value": None
                                    },
                                    {
                                        "name": "TELEGRAM_BOT_TOKEN",
                                        "value": None
                                    },
                                    {
                                        "name": "TELEGRAM_CHAT_ID",
                                        "value": None
                                    },
                                    {
                                        "name": "GRULE_AUDIT_ENABLED",
                                        "value": None
                                    },
                                    {
                                        "name": "GRULE_AUDIT_LEVEL",
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
                    "name": "grule",
                    "namespace": None
                },
                "spec": {
                    "type": "LoadBalancer",
                    "selector": {
                        "app": "grule"
                    },
                    "ports": [{
                        "name": "tcp-portal",
                        "port": None,
                        "protocol": "TCP",
                        "targetPort": 8081,
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
            if env["name"] == "PORTAL_ENDPOINT":
                env["value"] = params["portal_endpoint"]
            if env["name"] == "TELEGRAM_BOT_TOKEN":
                env["value"] = params["telegram_bot_token"]
            if env["name"] == "TELEGRAM_CHAT_ID":
                env["value"] = params["telegram_chat_id"]
            if env["name"] == "GRULE_AUDIT_ENABLED":
                env["value"] = params["grule_audit_enabled"]
            if env["name"] == "GRULE_AUDIT_LEVEL":
                env["value"] = params["grule_audit_level"]
        
        # Update service namespace and port
        if "service" in templates:
            service = templates["service"]
            templates["service"]["metadata"]["namespace"] = client_name
            service["spec"]["ports"][0].update({
                "port": params["grule_web_port"],
                "protocol": "TCP",
                "targetPort": 8081
            })
