package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/alphatr/acme-lego/common/bootstrap"
	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
	"github.com/alphatr/acme-lego/controller/account"
	"github.com/alphatr/acme-lego/controller/certificate"
)

const (
	defaultConfigPath string = "/etc/lego/config.toml"
)

func main() {
	app := cli.NewApp()
	cli.VersionPrinter = versionPrinter
	app.Version = GitVersionTag
	app.Name = "lego"
	app.Usage = "A Let's Encrypt Client"
	app.Authors = []*cli.Author{{Name: "AlphaTr"}}

	app.Commands = []*cli.Command{
		{
			Name:   "reg",
			Usage:  "create account",
			Action: account.Register,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "mail",
					Usage: "account email",
				},
			},
			Before: beforeCommand,
		},

		{
			Name:   "run",
			Usage:  "run get certificate",
			Action: certificate.Obtain,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "domain",
					Aliases: []string{"d"},
					Usage:   "certificate domain",
				},
				&cli.StringFlag{
					Name:  "http-path",
					Usage: ".well-known/acme-challenge path",
				},
			},
			Before: beforeCommand,
		},

		{
			Name:   "renew",
			Usage:  "renew certificate",
			Action: certificate.Renew,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "domain",
					Aliases: []string{"d"},
					Usage:   "certificate domain",
				},
			},
			Before: beforeCommand,
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   defaultConfigPath,
			Usage:   "config `FILE` to use for config",
			EnvVars: []string{"LEGO_CONFIG"},
		},
	}

	app.Run(os.Args)
}

func beforeCommand(ctx *cli.Context) error {
	configFile := ctx.String("config")

	if err := config.InitConfig(configFile); err != nil {
		err := errors.NewError(errors.ConfigInitErrno, err)
		return cli.NewExitError(err.Error(), 101)
	}

	if err := bootstrap.InitBootstrap(); err != nil {
		err := errors.NewError(errors.BootstrapInitErrno, err)
		return cli.NewExitError(err.Error(), 102)
	}

	config.Config.UserAgent = fmt.Sprintf("alphatr-lego-cli/%s", ctx.App.Version)
	return nil
}

func versionPrinter(ctx *cli.Context) {
	buildString := fmt.Sprintf("%s/%s; %s; git:%s", runtime.GOOS, runtime.GOARCH, runtime.Version(), GitHash)
	buildTime := fmt.Sprintf("build-time: %s", CurrentTime)
	fmt.Fprintf(ctx.App.Writer, "%s/v%s (%s) %s\n", ctx.App.Name, ctx.App.Version, buildString, buildTime)
}
