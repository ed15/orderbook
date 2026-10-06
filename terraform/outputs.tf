output "cluster_endpoint" {
  description = "Endpoint del API Server de EKS"
  value       = module.eks.cluster_endpoint
}

output "cluster_name" {
  description = "Nombre del clúster EKS"
  value       = module.eks.cluster_name
}

output "configure_kubectl" {
  description = "Comando para conectar kubectl con AWS"
  value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${module.eks.cluster_name}"
}