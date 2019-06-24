output "KeyPairName" {
  value       = "${aws_key_pair.bastion_key.key_name}"
  description = "Name of Bastion SSH key in KMS"
}
output "BastionPublicIp" {
  value       = "${aws_instance.bastion.public_ip}"
  description = "Bastion server public ip address"
}

