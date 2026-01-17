# k8s-poc -- Runbook de Execução e Operação (Versão Aprimorada)

Este repositório é uma **POC educacional** para aprender fundamentos de
Kubernetes utilizando uma aplicação Go (Gin), exposta via Service, com
métricas coletadas por Prometheus e visualizadas no Grafana.

O foco **não é produção**, e sim **entendimento profundo e consciente**
do fluxo completo:

-   build da aplicação
-   criação da imagem Docker
-   deploy no Kubernetes
-   exposição de serviços
-   observabilidade básica
-   operação, troubleshooting e reset do ambiente

Nada aqui é implícito. Tudo é **intencionalmente explícito** para evitar
modelos mentais errados.

------------------------------------------------------------------------

## Visão Geral da Arquitetura

    [ Usuário / curl ]
            |
         NodePort
            |
         Service
            |
       Pods (3 réplicas)
            |
         /metrics
            |
       Prometheus
            |
         Grafana

------------------------------------------------------------------------

## Estrutura do Projeto

    .
    ├── cmd/
    │   └── server/                # entrypoint da aplicação Go
    ├── internal/                  # código da aplicação
    ├── k8s/
    │   ├── kustomization.yaml     # aplica todos os manifests
    │   ├── namespaces.yaml        # namespace monitoring
    │   ├── app/
    │   │   ├── configmap.yaml     # variáveis de ambiente da app
    │   │   ├── deployment.yaml    # deployment da aplicação (3 réplicas)
    │   │   └── service.yaml       # Service (NodePort)
    │   ├── prometheus/
    │   │   ├── configmap.yaml     # config do Prometheus (scrape)
    │   │   ├── deployment.yaml    # Prometheus
    │   │   └── service.yaml       # Service do Prometheus
    │   └── grafana/
    │       ├── configmap-datasource.yaml # datasource Prometheus
    │       ├── deployment.yaml           # Grafana
    │       └── service.yaml              # Service do Grafana
    ├── Dockerfile
    └── go.mod

------------------------------------------------------------------------

## O Que Foi Implementado

### Aplicação

-   Aplicação Go usando Gin
-   Endpoints:
    -   `/health`: health check + métricas
    -   `/whoami`: informações do Pod (Downward API) + ambiente
        (ConfigMap)
    -   `/metrics`: métricas Prometheus

### Kubernetes

-   Deployment com **3 réplicas**
-   Service **NodePort** para acesso externo previsível
-   ConfigMap para configuração sem rebuild
-   Downward API para identidade do Pod
-   Probes de liveness e readiness

### Observabilidade

-   Prometheus raspando métricas via Service
-   Grafana com datasource provisionado automaticamente
-   Dashboard validando:
    -   taxa de health checks
    -   latência
    -   distribuição por Pod

------------------------------------------------------------------------

## Função de Cada Arquivo (Kubernetes)

### `k8s/app/configmap.yaml`

-   Centraliza configurações da aplicação
-   Alterações **não exigem rebuild da imagem**
-   Exemplo: `ENVIRONMENT`, usado no `/whoami`

### `k8s/app/deployment.yaml`

Define: - imagem Docker da aplicação - número de réplicas - injeção de
envs via ConfigMap e Downward API - políticas de health
(liveness/readiness)

### `k8s/app/service.yaml`

-   Expõe a aplicação via **NodePort**
-   Permite acesso direto em `http://localhost:30080`
-   Evita uso contínuo de `kubectl port-forward`

------------------------------------------------------------------------

### `k8s/prometheus/configmap.yaml`

-   Define scrape do endpoint `/metrics`
-   Aponta para o Service da aplicação, não para Pods

### `k8s/prometheus/deployment.yaml`

-   Sobe o Prometheus
-   Monta configuração via ConfigMap

### `k8s/prometheus/service.yaml`

-   Expõe Prometheus apenas internamente (ClusterIP)

------------------------------------------------------------------------

### `k8s/grafana/configmap-datasource.yaml`

-   Provisiona o datasource Prometheus automaticamente
-   Elimina configuração manual via UI

### `k8s/grafana/deployment.yaml`

-   Sobe o Grafana
-   Credenciais default apenas para POC

### `k8s/grafana/service.yaml`

-   Pode ser acessado via port-forward ou NodePort (consciente)

------------------------------------------------------------------------

## Pré-requisitos

-   Docker Desktop com Kubernetes habilitado
-   Docker CLI
-   `kubectl`
-   Go 1.25

------------------------------------------------------------------------

## Build da Aplicação

Sempre que **o código Go mudar**:

``` bash
go mod tidy
docker build -t k8s-poc:latest .
```

> Kubernetes **não recompila código**. Ele apenas executa imagens.

------------------------------------------------------------------------

## Subir Tudo do Zero

``` bash
kubectl apply -k k8s/
```

Validação:

``` bash
kubectl get pods
kubectl get pods -n monitoring
kubectl get svc
```

------------------------------------------------------------------------

## Como Acessar

### Aplicação

    http://localhost:30080/health
    http://localhost:30080/whoami
    http://localhost:30080/metrics

### Prometheus

``` bash
kubectl port-forward -n monitoring svc/prometheus 9090:9090
```

Acesse: http://localhost:9090

### Grafana

``` bash
kubectl port-forward -n monitoring svc/grafana 3000:3000
```

Acesse: http://localhost:3000\
Login: `admin / admin`

------------------------------------------------------------------------

## Como Descobrir Portas Expostas

``` bash
kubectl get svc
kubectl get svc --all-namespaces
```

Exemplo:

    80:30080/TCP

-   80 → porta interna do Service
-   30080 → porta exposta no Node

Para detalhes:

``` bash
kubectl describe svc k8s-poc
```

------------------------------------------------------------------------

## Quando Algo Mudar, O Que Fazer

### Código da aplicação

``` bash
docker build -t k8s-poc:latest .
kubectl rollout restart deployment k8s-poc
```

### ConfigMap

``` bash
kubectl apply -f k8s/app/configmap.yaml
kubectl rollout restart deployment k8s-poc
```

### YAMLs

``` bash
kubectl apply -f <arquivo>
```

------------------------------------------------------------------------

## Como Derrubar Tudo

### Apenas a aplicação

``` bash
kubectl delete deployment k8s-poc
kubectl delete svc k8s-poc
kubectl delete configmap k8s-poc-config
```

### Observabilidade

``` bash
kubectl delete namespace monitoring
```

### Reset completo

``` bash
kubectl delete -k k8s/
```

------------------------------------------------------------------------

## Comandos Úteis de Debug

``` bash
kubectl get pods -w
kubectl describe pod <nome>
kubectl logs <pod>
kubectl get svc --all-namespaces
```

------------------------------------------------------------------------

## Observações Importantes

-   NodePort é **global no cluster**
-   Port-forward é ferramenta de **debug**
-   Métricas com alta cardinalidade são didáticas
-   Resetar namespace em POC é normal

------------------------------------------------------------------------

## Objetivo Final da POC

Aprender, de forma consciente: - como Kubernetes gerencia Pods - como
Services distribuem tráfego - como métricas são coletadas - como
Prometheus e Grafana se conectam - como operar o ciclo build → deploy →
observar → ajustar

Nada aqui é mágico. Tudo é **controlado, observável e repetível**.

Fim do runbook.
