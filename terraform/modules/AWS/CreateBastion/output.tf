output "KeyPairName" {
  value       = "${aws_key_pair.simulator_bastion_key.key_name}"
  description = "Name of Bastion SSH key in KMS"
}
output "BastionPublicIp" {
  value       = "${aws_instance.simulator_bastion.public_ip}"
  description = "Bastion server public ip address"
}

