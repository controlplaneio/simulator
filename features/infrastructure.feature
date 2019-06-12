Feature: Infrastructure
  In order to launch simulations
  As an operator
  I need to stand up infrastructure in AWS

  Scenario: errors for invalid AWS credentials
    Given I set "SIMULATOR_AWS_CREDS_PATH" in the environment to garbage
    When I run the standup script
    Then I get error output (to avoid accidental leak of creds)

  Scenario: can connect to the bastion from the container
    Given I set "SIMULATOR_ACCESS_CIDR" to include the test machine
    When I run the standup script
    Then I can connect to the bastion from the container

  Scenario: cannot connect to the bastion from outside the container
    Given I set "SIMULATOR_ACCESS_CIDR" to not include the test machine
    When I run the standup script
    Then I cannot connect to the bastion from the host

  Scenario: outputs an rc script to be used by perturb.sh
    Given I set good env vars
    When I run the standup script
    Then the output has "MASTER_IP"
    And the output has "SLAVE_IP"
    And the output has "BASTION_IP"

