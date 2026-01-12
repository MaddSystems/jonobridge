from pydantic import BaseModel
from typing import List, Optional

# --------------------------
# Helper Classes
# --------------------------

class PodSummary(BaseModel):
    name: str
    ready: str
    status: str
    restarts: int
    age: str

class ServiceSummary(BaseModel):
    name: str
    type: str
    cluster_ip: Optional[str]
    external_ip: Optional[str]
    ports: str
    age: str

class DeleteDeploymentsResponse(BaseModel):
    deleted_deployments: List[str]
    namespace: str
    message: Optional[str] = None

class DeleteServicesResponse(BaseModel):
    deleted_services: List[str]
    namespace: str
    message: Optional[str] = None

class NamespaceSummary(BaseModel):
    name: str
    status: str
    age: str






