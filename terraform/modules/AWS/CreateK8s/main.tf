
data "template_file" "init-script-master" {
  template = "${file("cloud-init/cloud-init-master.cfg")}"
  vars = {
    REGION = "${var.region}"
  }
}


data "template_cloudinit_config" "cloudinit-securus-master" {

  gzip = false
  base64_encode = false
  
  part {
    filename     = "cloud-init.cfg"
    content_type = "text/cloud-config"
    content      = "${data.template_file.init-script-master.rendered}"
  }

}

data "template_file" "init-script" {
  template = "${file("cloud-init/cloud-init.cfg")}"
  vars = {
    REGION = "${var.region}"
  }
}

data "template_cloudinit_config" "cloudinit-securus" {

  gzip = false
  base64_encode = false
  
  part {
    filename     = "cloud-init.cfg"
    content_type = "text/cloud-config"
    content      = "${data.template_file.init-script.rendered}"
  }

}

resource "aws_instance" "controlplane" {
  count                       = "${var.number_of_master_instances}"
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.bastion_key.key_name}"
  instance_type               = "${var.master_instance_type}"
  security_groups             = ["${aws_security_group.controlplane-sg.id}"]
  associate_public_ip_address = false
  subnet_id                   = "${aws_subnet.private_subnet.id}"
  user_data                   = "${data.template_cloudinit_config.cloudinit-securus-master.rendered}"
  iam_instance_profile        = "${aws_iam_instance_profile.instance_profile.id}"
}

resource "aws_instance" "cluster_nodes" {
  count                       = "${var.number_of_cluster_instances}"
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.bastion_key.key_name}"
  instance_type               = "${var.cluster_nodes_instance_type}"
  security_groups             = ["${aws_security_group.controlplane-sg.id}"]
  associate_public_ip_address = false
  subnet_id                   = "${aws_subnet.private_subnet.id}"
  user_data                   = "${data.template_cloudinit_config.cloudinit-securus.rendered}"
  depends_on                  = ["aws_instance.controlplane"]
  iam_instance_profile        = "${aws_iam_instance_profile.instance_profile.id}"
}

