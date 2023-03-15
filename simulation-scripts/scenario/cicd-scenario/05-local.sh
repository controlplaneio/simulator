
set -x

USER="developer"
PASS="zah2sahvojae2aishei6DeivuzeeV7lu"
ORG="supersecureorg"
REPO="production-image-build"

SCP_FLAGS="-pr"
run_scp "./_gitrepo/*" "$(get_master):/tmp/gitrepo"

run_ssh "$(get_master)" bash <<EOF
cd /tmp/gitrepo/
git config --global init.defaultBranch main
git init
git config --local user.name "$USER"
git config --local user.email "$USER@localhost.ctf"
git remote add localhost http://$USER:$PASS@localhost:30080/$ORG/$REPO
git add .
git commit -m "Initial Commit"
git push localhost
EOF
