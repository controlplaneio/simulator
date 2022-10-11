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

resource "aws_s3_bucket_acl" "k8sjoin_acl" {
  bucket = aws_s3_bucket.k8sjoin.id
  acl    = "private"
}
