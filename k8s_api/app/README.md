# Kubernetes Management API

## Overview

The **Kubernetes Management API** is a Python-based RESTful API built with [FastAPI](https://fastapi.tiangolo.com/) that mirrors common `kubectl` commands. It allows you to manage your Kubernetes clusters, namespaces, pods, deployments, services, and more through a simple and intuitive HTTP interface. The API is documented using Swagger, providing an interactive UI for exploration and testing.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the API](#running-the-api)
- [Accessing the API Documentation](#accessing-the-api-documentation)
- [API Endpoints](#api-endpoints)
  - [Cluster Management](#cluster-management)
  - [Working with Namespaces](#working-with-namespaces)
  - [Working with Pods](#working-with-pods)
  - [Working with Deployments](#working-with-deployments)
  - [Working with Services](#working-with-services)
  - [Working with ConfigMaps and Secrets](#working-with-configmaps-and-secrets)
  - [Working with Persistent Volumes (PV) and Persistent Volume Claims (PVC)](#working-with-persistent-volumes-pv-and-persistent-volume-claims-pvc)
  - [Working with ReplicaSets and StatefulSets](#working-with-replicasets-and-statefulsets)
  - [Working with Jobs and CronJobs](#working-with-jobs-and-cronjobs)
  - [Working with Resources](#working-with-resources)
  - [Accessing and Debugging](#accessing-and-debugging)
- [Nginx Configuration](#nginx-configuration)
- [Security Considerations](#security-considerations)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Cluster Management**: Retrieve cluster information and manage nodes.
- **Namespaces**: List, create, and delete namespaces.
- **Pods**: Manage pods, retrieve logs, and execute commands within pods.
- **Deployments**: Handle deployments, scale replicas, and delete deployments.
- **Services**: Manage Kubernetes services, expose pods, and delete services.
- **ConfigMaps and Secrets**: Create, list, describe, and delete ConfigMaps and Secrets.
- **Persistent Volumes and Claims**: Manage Persistent Volumes (PV) and Persistent Volume Claims (PVC).
- **ReplicaSets and StatefulSets**: Handle ReplicaSets and StatefulSets.
- **Jobs and CronJobs**: Manage batch jobs and scheduled CronJobs.
- **Resources**: Apply, delete, edit, and replace Kubernetes resources using YAML files.
- **Accessing and Debugging**: Port forwarding, resource usage metrics, event logs, and version information.

## Prerequisites

- **Python 3.7+**: Ensure you have Python installed. You can download it from [python.org](https://www.python.org/downloads/).
- **Kubernetes Cluster**: Access to a Kubernetes cluster. The API can be run locally or within the cluster.
- **Kubeconfig**: Properly configured kubeconfig file for accessing the Kubernetes cluster.
- **Nginx**: Installed and configured to proxy requests to the API.

## Installation

1. **Clone the Repository**
   ```bash
   git git clone git@bitbucket.org:maddisontest/jonobridge.git

   ```
   
2. **Run locally**
  ```
  cd jonobridge/k8s_api/app
  uvicorn main:app --host 0.0.0.0 --port 8000
  ```
