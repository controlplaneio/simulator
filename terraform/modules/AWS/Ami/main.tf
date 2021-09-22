// Simply search for AMIs in region own by Canonical
// Restrict search to x86_64, ssd and hvm virt type

data "aws_ami" "find_ami" {

  owners = ["self"]
  most_recent = true
  name_regex = "^kubernetes-simulator-20.*"

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

data "aws_ami" "find_ami_ubuntu_upstream" {

  owners      = ["099720109477"] // this is Canonical's id
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}
