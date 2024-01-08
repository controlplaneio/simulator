resource "terraform_data" "known_hosts" {
  provisioner "local-exec" {
    command     = "rm ${var.ssh_known_hosts_file}"
    working_dir = var.ssh_config_dir
    on_failure  = continue
  }

  provisioner "local-exec" {
    command     = "ssh-keyscan -T 30 -t ed25519 -H ${var.bastion_ip} 22 > ${var.ssh_known_hosts_file}"
    working_dir = var.ssh_config_dir
    quiet       = true
  }

  triggers_replace = [
    var.bastion_ip,
    var.instances,
  ]
}

resource "terraform_data" "known_hosts_instances" {
  for_each = var.instances

  provisioner "local-exec" {
    command     = "ssh -F ${var.ssh_config_file} bastion -C \"ssh-keyscan -t ed25519 -H ${each.value} 22\" >> ${var.ssh_known_hosts_file}"
    working_dir = var.ssh_config_dir
    quiet       = true
  }

  triggers_replace = [
    var.bastion_ip,
    var.instances,
  ]

  depends_on = [
    terraform_data.known_hosts
  ]
}
