module "instances" {
  source = "./modules/instance-group"

  for_each = {
    for index, group in var.instance_groups : index => group
  }

  name                        = format("%s %s", title(var.name), each.value.name)
  group                       = each.value.name
  instance_count              = each.value.count
  ami_id                      = each.value.ami_id
  instance_type               = each.value.instance_type
  key_name                    = aws_key_pair.admin.id
  availability_zone           = var.availability_zone
  subnet_id                   = each.value.public ? var.public_subnet_id : var.private_subnet_id
  associate_public_ip_address = each.value.public
  security_group_id           = aws_security_group.instances.id
  iam_instance_profile        = each.value.iam_instance_profile
  volume_type                 = each.value.volume_type
  volume_size                 = each.value.volume_size
  user_data                   = data.template_file.cloud_config.rendered
  tags = merge(
    var.tags,
    {
      "Name"          = format("%s %s", title(var.name), title(each.value.name))
      "InstanceGroup" = each.value.name
    }
  )
}

resource "aws_security_group" "instances" {
  name   = local.instances_sg_name
  vpc_id = var.network_id

  tags = merge(
    var.tags,
    {
      "Name" = format("%s Instances", title(var.name))
    }
  )
}

resource "aws_security_group_rule" "instances_vpc_ingress" {
  security_group_id = aws_security_group.instances.id
  type              = "ingress"
  protocol          = -1
  from_port         = 0
  to_port           = 0
  cidr_blocks = [
    data.aws_vpc.target.cidr_block,
  ]
}

resource "aws_security_group_rule" "instances_open_egress" {
  security_group_id = aws_security_group.instances.id
  type              = "egress"
  protocol          = -1
  from_port         = 0
  to_port           = 0
  cidr_blocks = [
    "0.0.0.0/0",
  ]
}

data "aws_vpc" "target" {
  id = var.network_id
}
