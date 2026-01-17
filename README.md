# k8s-poc -- Execution and Operations Runbook (Enhanced Version)

This repository is an **educational POC** to learn Kubernetes
fundamentals using a Go (Gin) application, exposed via a Service, with
metrics collected by Prometheus and visualized in Grafana.

The focus is **not production**, but **deep and conscious
understanding** of the complete lifecycle:

-   application build
-   Docker image creation
-   Kubernetes deployment
-   service exposure
-   basic observability
-   operation, troubleshooting, and environment reset

Nothing is implicit. Everything is **explicit by design**.

------------------------------------------------------------------------

## Architecture Overview

    [ User / curl ]
            |
         NodePort
            |
         Service
            |
       Pods (3 replicas)
            |
         /metrics
            |
       Prometheus
            |
         Grafana

------------------------------------------------------------------------

## Project Structure

    .
    ├── cmd/
    │   └── server/                # Go application entrypoint
    ├── internal/                  # application code
    ├── k8s/
    │   ├── kustomization.yaml     # applies all manifests
    │   ├── namespaces.yaml        # monitoring namespace
    │   ├── app/
    │   │   ├── configmap.yaml     # application environment variables
    │   │   ├── deployment.yaml    # application deployment (3 replicas)
    │   │   └── service.yaml       # Service (NodePort)
    │   ├── prometheus/
    │   │   ├── configmap.yaml     # Prometheus scrape config
    │   │   ├── deployment.yaml    # Prometheus
    │   │   └── service.yaml       # Prometheus Service
    │   └── grafana/
    │       ├── configmap-datasource.yaml # Prometheus datasource
    │       ├── deployment.yaml           # Grafana
    │       └── service.yaml              # Grafana Service
    ├── Dockerfile
    └── go.mod

------------------------------------------------------------------------

## What Was Implemented

### Application

-   Go application using Gin
-   Endpoints:
    -   `/health`: health check + metrics
    -   `/whoami`: Pod identity (Downward API) + environment (ConfigMap)
    -   `/metrics`: Prometheus metrics

### Kubernetes

-   Deployment with **3 replicas**
-   **NodePort** Service for predictable access
-   ConfigMap for configuration without rebuild
-   Downward API for Pod identity
-   Liveness and readiness probes

### Observability

-   Prometheus scraping metrics via Service
-   Grafana with automatically provisioned datasource
-   Dashboard validating:
    -   health request rate
    -   latency
    -   per-Pod traffic distribution

------------------------------------------------------------------------

## Build the Application

Whenever **Go code changes**:

``` bash
go mod tidy
docker build -t k8s-poc:latest .
```

> Kubernetes **does not compile code**. It only runs images.

------------------------------------------------------------------------

## Deploy Everything

``` bash
kubectl apply -k k8s/
```

Validation:

``` bash
kubectl get pods
kubectl get pods -n monitoring
kubectl get svc
```

------------------------------------------------------------------------

## Access

### Application

    http://localhost:30080/health
    http://localhost:30080/whoami
    http://localhost:30080/metrics

### Prometheus

``` bash
kubectl port-forward -n monitoring svc/prometheus 9090:9090
```

### Grafana

``` bash
kubectl port-forward -n monitoring svc/grafana 3000:3000
```

Login:

    admin / admin

------------------------------------------------------------------------

## How to Discover Exposed Ports

``` bash
kubectl get svc
kubectl get svc --all-namespaces
```

Example:

    80:30080/TCP

-   80 → internal Service port
-   30080 → Node exposed port

Details:

``` bash
kubectl describe svc k8s-poc
```

------------------------------------------------------------------------

## When Something Changes

### Application code

``` bash
docker build -t k8s-poc:latest .
kubectl rollout restart deployment k8s-poc
```

### ConfigMap

``` bash
kubectl apply -f k8s/app/configmap.yaml
kubectl rollout restart deployment k8s-poc
```

### YAML files

``` bash
kubectl apply -f <file>
```

------------------------------------------------------------------------

## Tear Down

### Application only

``` bash
kubectl delete deployment k8s-poc
kubectl delete svc k8s-poc
kubectl delete configmap k8s-poc-config
```

### Observability

``` bash
kubectl delete namespace monitoring
```

### Full reset

``` bash
kubectl delete -k k8s/
```

------------------------------------------------------------------------

## Reinforced Concepts

-   NodePort is **cluster-wide**
-   Port-forward is a **debug tool**
-   High-cardinality metrics are educational
-   Namespace reset is normal in POCs

------------------------------------------------------------------------

## Final Goal

Learn, consciously: - how Kubernetes manages Pods - how Services
distribute traffic - how metrics are collected - how Prometheus and
Grafana integrate - how to operate the build → deploy → observe → adjust
cycle

Nothing here is magic. Everything is **controlled, observable, and
repeatable**.

End of runbook.
