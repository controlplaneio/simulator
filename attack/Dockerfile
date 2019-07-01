FROM kalilinux/kali-linux-docker:latest AS attack

RUN apt -y update        \
    && apt -y install    \
    metasploit-framework \
    exploitdb            \
    amap                 \
    dirb                 \
    dnschef              \
    dnsenum              \
    dnsmap               \
    dnsrecon             \
    dnswalk              \
    hping3               \ 
    ndiff                \ 
    netdiscover          \ 
    nikto                \
    nmap                 \
    padbuster            \
    reaver               \
    responder            \
    sqlmap               \
    socat                \
    tshark               \
    volafox              \
    volatility           \
    wireshark            \
    wpscan               \
    yersinia

CMD ["bash"]

