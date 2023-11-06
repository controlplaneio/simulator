
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


BOTUSER="kgarner"
BOTPASS="lqeQkqiU3GnW"
BOTREPO="chatbot"


cd /tmp/chatbotrepo/
git config --global init.defaultBranch main
git init
git config --local user.name "$BOTUSER"
git config --local user.email "$BOTUSER@$DOMAIN"
git remote add localhost http://$BOTUSER:$BOTPASS@localhost:30080/$BOTUSER/$BOTREPO
git add .
git commit -m "wip: template for chatbot"
git push localhost
