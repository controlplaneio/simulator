variable "ami_id" {
  description = "ami to use with Bastion host"
  default = "ami-09d38086eb2b23925"
}

variable "instance_type" {
  description = "instance type for Baston host"
  default = "t1.micro"
}

