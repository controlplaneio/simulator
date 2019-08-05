## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| access\_key\_name | Name of ssh key held in KMS | string | n/a | yes |
| ami\_id | ami to use | string | `"ami-09d38086eb2b23925"` | no |
| bastion\_public\_ip | IP address of the bastion for connecting to run tests | string | n/a | yes |
| cluster\_nodes\_instance\_type | instance type for k8s nodes | string | `"t1.micro"` | no |
| control\_plane\_sg\_id | Control plane (private) security group id | string | n/a | yes |
| default\_tags | Default tags for all resources | map | n/a | yes |
| iam\_instance\_profile\_id | IAM instance S3 access profile id | string | n/a | yes |
| master\_instance\_type | instance type for master node(s) | string | `"t2.medium"` | no |
| number\_of\_cluster\_instances | number of nodes to create | string | `"1"` | no |
| number\_of\_master\_instances | number of master instances to create | string | `"1"` | no |
| private\_subnet\_id | Private subnet id | string | n/a | yes |
| s3\_bucket\_name | Name  of s3 bucket | string | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| K8sMasterPrivateIp | Kubernetes master private ip address |
| K8sNodesPrivateIp | Kubernetes node(s) private ip address |

