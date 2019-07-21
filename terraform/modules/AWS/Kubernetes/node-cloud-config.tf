data "template_file" "node_cloud_config" {
  count = "${var.number_of_cluster_instances}"
  template = "${file("${path.module}/node-cloud-config.yaml")}"
  vars = {
    hostname = "kubernetes-node-${count.index}"
    s3_bucket_name = "${var.s3_bucket_name}"
  }
}


