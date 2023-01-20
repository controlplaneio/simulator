resource "aws_instance" "simulator_master_instances" {
  count                       = var.number_of_master_instances
  ami                         = var.ami_id
  key_name                    = var.access_key_name
  instance_type               = var.master_instance_type
  vpc_security_group_ids      = [var.control_plane_sg_id]
  associate_public_ip_address = false
  subnet_id                   = var.private_subnet_id
  user_data = element(
    data.template_file.master_cloud_config.*.rendered,
    count.index,
  )
  iam_instance_profile = var.iam_instance_profile_id
  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
  }

  root_block_device {
    volume_size = 20
  }

  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Kubernetes Master"
    },
  )
}

resource "aws_instance" "simulator_node_instances" {
  count                       = var.number_of_cluster_instances
  ami                         = var.ami_id
  key_name                    = var.access_key_name
  instance_type               = var.cluster_nodes_instance_type
  vpc_security_group_ids      = [var.control_plane_sg_id]
  associate_public_ip_address = false
  subnet_id                   = var.private_subnet_id
  user_data = element(
    data.template_file.node_cloud_config.*.rendered,
    count.index
  )
  iam_instance_profile = var.iam_instance_profile_id
  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
  }

  root_block_device {
    volume_size = 20
  }

  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Kubernetes Node"
    },
  )
}

locals {
  access_github_usernames = join(" ", var.access_github_usernames)
  version_minor           = split(".", var.kubernetes_version)[1]
  version_major_minor     = join(".", ["1", local.version_minor])
}
