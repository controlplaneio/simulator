variable "ami_id" {
  description = "ami to use with hosts"
  type = string
  default = "ami-09d38086eb2b23925"
}

variable "instance_type" {
  description = "instance type for Baston host"
  type = string
  default = "t1.micro"
}
variable "access_key" {
  description = "ssh public key"
  type = string
}

variable "access_key_name" {
  description = "ssh key name"
  type = string
  default = "ssh_access_key"
}

variable "security_group" {
  description = "configure security group"
}

variable "subnet_id" {
  description = "configure subnet id"
}

