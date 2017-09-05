package options

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	EnvGoogleAuthJSON = "GOOGLE_AUTH_JSON"

	EnvGoogleApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"
	EnvAccount                      = "PRESIGNER_ACCOUNT"
	EnvBucket                       = "PRESIGNER_BUCKET"
	EnvDuration                     = "PRESIGNER_DURATION"
	EnvHost                         = "PRESIGNER_HOST"
	EnvPort                         = "PRESIGNER_PORT"
	EnvPrefix                       = "RRESIGNER_PREFIX"
	EnvVerbose                      = "PRESIGNER_VERBOSE"

	FlagAccount  = "account"
	FlagBucket   = "bucket"
	FlagDuration = "duration"
	FlagHost     = "host"
	FlagPort     = "port"
	FlagPrefix   = "prefix"
	FlagVerbose  = "verbose"
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
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  presigner [options]\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}
	fs.Var(&o.Account, FlagAccount, `Path to the file of Google service account JSON.`)
	fs.StringVar(&o.Bucket, FlagBucket, "", `Bucket name of Google Cloud Storage to upload files.`)
	fs.DurationVar(&o.Duration, FlagDuration, time.Minute, `Available duration of published signature.`)
	fs.Var(&o.Hosts, FlagHost, `Hosts of the image that is allowed to resize.
        When this value isn't specified, all hosts are allowed.
        Multiple hosts can be specified with:
          $ presigner -host a.com,b.com
          $ presigner -host a.com -host b.com`)
	fs.IntVar(&o.Port, FlagPort, 80, `Port to be listened.`)
	fs.StringVar(&o.Prefix, FlagPrefix, "", `Prefix of object name like 'uploads/'.`)
	fs.BoolVar(&o.Verbose, FlagVerbose, false, `Verbose output.
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
