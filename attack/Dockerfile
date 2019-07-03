#FROM kalilinux/kali-linux-docker:latest AS attack
FROM ubuntu:latest

RUN apt -y update        \
    && apt -y install    \
    openssh-server       \
#    metasploit-framework \
#    exploitdb            \
#    amap                 \
#    dirb                 \
#    dnschef              \
#    dnsenum              \
#    dnsmap               \
#    dnsrecon             \
#    dnswalk              \
    hping3               \ 
#    ndiff                \ 
#    netdiscover          \ 
#    nikto                \
    nmap                 \
#    padbuster            \
#    reaver               \
#    responder            \
#    sqlmap               \
#    socat                \
#    tshark               \
#    volafox              \
#    volatility           \
#    wireshark            \
#    wpscan               \
#    yersinia             \
    curl

RUN mkdir -p /root/.ssh
RUN mkdir /var/run/sshd
ADD ./id_rsa.pub /root/.ssh/authorized_keys
RUN chown root:root /root/.ssh/authorized_keys
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]

