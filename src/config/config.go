package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/pelletier/go-toml"
)

type domainTOML struct {
	Domains   []string          `toml:"domains"`
	KeyType   []string          `toml:"key-type"`
	Challenge string            `toml:"challenge"`
	Options   map[string]string `toml:"options"`
}

type mainTOML struct {
	Email       string                `toml:"email"`
	KeyType     []string              `toml:"key-type"`
	Challenge   string                `toml:"challenge"`
	DomainGroup map[string]domainTOML `toml:"domain-group"`
	AcmeURL     string                `toml:"acme-url"`
	ExpireDays  int                   `toml:"expire-days"`
	AfterRenew  string                `toml:"after-renew"`
}

// DomainConfig 单个域名的配置
type DomainConfig struct {
	Domains   []string
	KeyType   []certcrypto.KeyType
	Challenge string
	Options   map[string]string
}

// MainConfig 主配置
type MainConfig struct {
	Email       string
	DomainGroup map[string]*DomainConfig
	HTTPTimeout time.Duration
	AcmeURL     string
	RootPath    string
	Expires     time.Duration
	AfterRenew  string
}

// Config 主配置
var Config MainConfig

// InitConfig 配置初始化
func InitConfig(rootPath string) error {
	conf, err := parseTOML(rootPath)
	if err != nil {
		return err
	}

	mainKeyType := parseKeyTypes(conf.KeyType)
	domainGroup := map[string]*DomainConfig{}

	for key, value := range conf.DomainGroup {
		challenge := "http-path"
		if len(value.Challenge) > 0 {
			challenge = value.Challenge
		} else if len(conf.Challenge) > 0 {
			challenge = conf.Challenge
		}

		keyType := []certcrypto.KeyType{certcrypto.RSA2048}
		valueKeyType := parseKeyTypes(value.KeyType)
		if len(valueKeyType) > 0 {
			keyType = valueKeyType
		} else if len(mainKeyType) > 0 {
			keyType = mainKeyType
		}

		domainConf := &DomainConfig{
			Domains:   makeSet(key, value.Domains),
			KeyType:   keyType,
			Challenge: challenge,
			Options:   value.Options,
		}

		domainGroup[strings.ToLower(key)] = domainConf
	}

	Config.Email = conf.Email
	Config.DomainGroup = domainGroup
	Config.HTTPTimeout = 30
	Config.Expires = time.Duration(conf.ExpireDays) * time.Hour * 24
	Config.AfterRenew = conf.AfterRenew

	if len(conf.AcmeURL) > 0 {
		Config.AcmeURL = conf.AcmeURL
	} else {
		Config.AcmeURL = "https://acme-v02.api.letsencrypt.org/directory"
	}

	Config.RootPath = rootPath
	return nil
}

func parseTOML(rootPath string) (*mainTOML, error) {
	configPath := path.Join(rootPath, "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Config file is not exist: %s", configPath)
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Could not load file for account %s: %s", configPath, err.Error())
	}

	var conf mainTOML
	if err := toml.Unmarshal(configBytes, &conf); err != nil {
		return nil, fmt.Errorf("Could not parse file for account %s: %s", configPath, err.Error())
	}

	return &conf, nil
}

func parseKeyTypes(typeString []string) []certcrypto.KeyType {
	var typeList []certcrypto.KeyType
	for _, value := range typeString {
		keyType, ok := stringToKeyType(value)
		if ok {
			typeList = append(typeList, keyType)
		}
	}

	return typeList
}

func stringToKeyType(typeString string) (certcrypto.KeyType, bool) {
	typeString = strings.ToUpper(typeString)
	switch typeString {
	case "EC256":
		return certcrypto.EC256, true
	case "EC384":
		return certcrypto.EC384, true
	case "RSA2048":
		return certcrypto.RSA2048, true
	case "RSA4096":
		return certcrypto.RSA4096, true
	case "RSA8192":
		return certcrypto.RSA8192, true
	}

	return certcrypto.RSA2048, false
}

func makeSet(keyDomain string, domainList []string) []string {
	retList := []string{}
	tempMap := map[string]bool{}

	insertToMap := func(domain string) {
		if len(domain) > 0 {
			tempMap[strings.ToLower(domain)] = true
		}
	}

	insertToMap(keyDomain)
	for _, domain := range domainList {
		insertToMap(domain)
	}

	for domain, exist := range tempMap {
		if exist {
			retList = append(retList, domain)
		}
	}

	return retList
}
