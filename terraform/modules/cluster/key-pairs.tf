resource "aws_key_pair" "admin" {
  key_name   = local.admin_key_name
  public_key = tls_private_key.admin.public_key_openssh

  tags = merge(
    var.tags,
    {
      "Name" = format("%s Simulator Admin", var.name)
    }
  )
}

resource "tls_private_key" "admin" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "player" {
  key_name   = local.player_key_name
  public_key = tls_private_key.player.public_key_openssh

  tags = merge(
    var.tags,
    {
      "Name" = format("%s Simulator User", var.name)
    }
  )
}

resource "tls_private_key" "player" {
  algorithm = "RSA"
  rsa_bits  = 4096
}