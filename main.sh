#!/bin/bash

#####################################
# SCRIPT PRINCIPAL DE DESPLIEGUE
#####################################
set -e

echo -e ""
echo -e "-------------------------------"
echo -e "      APLICACION ORDERBOOK"
echo -e "-------------------------------"
echo -e ""
echo -e "#################################"
echo -e "# MENU DE DESPLIEGUE, SELECCIONE:"
echo -e "#################################"

echo -e ""

echo -e "1. Local ---> Kubernetes Nativo"
echo -e "2. Nube ---> AWS"
echo -e ""

echo "Digite:"
read result

if [ $result -eq 1 ]; then
    echo "Ejecución Kubernetes nativo."
    echo -e ""
    echo -e "Digite el puerto para el servicio"
    read por    
    
    ./deploy_local.sh $por

elif [ $result -eq 2 ]; then
    echo "AWS, en ejecución..."
    echo "Digite aws_region:"
    read awsregion
    echo "Digite environment:"
    read environment
    echo "Digite el cluster_name:"
    read cluster_name

    printf "%s\n" \
    "aws_region   = \"$awsregion\"" \
    "environment  = \"$environment\"" \
    "cluster_name = \"$cluster_name\"" > terraform/terraform.tfvars

    echo "Digite el puerto para el servicio"
    read port

    perl -pe 'chomp if eof' -i terraform/terraform.tfvars

    # Script Terraform
    ./deploy_aws.sh $port
else
    echo "Opcion incorrecta"
fi
