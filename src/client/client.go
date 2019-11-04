package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/certificate"
	"github.com/go-acme/lego/v3/lego"

	"alphatr.com/acme-lego/src/account"
	"alphatr.com/acme-lego/src/config"
	"alphatr.com/acme-lego/src/misc"
)

type certFilePath struct {
	Cert   string
	Prev   string
	Meta   string
	Issuer string
}

// ProviderMap ProviderMap
var ProviderMap = map[string]func(string, *lego.Client, *config.DomainConfig) error{}

func Run(domainKey string, acc *account.Account, domainConf *config.DomainConfig) error {
	mainConf := config.Config

	for _, keyType := range domainConf.KeyType {
		cli, err := account.NewClient(acc)
		if err != nil {
			return fmt.Errorf("init-client[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		provider, ok := ProviderMap[domainConf.Challenge]
		if !ok {
			return fmt.Errorf("unknow-provider[%s]: %s", domainKey, domainConf.Challenge)
		}

		if err := provider(domainKey, cli, domainConf); err != nil {
			return fmt.Errorf("apply-provider-error[%s, %s]: %s", domainKey, domainConf.Challenge, err.Error())
		}

		request := certificate.ObtainRequest{
			Domains:    domainConf.Domains,
			Bundle:     true,
			MustStaple: true,
		}

		cert, err := cli.Certificate.Obtain(request)
		if err != nil {
			return fmt.Errorf("obtain-certificate[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		certPath := path.Join(mainConf.RootPath, "certificates", domainKey)
		err = checkFolder(certPath)
		if err != nil {
			return fmt.Errorf("check-folder[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		err = saveCertRes(cert, certPath, keyType)
		if err != nil {
			return fmt.Errorf("save-cert-res[%s, %s]: %s", domainKey, keyType, err.Error())
		}
	}

	return nil
}

func Renew(domainKey string, acc *account.Account, domainConf *config.DomainConfig) error {
	mainConf := config.Config

	for _, keyType := range domainConf.KeyType {
		cli, err := account.NewClient(acc)
		if err != nil {
			return fmt.Errorf("init-client[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		provider, ok := ProviderMap[domainConf.Challenge]
		if !ok {
			return fmt.Errorf("unknow-provider[%s]: %s", domainKey, domainConf.Challenge)
		}

		if err := provider(domainKey, cli, domainConf); err != nil {
			return fmt.Errorf("apply-provider-error[%s, %s]: %s", domainKey, domainConf.Challenge, err.Error())
		}

		certPath := path.Join(mainConf.RootPath, "certificates", domainKey)
		certFiles := generateFilePath(certPath, keyType)

		certBytes, err := ioutil.ReadFile(certFiles.Cert)
		if err != nil {
			return fmt.Errorf("read-cert-file[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		cert, err := certcrypto.ParsePEMCertificate(certBytes)
		if err != nil {
			return fmt.Errorf("parse-pem-certificate[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		if cert.NotAfter.Sub(time.Now()) > mainConf.Expires {
			// 时间未到
			return nil
		}

		keyBytes, err := ioutil.ReadFile(certFiles.Prev)
		if err != nil {
			return fmt.Errorf("read-privatekey-file[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		request := certificate.ObtainRequest{
			Domains:    domainConf.Domains,
			PrivateKey: keyBytes,
			Bundle:     true,
			MustStaple: true,
		}

		newCert, err := cli.Certificate.Obtain(request)
		if err != nil {
			return fmt.Errorf("renew-obtain-certificate[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		err = checkFolder(certPath)
		if err != nil {
			return fmt.Errorf("check-folder[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		err = saveCertRes(newCert, certPath, keyType)
		if err != nil {
			return fmt.Errorf("save-cert-res[%s, %s]: %s", domainKey, keyType, err.Error())
		}
	}

	if len(mainConf.AfterRenew) > 0 {
		result, err := misc.RunCmd(mainConf.AfterRenew)
		if err != nil {
			return fmt.Errorf("misc-runcmd: %s", err.Error())
		}

		fmt.Println(result)
	}

	return nil
}

func checkFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0700)
	}
	return nil
}

func saveCertRes(certRes *certificate.Resource, certPath string, keyType certcrypto.KeyType) error {
	certFiles := generateFilePath(certPath, keyType)

	err := ioutil.WriteFile(certFiles.Cert, certRes.Certificate, 0600)
	if err != nil {
		return fmt.Errorf("write-cert-file: %s", err.Error())
	}

	if certRes.IssuerCertificate != nil {
		err = ioutil.WriteFile(certFiles.Issuer, certRes.IssuerCertificate, 0600)
		if err != nil {
			return fmt.Errorf("write-issuer-file: %s", err.Error())
		}
	}

	// 提供 CSR 就不知道私钥了
	if certRes.PrivateKey != nil {
		err = ioutil.WriteFile(certFiles.Prev, certRes.PrivateKey, 0600)
		if err != nil {
			return fmt.Errorf("write-privatekey-file: %s", err.Error())
		}
	}

	jsonBytes, err := json.MarshalIndent(certRes, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal-meta-info: %s", err.Error())
	}

	err = ioutil.WriteFile(certFiles.Meta, jsonBytes, 0600)
	if err != nil {
		return fmt.Errorf("write-meta-file: %s", err.Error())
	}

	return nil
}

func generateFilePath(certPath string, keyType certcrypto.KeyType) *certFilePath {
	keyTypeMap := map[string]string{
		"p256": "ecdsa-256",
		"p384": "ecdsa-384",
		"2048": "rsa-2048",
		"4096": "rsa-4096",
		"8192": "rsa-8192",
	}

	keyTStr := keyTypeMap[strings.ToLower(string(keyType))]

	return &certFilePath{
		Cert:   path.Join(certPath, fmt.Sprintf("fullchain.%s.crt", keyTStr)),
		Prev:   path.Join(certPath, fmt.Sprintf("privkey.%s.key", keyTStr)),
		Meta:   path.Join(certPath, fmt.Sprintf("meta.%s.json", keyTStr)),
		Issuer: path.Join(certPath, fmt.Sprintf("issuer.%s.crt", keyTStr)),
	}
}
