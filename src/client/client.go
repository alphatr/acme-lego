package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"alphatr.com/acme-lego/src/account"
	"alphatr.com/acme-lego/src/config"
	acme "github.com/xenolf/lego/acmev2"
	"github.com/xenolf/lego/providers/http/webroot"
)

type certFilePath struct {
	Cert   string
	Prev   string
	Meta   string
	Issuer string
}

func Run(domainKey string, acc *account.Account, domainConf *config.DomainConfig) error {
	mainConf := config.Config

	for _, keyType := range domainConf.KeyType {
		cli, err := acme.NewClient(mainConf.AcmeURL, acc, keyType)
		if err != nil {
			return fmt.Errorf("init-client[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		if domainConf.Challenge == "http-path" {
			provider, err := webroot.NewHTTPProvider(domainConf.HTTPPath)
			if err != nil {
				return fmt.Errorf("init-http-provider[%s, %s]: %s", domainKey, keyType, err.Error())
			}

			cli.SetChallengeProvider(acme.HTTP01, provider)
		}

		if domainConf.Challenge == "http-port" {
			cli.SetHTTPAddress(":8057")
		}

		cert, failures := cli.ObtainCertificate(domainConf.Domains, true, nil, true)
		if len(failures) > 0 {
			for _, value := range failures {
				return fmt.Errorf("obtain-certificate[%s, %s]: %s", domainKey, keyType, value.Error())
			}
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
		cli, err := acme.NewClient(mainConf.AcmeURL, acc, keyType)
		if err != nil {
			return fmt.Errorf("init-client[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		if domainConf.Challenge == "http-path" {
			provider, err := webroot.NewHTTPProvider(domainConf.HTTPPath)
			if err != nil {
				return fmt.Errorf("init-http-provider[%s, %s]: %s", domainKey, keyType, err.Error())
			}

			cli.SetChallengeProvider(acme.HTTP01, provider)
		}

		if domainConf.Challenge == "http-port" {
			cli.SetHTTPAddress(":8057")
		}

		certPath := path.Join(mainConf.RootPath, "certificates", domainKey)
		certFiles := generateFilePath(certPath, keyType)

		certBytes, err := ioutil.ReadFile(certFiles.Cert)
		if err != nil {
			return fmt.Errorf("read-cert-file[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		expTime, err := acme.GetPEMCertExpiration(certBytes)
		if err != nil {
			return fmt.Errorf("get-cert-expire[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		if expTime.Sub(time.Now()) > mainConf.Expires {
			// 时间未到
			return nil
		}

		metaBytes, err := ioutil.ReadFile(certFiles.Meta)
		if err != nil {
			return fmt.Errorf("read-meta-file[%s, %s]: %s", domainKey, keyType, err.Error())
		}

		var certRes acme.CertificateResource
		err = json.Unmarshal(metaBytes, &certRes)
		if err != nil {
			return fmt.Errorf("unmarshal-meta-info[%s, %s]: %s", domainKey, keyType, err.Error())
		}
		certRes.Certificate = certBytes

		keyBytes, err := ioutil.ReadFile(certFiles.Prev)
		if err != nil {
			return fmt.Errorf("read-privatekey-file[%s, %s]: %s", domainKey, keyType, err.Error())
		}
		certRes.PrivateKey = keyBytes

		newCert, err := cli.RenewCertificate(certRes, true, true)
		if err != nil {
			return fmt.Errorf("renew-certificate[%s, %s]: %s", domainKey, keyType, err.Error())
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

	return nil
}

func checkFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0700)
	}
	return nil
}

func saveCertRes(certRes acme.CertificateResource, certPath string, keyType acme.KeyType) error {
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

func generateFilePath(certPath string, keyType acme.KeyType) *certFilePath {
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
