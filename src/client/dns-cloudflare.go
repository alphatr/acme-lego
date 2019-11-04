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

func ApplyDNSCloudflareProvider(domain string, cli *lego.Client, conf *config.DomainConfig) error {
	config := cloudflare.NewDefaultConfig()
	config.AuthEmail = "CLOUDFLARE_EMAIL"
	config.AuthKey = "CLOUDFLARE_API_KEY"
	config.AuthToken = "CLOUDFLARE_DNS_API_TOKEN"
	config.ZoneToken = "CLOUDFLARE_ZONE_API_TOKEN"

	provider, err := cloudflare.NewDNSProviderConfig(config)
	if err != nil {
		return fmt.Errorf("init-http-provider[%s, %s]: %s", domain, conf.KeyType, err.Error())
	}

	cli.Challenge.SetDNS01Provider(provider)
	return nil
}
