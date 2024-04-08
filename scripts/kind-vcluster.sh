#!/bin/sh

set -x

# Check requirements
if ! kubectl >/dev/null 2>&1 ; then
  echo "Missing Kubectl binary, please install it from https://kubernetes.io/docs/tasks/tools/"
  exit 1
fi

if ! kind --version >/dev/null 2>&1 ; then
  echo "Missing Kind binary, please install it from https://github.com/kubernetes-sigs/kind"
  exit 1
fi

helm repo add krateo https://charts.krateo.io

helm repo update krateo

kind create cluster \
  --wait 120s \
  --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: krateo-quickstart
nodes:
- role: control-plane
- role: worker
  extraPortMappings:
  - containerPort: 30080 # Krateo Portal
    hostPort: 30080
  - containerPort: 30081 # Krateo BFF
    hostPort: 30081
  - containerPort: 30082 # Krateo AuthN Service
    hostPort: 30082
  - containerPort: 30443 # Krateo Gateway
    hostPort: 30443
  - containerPort: 31443 # vCluster API Server Port
    hostPort: 31443
networking:
  # By default the API server listens on a random open port.
  # You may choose a specific port but probably don't need to in most cases.
  # Using a random port makes it easier to spin up multiple clusters.
  apiServerPort: 6443
EOF

helm upgrade installer installer \
  --repo https://charts.krateo.io \
  --version 0.1.33 \
  --namespace krateo-system \
  --create-namespace \
  --install \
  --wait
