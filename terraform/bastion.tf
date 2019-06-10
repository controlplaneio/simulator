resource "aws_instance" "bastion" {
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.bastion_key.key_name}"
  instance_type               = "${var.instance_type}"
  security_groups             = ["${aws_security_group.bastion-sg.id}"]
  associate_public_ip_address = true
  subnet_id                   = "${aws_subnet.public_subnet.id}"
}

resource "aws_key_pair" "bastion_key" {
  key_name   = "access_key"
  public_key = "${var.access_key}"
}
