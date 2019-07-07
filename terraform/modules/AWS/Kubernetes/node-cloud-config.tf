data "template_file" "node_cloud_config" {
  template = "${ file( "${ path.module }/node-cloud-config.yaml" )}"
  vars = {
    REGION = "${var.region}"
    s3_bucket_name = "${var.s3_bucket_name}"
  }
}


