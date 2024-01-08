---
all:
  children:
    ungrouped:
%{ for group in keys(hosts_by_group) ~}
    ${group}s:
%{ endfor ~}
ungrouped:
  hosts:
    bastion:
%{ for group in keys(hosts_by_group) ~}
${group}s:
  hosts:
%{ for host in hosts_by_group[group] ~}
    ${host}:
%{ endfor ~}
%{ endfor ~}
