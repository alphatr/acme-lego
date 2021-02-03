package account

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/alphatr/acme-lego/common"
	"github.com/alphatr/acme-lego/common/errors"
)

// GetAccount 获取用户账户
func GetAccount(rootDir string) (*Account, *errors.Error) {
	accDir := path.Join(rootDir, "account")

	keyPath := path.Join(accDir, "account.key")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return nil, errors.NewError(errors.CommonFileNotExistErrno, nil, keyPath)
	}

	key, err := common.LoadPrivateKey(keyPath)
	if err != nil {
		return nil, errors.NewError(errors.ModelAccLoadPrivateErrno, err)
	}

	ecc, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.NewError(errors.ModelAccNotECDSAErrno, nil)
	}

	configPath := path.Join(accDir, "account.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.NewError(errors.CommonFileNotExistErrno, nil, configPath)
	}

	content, errs := ioutil.ReadFile(configPath)
	if errs != nil {
		return nil, errors.NewError(errors.CommonFileReadErrno, errs)
	}

	acc := Account{}
	if errs := json.Unmarshal(content, &acc); errs != nil {
		return nil, errors.NewError(errors.CommonJSONUnmarshalErrno, errs)
	}

	acc.path = accDir
	acc.secret = ecc
	return &acc, nil
}
