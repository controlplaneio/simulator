output "player_private_key" {
  value = base64encode(tls_private_key.player.private_key_pem)
}

output "player_ssh_config" {
  value = base64encode(local.ssh_config_player)
}

output "admin_private_key" {
  value = base64encode(tls_private_key.admin.private_key_pem)
}

output "admin_ssh_config" {
  value = base64encode(local.ssh_config_admin)
}

output "ansible_config" {
  value = base64encode(local.ansible_config)
}

output "ansible_inventory" {
  value = base64encode(local.ansible_inventory)
}

output "bastion_ip" {
  value = aws_instance.bastion.public_ip
}
