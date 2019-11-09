# Contributing

> Work in progress instructions on how to contribute to various parts of simulator 

## Attack Container

- edit code in the `./attack` directory
- commit changes on branch
- From the root of the repo, `cd attack && CONTAINER_TAG=super-cool-feature make docker-push` to push the tagged attack container
- `cd .. && make run` to run the launch container
- `simulator infra create --attack-container-tag=super-cool-feature`

The tag is defined by a Terraform variable called "attack_container_tag". The variable is threaded through from the deployment to the bastion and is then templated into the cloud-config to pull the appropriate tag and launch that tag when the
ubuntu user logs in (done by `simulator ssh attack`).

The golang binary has a corresponding `--attack-container-tag` flag and configuration variable to control what this is set to. This is written to `tfvars` during initialisation so that it propagates all the way through.
