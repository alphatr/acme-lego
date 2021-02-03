package config

import (
	"strings"

	"github.com/go-acme/lego/v3/certcrypto"

	"github.com/alphatr/acme-lego/common"
	"github.com/alphatr/acme-lego/common/errors"
)

// DomainConf 单个域名的配置
type DomainConf struct {
	Domains   []string
	KeyType   []certcrypto.KeyType
	Challenge string
	Options   map[string]string
}

const defaultChallenge = "http-path"

func initDomainConfig(domain string, conf *domainTOML, types []certcrypto.KeyType, base *baseTOML) (*DomainConf, *errors.Error) {
	result := &DomainConf{
		Domains:   buildDomains(domain, conf.Domains),
		Challenge: common.DefaultString(common.DefaultString(conf.Challenge, base.Challenge), defaultChallenge),
		Options:   conf.Options,
	}

	keyType := []certcrypto.KeyType{certcrypto.RSA2048}
	domainKeyTypes := KeyTypeList(conf.KeyType)
	if len(domainKeyTypes) > 0 {
		keyType = domainKeyTypes
	} else if len(types) > 0 {
		keyType = types
	}

	result.KeyType = keyType
	return result, nil
}

func buildDomains(key string, list []string) []string {
	list = append([]string{key}, list...)

	result := []string{}
	exists := map[string]bool{}

	for _, domain := range list {
		domain = strings.ToLower(domain)
		if !exists[domain] {
			result = append(result, domain)
			exists[domain] = true
		}
	}

	return result
}
