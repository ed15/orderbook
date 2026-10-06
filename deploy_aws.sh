#!/bin/bash
set -e

port="$1"

# Cambiar port del servicio
sed -i '/port/s/$port/'$port'/g' k8s/svc.yaml
# Cambiar tipo service a ClusterIP
sed -i 's/$type/ClusterIP/g' k8s/svc.yaml

# Directorio Terraform
cd terraform/

# Descarga Modulos
terraform init

# Validar sintaxis
terraform validate

# Ver que se va a crear
terraform plan

# Crear los recursos en AWS
terraform apply -auto-approve

# Configurar acceso a Kubernetes
eval $(terraform output -raw configure_kubectl)

# Verificando nodos
kubectl get nodes

# Aplica manifiestos
kubectl apply -f k8s/

# Espera que los pods esten en running
kubectl get pods -w -o wide

# Cambiar el servicio como estaba
sed -i 's/ClusterIP/$type/g' k8s/svc.yaml
sed -i 's/ClusterIP/$type/g' k8s/svc.yaml



