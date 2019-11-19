//
// Do not hard code credentails in this file
// Do not place aws credentails file into this repo
//
provider "aws" {}

terraform {
  backend "s3" {
    key = "simulator.tfstate"
    // 'bucket='' must have this exact number of spaces for simulator to replace it properly
    bucket = "###REPLACED-BY-SIMULATOR###"
    // Optional, S3 Bucket Server Side Encryption
    encrypt = false
  }
}

data "terraform_remote_state" "state" {
  backend = "s3"
  config  = {
    key     = "simulator.tfstate"
    // 'bucket='' must have this exact number of spaces for simulator to replace it properly
    bucket = "###REPLACED-BY-SIMULATOR###"
    // Optional, S3 Bucket Server Side Encryption
    encrypt = false
  }
}
