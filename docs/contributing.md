# Contributing

> Work in progress instructions on how to contribute to various parts of simulator

<!-- toc -->

- [Bugs](#bugs)
- [Pull Requests](#pull-requests)
- [Components](#components)
  * [Attack Container](#attack-container)

<!-- tocstop -->

## Bugs

If you think you have found a bug please follow the instructions below.

- Please spend a small amount of time giving due diligence to the issue tracker. Your issue might be a duplicate.
- Open a [new issue](https://github.com/kubernetes-simulator/simulator/issues/new) if a duplicate doesn't already exist.
- Note the version of simulator you are running and the command line options you are using.
- Note the version of Kubernetes you are running.
- Remember users might be searching for your issue in the future, so please give it a meaningful title to help others.

## Pull Requests 

We welcome pull requests! 

- Your PR is more likely to be accepted if it focuses on just one change.
- Please include a comment with the results before and after your change. 
- Your PR is more likely to be accepted if it includes tests.
- You're welcome to submit a draft PR if you would like early feedback on an idea or an approach. 
- Happy coding!


## Components

### Attack Container

- edit code in the `./attack` directory
- commit changes on branch
- From the root of the repo, `cd attack && CONTAINER_TAG=super-cool-feature make docker-push` to push the tagged attack container
- `cd .. && make run` to run the launch container
- `simulator infra create --attack-container-tag=super-cool-feature`

The tag is defined by a Terraform variable called "attack_container_tag". The variable is threaded through from the deployment to the bastion and is then templated into the cloud-config to pull the appropriate tag and launch that tag when the
ubuntu user logs in (done by `simulator ssh attack`).

The golang binary has a corresponding `--attack-container-tag` flag and configuration variable to control what this is set to. This is written to `tfvars` during initialisation so that it propagates all the way through.
