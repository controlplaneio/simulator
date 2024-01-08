locals {
  bastion_sg_name   = format("%s-bastion", var.name)
  instances_sg_name = format("%s-instances", var.name)
  admin_key_name    = format("%s-admin-key", var.name)
  player_key_name   = format("%s-player-key", var.name)

  ssh_config_instances = merge([for i in module.instances : i.instances]...)
  ssh_config_player = templatefile("${path.module}/templates/ssh_config", {
    bastion_ip        = aws_instance.bastion.public_ip,
    ssh_user          = "player"
    ssh_force_tty     = true
    ssh_identity_file = var.ssh_identity_filename
    ssh_known_hosts   = var.ssh_known_hosts_filename
    instances         = {}
  })
  ssh_config_admin = templatefile("${path.module}/templates/ssh_config", {
    bastion_ip        = aws_instance.bastion.public_ip,
    ssh_user          = "ubuntu"
    ssh_force_tty     = false
    ssh_identity_file = var.ssh_identity_filename
    ssh_known_hosts   = var.ssh_known_hosts_filename
    instances         = local.ssh_config_instances
  })
  ]...)
}
