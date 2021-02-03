package account

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"os"
	"path"

	"github.com/alphatr/acme-lego/common/errors"
)

// CreateAccount 创建用户账户
func CreateAccount(email string, rootDir string) (*Account, *errors.Error) {
	accountPath := path.Join(rootDir, "account")

	if err := os.MkdirAll(accountPath, 0755); err != nil {
		return nil, errors.NewError(errors.CommonMakeDirErrno, err, accountPath)
	}

	keyPath := path.Join(accountPath, "account.key")
	if _, err := os.Stat(keyPath); err == nil || os.IsExist(err) {
		return nil, errors.NewError(errors.CommonFileNotExistErrno, nil, keyPath)
	}

	configPath := path.Join(accountPath, "account.json")
	if _, err := os.Stat(configPath); err == nil || os.IsExist(err) {
		return nil, errors.NewError(errors.CommonFileNotExistErrno, nil, configPath)
	}

	secret, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, errors.NewError(errors.ModelAccGenerateKeyErrno, err)
	}

	return &Account{Email: email, secret: secret, path: accountPath}, nil
}
