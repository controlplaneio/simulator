add-apt-repository ppa:rmescandon/yq &> /dev/null
apt update &> /dev/null
apt install yq -y &> /dev/null

yq w -i /var/lib/kubelet/config.yaml authentication.anonymous.enabled true
yq w -i /var/lib/kubelet/config.yaml authorization.mode AlwaysAllow

systemctl daemon-reload
systemctl restart kubelet.service
