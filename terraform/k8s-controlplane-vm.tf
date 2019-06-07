resource "aws_instance" "controlplane" {
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.bastion_key.key_name}"
  instance_type               = "${var.instance_type}"
  security_groups             = ["${aws_security_group.controlplane-sg.id}","${aws_security_group.private-subnet-comms.id}"]
  associate_public_ip_address = false
  subnet_id                   = "${aws_subnet.private_subnet.id}"
}

resource "aws_key_pair" "controlplane_key" {
  key_name   = "controlplane_key"
  public_key = "${var.access_key}"
}
