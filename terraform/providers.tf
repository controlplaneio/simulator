#
# Do not hard code credentails in this file
# Do not place aws credentails file into this repo
#
provider "aws" {
  region                  = "${var.region}"
  shared_credentials_file = "${var.shared_credentials_file}"
  profile                 = "${var.aws_profile}" 
}

