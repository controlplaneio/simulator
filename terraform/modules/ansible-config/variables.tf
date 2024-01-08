variable "ansible_config_dir" {
  description = "The directory to create the ansible config file."
}

variable "ansible_config_filename" {
  description = "The name of the ansible config file."
  default     = "ansible.cfg"
}

variable "ansible_roles_path" {
  description = "The ansible roles path."
}

variable "ssh_config_filename" {
  description = "The relative path to the ssh config file to use."
}

variable "ansible_inventory_filename" {
  description = "The name of the ansible inventory file."
  default     = "inventory.yaml"
}

variable "hosts_by_group" {
  description = "Map of group to hostnames."
  type        = map(list(string))
}
