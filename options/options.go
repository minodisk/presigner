package options

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Options struct {
	Buckets         Buckets
	Duration        time.Duration
	GoogleAuthEmail string
	GoogleAuthKey   string
	Port            int
}

func New(args []string) (Options, error) {
	o := Options{Buckets: Buckets{}}
	fs := flag.NewFlagSet("presigner", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  presigner [options]\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}
	fs.Var(&o.Buckets, "bucket", `Allowed buckets to publish pre-signed URL.
         When this flag is empty, allows any buckets to publish.
         You can set multi bucket with:
            $ presigner -bucket foo -bucket bar`)
	fs.DurationVar(&o.Duration, "duration", time.Minute, `Available duration of published signature.
         `)
	fs.StringVar(&o.GoogleAuthEmail, "email", "", `Google service account client email address
         from the Google Developers Console in the form of
         "xxx@developer.gserviceaccount.com".`)
	fs.StringVar(&o.GoogleAuthKey, "key", "", `Google service account private key
         generated from P12 file with:
            $ openssl pkcs12 -in key.p12 -passin pass:notasecret -out key.pem -nodes`)
	var keyPath string
	fs.StringVar(&keyPath, "keypath", "", `Path to Google service account private key.
         When -key isn't specified, load -keypath file.`)
	fs.IntVar(&o.Port, "port", 80, `Listening port.
         `)
	if err := fs.Parse(args); err != nil {
		return o, err
	}

	return o.InitializeGoogleAuthKey(keyPath)
}

func (o Options) InitializeGoogleAuthKey(keyPath string) (Options, error) {
	if o.GoogleAuthKey != "" {
		o.GoogleAuthKey = strings.Replace(o.GoogleAuthKey, `\n`, "\n", -1)
		return o, nil
	}
	if keyPath != "" {
		key, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return o, errors.Wrap(err, "fail to read key path")
		}
		o.GoogleAuthKey = string(key)
		return o, nil
	}
	return o, nil
}
