resource "random_uuid" "key_uuid" {  }

resource "aws_instance" "simulator_bastion" {
  ami                         = "${var.ami_id}"
  key_name                    = "${aws_key_pair.simulator_bastion_key.key_name}"
  instance_type               = "${var.instance_type}"
  security_groups             = ["${var.security_group}"]
  associate_public_ip_address = true
  subnet_id                   = "${var.subnet_id}"
  user_data                   = "${data.template_file.cloud_config.rendered}"
}

resource "aws_key_pair" "simulator_bastion_key" {
  key_name                    = "${var.access_key_name}-${random_uuid.key_uuid.result}"
  public_key                  = "${var.access_key}"
}
