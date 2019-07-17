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

variable "security_group" {
  description = "configure security group"
}

variable "subnet_id" {
  description = "configure subnet id"
}

variable "master_ip_addresses" {
  description = "Kubernetes master private IP addresses"
}

variable "node_ip_addresses" {
  description = "Kubernetes node private IP addresses"
}
