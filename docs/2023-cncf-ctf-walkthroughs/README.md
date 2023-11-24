# 2023 Cloud Native Computing Foundation - Capture the Flag (CTF) Scenarios

## 2023 Capture the Flag Scenarios

ControlPlane has developed 9 CTF Scenarios for the Cloud Native Computing Foundation (CNCF) providing hands-on experience of containers, Kubernetes and CI/CD infrastructure. The scenarios are designed to teach the participants about the how cloud native services work, the impact of vulnerable resources and how to remediate issues.

## Scenario Description

The table below outlines each scenario, learning objectives, technology used and the difficulty.

| Scenario | Scenario ID | Scenario Description | Learning Objective | Technology Used | Difficulty | No of Flags |
| --- | --- | --- | --- | --- | --- | --- |
| [Seven Seas](seven-seas/README.md) | seven-seas | Sail the Seven Seas, find all the missing map pieces and plunder the Royal Fortune | Kubernetes Fundamentals, Container Enumeration and Exploitation | Kubernetes Secrets, Container Images, Pod Security Standards, Network Policy, Pod Logs, Service Accounts and RBAC, Sidecar Containers | Easy | 2 |
| [Commandeer Container](commandeer-container/README.md) | commandeer-container | Use Kubernetes to Smuggle aboard and find the hidden treasure | Accessing Containers without `kubectl exec` | Kubernetes Secrets, Container Images, Service Accounts and RBAC | Easy | 1 |
| [CI Runner Next-Generation Breakout](ci-runner-ng-breakout/README.md) | ci-runner-ng-breakout |  | Container breakout via containerd | Docker, Containerd | Easy | 1 |
| [PSS Misconfiguration](pss-misconfiguration/README.md) | pss-misconfiguration | | Pod Security Standards, Pod Security Admission | Pod Security Standards, Pod Security Admission | Medium | 3 |
| [Build a Backdoor](build-a-backdoor/README.md) | build-a-backdoor | | Kubernetes Ingress, Services and Network Policies | Kubernetes Ingress, Services, Network Policies, Kyverno | Medium | 2 |
| [Cease and Desist](cease-and-desist/README.md) | cease-and-desist | | Cilium Network Policies | Kubernetes Secrets, Cilium Network Policies | Medium | 2 |
| [Devious Developer Data Dump](devious-developer-data-dump/README.md) | devious-developer-data-dump | | From secret discovery in a code repository to full cluster compromise | Gitea, GitHub Action Runners, Zot, SQL Database | Complex | 2 |
| [Identity Theft](identity-theft/README.md) | identity-theft |  | Realistic adversary initial access and OIDC token abuse | custom vulnerable application (pod schema validation), Dex, Kubernetes Services, Service Accounts and RBAC | Complex | 2 |
| [Coastline Cluster Attack](coastline-cluster-attack/README.md) | coastline-cluster-attack |  | Leveraging ephemeral containers for initial access, service account enumeration and privilege escalation, service account token abuse, vulnerable daemonsets | Ephemeral containers, Service Accounts and RBAC, Service Account Tokens, Custom "red herring" applications, Elasticsearch, Fluentbit Daemonsets | Complex | 3 |

### Difficulty Rating

**Easy** - A person with limited Kubernetes knowledge and is able to research how technology works, can complete these scenarios

**Medium** - More challenging scenarios which require a greater understanding or experience of technology and to problem solve

**Complex** - These are the most challenging scenarios, very little breadcrumbs are given and requires pre-existing knowledge or deep research to be conducted to complete
