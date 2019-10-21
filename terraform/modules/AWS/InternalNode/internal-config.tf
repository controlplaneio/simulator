data "template_file" "internal_config" {
  template = "${file("${path.module}/internal-config.yaml")}"
}


