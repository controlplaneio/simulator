data "template_file" "cloud_config" {
  template = "${ file( "${ path.module }/cloud-config.yaml" )}"
  vars = {
    master_ip_addresses = "${var.master_ip_addresses}"
    node_ip_addresses = "${var.node_ip_addresses}"
  }
}
