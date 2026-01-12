from typing import Dict, Any, List
from ..base_service import BaseService

class meitrackService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "meitrack"  # This must match the database table name
    
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
            "meitrack_host":"varchar(255)",
            "meitrack_fwd_only_adas": "varchar(255)",
            "meitrack_mock_imei": "varchar(1)",
            "meitrack_mock_value": "varchar(255)",
            "elastic_doc_name": "varchar(255)",
            "elastic_url": "varchar(255)",
            "elastic_user": "varchar(255)",
            "elastic_password": "varchar(255)",
            "send_to_elastic": "varchar(255)",
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "meitrack_host": ["server1.gpscontrol.com.mx:8500,ms03.trackermexico.com.mx:10005"],
            "meitrack_fwd_only_adas": ["N"],
            "meitrack_mock_imei": ["N"],
            "meitrack_mock_value": ["5"],
            "elastic_doc_name": ["xxxx_meitrack"],
            "elastic_url": ["https://opensearch.madd.com.mx:9200"],
            "elastic_user": ["admin"],
            "elastic_password": ["GPSc0ntr0l1"],
            "send_to_elastic": ["Y"],
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
                    "name": "meitrack",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "meitrack"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "meitrack"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "meitrack",
                                "image": "maddsystems/meitrack:1.0.0",
                                "imagePullPolicy": "Always",
                                "env": [
                                    {
                                        "name": "MQTT_BROKER_HOST",
                                        "value": "mosquitto"
                                    },
                                    {
                                        "name": "MEITRACK_HOST",
                                        "value": None
                                    },
                                                                        {
                                        "name": "MEITRACK_FWD_ONLY_ADAS",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "MEITRACK_MOCK_IMEI",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "MEITRACK_MOCK_VALUE",
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
                                        "name": "SEND_TO_ELASTIC",
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
                    "name": "meitrack",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "meitrack"
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
            if env["name"] == "CLIENT_ID":
                env["value"] = client_name
            if env["name"] == "MEITRACK_HOST":
                env["value"] = params["meitrack_host"]
            if env["name"] == "MEITRACK_FWD_ONLY_ADAS": 
                env["value"] = params["meitrack_fwd_only_adas"]
            if env["name"] == "MEITRACK_MOCK_IMEI":
                env["value"] = params["meitrack_mock_imei"]
            if env["name"] == "MEITRACK_MOCK_VALUE":
                env["value"] = params["meitrack_mock_value"]
            if env["name"] == "ELASTIC_DOC_NAME":
                env["value"] = params["elastic_doc_name"]
            if env["name"] == "ELASTIC_URL":
                env["value"] = params["elastic_url"]
            if env["name"] == "ELASTIC_USER":
                env["value"] = params["elastic_user"]
            if env["name"] == "ELASTIC_PASSWORD":
                env["value"] = params["elastic_password"]
            if env["name"] == "SEND_TO_ELASTIC":
                env["value"] = params["send_to_elastic"]
        # Update service namespace
        if "service" in templates:
            templates["service"]["metadata"]["namespace"] = client_name
