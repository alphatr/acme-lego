package client

import (
	"fmt"

	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/providers/http/webroot"

	"alphatr.com/acme-lego/src/config"
)

func init() {
	ProviderMap["http-path"] = ApplyHTTPPathProvider
}

func ApplyHTTPPathProvider(domain string, cli *lego.Client, conf *config.DomainConfig) error {
	provider, err := webroot.NewHTTPProvider(conf.HTTPPath)
	if err != nil {
		return fmt.Errorf("init-http-provider[%s, %s]: %s", domain, conf.KeyType, err.Error())
	}

	cli.Challenge.SetHTTP01Provider(provider)
	return nil
}
