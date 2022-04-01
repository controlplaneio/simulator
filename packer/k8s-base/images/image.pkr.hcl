variable "region" {
  type = string
  default = "eu-west-2"
}

locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

# source blocks are generated from your builders; a source can be referenced in
# build blocks. A build block runs provisioners and post-processors on a
# source.
source "amazon-ebs" "example" {
  ami_name = "kubernetes-simulator-${local.timestamp}"
  //  instance_type = "t2.micro"
  instance_type = "t2.medium"
  region = var.region
  source_ami_filter {
    filters = {
      name = "ubuntu/images/*ubuntu-focal-20.04-amd64-server-*"
      root-device-type = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners = [
      "099720109477"
    ]
  }
  ssh_username = "ubuntu"
}

# a build block invokes sources and runs provisioning steps on them.
build {
  sources = [
    "source.amazon-ebs.example"
  ]

  provisioner "file" {
    source = "../../common/authorized_keys.sh"
    destination = "/tmp/authorized_keys.sh"
  }

  provisioner "shell" {
    script = "../scripts/setup.sh"
  }
}
