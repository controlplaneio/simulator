#FROM kalilinux/kali-linux-docker:latest AS attack
FROM debian:buster-slim

RUN apt -y update                       \
    && apt -y install --no-recommends   \
    openssh-server                      \
    hping3                              \
    nmap                                \
    curl

RUN mkdir -p /root/.ssh
RUN mkdir /var/run/sshd
COPY ./id_rsa.pub /root/.ssh/authorized_keys
RUN chown root:root /root/.ssh/authorized_keys
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]

