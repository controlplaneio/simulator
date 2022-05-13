resource "aws_instance" "simulator_bastion" {
  ami                         = var.ami_id
  key_name                    = var.access_key_name
  instance_type               = var.instance_type
  vpc_security_group_ids      = [var.security_group]
  associate_public_ip_address = true
  subnet_id                   = var.subnet_id
  user_data = templatefile(
    "${path.module}/cloud-config.yaml",
    {
      master_ip_addresses      = var.master_ip_addresses
      node_ip_addresses        = var.node_ip_addresses
      internal_host_private_ip = var.internal_host_private_ip
      attack_container_tag     = var.attack_container_tag
      attack_container_repo    = var.attack_container_repo
      github_usernames         = local.access_github_usernames
      bastion_bashrc           = filebase64("${path.module}/bashrc")
      bastion_inputrc          = filebase64("${path.module}/inputrc")
      bastion_aliases          = filebase64("${path.module}/bash_aliases")
      bastion_auth_keys        = filebase64("${path.module}/authorized_keys.sh")
    }
  )

  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Bastion"
    },
  )

}
