#!/bin/bash

# Verificar que los argumentos requeridos estén presentes
if [ "$#" -lt 2 ]; then
  echo "Uso: $0 <module_name> <docker_repository> [<dockerfile_path>]"
  exit 1
fi

# Parámetros
MODULE_NAME=$1
DOCKER_REPOSITORY=$2
DOCKERFILE_PATH=${3:-"./Dockerfile"} # Usar un Dockerfile por defecto si no se especifica
VERSION_FILE="${MODULE_NAME}_VERSION"

# Leer o inicializar la versión
if [ ! -f "$VERSION_FILE" ]; then
  echo "1.0.0" > "$VERSION_FILE"
fi
CURRENT_VERSION=$(cat $VERSION_FILE)

# Incrementar el patch
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"
PATCH=$((PATCH + 1))
NEW_VERSION="$MAJOR.$MINOR.$PATCH"
echo "$NEW_VERSION" > "$VERSION_FILE"

# Mostrar información
echo "Construyendo módulo: $MODULE_NAME"
echo "Repositorio Docker: $DOCKER_REPOSITORY"
echo "Versión: $NEW_VERSION"
echo "Dockerfile: $DOCKERFILE_PATH"

# Construir la imagen
docker build -t "$DOCKER_REPOSITORY/$MODULE_NAME:$NEW_VERSION" -f "$DOCKERFILE_PATH" .
if [ $? -ne 0 ]; then
  echo "Error al construir la imagen Docker."
  exit 1
fi

# Etiquetar como "latest"
docker tag "$DOCKER_REPOSITORY/$MODULE_NAME:$NEW_VERSION" "$DOCKER_REPOSITORY/$MODULE_NAME:latest"

# Subir al repositorio Docker
docker push "$DOCKER_REPOSITORY/$MODULE_NAME:$NEW_VERSION"
docker push "$DOCKER_REPOSITORY/$MODULE_NAME:latest"

echo "Construcción y publicación completadas para $MODULE_NAME:$NEW_VERSION"
