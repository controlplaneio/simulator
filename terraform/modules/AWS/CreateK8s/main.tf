data "template_file" "init_script_master" {
  template = "${file("cloud-init/cloud-init-master.cfg")}"
  vars = {
    REGION = "${var.region}"
    s3_bucket_name = "${var.s3_bucket_name}"
  }
}


data "template_cloudinit_config" "simulator_cloudinit_master" {

  gzip = false
  base64_encode = false

  part {
    filename     = "cloud-init.cfg"
    content_type = "text/cloud-config"
    content      = "${data.template_file.init_script_master.rendered}"
  }

}

data "template_file" "init_script_cluster" {
  template = "${file("cloud-init/cloud-init.cfg")}"
  vars = {
    REGION = "${var.region}"
    s3_bucket_name = "${var.s3_bucket_name}"
  }
}

data "template_cloudinit_config" "simulator_cloudinit_cluster" {

  gzip = false
  base64_encode = false

  part {
    filename     = "cloud-init.cfg"
    content_type = "text/cloud-config"
    content      = "${data.template_file.init_script_cluster.rendered}"
  }

}

resource "aws_instance" "simulator_controlplane_instances" {
  count                       = "${var.number_of_master_instances}"
  ami                         = "${var.ami_id}"
  key_name                    = "${var.key_pair_name}"
  instance_type               = "${var.master_instance_type}"
  security_groups             = ["${var.control_plane_sg_id}"]
  associate_public_ip_address = false
  subnet_id                   = "${var.private_subnet_id}"
  user_data                   = "${data.template_cloudinit_config.simulator_cloudinit_master.rendered}"
  iam_instance_profile        = "${var.iam_instance_profile_id}"
}

resource "aws_instance" "simulator_cluster_instances" {
  count                       = "${var.number_of_cluster_instances}"
  ami                         = "${var.ami_id}"
  key_name                    = "${var.key_pair_name}"
  instance_type               = "${var.cluster_nodes_instance_type}"
  security_groups             = ["${var.control_plane_sg_id}"]
  associate_public_ip_address = false
  subnet_id                   = "${var.private_subnet_id}"
  user_data                   = "${data.template_cloudinit_config.simulator_cloudinit_cluster.rendered}"
  depends_on                  = ["aws_instance.simulator_controlplane_instances"]
  iam_instance_profile        = "${var.iam_instance_profile_id}"
}

