package options

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

const (
	EnvGoogleAuthJSON = "GOOGLE_AUTH_JSON"

	EnvGoogleApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"
	EnvAccount                      = "PRESIGNER_ACCOUNT"
	EnvBucket                       = "PRESIGNER_BUCKET"
	EnvHost                         = "PRESIGNER_HOST"
	EnvPort                         = "PRESIGNER_PORT"
	EnvPrefix                       = "RRESIGNER_PREFIX"
	EnvVerbose                      = "PRESIGNER_VERBOSE"

	FlagAccount = "account"
	FlagBucket  = "bucket"
	FlagHost    = "host"
	FlagPort    = "port"
	FlagPrefix  = "prefix"
	FlagVerbose = "verbose"
)

var (
	Envs = []string{
		EnvGoogleApplicationCredentials,
		EnvAccount,
		EnvBucket,
		EnvHost,
		EnvPort,
		EnvPrefix,
		EnvVerbose,
	}
	Flags = []string{
		FlagAccount,
		FlagAccount,
		FlagBucket,
		FlagHost,
		FlagPort,
		FlagPrefix,
		FlagVerbose,
	}
	EnvFlagMap = map[string]string{}
)

func init() {
	for i, env := range Envs {
		EnvFlagMap[env] = Flags[i]
	}
}

type Options struct {
	Account  Account
	Bucket   string
	Duration time.Duration
	Hosts    Hosts
	Port     int
	Prefix   string
	Verbose  bool
}

func (o *Options) Parse(args []string) error {
	fs := flag.NewFlagSet("presigner", flag.ContinueOnError)
	fs.Var(&o.Account, "account", `Path to the file of Google service account JSON.`)
	fs.StringVar(&o.Bucket, "bucket", "", `Bucket name of Google Cloud Storage to upload files.`)
	fs.DurationVar(&o.Duration, "duration", time.Minute, `Available duration of published signature.
         `)
	fs.IntVar(&o.Port, "port", 80, `TCP address to listen on.
         `)
	fs.StringVar(&o.Prefix, "prefix", "", `Prefix of object`)
	fs.BoolVar(&o.Verbose, "verbose", false, `Verbose output.
         `)

	if v := os.Getenv(EnvGoogleAuthJSON); v != "" {
		b := []byte(v)
		if err := json.Unmarshal(b, &o.Account); err != nil {
			return err
		}
		// o.Account.Path = filepath.Join(os.TempDir(), "presigner-google-auth.json")
		// if err := ioutil.WriteFile(o.Account.Path, b, 0644); err != nil {
		// 	return err
		// }
	}
	for _, env := range Envs {
		flag := EnvFlagMap[env]
		if v := os.Getenv(env); v != "" {
			if err := fs.Set(flag, v); err != nil {
				return err
			}
		}
	}
	return fs.Parse(args)
}
