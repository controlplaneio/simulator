output "IamInstanceProfileId" {
  value       = "${aws_iam_instance_profile.instance_profile.id}"
  description = "IAM instance profile id"
}
