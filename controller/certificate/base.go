package certificate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/certificate"

	"github.com/alphatr/acme-lego/common/errors"
)

type certFilePath struct {
	Cert   string
	Prev   string
	Meta   string
	Issuer string
}

func checkFolder(path string) *errors.Error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0700); err != nil {
			return errors.NewError(errors.CommonMakeDirErrno, err, path)
		}
	}
	return nil
}

func saveCertRes(certRes *certificate.Resource, certPath string, keyType certcrypto.KeyType) *errors.Error {
	files := generateFilePath(certPath, keyType)

	if err := ioutil.WriteFile(files.Cert, certRes.Certificate, 0600); err != nil {
		return errors.NewError(errors.CommonFileWriteErrno, err, files.Cert)
	}

	if certRes.IssuerCertificate != nil {
		if err := ioutil.WriteFile(files.Issuer, certRes.IssuerCertificate, 0600); err != nil {
			return errors.NewError(errors.CommonFileWriteErrno, err, files.Issuer)
		}
	}

	// 提供 CSR 就不知道私钥了
	if certRes.PrivateKey != nil {
		if err := ioutil.WriteFile(files.Prev, certRes.PrivateKey, 0600); err != nil {
			return errors.NewError(errors.CommonFileWriteErrno, err, files.Prev)
		}
	}

	content, err := json.MarshalIndent(certRes, "", "\t")
	if err != nil {
		return errors.NewError(errors.CommonJSONMarshalErrno, err)
	}

	if err := ioutil.WriteFile(files.Meta, content, 0600); err != nil {
		return errors.NewError(errors.CommonFileWriteErrno, err, files.Meta)
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
