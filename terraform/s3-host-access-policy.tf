resource "aws_iam_policy" "securus-s3-access" {
  name        = "securus_s3_host_access"
  path        = "/"
  description = "Host based access to S3 bucket"

  policy = <<EOF

{
 "Version": "2012-10-17",
 "Statement": [
    {
    "Effect": "Allow",
    "Action": [
      "s3:ListBucket"
       ],
    "Resource": [
       "arn:aws:s3:::securus-config"
       ]
    },
    {
    "Effect": "Allow",
    "Action": [
      "s3:PutObject",
      "s3:GetObject",
      "s3:DeleteObject",
      "s3:ListObject"
       ],
    "Resource": [
       "arn:aws:s3:::securus-config/*"
       ]
    }
  ]
}
EOF
}

