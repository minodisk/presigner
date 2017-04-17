package options

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Options struct {
	Buckets        Buckets
	Duration       time.Duration
	Port           int
	ServiceAccount Account
	Verbose        bool
}

func Parse(args []string) (Options, error) {
	o := Options{
		ServiceAccount: Account{},
		Buckets:        Buckets{},
	}

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
	fs.IntVar(&o.Port, "port", 80, `TCP address to listen on.
         `)
	fs.Var(&o.ServiceAccount, "account", `Path to the file of Google service account JSON.`)
	fs.BoolVar(&o.Verbose, "verbose", false, `Verbose output.
         `)

	return o, fs.Parse(args)
}
