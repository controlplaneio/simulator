//
// Do not hard code credentails in this file
// Do not place aws credentails file into this repo
//
terraform {
  backend "s3" {
    key = "simulator.tfstate"
    // Optional, S3 Bucket Server Side Encryption
    encrypt = false
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.51"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.2"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.4"
    }
    cloudinit = {
      source  = "hashicorp/cloudinit"
      version = "~> 2.2"
    }
  }
}
