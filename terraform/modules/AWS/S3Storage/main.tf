// Create S3 bucket

resource "aws_s3_bucket" "k8sjoin" {
  bucket_prefix = "k8sjoin"
  acl           = "private"
  force_destroy = true

  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Kubernetes S3 Bucket"
    },
  )
}


