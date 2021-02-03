package account

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-acme/lego/v3/registration"

	"github.com/alphatr/acme-lego/common/errors"
)

// Account 用户账户
type Account struct {
	Email        string                 `json:"email"`
	Registration *registration.Resource `json:"registration"`
	secret       *ecdsa.PrivateKey
	path         string
}

// GetEmail 获取邮箱
func (acc *Account) GetEmail() string {
	return acc.Email
}

// GetRegistration 获取注册源
func (acc *Account) GetRegistration() *registration.Resource {
	return acc.Registration
}

// GetPrivateKey 获取用户私钥
func (acc *Account) GetPrivateKey() crypto.PrivateKey {
	return acc.secret
}

// Save 保存用户信息
func (acc *Account) Save() *errors.Error {
	if err := acc.saveAccountConfig(); err != nil {
		return errors.NewError(errors.ModelAccSaveConfigErrno, err)
	}

	if err := acc.saveAccountPrivateKey(); err != nil {
		return errors.NewError(errors.ModelAccSavePrivateErrno, err)
	}

	return nil
}

func (acc *Account) saveAccountConfig() *errors.Error {
	content, err := json.MarshalIndent(acc, "", "    ")
	if err != nil {
		return errors.NewError(errors.CommonJSONMarshalErrno, err)
	}

	accountFile := path.Join(acc.path, "account.json")
	if err := ioutil.WriteFile(accountFile, content, 0600); err != nil {
		return errors.NewError(errors.CommonFileWriteErrno, err)
	}

	return nil
}

func (acc *Account) saveAccountPrivateKey() *errors.Error {
	bytes, err := x509.MarshalECPrivateKey(acc.secret)
	if err != nil {
		return errors.NewError(errors.CommonMarshalPrivateErrno, err, "ecc")
	}

	accKeyPath := path.Join(acc.path, "account.key")
	block := pem.Block{Type: "EC PRIVATE KEY", Bytes: bytes}
	output, err := os.Create(accKeyPath)
	if err != nil {
		return errors.NewError(errors.CommonFileCreateErrno, err, accKeyPath)
	}

	pem.Encode(output, &block)
	if err := output.Close(); err != nil {
		return errors.NewError(errors.CommonFileCloseErrno, err, accKeyPath)
	}

	return nil
}
