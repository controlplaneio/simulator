resource "aws_instance" "controlplane" {
  count                       = "${var.number_of_master_instances}"
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.bastion_key.key_name}"
  instance_type               = "${var.master_instance_type}"
  security_groups             = ["${aws_security_group.controlplane-sg.id}"]
  associate_public_ip_address = false
  subnet_id                   = "${aws_subnet.private_subnet.id}"
}

