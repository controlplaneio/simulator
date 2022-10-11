resource "aws_instance" "simulator_internal_host" {
  ami                         = var.ami_id
  key_name                    = var.access_key_name
  instance_type               = var.instance_type
  vpc_security_group_ids      = [var.control_plane_sg_id]
  associate_public_ip_address = false
  subnet_id                   = var.private_subnet_id
  user_data                   = data.template_file.internal_config.rendered
  iam_instance_profile        = var.iam_instance_profile_id
  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
  }

  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Internal Host"
    },
  )
}

