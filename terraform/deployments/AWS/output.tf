output "bastion_public_ip" {
  value       = "${module.Bastion.BastionPublicIp}"
  description = "Bastion public IP"
}
output "master_nodes_private_ip" {
  value       = "${module.Kubernetes.K8sMasterPrivateIp}"
  description = "Master node private IP"
}
output "cluster_nodes_private_ip" {
  value       = "${module.Kubernetes.K8sNodesPrivateIp}"
  description = "Cluster node private IPs"
}
output "internal_node_private_ip" {
  value       = "${module.InternalNode.InternalNodePrivateIp}"
  description = "Private Subnet node"
}
output "access_cidr" {
  value       = "${var.access_cidr}"
  description = "Remote access IP"
}
output "ami_id" {
  value       = "${module.Ami.AmiId}"
  description = "AMI used for all instances"
}

