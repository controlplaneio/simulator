terraform {
  # pin to the last MPL2.0 releases (1.5.x)
  required_version = "~> 1.5.7"
}

variable "name" {
  description = "Name to identity the cluster."
  default     = "simulator"
}

variable "bastion_ssh_ingress" {
  description = "List of CIDR blocks to grant ssh access to bastion."
  type        = list(string)
  default     = []
}

variable "ansible_playbook_dir" {
  description = "The full path to the directory containing the Ansible Playbooks."
  default     = "/simulator/ansible/playbooks"
}

variable "ansible_roles_dir" {
  description = "The full path to the directory containing the Ansible Roles."
  default     = "/simulator/ansible/roles"
}

variable "ansible_roles_path_extra" {
  description = "Extra directories to search for Ansible Roles."
  default     = ""
}

variable "admin_bundle_dir" {
  description = "The full path to the directory where the admin bundle files will be written."
  default     = "/simulator/config/admin"
}

variable "player_bundle_dir" {
  description = "The full path to the directory where the player bundle files will be written."
  default     = "/simulator/config/player"
}

locals {
  ssh_identity_filename    = "simulator_rsa"
  ssh_config_filename      = "simulator_config"
  ssh_known_hosts_filename = "simulator_known_hosts"

  ansible_config_filename       = "ansible.cfg"
  ansible_config_path           = format("%s/%s", var.admin_bundle_dir, local.ansible_config_filename)
  ansible_inventory_filename    = "inventory.yaml"
  ansible_roles_path            = length(var.ansible_roles_path_extra) > 0 ? format("%s:%s", var.ansible_roles_path_extra, var.ansible_roles_dir) : var.ansible_roles_dir
  ansible_playbook_init_cluster = "${var.ansible_playbook_dir}/init-cluster.yaml"

  bastion_ami_id        = data.aws_ami.bastion.id
  bastion_instance_type = "t2.small"
  bastion_ssh_ingress = length(var.bastion_ssh_ingress) > 0 ? var.bastion_ssh_ingress : [
    format("%s/32", trim(data.http.player_ip.response_body, "\n")),
  ]

  instance_groups = [
    {
      name                 = "master"
      count                = 1
      ami_id               = data.aws_ami.k8s.id
      public               = false
      instance_type        = "t2.medium"
      iam_instance_profile = ""
      volume_type          = "gp2"
      volume_size          = "20"
    },
    {
      name                 = "node"
      count                = 2
      ami_id               = data.aws_ami.k8s.id
      public               = false
      instance_type        = "t2.medium"
      iam_instance_profile = ""
      volume_type          = "gp2"
      volume_size          = "20"
    },
    {
      name                 = "internal"
      count                = 1
      ami_id               = data.aws_ami.k8s.id
      public               = false
      instance_type        = "t2.small"
      iam_instance_profile = ""
      volume_type          = "gp2"
      volume_size          = "20"
    }
  ]
  tags = {
    Name : title(var.name)
  }
}

module "network" {
  source = "../../modules/network"

  name              = var.name
  availability_zone = random_shuffle.availability_zones.result[0]
  tags              = local.tags
}

module "cluster" {
  source = "../../modules/cluster"

  name                  = var.name
  network_id            = module.network.network_id
  public_subnet_id      = module.network.public_subnet_id
  private_subnet_id     = module.network.private_subnet_id
  availability_zone     = random_shuffle.availability_zones.result[0]
  bastion_ami_id        = local.bastion_ami_id
  bastion_instance_type = local.bastion_instance_type
  bastion_ssh_ingress   = local.bastion_ssh_ingress
  instance_groups       = local.instance_groups
  tags                  = local.tags
}

resource "local_file" "admin_private_key" {
  content_base64  = module.cluster.admin_private_key
  filename        = format("%s/%s", var.admin_bundle_dir, local.ssh_identity_filename)
  file_permission = "0600"
}

module "admin_ssh_config" {
  source = "../../modules/ssh-config"

  bastion_ip           = module.cluster.bastion_ip
  instances            = module.cluster.instances
  ssh_config_dir       = var.admin_bundle_dir
  ssh_config_file      = local.ssh_config_filename
  ssh_user             = "ubuntu"
  ssh_identity_file    = local.ssh_identity_filename
  ssh_known_hosts_file = local.ssh_known_hosts_filename

  depends_on = [
    module.network,
    module.cluster,
    local_file.admin_private_key,
  ]
}

resource "local_file" "player_private_key" {
  content_base64  = module.cluster.player_private_key
  filename        = format("%s/%s", var.player_bundle_dir, local.ssh_identity_filename)
  file_permission = "0600"
}

module "player_ssh_config" {
  source = "../../modules/ssh-config"

  bastion_ip           = module.cluster.bastion_ip
  ssh_config_dir       = var.player_bundle_dir
  ssh_config_file      = local.ssh_config_filename
  ssh_user             = "player"
  ssh_force_tty        = true
  ssh_identity_file    = local.ssh_identity_filename
  ssh_known_hosts_file = local.ssh_known_hosts_filename

  depends_on = [
    module.network,
    module.cluster,
    local_file.player_private_key,
  ]
}

module "ansible_config" {
  source = "../../modules/ansible-config"

  ansible_config_dir      = var.admin_bundle_dir
  ansible_config_filename = local.ansible_config_filename
  ansible_roles_path      = length(var.ansible_roles_path_extra) > 0 ? var.ansible_roles_path_extra : var.ansible_roles_dir
  ssh_config_filename     = local.ssh_config_filename
  hosts_by_group          = module.cluster.hosts_by_group
}

resource "null_resource" "kubeadm_init" {
  provisioner "local-exec" {
    command = format("ansible-playbook %s", local.ansible_playbook_init_cluster)
    environment = {
      ANSIBLE_CONFIG = local.ansible_config_path
    }
    working_dir = var.admin_bundle_dir
  }

  depends_on = [
    module.admin_ssh_config,
    module.ansible_config,
  ]
}

resource "random_shuffle" "availability_zones" {
  input        = data.aws_availability_zones.available.names
  result_count = 1
}

data "aws_availability_zones" "available" {
  state = "available"
}

// TODO: filter these on K8s, containerd, runc versions?

data "aws_ami" "bastion" {
  most_recent = true
  owners = [
    "self",
  ]
  filter {
    name = "name"
    values = [
      "simulator-bastion-*"
    ]
  }
}

data "aws_ami" "k8s" {
  most_recent = true
  owners = [
    "self",
  ]
  filter {
    name = "name"
    values = [
      "simulator-k8s-*"
    ]
  }
}

terraform {
  backend "s3" {
  }
}

data "http" "player_ip" {
  url = "https://icanhazip.com/"
}
