add-apt-repository ppa:rmescandon/yq
apt update
apt install yq -y

yq w -i /var/lib/kubelet/config.yaml readOnlyPort 10255

systemctl daemon-reload
systemctl restart kubelet.service
