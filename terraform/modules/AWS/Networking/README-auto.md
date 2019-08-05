## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| default\_tags | Default tags for all resources | map | n/a | yes |
| private\_subnet\_cidr | cidr range for private subnet | string | `"172.31.2.0/24"` | no |
| public\_subnet\_cidr | cidr range for public subnet | string | `"172.31.1.0/24"` | no |
| vpc\_cidr | cidr block for vpc | string | `"172.31.0.0/16"` | no |

## Outputs

| Name | Description |
|------|-------------|
| PrivateSubnetCidrBlock | Private subnet cidr block |
| PrivateSubnetId | ID of private subnet |
| PublicSubnetCidrBlock | Public subnet cidr block |
| PublicSubnetId | ID of public subnet |
| VpcId | ID of VPC |

