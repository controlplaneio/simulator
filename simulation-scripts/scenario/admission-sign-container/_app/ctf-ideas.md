# Capture the Flag - Locked Out (admission restriction with container signing)

## Overview

* safe-webpage (dir)
* main.go - webserver
    * keys (dir)
        * default.pub (easily searchable keys)
        * hashjack.pub
    * safe-signing-keys (dir)
        * default-signing.key (default signing key)
        * hashjack-super-secret.key
    * static (dir)
        * css (dir)
        * fonts (dir)
        * img (dir)
        * js (dir)
        * about.html (hint page)
        * development (vulnerable page)
        * index.html (landing page)
        * sign.html (container signing interface)

## Notes
The main.go app requires the cosign binary.


> cosign generate-key-pair requires a password. This can be skipped with an env var of COSIGN_PASSWORD. This is essential for having a non-interactive signing and verification.

tll.sh
```bash
$ IMAGE_NAME=$(uuidgen)
$ docker build -t ttl.sh/${IMAGE_NAME}:1h .
$ docker push ttl.sh/${IMAGE_NAME}:1h
```