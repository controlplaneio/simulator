# Backing out AWS environment when Terraform state file lost

Most of the changes that Terraform has created  can be backed out via the AWS console, however the following may need to be actioned manually:

# Instance profile

To delete the instance profile run:

aws iam delete-instance-profile --instance-profile-name=simulator-instance-profile  --profile [aws profile if default not used]

# Delete ssh key pair

aws ec2 delete-key-pair --key-name [key name] --profile [aws profile if default not used] --region eu-west-2

