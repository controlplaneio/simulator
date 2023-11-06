
set -x

USER="jcastillo"
PASS="9zvB2cQf2tdC"
ORG="rescue-drop"
DOMAIN="rescue.drop"
REPO="production-image-build"
MASTER_IP="{{ master_ip }}"
NODE_IP="{{ node1_ip }}"

BADUSER="iramos"
BADPASS="B6ejWLgmcmE1"
BADREPO="test-ci"

TOKEN=$(cat /tmp/iramos-token)

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
sed -i -e "s/$TOKEN/\\\${{'{{'}} secrets.TOKEN {{'}}'}}/g" .github/workflows/build.yaml
git add -u
git commit -m "Make token more secure"
sed -i -e "s/gcc/gcc make/g" Dockerfile
git add -u
git commit -m "Add make to Dockerfile"
git push localhost
#shred -u /tmp/iramos-token
