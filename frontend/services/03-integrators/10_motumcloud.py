from typing import Dict, Any, List
from ..base_service import BaseService

class MotumcloudService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "motumcloud"  # This must match the database table name
    
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
            "motum_user": "varchar(255)",        
            "motum_password": "varchar(255)",        
            "motum_referer": "varchar(255)",
            "motum_apikey": "varchar(255)", 
            "motum_carrier": "varchar(255)",     
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "plates_url": ["https://pluto.dudewhereismy.com.mx/imei/search?appId=2834"],
            "elastic_doc_name": ["dafnegarrido_motumcloud"],
            "motum_user": ["mcloud-dafgarri-gpsctrl-cfort@dafgarri-gpsctrl.com"],        
            "motum_password": ["Ap1CeCl0udAfG4rGp5Ctrl"],     
            "motum_referer": ["mcloud-dafgarri-gpsctrl-cfort"],
            "motum_apikey": ["AIzaSyAgO6dk0ZKDO15_M6dJ7fClmcJ4_cHVm8c"],
            "motum_carrier": ["Dafne Garrido"],
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
                    "name": "motumcloud",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "motumcloud"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "motumcloud"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "motumcloud",
                                "image": "maddsystems/motumcloud:1.0.0",
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
                                        "name": "MOTUMCLOUD_USER",
                                        "value": None
                                    },
                                    {
                                        "name": "MOTUMCLOUD_PASSWORD",
                                        "value": None
                                    },
                                    {
                                        "name": "MOTUMCLOUD_REFERER",
                                        "value": None
                                    },
                                    {
                                        "name": "MOTUMCLOUD_APIKEY",
                                        "value": None
                                    },
                                    {
                                        "name": "MOTUMCLOUD_CARRIER",
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
                    "name": "motumcloud",
                    "namespace": None
                },
                "spec": {
                    "selector": {
                        "app": "motumcloud"
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
            if env["name"] == "MOTUMCLOUD_USER":
                env["value"] = params["motum_user"]
            if env["name"] == "MOTUMCLOUD_PASSWORD":
                env["value"] = params["motum_password"]     
            if env["name"] == "MOTUMCLOUD_REFERER":
                env["value"] = params["motum_referer"]
            if env["name"] == "MOTUMCLOUD_APIKEY":
                env["value"] = params["motum_apikey"]
            if env["name"] == "MOTUMCLOUD_CARRIER":
                env["value"] = params["motum_carrier"]
        # Update service namespace
        if "service" in templates:
            templates["service"]["metadata"]["namespace"] = client_name
