resource "random_uuid" "s3_iam_role_uuid" {}

// Create S3 bucket

resource "aws_s3_bucket" "k8sjoin" {
  bucket        = "k8sjoin-${random_uuid.s3_iam_role_uuid.result}"
  acl           = "private"
  force_destroy = true

  tags = {
    Name = "K8S Config"
  }
}

// Create IAM role, policy and instance profile
// used to assign to instances to access S3 bucket

resource "aws_iam_role" "simulator_s3_access_role" {
  name = "simulator-s3-host-role-${random_uuid.s3_iam_role_uuid.result}"

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


resource "aws_iam_role_policy" "simulator_s3_access_policy" {
  name = "simulator-s3-host-policy-${random_uuid.s3_iam_role_uuid.result}"
  role = "${aws_iam_role.simulator_s3_access_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": ["${aws_s3_bucket.k8sjoin.arn}"]
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject"
      ],
      "Resource": ["${aws_s3_bucket.k8sjoin.arn}/*"]
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "simulator_instance_profile" {
  name = "simulator-instance-profile-${random_uuid.s3_iam_role_uuid.result}"
  role = "${aws_iam_role.simulator_s3_access_role.name}"
}


