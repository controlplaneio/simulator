#!/bin/bash

set -e

terraform  destroy -var-file=settings/bastion.tfvars

