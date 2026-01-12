from typing import Dict, Any, List
from ..base_service import BaseService
import sys
sys.path.append('..')
from main import MINIKUBE_IP

class Gpsgatehttpervice(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "gpsgatehttp"  # This must match the database table name
    
    @property
    def service_type(self) -> str:
        """Return the name of the type of service, which matches the database table name."""
        return "interpreter"  # This must match the database table name
    
    @property
    def inputs(self) -> List[str]:
        """Return the inputs of the service."""
        return [
            "httpinput/get",
        ]   

    @property
    def outputs(self) -> List[str]:
        """Return the outputs of the service."""
        return []
         
    @property
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "elastic_doc_name": "varchar(255)",
            "elastic_url": "varchar(255)",
            "elastic_user": "varchar(255)",
            "elastic_password": "varchar(255)",
            "gpsgate_test_telegram": "varchar(255)",
            "telegram_api_url": "varchar(255)",
            "telegram_bot_token": "varchar(255)",
            "gpsgate_test_chat_id": "varchar(255)",
            "telegram_message_header": "varchar(255)",
            "telegram_additional_fields": "varchar(255)",
            "ego_api_url": "varchar(255)",
            "replicas": "int"
        }

    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "elastic_doc_name": ["gpsgatehttp"],
            "elastic_url": ["https://opensearch.madd.com.mx:9200"],
            "elastic_user": ["admin"],
            "elastic_password": ["GPSc0ntr0l1"],
            "gpsgate_test_telegram": ["Y"],
            "telegram_api_url": ["https://cygnus.dudewhereismy.com.mx/telegramgroups/get_telegra_by_appid"],
            "telegram_bot_token": ["1437635839:AAEVHvYYzzBaoya42zA1_1X9X1RQjlpdlUo"],
            "gpsgate_test_chat_id": ["-1002135388607"],
            "telegram_message_header": ["🚨 Alerta: %s"],
            "telegram_additional_fields": ["GEOFENCE_NAME,POS_ADDRESS,GEOFENCE_TAG_ID"],
            "ego_api_url": ["https://api2ego.elisasoftware.com.mx/catalog/doc/fac/6P2M3x1C024/347"],
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
                    "name": "gpsgatehttp",
                    "namespace": None
                },
                "spec": {
                    "revisionHistoryLimit": 2,  # Keep only 2 old ReplicaSets
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "gpsgatehttp"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "gpsgatehttp"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "gpsgatehttp",
                                "image": "maddsystems/gpsgatehttp:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
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
                                        "name": "GPSGATE_TEST_TELEGRAM",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "TELEGRAM_API_URL",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "TELEGRAM_BOT_TOKEN",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "GPSGATE_TEST_CHAT_ID",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "TELEGRAM_MESSAGE_HEADER",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "TELEGRAM_ADDITIONAL_FIELDS",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "EGO_API_URL",
                                        "value": None
                                    }
                                ],
                                "ports": [
                                        {
                                            "containerPort": 1883
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
                    "name": "gpsgatehttp",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "gpsgatehttp"
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
            if env["name"] == "ELASTIC_DOC_NAME":
                env["value"] = params["elastic_doc_name"]
            if env["name"] == "ELASTIC_URL":
                env["value"] = params["elastic_url"]
            if env["name"] == "ELASTIC_USER":
                env["value"] = params["elastic_user"]
            if env["name"] == "ELASTIC_PASSWORD":
                env["value"] = params["elastic_password"]
            if env["name"] == "GPSGATE_TEST_TELEGRAM":
                env["value"] = params["gpsgate_test_telegram"]
            if env["name"] == "TELEGRAM_API_URL":
                env["value"] = params["telegram_api_url"]
            if env["name"] == "TELEGRAM_BOT_TOKEN":
                env["value"] = params["telegram_bot_token"]
            if env["name"] == "GPSGATE_TEST_CHAT_ID":
                env["value"] = params["gpsgate_test_chat_id"]
            if env["name"] == "TELEGRAM_MESSAGE_HEADER":
                env["value"] = params.get("telegram_message_header")
            if env["name"] == "TELEGRAM_ADDITIONAL_FIELDS":
                env["value"] = params.get("telegram_additional_fields")
            if env["name"] == "EGO_API_URL":
                env["value"] = params["ego_api_url"]


        # Update service namespace
        if "service" in templates:
            templates["service"]["metadata"]["namespace"] = client_name
