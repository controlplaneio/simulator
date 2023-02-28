resource "aws_instance" "simulator_bastion" {
  ami                         = var.ami_id
  key_name                    = var.access_key_name
  instance_type               = var.instance_type
  vpc_security_group_ids      = [var.security_group]
  associate_public_ip_address = true
  subnet_id                   = var.subnet_id
  user_data                   = data.cloudinit_config.cc.rendered
  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
  }

  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Bastion"
    },
  )

}
