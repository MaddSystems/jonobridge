#!/bin/bash
# ============================================================================
# Script de deployment para Grule Universal
# ============================================================================
set -e

echo "=== Deployment de Grule Universal ==="

# Variables
NAMESPACE="testgrule"
IMAGE_NAME="grule-universal"
IMAGE_TAG="latest"
REGISTRY="your-registry"  # Cambiar por tu registry (e.g., gcr.io/project-id)

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Función de log
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# 1. Build de la imagen Docker
log "Building Docker image..."
docker build -t ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} . || error "Docker build failed"

# 2. Push a registry
log "Pushing image to registry..."
docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} || error "Docker push failed"

# 3. Verificar que el namespace existe
log "Verificando namespace ${NAMESPACE}..."
kubectl get namespace ${NAMESPACE} > /dev/null 2>&1 || {
    warn "Namespace ${NAMESPACE} no existe, creando..."
    kubectl create namespace ${NAMESPACE}
}

# 4. Verificar que los secrets existen
log "Verificando secrets..."
kubectl get secret grule-secrets -n ${NAMESPACE} > /dev/null 2>&1 || {
    warn "Secret 'grule-secrets' no existe. Debes crearlo con:"
    echo "kubectl create secret generic grule-secrets -n ${NAMESPACE} \\"
    echo "  --from-literal=mysql-host=HOST \\"
    echo "  --from-literal=mysql-user=USER \\"
    echo "  --from-literal=mysql-password=PASS \\"
    echo "  --from-literal=telegram-bot-token=TOKEN \\"
    echo "  --from-literal=mqtt-broker-host=BROKER"
    error "Secrets requeridos no encontrados"
}

# 5. Actualizar la imagen en el deployment
log "Actualizando deployment YAML con nueva imagen..."
sed -i "s|image:.*grule-universal.*|image: ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}|g" k8s-deployment.yaml

# 6. Apply del deployment
log "Aplicando deployment en Kubernetes..."
kubectl apply -f k8s-deployment.yaml -n ${NAMESPACE} || error "kubectl apply failed"

# 7. Esperar a que el deployment esté ready
log "Esperando a que el pod esté listo..."
kubectl rollout status deployment/grule-universal -n ${NAMESPACE} --timeout=300s || error "Rollout timeout"

# 8. Obtener info del pod
log "Deployment completado exitosamente!"
POD_NAME=$(kubectl get pods -n ${NAMESPACE} -l app=grule-universal -o jsonpath='{.items[0].metadata.name}')
log "Pod: ${POD_NAME}"

# 9. Mostrar logs
log "Mostrando logs recientes..."
kubectl logs -n ${NAMESPACE} ${POD_NAME} --tail=50

# 10. Info de acceso
log "=== Información de acceso ==="
echo "Backend API: http://localhost:8080 (port-forward)"
echo "Frontend Web: http://localhost:5000 (port-forward)"
echo ""
echo "Para acceder localmente:"
echo "kubectl port-forward -n ${NAMESPACE} ${POD_NAME} 8080:8080 5000:5000"

log "✅ Deployment completado"
