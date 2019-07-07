data "template_file" "cloud-config" {
  template = "${ file( "${ path.module }/cloud-config.yaml" )}"
  vars = {
    REGION = "${var.region}"
  }
}
