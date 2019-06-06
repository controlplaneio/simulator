variable "security_group_name" {
  description = "Name for security group"
  default     = "bastion-security-group"
}

variable "access_key" {
  description = "ssh public key"
}

variable "private_vpc_cidr" {
  description = "private cidr block for vpc"
}

variable "public_vpc_cidr" {
  description = "public cidr block for vpc"
}

variable "access_cidr" {
  description = "cidr range of client comnection"
}

variable "public_cidr" {
  description = "cidr range for public subnet"
}

variable "private_cidr" {
  description = "cidr range for private subnet"
}

variable "ami_id" {
  description = "ami to use with Bastion host"
}

variable "instance_type" {
  description = "instance type for Baston host"
  default = "t2.micro"
}

variable "private_avail_zone" {
  description = "availability zone for private subnet"
}

variable "public_avail_zone" {
  description = "availability zone for public subnet"
}

