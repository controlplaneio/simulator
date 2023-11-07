resource "local_file" "player_private_key" {
  content_base64  = module.cluster.player_private_key
  filename        = format("%s/%s", var.player_ssh_bundle_dir, local.ssh_identity_filename)
  file_permission = "0600"
}

resource "local_file" "player_ssh_config" {
  content_base64  = module.cluster.player_ssh_config
  filename        = format("%s/%s", var.player_ssh_bundle_dir, local.ssh_config_filename)
  file_permission = "0600"
}

resource "null_resource" "player_ssh_known_hosts" {
  provisioner "local-exec" {
    command     = "ssh-keyscan -t rsa -H ${module.cluster.bastion_ip} 22 > ${local.ssh_known_hosts_filename}"
    working_dir = var.player_ssh_bundle_dir
  }

  triggers = {
    always_run = timestamp()
  }

  depends_on = [
    null_resource.admin_ssh_known_hosts,
  ]
}
