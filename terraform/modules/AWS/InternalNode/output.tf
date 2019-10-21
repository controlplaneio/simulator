output "InternalNodePrivateIp" {
  value       = "${aws_instance.simulator_internal_node.private_ip}"
  description = "Internal node private ip address"
}

