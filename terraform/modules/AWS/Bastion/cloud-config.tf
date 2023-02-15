data "cloudinit_config" "cc" {
  part {
    content_type = "text/cloud-config"
    merge_type   = var.cloudinit_merge_strategy
    content      = var.cloudinit_common
  }
  part {
    content_type = "text/cloud-config"
    merge_type   = var.cloudinit_merge_strategy
    content = templatefile("${path.module}/cloud-config.yaml", {
      master_ip_addresses      = var.master_ip_addresses
      node_ip_addresses        = var.node_ip_addresses
      internal_host_private_ip = var.internal_host_private_ip
      attack_container_tag     = var.attack_container_tag
      attack_container_repo    = var.attack_container_repo
      bastion_motd             = filebase64("${path.module}/motd.sh")
    })
  }
}

