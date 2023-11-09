resource "aws_instance" "instance" {
  count = var.instance_count

  ami                         = var.ami_id
  instance_type               = var.instance_type
  key_name                    = var.key_name
  availability_zone           = var.availability_zone
  subnet_id                   = var.subnet_id
  associate_public_ip_address = var.associate_public_ip_address
  hibernation                 = true
  vpc_security_group_ids      = [
    var.security_group_id,
  ]

  #  metadata_options {
  #    http_endpoint = local.http_endpoint
  #  }

  iam_instance_profile = var.iam_instance_profile

  root_block_device {
    volume_type = var.volume_type
    volume_size = var.volume_size
    encrypted   = true
  }

  user_data = <<EOF
${var.user_data}
hostname: ${format("%s-%s", var.group, count.index + 1)}
EOF

  tags = merge(
    var.tags,
    {
      "ID" : format("%s-%s", var.group, count.index + 1)
    }
  )
  volume_tags = merge(
    var.tags,
    {
      "ID" : format("%s-%s", var.group, count.index + 1)
    }
  )
}
