data "template_file" "master_cloud_config" {
  template = "${ file( "${ path.module }/master-cloud-config.yaml" )}"
  vars = {
    REGION = "${var.region}"
    s3_bucket_name = "${var.s3_bucket_name}"
  }
}
