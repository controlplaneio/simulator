# IP Detection

The clusters created by simulator are insecure by design. As such, it would be irresponsible to leave them open to the entire internet as they could be a vector into a users' infrastructure. To network sandbox the clusters, simulator uses AWS security groups to whitelist the current users' IP for SSH to the bastion. All other external connections are blocked.

## Go External IP

The most challenging part of whitelisting the users' IP is accurately and repeatedly finding it. To solve this problem the [github.com/GlenDC/go-external-ip](https://github.com/GlenDC/go-external-ip) library is used. The `consensus` function from the library polls a selection of remote IP detection sources to build a consensus on what the calling IP is.

The default IP sources can be seen [in the source code](https://github.com/GlenDC/go-external-ip/blob/139229dcdddd5ad18f5c4912fcb905a4079e2a36/consensus.go#L23) where the HTTPS sources are:

* https://icanhazip.com/
* https://myexternalip.com/raw/

These HTTPS sources are given three votes to the following HTTP sources single vote:

* http://ifconfig.io/ip
* http://checkip.amazonaws.com/
* http://ident.me/
* http://whatismyip.akamai.com/
* http://tnx.nl/ip
* http://myip.dnsomatic.com/
* http://diagnostic.opendns.com/myip

Once all sources have returned an IP or timed out, the IP with most votes wins and is returned as the IP for the user.