variable "name" {
  type    = string
  default = "simulator-k8s"
}

variable "region" {
  type    = string
  default = "${env("AWS_REGION")}"
}

variable "containerd_version" {
  type    = string
  default = "1.7.7"
}

variable "runc_version" {
  type    = string
  default = "1.1.9"
}

variable "cni_version" {
  type    = string
  default = "1.3.0"
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
  name = "simulator-k8s"
  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner "shell" {
    inline = [
      "sudo apt update",
      "sudo apt install -y apt-transport-https ca-certificates figlet curl jq",
    ]
  }

  provisioner "shell" {
    script = "scripts/common"
  }

  provisioner "shell" {
    environment_vars = [
      "CONTAINERD_VERSION=${var.containerd_version}",
      "RUNC_VERSION=${var.runc_version}",
      "CNI_VERSION=${var.cni_version}",
    ]
    script = "scripts/containerd"
  }

  provisioner "shell" {
    environment_vars = [
      "KUBE_VERSION=${var.kube_version}",
      "PACKAGES=kubelet kubeadm kubectl",
      "PULL_IMAGES=true"
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
      name                = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-20231117"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
  tags = {
    Containerd_Version = "${var.containerd_version}"
    Runc_Version       = "${var.runc_version}"
    CNI_Version        = "${var.cni_version}"
    K8s_Version        = "${var.kube_version}"
    Base_AMI_Name      = "{{ .SourceAMIName }}"
    Kind               = "K8s"
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
