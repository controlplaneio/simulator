resource "aws_instance" "simulator_bastion" {
  ami                         = var.ami_id
  key_name                    = var.access_key_name
  instance_type               = var.instance_type
  vpc_security_group_ids      = [var.security_group]
  associate_public_ip_address = true
  subnet_id                   = var.subnet_id
  user_data                   = data.template_file.cloud_config.rendered
  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Bastion"
    },
  )
  provisioner "file" {
    source      = "../../../../simulation-scripts/scenario/authorize-keys.sh"
    destination = "/tmp/authorize-keys.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "/tmp/authorize-keys.sh ${var.access_github_usernames}",
    ]
  }
}

