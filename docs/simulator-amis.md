# Simulator AMIs

Simulator uses two AMIs; one for the bastion, one for the Kubernetes instances.

These must be created in the target AWS account before launching the Simulator infrastructure by running the following
commands.

```shell
simulator image build bastion
simulator image build k8s
```
