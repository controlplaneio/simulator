locals {
  ansible_config = templatefile("${path.module}/templates/ansible.cfg", {
    roles_path = var.ansible_roles_path
    ssh_config_filename = var.ssh_config_filename
  })
}

resource "local_file" "ansible_cfg" {
  content = local.ansible_config
  filename = format("%s/%s", var.ansible_config_dir, var.ansible_config_filename)
}
