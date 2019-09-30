# Quick help

# Simulator init error

If you get the following error when running ```simulator init``` :

```
WARN	cmd/init.go:63	Simulator is already configured to use an S3 bucket named wildgeek
WARN	cmd/init.go:64	Please remove the state-bucket from simulator.yaml to create another
```

Run the following to resolve it :

```
cat /dev/null > ~/.kubesim/simulator.yaml`
```

Now re-run ```simulator init```

