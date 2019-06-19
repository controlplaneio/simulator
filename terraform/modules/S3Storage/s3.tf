resource "aws_s3_bucket" "k8sjoin" {
  bucket        = "securus-config"
  acl           = "private"
  force_destroy = true

  tags = {
    Name        = "K8S Config"
  }
}

