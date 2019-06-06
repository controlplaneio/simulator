# securus-aws-k8s
Terraform templates to create bastion host and prublic private vpcs/subnets

# MacOS Notes

This requires that the following installed:

* terraform
* awscli

# Validate your aws credentials

Easiest way to validate what you creds are and that they are working is to run:
```
aws sts get-caller-identity
```
