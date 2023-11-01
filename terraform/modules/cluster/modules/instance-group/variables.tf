variable "name" {
  description = ""
}

variable "group" {
  description = ""
}

variable "instance_count" {
  description = ""
}

variable "ami_id" {
  description = ""
}

variable "instance_type" {
  description = ""
}

variable "key_name" {
  description = ""
}

variable "availability_zone" {
  description = ""
}

variable "subnet_id" {
  description = ""
}

variable "associate_public_ip_address" {
  description = ""
  type        = bool
}

variable "security_group_id" {
  description = ""
}

variable "iam_instance_profile" {
  description = ""
  default     = ""
}

variable "volume_type" {
  description = ""
}

variable "volume_size" {
  description = ""
}

variable "user_data" {
  description = ""
  default     = ""
}

variable "tags" {
  description = ""
  type        = map(string)
}
