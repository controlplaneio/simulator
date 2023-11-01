variable "name" {
  description = ""
}

variable "vpc_cidr" {
  description = "The ipv4 cidr block for the vpc."
  default     = "10.0.0.0/16"
}

variable "availability_zone" {
  description = ""
}

variable "tags" {
  description = ""
  type        = map(string)
  default     = {}
}
