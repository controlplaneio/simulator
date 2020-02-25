variable "access_cidr" {
  description = "cidr range of client connection"
  type        = list(string)
}

variable "vpc_id" {
  description = "VPC ID"
}

variable "public_subnet_cidr_block" {
  description = "Public subnet cidr block"
}

variable "private_subnet_cidr_block" {
  description = "Private subnet cidr block"
}

variable "default_tags" {
  description = "Default tags for all resources"
  type        = map(string)
}
