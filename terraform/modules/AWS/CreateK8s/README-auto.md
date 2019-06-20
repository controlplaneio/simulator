## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| ami\_id | ami to use | string | `"ami-09d38086eb2b23925"` | no |
| cluster\_nodes\_instance\_type | instance type for k8s nodes | string | `"t1.micro"` | no |
| control\_plane\_sg\_id | Control plane (private) security group id | string | n/a | yes |
| iam\_instance\_profile\_id | IAM instance S3 access profile id | string | n/a | yes |
| key\_pair\_name | Name of ssh key held in KMS | string | n/a | yes |
| master\_instance\_type | instance type for master node(s) | string | `"t2.medium"` | no |
| number\_of\_cluster\_instances | number of nodes to create | string | `"1"` | no |
| number\_of\_master\_instances | number of master instances to create | string | `"1"` | no |
| private\_subnet\_id | Private subnet id | string | n/a | yes |
| region | aws region | string | `"eu-west-1"` | no |

