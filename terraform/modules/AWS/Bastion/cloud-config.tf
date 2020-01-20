data "template_file" "cloud_config" {
  template = "${file("${path.module}/cloud-config.yaml")}"
  vars = {
    master_ip_addresses      = "${var.master_ip_addresses}"
    node_ip_addresses        = "${var.node_ip_addresses}"
    internal_host_private_ip = "${var.internal_host_private_ip}"
    attack_container_tag     = "${var.attack_container_tag}"
    bastion_bashrc           = "${filebase64("${path.module}/bashrc")}"
    bastion_inputrc          = "${filebase64("${path.module}/inputrc")}"
    bastion_aliases          = "${filebase64("${path.module}/bash_aliases")}"
  }
}
