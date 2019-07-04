output "IamInstanceProfileId" {
  value       = "${aws_iam_instance_profile.simulator_instance_profile.id}"
  description = "IAM instance profile id"
}
