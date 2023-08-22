FROM ubuntu:latest

RUN apt update && apt install -y openssh-server

COPY authorized_keys /root/.ssh/authorized_keys

# System Config
RUN chmod -x /etc/update-motd.d/*

COPY 01-royal-fortune /etc/update-motd.d/01-royal-fortune
RUN chmod +x /etc/update-motd.d/01-royal-fortune

RUN sed -i '/#Port 22/c\Port 2222' /etc/ssh/sshd_config
RUN sed -i '/#PermitRootLogin prohibit-password/c\PermitRootLogin prohibit-password' /etc/ssh/sshd_config
RUN sed -i '/#AuthorizedKeysFile	.ssh\/authorized_keys .ssh\/authorized_keys2/c\AuthorizedKeysFile    .ssh\/authorized_keys' /etc/ssh/sshd_config

RUN service ssh start

EXPOSE 2222

CMD ["/usr/sbin/sshd","-D"]