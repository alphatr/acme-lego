package challenge

import (
	"github.com/go-acme/lego/v3/challenge"

	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
)

// ProviderType 类型
type ProviderType string

// ProviderType 支持的三种类型
const (
	ProviderDNS  ProviderType = "dns"
	ProviderHTTP ProviderType = "http"
	ProviderTLS  ProviderType = "tls"
)

// Provider 解决方案
type Provider interface {
	Type() ProviderType
	Provider(string, *config.DomainConf) (challenge.Provider, *errors.Error)
}

// ProviderMap ProviderMap
var ProviderMap = map[string]Provider{}
