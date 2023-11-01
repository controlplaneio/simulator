data "template_file" "cloud_config" {
  template = file("${path.module}/templates/cloud-config.yaml")
  vars = {
    player_public_key = tls_private_key.player.public_key_openssh
  }
}


#admin_private_key = base64encode(tls_private_key.admin.private_key_openssh)
#ssh_config = base64encode(templatefile("${path.module}/templates/ssh_config_internal", {
#  user      = local.ssh_config_user
#  key_name  = "id_rsa"
#  instances = local.ssh_config_instances
#}))
#ansible_config = base64encode(templatefile("${path.module}/templates/ansible.cfg", {
#  ansible_inventory = local.ansible_inventory_file
#}))
#ansible_inventory = base64encode(templatefile("${path.module}/templates/inventory.yaml", {
#  ansible_inventory_instances = local.ansible_inventory_instances
#}))
#sudo = false