resource "local_file" "player_private_key" {
  content_base64  = module.cluster.player_private_key
  filename        = format("%s/%s", var.player_config_dir, local.ssh_identity_filename)
  file_permission = "0600"
}

// TODO: add vars for config path in S3

resource "aws_s3_object" "player_private_key" {
  content_base64 = module.cluster.player_private_key
  bucket         = var.bucket
  key            = format("config/player/%s", local.ssh_identity_filename)
}

resource "local_file" "player_ssh_config" {
  content_base64  = module.cluster.player_ssh_config
  filename        = format("%s/%s", var.player_config_dir, local.ssh_config_filename)
  file_permission = "0600"
}

resource "aws_s3_object" "player_ssh_config" {
  content_base64 = module.cluster.player_ssh_config
  bucket         = var.bucket
  key            = format("config/player/%s", local.ssh_config_filename)
}

resource "null_resource" "player_ssh_known_hosts" {
  provisioner "local-exec" {
    command     = "ssh-keyscan -t rsa -H ${module.cluster.bastion_ip} 22 > ${local.ssh_known_hosts_filename}"
    working_dir = var.player_config_dir
  }

  depends_on = [
    null_resource.admin_ssh_known_hosts,
  ]
}

resource "aws_s3_object" "player_ssh_known_hosts" {
  source = format("%s/%s", var.player_config_dir, local.ssh_known_hosts_filename)
  bucket = var.bucket
  key    = format("config/player/%s", local.ssh_known_hosts_filename)

  depends_on = [
    null_resource.player_ssh_known_hosts,
  ]
}
