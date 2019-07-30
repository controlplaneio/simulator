// Variables with no default to be place at the top of
// this file as per coding standards

variable "access_key" {
  description = "ssh public key"
}

variable "access_cidr" {
  description = "cidr range of client connection"
}

variable "s3_bucket_name" {
  description = "name of S3 bucket"
}

// Variables below are to have defined defaults

variable "shared_credentials_file" {
  description = "location of aws credentials file"
  default     = "~/.aws/credentials"
}

variable "access_key_name" {
  description = "ssh key name"
  default     = "simulator_bastion_access_key"
}

variable "vpc_cidr" {
  description = "cidr block for vpc"
  default     = "172.31.0.0/16"
}

variable "public_subnet_cidr" {
  description = "cidr range for public subnet"
  default     = "172.31.1.0/24"
}

variable "private_subnet_cidr" {
  description = "cidr range for private subnet"
  default     = "172.31.2.0/24"
}

variable "instance_type" {
  description = "instance type for Baston host"
  default     = "t2.micro"
}

variable "master_instance_type" {
  description = "instance type for master node(s) "
  default     = "t2.medium"
}

variable "number_of_master_instances" {
  description = "number of master instances to create"
  default     = "1"
}

variable "cluster_nodes_instance_type" {
  description = "instance type for k8s nodes"
  default     = "t2.medium"
}

variable "number_of_cluster_instances" {
  description = "number of nodes to create"
  default     = "2"
}

variable "default_tags" {
  description = "Default tags for all resources"
  type        = "map"
  default = {
    Product = "simulation-standalone"
  }
}

