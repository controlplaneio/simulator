
set -x

USER="jcastillo"
PASS="9zvB2cQf2tdC"
ORG="rescue-drop"
DOMAIN="rescue.drop"
REPO="production-image-build"
MASTER_IP="$(get_master)"
NODE_IP="$(get_node 1)"

run_scp "./_git-repo-orders-processor/*" "$MASTER_IP:/tmp/gitrepo"
run_scp "./_git-repo-orders-processor/.github" "$MASTER_IP:/tmp/gitrepo"

run_ssh "$MASTER_IP" bash <<EOF
cd /tmp/gitrepo/
git config --global init.defaultBranch main
git init
git config --local user.name "$USER"
git config --local user.email "$USER@$DOMAIN"
git remote add localhost http://$USER:$PASS@localhost:30080/$ORG/$REPO
sed -i -e "s/__REGISTRY_IP__/$MASTER_IP/g" .github/workflows/build.yaml
git add .
git commit -m "Initial Commit"
git push localhost
sed -i -e "s/__REGISTRY_IP__/$MASTER_IP/g" /etc/containerd/certs.d/reg.rescue.drop/hosts.toml
systemctl restart containerd
EOF

run_ssh "$NODE_IP" bash <<EOF
sed -i -e "s/__REGISTRY_IP__/$MASTER_IP/g" /etc/containerd/certs.d/reg.rescue.drop/hosts.toml
systemctl restart containerd
EOF

BADUSER="iramos"
BADPASS="B6ejWLgmcmE1"
BADREPO="test-ci"

run_scp "./_git-repo-test-ci/*" "$MASTER_IP:/tmp/cirepo"
run_scp "./_git-repo-test-ci/.github" "$MASTER_IP:/tmp/cirepo"

TOKEN=$(echo "cat /tmp/iramos-token" | run_ssh "$MASTER_IP")

run_ssh "$MASTER_IP" bash <<EOF
cd /tmp/cirepo/
sed -i -e "s/__TOKEN__/$TOKEN/g" .github/workflows/build.yaml
git config --global init.defaultBranch main
git init
git config --local user.name "$BADUSER"
git config --local user.email "$BADUSER@$DOMAIN"
git remote add localhost http://$BADUSER:$BADPASS@localhost:30080/$BADUSER/$BADREPO
git add .
git commit -m "Initial Commit"
git push localhost
sed -i -e "s/$TOKEN/\\\${{ secrets.TOKEN }}/g" .github/workflows/build.yaml
git add -u
git commit -m "Make token more secure"
sed -i -e "s/gcc/gcc make/g" Dockerfile
git add -u
git commit -m "Add make to Dockerfile"
git push localhost
#shred -u /tmp/iramos-token
EOF

BOTUSER="kgarner"
BOTPASS="lqeQkqiU3GnW"
BOTREPO="chatbot"

run_scp "./_git-repo-chatbot/" "$MASTER_IP:/tmp/chatbotrepo"
run_scp "./_git-repo-chatbot/.config" "$MASTER_IP:/tmp/chatbotrepo"
run_scp "./_git-repo-chatbot/.data" "$MASTER_IP:/tmp/chatbotrepo"
run_scp "./_git-repo-chatbot/.env" "$MASTER_IP:/tmp/chatbotrepo"

run_ssh "$MASTER_IP" bash <<EOF
cd /tmp/chatbotrepo/
git config --global init.defaultBranch main
git init
git config --local user.name "$BOTUSER"
git config --local user.email "$BOTUSER@$DOMAIN"
git remote add localhost http://$BOTUSER:$BOTPASS@localhost:30080/$BOTUSER/$BOTREPO
git add .
git commit -m "wip: template for chatbot"
git push localhost
EOF