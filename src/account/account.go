package account

import (
	"crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	acme "github.com/xenolf/lego/acmev2"

	"alphatr.com/acme-lego/src/config"
	"alphatr.com/acme-lego/src/misc"
)

type Account struct {
	Email        string                     `json:"email"`
	Registration *acme.RegistrationResource `json:"registration"`
	key          crypto.PrivateKey
	path         string
}

func GetAccount(rootPath string) (*Account, error) {
	accountPath := path.Join(rootPath, "account")

	accKeyPath := path.Join(accountPath, "account.key")
	if _, err := os.Stat(accKeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("key file is not exist[%s]", accKeyPath)
	}

	privateKey, err := misc.LoadPrivateKey(accKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load-private-key[%s]: %s", accKeyPath, err.Error())
	}

	accountFile := path.Join(accountPath, "account.json")
	if _, err := os.Stat(accountFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("account file is not exist[%s]", accountFile)
	}

	fileBytes, err := ioutil.ReadFile(accountFile)
	if err != nil {
		return nil, fmt.Errorf("read-file[%s]: %s", accountFile, err.Error())
	}

	var account Account
	err = json.Unmarshal(fileBytes, &account)
	if err != nil {
		return nil, fmt.Errorf("unmarshal-json[%s]: %s", accountFile, err.Error())
	}

	account.path = accountPath
	account.key = privateKey
	return &account, nil
}

func CreateAccount(email string, rootPath string) (*Account, error) {
	accountPath := path.Join(rootPath, "account")

	accKeyPath := path.Join(accountPath, "account.key")
	if _, err := os.Stat(accKeyPath); err == nil || os.IsExist(err) {
		return nil, fmt.Errorf("account key file is exist, please delete it. [%s]", accKeyPath)
	}

	accountFile := path.Join(accountPath, "account.json")
	if _, err := os.Stat(accountFile); err == nil || os.IsExist(err) {
		return nil, fmt.Errorf("account file is exist, please delete it. [%s]", accountFile)
	}

	privateKey, err := misc.GeneratePrivateKey(accKeyPath)
	if err != nil {
		return nil, fmt.Errorf("generate-privatekey[%s]: %s", accKeyPath, err.Error())
	}

	acc := &Account{
		Email: email,
		key:   privateKey,
		path:  accountPath,
	}

	cli, err := acme.NewClient(config.Config.AcmeURL, acc, acme.RSA2048)
	if err != nil {
		return nil, fmt.Errorf("init-client: %s", err.Error())
	}

	reg, err := cli.Register(true)
	if err != nil {
		return nil, fmt.Errorf("cli-register: %s", err.Error())
	}

	acc.Registration = reg
	return acc, acc.Save()
}

func (acc Account) GetEmail() string {
	return acc.Email
}

func (acc Account) GetRegistration() *acme.RegistrationResource {
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
