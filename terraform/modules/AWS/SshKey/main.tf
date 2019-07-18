resource "random_uuid" "key_uuid" {}

resource "aws_key_pair" "simulator_bastion_key" {
  key_name   = "${var.access_key_name}-${random_uuid.key_uuid.result}"
  public_key = "${var.access_key}"
}

