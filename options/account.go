package options

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type Account struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	ProjectID   string `json:"project_id"`
}

func (a *Account) String() string {
	return a.ClientEmail
}

func (a *Account) Set(path string) error {
	if path == "" {
		return errors.New("path to Google service account JSON isn't specified")
	}
	var b []byte
	var err error
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "fail to read the file of Google service account JSON")
	}
	if err := json.Unmarshal(b, a); err != nil {
		return errors.Wrap(err, "fail to unmarshal JSON")
	}
	return nil
}

func (a *Account) UnmarshalJSON(data []byte) error {
	type Alias Account
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*a = Account(alias)
	a.PrivateKey = strings.Replace(a.PrivateKey, `\n`, "\n", -1)
	return nil
}
