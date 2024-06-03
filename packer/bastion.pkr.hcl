packer {
  # pin to the last MPL2.0 releases (1.9.x)
  required_version = "~> 1.9.5"
}

variable "name" {
  type    = string
  default = "simulator-bastion"
}

variable "region" {
  type    = string
  default = "${env("AWS_REGION")}"
}

variable "kube_version" {
  type    = string
  default = "1.28"
}

locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
  name      = "${var.name}-${var.kube_version}-${local.timestamp}"
}

build {
  name = "simulator-bastion"
  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner "shell" {
    inline = [
      "sudo apt-get update",
      "sudo apt-get upgrade -y",
      "sudo apt-get install -y apt-transport-https ca-certificates figlet curl jq python3-pip ansible socat",
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
  ami_name      = "${local.name}"
  instance_type = "t2.micro"
  region        = "${var.region}"
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20231117"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
  tags = {
    K8s_Version   = "${var.kube_version}"
    Base_AMI_Name = "{{ .SourceAMIName }}"
    Type          = "Bastion"
  }
  snapshot_tags = {
    AMI_Name = "${local.name}"
  }
}

packer {
  required_plugins {
    amazon = {
      version = "~> 1"
      source  = "github.com/hashicorp/amazon"
    }
  }
}
