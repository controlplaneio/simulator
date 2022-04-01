resource "aws_iam_role" "simulator_s3_access_role" {
  name_prefix = "simulator-instance-role"

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


  tags = merge(
    var.default_tags,
    {
      "Name" = "Simulator S3 Bucket Role"
    },
  )
}

resource "aws_iam_role_policy" "simulator_s3_access_policy" {
  name_prefix = "simulator-s3-host-policy"
  role        = aws_iam_role.simulator_s3_access_role.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": ["${var.join_bucket_arn}"]
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject"
      ],
      "Resource": ["${var.join_bucket_arn}/*"]
    }
  ]
}
EOF

}

resource "aws_iam_instance_profile" "simulator_instance_profile" {
  name_prefix = "simulator-instance-profile"
  role        = aws_iam_role.simulator_s3_access_role.name
}

# Add ECR pull rights to nodes
resource "aws_iam_role_policy_attachment" "simulator_ecr_access" {
  role       = aws_iam_role.simulator_s3_access_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy"
}

