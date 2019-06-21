output "KeyPairName" {
  value = "${aws_key_pair.bastion_key.key_name}"
}
output "BastionPublicIp" {
  value = "${aws_instance.bastion.public_ip}"
}

