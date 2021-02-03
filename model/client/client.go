package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/log"

	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
	"github.com/alphatr/acme-lego/model/account"
)

// Client 客户端
type Client struct {
	lego    *lego.Client
	config  *lego.Config
	account *account.Account
}

// NewClient 创建新客户端
func NewClient(acc *account.Account) (*Client, *errors.Error) {
	log.Logger = &clientLogger{}

	conf := lego.NewConfig(acc)
	conf.CADirURL = config.Config.AcmeURL
	conf.UserAgent = config.Config.UserAgent
	conf.Certificate = lego.CertificateConfig{
		KeyType: certcrypto.RSA2048,
		Timeout: config.Config.HTTPTimeout * time.Second,
	}

	conf.HTTPClient.Timeout = config.Config.HTTPTimeout * time.Second
	if config.Config.Dev {
		conf.HTTPClient.Transport = createInsecureTransport()
	}

	client, err := lego.NewClient(conf)
	if err != nil {
		return nil, errors.NewError(errors.ModelClientInitErrno, err)
	}

	return &Client{lego: client}, nil
}

func createInsecureTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}
