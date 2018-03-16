package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	acme "github.com/xenolf/lego/acmev2"
)

type domainJSON struct {
	Domains   []string `json:"domains"`
	KeyType   []string `json:"key-type"`
	Challenge string   `json:"challenge"`
	HTTPPath  string   `json:"http-path"`
}

type rootJSON struct {
	Email       string                `json:"email"`
	KeyType     []string              `json:"key-type"`
	Challenge   string                `json:"challenge"`
	DomainGroup map[string]domainJSON `json:"domain-group"`
	AcmeURL     string                `json:"acme-url"`
}

type DomainConfig struct {
	Domains   []string
	KeyType   []acme.KeyType
	Challenge string
	HTTPPath  string
}

type MainConfig struct {
	Email       string
	DomainGroup map[string]*DomainConfig
	HTTPTimeout time.Duration
	DNSTimeout  time.Duration
	AcmeURL     string
	RootPath    string
}

var Config MainConfig

func InitConfig(rootPath string) error {
	jsonConfig, err := parseJSON(rootPath)
	if err != nil {
		return err
	}

	jsonKeyType := parseKeyTypes(jsonConfig.KeyType)
	domainGroup := map[string]*DomainConfig{}

	for key, value := range jsonConfig.DomainGroup {
		challenge := "http-file"
		if len(value.Challenge) > 0 {
			challenge = value.Challenge
		} else if len(jsonConfig.Challenge) > 0 {
			challenge = jsonConfig.Challenge
		}

		keyType := []acme.KeyType{acme.RSA2048}
		valueKeyType := parseKeyTypes(value.KeyType)
		if len(valueKeyType) > 0 {
			keyType = valueKeyType
		} else if len(jsonKeyType) > 0 {
			keyType = jsonKeyType
		}

		domainConf := &DomainConfig{
			Domains:   makeSet(key, value.Domains),
			KeyType:   keyType,
			Challenge: challenge,
			HTTPPath:  value.HTTPPath,
		}

		domainGroup[strings.ToLower(key)] = domainConf
	}

	Config.Email = jsonConfig.Email
	Config.DomainGroup = domainGroup
	Config.HTTPTimeout = 30
	Config.DNSTimeout = 10
	if len(jsonConfig.AcmeURL) > 0 {
		Config.AcmeURL = jsonConfig.AcmeURL
	} else {
		Config.AcmeURL = "https://acme-v02.api.letsencrypt.org/directory"
	}

	Config.RootPath = rootPath
	return nil
}

func parseJSON(rootPath string) (*rootJSON, error) {
	configPath := path.Join(rootPath, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Config file is not exist: %s", configPath)
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Could not load file for account %s: %s", configPath, err.Error())
	}

	var jsonConfig rootJSON
	err = json.Unmarshal(configBytes, &jsonConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not parse file for account %s: %s", configPath, err.Error())
	}

	return &jsonConfig, nil
}

func parseKeyTypes(typeString []string) []acme.KeyType {
	var typeList []acme.KeyType
	for _, value := range typeString {
		keyType, ok := stringToKeyType(value)
		if ok {
			typeList = append(typeList, keyType)
		}
	}

	return typeList
}

func stringToKeyType(typeString string) (acme.KeyType, bool) {
	typeString = strings.ToUpper(typeString)
	switch typeString {
	case "EC256":
		return acme.EC256, true
	case "EC384":
		return acme.EC384, true
	case "RSA2048":
		return acme.RSA2048, true
	case "RSA4096":
		return acme.RSA4096, true
	case "RSA8192":
		return acme.RSA8192, true
	}

	return acme.RSA2048, false
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
