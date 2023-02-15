data "cloudinit_config" "master" {
  count = var.number_of_master_instances
  part {
    content_type = "text/cloud-config"
    merge_type   = var.cloudinit_merge_strategy
    content      = var.cloudinit_common
  }
  part {
    content_type = "text/cloud-config"
    merge_type   = var.cloudinit_merge_strategy
    content = templatefile("${path.module}/master-cloud-config.yaml", {
      hostname            = "k8s-master-${count.index}"
      s3_bucket_name      = var.s3_bucket_name
      version             = var.kubernetes_version
      version_major_minor = local.version_major_minor
    })
  }
}

