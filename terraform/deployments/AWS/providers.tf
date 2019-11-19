//
// Do not hard code credentails in this file
// Do not place aws credentails file into this repo
//
provider "aws" {}

terraform {
  backend "s3" {
    key = "simulator.tfstate"
    bucket = "###REPLACED-BY-SIMULATOR###" # Must have exact number of spaces for simulator to replace properly
    encrypt = false # Optional, S3 Bucket Server Side Encryption
  }
}

data "terraform_remote_state" "state" {
  backend = "s3"
  config  = {
    key     = "simulator.tfstate"
    bucket = "###REPLACED-BY-SIMULATOR###"  # Must have exact number of spaces for simulator to replace properly
    encrypt = false # Optional, S3 Bucket Server Side Encryption
  }
}
