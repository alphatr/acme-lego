package challenge

import (
	"github.com/go-acme/lego/v3/challenge"
	"github.com/go-acme/lego/v3/providers/dns/cloudflare"

	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
)

func init() {
	ProviderMap["dns-cloudflare"] = &DNSCloudflareProvider{}
}

// DNSCloudflareProvider Cloudflare DNS
type DNSCloudflareProvider struct {
}

// Type 返回注册的类型
func (ins *DNSCloudflareProvider) Type() []ProviderType {
	return []ProviderType{ProviderDNS}
}

// Provider Provider 实体
func (ins *DNSCloudflareProvider) Provider(domain string, conf *config.DomainConf) (challenge.Provider, *errors.Error) {
	config := cloudflare.NewDefaultConfig()
	config.AuthToken = conf.Options["token"]

	provider, err := cloudflare.NewDNSProviderConfig(config)
	if err != nil {
		return nil, errors.NewError(errors.ModelChalDNSConfigErrno, err, "cloudflare")
	}

	return provider, nil
}
