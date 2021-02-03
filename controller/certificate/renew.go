package certificate

import (
	"io/ioutil"
	"path"
	"time"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/urfave/cli/v2"

	"github.com/alphatr/acme-lego/common"
	"github.com/alphatr/acme-lego/common/bootstrap"
	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
	"github.com/alphatr/acme-lego/model/account"
	"github.com/alphatr/acme-lego/model/client"
)

// Renew 续期域名证书
func Renew(ctx *cli.Context) error {
	acc, err := account.GetAccount(config.Config.RootDir)
	if err != nil {
		err := errors.NewError(errors.ConGetAccountErrno, err)
		return cli.NewExitError(err.Error(), 401)
	}

	lego, err := client.NewClient(acc)
	if err != nil {
		err := errors.NewError(errors.ConInitClientErrno, err)
		return cli.NewExitError(err.Error(), 402)
	}

	domain := ctx.String("domain")

	if len(domain) > 0 {
		conf, ok := config.Config.DomainGroup[domain]
		if !ok {
			err := errors.NewError(errors.ConErrorParamErrno, nil, "domain")
			return cli.NewExitError(err.Error(), 403)
		}

		if err := renewDomain(domain, lego, conf); err != nil {
			err := errors.NewError(errors.ConCertRenewDomainErrno, err, domain)
			return cli.NewExitError(err.Error(), 404)
		}

		bootstrap.Log.Infof("[success] renew-certificate: %s\n", domain)
	} else {
		for domain, conf := range config.Config.DomainGroup {
			if err := renewDomain(domain, lego, conf); err != nil {
				err := errors.NewError(errors.ConCertRenewDomainErrno, err, domain)
				return cli.NewExitError(err.Error(), 404)
			}

			bootstrap.Log.Infof("[success] renew-certificate: %s\n", domain)
		}
	}

	if len(config.Config.AfterRenew) > 0 {
		result, err := common.RunCommand(config.Config.AfterRenew)
		if err != nil {
			return errors.NewError(errors.ConCertRunAfterRenewErrno, err)
		}

		bootstrap.Log.Debugf("after-renew: %s\n", result)
	}

	return nil
}

func renewDomain(domain string, cli *client.Client, conf *config.DomainConf) *errors.Error {
	for _, keyType := range conf.KeyType {
		if err := cli.SetupChallenge(conf.Challenge, domain, conf); err != nil {
			return errors.NewError(errors.ConCertSetupChallengeErrno, err)
		}

		certPath := path.Join(config.Config.RootDir, "certificates", domain)
		files := generateFilePath(certPath, keyType)

		content, errs := ioutil.ReadFile(files.Cert)
		if errs != nil {
			return errors.NewError(errors.CommonFileReadErrno, errs, files.Cert)
		}

		cert, errs := certcrypto.ParsePEMCertificate(content)
		if errs != nil {
			return errors.NewError(errors.CommonParseCertificateErrno, errs, domain, keyType)
		}

		if cert.NotAfter.After(time.Now().Add(config.Config.Expires)) {
			bootstrap.Log.Debugf("cert-not-expires(%s)", domain)
			return nil
		}

		privateKey, err := common.LoadPrivateKey(files.Prev)
		if err != nil {
			return errors.NewError(errors.ConCertLoadPrivateErrno, err, domain, keyType)
		}

		newCert, errs := cli.CertificateObtain(conf.Domains, privateKey)
		if errs != nil {
			return errors.NewError(errors.ConCertObtainErrno, errs, domain, keyType)
		}

		if err := checkFolder(certPath); err != nil {
			return errors.NewError(errors.ConCertCheckFolderErrno, err, domain)
		}

		if err := saveCertRes(newCert, certPath, keyType); err != nil {
			return errors.NewError(errors.ConCertSaveCertErrno, err, domain, keyType)
		}
	}

	return nil
}
