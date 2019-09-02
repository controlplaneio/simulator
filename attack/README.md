# Attack Container

This container is run on the bastion host.  When you SSH onto the bastion using
`simulator ssh attack` you will land in here.  It comes with the set of tools
needed to complete a scenario.  

The container will also provide instructions for a scenario after the cluster
has been perturbed and it will provide a `start` script to get the user to the
point for the scenario as well as helper scripts to navigate around the cluster.

In the future all scenario related tasks should be here.  E.g. hints, finding
out how long you've been trying the scenario for etc etc.
