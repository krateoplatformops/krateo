#!/bin/sh

set -x

helm repo add krateo https://charts.krateo.io

helm repo update krateo

helm repo add loft-sh https://charts.loft.sh

helm repo update loft-sh

helm upgrade krateo-vcluster loft-sh/vcluster-k8s \
  --version 0.19.4 \
  --namespace krateo-system \
  --create-namespace \
  --set service.type=LoadBalancer \
  --install \
  --wait

export KUBECONFIG_CACRT=$(kubectl get secret krateo-vcluster-certs -n krateo-system -o json | jq -r '.data."ca.crt"')

export KUBECONFIG_CAKEY=$(kubectl get secret krateo-vcluster-certs -n krateo-system -o json | jq -r '.data."ca.key"')

export KUBECONFIG_KUBERNETES_IP=$(kubectl get svc krateo-vcluster-lb -n krateo-system -o json | jq -r '.status.loadBalancer.ingress[0].ip')

curl -L -o vcluster "https://github.com/loft-sh/vcluster/releases/latest/download/vcluster-darwin-amd64" && chmod +x vcluster

./vcluster connect krateo-vcluster

kubectl create ns krateo-system

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: krateo-gateway
  namespace: krateo-system
type: Opaque
stringData:
  KRATEO_GATEWAY_CAKEY: $KUBECONFIG_CAKEY
EOF

helm upgrade krateo-gateway krateo-gateway \
  --repo https://charts.krateo.io \
  --version 0.3.12 \
  --namespace krateo-system \
  --create-namespace \
  --set service.type=LoadBalancer \
  --set livenessProbe=null \
  --set readinessProbe=null \
  --set env.KRATEO_GATEWAY_CACRT=$KUBECONFIG_CACRT \
  --set env.KRATEO_BFF_SERVER=http://krateo-bff.krateo-system.svc:8081 \
  --set env.KRATEO_GATEWAY_DEBUG=true \
  --set env.KRATEO_GATEWAY_DUMP_ENV=true \
  --install \
  --wait

export KRATEO_GATEWAY_LOADBALANCER_IP=$(kubectl get svc krateo-gateway -n krateo-system -o json | jq -r '.status.loadBalancer.ingress[0].ip')

helm upgrade krateo-gateway krateo-gateway \
  --repo https://charts.krateo.io \
  --version 0.3.12 \
  --namespace krateo-system \
  --create-namespace \
  --set service.type=LoadBalancer \
  --set livenessProbe=null \
  --set readinessProbe=null \
  --set env.KRATEO_GATEWAY_CACRT=$KUBECONFIG_CACRT \
  --set env.KRATEO_BFF_SERVER=http://krateo-bff.krateo-system.svc:8081 \
  --set env.KRATEO_GATEWAY_IP_ADDRESSES=$KRATEO_GATEWAY_LOADBALANCER_IP \
  --set env.KRATEO_GATEWAY_DEBUG=true \
  --set env.KRATEO_GATEWAY_DUMP_ENV=true \
  --install \
  --wait

helm upgrade authn-service authn-service \
  --repo https://charts.krateo.io \
  --version 0.10.1 \
  --namespace krateo-system \
  --create-namespace \
  --set service.type=LoadBalancer \
  --set env.AUTHN_CORS=true \
  --set env.AUTHN_KUBERNETES_URL=https://$KUBECONFIG_KUBERNETES_IP \
  --set env.AUTHN_KUBECONFIG_PROXY_URL=http://$KRATEO_GATEWAY_LOADBALANCER_IP:8443 \
  --set env.AUTHN_KUBECONFIG_CACRT=$KUBECONFIG_CACRT \
  --set env.AUTHN_DUMP_ENV=true \
  --install \
  --wait

export AUTHN_SERVICE_LOADBALANCER_IP=$(kubectl get svc authn-service -n krateo-system -o json | jq -r '.status.loadBalancer.ingress[0].ip')

helm upgrade krateo-bff krateo-bff \
  --repo https://charts.krateo.io \
  --version 0.14.3 \
  --namespace krateo-system \
  --create-namespace \
  --set service.type=LoadBalancer \
  --set env.KRATEO_BFF_CORS=true \
  --set env.KRATEO_BFF_DUMP_ENV=true \
  --set env.KRATEO_BFF_DEBUG=true \
  --install \
  --wait

export KRATEO_BFF_LOADBALANCER_IP=$(kubectl get svc krateo-bff -n krateo-system -o json | jq -r '.status.loadBalancer.ingress[0].ip')

helm upgrade krateo-frontend krateo-frontend \
  --repo https://charts.krateo.io \
  --version 2.0.6 \
  --namespace krateo-system \
  --create-namespace \
  --set service.type=LoadBalancer \
  --set env.AUTHN_API_BASE_URL=http://$AUTHN_SERVICE_LOADBALANCER_IP:8082 \
  --set env.BFF_API_BASE_URL=http://$KRATEO_BFF_LOADBALANCER_IP:8081 \
  --install \
  --wait

helm install core-provider core-provider \
  --repo https://charts.krateo.io \
  --version 0.9.0 \
  --namespace krateo-system \
  --create-namespace \
  --wait

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
type: kubernetes.io/basic-auth
metadata:
  name: cyberjoker-password
  namespace: krateo-system
stringData:
  password: "123456"
---
apiVersion: basic.authn.krateo.io/v1alpha1
kind: User
metadata:
  name: cyberjoker
  namespace: krateo-system
spec:
  displayName: Cyber Joker
  avatarURL: https://i.pravatar.cc/256?img=70
  groups:
    - devs
  passwordRef:
    namespace: krateo-system
    name: cyberjoker-password
    key: password
---
apiVersion: v1
kind: Namespace
metadata:
  name: demo-system
---
apiVersion: layout.ui.krateo.io/v1alpha1
kind: Row
metadata:
  name: two
  namespace: demo-system
spec:
  columnListRef:
    - name: three
      namespace: demo-system
    - name: eleven
      namespace: demo-system
---
apiVersion: layout.ui.krateo.io/v1alpha1
kind: Column
metadata:
  name: three
  namespace: demo-system
spec:
  app:
    props:
      width: "12"
  cardTemplateListRef:
    - name: three
      namespace: demo-system
---
apiVersion: layout.ui.krateo.io/v1alpha1
kind: Column
metadata:
  name: eleven
  namespace: demo-system
spec:
  app:
    props:
      width: "12"
  cardTemplateListRef:
    - name: one
      namespace: demo-system
    - name: ten
      namespace: demo-system
---
apiVersion: widgets.ui.krateo.io/v1alpha1
kind: CardTemplate
metadata:
  name: one
  namespace: demo-system
spec:
  app:
    icon: fa-solid fa-truck-fast
    color: green
    title: \${ .api2.items[0] | (.name  + " -> " + .email) }
    content: \${ .api2.items[0].body }
    date: Sep 15th 2023 08:15:43
    actions:
    - name: remove
      verb: DELETE
      endpointRef:
        name: typicode-endpoint
        namespace: demo-system
      path: \${ "/todos/1/comments/" + (.api2.items[0].id|tostring) }
  api:
  - name: api1
    path: "/todos/1"
    endpointRef:
      name: typicode-endpoint
      namespace: demo-system
    verb: GET
    headers:
    - 'Accept: application/json'
  - name: api2
    dependOn: api1
    path: \${ "/todos/" + (.api1.id|tostring) +  "/comments" }
    endpointRef:
      name: typicode-endpoint
      namespace: demo-system
    verb: GET
    headers:
    - 'Accept: application/json'
---
apiVersion: widgets.ui.krateo.io/v1alpha1
kind: CardTemplate
metadata:
  name: three
  namespace: demo-system
spec:
  iterator: .api1.products[:3]
  app:
    icon: fa-solid fa-mobile-button
    color: blue
    title: \${ .title }
    content: \${ .description }
    tags: \${ .brand }
  api:
  - name: api1
    path: "/products"
    endpointRef:
      name: dummyjson-endpoint
      namespace: demo-system
    verb: GET
    headers:
    - 'Accept: application/json'
---
apiVersion: widgets.ui.krateo.io/v1alpha1
kind: CardTemplate
metadata:
  name: ten
  namespace: demo-system
spec:
  iterator: .api2.items[:10]
  app:
    icon: \${ "fa-solid fa-" + (.id|tostring)}
    color: darkBlue
    title: \${ .name }
    content: \${ .body }
    tags: \${ .email }
    actions:
    - name: view
      endpointRef:
        name: typicode-endpoint
        namespace: demo-system
      path: \${ "/todos/1/comments/" + (.id|tostring) }
  api:
  - name: api1
    path: "/todos/1"
    endpointRef:
      name: typicode-endpoint
      namespace: demo-system
    verb: GET
    headers:
    - 'Accept: application/json'
  - name: api2
    dependOn: api1
    path: \${ "/todos/" + (.api1.id|tostring) +  "/comments" }
    endpointRef:
      name: typicode-endpoint
      namespace: demo-system
    verb: GET
    headers:
    - 'Accept: application/json'
---
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: typicode-endpoint
  namespace: demo-system
stringData:
  server-url: https://jsonplaceholder.typicode.com
---
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: dummyjson-endpoint
  namespace: demo-system
stringData:
  server-url: https://dummyjson.com
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: dev-get-any-layout-in-demosystem-namespace
  namespace: demo-system
rules:
- apiGroups:
  - layout.ui.krateo.io
  resources:
  - rows
  - columns
  verbs:
  - get
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: dev-get-any-layout-in-demosystem-namespace
  namespace: demo-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: dev-get-any-layout-in-demosystem-namespace
subjects:
- kind: Group
  name: devs
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: dev-get-cardtemplate-one-three-ten-in-demosystem-namespace
  namespace: demo-system
rules:
- apiGroups:
  - widgets.ui.krateo.io
  resources:
  - cardtemplates
  resourceNames:
  - one
  - three
  - ten
  verbs:
  - get
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: dev-get-cardtemplate-one-three-ten-in-demosystem-namespace
  namespace: demo-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: dev-get-cardtemplate-one-three-ten-in-demosystem-namespace
subjects:
- kind: Group
  name: devs
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: dev-delete-cardtemplate-one-in-demosystem-namespace
  namespace: demo-system
rules:
- apiGroups:
  - widgets.ui.krateo.io
  resources:
  - cardtemplates
  resourceNames:
  - one
  verbs:
  - delete
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: dev-delete-cardtemplate-one-in-demosystem-namespace
  namespace: demo-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: dev-delete-cardtemplate-one-in-demosystem-namespace
subjects:
- kind: Group
  name: devs
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: core.krateo.io/v1alpha1
kind: SchemaDefinition
metadata:
  annotations:
     "krateo.io/connector-verbose": "true"
  name: fireworksapp
  namespace: demo-system
spec:
  schema:
    version: v1alpha1
    kind: Fireworksapp
    url: https://raw.githubusercontent.com/krateoplatformops/krateo-v2-template-fireworksapp/main/chart/values.schema.json
---
apiVersion: widgets.ui.krateo.io/v1alpha1
kind: FormTemplate
metadata:
  name: fireworksapp
  namespace: demo-system
spec:
  schemaDefinitionRef:
    name: fireworksapp
    namespace: demo-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: apps-viewer
rules:
- apiGroups:
  - apps.krateo.io
  resources:
  - '*'
  verbs:
  - get
  - list
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: apps-viewer
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name:  apps-viewer
subjects:
- kind: Group
  name: devs
  apiGroup: rbac.authorization.k8s.io
EOF

curl http://$AUTHN_SERVICE_LOADBALANCER_IP:8082/basic/login -H "Authorization: Basic Y3liZXJqb2tlcjoxMjM0NTY=" | jq -r .data > cyberjoker.kubeconfig
