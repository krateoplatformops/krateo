#!/bin/sh

set -x

# Check requirements
if ! kubectl >/dev/null 2>&1 ; then
  echo "Missing Kubectl binary, please install it from https://kubernetes.io/docs/tasks/tools/"
  exit 1
fi

helm repo add krateo https://charts.krateo.io

helm repo update krateo

helm upgrade installer installer \
  --repo https://charts.krateo.io \
  --namespace krateo-system \
  --create-namespace \
  --set krateoplatformops.service.type=LoadBalancer \
  --set krateoplatformops.service.externalIpAvailable=true \
  --install \
  --wait
