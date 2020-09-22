terraform {
  version = "0.13.3"
}

providers {
  aws = {
    versions = ["~> 3.6"]
  }
  null = {
    versions = ["~> 2.1"]
  }
  random = {
    versions = ["~> 2.1"]
  }
  template = {
    versions = ["~> 2.1"]
  }
}

# vim: ft=tf
