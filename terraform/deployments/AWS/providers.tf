//
// Do not hard code credentails in this file
// Do not place aws credentails file into this repo
//
provider "aws" {}

terraform {
  backend "s3" {
    key = "simulator.tfstate"
    // 'bucket='' must have this exact number of spaces for simulator to replace it properly
    bucket = "271119-pi1"
    // Optional, S3 Bucket Server Side Encryption
    encrypt = false
  }
}
