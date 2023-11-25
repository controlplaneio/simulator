variable "name" {
  description = "Name to identity the cluster."
  default     = "simulator"
}

variable "ansible_playbook_dir" {
  description = "The full path to the directory containing the Ansible Playbooks."
  default = "/simulator/ansible/playbooks"
}

variable "ansible_roles_dir" {
  description = "The full path to the directory containing the Ansible Roles."
  default = "/simulator/ansible/roles"
}

variable "admin_bundle_dir" {
  description = "The full path to the directory where the admin bundle files will be written."
  default = "/simulator/config/admin"
}

variable "player_bundle_dir" {
  description = "The full path to the directory where the player bundle files will be written."
  default = "/simulator/config/player"
}

# TODO: add switch to turn of ip lookup and ingress control

locals {
  ssh_identity_filename    = "simulator_rsa"
  ssh_config_filename      = "simulator_config"
  ssh_known_hosts_filename = "simulator_known_hosts"

  ansible_config_filename             = "ansible.cfg"
  ansible_inventory_filename          = "inventory.yaml"
  ansible_playbook_update_known_hosts = "${var.ansible_playbook_dir}/update-known-hosts.yaml"
  ansible_playbook_init_cluster       = "${var.ansible_playbook_dir}/init-cluster.yaml"

  bastion_ami_id        = data.aws_ami.bastion.id
  bastion_instance_type = "t2.small"

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

  name                     = var.name
  network_id               = module.network.network_id
  public_subnet_id         = module.network.public_subnet_id
  private_subnet_id        = module.network.private_subnet_id
  availability_zone        = random_shuffle.availability_zones.result[0]
  ssh_identity_filename    = local.ssh_identity_filename
  ssh_known_hosts_filename = local.ssh_known_hosts_filename
  bastion_ami_id           = local.bastion_ami_id
  bastion_instance_type    = local.bastion_instance_type
  instance_groups          = local.instance_groups
  ssh_config_filename      = local.ssh_config_filename
  tags                     = local.tags
  ansible_roles_dir        = var.ansible_roles_dir
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
