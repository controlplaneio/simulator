############################
# Create S3 bucket
#

resource "aws_s3_bucket" "k8sjoin" {
  bucket        = "${var.s3_bucket_name}"
  acl           = "private"
  force_destroy = true

  tags = {
    Name        = "K8S Config"
  }
}

############################
# Create IAM role, policy and instance profile
# used to assign to instances to access S3 bucket
#

resource "aws_iam_role" "securus_s3_access_role" {
  name = "securus_s3_host_access_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}


resource "aws_iam_role_policy" "securus_s3_access_policy" {
  name        = "securus_s3_host_access_policy"
  role        = "${aws_iam_role.securus_s3_access_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": ["arn:aws:s3:::${var.s3_bucket_name}"]
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject"
      ],
      "Resource": ["arn:aws:s3:::${var.s3_bucket_name}/*"]
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "instance_profile" {
  name = "instance_profile"
  role = "${aws_iam_role.securus_s3_access_role.name}"
}


