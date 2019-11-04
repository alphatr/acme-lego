package account

import (
	"crypto"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/go-acme/lego/v3/registration"
)

type Account struct {
	Email        string                 `json:"email"`
	Registration *registration.Resource `json:"registration"`
	key          crypto.PrivateKey
	path         string
}

func (acc Account) GetEmail() string {
	return acc.Email
}

func (acc Account) GetRegistration() *registration.Resource {
	return acc.Registration
}

func (acc Account) GetPrivateKey() crypto.PrivateKey {
	return acc.key
}

func (acc Account) Save() error {
	jsonBytes, err := json.MarshalIndent(acc, "", "    ")
	if err != nil {
		return err
	}

	accountFile := path.Join(acc.path, "account.json")
	return ioutil.WriteFile(accountFile, jsonBytes, 0600)
}
