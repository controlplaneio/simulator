output "AmiId" {
  value       = "${aws_ami.find_ami.id}"
  description = "AMI image id"
}

