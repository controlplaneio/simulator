output "S3BucketName" {
  value       = aws_s3_bucket.k8sjoin.id
  description = "Name of S3 bucket created"
}

output "JoinBucketArn" {
  value       = aws_s3_bucket.k8sjoin.arn
  description = "ARN of S3 bucket created"
}

