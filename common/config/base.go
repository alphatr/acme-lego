package config

import (
	"path"
	"strings"
	"time"

	"github.com/go-acme/lego/v3/certcrypto"

	"github.com/alphatr/acme-lego/common"
	"github.com/alphatr/acme-lego/common/errors"
)

const defaultAcmeURL = "https://acme-v02.api.letsencrypt.org/directory"

// BaseConf 配置
type BaseConf struct {
	Name        string
	Dev         bool
	LogLevel    string
	Email       string
	UserAgent   string
	DomainGroup map[string]*DomainConf
	HTTPTimeout time.Duration
	AcmeURL     string
	RootDir     string
	Expires     time.Duration
	AfterRenew  string
}

// Config 配置
var Config BaseConf

func initBaseConfig(conf *baseTOML, configPath string) *errors.Error {
	Config.Name = "alphatr-lego"
	Config.Dev = conf.Dev
	Config.LogLevel = conf.LogLevel
	Config.Email = conf.Email
	Config.HTTPTimeout = 30
	Config.Expires = time.Duration(conf.ExpireDays) * time.Hour * 24
	Config.AfterRenew = conf.AfterRenew
	Config.RootDir = common.DefaultString(conf.RootDir, path.Dir(configPath))

	Config.AcmeURL = defaultAcmeURL
	if conf.Dev {
		Config.AcmeURL = common.DefaultString(conf.AcmeURL, defaultAcmeURL)
	}

	types := KeyTypeList(conf.KeyType)

	domainGroup := map[string]*DomainConf{}
	for domain, value := range conf.DomainGroup {
		conf, err := initDomainConfig(domain, &value, types, conf)
		if err != nil {
			return errors.NewError(errors.ConfigDomainInitErrno, err)
		}

		domainGroup[strings.ToLower(domain)] = conf
	}

	Config.DomainGroup = domainGroup
	return nil
}

// KeyTypeList 返回证书类型列表
func KeyTypeList(input []string) []certcrypto.KeyType {
	getKeyType := func(input string) (certcrypto.KeyType, bool) {
		keysMap := map[string]certcrypto.KeyType{
			"EC256":   certcrypto.EC256,
			"EC384":   certcrypto.EC384,
			"RSA2048": certcrypto.RSA2048,
			"RSA4096": certcrypto.RSA4096,
			"RSA8192": certcrypto.RSA8192,
		}

		if res, ok := keysMap[strings.ToUpper(input)]; ok {
			return res, true
		}

		return certcrypto.RSA2048, false
	}

	result := []certcrypto.KeyType{}

	for _, value := range input {
		if item, ok := getKeyType(value); ok {
			result = append(result, item)
		}
	}

	return result
}
