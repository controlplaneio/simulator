## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| access\_key\_name | ssh key name | string | `"simulator_ssh_access_key"` | no |
| ami\_id | ami to use with hosts | string | n/a | yes |
| bastion\_public\_ip | IP address of the bastion for connecting to run tests | string | n/a | yes |
| control\_plane\_sg\_id | configure security group | string | n/a | yes |
| default\_tags | Default tags for all resources | map | n/a | yes |
| instance\_type | instance type for Baston host | string | `"t2.micro"` | no |
| private\_subnet\_id | private subnet id | string | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| InternalNodePrivateIp | Internal node private ip address |

