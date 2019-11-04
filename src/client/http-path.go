package client

import (
	"fmt"

	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/providers/http/webroot"

	"alphatr.com/acme-lego/src/config"
)

type HTTPPathOptions struct {
	Public string
}

func init() {
	ProviderMap["http-path"] = ApplyHTTPPathProvider
}

func ApplyHTTPPathProvider(domain string, cli *lego.Client, conf *config.DomainConfig) error {
	options := parseOptions(conf.Options)
	provider, err := webroot.NewHTTPProvider(options.Public)
	if err != nil {
		return fmt.Errorf("init-http-provider[%s, %s]: %s", domain, conf.KeyType, err.Error())
	}

	cli.Challenge.SetHTTP01Provider(provider)
	return nil
}

func parseOptions(opt map[string]string) *HTTPPathOptions {
	options := &HTTPPathOptions{}
	if publicPath, ok := opt["http-path"]; ok {
		options.Public = publicPath
	}

	return options
}
