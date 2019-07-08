## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| access\_key | ssh public key | string | n/a | yes |
| access\_key\_name | ssh key name | string | `"simulator_ssh_access_key"` | no |
| ami\_id | ami to use with hosts | string | `"ami-09d38086eb2b23925"` | no |
| instance\_type | instance type for Baston host | string | `"t1.micro"` | no |
| master\_ip\_addresses | Kubernetes master private IP addresses | string | n/a | yes |
| node\_ip\_addresses | Kubernetes node private IP addresses | string | n/a | yes |
| region | aws region | string | `"eu-west-1"` | no |
| security\_group | configure security group | string | n/a | yes |
| subnet\_id | configure subnet id | string | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| BastionPublicIp | Bastion server public ip address |
| KeyPairName | Name of Bastion SSH key in KMS |

