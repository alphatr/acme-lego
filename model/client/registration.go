package client

import (
	"github.com/go-acme/lego/v3/registration"

	"github.com/alphatr/acme-lego/common/errors"
)

// AccountRegister 注册账户
func (cli *Client) AccountRegister() (*registration.Resource, *errors.Error) {
	reg, errs := cli.lego.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if errs != nil {
		return nil, errors.NewError(errors.ModelClientRegisterErrno, errs)
	}

	return reg, nil
}
