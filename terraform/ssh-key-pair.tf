resource "aws_key_pair" "bastion_key" {
  key_name   = "${var.access_key_name}"
  public_key = "${var.access_key}"
}
