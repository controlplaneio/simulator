data "aws_ami" "find_ami" {
  most_recent = true

  // need to check the below
  filter {
    name   = "name"
    values = ["18.04 LTS"]
  }
  // should this be amd64, needs clarification
  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs-ssd"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

