data "template_file" "master_cloud_config" {
  count = "${var.number_of_master_instances}"
  template = "${file("${path.module}/master-cloud-config.yaml")}"
  vars = {
    hostname = "kubernetes-master-${count.index}"
    s3_bucket_name = "${var.s3_bucket_name}"
  }
}
