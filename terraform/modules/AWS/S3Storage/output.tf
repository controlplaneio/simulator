output "IamInstanceProfileId" {
  value       = "${aws_iam_instance_profile.simulator_instance_profile.id}"
  description = "IAM instance profile id"
}
output "IamAccessRoleName" {
  value       = "${aws_iam_role.simulator_s3_access_role.name}"
  description = "IAM simulator s3 access role"
}

output "S3BucketName" {
  value       = "${aws_s3_bucket.k8sjoin.id}"
  description = "Name of S3 bucket created"
}
