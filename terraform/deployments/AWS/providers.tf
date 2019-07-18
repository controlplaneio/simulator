#
# Do not hard code credentails in this file
# Do not place aws credentails file into this repo
#
provider "aws" {}

terraform {
  backend "s3" {
    key = "simulator.tfstate"
    //    region = "eu-west-2"
    bucket = "jk-tf-remote-state"
    //    profile = "controlplane"
    encrypt = false # Optional, S3 Bucket Server Side Encryption
  }
}

