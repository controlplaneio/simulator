// Create S3 bucket

resource "aws_s3_bucket" "k8sjoin" {
  bucket_prefix = "k8sjoin"
  force_destroy = true

  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator Kubernetes S3 Bucket"
    },
  )
}

resource "aws_s3_bucket_public_access_block" "k8sjoin" {
  bucket = aws_s3_bucket.k8sjoin.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
