data "template_file" "cloud_config" {
  template = "${file("${path.module}/cloud-config.yaml")}"
  vars = {
    master_ip_addresses      = "${var.master_ip_addresses}"
    node_ip_addresses        = "${var.node_ip_addresses}"
    internal_node_private_ip = "${var.internal_node_private_ip}"
    attack_container_tag     = "${var.attack_container_tag}"
  }
}
