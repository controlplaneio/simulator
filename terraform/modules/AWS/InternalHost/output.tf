output "InternalHostPrivateIp" {
  value       = "${aws_instance.simulator_internal_host.private_ip}"
  description = "Internal host private ip address"
}

