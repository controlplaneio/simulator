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
      version = "~> 3.75.2"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 2.1"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 2.3"
    }
    template = {
      source  = "hashicorp/template"
      version = "~> 2.1"
    }
  }
}
