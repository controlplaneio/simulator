variable "name" {
  description = ""
}

variable "network_id" {
  description = ""
}

variable "public_subnet_id" {
  description = ""
}

variable "private_subnet_id" {
  description = ""
}

variable "availability_zone" {
  description = ""
}

variable "bastion_ami_id" {
  description = ""
}

variable "bastion_instance_type" {
  description = "The instance type to use for the bastion."
  default     = "t3.micro"
}

variable "bastion_volume_type" {
  description = "The type of the root volume in for the bastion."
  default     = "gp2"
}

variable "bastion_volume_size" {
  description = "The size of the root volume in for the bastion."
  default     = "8"
}

variable "bastion_ssh_ingress" {
  description = "List of CIDR blocks to grant ssh access to bastion."
  type        = list(string)
}

variable "instance_groups" {
  description = ""
  type = list(object({
    name                 = string
    count                = number
    ami_id               = string
    public               = optional(bool, false)
    instance_type        = optional(string, "t3.micro")
    iam_instance_profile = optional(string, "")
    volume_type          = optional(string, "gp2")
    volume_size          = optional(string, "8")
  }))
  default = []
}

variable "tags" {
  description = "The common tags to apply to resources that support tagging."
  type        = map(string)
  default     = {}
}
