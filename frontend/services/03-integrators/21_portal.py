from typing import Dict, Any, List
from ..base_service import BaseService

class PortalService(BaseService):
    @property
    def service_name(self) -> str:
        """Return the name of the service, which matches the database table name."""
        return "portal"  # This must match the database table name
    
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
            "plates_url": "varchar(255)",   
            "portal_endpoint": "varchar(255)",    
            "portal_port": "int",    
            "portal_user": "varchar(255)",
            "portal_password": "text",
            "portal_script": "text",
            "replicas": "int"
        }
    
    @property
    def parameters_helpers(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        return {
            "mysql_user": ["gpscontrol"],         # MYSQL_USE
            "mysql_password": ["qazwsxedc"],     # MYSQL_PASS
            "mysql_database": ["bridge"],     # MYSQL_DB
            "plates_url": ["https://pluto.dudewhereismy.com.mx/imei/search?appId=179","https://pluto.dudewhereismy.com.mx/imei/search?appId=3059"],
            "portal_endpoint": ["api/v1.0/semov"],
            "portal_port": ["10000-19999"],
            "portal_user": [""],
            "portal_password": [""],
            "portal_script": ["""{ "Hora": "Hora", "Latitud": "Latitud", "Empresa": "Empresa", "Velocidad": "Velocidad", "IDEmpresa": "IDEmpresa", "Fecha": "Fecha", "Altitud": "Altitud", "UrlCamara": "UrlCamara", "NombreProveedor": "NombreProveedor", "BotonPanico": "BotonPanico", "IMEI": "IMEI", "NumeroEconomico": "NumeroEconomico", "SerieVehicularVIN": "SerieVehicularVIN", "Direccion": "Direccion", "Placas": "Placas", "Longitud": "Longitud", "NombreRuta": "NombreRuta" }"""],
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
                    "name": "portal",
                    "namespace": None
                },
                "spec": {
                    "replicas": None,
                    "selector": {
                        "matchLabels": {
                            "app": "portal"
                        }
                    },
                    "template": {
                        "metadata": {
                            "labels": {
                                "app": "portal"
                            }
                        },
                        "spec": {
                            "initContainers": [{
                                "name": "wait-for-mosquitto",
                                "image": "busybox",
                                "command": ["sh", "-c", "until nc -z mosquitto 1883; do echo waiting for mosquitto; sleep 2; done"]
                            }],
                            "containers": [{
                                "name": "portal",
                                "image": "maddsystems/portal:1.0.0",
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
                                        "name": "PLATES_URL",
                                        "value": None
                                    },
                                    {
                                        "name": "PORTAL_ENDPOINT",
                                        "value": None
                                    },
                                    {
                                        "name": "PORTAL_USER",
                                        "value": None
                                    },
                                    {
                                        "name": "PORTAL_PASSWORD",
                                        "value": None
                                    }
                                    ,
                                    {
                                        "name": "PORTAL_SCRIPT",
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
                    "name": "portal",
                    "namespace": None
                },
                "spec": {
                    "type": "LoadBalancer",
                    "selector": {
                        "app": "portal"
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
            if env["name"] == "PLATES_URL":
                env["value"] = params["plates_url"]
            if env["name"] == "PORTAL_ENDPOINT":
                env["value"] = params["portal_endpoint"]    
            if env["name"] == "PORTAL_USER":
                env["value"] = params["portal_user"]
            if env["name"] == "PORTAL_PASSWORD":
                env["value"] = params["portal_password"]
            if env["name"] == "PORTAL_SCRIPT":
                env["value"] = params["portal_script"]

        # Update service namespace
        if "service" in templates:
            service = templates["service"]
            templates["service"]["metadata"]["namespace"] = client_name
            service["spec"]["ports"][0].update({
                "port": params["portal_port"],
                "protocol": "TCP",
                "targetPort": 8081
            })
