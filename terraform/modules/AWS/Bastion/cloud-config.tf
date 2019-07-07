data "template_file" "cloud_config" {
  template = "${ file( "${ path.module }/cloud-config.yaml" )}"
  vars = {
    REGION = "${var.region}"
  }
}
