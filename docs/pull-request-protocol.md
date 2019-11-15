# How we do pull requests

The following is a guide on contributing to the code base in pull request terms.

## Naming conventions

The name of the commit should hold the following format; <Type>(<scope>): commit name, e.g. Fix(scen): container-ambush typo.

- _type_: What type of commit is this? An added feature? Documentation? A fix? See convention examples [here](https://www.conventionalcommits.org/en/v1.0.0/)
- _scope_: Where did your commit apply? Options are:
  - "sim"; for changes to the golang binary
  - "tf"; for changes to terraform
  - "scen"; for scenario changes
  - "peturb"; for changes to the perturb script
  - "attack"; for changes to the attack container
  - "tools"; for changes to the scenario helper and/or migrationtools
  but feel free to add your own or combine if you have worked on several areas in the commit.
- _text_: What was the commit? This does not need to be long.

## Labels

Labels are a PR feature of github and should at the very minimum be set to the type and scope to inform the pick of the reviewer.

## Reviewing

A reviewer should do more than sanity check a PR unless it is just metadata or documentation. If there is any functional changes or fixes, the reviewer should run the code as it sits in the PR and check it passes suitably. 
