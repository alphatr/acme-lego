package account

import (
	"github.com/urfave/cli/v2"

	"github.com/alphatr/acme-lego/common"
	"github.com/alphatr/acme-lego/common/bootstrap"
	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
	"github.com/alphatr/acme-lego/model/account"
	"github.com/alphatr/acme-lego/model/client"
)

// Register 注册账号
func Register(ctx *cli.Context) error {
	mail := common.DefaultString(ctx.String("mail"), config.Config.Email)

	acc, err := account.CreateAccount(mail, config.Config.RootDir)
	if err != nil {
		err := errors.NewError(errors.ConAccCreateErrno, err)
		return cli.NewExitError(err.Error(), 201)
	}

	client, err := client.NewClient(acc)
	if err != nil {
		err := errors.NewError(errors.ConInitClientErrno, err)
		return cli.NewExitError(err.Error(), 202)
	}

	reg, err := client.AccountRegister()
	if err != nil {
		err := errors.NewError(errors.ConAccRegisterErrno, err)
		return cli.NewExitError(err.Error(), 203)
	}

	acc.Registration = reg
	if err := acc.Save(); err != nil {
		err := errors.NewError(errors.ConAccSaveErrno, err)
		return cli.NewExitError(err.Error(), 203)
	}

	bootstrap.Log.Infof("[success] registering-account: %s\n", acc.Email)
	return nil
}
