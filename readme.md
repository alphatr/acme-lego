logo reg --email="baipan@baipan.me"
logo run --domain="baipan.me" --type ecc
logo renew --domain="baipan.me"


data/
    account/
    domains/
        baipan.me/ecc.key
        baipan.me/config.json
        baipan.me/

## build

```bash
go build -o bin/lego main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/lego main.go
```
