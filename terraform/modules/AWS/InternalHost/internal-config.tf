locals {
  access_github_usernames = "${join(" ", var.access_github_usernames)}"
}

data "template_file" "internal_config" {
  template = file("${path.module}/internal-config.yaml")
  vars = {
    s3_bucket_name = var.s3_bucket_name
    github_usernames       = local.access_github_usernames
    host_bashrc            = filebase64("${path.module}/bashrc")
    host_inputrc           = filebase64("${path.module}/inputrc")
    host_aliases           = filebase64("${path.module}/bash_aliases")
    authorized_keys_script = filebase64("${path.module}/authorized_keys.sh")
  }
}

