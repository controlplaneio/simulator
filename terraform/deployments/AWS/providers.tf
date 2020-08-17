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
      version = "~> 2.70"
    }
    null = {
      version = "~> 2.1"
    }
    random = {
      version = "~> 2.3"
    }
    template = {
      version = "~> 2.1"
    }
  }
}
