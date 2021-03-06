package challenge

import (
	"github.com/go-acme/lego/v3/challenge"
	"github.com/go-acme/lego/v3/providers/http/webroot"

	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
)

func init() {
	ProviderMap["http-path"] = &HTTPPathProvider{isHTTPS: false}
	ProviderMap["https-path"] = &HTTPPathProvider{isHTTPS: true}
}

// HTTPPathProvider HTTPPathProvider
type HTTPPathProvider struct {
	isHTTPS bool
}

// Type 返回注册的类型
func (ins *HTTPPathProvider) Type() ProviderType {
	if ins.isHTTPS {
		return ProviderTLS
	}

	return ProviderHTTP
}

// Provider Provider 实体
func (ins *HTTPPathProvider) Provider(domain string, conf *config.DomainConf) (challenge.Provider, *errors.Error) {
	provider, err := webroot.NewHTTPProvider(conf.Options["public"])
	if err != nil {
		return nil, errors.NewError(errors.ModelChalHTTPInitErrno, err)
	}

	return provider, nil
}
