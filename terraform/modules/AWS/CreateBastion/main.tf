resource "aws_instance" "bastion" {
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.bastion_key.key_name}"
  instance_type               = "${var.instance_type}"
  security_groups             = ["${var.security_group}"]
  associate_public_ip_address = true
  subnet_id                   = "${var.subnet_id}"
}
resource "aws_key_pair" "bastion_key" {
  key_name                    = "${var.access_key_name}"
  public_key                  = "${var.access_key}"
}
