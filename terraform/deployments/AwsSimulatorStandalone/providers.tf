#
# Do not hard code credentails in this file
# Do not place aws credentails file into this repo
#
provider "aws" {
  region                  = "${var.region}"
  shared_credentials_file = "${var.shared_credentials_file}"
  profile                 = "${var.aws_profile}"
}

terraform {
  backend "s3" {
    key = "simulator.tfstate"
    region = "eu-west-1"
    bucket = "simulator-standalone-terraform-state"
    dynamodb_table = "simulator-standalone-terraform-state-locking"
    encrypt = true # Optional, S3 Bucket Server Side Encryption
  }
}
