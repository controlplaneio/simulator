locals {
  ssh_config = templatefile("${path.module}/templates/ssh_config", {
    bastion_ip = var.bastion_ip
    instances = var.instances
    ssh_user = var.ssh_user
    ssh_force_tty = var.ssh_force_tty
    ssh_identity_file = var.ssh_identity_file
    ssh_known_hosts = var.ssh_known_hosts_file
  })
}

resource "local_file" "ssh_config" {
  content = local.ssh_config
  filename = format("%s/%s", var.ssh_config_dir, var.ssh_config_file)
  file_permission = "0600"
}
