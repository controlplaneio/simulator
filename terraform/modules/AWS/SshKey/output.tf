output "KeyPairName" {
  value       = "${aws_key_pair.simulator_bastion_key.key_name}"
  description = "Name of Bastion SSH key in KMS"
}
