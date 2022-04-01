output "IamInstanceProfileId" {
  value       = aws_iam_instance_profile.simulator_instance_profile.id
  description = "IAM instance profile id"
}

output "IamAccessRoleName" {
  value       = aws_iam_role.simulator_s3_access_role.name
  description = "IAM simulator s3 access role"
}

