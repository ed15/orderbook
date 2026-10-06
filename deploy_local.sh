#!/bin/bash
set -e

echo -e ""

por="$1"

# Construir la imagen del orderbook
echo -e "Contruyendo imagen docker....."
docker build -t orderbook .

# Guardar en imagenes de Kubernetes
echo -e ""
echo -e "Guardando imagen en kubernetes local ....."
docker save orderbook:latest | sudo ctr -n=k8s.io images import -

# Cambiar imagen deploy
sed -i 's/$imagen/orderbook:latest/g' k8s/deploy.yaml

# Aplicar manifiestos de kubernetes (Deployment)
echo -e ""
echo -e "Aplicando deployment App Orderbook....."
kubectl apply -f k8s/deploy.yaml

# Aplicar HPA Kubernetes
echo -e ""
echo -e "Aplicando HPA App Orderbook....."
kubectl apply -f k8s/hpa.yaml

# Cambiar tipo service a ClusterIP
sed -i 's/$type/ClusterIP/g' k8s/svc.yaml
# Cambiar port del servicio
sed -i '/port/s/$port/'$por'/g' k8s/svc.yaml

# Aplicar servicio ClusterIP kubernetes
echo -e ""
echo -e "Aplicando SVC App Orderbook....."
kubectl apply -f k8s/svc.yaml

# Cambiar el servicio como estaba
sed -i 's/ClusterIP/$type/g' k8s/svc.yaml
sed -i '/port/s/'$por'/$port/g' k8s/svc.yaml

# Cambiar imagen deploy como estaba
sed -i 's/orderbook:latest/$imagen/g' k8s/deploy.yaml

#sleep 10
# Port forward servicio kubernetes puerto 8080
#echo -e ""
#echo -e "Port Forward....."
#kubectl port-forward svc/orderbook-service 8080:8080