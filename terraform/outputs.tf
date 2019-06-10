output "bastion_public_ip" {
  value = "${aws_instance.bastion.public_ip}"
}
output "master_nodes_private_ip" {
  value = "${aws_instance.controlplane.*.private_ip}"
}
output "cluster_nodes_private_ip" {
  value = "${aws_instance.cluster_nodes.*.private_ip}"
}
