# ACME Lego

Go 版本 ACME 客户端，支持 ACMEv2 协议，支持 ECC 证书，支持泛域名证书

![GitHub release (latest by date)](https://img.shields.io/github/v/release/alphatr/acme-lego?label=acme-lego&style=flat-square)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/alphatr/acme-lego/Build-Release?style=flat-square)
![GitHub](https://img.shields.io/github/license/alphatr/acme-lego?style=flat-square)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/alphatr/acme-lego?style=flat-square)

验证方式支持 HTTP/HTTPS 文件路径，HTTP/HTTPS 端口转发，DNS 验证（泛域名证书）

### 开始使用

1、下载对应平台的二进制文件 [acme-lego/releases/latest](https://github.com/alphatr/acme-lego/releases/latest) 到可执行目录中(如 Linux 下的 `/usr/local/bin/` 目录)，命名 `lego`

2、建立 Lego 默认配置文件 `/etc/lego/config.toml`，或者建立在其他目录的话，执行 lego 调用需要传入 `config` 参数为配置文件路径（或者说 `config` 参数默认为 `/etc/lego/config.toml` ）

在配置目录下建立 `config.toml` 配置文件，如下

```toml
### 基础配置
email = "acme@example.com" # 用于账户注册的邮箱
expire-days = 30 # 在临过期多少天执行续签

key-type = ["rsa2048", "ec256"] # 全局支持的证书类型
challenge = "http-path" # 全局支持的验证方式
after-renew = "systemctl reload nginx" # 整体续签成功后执行的命令

# 域名配置
[domain-group."a.example.com"]
options.public = "/web-path/certificate/acme" # 针对 http-path 验证，验证文件的位置
```

3、建立配置文件后，依次执行

```bash
lego reg # 注册账户
lego run # 执行申请证书
```

就可以将配置的域名申请证书

4、也可以在配置文件不存在情况下手动执行

```bash
lego reg --email="acme@example.com" # 用传入的邮箱执行账户申请
```

5、或者忽略配置，为单个域名执行单个证书的申请

```bash
lego run --domain="c.example.com" # 执行 c.example.com 域名的证书申请
```

6、域名续签，由于执行文件没有做为服务在后台守护，所以需要手动添加 crontab 任务定时执行下面命令

```bash
lego renew
```

或者单个域名续签

```bash
lego renew --domain="c.example.com"
```

### 支持的验证方式

#### `http-path`: HTTP 请求的路径验证

lego 直接将 ACME 服务器返回的验证文件写入配置的路径下，用来支持 HTTP 访问 `http://a.example.com/.well-known/acme-challenge/xxxxxx` 来验证

例如有 Nginx 配置

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

则需要配置对应域名的 `options.public` 配置为 nginx 中对应的 root 目录

```toml
options.public = "/public/demo/challenge"
```

#### `https-path`: HTTPS 请求的路径验证

HTTPS 请求和上面类似，只是验证变为了访问 `https://a.example.com/.well-known/acme-challenge/xxxxxx` 来验证，配置不变

#### `http-port`: HTTP 请求验证的端口转发

lego 建立独立 Web 服务器用来支持验证，建议在 Nginx 上配置将 `/.well-known/acme-challenge/` 请求转发到 lego 服务器

nginx 配置如下，其中 `proxy_pass` 根据转发需求来修改

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

则需要配置对应域名的 `options.server` 配置为 nginx 中对应的转发服务器端口

```toml
options.server = ":8013"
```

#### `https-port`: HTTPS 请求验证的端口转发

HTTPS 请求和上面类似，只是验证变为了访问 `https://a.example.com/.well-known/acme-challenge/xxxxxx` 来验证，配置不变

#### `dns-cloudflare`: 通过 Cloudflare API 进行 DNS 修改的验证

如果域名在 Cloudflare 管理，则可以通过配置 `options.token` 为 Cloudflare Token 来进行验证，[Cloudflare Token 文档](https://blog.cloudflare.com/api-tokens-general-availability/)

### 高级配置

默认配置文件 `$PATH/config.toml` 所在的目录 `$PATH/` 即为所有证书及账户信息的配置目录，可以通过 `root-dir` 参数修改到其他目录

配置文件还支持的配置参数有

```toml
root-dir = "/etc/lego" # 配置目录，默认为配置文件所在目录
log-level = "info" # 日志级别，可取值依次从高到低有 panic, fatal, error, warn, info, debug
dev = true # 是否是开发模式，开发模式可以自定义 ACME 服务地址，请求 ACME 地址会忽略 HTTPS 的证书验证
acme-url = "https://127.0.0.1:14000/dir" # 开发模式下生效，请求 ACME 的服务地址
```

dev 下 log 默认等级为 debug，非 dev 下为 info

### 配置目录结构

```
lego/
    account/
        account.json # 账户信息
        account.key # 账户私钥
    certificates/
        a.example.com/ # 每个域名单独目录
            fullchain.ecdsa-256.crt # ecc 公钥
            fullchain.rsa-2048.crt # rsa 公钥
            meta.ecdsa-256.json # ecc 数据文件
            meta.rsa-2048.json # rsa 数据文件
            privkey.ecdsa-256.key # ecc 私钥
            privkey.rsa-2048.key # rsa 私钥
        b.example.com/

```

更多功能参考 [/config/config.template.toml](/config/config.template.toml) 文件

### 为什么是 Go 版本

Go 版本用来写类似的工具使用起来及其方便，只需将编译后的二进制文件下载放在服务器执行就可以，不额外依赖其他环境，不需要安装 Go 语言环境

### 感谢

基于 [github.com/go-acme/lego](https://github.com/go-acme/lego) 开发

### License

[MIT](https://opensource.org/licenses/MIT)

Copyright (c) 2020-present, AlphaTr
