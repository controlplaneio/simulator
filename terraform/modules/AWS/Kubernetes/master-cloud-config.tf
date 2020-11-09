data "template_file" "master_cloud_config" {
  count    = var.number_of_master_instances
  template = file("${path.module}/master-cloud-config.yaml")
  vars = {
    hostname               = "k8s-master-${count.index}"
    s3_bucket_name         = var.s3_bucket_name
    github_usernames       = local.access_github_usernames
    master_bashrc          = filebase64("${path.module}/bashrc")
    master_inputrc         = filebase64("${path.module}/inputrc")
    master_aliases         = filebase64("${path.module}/bash_aliases")
    authorized_keys_script = filebase64("${path.module}/authorized_keys.sh")
  }
}

