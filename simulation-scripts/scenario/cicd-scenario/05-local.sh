
set -x

USER="jcastillo"
PASS="9zvB2cQf2tdC"
ORG="rescue-drop"
DOMAIN="rescue.drop"
REPO="production-image-build"
MASTER_IP="$(get_master)"

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
EOF
