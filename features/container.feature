Feature: Launch Container
  In order to launch simulations
  As an operator
  I want to isolate all scripts and environment from the host machine

  Scenario: configuration propagates to the container
    Given I set good env vars
    When I run the script
    Then the launch container has the env vars

  Scenario: outputs propagate through the scripts
    Given I set good env vars
    And I mount a fake standup script
    And I mount a fake perturb script
    And I mount a fake scenario-select menu
    When I run the launch container
    Then it threads outputs through from the standup script to the perturb script

  Scenario: container mints an SSH key
    Given I exec into the container
    When I cat ~/.ssh/id_rsa
    Then I have a fresh ssh key pair
