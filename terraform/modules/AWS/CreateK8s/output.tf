output "K8sMasterPrivateIp" {
  value       = "${aws_instance.controlplane.*.private_ip}"
  description = "Kubernetes master private ip address"
}
output "K8sNodesPrivateIp" {
  value       = "${aws_instance.cluster_nodes.*.private_ip}"
  description = "Kubernetes node(s) private ip address"
}

