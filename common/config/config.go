package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"

	"github.com/alphatr/acme-lego/common/errors"
)

type domainTOML struct {
	Domains   []string          `toml:"domains"`
	KeyType   []string          `toml:"key-type"`
	Challenge string            `toml:"challenge"`
	Options   map[string]string `toml:"options"`
}

type baseTOML struct {
	Dev         bool                  `toml:"dev"`
	RootDir     string                `toml:"root-dir"`
	AcmeURL     string                `toml:"acme-url"`
	LogLevel    string                `toml:"log-level"`
	Email       string                `toml:"email"`
	KeyType     []string              `toml:"key-type"`
	Challenge   string                `toml:"challenge"`
	DomainGroup map[string]domainTOML `toml:"domain-group"`
	ExpireDays  int                   `toml:"expire-days"`
	AfterRenew  string                `toml:"after-renew"`
}

// InitConfig 配置初始化
func InitConfig(configPath string) *errors.Error {
	configPath, errs := filepath.Abs(configPath)
	if errs != nil {
		return errors.NewError(errors.CommonPathAbsErrno, errs, configPath)
	}

	conf, err := parseTOML(configPath)
	if err != nil {
		return errors.NewError(errors.ConfigParseTOMLErrno, err)
	}

	if err := initBaseConfig(conf, configPath); err != nil {
		return errors.NewError(errors.ConfigBaseInitErrno, err)
	}

	return nil
}

func parseTOML(configPath string) (*baseTOML, *errors.Error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.NewError(errors.CommonFileNotExistErrno, nil, configPath)
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.NewError(errors.CommonFileReadErrno, err, configPath)
	}

	var conf baseTOML
	if err := toml.Unmarshal(configBytes, &conf); err != nil {
		return nil, errors.NewError(errors.CommonTOMLUnmarshalErrno, err)
	}

	return &conf, nil
}
