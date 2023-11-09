variable "name" {
  type    = string
  default = "simulator-bastion"
}

variable "region" {
  type    = string
  default = "eu-west-2"
}

variable "kube_version" {
  type    = string
  default = "1.28"
}

locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

build {
  name    = "simulator-bastion"
  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner "shell" {
    inline = [
      "sudo apt update",
      "sudo apt install -y apt-transport-https ca-certificates figlet curl jq python3-pip ansible socat",
      "ansible-galaxy collection install kubernetes.core",
      "pip install kubernetes",
    ]
  }

  provisioner "shell" {
    script = "scripts/common"
  }

  provisioner "shell" {
    environment_vars = [
      "KUBE_VERSION=${var.kube_version}",
      "PACKAGES=kubectl",
    ]
    script = "scripts/kubernetes"
  }


  provisioner "shell" {
    inline = [
      "rm .ssh/authorized_keys",
    ]
  }
}

source "amazon-ebs" "ubuntu" {
  ami_name      = "${var.name}-${var.kube_version}-${local.timestamp}"
  instance_type = "t2.micro"
  region        = "${var.region}"
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
}

packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}
