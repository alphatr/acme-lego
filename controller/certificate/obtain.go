package certificate

import (
	"path"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/urfave/cli/v2"

	"github.com/alphatr/acme-lego/common/bootstrap"
	"github.com/alphatr/acme-lego/common/config"
	"github.com/alphatr/acme-lego/common/errors"
	"github.com/alphatr/acme-lego/model/account"
	"github.com/alphatr/acme-lego/model/client"
)

// Obtain 获取域名证书
func Obtain(ctx *cli.Context) error {
	acc, err := account.GetAccount(config.Config.RootDir)
	if err != nil {
		err := errors.NewError(errors.ConGetAccountErrno, err)
		return cli.NewExitError(err.Error(), 301)
	}

	lego, err := client.NewClient(acc)
	if err != nil {
		err := errors.NewError(errors.ConInitClientErrno, err)
		return cli.NewExitError(err.Error(), 302)
	}

	domain := ctx.String("domain")
	httpPath := ctx.String("http-path")

	if len(domain) > 0 {
		conf, ok := config.Config.DomainGroup[domain]
		if !ok {
			if len(httpPath) == 0 {
				err := errors.NewError(errors.ConRequireParamErrno, nil, "http-path")
				return cli.NewExitError(err.Error(), 303)
			}

			conf = &config.DomainConf{
				Domains:   []string{domain},
				KeyType:   []certcrypto.KeyType{certcrypto.RSA2048},
				Challenge: "http-path",
				Options:   map[string]string{"public": httpPath},
			}
		}

		if err := obtainDomain(domain, lego, conf); err != nil {
			err := errors.NewError(errors.ConCertObtainDomainErrno, err, domain)
			return cli.NewExitError(err.Error(), 304)
		}

		bootstrap.Log.Infof("[success] request-certificate: %s\n", domain)
		return nil
	}

	for domain, conf := range config.Config.DomainGroup {
		if err := obtainDomain(domain, lego, conf); err != nil {
			err := errors.NewError(errors.ConCertObtainDomainErrno, err, domain)
			return cli.NewExitError(err.Error(), 304)
		}

		bootstrap.Log.Infof("[success] request-certificate: %s\n", domain)
	}

	return nil
}

func obtainDomain(domain string, cli *client.Client, conf *config.DomainConf) *errors.Error {
	for _, keyType := range conf.KeyType {
		if err := cli.SetupChallenge(conf.Challenge, domain, conf); err != nil {
			return errors.NewError(errors.ConCertSetupChallengeErrno, err)
		}

		secret, errs := certcrypto.GeneratePrivateKey(keyType)
		if errs != nil {
			return errors.NewError(errors.ConCertGenerateKeyErrno, errs, keyType)
		}

		cert, err := cli.CertificateObtain(conf.Domains, secret)
		if err != nil {
			return errors.NewError(errors.ConCertObtainErrno, err, domain, keyType)
		}

		certPath := path.Join(config.Config.RootDir, "certificates", domain)
		if err := checkFolder(certPath); err != nil {
			return errors.NewError(errors.ConCertCheckFolderErrno, err, domain)
		}

		if err := saveCertRes(cert, certPath, keyType); err != nil {
			return errors.NewError(errors.ConCertSaveCertErrno, err, domain, keyType)
		}
	}

	return nil
}
