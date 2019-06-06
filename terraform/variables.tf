
variable "region" {
  description = "aws region"
  default = "eu-west-1"
}

variable "shared_credentials_file" {
  description = "location of aws credentiala file"
}

variable "aws_profile" {
  description = "aws profile"
  default = "default"
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
  description = "cidr range of client connection"
}

variable "public_subnet_cidr" {
  description = "cidr range for public subnet"
}

variable "private_subnet_cidr" {
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

