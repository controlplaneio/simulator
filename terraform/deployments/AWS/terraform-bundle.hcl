terraform {
  version = "0.13.7"
}

providers {
  aws = {
    versions = ["~> 4.51"]
  }
  null = {
    versions = ["~> 3.2"]
  }
  random = {
    versions = ["~> 3.4"]
  }
  cloudinit = {
    versions = ["~> 2.2"]
  }
}

# vim: ft=tf
