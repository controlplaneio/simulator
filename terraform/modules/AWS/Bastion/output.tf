output "BastionPublicIp" {
  value       = "${aws_instance.simulator_bastion.public_ip}"
  description = "Bastion server public ip address"
}

