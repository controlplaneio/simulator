resource "aws_instance" "simulator_master_instances" {
  count                       = var.number_of_master_instances
  ami                         = var.ami_id
  key_name                    = var.access_key_name
  instance_type               = var.master_instance_type
  vpc_security_group_ids      = [var.control_plane_sg_id]
  associate_public_ip_address = false
  subnet_id                   = var.private_subnet_id
  user_data = templatefile(
    "${path.module}/master-cloud-config.yaml",
    {
      hostname       = "k8s-master-${count.index}"
      s3_bucket_name = var.s3_bucket_name
      master_bashrc  = filebase64("${path.module}/bashrc")
      master_inputrc = filebase64("${path.module}/inputrc")
      master_aliases = filebase64("${path.module}/bash_aliases")
      version        = var.kubernetes_version
    }
  )
  iam_instance_profile = var.iam_instance_profile_id
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
  user_data = templatefile(
    "${path.module}/node-cloud-config.yaml",
    {
      hostname       = "k8s-node-${count.index}"
      s3_bucket_name = var.s3_bucket_name
      node_bashrc    = filebase64("${path.module}/bashrc")
      node_inputrc   = filebase64("${path.module}/inputrc")
      node_aliases   = filebase64("${path.module}/bash_aliases")
      version        = var.kubernetes_version
    }
  )
  iam_instance_profile = var.iam_instance_profile_id
  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Kubernetes Node"
    },
  )
}
