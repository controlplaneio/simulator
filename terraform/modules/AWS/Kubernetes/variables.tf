variable "number_of_master_instances" {
  description = "number of master instances to create"
  default     = "1"
}
variable "ami_id" {
  description = "ami to use"
  // Ensure we can SSH as root for the goss tests and also for preturb.sh
  default     = "ami-09d38086eb2b23925"
}
variable "bastion_public_ip" {
  description = "IP address of the bastion for connecting to run tests"
}
variable "master_instance_type" {
  description = "instance type for master node(s) "
  default     = "t2.medium"
}
variable "number_of_cluster_instances" {
  description = "number of nodes to create"
  default     = "1"
}
variable "cluster_nodes_instance_type" {
  description = "instance type for k8s nodes"
  default     = "t1.micro"
}
variable "access_key_name" {
  description = "Name of ssh key held in KMS"
}
variable "control_plane_sg_id" {
  description = "Control plane (private) security group id"
}
variable "private_subnet_id" {
  description = "Private subnet id"
}
variable "iam_instance_profile_id" {
  description = "IAM instance S3 access profile id"
}
variable "s3_bucket_name" {
  description = "Name  of s3 bucket"
}

