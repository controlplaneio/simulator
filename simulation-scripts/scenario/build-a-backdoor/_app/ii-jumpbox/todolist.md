# To Do List

1. SMTP forwarder for contact form
2. Create an backlog item for the development team to remove the operations mgmt port on the ii management service
3. Clean up any misconfigured and failed helm deployments, noticed ii-pord was hanging around
4. Research using a load balancer instead of using nodeport for ingress
5. Speak to SecOps about getting access to the Kyverno policies they've deployed
6. Review AI engine with development team and figure out cluster resource requirements
7. Use the bastion ssh access to nodeport 30080 to check website is working `ssh -F cp_simulator_config -L 8080:localhost:8080 -N bastion`