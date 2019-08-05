## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| access\_cidr | cidr range of client connection | string | n/a | yes |
| default\_tags | Default tags for all resources | map | n/a | yes |
| private\_subnet\_cidr\_block | Private subnet cidr block | string | n/a | yes |
| public\_subnet\_cidr\_block | Public subnet cidr block | string | n/a | yes |
| vpc\_id | VPC ID | string | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| BastionSecurityGroupID | Bastion security group id |
| ControlPlaneSecurityGroupID | Controlplane (Kubernetes) security group id |

