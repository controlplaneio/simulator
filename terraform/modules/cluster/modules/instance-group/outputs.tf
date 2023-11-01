output "instances" {
  value = zipmap(aws_instance.instance.*.tags.ID, aws_instance.instance.*.private_ip)
}

