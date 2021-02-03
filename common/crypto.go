package common

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/alphatr/acme-lego/common/errors"
)

// LoadPrivateKey 加载私钥
func LoadPrivateKey(file string) (crypto.PrivateKey, *errors.Error) {
	keyBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.NewError(errors.CommonFileReadErrno, err, file)
	}

	keyBlock, _ := pem.Decode(keyBytes)

	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		prevate, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err != nil {
			return nil, errors.NewError(errors.CommonParsePrivateErrno, err, "rsa")
		}

		return prevate, nil
	case "EC PRIVATE KEY":
		prevate, err := x509.ParseECPrivateKey(keyBlock.Bytes)
		if err != nil {
			return nil, errors.NewError(errors.CommonParsePrivateErrno, err, "ecc")
		}

		return prevate, nil
	}

	return nil, errors.NewError(errors.CommonUnknowBlockErrno, nil, keyBlock.Type)
}
