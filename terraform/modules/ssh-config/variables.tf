variable "bastion_ip" {
  description = "The IP address of the bastion host."
}

variable "instances" {
  description = "Map of instance name to ip address."
  type        = map(string)
  default     = {}
}

variable "ssh_config_dir" {
  description = "The directory to write the ssh config files."
}

variable "ssh_config_file" {
  description = "The name of the ssh config file."
}

variable "ssh_user" {
  description = "The name of the ssh user."
}

variable "ssh_force_tty" {
  description = "Set RequestTTY Force in the ssh config file."
  type        = bool
  default     = false
}

variable "ssh_identity_file" {
  description = "The name of the private key."
}

variable "ssh_known_hosts_file" {
  description = "The of the known hosts file."
}
