package options

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Options struct {
	Account  Account
	Buckets  Buckets
	Duration time.Duration
	Port     int
	Verbose  bool
}

type Account struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
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

func Parse(args []string) (Options, error) {
	o := Options{
		Account: Account{},
		Buckets: Buckets{},
	}

	fs := flag.NewFlagSet("presigner", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  presigner [options]\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}
	var j string
	fs.StringVar(&j, "account", "", `Google service account JSON.`)
	var f string
	fs.StringVar(&f, "accountfile", "", `Path to Google service account JSON file.
         When -account isn't specified, load file at -accountfile.`)
	fs.Var(&o.Buckets, "bucket", `Allowed buckets to publish pre-signed URL.
         When this flag is empty, allows any buckets to publish.
         You can set multi bucket with:
            $ presigner -bucket foo -bucket bar`)
	fs.DurationVar(&o.Duration, "duration", time.Minute, `Available duration of published signature.
         `)
	fs.IntVar(&o.Port, "port", 80, `Listening port.
         `)
	fs.BoolVar(&o.Verbose, "verbose", false, `Verbose
         `)
	if err := fs.Parse(args); err != nil {
		return o, err
	}

	var b []byte
	if j != "" {
		b = []byte(j)
	} else if f != "" {
		var err error
		b, err = ioutil.ReadFile(f)
		if err != nil {
			return o, errors.Wrap(err, "fail to read the file of Google service account JSON")
		}
	} else {
		return o, errors.New("Google service account JSON isn't specified")
	}
	return o.FillAccountWithJSON(b)
}

func (o Options) FillAccountWithJSON(b []byte) (Options, error) {
	return o, json.Unmarshal(b, &o.Account)
}
