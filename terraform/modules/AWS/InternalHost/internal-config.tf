data "cloudinit_config" "cc" {
  part {
    content_type = "text/cloud-config"
    merge_type   = var.cloudinit_merge_strategy
    content      = var.cloudinit_common
  }
  part {
    content_type = "text/cloud-config"
    merge_type   = var.cloudinit_merge_strategy
    content = templatefile("${path.module}/internal-config.yaml", {
      s3_bucket_name = var.s3_bucket_name
      version        = var.kubernetes_version
      internal_motd  = filebase64("${path.module}/motd.sh")
    })
  }
}

