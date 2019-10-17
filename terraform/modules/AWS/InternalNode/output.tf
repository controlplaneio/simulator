output "InternalNodePrivateIp" {
  value       = "${aws_instance.simulator_bastion.public_ip}"
  description = "Internal node private ip address"
}

