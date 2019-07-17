output "AmiId" {
  value       = "${data.aws_ami.find_ami.id}"
  description = "AMI image id"
}

