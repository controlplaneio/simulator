data "template_file" "init-script" {
  template = "${file("cloud-init/bastion-cloud-init.cfg")}"
  vars = {
    REGION = "${var.region}"
  }
}

data "template_cloudinit_config" "simulator_cloudinit_bastion" {

  gzip = false
  base64_encode = false

  part {
    filename     = "cloud-init.cfg"
    content_type = "text/cloud-config"
    content      = "${data.template_file.init-script.rendered}"
  }

}

resource "aws_instance" "simulator_bastion" {
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.simulator_bastion_key.key_name}"
  instance_type               = "${var.instance_type}"
  security_groups             = ["${var.security_group}"]
  associate_public_ip_address = true
  subnet_id                   = "${var.subnet_id}"
  user_data                   = "${data.template_cloudinit_config.simulator_cloudinit_bastion.rendered}"
}
resource "aws_key_pair" "simulator_bastion_key" {
  key_name                    = "${var.access_key_name}"
  public_key                  = "${var.access_key}"
}
