## Well-read

To test that the sample app returns no 502s the engineering team has run this ab
command with a post payload of both 10KB and 4KB.

### Generate a large json payload

```bash
dd if=/dev/urandom of=10k.json bs=1024 count=10
```

### Running the bench

```bash
ab -T "application/json" -n 1000 -c 100 -l -p 10k.json http://<your-app-name>.<your-domain>/api/boomerangnsq
```
