Feature: perturb
  In order to run a simulation
  As an operator
  I need to perturb the kubernetes cluster in some way

  Scenario: perturbs an AWS cluster
    Given no clustername and a non digital ocean master IP and slave IP
    When I run perturb
    Then the output has "MASTER_IP"
    And the output has "SLAVE_IP"
    And the output has "BASTION_IP"
    Then I ssh into the master
    And The cluster is perturbed

  Scenario: perturbs a Digital Ocean cluster
    Given a clustername, no master IP and no slave IPs
    When I run perturb
    Then I ssh into the master
    And The cluster is perturbed

