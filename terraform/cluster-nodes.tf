resource "aws_instance" "clusterNodes" {
  count                       = "${var.number_of_cluster_instances}"
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.bastion_key.key_name}"
  instance_type               = "${var.clusterNodes_instance_type}"
  security_groups             = ["${aws_security_group.controlplane-sg.id}"]
  associate_public_ip_address = false
  subnet_id                   = "${aws_subnet.private_subnet.id}"
}

