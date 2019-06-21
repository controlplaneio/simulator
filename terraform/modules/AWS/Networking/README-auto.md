## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| private\_avail\_zone | availability zone for private subnet | string | `"eu-west-1a"` | no |
| private\_subnet\_cidr | cidr range for private subnet | string | `"172.31.2.0/24"` | no |
| public\_avail\_zone | availability zone for public subnet | string | `"eu-west-1a"` | no |
| public\_subnet\_cidr | cidr range for public subnet | string | `"172.31.1.0/24"` | no |
| vpc\_cidr | cidr block for vpc | string | `"172.31.0.0/16"` | no |

## Outputs

| Name | Description |
|------|-------------|
| PrivateSubnetCidrBlock |  |
| PrivateSubnetId |  |
| PublicSubnetCidrBlock |  |
| PublicSubnetId |  |
| VpcId |  |

