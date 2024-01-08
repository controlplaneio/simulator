locals {
  bastion_sg_name   = format("%s-bastion", var.name)
  instances_sg_name = format("%s-instances", var.name)
  admin_key_name    = format("%s-admin-key", var.name)
  player_key_name   = format("%s-player-key", var.name)

  instances = merge([for i in module.instances : i.instances]...)

  hosts_by_group = merge([
    for i, g in var.instance_groups :
    { format("%s", lower(var.instance_groups[i].name)) = keys(module.instances[i].instances) }
  ]...)
}



