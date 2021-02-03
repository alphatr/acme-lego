package client

import (
	"crypto"

	"github.com/go-acme/lego/v3/certificate"

	"github.com/alphatr/acme-lego/common/errors"
)

// CertificateObtain 证书获取
func (cli *Client) CertificateObtain(domains []string, secret crypto.PrivateKey) (*certificate.Resource, *errors.Error) {
	request := certificate.ObtainRequest{
		Domains:    domains,
		PrivateKey: secret,
		Bundle:     true,
		MustStaple: true,
	}

	cert, errs := cli.lego.Certificate.Obtain(request)
	if errs != nil {
		return nil, errors.NewError(errors.ModelClientObtainErrno, errs)
	}

	return cert, nil
}
