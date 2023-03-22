
set -x

USER="jcastillo"
PASS="9zvB2cQf2tdC"
ORG="rescue-drop"
DOMAIN="rescue.drop"
REPO="production-image-build"

SCP_FLAGS="-pr"
run_scp "./_git-repo-orders-processor/*" "$(get_master):/tmp/gitrepo"
run_scp "./_git-repo-orders-processor/.github" "$(get_master):/tmp/gitrepo"

run_ssh "$(get_master)" bash <<EOF
cd /tmp/gitrepo/
git config --global init.defaultBranch main
git init
git config --local user.name "$USER"
git config --local user.email "$USER@$DOMAIN"
git remote add localhost http://$USER:$PASS@localhost:30080/$ORG/$REPO
git add .
git commit -m "Initial Commit"
git push localhost
EOF
