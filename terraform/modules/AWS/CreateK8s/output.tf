output "K8sMasterPrivateIp" {
  value = "${aws_instance.controlplane.*.private_ip}"
}
output "K8sNodesPrivateIp" {
  value = "${aws_instance.cluster_nodes.*.private_ip}"
}

