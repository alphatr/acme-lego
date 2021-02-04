package client

import (
	"github.com/go-acme/lego/v3/challenge"

	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
	chall "github.com/alphatr/acme-lego/model/challenge"
)

// SetupChallenge 初始化 Challenge
func (cli *Client) SetupChallenge(name string, domain string, conf *config.DomainConf) *errors.Error {
	item, ok := chall.ProviderMap[name]
	if !ok {
		return errors.NewError(errors.ModelClientUnknowProviderErrno, nil, name)
	}

	provider, err := item.Provider(domain, conf)
	if err != nil {
		return errors.NewError(errors.ModelClientProviderErrno, err)
	}

	if err := cli.setProvider(item.Type(), provider); err != nil {
		return errors.NewError(errors.ModelClientSetProviderErrno, err)
	}

	return nil
}

func (cli *Client) setProvider(input chall.ProviderType, provider challenge.Provider) *errors.Error {
	switch input {
	case chall.ProviderHTTP:
		if err := cli.lego.Challenge.SetHTTP01Provider(provider); err != nil {
			return errors.NewError(errors.ModelClientTypeProviderErrno, err, input)
		}

	case chall.ProviderTLS:
		if err := cli.lego.Challenge.SetTLSALPN01Provider(provider); err != nil {
			return errors.NewError(errors.ModelClientTypeProviderErrno, err, input)
		}

	case chall.ProviderDNS:
		if err := cli.lego.Challenge.SetDNS01Provider(provider); err != nil {
			return errors.NewError(errors.ModelClientTypeProviderErrno, err, input)
		}
	}

	return nil
}
