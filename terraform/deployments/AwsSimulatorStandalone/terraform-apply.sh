#!/bin/bash

set -e

terraform init
terraform apply -var-file=settings/bastion.tfvars -auto-approve

