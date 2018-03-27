package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"alphatr.com/acme-lego/src/account"
	"alphatr.com/acme-lego/src/client"
	"alphatr.com/acme-lego/src/config"
	"github.com/urfave/cli"
	acme "github.com/xenolf/lego/acmev2"
)

const (
	defaultPath string = "/etc/lego"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:   "reg",
			Usage:  "注册帐号",
			Action: registerAccount,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "mail",
					Usage: "账户邮件",
				},
			},
		},
		{
			Name:   "run",
			Usage:  "执行证书获取",
			Action: runClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "获取证书的域名",
				},
				cli.StringFlag{
					Name:  "http-path",
					Usage: "验证 .well-known/acme-challenge 文件目录",
				},
			},
		},
		{
			Name:   "renew",
			Usage:  "执行证书续期",
			Action: renewClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "domain, d",
					Usage: "获取证书的域名",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Value: defaultPath,
			Usage: "Directory `FILE` to use for storing the data",
		},
	}

	app.Run(os.Args)
}

func registerAccount(ctx *cli.Context) error {
	rootPath := ctx.Parent().String("path")
	mail := ctx.String("mail")

	err := config.InitConfig(rootPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error: init-config: %s", err.Error()), 101)
	}

	conf := config.Config
	if len(mail) > 0 {
		conf.Email = mail
	}

	acc, err := account.CreateAccount(conf.Email, rootPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error: create-account: %s", err.Error()), 102)
	}

	fmt.Printf("Success: create-account: %s\n", acc.Registration.URI)
	return nil
}

func runClient(ctx *cli.Context) error {
	rootPath := ctx.Parent().String("path")

	err := config.InitConfig(rootPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error: init-config: %s", err.Error()), 201)
	}

	conf := config.Config
	acme.HTTPClient = http.Client{Timeout: conf.HTTPTimeout * time.Second}
	acme.DNSTimeout = conf.DNSTimeout * time.Second

	acc, err := account.GetAccount(rootPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error: get-account: %s", err.Error()), 202)
	}

	keyDomain := ctx.String("domain")
	httpPath := ctx.String("http-path")

	if len(keyDomain) > 0 {
		domainConf, ok := conf.DomainGroup[keyDomain]
		if !ok {
			if len(httpPath) == 0 {
				return cli.NewExitError(fmt.Sprintf("Error: http-path-empty: %s", err.Error()), 203)
			}

			domainConf = &config.DomainConfig{
				Domains:   []string{keyDomain},
				KeyType:   []acme.KeyType{acme.RSA2048},
				Challenge: "http-file",
				HTTPPath:  httpPath,
			}
		}

		err := client.Run(keyDomain, acc, domainConf)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error: run-client: %s", err.Error()), 204)
		}

		return nil
	}

	for key, domainConf := range conf.DomainGroup {
		err := client.Run(key, acc, domainConf)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error: run-client: %s", err.Error()), 205)
		}
	}

	fmt.Printf("Success: run-client\n")
	return nil
}

func renewClient(ctx *cli.Context) error {
	rootPath := ctx.Parent().String("path")

	err := config.InitConfig(rootPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error: init-config: %s", err.Error()), 301)
	}

	conf := config.Config
	acme.HTTPClient = http.Client{Timeout: conf.HTTPTimeout * time.Second}
	acme.DNSTimeout = conf.DNSTimeout * time.Second

	acc, err := account.GetAccount(rootPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error: get-account: %s", err.Error()), 302)
	}

	keyDomain := ctx.String("domain")

	if len(keyDomain) > 0 {
		domainConf, ok := conf.DomainGroup[keyDomain]
		if !ok {
			return cli.NewExitError(fmt.Sprintf("Error: get-domain-info: %s", err.Error()), 303)
		}

		err := client.Renew(keyDomain, acc, domainConf)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error: run-client: %s", err.Error()), 304)
		}

		return nil
	}

	for key, domainConf := range conf.DomainGroup {
		err := client.Renew(key, acc, domainConf)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Error: run-client: %s", err.Error()), 305)
		}
	}

	fmt.Printf("Success: renew-client\n")
	return nil
}
