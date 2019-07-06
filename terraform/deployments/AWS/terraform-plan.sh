#!/bin/bash

set -e

terraform init
terraform plan -var-file=settings/bastion.tfvars

