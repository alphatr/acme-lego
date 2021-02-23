# ACME Lego

English | [简体中文](./readme.zh-cn.md)

ACME client for Go version, supports ACMEv2 protocol, supports ECC certificate, supports wildcard SSL certificate

![GitHub release (latest by date)](https://img.shields.io/github/v/release/alphatr/acme-lego?label=acme-lego&style=flat-square)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/alphatr/acme-lego/Build-Release?style=flat-square)
![GitHub](https://img.shields.io/github/license/alphatr/acme-lego?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/alphatr/acme-lego?style=flat-square)

Challenge method supports HTTP/HTTPS file path, HTTP/HTTPS port forwarding, DNS challenge (wildcard SSL certificate)

### Start Using

1. Download the binary file of the corresponding platform [acme-lego/releases/latest](https://github.com/alphatr/acme-lego/releases/latest) into the executable directory(for Linux is `/usr/local/bin/` directory), and rename `lego`

2. Create Lego default config file `/etc/lego/config.toml`, or if it is created in another directory, to execute `lego`, you need to pass in the `config` parameter as the configuration file path(the `config` parameter default value is `/etc/lego/config.toml`)

Create a `config.toml` configuration file in the configuration directory as follows

```toml
### Base Config
email = "acme@example.com" # Email used for account registration
expire-days = 30 # How many days will the renewal be executed before expiration

key-type = ["rsa2048", "ec256"] # Globally supported certificate types
challenge = "http-path" # Globally supported challenge methods
after-renew = "systemctl reload nginx" # The command executed after the overall renewal is successful

# Domain Config
[domain-group."a.example.com"]
options.public = "/web-path/certificate/acme" # For http-path challenge, verify the location of the file
```

3. After create the configuration file, then execute

```bash
lego reg # register account
lego run # Execution obtain certificate
```

You can obtain for a certificate with the configured domain

4. It can also be executed manually if the configuration file does not exist

```bash
lego reg --email="acme@example.com" # Execution register account using the params email
```

5. Ignore the configuration and execution a single certificate obtain for a single domain

```bash
lego run --domain="c.example.com" # execution c.example.com domain certificate obtain
```

6. Domain certificate renewal, because the execution file is not guarded in the background as a service, you need to manually add crontab tasks and execute the following commands regularly

```bash
lego renew
```

Or single domain certificate renewal

```bash
lego renew --domain="c.example.com"
```

### Supported challenge methods

#### `http-path`: Path challenge for HTTP requests

`lego` writes the challenge file returned by the ACME server into the configured path to support HTTP request to `http://a.example.com/.well-known/acme-challenge/xxxxxx` to verify

For example, there is `Nginx` configuration

```nginx
server {
    listen       80;
    server_name  a.example.com;
    root         /public/demo/www;

    location ^~ /.well-known/acme-challenge/ {
        default_type  "text/plain";
        root     /public/demo/challenge/;
    }
}
```

You need to configure the `options.public` of the corresponding domain as the corresponding root directory in nginx

```toml
options.public = "/public/demo/challenge"
```

#### `https-path`: Path challenge for HTTPS requests

The HTTPS request is similar to the above, except that the challenge has changed to visit `https://a.example.com/.well-known/acme-challenge/xxxxxx` to verify, the configuration remains unchanged

#### `http-port`: Port forwarding for HTTP request challenge

`lego` start a web server to support challenge. It is recommended to configure on Nginx to forward `/.well-known/acme-challenge/` requests to the lego server

The nginx configuration is as follows, where `proxy_pass` is modified according to forwarding requirements

```
server {
    listen       80;
    server_name  a.example.com;
    root         /public/demo/www;

    location ^~ /.well-known/acme-challenge/ {
        proxy_http_version  1.1;
        proxy_redirect      off;
        proxy_set_header    Host $http_host;
        proxy_pass          http://127.0.0.1:8013$request_uri;
    }
}
```

You need to configure the `options.server` of the corresponding domain as the corresponding forwarding server port in nginx

```toml
options.server = ":8013"
```

#### `https-port`: Port forwarding for HTTPS request challenge

The HTTPS request is similar to the above, except that the challenge has changed to visit `https://a.example.com/.well-known/acme-challenge/xxxxxx` to verify, the configuration remains unchanged

#### `dns-cloudflare`: Verify DNS challenge through Cloudflare API

If the domain is managed by Cloudflare, it can be verified by configuring `options.token` as Cloudflare Token, [Cloudflare Token Docs](https://blog.cloudflare.com/api-tokens-general-availability/)

### Advanced configuration

The directory where the default configuration file `$PATH/config.toml` is located. `$PATH/` is the configuration directory for all certificates and account information, which can be modified to other directories through the `root-dir` parameter

The configuration parameters that the configuration file also supports are

```toml
root-dir = "/etc/lego" # Configuration directory, the default is the directory where the configuration file is located
log-level = "info" # Log level, the possible values from high to low are panic, fatal, error, warn, info, debug
dev = true # development mode, the development mode can customize the ACME service address, requesting the ACME address will ignore the HTTPS certificate verification
acme-url = "https://127.0.0.1:14000/dir" # Effective in development mode, request the service address of ACME
```

The default level of log under dev is debug, and under non-dev is info

### Configuration directory structure

```
lego/
    account/
        account.json # Account information
        account.key # Account private key
    certificates/
        a.example.com/ # Separate directory for each domain
            fullchain.ecdsa-256.crt # ecc public key
            fullchain.rsa-2048.crt # rsa public key
            meta.ecdsa-256.json # ecc data file
            meta.rsa-2048.json # rsa data file
            privkey.ecdsa-256.key # ecc private key
            privkey.rsa-2048.key # rsa private key
        b.example.com/

```

More function reference [/config/config.template.toml](/config/config.template.toml) file

### Why use Go language

The Go language is used to write similar tools and it is very convenient to use. You only need to download the compiled binary file and place it on the server for execution. It does not rely on other environments and does not need to install the Go language environment.

### Thanks

[github.com/go-acme/lego](https://github.com/go-acme/lego) based development

### License

[MIT](https://opensource.org/licenses/MIT)

Copyright (c) 2020-present, AlphaTr
