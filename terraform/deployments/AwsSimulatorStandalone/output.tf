output "bastion_public_ip" {
  value = "${module.Bastion.BastionPublicIp}"
}
output "master_nodes_private_ip" {
  value = "${module.Kubernetes.K8sMasterPrivateIp}"
}
output "cluster_nodes_private_ip" {
  value = "${module.Kubernetes.K8sNodesPrivateIp}"
}
output "access_cidr" {
  value = "${var.access_cidr}"
}
