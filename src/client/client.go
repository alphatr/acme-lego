package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"alphatr.com/acme-lego/src/account"
	"alphatr.com/acme-lego/src/config"
	"github.com/xenolf/lego/acme"
	"github.com/xenolf/lego/providers/http/webroot"
)

func Run(key string, acc *account.Account, domainConf *config.DomainConfig) error {
	mainConf := config.Config

	fmt.Printf("%#v\n", domainConf.KeyType)
	for _, keyType := range domainConf.KeyType {
		cli, err := acme.NewClient(mainConf.AcmeURL, acc, keyType)
		if err != nil {
			return fmt.Errorf("Create Client Error: %s", err.Error())
		}

		if domainConf.Challenge == "http-path" {
			provider, err := webroot.NewHTTPProvider(domainConf.HTTPPath)
			if err != nil {
				return fmt.Errorf("Init HTTP Provider Error: %s", err.Error())
			}

			cli.SetChallengeProvider(acme.HTTP01, provider)
		}

		if domainConf.Challenge == "http-port" {
			cli.SetHTTPAddress(":8057")
		}

		cert, failures := cli.ObtainCertificate(domainConf.Domains, true, nil, true)
		if len(failures) > 0 {
			for key, value := range failures {
				fmt.Printf("[%s] Could not obtain certificates\n\t%s", key, value.Error())
			}

			os.Exit(1)
		}

		certPath := path.Join(mainConf.RootPath, "certificates", key)
		err = checkFolder(certPath)
		if err != nil {
			return fmt.Errorf("Could not check/create path: %s", err.Error())
		}

		err = saveCertRes(cert, certPath, keyType)
		if err != nil {
			return fmt.Errorf("Save Cert Res Error: %s", err.Error())
		}
	}

	return nil
}

func checkFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0700)
	}
	return nil
}

func saveCertRes(certRes acme.CertificateResource, certPath string, keyType acme.KeyType) error {
	keyTypeString := strings.ToLower(string(keyType))
	certOut := path.Join(certPath, keyTypeString+".crt")
	privOut := path.Join(certPath, keyTypeString+".key")
	metaOut := path.Join(certPath, keyTypeString+".json")
	issuerOut := path.Join(certPath, keyTypeString+".issuer.crt")

	err := ioutil.WriteFile(certOut, certRes.Certificate, 0600)
	if err != nil {
		return fmt.Errorf("[%s] Save Certificate Error: %s", certRes.Domain, err.Error())
	}

	if certRes.IssuerCertificate != nil {
		err = ioutil.WriteFile(issuerOut, certRes.IssuerCertificate, 0600)
		if err != nil {
			return fmt.Errorf("[%s] Save IssuerCertificate Error: %s", certRes.Domain, err.Error())
		}
	}

	// 提供 CSR 就不知道私钥了
	if certRes.PrivateKey != nil {
		err = ioutil.WriteFile(privOut, certRes.PrivateKey, 0600)
		if err != nil {
			return fmt.Errorf("[%s] Unable to save PrivateKey: %s", certRes.Domain, err.Error())
		}
	}

	jsonBytes, err := json.MarshalIndent(certRes, "", "\t")
	if err != nil {
		return fmt.Errorf("[%s] Unable to marshal CertResource: %s", certRes.Domain, err.Error())
	}

	err = ioutil.WriteFile(metaOut, jsonBytes, 0600)
	if err != nil {
		return fmt.Errorf("[%s] Unable to save CertResource: %s", certRes.Domain, err.Error())
	}

	return nil
}
