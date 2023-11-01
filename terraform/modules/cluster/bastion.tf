resource "aws_instance" "bastion" {
  ami                         = var.bastion_ami_id
  instance_type               = var.bastion_instance_type
  key_name                    = aws_key_pair.admin.id
  subnet_id                   = var.public_subnet_id
  associate_public_ip_address = true
  vpc_security_group_ids = [
    aws_security_group.bastion.id,
  ]

  #  metadata_options {
  #    http_endpoint = "disabled"
  #  }

  user_data = <<EOF
${data.template_file.cloud_config.rendered}
hostname: bastion
EOF

  root_block_device {
    volume_type = var.bastion_volume_type
    volume_size = var.bastion_volume_size
  }

  tags = merge(
    var.tags,
    {
      "Name" = format("%s Bastion", title(var.name))
    }
  )

  volume_tags = merge(
    var.tags,
    {
      "Name" = format("%s Bastion", title(var.name))
    }
  )
}

resource "aws_security_group" "bastion" {
  name   = local.bastion_sg_name
  vpc_id = var.network_id

  tags = merge(
    var.tags,
    {
      "Name" = format("%s Bastion", title(var.name))
    }
  )
}

resource "aws_security_group_rule" "bastion_ssh_ingress" {
  security_group_id = aws_security_group.bastion.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 22
  to_port           = 22
  cidr_blocks = [
    format("%s/32", trim(data.http.player_ip.response_body, "\n")),
    #    "0.0.0.0/0",
  ]
}

resource "aws_security_group_rule" "bastion_open_egress" {
  security_group_id = aws_security_group.bastion.id
  type              = "egress"
  protocol          = -1
  from_port         = 0
  to_port           = 0
  cidr_blocks = [
    "0.0.0.0/0",
  ]
}

data "http" "player_ip" {
  url = "https://icanhazip.com/"
}
