# Pull official Kali Linux container image
FROM kalilinux/kali-linux-docker:latest

RUN apt -y update \
    && apt -y install nmap \
    && apt -y install metasploit-framework

CMD ["bash"]

