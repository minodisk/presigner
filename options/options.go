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
	EnvPort                         = "PRESIGNER_PORT"
	EnvPrefix                       = "PRESIGNER_PREFIX"
	EnvVerbose                      = "PRESIGNER_VERBOSE"

	FlagAccount  = "account"
	FlagBucket   = "bucket"
	FlagDuration = "duration"
	FlagPort     = "port"
	FlagPrefix   = "prefix"
	FlagVerbose  = "verbose"
)

var (
	Envs = []string{
		EnvGoogleApplicationCredentials,
		EnvAccount,
		EnvBucket,
		EnvPort,
		EnvPrefix,
		EnvVerbose,
	}
	Flags = []string{
		FlagAccount,
		FlagAccount,
		FlagBucket,
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
	Port     int
	Prefix   string
	Verbose  bool
}

func (o *Options) Parse(args []string) error {
	// Setup flag set.
	fs := flag.NewFlagSet("presigner", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage:
  presigner [options]

Options:
`)
		fs.PrintDefaults()
	}
	fs.Var(&o.Account, FlagAccount, `Path to the file of Google service account JSON.`)
	fs.StringVar(&o.Bucket, FlagBucket, "", `Bucket name of Google Cloud Storage to upload files.`)
	fs.DurationVar(&o.Duration, FlagDuration, time.Minute, `Available duration of published signature.`)
	fs.IntVar(&o.Port, FlagPort, 80, `Port to be listened.`)
	fs.StringVar(&o.Prefix, FlagPrefix, "", `Prefix of object name like 'uploads/'.`)
	fs.BoolVar(&o.Verbose, FlagVerbose, false, `Verbose output.`)

	// Parse service account JSON in environment variable
	// if that is specified.
	if v := os.Getenv(EnvGoogleAuthJSON); v != "" {
		b := []byte(v)
		if err := json.Unmarshal(b, &o.Account); err != nil {
			return err
		}
	}

	// Set other environment variables to options.
	for _, env := range Envs {
		flag := EnvFlagMap[env]
		if v := os.Getenv(env); v != "" {
			if err := fs.Set(flag, v); err != nil {
				return err
			}
		}
	}

	// Overwrite options with command line flags.
	return fs.Parse(args)
}
