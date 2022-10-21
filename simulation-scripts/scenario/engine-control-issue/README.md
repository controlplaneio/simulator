# Ingress-NGINX Exploitation Scenario

This scenarios aims to gain familiarity with 3 ingress-nginx CVEs and how to exploit them:
* CVE-2021-25742 - add snippets via annotation in versions v1.0.0
* CVE-2021-25745 - path exploit <v1.2.0
    * fixed via "deep inspector" regex...
* CVE-2021-25748 - pathsanitation can be bypassed in versions <v1.2.1


## Notes
- use PS baseline profile to prevent privesc

## Refs
- K8s Blog - https://kubernetes.io/blog/2022/04/28/ingress-nginx-1-2-0/
- Hackerone
    - 742 - https://hackerone.com/reports/1249583
    - 742 - https://hackerone.com/reports/1378175
    - 745/748 - https://hackerone.com/reports/1382919
