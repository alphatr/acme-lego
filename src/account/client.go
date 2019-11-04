package account

import (
	"fmt"
	"time"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/lego"

	"alphatr.com/acme-lego/src/config"
)

func NewClient(acc *Account) (*lego.Client, error) {
	conf := lego.NewConfig(acc)
	conf.CADirURL = config.Config.AcmeURL
	conf.UserAgent = fmt.Sprintf("alphatr-lego-cli")
	conf.Certificate = lego.CertificateConfig{
		KeyType: certcrypto.RSA2048,
		Timeout: config.Config.HTTPTimeout * time.Second,
	}
	conf.HTTPClient.Timeout = config.Config.HTTPTimeout * time.Second

	client, err := lego.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("init-client: %s", err.Error())
	}

	return client, nil
}
