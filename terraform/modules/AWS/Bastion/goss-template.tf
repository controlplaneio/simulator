data "template_file" "goss_template" {
  template = "${file("${path.module}/goss.yaml")}"
  vars = {
    attack_container_tag = "${var.attack_container_tag}"
  }
}

