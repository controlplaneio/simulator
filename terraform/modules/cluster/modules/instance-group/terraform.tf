terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      version = "~> 5.20"
    }
    tls = {
      version = "~> 4.0"
    }
  }
}
