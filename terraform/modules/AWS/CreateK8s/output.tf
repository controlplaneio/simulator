output "K8sMasterPrivateIp" {
  value       = "${aws_instance.simulator_controlplane_instances.*.private_ip}"
  description = "Kubernetes master private ip address"
}
output "K8sNodesPrivateIp" {
  value       = "${aws_instance.simulator_cluster_instances.*.private_ip}"
  description = "Kubernetes node(s) private ip address"
}

