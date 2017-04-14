package options

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type Account struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

func (a *Account) String() string {
	return fmt.Sprintf("%+v", *a)
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
	return json.Unmarshal(b, a)
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
