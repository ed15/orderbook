variable "aws_region" {
  description = "Región de AWS para desplegar la infraestructura"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Ambiente de despliegue"
  type        = string
  default     = "test"
}

variable "cluster_name" {
  description = "Nombre del clúster de EKS"
  type        = string
  default     = "orderbook-eks-cluster"
}