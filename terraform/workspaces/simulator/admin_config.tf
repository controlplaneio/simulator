resource "local_file" "admin_private_key" {
  content_base64  = module.cluster.admin_private_key
  filename        = format("%s/%s", var.admin_config_dir, local.ssh_identity_filename)
  file_permission = "0600"
}

resource "aws_s3_object" "admin_private_key" {
  content_base64 = module.cluster.admin_private_key
  bucket         = var.bucket
  key            = format("config/admin/%s", local.ssh_identity_filename)
}

resource "local_file" "admin_ssh_config" {
  content_base64  = module.cluster.admin_ssh_config
  filename        = format("%s/%s", var.admin_config_dir, local.ssh_config_filename)
  file_permission = "0600"
}

resource "aws_s3_object" "admin_ssh_config" {
  content_base64 = module.cluster.admin_ssh_config
  bucket         = var.bucket
  key            = format("config/admin/%s", local.ssh_config_filename)
}

resource "local_file" "ansible_config" {
  content_base64 = module.cluster.ansible_config
  filename       = format("%s/%s", var.admin_config_dir, local.ansible_config_filename)
}

resource "aws_s3_object" "ansible_config" {
  content_base64 = module.cluster.ansible_config
  bucket         = var.bucket
  key            = format("config/admin/%s", local.ansible_config_filename)
}

resource "local_file" "ansible_inventory" {
  content_base64 = module.cluster.ansible_inventory
  filename       = format("%s/%s", var.admin_config_dir, local.ansible_inventory_filename)
}

resource "aws_s3_object" "ansible_inventory" {
  content_base64 = module.cluster.ansible_inventory
  bucket         = var.bucket
  key            = format("config/admin/%s", local.ansible_inventory_filename)
}

resource "null_resource" "admin_ssh_known_hosts" {
  provisioner "local-exec" {
    command     = format("ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook %s/%s", var.admin_config_dir, local.ansible_playbook_update_known_hosts)
    working_dir = var.admin_config_dir
  }

  depends_on = [
    module.network,
    module.cluster,
    local_file.admin_private_key,
    local_file.admin_ssh_config,
  ]
}

resource "null_resource" "kubeadm_init" {
  provisioner "local-exec" {
    command     = format("ansible-playbook %s/%s", var.admin_config_dir, local.ansible_playbook_init_cluster)
    working_dir = var.admin_config_dir
  }

  depends_on = [
    null_resource.admin_ssh_known_hosts,
  ]
}

resource "aws_s3_object" "admin_ssh_known_hosts" {
  source = format("%s/%s", var.admin_config_dir, local.ssh_known_hosts_filename)
  bucket = var.bucket
  key    = format("config/admin/%s", local.ssh_known_hosts_filename)

  depends_on = [
    null_resource.admin_ssh_known_hosts,
  ]
}
