package client

import (
	"fmt"

	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/providers/dns/cloudflare"

	"alphatr.com/acme-lego/src/config"
)

func init() {
	ProviderMap["dns-cloudflare"] = ApplyDNSCloudflareProvider
}

// ApplyDNSCloudflareProvider 应用 Cloudflare Provider
func ApplyDNSCloudflareProvider(domain string, cli *lego.Client, conf *config.DomainConfig) error {
	config := cloudflare.NewDefaultConfig()
	config.AuthToken = conf.Options["token"]

	provider, err := cloudflare.NewDNSProviderConfig(config)
	if err != nil {
		return fmt.Errorf("init-http-provider[%s, %s]: %s", domain, conf.KeyType, err.Error())
	}

	cli.Challenge.SetDNS01Provider(provider)
	return nil
}
