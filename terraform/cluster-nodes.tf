resource "aws_instance" "cluster_nodes" {
  count                       = "${var.number_of_cluster_instances}"
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.bastion_key.key_name}"
  instance_type               = "${var.cluster_nodes_instance_type}"
  security_groups             = ["${aws_security_group.controlplane-sg.id}"]
  associate_public_ip_address = false
  subnet_id                   = "${aws_subnet.private_subnet.id}"
  user_data                   = "${data.template_cloudinit_config.cloudinit-securus.rendered}"
}

