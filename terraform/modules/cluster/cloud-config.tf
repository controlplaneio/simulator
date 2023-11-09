data "template_file" "cloud_config" {
  template = file("${path.module}/templates/cloud-config.yaml")
  vars     = {
    player_public_key = tls_private_key.player.public_key_openssh
  }
}
