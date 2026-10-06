# APP - Libro de Ofertas - Orderbook

Aplicación desarrolada en GO que implementa un motor de libro de órdenes (Order Book), construido en Docker, orquestado en Kubernetes, desplegado sobre infraestructura Kubernetes, adecuado tambien para desplegar en un cluster AWS EKS con VPC por medio de código IAC con Terraform.

# Estructura del Repositorio

- Dockerfile      # Configuración para construir la imagen del contenedor en Go
- main.sh         # Script orquestador principal
- deploy_local.sh # Script para despliegue local (Docker - Kubernetes)
- deploy_aws.sh   # Script para despliegue en AWS (Terraform + EKS)
- go.mod          # Declaración del módulo y lista de dependencias
- go.sum          # Hashes de seguridad para verificar las dependencias
- README.md       # Documentación y guía del proyecto

# Prerrequisitos

- Go (v1.22+)
- Docker
- kubectl (v1.30+)
- Terraform (v1.5+)
- AWS CLI

# Componentes

# Guía de Despliegue

## 1. Despliegue local (Kubernetes/Docker) Automatico (Script)
Para probar la aplicación en el cluster local kubernetes, realizar los siguientes pasos en una consola de comandos Linux.

- Ejecutar los siguientes comandos:

chmod +x main.sh
chmod +x deploy_local.sh
chmod +x deploy_aws.sh

- Ejecutar el script main.sh para inciar el despliegue:

./main.sh

 MENU DE DESPLIEGUE, SELECCIONE:

1. Local ---> Kubernetes Nativo
2. Nube ---> AWS

Digite:

- Digitar la opción numero 1.
- Luego digitar el puerto para el servicio que sera expuesto.
- Despues se construira la imagen Docker con el nombre (orderbook:lastest).
- Guardara la imagen en las imagenes guardadas de Kubernetes.
- Luego aplicara el deployment de la aplicación.
- Aplicara el HPA (Horizontal pod Autoscaler).
- Luego aplicara el manifiesto del Servicio, para este caso (ClusterIP).
- Por ultimo dejara los manifiestos de deploy y svc como estaban antes de ejecutar el script main.sh.

Para hacer el port forward tener el cuenta el puerto digitado y ejecutar el siguiente comando:

kubectl port-forward svc/orderbook-service 8080:$puerto_digitado


## 2. Despliegue en aws (EKS - Terraform) Automatico (Script):

Para probar la aplicación en EKS de AWS, realizar los siguientes pasos en una consola de comandos Linux.

- Ejecutar los siguientes comandos:

chmod +x main.sh
chmod +x deploy_local.sh
chmod +x deploy_aws.sh

- Ejecutar el script main.sh para inciar el despliegue:

./main.sh

 MENU DE DESPLIEGUE, SELECCIONE:


1. Local ---> Kubernetes Nativo
2. Nube ---> AWS

Digite:

- Digitar la opción numero 2.
- Luego digitar la región de AWS para desplegar la infraestructura.
- Luego digitar el environment.
- Despues ejecutar el nombre del cluster.
- Digitar el puerto para el servicio.
- Se aplicaran los manifiestos de kubernetes.
- Luego aplicara el manifiesto del Servicio, para este caso (ClusterIP).
- Por ultimo dejara los manifiestos de deploy y svc como estaban antes de ejecutar el script main.sh.

- Para hacer el port forward tener el cuenta el puerto digitado y ejecutar el siguiente comando:

kubectl port-forward svc/orderbook-service 8080:$puerto_digitado

## 3. Despliegue manual en entorno local (Kubernetes/Docker)

Paso 1: Construir la imagen de Docker y subirla algun registro privado o publico
docker build -t orderbook:latest .
docker tag orderbook:latest mi-registro:puerto/orderbook:latest

Paso 2: Modificar los manifiestos de kubernetes
deploy.yaml ------> Cambiar $image a la dirección de la imagen donde quedo guardada en el paso anterior.
svc.yanl ---------> Cambiar $type al tipo de servicio a crear y $port al puerto del servicio.

Paso 3: Desplegar los manifiestos de kubernetes.
kubectl apply -f k8s/

Paso 4: Ejecutar el comando para hacer forward del puerto del servicio
kubectl port-forward svc/orderbook-service 8080:$puerto

## 4. Despliegue manual en entorno AWS (EKS)

Paso 1: Construir la imagen de Docker y subirla algun registro privado o publico
docker build -t orderbook:latest .
docker tag orderbook:latest mi-registro:puerto/orderbook:latest

Paso 2: Modificar los manifiestos de kubernetes
deploy.yaml ------> Cambiar $image a la dirección de la imagen donde quedo guardada en el paso anterior.
svc.yanl ---------> Cambiar $type al tipo de servicio a crear y $port al puerto del servicio.

Paso 3: Ajustar archivo terraform/terraform.tfvars

Paso 4: Inicializar módulos de Terraform
cd terraform
terraform init

Paso 5: Planificar la infraestructura
terraform plan

Paso 6: Crear la VPC y cluster EKS
terraform apply -auto-approve

Paso 7: conectar kubectl con el cluster de AWS
eval $(terraform output -raw configure_kubectl)

Paso 8: Ver nodos EC2
kubectl get nodes

Paso 9: Desplegar los manifiestos de kubernetes
cd ..
kubectl apply -f k8s/

Paso 10: Consultar los pods y el servicio
kubectl get pods
kubectl get svc orderbook-service

# Guia peticiones API

Para consultar la ip: (kubectl get svc)
localhost = ip del servicio

Para consultar las billeteras:
curl -s http://localhost:8080/wallets | jq .

{
  "user_buyer": {
    "user_id": "user_buyer",
    "brl_available": 1000000,
    "brl_locked": 0,
    "vibranium_available": 0,
    "vibranium_locked": 0
  },
  "user_seller": {
    "user_id": "user_seller",
    "brl_available": 0,
    "brl_locked": 0,
    "vibranium_available": 10000,
    "vibranium_locked": 0
  }
}

Para enviar orden de venta para el usuario user_seller:

curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user_seller", "type": "SELL", "price": 100.0, "amount": 10.0}'

Para enviar orden de compra desde el usuario user_buyer:

curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user_buyer", "type": "BUY", "price": 100.0, "amount": 10.0}'

Para consultar las ultima transaccion:

curl -s http://10.101.173.18:8080/trades | jq '.[-1:]'

Para consultar las dos ultimas transacciones:

curl -s http://10.101.173.18:8080/trades | jq '.[-2:]'

Para consultar las metricas:

curl -s http://localhost:8080/metrics

Para consultar el endpoint utilizado por kubernetes:

curl -s http://localhost:8080/health


