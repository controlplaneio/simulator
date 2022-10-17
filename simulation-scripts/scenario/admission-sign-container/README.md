# Capture the Flag - Locked Out (admission restriction with container signing)

## Overview

* _app (dir)
    * safe-webpage (dir)
    * main.go - webserver
        * keys (dir)
            * default.pub (easily searchable keys)
            * hashjack.pub
        * safe-signing-keys (dir)
            * default-signing.key (default signing key)
            * hashjack-super-secret.key
            * flag_ctf{LFI_4_KEY_NOW_SIGN} (first flag)
        * static (dir)
            * css (dir)
            * fonts (dir)
            * img (dir)
            * js (dir)
            * about.html (hint page)
            * development (vulnerable page)
            * index.html (landing page)
            * sign.html (container signing interface)
* apply (dir)
    * 01-kyverno-install.yaml (kyverno installation resources)
    * 02-signed-admission-policy.yaml (kyvenro clusterpolicy for challenge)
    * 03-scenario.yaml (The resources players will interact with for the challenge)
    * 04-flag-job.yaml (The job for providing the final flag to players)


## Notes
The main.go app requires the cosign binary.

> cosign generate-key-pair requires a password. This can be skipped with an env var of COSIGN_PASSWORD. This is essential for having a non-interactive signing and verification.

Originally the scenario was only going to allow use of ttl.sh for a registry but the way Kyverno policy is defined it is now defined for any registry.

To view the kyverno clusterpolicy use: ```kubectl describe clusterpolicies -A```

