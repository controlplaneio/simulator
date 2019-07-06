## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| access\_cidr | cidr range of client connection | string | n/a | yes |
| access\_key | ssh public key | string | n/a | yes |
| access\_key\_name | ssh key name | string | `"bastion_access_key"` | no |
| ami\_id | ami to use with Bastion host | string | `"ami-09d38086eb2b23925"` | no |
| aws\_profile | aws profile | string | `"default"` | no |
| cluster\_nodes\_instance\_type | instance type for k8s nodes | string | `"t1.micro"` | no |
| instance\_type | instance type for Baston host | string | `"t1.micro"` | no |
| master\_instance\_type | instance type for master node(s) | string | `"t2.medium"` | no |
| number\_of\_cluster\_instances | number of nodes to create | string | `"1"` | no |
| number\_of\_master\_instances | number of master instances to create | string | `"1"` | no |
| private\_avail\_zone | availability zone for private subnet | string | `"eu-west-1a"` | no |
| private\_subnet\_cidr | cidr range for private subnet | string | `"172.31.2.0/24"` | no |
| public\_avail\_zone | availability zone for public subnet | string | `"eu-west-1a"` | no |
| public\_subnet\_cidr | cidr range for public subnet | string | `"172.31.1.0/24"` | no |
| region | aws region | string | `"eu-west-1"` | no |
| shared\_credentials\_file | location of aws credentials file | string | `"~/.aws"` | no |
| vpc\_cidr | cidr block for vpc | string | `"172.31.0.0/16"` | no |

