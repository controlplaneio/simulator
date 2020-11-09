data "template_file" "node_cloud_config" {
  count    = var.number_of_cluster_instances
  template = file("${path.module}/node-cloud-config.yaml")
  vars = {
    hostname               = "k8s-node-${count.index}"
    s3_bucket_name         = var.s3_bucket_name
    github_usernames       = local.access_github_usernames
    node_bashrc            = filebase64("${path.module}/bashrc")
    node_inputrc           = filebase64("${path.module}/inputrc")
    node_aliases           = filebase64("${path.module}/bash_aliases")
    authorized_keys_script = filebase64("${path.module}/authorized_keys.sh")
  }
}

