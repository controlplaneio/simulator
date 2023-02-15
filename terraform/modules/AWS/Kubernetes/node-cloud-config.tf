data "cloudinit_config" "node" {
  count = var.number_of_cluster_instances
  part {
    content_type = "text/cloud-config"
    merge_type   = var.cloudinit_merge_strategy
    content      = var.cloudinit_common
  }
  part {
    content_type = "text/cloud-config"
    merge_type   = var.cloudinit_merge_strategy
    content = templatefile("${path.module}/node-cloud-config.yaml", {
      hostname       = "k8s-node-${count.index}"
      s3_bucket_name = var.s3_bucket_name
      version        = var.kubernetes_version
    })
  }
}

