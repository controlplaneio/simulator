locals {
  ansible_inventory = templatefile("${path.module}/templates/inventory.yaml.tpl", {
    hosts_by_group = var.hosts_by_group
  })
}

resource "local_file" "ansible_inventory" {
  content = local.ansible_inventory
  filename = format("%s/%s", var.ansible_config_dir, var.ansible_inventory_filename)
}
