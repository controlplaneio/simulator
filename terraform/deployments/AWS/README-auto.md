## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| access\_cidr | cidr range of client connection | string | n/a | yes |
| access\_key | ssh public key | string | n/a | yes |
| access\_key\_name | ssh key name | string | `"simulator_bastion_access_key"` | no |
| attack\_container\_tag | the docker tag of the attack container to use | string | `"latest"` | no |
| cluster\_nodes\_instance\_type | instance type for k8s nodes | string | `"t2.medium"` | no |
| default\_tags | Default tags for all resources | map | `<map>` | no |
| instance\_type | instance type for Baston host | string | `"t2.micro"` | no |
| master\_instance\_type | instance type for master node(s) | string | `"t2.medium"` | no |
| number\_of\_cluster\_instances | number of nodes to create | string | `"2"` | no |
| number\_of\_master\_instances | number of master instances to create | string | `"1"` | no |
| private\_subnet\_cidr | cidr range for private subnet | string | `"172.31.2.0/24"` | no |
| public\_subnet\_cidr | cidr range for public subnet | string | `"172.31.1.0/24"` | no |
| shared\_credentials\_file | location of aws credentials file | string | `"~/.aws/credentials"` | no |
| vpc\_cidr | cidr block for vpc | string | `"172.31.0.0/16"` | no |

## Outputs

| Name | Description |
|------|-------------|
| access\_cidr | Remote access IP |
| ami\_id | AMI used for all instances |
| bastion\_public\_ip | Bastion public IP |
| cluster\_nodes\_private\_ip | Cluster node private IPs |
| internal\_node\_private\_ip | Private Subnet node IP |
| master\_nodes\_private\_ip | Master node private IP |

