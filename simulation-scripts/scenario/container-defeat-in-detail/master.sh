#!/bin/bash

docker exec -e "ETCDCTL_API=3" $(docker ps -f "label=io.kubernetes.container.name=etcd" -q) etcdctl --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key put /super/secret/flag ohnomyetcd

kubectl taint nodes --all node-role.kubernetes.io/master- || true
