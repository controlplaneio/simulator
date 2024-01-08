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
output "bastion_ip" {
  value = aws_instance.bastion.public_ip
}

output "hosts_by_group" {
  value = local.hosts_by_group
}
