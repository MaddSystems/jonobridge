from abc import ABC, abstractmethod
from typing import Dict, Any, List

class BaseService(ABC):
    """Base class for all services."""
    
    @property
    @abstractmethod
    def service_name(self) -> str:
        """Return the name of the service."""
        pass
    
    @property
    @abstractmethod
    def service_type(self) -> str:
        """Return the kind of the service."""
        pass

    @property
    @abstractmethod
    def inputs(self) -> List[str]:
        """Return the inputs of the service."""
        pass

    @property
    @abstractmethod
    def outputs(self) -> List[str]:
        """Return the outputs of the service."""
        pass

    @property
    @abstractmethod
    def parameters(self) -> Dict[str, str]:
        """Return the parameters schema for the service."""
        pass
    
    @property
    @abstractmethod
    def parameters_helpers(self) -> Dict[str, List[str]]:
        """Return the parameters schema for the service."""
        pass

    @property
    @abstractmethod
    def kubernetes_templates(self) -> Dict[str, Any]:
        """Return the Kubernetes templates for the service."""
        pass
    
    def validate_parameters(self, params: Dict[str, Any]) -> bool:
        """Validate that the provided parameters match the schema."""
        required_params = self.parameters.keys()
        return all(param in params for param in required_params)
    
    def get_kubernetes_manifests(self, client_name: str, params: Dict[str, Any]) -> Dict[str, Any]:
        """Generate Kubernetes manifests for this service."""
        if not self.validate_parameters(params):
            raise ValueError(f"Invalid parameters for service {self.service_name}")
        
        templates = self.kubernetes_templates.copy()
        self._customize_manifests(templates, client_name, params)
        return templates
    
    @abstractmethod
    def _customize_manifests(self, templates: Dict[str, Any], client_name: str, params: Dict[str, Any]) -> None:
        """Customize the Kubernetes manifests with service-specific logic."""
        pass
