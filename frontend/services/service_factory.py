import importlib
import os
import pkgutil
from typing import Dict, Type, Optional
from .base_service import BaseService

class ServiceFactory:
    _instance = None
    _services: Dict[str, Type[BaseService]] = {}
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super(ServiceFactory, cls).__new__(cls)
            cls._instance._load_services()
        return cls._instance
    
    def _load_services(self) -> None:
        services_dir = os.path.dirname(__file__)
        def load_from_dir(path: str, package_prefix: str = "services") -> None:
            for finder, name, _ in pkgutil.iter_modules([path]):
                if name != "base_service" and name != "service_factory":
                    full_path = os.path.join(path, name)
                    if os.path.isdir(full_path):
                        # Recursively load from subdirectory
                        load_from_dir(full_path, f"{package_prefix}.{name}")
                    else:
                        # Load the module
                        module = importlib.import_module(f".{name}", package=package_prefix)
                        for attr_name in dir(module):
                            attr = getattr(module, attr_name)
                            if (isinstance(attr, type) and 
                                issubclass(attr, BaseService) and 
                                attr != BaseService):
                                service = attr()
                                self._services[service.service_name] = attr
        load_from_dir(services_dir)

    def get_service(self, service_name: str) -> Optional[BaseService]:
        """Get a service instance by name."""
        service_class = self._services.get(service_name)
        if service_class:
            return service_class()
        return None
    
    def get_all_services(self) -> Dict[str, BaseService]:
        """Get all available services."""
        return {name: cls() for name, cls in self._services.items()}

    def get_all_protocols(self) -> Dict[str, BaseService]:
        """Get all available services."""
        return {name: cls() for name, cls in self._services.items()}
