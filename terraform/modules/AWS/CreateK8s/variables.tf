variable "region" {
  description = "aws region"
  default = "eu-west-1"
}
variable "number_of_master_instances" {
  description = "number of master instances to create"
  default = "1"
}
variable "ami_id" {
  description = "ami to use"
  default = "ami-09d38086eb2b23925"
}
variable "master_instance_type" {
  description = "instance type for master node(s) "
  default = "t2.medium"
}
variable "number_of_cluster_instances" {
  description = "number of nodes to create"
  default = "1"
}
variable "cluster_nodes_instance_type" {
  description = "instance type for k8s nodes"
  default = "t1.micro"
}
