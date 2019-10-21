variable "ami_id" {
  description = "ami to use with hosts"
}

variable "instance_type" {
  description = "instance type for Baston host"
  default     = "t2.micro"
}

variable "access_key_name" {
  description = "ssh key name"
  default     = "simulator_ssh_access_key"
}

variable "control_plane_sg_id" {
  description = "configure security group"
}

variable "private_subnet_id" {
  description = "private subnet id"
}

variable "default_tags" {
  description = "Default tags for all resources"
  type        = "map"
}

variable "bastion_public_ip" {
  description = "IP address of the bastion for connecting to run tests"
}

